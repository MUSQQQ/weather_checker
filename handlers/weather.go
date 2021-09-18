package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

const geocodeURL = "https://geocode.xyz/?json=1&scantext="

// "/weather/:cityname"
func Weatherhandler(ctx *fasthttp.RequestCtx) {
	lon, lat, status, err := getCoordinates(fmt.Sprintf("%s", ctx.UserValue("cityname")))
	if err != nil {
		log.Printf("error occured while trying to get coordinates: %s", err)
		return
	}
	if status >= 500 {
		fmt.Fprintf(ctx, "Unfortunately, geocode services are unavailable at the moment. Please try again later.")
		return
	}
	if lat == -1 && lon == -1 {
		fmt.Fprintf(ctx, "Geocode service got too many requests and have not processed your search. Please try again.")
		return
	}

	fmt.Fprintf(ctx, "Coordinates of the looked up city:\n")
	fmt.Fprintf(ctx, "lat: %f, lon: %f", lat, lon)
}

func getCoordinates(searchText string) (lat float32, lon float32, status int, err error) {
	var geocodingRequest []byte
	URI := geocodeURL + searchText
	status, geocodingRequest, err = fasthttp.Get(geocodingRequest, URI)
	if err != nil {
		errors.Errorf("error while requesting coordinates from geocoding")
		return -1, -1, status, err
	}
	if status >= 500 {
		log.Printf("geocode service unvailable")
		return -1, -1, status, nil
	}
	if status != 200 {
		errors.Errorf("status not OK in geocoding response")
		return -1, -1, status, nil
	}

	c := make(map[string]json.RawMessage)

	err = json.Unmarshal(geocodingRequest, &c)
	if err != nil {
		errors.Errorf("error while unmarshaling request")
		return -1, -1, status, err
	}

	longt, err := byteArrayToFloat(c["longt"])
	if err != nil {
		log.Printf("error while converting coordinates to float: %s", err)
		return -1, -1, status, err
	}
	latt, err := byteArrayToFloat(c["latt"])
	if err != nil {
		log.Printf("error while converting coordinates to float: %s", err)
		return -1, -1, status, err
	}

	return longt, latt, 200, nil
}

func byteArrayToFloat(bytes []byte) (result float32, err error) {
	strByte := string(bytes)
	var i int
	for i = 0; i < len(strByte); i++ {
		if strByte[i] == '.' {
			break
		}
	}

	fmt.Println(strByte[1:i])
	fmt.Println(strByte[i+1 : len(strByte)-1])
	intResultPart, err := strconv.Atoi(strByte[1:i])
	if err != nil {
		log.Printf("error while converting: %s", err)
		return -1, err
	}
	mantissaPart, err := strconv.Atoi(strByte[i+1 : len(strByte)-1])
	if err != nil {
		log.Printf("error while converting: %s", err)
		return -1, err
	}
	result = float32(intResultPart) + float32(mantissaPart)/(float32(math.Pow10(i+2)))
	fmt.Println(result)
	return result, nil
}
