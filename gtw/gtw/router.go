package gtw

// 含路由方法的map与中间件（方法）的切片
type router struct {
	Handlers   map[string]HandlerFunc
	Middleware []MiddlewareFunc
}

// 路由组，主要用于path处理
// 和Routers互相嵌套
type RoutersGroup struct {
	path    string
	routers *Routers
}

// 等价部分
/* 因四个restful请求方法实现的代码基本相同，用后面的代码代替
func (routersGroup *RoutersGroup) GET(relativePath string, handlerFunc HandlerFunc) {
	if relativePath[0] != '/' {
		relativePath = "/" + relativePath
	}
	methodAndPath := "GET" + "." + relativePath
	routersGroup.routers.router.Handlers[methodAndPath] = handlerFunc
}

func (routersGroup *RoutersGroup) POST(relativePath string, handlerFunc HandlerFunc) {
	if relativePath[0] != '/' {
		relativePath = "/" + relativePath
	}
	methodAndPath := "POST" + "." + relativePath
	routersGroup.routers.router.Handlers[methodAndPath] = handlerFunc
}

func (routersGroup *RoutersGroup) PUT(relativePath string, handlerFunc HandlerFunc) {
	if relativePath[0] != '/' {
		relativePath = "/" + relativePath
	}
	methodAndPath := "PUT" + "." + relativePath
	routersGroup.routers.router.Handlers[methodAndPath] = handlerFunc
}

func (routersGroup *RoutersGroup) DELETE(relativePath string, handlerFunc HandlerFunc) {
	if relativePath[0] != '/' {
		relativePath = "/" + relativePath
	}
	methodAndPath := "DELETE" + "." + relativePath
	routersGroup.routers.router.Handlers[methodAndPath] = handlerFunc
}*/

// 路由组，仅适用于分组，未添加公用中间件
func (routersGroup *RoutersGroup) Group(relativePath string) *RoutersGroup {
	if relativePath == "" {
		relativePath = "/" + relativePath
	}
	if relativePath[0] != '/' {
		relativePath = "/" + relativePath
	}
	return &RoutersGroup{
		path:    relativePath,
		routers: routersGroup.routers,
	}
}

// GET请求方式
func (routersGroup *RoutersGroup) GET(relativePath string, handlerFunc HandlerFunc) {
	routersGroup.routers.addHandler("GET", relativePath, handlerFunc)
}

// POST请求方式
func (routersGroup *RoutersGroup) POST(relativePath string, handlerFunc HandlerFunc) {
	routersGroup.routers.addHandler("POST", relativePath, handlerFunc)
}

// PUT请求方式
func (routersGroup *RoutersGroup) PUT(relativePath string, handlerFunc HandlerFunc) {
	routersGroup.routers.addHandler("PUT", relativePath, handlerFunc)
}

// DELETE请求方式
func (routersGroup *RoutersGroup) DELETE(relativePath string, handlerFunc HandlerFunc) {
	routersGroup.routers.addHandler("DELETE", relativePath, handlerFunc)
}

// 使用中间件
func (routersGroup *RoutersGroup) Use(middleware ...MiddlewareFunc) {
	routersGroup.routers.router.Middleware = append(routersGroup.routers.router.Middleware, middleware...)
}
