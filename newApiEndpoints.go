package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

var endPointSlice []current
var curr current

func addLocationzz(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if len(endPointSlice) > 0 {
		for _, v := range endPointSlice {
			if v.CityName == params["location"] {
				w.WriteHeader(http.StatusForbidden)
			} else {
				curr.updateCurrentWeather(params["location"])
				endPointSlice = append(endPointSlice, curr)
				//	w.WriteHeader(http.StatusOK)
			}
		}
	} else {
		curr.updateCurrentWeather(params["location"])
		endPointSlice = append(endPointSlice, curr)
	}
}
func showLocation(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, v := range endPointSlice {
		if v.CityName == params["location"] {
			json.NewEncoder(w).Encode(v)
			return
		}
	}
	http.Error(w, "Location not added yet!", http.StatusForbidden)
}

func deleteCities(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for i, v := range endPointSlice {
		if params["location"] == v.CityName {
			endPointSlice = append(endPointSlice[:i], endPointSlice[i+1:]...)
			w.WriteHeader(http.StatusOK)
			return
		}
		http.Error(w, "Cannot delete the city, are you sure you added it?", http.StatusForbidden)
	}
	http.Error(w, "Cannot delete the city, are you sure you added it?", http.StatusForbidden)

}
func getLocations(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(endPointSlice)
}
