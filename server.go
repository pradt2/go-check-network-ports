package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"strings"
)

type serverConfig struct {
	portRangeStart   uint16
	portRangeEnd     uint16
	ip               net.IP
	networks         []network
	datagramSize     uint
	ignoreBindErrors bool
	tlsCert          *tls.Certificate
}

var defaultServerConfig = serverConfig{
	portRangeStart:   1,
	portRangeEnd:     65535,
	ip:               net.ParseIP("0.0.0.0"),
	networks:         []network{tcp4},
	ignoreBindErrors: false,
	datagramSize:     4,
	tlsCert:          nil,
}

type serverX interface {
	start() error
	stop()
}

type serverImpl struct {
	config    *serverConfig
	conns     []networkConn
	isStarted bool
}

func newServer(config *serverConfig) serverX {
	return &serverImpl{
		config: config,
		conns:  make([]networkConn, 0),
	}
}

var pingPong onNewConnCallback = func(network network, ip net.IP, port int, data []byte) []byte {
	if string(data) != string(PING) {
		log.Debug("Data received is not PING.")
		return nil
	}
	log.Infof("Received %s PING from %s at port %d", strings.ToUpper(string(network)), ip.String(), port)
	return PONG
}

func (s *serverImpl) start() error {
	if s.isStarted {
		return errors.New("server has been started. create a new instance if needed")
	}
	if s.config.portRangeEnd < s.config.portRangeStart {
		return errors.New(fmt.Sprintf("Port range is invalid."))
	}
	if s.config.networks == nil || len(s.config.networks) == 0 {
		return errors.New("no network types provided")
	}
	for port := s.config.portRangeStart; port <= s.config.portRangeEnd; port++ {
		for _, networkType := range s.config.networks {
			conn := newConn(pingPong, s.config.datagramSize)
			err := conn.bind(networkType, s.config.ip, int(port))
			if err != nil {
				log.Warningf("Could not bind to %s on port %d", networkType, port)
				if s.config.ignoreBindErrors {
					continue
				} else {
					return errors.New("could not bind to required port")
				}
			}
			s.isStarted = true
			s.conns = append(s.conns)
		}
	}
	if len(s.conns) == 0 {
		return errors.New("could not bind to any port")
	}
	return nil
}

func (s *serverImpl) stop() {
	for _, conn := range s.conns {
		conn.unbind()
	}
}
