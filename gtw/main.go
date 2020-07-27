package main

import (
	"gtw/gtw"
)

// 测试
func main() {
	r := gtw.New()
	r.GET("/", func(ctx *gtw.Context) {

		ctx.String(200, "123")
	})
	g := r.Group("/s")
	{
		g.GET("/s", func(ctx *gtw.Context) {
			ctx.String(200, "asd")
		})
	}
	r.Run(":8080")
}
