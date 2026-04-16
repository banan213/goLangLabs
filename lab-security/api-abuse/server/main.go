package main

import (
	"fmt"
	"net/http"
)

var credits = 0

// помилка
func rewardHandler(w http.ResponseWriter, r *http.Request) {
	credits += 5
	fmt.Fprintf(w, "Credits: %d\n", credits)
}

/*
	виправлення

func rewardHandler(w http.ResponseWriter, r *http.Request) {

	now := time.Now()

	if now.Sub(lastRequestTime) < 3*time.Second {
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
		return
	}

	lastRequestTime = now

	credits += 5
	fmt.Fprintf(w, "Credits: %d\n", credits)
*/
func main() {
	http.HandleFunc("/reward", rewardHandler)

	fmt.Println("API Abuse server on :8083")
	http.ListenAndServe(":8083", nil)
}
