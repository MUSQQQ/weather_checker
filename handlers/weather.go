package handlers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/valyala/fasthttp"
)

/*
TODO
dodac wysylanie jsona w response body przy odpytywaniu /coordinates/:cityname
nastepnie dodac obsluge odpytywnaia o pogode
*/

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

	///udane zwracanie latt i longt w formie prawilnego jsona
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
	request := ctx.Request
	request.Header.Add("User-Agent", "golang weather_checker app")

	toExclude := "minutely,hourly,daily,alerts"
	var openWeatherRequest []byte
	URI := fmt.Sprintf("%s?lat=%f&lon=%f&exclude=%s&appid=%s", openWeatherURL, latt, longt, toExclude, openWeatherAPIKey)

	status, openWeatherRequest, err = fasthttp.Get(openWeatherRequest, URI)

	if err != nil {
		log.Printf("error while requesting coordinates from geocoding")
		fmt.Fprintf(ctx, "Our server encountered problem while trying to get weather data")
		return
	}
	if status >= 500 {
		log.Printf("geocode service unvailable")
		fmt.Fprintf(ctx, "Unfortunately, openweather services are unavailable at the moment. Please try again later.")
		return
	}
	if status != 200 {
		log.Printf("status not OK in geocoding response")
		fmt.Fprintf(ctx, "Looked up term probably cannot be recognized")
		return
	}

	fmt.Fprint(ctx, string(openWeatherRequest))

}
