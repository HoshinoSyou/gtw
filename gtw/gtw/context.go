package gtw

import (
	"encoding/json"
	"encoding/xml"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"time"
)

// Param是url中path参数的键值对
// 键key表示该path参数的名称
// 值value是该path参数的值
type Param struct {
	key   string
	value string
}

// 因path参数可能不止一个，所以此处定义Params切片
type Params []Param

// 上下文，使用原生http标准库，封装http.Request和http.ResponseWriter方法，以及含有用户数据和Params切片
type Context struct {
	responseWriter http.ResponseWriter
	request        *http.Request
	data           map[string]interface{}
	Params         Params
	path           string
	httpMethod     string
	statusCode     int
}

// 获取url中的path参数
func (ctx *Context) Param(key string) (value string) {
	for _, param := range ctx.Params {
		if param.key == key {
			value = param.value
			return
		}
	}
	return
}

// 获取form表单参数
func (ctx *Context) PostForm(key string) string {
	value := ctx.request.FormValue(key)
	return value
}

// 获取url中？后的querystring参数
func (ctx *Context) Query(key string) string {
	value := ctx.request.URL.Query().Get(key)
	return value
}

// 获取int型数据
func (ctx *Context) GetInt(key string) (bool, int) {
	data, res := ctx.data[key]
	if res && reflect.TypeOf(data).String() == "int" {
		return true, data.(int)
	}
	return false, 0
}

// 获取string型数据
func (ctx *Context) GetString(key string) (bool, string) {
	data, res := ctx.data[key]
	if res && reflect.TypeOf(data).String() == "string" {
		return true, data.(string)
	}
	return false, ""
}

// 获取bool型数据
func (ctx *Context) GetBool(key string) (bool, bool) {
	data, res := ctx.data[key]
	if res && reflect.TypeOf(data).String() == "bool" {
		return true, data.(bool)
	}
	return false, false
}

// 获取time.Time型数据
func (ctx *Context) GetTime(key string) (bool, time.Time) {
	data, res := ctx.data[key]
	if res && reflect.TypeOf(data).String() == "time.Time" {
		return true, data.(time.Time)
	}
	return false, time.Time{}
}

// 等价
// 使用Get方法获取数据
/* 因上面的GetXXX方法代码基本相同可改用一下方法
func (ctx *Context) Get(key string) {
	data, res := ctx.data[key]
	if res {
		switch reflect.TypeOf(data).String() {
		case "string":
			ctx.GetString(data.(string))
		case "int":
			ctx.GetInt(data.(int))
		case "bool":
			ctx.GetBool(data.(bool))
		case "time.Time":
			ctx.GetTime(data.(time.Time))
		default:
			log.Printf("获取数据类型失败！")
		}
	} else {
		log.Printf("获取数据失败！")
	}
}

func (ctx *Context) GetString(str string) string {
	return str
}

func (ctx *Context) GetInt(i int) int {
	return i
}

func (ctx *Context) GetBool(b bool) bool {
	return b
}

func (ctx *Context) GetTime(t time.Time) time.Time {
	return t
}*/

// 返回string类型数据
func (ctx *Context) String(code int, str string) {
	ctx.responseWriter.WriteHeader(code)
	ctx.responseWriter.Write([]byte(str))
}

// 返回[]byte类型数据
func (ctx *Context) Byte(code int, byte []byte) {
	ctx.responseWriter.WriteHeader(code)
	ctx.responseWriter.Write(byte)
}

// 使用render方法判断格式进行渲染
// 传入form来判断渲染为哪种格式
// 可改为Render方法单独调用
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

// JSON渲染
func (ctx *Context) JSON(code int, obj interface{}) {
	ctx.render("json", code, obj)
}

// XML渲染
func (ctx *Context) XML(code int, obj interface{}) {
	ctx.render("xml", code, obj)
}

// 等价
/* 因JSON与XML方法除序列化外代码基本相同，因此合并为上面的方法
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

// SetCookies方法设置cookies
// 使用http包中的SetCookie方法与http.Cookie结构体
// Name是Cookie的名称
// Value是Cookie的值
// Domain是Cookie作用的域名
// Path是Cookie作用的路径
// MaxAge是Cookie的作用时间
// Secure是Cookie的Secure属性，即Cookie是否仅在浏览器进行安全或加密连接时才能被使用
// HttpOnly是Cookie的HttpOnly属性，即Cookie是否仅在浏览器进行HTTP与HTTPS请求时才会被暴露
func (ctx *Context) SetCookies(name string, value string, path string, domain string, maxAge int, secure bool, httpOnly bool) {
	if path[0] != '/' || path == "" {
		path = "/" + path
	}
	http.SetCookie(ctx.responseWriter, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     path,
		Domain:   domain,
		MaxAge:   maxAge,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
}

// 返回请求中名为name的cookie，如果未找到该cookie会返回nil, ErrNoCookie
func (ctx *Context) Cookies(name string) (value string, err error) {
	var cookie *http.Cookie
	cookie, err = ctx.request.Cookie(name)
	if err != nil {
		return "", err
	}
	value, err = url.QueryUnescape(cookie.Value)
	return value, nil
}

// 返回状态码
func (ctx *Context) Status(code int) {
	ctx.statusCode = code
	ctx.responseWriter.WriteHeader(code)
}

// 重定向（未完成）
// func (ctx *Context) Redirect(code int, location string) {}

// 上传文件（未完成）
// func (ctx *Context) UploadFile(name string) (*multipart.FileHeader, error) {}
