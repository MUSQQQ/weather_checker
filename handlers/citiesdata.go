package handlers

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

// "/v1/cities/:name"
func CityHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Requested city is %s", ctx.UserValue("name"))
}
