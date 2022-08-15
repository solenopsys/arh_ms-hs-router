package utils

import (
	"encoding/json"
	"k8s.io/klog/v2"
	"net/http"
)

func SetCorsHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func PrintMap(s string, m map[string]uint16) {
	bs, _ := json.Marshal(m)

	klog.Infof(s, string(bs))
}

func RemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}
