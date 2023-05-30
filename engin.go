package feng

import (
	"github.com/shavenfeng/feng/utils"
	"net"
	"net/http"
)

type HandlerFunc func(ctx *Context)

type Engine struct {
	*RouterGroup
	routerTrees trees
}

func NewEngine() *Engine {
	engine := &Engine{
		routerTrees: make(trees),
		RouterGroup: &RouterGroup{bathPath: "/", handlers: []HandlerFunc{}},
	}
	engine.RouterGroup.engin = engine
	return engine
}

func (engine *Engine) Start(address string) error {
	listen, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	return http.Serve(listen, engine)
}

func (engine *Engine) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	engine.handle(responseWriter, request)
}

func (engine *Engine) handle(responseWriter http.ResponseWriter, request *http.Request) {
	treeNode := engine.routerTrees[request.Method].findNode(request.Method, request.URL.Path)
	for _, handler := range treeNode.handlers {
		handler(&Context{
			request:  *request,
			response: responseWriter,
		})
	}
}

func (engine *Engine) addRoute(method, path string, handlers ...HandlerFunc) {
	utils.CheckPath(path)
	root := engine.routerTrees[method]
	if root == nil {
		root = &node{
			pattern: "/",
		}
		engine.routerTrees[method] = root
	}
	root.addNode(method, path, handlers...)
}