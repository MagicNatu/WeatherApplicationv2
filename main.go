package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

//Defining key environment variable in order to access weather data from openweathermap.org
//Also defining slices, structs and map to store retrieved data

var setenv = os.Setenv("OWM_API_KEY", "ae109bfa9bd34e64691b5b467fe631eb")
var apiKey = os.Getenv("OWM_API_KEY")
var cities []string
var locMap map[string]current
var forecastMap map[string][]forecastWeatherList
var jSONforecastMap map[string][]forecast
var g current
var f forecast
var maxTemp []float64
var minTemp []float64

// Defining routes and initializing router
func main() {
	jSONforecastMap = map[string][]forecast{} //Initializing map to store JSON data from RESTapi endpoint
	endPointSlice = []current{}
	router := mux.NewRouter()
	router.HandleFunc("/forecast/", getfcData).Methods("POST")
	router.HandleFunc("/main", firstPage).Methods("POST", "GET")
	router.HandleFunc("/current/", getCdata).Methods("POST", "GET")

	//Rest api endpoints, added 16.10.2018

	router.HandleFunc("/current2/api/{location}", addLocationzz).Methods("POST")
	router.HandleFunc("/current2/api/{location}", showLocation).Methods("GET")
	router.HandleFunc("/current2/api/cities/all", getLocations).Methods("GET")
	router.HandleFunc("/current2/api/{location}", deleteCities).Methods("DELETE")

	//added 16.10 ends here

	router.HandleFunc("/forecast/{city}/{days}/{min_temp}/{max_temp}", addJSONfcdata).Methods("POST", "GET")
	router.HandleFunc("/JSONforecast", jsonFCdata).Methods("GET")
	router.HandleFunc("/forecast/for/all/days", showForecastJSONdata).Methods("GET")
	router.HandleFunc("/current/all/cities", showCurrentJSONdata).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
}

