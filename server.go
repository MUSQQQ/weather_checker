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
	r.GET("/coordinates/{cityname}", handlers.CoordinatesHandler)
	r.GET("/weather/{cityname}", handlers.WeatherHandler)
	r.GET("/weather/checkcity/{cityname}", handlers.MainWeatherHandler)

	log.Fatal(fasthttp.ListenAndServe(":8080", r.Handler))
}
