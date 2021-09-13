# sjson

[![Go Doc](https://godoc.org/github.com/strugglerx/sjson?status.svg)](https://godoc.org/github.com/strugglerx/sjson)
[![Production Ready](https://img.shields.io/badge/production-ready-blue.svg)](https://github.com/strugglerx/sjson)
[![License](https://img.shields.io/github/license/strugglerx/sjson.svg?style=flat)](https://github.com/strugglerx/sjson)

convert `string json` to `json`
# Feature

```json
{
    "key1": "{\"key\":[\"中文\", \"english\", \"dog\", \"man\"]}",
    "key2": "{\"key\":[\"中文\", \"english\", \"dog\", \"man\"],\"key1\":[-1,2,3,4,5],\"key2\":[1,2,3,4,5]}",
    "key3": "[\"中文\", \"english\", \"dog\", \"man\"]",
    "key4": "[{\"url\":\"https://xxxxxxxx.com.cn/pic_2323.1-5-2png\", \"desc\": \"换行\\n换行\"},[\"中文\", \"english\", \"dog\", \"man\"],[-1,2,3,4,5],[1,2,3,4,5]]"
}
```
convert
```json
{   
  "key1":{"key":["中文", "english", "dog", "man"]},
  "key2":{"key":["中文", "english", "dog", "man"],
  "key1":[-1,2,3,4,5],"key2":[1,2,3,4,5]},
  "key3":["中文", "english", "dog", "man"],
  "key4":[{"url":"https://xxxxxxxx.com.cn/pic_2323.1-5-2png", "desc": "换行\n换行"},["中文", "english", "dog", "man"],[-1,2,3,4,5],[1,2,3,4,5]]
}
```



# Installation
```
go get -u -v github.com/strugglerx/sjson
```
suggested using `go.mod`:
```
require github.com/strugglerx/sjson
```

# Usage
```golang
package main 

import (
    "github.com/strugglerx/sjson"
    "fmt"
)

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


func main() {
    fmt.Println("eg1:",sjson.StringWithJsonToString(check))
    fmt.Println("eg2:", sjson.StringWithJsonToString(check1))
    
    fmt.Println("eg1:", sjson.StringWithJsonRegexToString(check))
    fmt.Println("eg2:", sjson.StringWithJsonRegexToString(check1))
    
    //suggest 
    fmt.Println("eg1:", sjson.StringWithJsonSafetyRegexToString(check))
    fmt.Println("eg2:", sjson.StringWithJsonSafetyRegexToString(check1))
}

```

# License

`sjson` is licensed under the [MIT License](LICENSE), 100% free and open-source, forever.