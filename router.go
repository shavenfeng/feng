package feng

import (
	"github.com/shavenfeng/feng/utils"
	"net/http"
)

type Route interface {
	GET(string, ...HandlerFunc) Route
	POST(string, ...HandlerFunc) Route
	DELETE(string, ...HandlerFunc) Route
	PATCH(string, ...HandlerFunc) Route
	PUT(string, ...HandlerFunc) Route
	OPTIONS(string, ...HandlerFunc) Route
	HEAD(string, ...HandlerFunc) Route
}

type RouterGroup struct {
	bathPath string
	handlers []HandlerFunc
	engin    *Engine
}

func (r *RouterGroup) Group(path string, handlers ...HandlerFunc) RouterGroup {
	utils.CheckPath(path)
	return RouterGroup{
		bathPath: path,
		handlers: handlers,
		engin:    r.engin,
	}
}

func (r *RouterGroup) addHandler(handlers ...HandlerFunc) {
	r.handlers = append(r.handlers, handlers...)
}

func (r *RouterGroup) getAbsolutePath(relativePath string) string {
	if r.bathPath == "/" {
		return relativePath
	}
	return r.bathPath + relativePath
}

func (r *RouterGroup) combineHandlers(handlers []HandlerFunc) []HandlerFunc {
	mergeHandlers := make([]HandlerFunc, len(r.handlers)+len(handlers))
	copy(mergeHandlers, r.handlers)
	copy(mergeHandlers[len(r.handlers):], handlers)
	return mergeHandlers
}

func (r *RouterGroup) handle(method, relativePath string, handlers ...HandlerFunc) {
	absolutePath := r.getAbsolutePath(relativePath)
	mergeHandlers := r.combineHandlers(handlers)
	r.engin.addRoute(method, absolutePath, mergeHandlers...)
}

func (r *RouterGroup) GET(relativePath string, handlers ...HandlerFunc) Route {
	r.handle(http.MethodGet, relativePath, handlers...)
	return r
}

func (r *RouterGroup) POST(relativePath string, handlers ...HandlerFunc) Route {
	r.handle(http.MethodPost, relativePath, handlers...)
	return r
}

func (r *RouterGroup) DELETE(relativePath string, handlers ...HandlerFunc) Route {
	r.handle(http.MethodDelete, relativePath, handlers...)
	return r
}

func (r *RouterGroup) PATCH(relativePath string, handlers ...HandlerFunc) Route {
	r.handle(http.MethodPatch, relativePath, handlers...)
	return r
}

func (r *RouterGroup) PUT(relativePath string, handlers ...HandlerFunc) Route {
	r.handle(http.MethodPut, relativePath, handlers...)
	return r
}

func (r *RouterGroup) OPTIONS(relativePath string, handlers ...HandlerFunc) Route {
	r.handle(http.MethodOptions, relativePath, handlers...)
	return r
}

func (r *RouterGroup) HEAD(relativePath string, handlers ...HandlerFunc) Route {
	r.handle(http.MethodHead, relativePath, handlers...)
	return r
}
