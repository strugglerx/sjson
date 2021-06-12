package sjson

import (
	"testing"
)

/**
 * @PROJECT_NAME sjson
 * @author  Moqi
 * @date  2021-06-12 11:41
 * @Email:str@li.cm
 **/

var check = map[string]interface{}{
	"key1":        "string",
	"key2":        0,
	"key3":        1.23,
	"nick":        "[天呐]",
	"description": "我叫“王二蛋\"个子不高本事不小",
}

var check1 = map[string]interface{}{
	"key1": "{\"key\":[\"中文\", \"english\", \"dog\", \"man\"]}",
	"key2": "{\"key\":[\"中文\", \"english\", \"dog\", \"man\"],\"key1\":[-1,2,3,4,5],\"key2\":[1,2,3,4,5]}",
	"key3": "[\"中文\", \"english\", \"dog\", \"man\"]",
	"key4": "[{\"url\":\"https://xxxxxxxx.com.cn/pic_2323.1-5-2png\", \"desc\": \"换行\\n换行\"},[\"中文\", \"english\", \"dog\", \"man\"],[-1,2,3,4,5],[1,2,3,4,5]]",
}

func TestJson_StringWithJsonMustToString(t *testing.T) {
	t.Log("eg1:", StringWithJsonToString(check))
	t.Log("eg2:", StringWithJsonToString(check1))
}

func TestStringWithJsonRegexToString(t *testing.T) {
	t.Log("eg1:", StringWithJsonRegexToString(check))
	t.Log("eg2:", StringWithJsonRegexToString(check1))
}
