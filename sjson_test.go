package sjson

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

var check = map[string]interface{}{
	"a": "string",
	"b": 0,
	"c": 1.23,
	"d": "[天呐]",
	"e": "我叫\u201c王二蛋\"个子不高本事不小",
	"f": "请认真完善相关信息(）单位为\"m\",幕墙面积等，完善，信息。",
	"g": "请认真完善相关信息(）单位为,\"m\",幕墙面积等，完善，信息。",
}

var check1 = map[string]interface{}{
	"a": "{\"key\":[\"中文\", \"english\", \"dog\", \"man\"]}",
	"b": "{\"key\":[\"中文\", \"english\", \"dog\", \"man\"],\"key1\":[-1,2,3,4,5],\"key2\":[1,2,3,4,5]}",
	"c": "[1,\"中文\",2,3, \"english\", \"dog\", \"man\"]",
	"d": "[{\"url\":\"请认真完善相关信息(）单{位}为\\\"m\\\",幕墙面积等，完善，信息。\",\"list\":[\"https://xxxxxxxx.com.cn/pic_2323.1-5-2png\",\"https://xxxxxxxx.com.cn/pic_2323.1-5-2png\",\"https://xxxxxxxx.com.cn/pic_2323.1-5-2png\"],\"desc\":\"换行\\n换行\"},[\"中文\",\"english\",\"dog\",\"man\"],[-1,2,3,4,5],[1,2,3,4,5]]",
	"e": "[]",
	"f": "{}",
	"g": "[[\"中文\", \"english\", \"dog\", \"man\"],[\"中文\", \"english\", \"dog\", \"man\"]]",
	"h": "[[1,2,3,4,14,5.3],[1,2,3,4,14,5.3]]",
	"i": "{\"url\":{\"url\":\"https://xxxxxxxx.com.cn/pic_2323.1-5-2png\", \"desc\": \"换行\\n换行\"}}",
	"j": "[\"https://xxxxxxxx.com.cn/pic_2323.1-5-2png\",\"https://xxxxxxxx.com.cn/pic_2323.1-5-2png\",\"https://xxxxxxxx.com.cn/pic_2323.1-5-2png\"]",
	"k": map[string]string{
		"a": "string",
		"c": "{天呐}",
		"d": "[天呐]",
		"e": "我叫\u201c王二蛋\"个子不高本事不小",
		"f": "请认真完善相关信息(）单位为\"m\",幕墙面积等，完善，信息。",
		"g": "请认真完善相关信息(）单位为,\"m\",幕墙面积等，完善，信息。",
	},
	"l": "[{\"gg\":\"我叫\u201c王二蛋个\\\"子不高本事不小\",\"url\":\"请认真完善相关信息(）单{位}为\\\"m\\\",幕墙面积等，完善，信息。\",\"list\":[\"https://xxxxxxxx.com.cn/pic_2323.1-5-2png\",\"https://xxxxxxxx.com.cn/pic_2323.1-5-2png\",\"https://xxxxxxxx.com.cn/pic_2323.1-5-2png\"],\"desc\":\"换行\\n换行\"},[\"中文\",\"english\",\"dog\",\"man\"],[-1,2,3,4,5],[1,2,3,4,5]]",
}

// assertValidJSON 验证字符串是否为合法的 JSON，不合法则 t.Error
func assertValidJSON(t *testing.T, label, s string) {
	t.Helper()
	var v interface{}
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		t.Errorf("[%s] ❌ 输出 JSON 不合法: %v\n内容:\n%s", label, err, s)
	} else {
		t.Logf("[%s] ✅ JSON 合法", label)
	}
}

func assertJSONEqual(t *testing.T, label, got, want string) {
	t.Helper()

	var gotV interface{}
	if err := json.Unmarshal([]byte(got), &gotV); err != nil {
		t.Fatalf("[%s] got is invalid JSON: %v\n%s", label, err, got)
	}

	var wantV interface{}
	if err := json.Unmarshal([]byte(want), &wantV); err != nil {
		t.Fatalf("[%s] expected value is invalid JSON: %v\n%s", label, err, want)
	}

	if !reflect.DeepEqual(gotV, wantV) {
		t.Fatalf("[%s] JSON mismatch.\nGot: %s\nExp: %s", label, got, want)
	}
}

