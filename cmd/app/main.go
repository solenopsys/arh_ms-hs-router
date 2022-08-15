package main

import (
	"solenopsys.org/zmq_router/internal/core"
	"solenopsys.org/zmq_router/pkg/utils"
)

var Stat *utils.Statistic

func main() { //todo сделать интерфейсы
	core.Run()
}
