package sjson

import (
	"bytes"
	"encoding/json"
	"errors"
	"regexp"
	"strings"
)

type Json struct {
	V interface{}
}

func New(v interface{}) *Json {
	return &Json{
		V: v,
	}
}

var replace = []string{`\\`, `\`, `\"`, `"`, `"[`, "[", `"{`, "{", `]",`, "],", `]"}`, "]}", `}",`, "},", `}"}`, "}}"}

var regexps = []*regexp.Regexp{
	regexp.MustCompile(`(:)"([\[\{])([\\"\"\{\[\d\-])`),    // { [
	regexp.MustCompile(`([\"\]\}\d\-][\]\}])\"([,\,\}]*)`), // }]
	regexp.MustCompile(`(:)"([\[\{])([\]\}])"`),            // [] {}
	regexp.MustCompile(`([,:\[ \{])\\(")`),                 // ,\"
	regexp.MustCompile(`\\(")([:,\]\}])`),                  // \",
	regexp.MustCompile(`\\(\\)`),                           // \\
}

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

func (j *Json) MustToJsonByte() []byte {
	buffer := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(buffer)
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.Encode(j.V)
	return buffer.Bytes()
}

func (j *Json) MustToJsonString() string {
	buffer := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(buffer)
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.Encode(j.V)
	return buffer.String()
}

func (j *Json) Struct(v interface{}) error {
	var err error
	switch vv := j.V.(type) {
	case string:
		err = json.Unmarshal([]byte(vv), &v)
	case []byte:
		err = json.Unmarshal(vv, &v)
	default:
		err = errors.New("type error")
	}
	return err
}

//`(:)"([\[\{])([\\"\"\{\d\-]) => $1$2$3 `       // { [
//
//`([\"\]\}\d\-][\]\}])\"([,\,\}]*)` => $1$2$3   // }]
//
//`([,:\[ \{])\\(")`                             //  [] {}
//`([,:\[ \{])\\(")` => $1$2$3                   // ,\"
//
//`\\(")([:,\]\}])`  => $1$2$3                   // \",
//
//`\\(\\)` => $1$2$3                             // \\
func (j *Json) StringWithJsonMustRegexToString() string {
	return j.ReplaceAllString(regexps, j.MustToJsonString())
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
	return
}

func (j *Json) StringWithJsonSafetyMustRegexToString() string {
	src := j.MustToJsonString()
	preElement := j.RemoveRepeatedElement(j.SearchStringWithJsons(src))
	for _, v := range preElement {
		src = strings.ReplaceAll(src, v, j.ReplaceAllString(regexpsSafety, v))
	}
	return src
}

//Mysql 如果含有json的数据 可以纠正转义符号转换成正常的结构
//
//几种情况:
//正常的字段没有转义符号但是json结构里有 \" => "
//
//json如果存的array 因为是字符串 "[ => [
//
//json如果存的array 因为是字符串且]在中间 ]", => ],
//
//json如果存的array 因为是字符串且]在最后 ]"} => ]},
//
//json如果存的object 因为是字符串 "{ => {
//
//json如果存的object 因为是字符串且}在中间 }", => },
//
//json如果存的object 因为是字符串且}在最后 }"} => }},
//
//json如果存的object 因为是字符串 \\ => \,
func (j *Json) StringWithJsonMustToString() string {
	return strings.NewReplacer(replace...).Replace(j.MustToJsonString())
}

func ToJsonByte(v interface{}) []byte {
	return New(v).MustToJsonByte()
}

func ToJsonString(v interface{}) string {
	return New(v).MustToJsonString()
}

func StringWithJsonToString(v interface{}) string {
	return New(v).StringWithJsonMustToString()
}

func StringWithJsonRegexToString(v interface{}) string {
	return New(v).StringWithJsonMustRegexToString()
}

func StringWithJsonSafetyRegexToString(v interface{}) string {
	return New(v).StringWithJsonSafetyMustRegexToString()
}
