package sjson

import (
	"encoding/json"
	"fmt"
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
