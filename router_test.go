package feng

import (
	"net/http"
	"testing"
)

func TestRouterGroup(t *testing.T) {
	t.Run("test custom router group", func(t *testing.T) {
		handlerFunc := func(ctx *Context) {
			ctx.response.Write([]byte(ctx.request.Host))
		}
		engine := NewEngine()
		engine.GET("/goods/photo", handlerFunc)
		group := engine.Group("/goods", handlerFunc)
		{
			group.GET("/detail", handlerFunc)
			group.GET("/list", handlerFunc)
		}
		group1 := engine.Group("/user", handlerFunc)
		{
			group1.GET("/detail", handlerFunc)
			group1.GET("/list", handlerFunc)
			group1.POST("/login", handlerFunc)
			group1.POST("/register", handlerFunc)
			group1.DELETE("/delete/:id", handlerFunc)
		}
		getRoot := group.engin.routerTrees[http.MethodGet]
		postRoot := group.engin.routerTrees[http.MethodPost]
		deleteRoot := group.engin.routerTrees[http.MethodDelete]
		goodsGetNode := getRoot.findNode(http.MethodGet, "/goods")
		userGetNode := getRoot.findNode(http.MethodGet, "/user")
		userPostNode := postRoot.findNode(http.MethodPost, "/user")
		userDeleteNode := deleteRoot.findNode(http.MethodDelete, "/user")
		if len(goodsGetNode.children) != 3 || goodsGetNode.children[0].pattern != "detail" || goodsGetNode.children[1].pattern != "list" {
			t.Fatal("未能正确添加GET:/goods")
		}
		if len(userGetNode.children) != 2 || userGetNode.children[0].pattern != "detail" || userGetNode.children[1].pattern != "list" {
			t.Fatal("未能正确添加GET:/goods")
		}
		if len(userPostNode.children) != 2 || userPostNode.children[0].pattern != "login" || userPostNode.children[1].pattern != "register" {
			t.Fatal("未能正确添加POST:/user")
		}
		if len(userDeleteNode.children) != 1 || userDeleteNode.children[0].pattern != ":id" {
			t.Fatal("未能正确添加DELETE:/user")
		}

	})
}
