package sjson

import (
	"bytes"
	"encoding/json"
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
	regexp.MustCompile(`(:)"([\[\{])([\\"\"\{\[\d\-])`),      // { [
	regexp.MustCompile(`([\"\]\}\d\-][\]\}])\"([,\,\}]*)`), // }]
	regexp.MustCompile(`(:)"([\[\{])([\]\}])"`),      // [] {}
	regexp.MustCompile(`([,:\[ \{])\\(")`),                 // ,\"
	regexp.MustCompile(`\\(")([:,\]\}])`),                  // \",
	regexp.MustCompile(`\\(\\)`),                           // \\
}

func (j *Json) ReplaceAllString(regexps []*regexp.Regexp, src string) string {
	substitution := `$1$2$3`
	for _, v := range regexps {
		src = v.ReplaceAllString(src, substitution)
	}
	return src
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
	j := &Json{
		V: v,
	}
	return j.MustToJsonByte()
}

func ToJsonString(v interface{}) string {
	j := &Json{
		V: v,
	}
	return j.MustToJsonString()
}

func StringWithJsonToString(v interface{}) string {
	j := &Json{
		V: v,
	}
	return j.StringWithJsonMustToString()
}

func StringWithJsonRegexToString(v interface{}) string {
	j := &Json{
		V: v,
	}
	return j.StringWithJsonMustRegexToString()
}
