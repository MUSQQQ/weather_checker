package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

const geocodeURL = "https://geocode.xyz/?json=1&scantext="

// "/weather/:cityname"
func Weatherhandler(ctx *fasthttp.RequestCtx) {
	lat, lon, err := getCoordinates(fmt.Sprintf("%s", ctx.UserValue("cityname")))

	fmt.Println(lat, lon, err)
}

func getCoordinates(searchText string) (lat float32, lon float32, err error) {
	var geocodingRequest []byte
	URI := geocodeURL + searchText
	status, geocodingRequest, err := fasthttp.Get(geocodingRequest, URI)
	if err != nil {
		errors.Errorf("error while requesting coordinates from geocoding")
		return -1, -1, err
	}
	if status != 200 {
		errors.Errorf("status not OK in geocoding response")
		return -1, -1, err
	}

	c := make(map[string]json.RawMessage)
	fmt.Println(len(geocodingRequest))
	err = json.Unmarshal(geocodingRequest, &c)

	fmt.Println(string(c["longt"]))
	fmt.Println(string(c["latt"]))
	//fmt.Println(string(geocodingRequest))

	//fmt.Printf(string(geocodingRequest))

	return
}
