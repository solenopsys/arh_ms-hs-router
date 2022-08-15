package utils

import (
	"k8s.io/klog/v2"
	"strconv"
)

func ConvertMapToInt(stringMap map[string]string) map[string]uint16 {
	numberMap := make(map[string]uint16)
	for key, value := range stringMap {
		parseUint, err := strconv.ParseUint(value, 10, 16)
		if err != nil {
			klog.Error("Error parce integer: %v\n", err)
		} else {
			numberMap[key] = uint16(parseUint)
		}
	}
	return numberMap
}

func ConvertMapToString(numberMap map[string]uint16) map[string]string {
	stringMap := make(map[string]string)
	for key, value := range numberMap {
		stringMap[key] = strconv.Itoa(int(value))
	}
	return stringMap
}
