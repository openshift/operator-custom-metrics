package metrics

import (
	"fmt"
	"net"
	"testing"
)

// GetFreePort asks the kernel for a free open port that is ready to use.
func getFreePort() (string, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return "", err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return "", err
	}
	defer l.Close()
	return fmt.Sprintf(":%d", l.Addr().(*net.TCPAddr).Port), nil
}

func TestCheckOpenTCPPortUsed(t *testing.T) {
	freePort, err := getFreePort()
	if err != nil {
		t.Errorf("failed to get free port on host: %v", err)
		return
	}
	_, err = net.Listen("tcp", freePort)
	if err != nil {
		t.Errorf("failed to listen on free port %s: %v", freePort, err)
		return
	}
	if isPortFree(freePort) {
		t.Errorf("address %s is is not actually free", freePort)
	}
}

func TestCheckOpenTCPPortFree(t *testing.T) {
	freePort, err := getFreePort()
	if err != nil {
		t.Errorf("failed to get free port")
		return
	}
	if !isPortFree(freePort) {
		t.Errorf("port %s is actually free", freePort)
	}
}
