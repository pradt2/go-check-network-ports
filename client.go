package main

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

func client(host string, startport uint, endport uint) {
	failedTcpPorts := make([]uint, 0, 10)
	failedUdpPorts := make([]uint, 0, 10)
	var i uint
	for i = startport; i < endport+1; i++ {
		if err := checkTcp(host, i); err != nil {
			failedTcpPorts = append(failedTcpPorts, i)
		}
		if err := checkUdp(host, i); err != nil {
			failedUdpPorts = append(failedUdpPorts, i)
		}
	}
	fmt.Printf("Blocked TCP ports: %v\n", failedTcpPorts)
	fmt.Printf("Blocked UDP ports: %v\n", failedUdpPorts)
}

func checkUdp(host string, port uint) error {
	raddr, err := net.ResolveUDPAddr(UDP, host+":"+strconv.Itoa(int(port)))
	if err != nil {
		log.Error("Could not resolve host address", err)
		return err
	}
	conn, err := net.DialUDP(UDP, nil, raddr)
	if err != nil {
		log.Warning("Could not open UDP connection", err)
		return err
	}
	conn.SetDeadline(time.Now().Add(2000 * time.Millisecond))
	_, err = conn.Write(PING)
	if err != nil {
		log.Warning("Could not write to UCP connection", err)
		return err
	}
	buf := make([]byte, len(PONG))
	_, _, err = conn.ReadFromUDP(buf)
	if err != nil {
		log.Warning("Could not read from UDP connection", err)
		return err
	}
	if string(buf) != string(PONG) {
		log.Warning("Received unexpected data from UDP connection", err)
		return err
	}
	return nil
}

func checkTcp(host string, port uint) error {
	conn, err := net.Dial(TCP, host+":"+strconv.Itoa(int(port)))
	if err != nil {
		log.Warning("Could not open TCP connection on port", port, err)
		return err
	}
	bytes := make([]byte, len(PONG))
	if err := conn.SetReadDeadline(time.Now().Add(2000 * time.Millisecond)); err != nil {
		log.Warning("Failed to set TCP connection read deadline", err)
	}
	conn.SetDeadline(time.Now().Add(2000 * time.Millisecond))
	_, err = conn.Write(PING)
	if err != nil {
		log.Warning("Could not write to TCP connection", err)
		return err
	}
	_, err = conn.Read(bytes)
	if err != nil {
		log.Warning("Could not read from TCP connection", err)
		return err
	}
	if string(bytes) != string(PONG) {
		log.Warning("Received unexpected data from TCP connection")
		return err
	}
	err = conn.Close()
	if err != nil {
		log.Warning("Could not close TCP connection", err)
		return err
	}
	return nil
}
