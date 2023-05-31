package feng

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"
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
		goodsGetNode := getRoot.findNode(http.MethodGet, "/goods", nil)
		userGetNode := getRoot.findNode(http.MethodGet, "/user", nil)
		userPostNode := postRoot.findNode(http.MethodPost, "/user", nil)
		userDeleteNode := deleteRoot.findNode(http.MethodDelete, "/user", nil)
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

func TestRouterQueryAndParams(t *testing.T) {
	engine := NewEngine()
	engine.GET("/user", func(ctx *Context) {
		query := ctx.Query()
		queryJson, _ := json.Marshal(query)
		ctx.response.Write([]byte(queryJson))
	})
	engine.GET("/user/:id/:name", func(ctx *Context) {
		paramsJson, _ := json.Marshal(ctx.Param())
		ctx.response.Write([]byte(paramsJson))
	})
	l, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		t.Fatal("tcp server started failed")
	}
	wg := &sync.WaitGroup{}
	t.Run("test router query", func(t *testing.T) {
		wg.Add(1)
		go func(l net.Listener) {
			time.Sleep(100 * time.Millisecond)
			resp, err := http.Get("http://localhost:8000/user?id=1&type=2&work=boss")
			defer func() {
				if resp != nil {
					resp.Body.Close()
				}
				wg.Done()
			}()
			if err != nil {
				t.Fatal("Get: http://localhost:8000/user?id=1&type=2&work=boss failed")
			}
			resData, _ := io.ReadAll(resp.Body)
			type s struct {
				Id   string
				Type string
				Work string
			}
			mapS := s{}
			if err := json.Unmarshal(resData, &mapS); err != nil {
				t.Errorf("failed to unmarshal json：%s", err)
				return
			}
			if mapS.Id != "1" || mapS.Type != "2" && mapS.Work != "boss" {
				t.Error("test router query failed")
				return
			}
		}(l)
	})
	t.Run("test router params", func(t *testing.T) {
		wg.Add(1)
		go func(l net.Listener) {
			time.Sleep(100 * time.Millisecond)
			resp, err := http.Get("http://localhost:8000/user/111/feng")
			defer func() {
				if resp != nil {
					resp.Body.Close()
				}
				wg.Done()
			}()
			if err != nil {
				t.Error("Get: http://localhost:8000/user/111/feng failed")
				return
			}
			bytes, _ := io.ReadAll(resp.Body)
			type Params struct {
				Id   string
				Name string
			}
			params := Params{}
			if err := json.Unmarshal(bytes, &params); err != nil {
				t.Errorf("failed to unmarshal json：%s", err)
				return
			}
			if params.Id != "111" || params.Name != "feng" {
				t.Error("test router params failed")
				return
			}
		}(l)
	})
	go func() {
		wg.Wait()
		l.Close()
	}()
	http.Serve(l, engine)
}
