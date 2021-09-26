package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/valyala/fasthttp"
)

const openWeatherURL = "https://api.openweathermap.org/data/2.5/onecall"

// "/coordinates/:cityname"
func WeatherHandler(ctx *fasthttp.RequestCtx) {
	lon, lat, status, err := getCoordinatesAsString(fmt.Sprintf("%s", ctx.UserValue("cityname")))
	if err != nil {
		log.Printf("error occured while trying to get coordinates: %s", err)
		return
	}
	if status >= 500 {
		fmt.Fprintf(ctx, "Unfortunately, geocode services are unavailable at the moment. Please try again later.")
		return
	}
	if lat == "" && lon == "" {
		fmt.Fprintf(ctx, "Geocode service got too many requests and have not processed your search. Please try again.")
		return
	}

	ctx.Response.Header.Add("Content-Type", "application-json")

	jsonMap := make(map[string]json.RawMessage)
	jsonMap["latt"] = []byte(lat)
	jsonMap["longt"] = []byte(lon)

	a, err := json.Marshal(jsonMap)
	if err != nil {
		fmt.Fprintf(ctx, "marshaling error")
	}
	ctx.Response.SetBody(a)

}

func getCoordinatesAsString(searchText string) (lat string, lon string, status int, err error) {
	var geocodingRequest []byte
	URI := geocodeURL + searchText

	status, geocodingRequest, err = fasthttp.Get(geocodingRequest, URI)

	if err != nil {
		log.Printf("error while requesting coordinates from geocoding")
		return "", "", status, err
	}
	if status >= 500 {
		log.Printf("geocode service unvailable")
		return "", "", status, nil
	}
	if status != 200 {
		log.Printf("status not OK in geocoding response")
		return "", "", status, nil
	}

	c := make(map[string]json.RawMessage)

	err = json.Unmarshal(geocodingRequest, &c)
	if err != nil {
		log.Printf("error while unmarshaling request")
		return "", "", status, err
	}

	longt := string(c["longt"])
	latt := string(c["latt"])

	return longt, latt, 200, nil
}

// /weather/checkcity/:cityname
func MainWeatherHandler(ctx *fasthttp.RequestCtx) {
	//https://api.weather.gov/points/

	longt, latt, status, err := getCoordinatesAsFloat(fmt.Sprintf("%s", ctx.UserValue("cityname")))
	if err != nil {
		log.Printf("error occured while trying to get coordinates: %s", err)
		return
	}
	if status >= 500 {
		fmt.Fprintf(ctx, "Unfortunately, geocode services are unavailable at the moment. Please try again later.")
		return
	}
	if latt == 500 && longt == 500 {
		fmt.Fprintf(ctx, "Geocode service got too many requests and have not processed your search. Please try again.")
		return
	}
	//request := ctx.Request
	//request.Header.Add("User-Agent", "golang weather_checker app")

	temp, pressure, humidity, sunrise, sunset, status, err := getWeatherData(latt, longt)
	if err != nil {
		fmt.Fprintf(ctx, "We could not process your request due to unidentified issues")
		return
	}
	if status >= 500 {
		fmt.Fprintf(ctx, "Openweather service is unavailable. Try again later.")
		return
	}
	if status != 200 {
		fmt.Fprintf(ctx, "Openweather could not process our request")
		return
	}

	temp -= 272.15 //convert to Celsius

	clientSunrise := time.Unix(sunrise, 0)
	clientSunset := time.Unix(sunset, 0)

	fmt.Fprintf(ctx, "Weather data for coordinates: %f, %f\n", latt, longt)
	fmt.Fprintf(ctx, "temp: %f, pressure: %f, humidity: %f, sunrise: %s, sunset: %s", temp, pressure, humidity, clientSunrise, clientSunset)
}

func getWeatherData(lat float32, lon float32) (temp float32, pressure float32, humidity float32, sunrise int64, sunset int64, status int, err error) {
	toExclude := "minutely,hourly,daily,alerts"
	var openWeatherRequest []byte
	URI := fmt.Sprintf("%s?lat=%f&lon=%f&exclude=%s&appid=%s", openWeatherURL, lat, lon, toExclude, openWeatherAPIKey)

	status, openWeatherRequest, err = fasthttp.Get(openWeatherRequest, URI)

	if err != nil {
		log.Printf("error while requesting coordinates from geocoding")
		return 0, 0, 0, 0, 0, 500, err
	}
	if status != 200 {
		log.Printf("openweather service unvailable or wrong request")
		return 0, 0, 0, 0, 0, status, nil
	}
	unmarshaledMap1 := make(map[string]json.RawMessage)

	err = json.Unmarshal(openWeatherRequest, &unmarshaledMap1)
	if err != nil {
		log.Printf("error while unmarshaling request")
		return 0, 0, 0, 0, 0, 590, err
	}
	unmarshaledMap2 := make(map[string]json.RawMessage)

	err = json.Unmarshal(unmarshaledMap1["current"], &unmarshaledMap2)
	if err != nil {
		log.Printf("error while unmarshaling request")
		return 0, 0, 0, 0, 0, 590, err
	}

	temp, err = byteArrayToFloat(unmarshaledMap2["temp"])
	if err != nil {
		log.Printf("error while converting coordinates to float: %s", err)
		return 0, 0, 0, 0, 0, status, err
	}

	pressure, err = byteArrayToFloat(unmarshaledMap2["pressure"])
	if err != nil {
		log.Printf("error while converting coordinates to float: %s", err)
		return 0, 0, 0, 0, 0, status, err
	}
	humidity, err = byteArrayToFloat(unmarshaledMap2["humidity"])
	if err != nil {
		log.Printf("error while converting coordinates to float: %s", err)
		return 0, 0, 0, 0, 0, status, err
	}
	err = json.Unmarshal(unmarshaledMap2["sunrise"], &sunrise)
	if err != nil {
		log.Printf("error while unmarshaling request")
		return 0, 0, 0, 0, 0, 590, err
	}
	err = json.Unmarshal(unmarshaledMap2["sunset"], &sunset)
	if err != nil {
		log.Printf("error while unmarshaling request")
		return 0, 0, 0, 0, 0, 590, err
	}

	return temp, pressure, humidity, sunrise, sunset, status, nil
}
