package main

import (
	"encoding/json"
	"net/http"
)

type Secret struct {
	Data string `json:"data"`
}

var secret = Secret{
	Data: "TOP SECRET TOKEN 12345",
}

// помилка
func corsHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

/*
	виправлення

	func corsHeaders(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if origin == "http://trusted.local" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
	}
*/
func secretHandler(w http.ResponseWriter, r *http.Request) {
	corsHeaders(w)
	json.NewEncoder(w).Encode(secret)
}

func main() {
	http.HandleFunc("/secret", secretHandler)

	println("CORS server running on :8081")
	http.ListenAndServe(":8081", nil)
}
