package main

import (
	"io"
	"net"
	"time"
)

type udpConn struct {
	udpConn    *net.UDPConn
	remoteAddr *net.UDPAddr
	datagram   []byte
}

func newUdpConn(conn *net.UDPConn, remoteAddr *net.UDPAddr, datagram []byte) net.Conn {
	return &udpConn{
		udpConn:    conn,
		remoteAddr: remoteAddr,
		datagram:   datagram,
	}
}

func (m *udpConn) Read(b []byte) (n int, err error) {
	if len(m.datagram) == 0 {
		return 0, io.EOF
	}
	n = copy(b, m.datagram)
	return n, nil
}

func (m *udpConn) Write(b []byte) (n int, err error) {
	return m.udpConn.WriteToUDP(b, m.remoteAddr)
}

func (m *udpConn) Close() error {
	return nil
}

func (m *udpConn) LocalAddr() net.Addr {
	return m.udpConn.LocalAddr()
}

func (m *udpConn) RemoteAddr() net.Addr {
	return m.remoteAddr
}

func (m *udpConn) SetDeadline(t time.Time) error {
	return m.SetWriteDeadline(t)
}

func (m *udpConn) SetReadDeadline(t time.Time) error {
	return m.udpConn.SetReadDeadline(t)
}

func (m *udpConn) SetWriteDeadline(t time.Time) error {
	return m.udpConn.SetWriteDeadline(t)
}

type udpListener struct {
	*net.UDPConn
	datagramSize uint
}

func newUdpListener(conn *net.UDPConn, datagramSize uint) net.Listener {
	return &udpListener{
		UDPConn:      conn,
		datagramSize: datagramSize,
	}
}

func (u *udpListener) Accept() (net.Conn, error) {
	b := make([]byte, u.datagramSize)
	_, remoteAddr, err := u.UDPConn.ReadFromUDP(b)
	if err != nil {
		return nil, err
	}
	return newUdpConn(u.UDPConn, remoteAddr, b), nil
}

func (u *udpListener) Addr() net.Addr {
	return u.UDPConn.LocalAddr()
}

func (u *udpListener) Close() error {
	return u.UDPConn.Close()
}
