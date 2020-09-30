package main

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

type clientConfig struct {
	host            string
	portRangeStart  uint16
	portRangeEnd    uint16
	networks        []network
	waitTime        time.Duration
	parallelisation uint
}

type triple struct {
	network network
	host    string
	port    int
}

func run(config *clientConfig) (map[network][]int, error) {
	failedPorts := make(map[network][]int)
	if config.portRangeEnd < config.portRangeStart {
		return nil, errors.New(fmt.Sprintf("Port range is invalid."))
	}
	if config.networks == nil || len(config.networks) == 0 {
		return nil, errors.New("no network types provided")
	}
	channel := make(chan triple)
	wg := sync.WaitGroup{}
	for i := uint(0); i < config.parallelisation; i++ {
		wg.Add(1)
		go func() {
			for triple := range channel {
				isSuccessful := check(triple.network, triple.host, triple.port, config.waitTime)
				if isSuccessful {
					continue
				}
				failedPorts[triple.network] = append(failedPorts[triple.network], triple.port)
			}
			wg.Done()
		}()
	}
	for _, network := range config.networks {
		for port := config.portRangeStart; port <= config.portRangeEnd; port++ {
			channel <- triple{
				network: network,
				host:    config.host,
				port:    int(port),
			}
		}
	}
	close(channel)
	wg.Wait()
	return failedPorts, nil
}

func check(network network, host string, port int, waitTime time.Duration) bool {
	conn, _ := net.Dial(string(network), fmt.Sprintf("%s:%d", host, port))
	_ = conn.SetDeadline(time.Now().Add(waitTime))
	_, _ = conn.Write(PING)
	buf := make([]byte, len(PONG))
	_, err := conn.Read(buf)
	if err != nil {
		log.Info("Error while reading server response", err)
	}
	if string(buf) != string(PONG) {
		log.Info("Incorrect server response", string(buf), len(buf))
		return false
	}
	return true
}
