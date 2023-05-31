package feng

import "net/http"

func NotFindMiddleware(ctx *Context) {
	ctx.response.WriteHeader(http.StatusNotFound)
	ctx.response.Write([]byte("Not Found!"))
}
