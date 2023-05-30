package feng

import (
	"fmt"
	"net/http"
	"testing"
)

func TestAddNode(t *testing.T) {
	path := "/user/member/list"
	path1 := "/user/member"
	rootNode := &node{
		pattern: "/",
	}
	handlers := []HandlerFunc{
		func(ctx *Context) {
			fmt.Println("111")
		},
		func(ctx *Context) {
			fmt.Println("222")
		},
	}
	t.Run("test addNode  case 1", func(t *testing.T) {
		rootNode.addNode(http.MethodGet, path, handlers...)
		userPatterNode := rootNode.children[0]
		memberPatterNode := userPatterNode.children[0]
		listPatterNode := memberPatterNode.children[0]
		if len(rootNode.children) != 1 || userPatterNode.pattern != "user" || memberPatterNode.pattern != "member" || listPatterNode.pattern != "list" {
			t.Fatal("test failed")
		}
		originHandlers, _ := fmt.Printf("%p", handlers)
		realHandlers, _ := fmt.Printf("%p", listPatterNode.handlers)
		if originHandlers != realHandlers {
			t.Fatal("test failed")
		}
	})
	t.Run("test addNode  case 2", func(t *testing.T) {
		rootNode.addNode(http.MethodGet, path1, handlers...)
		userPatterNode := rootNode.children[0]
		memberPatterNode := userPatterNode.children[0]
		if len(rootNode.children) != 1 || userPatterNode.pattern != "user" || memberPatterNode.pattern != "member" {
			t.Fatal("test failed")
		}
		originHandlers, _ := fmt.Printf("%p", handlers)
		realHandlers, _ := fmt.Printf("%p", memberPatterNode.handlers)
		if originHandlers != realHandlers {
			t.Fatal("test failed")
		}
	})
}

func TestFindNode(t *testing.T) {
	handlers := []HandlerFunc{
		func(ctx *Context) {
			fmt.Println("111")
		},
		func(ctx *Context) {
			fmt.Println("222")
		},
	}
	rootNode := &node{
		pattern:  "/",
		handlers: nil,
		children: []*node{
			&node{
				pattern:  "use",
				handlers: nil,
				children: []*node{
					&node{
						pattern:  "member",
						handlers: handlers,
						children: []*node{
							&node{
								pattern:  "list",
								handlers: nil,
								children: nil,
							},
						},
					},
				},
			},
		},
	}
	findNode := rootNode.findNode(http.MethodGet, "/user/member", nil)
	if findNode == nil || findNode.pattern != "member" {
		t.Fatal("TestFindNode failed")
	}
	for i, handler := range findNode.handlers {
		handlerAddr, _ := fmt.Printf("%p", handlers[i])
		findHandlerAddr, _ := fmt.Printf("%p", handler)
		if handlerAddr != findHandlerAddr {
			t.Fatal("TestFindNode failed")
		}
	}
}
