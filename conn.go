package main

import (
	"errors"
	"net"
	"time"
)

var zeroLenByteArr = make([]byte, 0, 0)

type onNewConnCallback func(network network, ip net.IP, port int, data []byte) []byte

type networkConn interface {
	bind(network network, ip net.IP, port int) error
	unbind()
}

func newConn(callback onNewConnCallback, bufferLength uint) networkConn {
	return &netConn{
		callback:     callback,
		bufferLength: bufferLength,
		isClosed:     false,
		deadline:     2 * time.Second,
		listener:     nil,
	}
}

type netConn struct {
	callback     onNewConnCallback
	bufferLength uint
	isClosed     bool
	deadline     time.Duration
	listener     net.Listener
}

func (n *netConn) bind(network network, ip net.IP, port int) error {
	listener, err := listen(network, ip, port, n.bufferLength)
	if err != nil {
		log.Warning("Failed to bind to TCP port.", port, err)
		return err
	}
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Debug("Failed to accept an incoming TCP connection.", err)
				if !n.isClosed {
					continue
				}
				break
			}
			remoteIp, err := getRemoteIp(conn.RemoteAddr())
			if err != nil {
				log.Error("Cannot get remote address.")
				break
			}
			err = conn.SetDeadline(time.Now().Add(n.deadline))
			if err != nil {
				log.Debug("Failed to set connection deadline.")
			}
			buf := make([]byte, n.bufferLength)
			readBytesCount, err := conn.Read(buf)
			if err != nil {
				log.Debug("Error while reading data.", err)
				continue
			}
			if readBytesCount != int(n.bufferLength) {
				log.Debug("Received data of unexpected length.", readBytesCount)
				continue
			}
			bytesToSend := n.callback(network, remoteIp, port, buf)
			if bytesToSend != nil {
				writtenBytesCount, err := conn.Write(bytesToSend)
				if writtenBytesCount != len(bytesToSend) {
					log.Debug("Could not send all bytes. Sending will not be resumed.")
				}
				if err != nil {
					log.Debug("Failed to write to an incoming TCP connection.", err)
				}
			}
			err = conn.Close()
			if err != nil {
				log.Debug("Failed to close an incoming TCP connection.", err)
			}
		}
	}()
	return nil
}

func getRemoteIp(addr net.Addr) (net.IP, error) {
	if tcpAddr, ok := addr.(*net.TCPAddr); ok {
		return tcpAddr.IP, nil
	} else if udpAddr, ok := addr.(*net.UDPAddr); ok {
		return udpAddr.IP, nil
	} else {
		return nil, errors.New("unknown address type")
	}

}

func (n *netConn) unbind() {
	if n.listener == nil {
		return
	}
	n.isClosed = true
	if err := n.listener.Close(); err != nil {
		log.Debug("Error while closing a connection.", err)
	}
}
