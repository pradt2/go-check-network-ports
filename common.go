package main

import "github.com/op/go-logging"

const TCP string = "tcp"
const UDP string = "udp"

var PING []byte = []byte("ping")
var PONG []byte = []byte("pong")

var log *logging.Logger = &logging.Logger{}
