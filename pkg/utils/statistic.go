package utils

import (
	"k8s.io/klog/v2"
)

func NewStatistic() *Statistic {
	s := &Statistic{values: make(map[string]uint16), Increment: make(chan string, 256)}
	go s.UpdateLoop()
	return s
}

type Statistic struct {
	values    map[string]uint16
	Increment chan string
}

func (s Statistic) Values() map[string]uint16 {
	return s.values
}

func (s Statistic) UpdateLoop() {
	for {
		key := <-s.Increment
		klog.Infof("ADD STAT:", key)
		if currentValue, ok := s.values[key]; ok {
			s.values[key] = currentValue + 1
		} else {
			s.values[key] = 1
		}
	}
}
