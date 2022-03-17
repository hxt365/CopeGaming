package socket

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net"
	"testing"
)

func TestFailOnPortInUse(t *testing.T) {
	l, err := NewSocket("udp", 1234)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	defer l.(*net.UDPConn).Close()
	_, err = NewSocket("udp", 1234)
	if err == nil {
		t.Errorf("expected busy port error, but got none")
	}
}

func TestListenerPortRoll(t *testing.T) {
	l, err := NewSocketPortRoll("udp", 1234)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	defer l.(*net.UDPConn).Close()
	l2, err := NewSocketPortRoll("udp", 1234)
	if err != nil {
		t.Errorf("expected no port error, but got one")
	}
	l2.(*net.UDPConn).Close()
}

func TestNewRandomUDPListener(t *testing.T) {
	l1, err := NewRandomUDPListener()
	require.NoError(t, err)
	defer l1.Close()

	l2, err := NewRandomUDPListener()
	require.NoError(t, err)
	defer l2.Close()

	assert.NotEqual(t, l1.LocalAddr().String(), l2.LocalAddr().String())
}

func TestNewRandomTCPListener(t *testing.T) {
	l1, err := NewRandomTCPListener()
	require.NoError(t, err)
	defer l1.Close()

	l2, err := NewRandomTCPListener()
	require.NoError(t, err)
	defer l2.Close()

	defer l1.Close()
	defer l2.Close()

	assert.NotEqual(t, l1.Addr().String(), l2.Addr().String())
}

func TestExtractPort(t *testing.T) {
	port, err := ExtractPort("[::]:123")
	require.NoError(t, err)
	assert.Equal(t, 123, port)

	port, err = ExtractPort("127.0.0.1:123")
	require.NoError(t, err)
	assert.Equal(t, 123, port)
}
