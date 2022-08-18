package fortest

import (
	"net/http"
)

func getServices() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func auth() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func RunServer() {
	http.HandleFunc("/services", getServices())
	http.HandleFunc("/auth", auth())
}