// buildLargeList 构建大数据集（用于 benchmark）
func buildLargeList() []map[string]interface{} {
	list := make([]map[string]interface{}, 0, 200)
	for k := 0; k < 100; k++ {
		middle := map[string]interface{}{
			"a": "{\"key\":[\"中文\", \"english\", \"dog\", \"man\"]}",
			"b": "{\"key\":[\"中文\", \"english\", \"dog\", \"man\"],\"key1\":[-1,2,3,4,5],\"key2\":[1,2,3,4,5]}",
			"i": fmt.Sprintf("{\"url\":{\"url\":\"https://xxxxxxxx.com.cn/pic-%d.png\", \"desc\": \"换行\\n换行\"}}", k),
		}
		list = append(list, check1)
		list = append(list, middle)
	}
	return list
}

// ===================== 功能测试 =====================

// TestSafetyRegex 旧版：正则实现，测试输出是否合法 JSON
func TestSafetyRegex(t *testing.T) {
	eg1 := StringWithJsonSafetyRegexToString(check)
	t.Log("eg1:", eg1)
	assertValidJSON(t, "Regex/check", eg1)

	eg2 := StringWithJsonSafetyRegexToString(check1)
	t.Log("eg2:", eg2)
	assertValidJSON(t, "Regex/check1", eg2)
}

// TestScanner 新版：单次扫描，测试输出是否合法 JSON
func TestScanner(t *testing.T) {
	eg1 := StringWithJsonScanToString(check)
	t.Log("eg1:", eg1)
	assertValidJSON(t, "Scanner/check", eg1)

	eg2 := StringWithJsonScanToString(check1)
	t.Log("eg2:", eg2)
	assertValidJSON(t, "Scanner/check1", eg2)

	//ioutil.WriteFile("test.json", []byte(StringWithJsonScanToString(buildLargeList())), 0666)
}

func TestJson_Struct(t *testing.T) {
	v := "[1,2,3,4,5,6,7,8]"
	vs := make([]int, 0)
	err := New(v).Struct(&vs)
	t.Log(err, vs)
}

func TestScannerPreservesPlainStrings(t *testing.T) {
	input := map[string]interface{}{
		"plainArrayLike":  "[天呐]",
		"plainObjectLike": "{not-json}",
		"realNested":      "{\"ok\":true,\"list\":[1,2,3]}",
	}

	got := StringWithJsonScanToString(input)

	var decoded map[string]interface{}
	if err := json.Unmarshal([]byte(got), &decoded); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}

	if decoded["plainArrayLike"] != "[天呐]" {
		t.Fatalf("plainArrayLike changed unexpectedly: %#v", decoded["plainArrayLike"])
	}
	if decoded["plainObjectLike"] != "{not-json}" {
		t.Fatalf("plainObjectLike changed unexpectedly: %#v", decoded["plainObjectLike"])
	}

	realNested, ok := decoded["realNested"].(map[string]interface{})
	if !ok {
		t.Fatalf("realNested was not expanded into object: %#v", decoded["realNested"])
	}
	if realNested["ok"] != true {
		t.Fatalf("unexpected realNested content: %#v", realNested)
	}
}

func TestErrorAPIsAndMustPanic(t *testing.T) {
	unsupported := map[string]interface{}{
		"bad": func() {},
	}

	if _, err := ToJsonByteE(unsupported); err == nil {
		t.Fatal("ToJsonByteE should return error for unsupported type")
	}
	if _, err := ToJsonStringE(unsupported); err == nil {
		t.Fatal("ToJsonStringE should return error for unsupported type")
	}
	if _, err := StringWithJsonScanToStringE(unsupported); err == nil {
		t.Fatal("StringWithJsonScanToStringE should return error for unsupported type")
	}
	if _, err := StringWithJsonScanToBytesE(unsupported); err == nil {
		t.Fatal("StringWithJsonScanToBytesE should return error for unsupported type")
	}

	assertPanicContains(t, "ToJsonByte", "unsupported type", func() {
		ToJsonByte(unsupported)
	})
	assertPanicContains(t, "ToJsonString", "unsupported type", func() {
		ToJsonString(unsupported)
	})
	assertPanicContains(t, "StringWithJsonScanToString", "unsupported type", func() {
		StringWithJsonScanToString(unsupported)
	})
}

func assertPanicContains(t *testing.T, label, want string, fn func()) {
	t.Helper()
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("%s should panic", label)
		}
		if !strings.Contains(fmt.Sprint(r), want) {
			t.Fatalf("%s panic mismatch: got %v, want substring %q", label, r, want)
		}
	}()
	fn()
}

