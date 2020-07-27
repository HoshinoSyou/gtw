package gtw

import (
	"log"
	"net/http"
)

type HandlerFunc func(ctx *Context)

type Middleware func(HandlerFunc)

type Routers struct {
	router *router
	path   string
}

//New方法创建一个路由引擎
func New() *Routers {
	return &Routers{&router{Handlers: make(map[string]HandlerFunc)}, nil}
}

/*
func (routers *Routers) GET(relativePath string, handlerFunc HandlerFunc) {
	if relativePath[0] != '/' {
		relativePath = "/" + relativePath
	}
	methodAndPath := "GET" + "." + relativePath
	routers.router.Handlers[methodAndPath] = handlerFunc
}

func (routers *Routers) POST(relativePath string, handlerFunc HandlerFunc) {
	if relativePath[0] != '/' {
		relativePath = "/" + relativePath
	}
	methodAndPath := "POST" + "." + relativePath
	routers.router.Handlers[methodAndPath] = handlerFunc
}

func (routers *Routers) PUT(relativePath string, handlerFunc HandlerFunc) {
	if relativePath[0] != '/' {
		relativePath = "/" + relativePath
	}
	methodAndPath := "PUT" + "." + relativePath
	routers.router.Handlers[methodAndPath] = handlerFunc
}

func (routers *Routers) DELETE(relativePath string, handlerFunc HandlerFunc) {
	if relativePath[0] != '/' {
		relativePath = "/" + relativePath
	}
	methodAndPath := "DELETE" + "." + relativePath
	routers.router.Handlers[methodAndPath] = handlerFunc
}*/

//路由组，仅适用于分组
func (routers *Routers) Group(relativePath string) *Routers {
	if relativePath[0] != '/' {
		relativePath = "/" + relativePath
	}
	routers.path = relativePath
	return &Routers{
		router: make(map[string]),
		path:   "",
	}
}

//写一个addHandle方法统一添加路由
func (routers *Routers) addHandle(method string, relativePath string, handlerFunc HandlerFunc) {
	if relativePath[0] != '/' {
		relativePath = "/" + relativePath
	}
	path := routers.path + relativePath
	methodAndPath := method + "." + path
	routers.router.Handlers[methodAndPath] = handlerFunc
}

//GET请求方式
func (routers *Routers) GET(relativePath string, handlerFunc HandlerFunc) {
	routers.addHandle("GET", relativePath, handlerFunc)
}

//POST请求方式
func (routers *Routers) POST(relativePath string, handlerFunc HandlerFunc) {
	routers.addHandle("POST", relativePath, handlerFunc)
}

//PUT请求方式
func (routers *Routers) PUT(relativePath string, handlerFunc HandlerFunc) {
	routers.addHandle("PUT", relativePath, handlerFunc)
}

//DELETE请求方式
func (routers *Routers) DELETE(relativePath string, handlerFunc HandlerFunc) {
	routers.addHandle("DELETE", relativePath, handlerFunc)
}

func (routers *Routers) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	methodAndPath := req.Method + "." + req.URL.Path
	handler, ok := routers.router.Handlers[methodAndPath]
	if !ok {
		log.Printf("%s :404 not found", methodAndPath)
	} else {
		handler(&Context{
			responseWriter: w,
			request:        req,
		})
	}
}

func (routers *Routers) Use(middleware ...Middleware)  {
	
}

func (routers *Routers) Run(port string) error {
	if port == "" {
		port = ":8080"
	}
	return http.ListenAndServe(port, routers)
}
