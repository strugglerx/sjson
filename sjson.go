package sjson

import (
	"bytes"
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"sync"
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

// ===================== JSON 编码 =====================

func (j *Json) MustToJsonByte() []byte {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.Encode(j.V)
	b := buf.Bytes()
	if len(b) > 0 && b[len(b)-1] == '\n' {
		b = b[:len(b)-1]
	}
	res := make([]byte, len(b))
	copy(res, b)
	return res
}

func (j *Json) MustToJsonString() string {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.Encode(j.V)
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
	// 直接处理字节，减少 string 拷贝开销
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.Encode(j.V)
	src := buf.Bytes()
	if len(src) > 0 && src[len(src)-1] == '\n' {
		src = src[:len(src)-1]
	}

	sb := sbPool.Get().(*strings.Builder)
	sb.Reset()
	sb.Grow(len(src))
	defer sbPool.Put(sb)

	scratch := scratchPool.Get().([]byte)
	scratch = scratch[:0]
	defer func() {
		if cap(scratch) < 4096 { // 限制复用大小，防止极端长数据撑大内存
			scratchPool.Put(scratch)
		}
	}()

	i, n := 0, len(src)
	for i < n {
		if src[i] == ':' && i+2 < n && src[i+1] == '"' {
			next := src[i+2]
			if next == '{' || next == '[' {
				if end := findStringEndBytes(src, i+2); end > 0 {
					// 复用 scratch 缓冲区进行脱义和验证
					scratch = unescapeIntoScratch(scratch[:0], src[i+2:end])
					if json.Valid(scratch) {
						sb.WriteByte(':')
						sb.Write(scratch)
						i = end + 1
						continue
					}
				}
			}
		}
		sb.WriteByte(src[i])
		i++
	}
	return sb.String()
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

func ToJsonString(v interface{}) string {
	return New(v).MustToJsonString()
}

func StringWithJsonScanToString(v interface{}) string {
	return New(v).StringWithJsonMustScanToString()
}
