package main

/* The purpose here is to test all functions related to this app. TODO: Create
tests for showCurrentJSONdata, jsonFCdata, showForecastJSONdata */

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Testing getForecast method
func TestGetForecast(t *testing.T) {
	forecastData := forecast{}
	location := "Turku"
	days := 2
	if forecastData.Days == days {
		t.Errorf("Error in forecast count")
	}
	forecastData = getForecast(location, days)
	if forecastData.Message == "404" || forecastData.Message == "400" {
		t.Errorf("Forecast weatherdata struct is empty? location not found")
	}
	if (forecastData.Days) != days*7 {
		t.Errorf("Error in forecast count")
	}
}

//Testing that getCurrent method is working correctly
func TestGetCurrent(t *testing.T) {
	location := "Turku"
	currentDataLoaded := getCurrent(location)
	if currentDataLoaded.CityName != location {
		t.Errorf("City names do not match!!")
	}
	if currentDataLoaded.Message == "city not found" {
		t.Errorf("Location does not exist?")
	}
}

// Testing that updateCurrent function is working correctly
func TestUpdateCurrentWeather(t *testing.T) {
	location := "Turku"
	c := getCurrent(location)
	c2 := getCurrent(location)
	c.updateCurrentWeather("Helsinki")
	if c == c2 {
		t.Errorf("Updatecurrentweather didn't work")
	}

}

// Testing that updateForecast function is working correctly
func TestUpdateForecast(t *testing.T) {
	location := "Turku"
	days := 1
	c := getForecast(location, days)
	c2 := getForecast(location, days)
	c.updateForecast("Helsinki", 3)

	if c.City.Name == c2.City.Name {
		t.Errorf("UpdateForecast didn't work")
	}
	if c.City.ID == c2.City.ID {
		t.Errorf("UpdateForecast didn't work")
	}
}

// Testing GET function in http.get; page is loaded successfully
func TestGetIndex(t *testing.T) {
	req, err := http.NewRequest("GET", "/main", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(firstPage)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK { //Testing statuscode when page loads successfully
		t.Errorf("Fatal error with GET")
	}
	// deliberately setting false url string to check if error handling works
	req2, err2 := http.NewRequest("GET", "/maien", nil)
	if err2 != nil {
		t.Fatal(err2)
	}
	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req2)
	if status := rr2.Code; status != 404 { //Checking statuscode when page loads unsuccessfully
		t.Errorf("Fatal error with URL path %d", status)
	}
}

// Testing http.post function in firstPage.go
func TestPostIndex(t *testing.T) {
	req, err := http.NewRequest("POST", "/main", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(firstPage)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Fatal error with POST %d", status)
	}
}

// Testing whether deleteLocation works. Adding locations and temps manually.
func TestDeleteLocation(t *testing.T) {
	minTurku := 1.0
	maxTurku := 10.0
	minHelsinki := 2.0
	maxHelsinki := 18.0

	cities = append(cities, "Turku")
	maxTemp = append(maxTemp, maxTurku)
	minTemp = append(minTemp, minTurku)

	cities = append(cities, "Helsinki")
	maxTemp = append(maxTemp, maxHelsinki)
	minTemp = append(minTemp, minHelsinki)

	deleteLocation("Turku")

	for i := 0; i < len(cities)-1; i++ {
		if cities[i] == "Turku" {
			t.Errorf("deleteLocation didn't work as intended")
		}
		if minTemp[i] == minTurku {
			t.Errorf("deleteLocation didn't work as intended")
		}
		if maxTemp[i] == maxTurku {
			t.Errorf("deleteLocation didn't work as intended")
		}

	}

}
