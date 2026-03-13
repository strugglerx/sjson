package sjson

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
)

const (
	maxPooledBufferCap  = 64 << 10
	maxPooledBuilderCap = 64 << 10
	maxPooledScratchCap = 4 << 10
)

// bufferPool 复用用于 JSON 序列化的缓冲
var bufferPool = sync.Pool{
	New: func() interface{} { return new(bytes.Buffer) },
}

// sbPool 复用用于构建最终结果的 Builder
var sbPool = sync.Pool{
	New: func() interface{} { return new(strings.Builder) },
}

// scratchPool 复用临时字节切片，用于存储脱义后的内容进行验证
var scratchPool = sync.Pool{
	New: func() interface{} { return make([]byte, 0, 1024) },
}

type Json struct {
	V interface{}
}

func New(v interface{}) *Json {
	return &Json{V: v}
}

func putBuffer(buf *bytes.Buffer) {
	if buf.Cap() <= maxPooledBufferCap {
		bufferPool.Put(buf)
	}
}

func putBuilder(sb *strings.Builder) {
	if sb.Cap() <= maxPooledBuilderCap {
		sbPool.Put(sb)
	}
}

func putScratch(scratch []byte) {
	if cap(scratch) <= maxPooledScratchCap {
		scratchPool.Put(scratch)
	}
}

func trimTrailingNewline(src []byte) []byte {
	if len(src) > 0 && src[len(src)-1] == '\n' {
		return src[:len(src)-1]
	}
	return src
}

func encodeIntoBuffer(v interface{}) (*bytes.Buffer, error) {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()

	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		putBuffer(buf)
		return nil, err
	}
	return buf, nil
}

func scanEncodedJSONBytes(src []byte) []byte {
	scratch := scratchPool.Get().([]byte)
	scratch = scratch[:0]
	defer func() {
		putScratch(scratch)
	}()

	out := make([]byte, 0, len(src))
	i, n := 0, len(src)
	segmentStart := 0
	for i < n {
		if src[i] == ':' && i+2 < n && src[i+1] == '"' {
			next := src[i+2]
			if next == '{' || next == '[' {
				if end := findStringEndBytes(src, i+2); end > 0 {
					if expanded, ok := expandCandidateJSON(src[i+2:end], scratch[:0]); ok {
						out = append(out, src[segmentStart:i]...)
						out = append(out, ':')
						out = append(out, expanded...)
						i = end + 1
						segmentStart = i
						continue
					}
				}
			}
		}
		i++
	}
	out = append(out, src[segmentStart:]...)
	return out
}

func scanEncodedJSON(src []byte) string {
	sb := sbPool.Get().(*strings.Builder)
	sb.Reset()
	sb.Grow(len(src))
	defer putBuilder(sb)

	scratch := scratchPool.Get().([]byte)
	scratch = scratch[:0]
	defer func() {
		putScratch(scratch)
	}()

	i, n := 0, len(src)
	segmentStart := 0
	for i < n {
		if src[i] == ':' && i+2 < n && src[i+1] == '"' {
			next := src[i+2]
			if next == '{' || next == '[' {
				if end := findStringEndBytes(src, i+2); end > 0 {
					if expanded, ok := expandCandidateJSON(src[i+2:end], scratch[:0]); ok {
						sb.Write(src[segmentStart:i])
						sb.WriteByte(':')
						sb.Write(expanded)
						i = end + 1
						segmentStart = i
						continue
					}
				}
			}
		}
		i++
	}
	sb.Write(src[segmentStart:])
	return sb.String()
}

func expandCandidateJSON(src, scratch []byte) ([]byte, bool) {
	if !looksLikeJSONContainer(src) {
		return scratch, false
	}

	if bytes.IndexByte(src, '\\') < 0 {
		return src, json.Valid(src)
	}

	scratch = unescapeIntoScratch(scratch, src)
	if !looksLikeJSONContainer(scratch) {
		return scratch, false
	}
	return scratch, json.Valid(scratch)
}

func looksLikeJSONContainer(src []byte) bool {
	if len(src) < 2 {
		return false
	}

	first := src[0]
	last := src[len(src)-1]
	return (first == '{' && last == '}') || (first == '[' && last == ']')
}

func encodeJSON(v interface{}) ([]byte, error) {
	buf, err := encodeIntoBuffer(v)
	if err != nil {
		return nil, err
	}
	defer putBuffer(buf)

	b := buf.Bytes()
	if len(b) > 0 && b[len(b)-1] == '\n' {
		b = b[:len(b)-1]
	}

	out := make([]byte, len(b))
	copy(out, b)
	return out, nil
}

// ===================== JSON 编码 =====================

func (j *Json) MustToJsonByte() []byte {
	b, err := encodeJSON(j.V)
	if err != nil {
		panic(fmt.Sprintf("sjson: encode json: %v", err))
	}
	return b
}

func (j *Json) MustToJsonString() string {
	buf, err := encodeIntoBuffer(j.V)
	if err != nil {
		panic(fmt.Sprintf("sjson: encode json: %v", err))
	}
	defer putBuffer(buf)

	return strings.TrimSuffix(buf.String(), "\n")
}

