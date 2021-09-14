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
	"a":        "string",
	"b":        0,
	"c":        1.23,
	"d":        "[天呐]",
	"e": "我叫“王二蛋\"个子不高本事不小",
	"f": "请认真完善相关信息(）单位为\"m\",幕墙面积等，完善，信息。",
	"g": "请认真完善相关信息(）单位为,\"m\",幕墙面积等，完善，信息。",
}

var check1 = map[string]interface{}{
	"a": "{\"key\":[\"中文\", \"english\", \"dog\", \"man\"]}",
	"b": "{\"key\":[\"中文\", \"english\", \"dog\", \"man\"],\"key1\":[-1,2,3,4,5],\"key2\":[1,2,3,4,5]}",
	"c": "[1,\"中文\",2,3, \"english\", \"dog\", \"man\"]",
	"d": "[{\"url\":\"[请认真完善相关信息(）单{位}为\"m\",幕墙面积等，完善，信息。]\",\"list\":[\"https://xxxxxxxx.com.cn/pic_2323.1-5-2png\",\"https://xxxxxxxx.com.cn/pic_2323.1-5-2png\",\"https://xxxxxxxx.com.cn/pic_2323.1-5-2png\"], \"desc\": \"换行\\n换行\"},[\"中文\", \"english\", \"dog\", \"man\"],[-1,2,3,4,5],[1,2,3,4,5]]",
	"e": "[]",
	"f": "{}",
	"g": "[[\"中文\", \"english\", \"dog\", \"man\"],[\"中文\", \"english\", \"dog\", \"man\"]]",
	"h": "[[1,2,3,4,14,5.3],[1,2,3,4,14,5.3]]",
	"i": "{\"url\":{\"url\":\"https://xxxxxxxx.com.cn/pic_2323.1-5-2png\", \"desc\": \"换行\\n换行\"}}",
	"j": "[\"https://xxxxxxxx.com.cn/pic_2323.1-5-2png\",\"https://xxxxxxxx.com.cn/pic_2323.1-5-2png\",\"https://xxxxxxxx.com.cn/pic_2323.1-5-2png\"]",
}

func TestJson_StringWithJsonMustToString(t *testing.T) {
	t.Log("eg1:", StringWithJsonToString(check))
	t.Log("eg2:", StringWithJsonToString(check1))
}

func TestStringWithJsonRegexToString(t *testing.T) {
	t.Log("eg1:", StringWithJsonRegexToString(check))
	t.Log("eg2:", StringWithJsonRegexToString(check1))
}

func TestStringWithJsonSafetyRegexToString(t *testing.T) {
	t.Log("eg1:", StringWithJsonSafetyRegexToString(check))
	t.Log("eg2:", StringWithJsonSafetyRegexToString(check1))
}