func TestScannerComplexCases(t *testing.T) {
	cases := []struct {
		name   string
		input  map[string]interface{}
		verify func(t *testing.T, decoded map[string]interface{})
	}{
		{
			name: "escaped nested object with quotes and slashes",
			input: map[string]interface{}{
				"payload": "{\"meta\":{\"title\":\"A \\\"quote\\\" here\",\"path\":\"C:\\\\temp\\\\demo\"},\"ok\":true}",
			},
			verify: func(t *testing.T, decoded map[string]interface{}) {
				t.Helper()
				payload, ok := decoded["payload"].(map[string]interface{})
				if !ok {
					t.Fatalf("payload not expanded: %#v", decoded["payload"])
				}
				meta, ok := payload["meta"].(map[string]interface{})
				if !ok {
					t.Fatalf("meta missing: %#v", payload["meta"])
				}
				if meta["title"] != `A "quote" here` {
					t.Fatalf("unexpected title: %#v", meta["title"])
				}
				if meta["path"] != `C:\temp\demo` {
					t.Fatalf("unexpected path: %#v", meta["path"])
				}
			},
		},
		{
			name: "nested array with mixed objects",
			input: map[string]interface{}{
				"payload": "[{\"name\":\"alice\",\"attrs\":{\"level\":3}},{\"name\":\"bob\",\"tags\":[\"x\",\"y\"]}]",
			},
			verify: func(t *testing.T, decoded map[string]interface{}) {
				t.Helper()
				payload, ok := decoded["payload"].([]interface{})
				if !ok || len(payload) != 2 {
					t.Fatalf("payload not expanded into array: %#v", decoded["payload"])
				}
				first, ok := payload[0].(map[string]interface{})
				if !ok || first["name"] != "alice" {
					t.Fatalf("unexpected first element: %#v", payload[0])
				}
				second, ok := payload[1].(map[string]interface{})
				if !ok || second["name"] != "bob" {
					t.Fatalf("unexpected second element: %#v", payload[1])
				}
			},
		},
		{
			name: "plain strings that look close to json remain strings",
			input: map[string]interface{}{
				"bad1": "{missing",
				"bad2": "[1,2,3",
				"bad3": "{key:value}",
				"bad4": `[not "real" json]`,
			},
			verify: func(t *testing.T, decoded map[string]interface{}) {
				t.Helper()
				for key, want := range map[string]interface{}{
					"bad1": "{missing",
					"bad2": "[1,2,3",
					"bad3": "{key:value}",
					"bad4": `[not "real" json]`,
				} {
					if decoded[key] != want {
						t.Fatalf("%s changed unexpectedly: got %#v want %#v", key, decoded[key], want)
					}
				}
			},
		},
		{
			name: "empty containers should expand",
			input: map[string]interface{}{
				"emptyArray":  "[]",
				"emptyObject": "{}",
			},
			verify: func(t *testing.T, decoded map[string]interface{}) {
				t.Helper()
				if got, ok := decoded["emptyArray"].([]interface{}); !ok || len(got) != 0 {
					t.Fatalf("emptyArray mismatch: %#v", decoded["emptyArray"])
				}
				if got, ok := decoded["emptyObject"].(map[string]interface{}); !ok || len(got) != 0 {
					t.Fatalf("emptyObject mismatch: %#v", decoded["emptyObject"])
				}
			},
		},
		{
			name: "unicode and escaped newline stay correct",
			input: map[string]interface{}{
				"payload": "{\"message\":\"你好\\n世界\",\"items\":[\"雪\",\"山\",\"海\"]}",
			},
			verify: func(t *testing.T, decoded map[string]interface{}) {
				t.Helper()
				payload, ok := decoded["payload"].(map[string]interface{})
				if !ok {
					t.Fatalf("payload not expanded: %#v", decoded["payload"])
				}
				if payload["message"] != "你好\n世界" {
					t.Fatalf("unexpected message: %#v", payload["message"])
				}
				items, ok := payload["items"].([]interface{})
				if !ok || len(items) != 3 {
					t.Fatalf("unexpected items: %#v", payload["items"])
				}
			},
		},
		{
			name: "double encoded json only expands one level",
			input: map[string]interface{}{
				"payload": "\"{\\\"deep\\\":true}\"",
			},
			verify: func(t *testing.T, decoded map[string]interface{}) {
				t.Helper()
				if decoded["payload"] != `"{"deep":true}"` && decoded["payload"] != `"{\"deep\":true}"` {
					t.Fatalf("unexpected payload: %#v", decoded["payload"])
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out := StringWithJsonScanToString(tc.input)
			assertValidJSON(t, tc.name, out)

			var decoded map[string]interface{}
			if err := json.Unmarshal([]byte(out), &decoded); err != nil {
				t.Fatalf("unmarshal result: %v", err)
			}
			tc.verify(t, decoded)
		})
	}
}

func FuzzStringWithJsonScanToString(f *testing.F) {
	seeds := []string{
		``,
		`[]`,
		`{}`,
		`[天呐]`,
		`{"a":1}`,
		`[1,2,3]`,
		`{\"a\":1}`,
		`{"meta":{"title":"A \"quote\""}}`,
		`"[{\"nested\":true}]"`,
		`{missing`,
		`[1,2,3`,
		`{key:value}`,
		`你好\n世界`,
		`C:\\temp\\demo`,
	}
	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		data := map[string]interface{}{
			"payload": input,
			"mirror":  input + "_suffix",
		}

		out := StringWithJsonScanToString(data)

		var decoded map[string]interface{}
		if err := json.Unmarshal([]byte(out), &decoded); err != nil {
			t.Fatalf("invalid json: %v\ninput=%q\nout=%s", err, input, out)
		}

		wantMirror := jsonRoundTripString(t, input+"_suffix")
		if decoded["mirror"] != wantMirror {
			t.Fatalf("mirror changed unexpectedly: input=%q got=%#v want=%#v", input, decoded["mirror"], wantMirror)
		}

		wantPayload := expectedPayloadAfterScan(t, input)
		if !reflect.DeepEqual(decoded["payload"], wantPayload) {
			t.Fatalf("payload changed unexpectedly: input=%q got=%#v want=%#v", input, decoded["payload"], wantPayload)
		}
	})
}

