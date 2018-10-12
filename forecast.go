package main

// Defining structs to store recieved JSON data in.

type forecastWeatherData struct {
	Unit    string
	Lang    string
	Key     string
	baseURL string
}

type forecast struct {
	Days    int                   `json:"cnt"`
	Dt      int                   `json:"dt"`
	Message string                `json:"cod"`
	List    []forecastWeatherList `json:"list"`
	City    city                  `json:"city"`
}

type city struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type forecastWeatherList struct {
	Dt         int       `json:"dt"`
	Main       Main      `json:"main"`
	Weather    []weather `json:"weather"`
	Date       string    `json:"dt_txt"`
	Userminmax minmaxTemp
}

type Main struct {
	Temp_min float64 `json:"temp_min"`
	Temp_max float64 `json:"temp_max"`
	Avg_Temp float64 `json:"temp"`
}

type weather struct {
	Sky string `json:"description"`
}

// loads forecastWeatherData with data
func NewWeatherData(unit string, lang string) forecastWeatherData {
	data := forecastWeatherData{}
	data.Key = apiKey
	data.Lang = lang
	data.Unit = unit
	data.baseURL = "http://api.openweathermap.org/data/2.5/forecast/?q="
	return data
}
