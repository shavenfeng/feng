package feng

import (
	"encoding/json"
	"github.com/go-playground/assert/v2"
	"io"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestUseMiddleware(t *testing.T) {
	callIndex := []int{}
	type info struct {
		Name string
		Work string
	}
	res := info{"feng", "boss"}
	logMiddleware := func(ctx *Context) {
		callIndex = append(callIndex, 1)
		ctx.Next()
		callIndex = append(callIndex, 7)
	}
	printHostMiddleware := func(ctx *Context) {
		callIndex = append(callIndex, 2)
		ctx.Next()
		callIndex = append(callIndex, 6)
	}
	printUrlMiddleware := func(ctx *Context) {
		callIndex = append(callIndex, 3)
		ctx.Next()
		callIndex = append(callIndex, 5)
	}
	engine := NewEngine()

	engine.Use(logMiddleware, printHostMiddleware).Use(printUrlMiddleware)

	engine.GET("/user", func(ctx *Context) {
		callIndex = append(callIndex, 4)
		ctx.Json(http.StatusOK, res)
	})

	listen, err := net.Listen("tcp", "localhost:5000")
	assert.Equal(t, err, nil)

	go func(l net.Listener) {
		time.Sleep(100 * time.Millisecond)
		resp, err := http.Get("http://localhost:5000/user")
		defer func() {
			if resp != nil {
				resp.Body.Close()
			}
			l.Close()
		}()
		if err != nil {
			t.Error(err)
			return
		}

		for i := 0; i < 7; i++ {
			assert.Equal(t, callIndex[i], i+1)
		}
		data, _ := io.ReadAll(resp.Body)
		i := info{}
		assert.Equal(t, json.Unmarshal(data, &i), nil)
		assert.Equal(t, i.Name, "feng")
		assert.Equal(t, i.Work, "boss")
	}(listen)

	http.Serve(listen, engine)
}

func TestAbortMiddleware(t *testing.T) {
	callIndex := []int{}
	middleware1 := func(c *Context) {
		callIndex = append(callIndex, 1)
	}
	middleware2 := func(c *Context) {
		callIndex = append(callIndex, 2)
		c.Abort()
	}
	middleware3 := func(c *Context) {
		callIndex = append(callIndex, 3)
	}
	ctx := &Context{
		index:    0,
		handlers: []HandlerFunc{middleware1, middleware2, middleware3},
	}
	ctx.execHandlers()
	assert.Equal(t, len(callIndex), 2)
	for index, v := range callIndex {
		assert.Equal(t, index+1, v)
	}
}

func TestNotFoundMiddleware(t *testing.T) {
	engine := NewEngine()
	engine.GET("/user", func(ctx *Context) {
	})
	listen, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		t.Error(err)
		return
	}
	go func(l net.Listener) {
		time.Sleep(100 * time.Millisecond)
		response, err := http.Get("http://localhost:8000/user/list")
		defer func() {
			if response != nil {
				response.Body.Close()
			}
			l.Close()
		}()
		assert.Equal(t, err, nil)
		data, err := io.ReadAll(response.Body)
		assert.Equal(t, err, nil)
		assert.Equal(t, string(data), "Not Found!")
	}(listen)
	http.Serve(listen, engine)
}