func jsonRoundTripString(t *testing.T, s string) string {
	t.Helper()

	b, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("marshal string: %v", err)
	}

	var out string
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal string: %v", err)
	}
	return out
}

func expectedPayloadAfterScan(t *testing.T, input string) interface{} {
	t.Helper()

	normalized := jsonRoundTripString(t, input)
	buf, err := encodeIntoBuffer(input)
	if err != nil {
		t.Fatalf("encode input: %v", err)
	}
	defer putBuffer(buf)

	encoded := trimTrailingNewline(buf.Bytes())
	if len(encoded) < 2 {
		return normalized
	}

	candidate := encoded[1 : len(encoded)-1]
	expanded, ok := expandCandidateJSON(candidate, nil)
	if !ok {
		return normalized
	}

	var v interface{}
	if err := json.Unmarshal(expanded, &v); err != nil {
		return normalized
	}

	switch v.(type) {
	case map[string]interface{}, []interface{}:
		return v
	default:
		return normalized
	}
}

// ===================== Benchmark：正则 vs 单次扫描 =====================

func BenchmarkSafetyRegex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		StringWithJsonSafetyRegexToString(check1)
	}
}

func BenchmarkScanner(b *testing.B) {
	for i := 0; i < b.N; i++ {
		StringWithJsonScanToString(check1)
	}
}

func BenchmarkScannerBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		StringWithJsonScanToBytes(check1)
	}
}

func BenchmarkSafetyRegex_Large(b *testing.B) {
	list := buildLargeList()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StringWithJsonSafetyRegexToString(list)
	}
}

func BenchmarkScanner_Large(b *testing.B) {
	list := buildLargeList()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StringWithJsonScanToString(list)
	}
}

func BenchmarkScannerBytes_Large(b *testing.B) {
	list := buildLargeList()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StringWithJsonScanToBytes(list)
	}
}

// ---------------- 1MB 数据对比 ----------------

func build1MBList() []map[string]interface{} {
	// 大约需要 400-500 个 check1 左右的对象达到 1MB
	list := make([]map[string]interface{}, 0, 500)
	for k := 0; k < 250; k++ {
		middle := map[string]interface{}{
			"a": "{\"key\":[\"中文\", \"english\", \"dog\", \"man\"]}",
			"b": "{\"key\":[\"中文\", \"english\", \"dog\", \"man\"],\"key1\":[-1,2,3,4,5],\"key2\":[1,2,3,4,5]}",
			"i": fmt.Sprintf("{\"url\":{\"url\":\"https://xxxxxxxx.com.cn/pic-%d.png\", \"desc\": \"测试长文本脱义验证验证验证\"}}", k),
		}
		list = append(list, check1)
		list = append(list, middle)
	}
	return list
}