//the firstPage function represents the initial page when a user navigates to
//{host}, port 8000
func firstPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/main" { //if url != localhost:8000/main
		http.Error(w, "404 not found.", http.StatusNotFound) //throw error
		return
	}
	switch r.Method {
	case "GET": //if http method == "GET"
		http.ServeFile(w, r, "index.html")
	case "POST": //if http method == "POST"
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		location := r.FormValue("location")
		if location == "" {
			http.Redirect(w, r, "/main", http.StatusSeeOther)
			break
		}
		max := r.FormValue("max_temp")
		min := r.FormValue("min_temp")
		maxx, _ := strconv.ParseFloat(max, 64)
		minn, _ := strconv.ParseFloat(min, 64)
		cities = append(cities, location) //load slices with form values so they can be used later
		maxTemp = append(maxTemp, maxx)
		minTemp = append(minTemp, minn)
		http.Redirect(w, r, "/current/", http.StatusSeeOther)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

//Function to delete a location that no longer requires monitoring
func deleteLocation(location string) {
	for i := 0; i < len(cities); i++ {
		if cities[i] == location {
			cities = append(cities[:i], cities[i+1:]...)
			maxTemp = append(maxTemp[:i], maxTemp[i+1:]...)
			minTemp = append(minTemp[:i], minTemp[i+1:]...)
		}
	}
}

// Writing current weather data to router.
//Executes template with a loaded map datastructure (locMap)
//locMap is loaded with a Current{} struct
//Template is loaded so we can access data from structs in html
func getCdata(w http.ResponseWriter, r *http.Request) {
	locMap = map[string]current{}

	for i := 0; i < len(cities); i++ {
		g.updateCurrentWeather(cities[i])
		if g.Message == "city not found" {
			http.Error(w, "Location not found!", http.StatusNotFound)
			deleteLocation(cities[i])
			return
		}
		g.Userminmax.Usermax = maxTemp[i]
		g.Userminmax.Usermin = minTemp[i]
		locMap[cities[i]] = g //loading map with entire current struct and corresponding citydata

	}

	switch r.Method {
	case "GET":
		t, err := template.ParseFiles("AddLocation.html")
		if err != nil { // if there is an error
			log.Print("template parsing error: ", err) // log it
		}
		err = t.Execute(w, locMap)
	case "POST":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		dl := r.FormValue("delete")
		deleteLocation(dl)
		http.Redirect(w, r, "/current/", 303)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

//Writing forecast weather data to router.
//Executes template with a loaded map datastructure
func getfcData(w http.ResponseWriter, r *http.Request) {
	forecastMap = map[string][]forecastWeatherList{}
	showTimeSlice := []forecastWeatherList{}
	switch r.Method {
	case "POST":
		t, err := template.ParseFiles("Forecast.html")
		if err != nil { //catch error
			log.Print("template parsing error: ", err) //log the error
		}
		days := r.FormValue("amt_days") //Loading variables with data from form
		daysInt, _ := strconv.Atoi(days)
		time := r.FormValue("time")
		for i := 0; i < len(cities); i++ {
			f.updateForecast(cities[i], daysInt)
			for j := 0; j < f.Days; j++ {
				if strings.Contains(f.List[j].Date, time) {
					f.List[j].Userminmax.Usermax = maxTemp[i]
					f.List[j].Userminmax.Usermin = minTemp[i]
					showTimeSlice = append(showTimeSlice, f.List[j])
				}
			}
			forecastMap[cities[i]] = showTimeSlice //loading forecastMap with a slice of type []ForecastWEatherList{}
			showTimeSlice = nil                    //setting slice to nil in order to clear slice for next run
		}

		err = t.Execute(w, forecastMap) //Executing template and loading it with forecastMap
		if err != nil {                 //catch error
			log.Print("template execute error: ", err) //log the error
		}
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

//Function to update currentTemp by calling getCurrent(City)
//Receiver is a pointervalue to a Current struct
func (c *current) updateCurrentWeather(City string) {
	updatedWeather := getCurrent(City)
	(*c) = updatedWeather
}

//Similar function to update forecastweather
func (f *forecast) updateForecast(City string, Days int) {
	updatedForecast := getForecast(City, Days)
	(*f) = updatedForecast
}

// Retrieves forecast data from openweathermap.org,
//takes location and amount of days as arguments. Structs defined in forecast.go
func getForecast(location string, days int) forecast {
	var f forecastWeatherData
	f = NewWeatherData("metric", "EN")
	actualDays := (days * 8) //days represent a forecast for every 3rd hour. actualDays converts it into a whole day
	daysString := strconv.Itoa(actualDays)
	response, err := http.Get(fmt.Sprintf("%s%s&appid=%s&units=%s&lang=%s&cnt=%s", f.baseURL, location, f.Key, f.Unit, f.Lang, daysString))
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	forecastData := forecast{}
	err = json.NewDecoder(response.Body).Decode(&forecastData)
	if err != nil {
		panic(err)
	}
	return forecastData
}

// retrieves current weather data from openweathermap.org, returns a current struct.
// Structs defined in currentWeather.go.
func getCurrent(location string) current {
	var c currentWeatherData
	c = NewCurrentWeatherData("metric", "EN")
	response, err := http.Get(fmt.Sprintf("%s%s&appid=%s&units=%s&lang=%s", c.baseURL, location, c.Key, c.Unit, c.Lang))
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	currentData := current{}
	err = json.NewDecoder(response.Body).Decode(&currentData)
	if err != nil {
		panic(err)
	}
	return currentData
}

//Provides current JSON data for n locations
func showCurrentJSONdata(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(locMap)
}

//provides forecast JSON data for n locations and n<5 days.
func showForecastJSONdata(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(forecastMap)
}

//Represents endpoint for adding weather data from any location.
//accessible from localhost:8000/forecast/{city}/{days}
func addJSONfcdata(w http.ResponseWriter, r *http.Request) {
	var fc forecast
	params := mux.Vars(r)
	days, _ := strconv.Atoi(params["days"])
	minTempp, _ := strconv.ParseFloat(params["min_temp"], 64)
	maxTempp, _ := strconv.ParseFloat(params["max_temp"], 64)
	fc.updateForecast(params["city"], days)
	jSONforecastMap[params["city"]] = append(jSONforecastMap[params["city"]], fc)
	for _, item := range jSONforecastMap {
		fc.updateForecast(params["city"], days)
		for a, p := range item {
			p.List[a].Userminmax.Usermin = minTempp
			p.List[a].Userminmax.Usermax = maxTempp
		}
	}
}

//Shows the JSON data added from addJSONfcdata endpoint
func jsonFCdata(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(jSONforecastMap)
}
