package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func getInfo(sm *map[string]uint16) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		setCorsHeaders(w)
		marshal, err := json.Marshal(sm)
		if err != nil {
			log.Println("read:", err)
		} else {
			w.Write(marshal)
		}
	}
}
