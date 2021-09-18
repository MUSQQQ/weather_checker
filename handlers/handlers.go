package handlers

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

func MainPageHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Hello")
}
