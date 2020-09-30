package main

import (
	"errors"
	"fmt"
	"net"
)

type network string

const tcp4 network = "tcp4"
const tcp6 network = "tcp6"
const udp4 network = "udp4"
const udp6 network = "udp6"

func listen(network network, ip net.IP, port int, datagramSize uint) (net.Listener, error) {
	switch network {
	case tcp4, tcp6:
		return net.ListenTCP(string(network), &net.TCPAddr{
			IP:   ip,
			Port: port,
			Zone: "",
		})
	case udp4, udp6:
		udpConn, err := net.ListenUDP(string(network), &net.UDPAddr{
			IP:   ip,
			Port: port,
			Zone: "",
		})
		if err != nil {
			return nil, err
		}
		return newUdpListener(udpConn, datagramSize), nil
	default:
		log.Error("Unknown network type.", network)
		return nil, errors.New(fmt.Sprintf("Unknown network type: %s.", network))
	}
}
