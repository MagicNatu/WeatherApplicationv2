package main

//stores user input data: minimum and maximum temperatures for
// a given location

type minmaxTemp struct {
	Usermax float64
	Usermin float64
}

// Defining structs to store received JSON data in.

type currentWeatherData struct {
	Unit    string
	Lang    string
	Key     string
	baseURL string
}
type current struct {
	Main       currentWeather `json:"main"`
	CityName   string         `json:"name"`
	Userminmax minmaxTemp
	Message    string `json:"message"`
}

type currentWeather struct {
	CurrentTemp float64 `json:"temp"`
}

// loads currentWeatherData struct with data
func NewCurrentWeatherData(unit string, lang string) currentWeatherData {
	data := currentWeatherData{}
	data.Key = apiKey
	data.Lang = lang
	data.Unit = unit
	data.baseURL = "http://api.openweathermap.org/data/2.5/weather/?q="
	return data
}
