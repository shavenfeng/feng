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