func (j *Json) Struct(v interface{}) error {
	switch vv := j.V.(type) {
	case string:
		return json.Unmarshal([]byte(vv), &v)
	case []byte:
		return json.Unmarshal(vv, &v)
	default:
		return errors.New("type error")
	}
}

// ===================== 方案一：正则版（Safety Regex）=====================

var regexpsSafety = []*regexp.Regexp{
	regexp.MustCompile(`\\(\"[\w]+)\\(\"\:)`),               //key
	regexp.MustCompile(`(:)"([\[\{][\d\{\["\\].*?[\]\}])"`), //{} []
	regexp.MustCompile(`(:)"([\[\{][\]\}])"`),               //{} []
	regexp.MustCompile(`\\(")`),                             // \"
	regexp.MustCompile(`\\(\\)`),                            // \\
}

func (j *Json) ReplaceAllString(regexps []*regexp.Regexp, src string) string {
	substitution := `$1$2$3`
	for _, v := range regexps {
		src = v.ReplaceAllString(src, substitution)
	}
	return src
}

func (j *Json) SearchStringWithJsons(src string) []string {
	reg := regexp.MustCompile(`:(\"[\{\[].*?[\}\]]\")`)
	return reg.FindAllString(src, -1)
}

func (j *Json) RemoveRepeatedElement(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return newArr
}

func (j *Json) StringWithJsonSafetyMustRegexToString() string {
	src := j.MustToJsonString()
	preElement := j.RemoveRepeatedElement(j.SearchStringWithJsons(src))
	for _, v := range preElement {
		src = strings.ReplaceAll(src, v, j.ReplaceAllString(regexpsSafety, v))
	}
	return src
}

func StringWithJsonSafetyRegexToString(v interface{}) string {
	return New(v).StringWithJsonSafetyMustRegexToString()
}

// ===================== 方案二：极致扫描版（单次 O(n)）=====================

func (j *Json) StringWithJsonMustScanToString() string {
	buf, err := encodeIntoBuffer(j.V)
	if err != nil {
		panic(fmt.Sprintf("sjson: encode json: %v", err))
	}
	defer putBuffer(buf)

	return scanEncodedJSON(trimTrailingNewline(buf.Bytes()))
}

func findStringEndBytes(s []byte, pos int) int {
	for i := pos; i < len(s); i++ {
		if s[i] == '\\' {
			i++
			continue
		}
		if s[i] == '"' {
			return i
		}
	}
	return -1
}

// unescapeIntoScratch 使用批量处理逻辑：寻找反斜杠，批量拷贝中间段
func unescapeIntoScratch(dst, src []byte) []byte {
	idx := bytes.IndexByte(src, '\\')
	if idx < 0 {
		return append(dst, src...)
	}

	for {
		dst = append(dst, src[:idx]...)
		src = src[idx:]
		if len(src) < 2 {
			dst = append(dst, src...)
			break
		}

		// 处理转义
		switch src[1] {
		case '"':
			dst = append(dst, '"')
			src = src[2:]
		case '\\':
			dst = append(dst, '\\')
			src = src[2:]
		default:
			dst = append(dst, src[0], src[1])
			src = src[2:]
		}

		idx = bytes.IndexByte(src, '\\')
		if idx < 0 {
			dst = append(dst, src...)
			break
		}
	}
	return dst
}

// ===================== 公开函数 =====================

func ToJsonByte(v interface{}) []byte {
	return New(v).MustToJsonByte()
}

func ToJsonByteE(v interface{}) ([]byte, error) {
	return encodeJSON(v)
}

func ToJsonString(v interface{}) string {
	return New(v).MustToJsonString()
}

func ToJsonStringE(v interface{}) (string, error) {
	buf, err := encodeIntoBuffer(v)
	if err != nil {
		return "", err
	}
	defer putBuffer(buf)

	return strings.TrimSuffix(buf.String(), "\n"), nil
}

func StringWithJsonScanToString(v interface{}) string {
	return New(v).StringWithJsonMustScanToString()
}

func StringWithJsonScanToBytes(v interface{}) []byte {
	buf, err := encodeIntoBuffer(v)
	if err != nil {
		panic(fmt.Sprintf("sjson: encode json: %v", err))
	}
	defer putBuffer(buf)

	return scanEncodedJSONBytes(trimTrailingNewline(buf.Bytes()))
}

func StringWithJsonScanToStringE(v interface{}) (string, error) {
	buf, err := encodeIntoBuffer(v)
	if err != nil {
		return "", err
	}
	defer putBuffer(buf)

	return scanEncodedJSON(trimTrailingNewline(buf.Bytes())), nil
}

func StringWithJsonScanToBytesE(v interface{}) ([]byte, error) {
	buf, err := encodeIntoBuffer(v)
	if err != nil {
		return nil, err
	}
	defer putBuffer(buf)

	return scanEncodedJSONBytes(trimTrailingNewline(buf.Bytes())), nil
}
