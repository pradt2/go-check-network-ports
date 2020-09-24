package main

import (
	"os"
	"strconv"
)

func main() {
	argsWithoutProg := os.Args[1:]
	if argsWithoutProg[0] == "server" {
		startport, _ := strconv.Atoi(argsWithoutProg[1])
		endport, _ := strconv.Atoi(argsWithoutProg[2])
		server(uint(startport), uint(endport))
	} else {
		host := argsWithoutProg[1]
		startport, _ := strconv.Atoi(argsWithoutProg[2])
		endport, _ := strconv.Atoi(argsWithoutProg[3])
		client(host, uint(startport), uint(endport))
	}
}
