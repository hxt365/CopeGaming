package socket

import (
	"errors"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

const listenAttempts = 42
const udpBufferSize = 16 * 1024 * 1024

// NewSocket creates either TCP or UDP socket listener on a given port.
// The proto param supports on of these values:
// udp, udp4, udp6, tcp, tcp4, tcp6
// The function result will be either *net.UDPConn for UDPs or
// *net.TCPListener for TCPs.
func NewSocket(proto string, port int) (interface{}, error) {
	if listener, err := socket(proto, port); err != nil {
		return nil, err
	} else {
		return listener, nil
	}
}

// NewSocketPortRoll creates either TCP or UDP socket listener on the next free port.
// See: NewSocket.
func NewSocketPortRoll(proto string, port int) (listener interface{}, err error) {
	if listener, err = NewSocket(proto, port); err == nil {
		return listener, nil
	}
	if IsPortBusyError(err) {
		for i := port + 1; i < port+listenAttempts; i++ {
			listener, err := socket(proto, i)
			if err == nil {
				return listener, nil
			}
		}
		return nil, errors.New("no available ports")
	}
	return nil, err
}

func socket(proto string, port int) (interface{}, error) {
	switch proto {
	case "udp", "udp4", "udp6":
		if l, err := net.ListenUDP(proto, &net.UDPAddr{Port: port}); err == nil {
			_ = l.SetReadBuffer(udpBufferSize)
			_ = l.SetWriteBuffer(udpBufferSize)
			return l, nil
		} else {
			return nil, err
		}
	case "tcp", "tcp4", "tcp6":
		if l, err := net.ListenTCP(proto, &net.TCPAddr{Port: port}); err == nil {
			return l, nil
		} else {
			return nil, err
		}
	}
	return nil, errors.New("socket error")
}

// IsPortBusyError tests if the given error is one of
// the port busy errors.
func IsPortBusyError(err error) bool {
	if err == nil {
		return false
	}
	var eOsSyscall *os.SyscallError
	if !errors.As(err, &eOsSyscall) {
		return false
	}
	var errErrno syscall.Errno
	if !errors.As(eOsSyscall, &errErrno) {
		return false
	}
	if errErrno == syscall.EADDRINUSE {
		return true
	}
	const WSAEADDRINUSE = 10048
	if runtime.GOOS == "windows" && errErrno == WSAEADDRINUSE {
		return true
	}
	return false
}

func newRandomAvailableListener(proto string) (interface{}, error) {
	switch proto {
	case "udp", "udp4", "udp6":
		l, err := net.ListenUDP(proto, &net.UDPAddr{Port: 0})
		if err != nil {
			return nil, err
		}
		return l, nil
	case "tcp", "tcp4", "tcp6":

	}

	return nil, nil
}

func NewRandomUDPListener() (*net.UDPConn, error) {
	l, err := net.ListenUDP("udp", &net.UDPAddr{Port: 0})
	if err != nil {
		return nil, err
	}
	return l, nil
}

func NewRandomTCPListener() (*net.TCPListener, error) {
	l, err := net.ListenTCP("tcp", &net.TCPAddr{Port: 0})
	if err != nil {
		return nil, err
	}
	return l, nil
}

func ExtractPort(addr string) (int, error) {
	port, err := strconv.Atoi(addr[strings.LastIndex(addr, ":")+1:])
	if err != nil {
		return 0, err
	}
	return port, nil
}
