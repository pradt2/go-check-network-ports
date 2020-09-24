package main

import (
	"net"
	"strconv"
)

func server(startport uint, endport uint) {
	var i uint
	for i = startport; i < endport+1; i++ {
		go handleTcp(i)
		go handleUdp(i)
	}
	select {}
}

func handleUdp(port uint) error {
	listener, err := net.ListenUDP(UDP, &net.UDPAddr{IP: []byte{0, 0, 0, 0}, Port: int(port), Zone: ""})
	if err != nil {
		log.Error("Failed to bind to TCP port", port, err)
		return err
	}
	buf := make([]byte, len(PING))
	for {
		_, addr, err := listener.ReadFromUDP(buf)
		if err != nil {
			log.Warning("Could not read from UDP port", port, err)
			continue
		}
		if string(buf) != string(PING) {
			log.Warning("Received unexpected data. Not responding.")
			continue
		}
		_, err = listener.WriteToUDP(PONG, addr)
		if err != nil {
			log.Warning("Could not write to UDP port", port, err)
			continue
		}
	}
}

func handleTcp(port uint) error {
	listener, err := net.Listen(TCP, ":"+strconv.Itoa(int(port)))
	if err != nil {
		log.Error("Failed to bind to TCP port", port, err)
		return err
	}
	buf := make([]byte, len(PING))
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Warning("Failed to accept an incoming TCP connection", err)
			continue
		}
		_, err = conn.Read(buf)
		if string(buf) != string(PING) {
			log.Warning("Received unexpected data. Not responding.")
		}
		_, err = conn.Write(PONG)
		if err != nil {
			log.Warning("Failed to write to an incoming TCP connection", err)
		}
		err = conn.Close()
		if err != nil {
			log.Warning("Failed to close an incoming TCP connection", err)
		}
	}
}
