package feng

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Params map[string]string

type Context struct {
	request  *http.Request
	response http.ResponseWriter
	query    map[string]string
	params   Params
}

func (ctx *Context) Query() map[string]string {
	queryMap := make(map[string]string)
	for _, s := range strings.Split(ctx.request.URL.RawQuery, "&") {
		kvSlice := strings.Split(s, "=")
		queryMap[kvSlice[0]] = kvSlice[1]
	}
	return queryMap
}

func (ctx *Context) GetQueryByKey(key string) any {
	queryMap := ctx.Query()
	return queryMap[key]
}

func (ctx *Context) setParam(key string, value string) {
	if ctx.params == nil {
		ctx.params = make(map[string]string)
	}
	ctx.params[key] = value
}

func (ctx *Context) Param() map[string]string {
	return ctx.params
}

func (ctx *Context) GeParamByKey(key string) string {
	return ctx.params[key]
}

func (ctx *Context) Status(code int) {
	ctx.response.WriteHeader(code)
}

func (ctx *Context) Json(code int, obj any) error {
	ctx.response.Header().Set("Content-Type", "application/json")
	ctx.Status(code)
	bytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	_, err = ctx.response.Write(bytes)
	return err
}