func build10MBList() []map[string]interface{} {
	// 在 1MB 数据集基础上放大约 10 倍，模拟更极端的大对象场景
	list := make([]map[string]interface{}, 0, 5000)
	for k := 0; k < 2500; k++ {
		middle := map[string]interface{}{
			"a": "{\"key\":[\"中文\", \"english\", \"dog\", \"man\"]}",
			"b": "{\"key\":[\"中文\", \"english\", \"dog\", \"man\"],\"key1\":[-1,2,3,4,5],\"key2\":[1,2,3,4,5]}",
			"i": fmt.Sprintf("{\"url\":{\"url\":\"https://xxxxxxxx.com.cn/pic-%d.png\", \"desc\": \"10MB benchmark payload\"}}", k),
		}
		list = append(list, check1)
		list = append(list, middle)
	}
	return list
}

func build100MBList() []map[string]interface{} {
	// 在 1MB 数据集基础上放大约 100 倍，用于极端大对象 benchmark
	list := make([]map[string]interface{}, 0, 50000)
	for k := 0; k < 25000; k++ {
		middle := map[string]interface{}{
			"a": "{\"key\":[\"中文\", \"english\", \"dog\", \"man\"]}",
			"b": "{\"key\":[\"中文\", \"english\", \"dog\", \"man\"],\"key1\":[-1,2,3,4,5],\"key2\":[1,2,3,4,5]}",
			"i": fmt.Sprintf("{\"url\":{\"url\":\"https://xxxxxxxx.com.cn/pic-%d.png\", \"desc\": \"100MB benchmark payload\"}}", k),
		}
		list = append(list, check1)
		list = append(list, middle)
	}
	return list
}

func BenchmarkSafetyRegex_1MB(b *testing.B) {
	list := build1MBList()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StringWithJsonSafetyRegexToString(list)
	}
}

func BenchmarkScanner_1MB(b *testing.B) {
	list := build1MBList()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StringWithJsonScanToString(list)
	}
}

func BenchmarkSafetyRegex_10MB(b *testing.B) {
	list := build10MBList()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StringWithJsonSafetyRegexToString(list)
	}
}

func BenchmarkScanner_10MB(b *testing.B) {
	list := build10MBList()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StringWithJsonScanToString(list)
	}
}

func BenchmarkScanner_100MB(b *testing.B) {
	list := build100MBList()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StringWithJsonScanToString(list)
	}
}

func BenchmarkSafetyRegex_100MB(b *testing.B) {
	list := build100MBList()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StringWithJsonSafetyRegexToString(list)
	}
}

func BenchmarkScannerBytes_1MB(b *testing.B) {
	list := build1MBList()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StringWithJsonScanToBytes(list)
	}
}

// ---------------------------------------------------------
// 测试 Depth 深度解析
// ---------------------------------------------------------
func TestDepthDecoding(t *testing.T) {
	input := map[string]interface{}{
		"data": `{"inner1": "{\"inner2\": \"{\\\"deep\\\": 1}\"}"}`,
	}

	out0 := StringWithJsonScanDepthToString(input, 0)
	expected0 := `{"data":{"inner1": "{\"inner2\": \"{\\\"deep\\\": 1}\"}"}}`
	assertJSONEqual(t, "Depth 0", out0, expected0)

	out1 := StringWithJsonScanDepthToString(input, 1)
	expected1 := `{"data":{"inner1":{"inner2": "{\"deep\": 1}"}}}`
	assertJSONEqual(t, "Depth 1", out1, expected1)

	out2 := StringWithJsonScanDepthToString(input, 2)
	expected2 := `{"data":{"inner1":{"inner2":{"deep": 1}}}}`
	assertJSONEqual(t, "Depth 2", out2, expected2)

	outInf := StringWithJsonScanDepthToString(input, -1)
	assertJSONEqual(t, "Depth -1", outInf, expected2)

	outBytes := StringWithJsonScanDepthToBytes(input, -1)
	assertJSONEqual(t, "Bytes Depth -1", string(outBytes), expected2)
}
