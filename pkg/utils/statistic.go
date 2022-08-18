package utils

import (
	"k8s.io/klog/v2"
)

type Statistic struct {
	Values    map[string]uint16
	Increment chan string
}

func (s Statistic) GetValues() map[string]uint16 {
	return s.Values
}

func (s Statistic) UpdateLoop() {
	for {
		key := <-s.Increment
		klog.Infof("ADD STAT:", key)
		if currentValue, ok := s.Values[key]; ok {
			s.Values[key] = currentValue + 1
		} else {
			s.Values[key] = 1
		}
	}
}
