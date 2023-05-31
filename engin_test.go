package feng

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestServerStart(t *testing.T) {
	server := NewEngine()
	t.Run("start server", func(t *testing.T) {
		if err := server.Start("localhost:8000"); err != nil {
			t.Fatal(err)
		}
	})
}

func TestServer(t *testing.T) {
	t.Run("add router", func(t *testing.T) {
		server := NewEngine()
		server.GET("/user/1", func(ctx *Context) {
			fmt.Println("URL: ", ctx.request.URL)
		}).GET("/user", func(ctx *Context) {
			fmt.Println("URL: ", ctx.request.URL)
		})
		server.addRoute(http.MethodPost, "/goods/add", func(ctx *Context) {
			fmt.Println("test add router")
		})
		findNode := server.routerTrees[http.MethodPost].findNode(http.MethodPost, "/goods/add", nil)
		if findNode.pattern != "add" {
			t.Fatal("add router failed")
		}
	})

	t.Run("test router handler", func(t *testing.T) {
		server := NewEngine()
		responseData := "test router handler: /user/list"
		server.GET("/user/list", func(ctx *Context) {
			ctx.response.Write([]byte(responseData))
		})
		listen, _ := net.Listen("tcp", "localhost:8000")
		go func(l net.Listener) {
			time.Sleep(1 * time.Second)
			res, err := http.Get("http://localhost:8000/user/list")
			if err != nil {
				t.Errorf("service error: %s", err)
				return
			}
			defer func() {
				res.Body.Close()
				l.Close()
			}()
			body, _ := io.ReadAll(res.Body)
			if string(body) != responseData {
				t.Error("failed to test router handler")
				return
			}
		}(listen)
		http.Serve(listen, server)
	})
}

func TestUseMiddleware(t *testing.T) {
	logMiddlewareCalled := false
	printHostMiddlewareCalled := false
	printUrlMiddlewareCalled := false
	logMiddleware := func(ctx *Context) {
		fmt.Println("logMiddleware: ", ctx)
		logMiddlewareCalled = true
	}
	printHostMiddleware := func(ctx *Context) {
		fmt.Println("printHostMiddleware: ", ctx.request.Host)
		printHostMiddlewareCalled = true
	}
	printUrlMiddleware := func(ctx *Context) {
		fmt.Println("printUrlMiddleware: ", ctx.request.URL.Path)
		printUrlMiddlewareCalled = true
	}
	engine := NewEngine()

	engine.Use(logMiddleware, printHostMiddleware).Use(printUrlMiddleware)

	engine.GET("/user", func(ctx *Context) {
		fmt.Println("this is user handler")
		ctx.Json(http.StatusOK, map[string]any{
			"name": "feng",
			"age":  18,
		})
	})

	listen, err := net.Listen("tcp", "localhost:5000")
	if err != nil {
		t.Error("tcp error: ", err)
	}

	go func(l net.Listener) {
		time.Sleep(100 * time.Millisecond)
		resp, err := http.Get("http://localhost:5000/user")
		if err != nil {
			t.Error(err)
			return
		}
		defer func() {
			resp.Body.Close()
			l.Close()
		}()
		if !logMiddlewareCalled || !printHostMiddlewareCalled || !printUrlMiddlewareCalled {
			t.Error("middlewares were not be called")
		}
		data, _ := io.ReadAll(resp.Body)
		fmt.Println(string(data))
	}(listen)

	http.Serve(listen, engine)
}
