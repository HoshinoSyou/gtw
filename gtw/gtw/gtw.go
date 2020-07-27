package gtw

import (
	"log"
	"net/http"
	"path"
	"strings"
)

// 用上下文定义路由方法
type HandlerFunc func(ctx *Context)

// 用路由方法定义中间件（方法）
type MiddlewareFunc func(HandlerFunc) HandlerFunc

// 定义一个含*router的结构体，即顶级框架实例
// 且和RoutersGroup互相嵌套
type Routers struct {
	router *router
	RoutersGroup
}

//定义一个空接口，用于传输各种键值对的数据
//在进行JSON、XML渲染时需要用到
type G map[string]interface{}

// New方法创建一个实例化一个Routers
func New() *Routers {
	routers := &Routers{
		router: &router{
			Handlers: make(map[string]HandlerFunc),
		},
	}
	routers.RoutersGroup = RoutersGroup{
		routers: routers,
	}
	return routers
}

// 写一个addHandle方法统一添加路由
func (routers *Routers) addHandler(method string, relativePath string, handlerFunc HandlerFunc) {
	if relativePath == "" {
		relativePath = "/" + relativePath
	}
	if relativePath[0] != '/' {
		relativePath = "/" + relativePath
	}
	finalPath := path.Join(routers.RoutersGroup.path, relativePath)
	finalPath = strings.Replace(finalPath, "\n", "", -1)
	methodAndPath := method + "." + finalPath
	routers.router.Handlers[methodAndPath] = handlerFunc
}

// 开启HTTP服务，将http.ResponseWriter与*http.Request传给上下文Context
func (routers *Routers) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	methodAndPath := req.Method + "." + req.URL.Path
	handler, ok := routers.router.Handlers[methodAndPath]
	if !ok {
		log.Printf("%s :404 not found", methodAndPath)
	} else {
		log.Printf("%s :200 succesefully", methodAndPath)
		handler(&Context{
			responseWriter: w,
			request:        req,
		})
	}
}

// 以port端口运行路由
func (routers *Routers) Run(port string) error {
	if port == "" {
		port = ":8080"
	}
	return http.ListenAndServe(port, routers)
}
