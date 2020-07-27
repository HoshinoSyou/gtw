package gtw

import (
	"encoding/json"
	"encoding/xml"
	"log"
	"net/http"
	"reflect"
	"time"
)

type Context struct {
	responseWriter http.ResponseWriter
	request        *http.Request
	data           map[string]interface{}
}

//返回string类型数据
func (ctx *Context) String(code int, str string) {
	ctx.responseWriter.WriteHeader(code)
	ctx.responseWriter.Write([]byte(str))
}

//使用render方法判断格式进行渲染
//传入form来判断渲染为哪种格式
//可改为Render方法单独调用
func (ctx *Context) render(form string, code int, obj interface{}) {
	var bytes []byte
	var err error
	switch form {
	case "json":
		bytes, err = json.Marshal(&obj)
	case "xml":
		bytes, err = xml.Marshal(&obj)
	}
	if err != nil {
		log.Fatalf("初始化%s格式失败，错误信息：%v", bytes, err)
	}
	ctx.responseWriter.WriteHeader(code)
	ctx.responseWriter.Write(bytes)
}

//JSON渲染
func (ctx *Context) JSON(code int, obj interface{}) {
	ctx.render("json", code, obj)
}

//XML渲染
func (ctx *Context) XML(code int, obj interface{}) {
	ctx.render("xml", code, obj)
}

//等价
/*
func (ctx *Context) JSON(code int, obj interface{}) {
	bytes, err := json.Marshal(&obj)
	if err != nil {
		log.Fatalf("初始化JSON格式失败，错误信息：%v\n", err.Error())
	}
	ctx.responseWriter.WriteHeader(code)
	ctx.responseWriter.Write(bytes)
}

func (ctx *Context) XML(code int, obj interface{}) {
	bytes, err := xml.Marshal(&obj)
	if err != nil {
		log.Fatalf("初始化XML格式失败，错误信息：%v\n", err.Error())
	}
	ctx.responseWriter.WriteHeader(code)
	ctx.responseWriter.Write(bytes)
}
*/

func (ctx *Context) GetInt(key string) (bool, int) {
	data, res := ctx.data[key]
	if res && reflect.TypeOf(data).String() == "int" {
		return true, data.(int)
	}
	return false, 0
}

func (ctx *Context) GetString(key string) (bool, string) {
	data, res := ctx.data[key]
	if res && reflect.TypeOf(data).String() == "string" {
		return true, data.(string)
	}
	return false, ""
}

func (ctx *Context) GetBool(key string) (bool, bool) {
	data, res := ctx.data[key]
	if res && reflect.TypeOf(data).String() == "bool" {
		return true, data.(bool)
	}
	return false, false
}

func (ctx *Context) GetTime(key string) (bool, time.Time) {
	data, res := ctx.data[key]
	if res && reflect.TypeOf(data).String() == "time.Time" {
		return true, data.(time.Time)
	}
	return false, time.Time{}
}
