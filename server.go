package main

import (
	"log"
	"weather_checker/handlers"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func main() {
	r := router.New()
	r.GET("/", handlers.MainPageHandler)
	r.GET("/cities/{name}", handlers.CityHandler)
	r.GET("/weather/{cityname}", handlers.Weatherhandler)

	log.Fatal(fasthttp.ListenAndServe(":8080", r.Handler))
}
