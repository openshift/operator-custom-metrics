package metrics

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"testing"
	"time"
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
	return strconv.Itoa(l.Addr().(*net.TCPAddr).Port), nil
}

func TestCheckOpenTCPPortUsed(t *testing.T) {
	port, err := getFreePort()
	if err != nil {
		t.Errorf("failed to get free port on host: %v", err)
		return
	}
	freePort := fmt.Sprintf(":%s", port)
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
	port, err := getFreePort()
	if err != nil {
		t.Errorf("failed to get free port")
		return
	}
	freePort := fmt.Sprintf(":%s", port)
	if !isPortFree(freePort) {
		t.Errorf("port %s is actually free", freePort)
	}
}

func getMetricsResponse(metricsAddr string) (int, error) {
	//#nosec G107 : the variable part of the URL is the port which is provided by Kernel at runtime for the test so cannot be const
	resp, err := http.Get(metricsAddr)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}

func TestStartMetricsSuccess(t *testing.T) {
	freePort, err := getFreePort()
	if err != nil {
		t.Errorf("failed to get free port on host: %v", err)
		return
	}
	err = StartMetrics(metricsConfig{
		metricsPort: freePort,
		metricsPath: "/metrics",
	})
	if err != nil {
		t.Errorf("failed to start metrics listener: %v", err)
		return
	}
	metricsAddr := fmt.Sprintf("http://localhost:%s/metrics", freePort)
	var success bool
	for i := 0; i < 12; i++ {
		time.Sleep(time.Second * 5)
		statusCode, err := getMetricsResponse(metricsAddr)
		if err != nil || statusCode != http.StatusOK {
			continue
		}
		success = true
		break
	}
	if !success {
		t.Error("failed to get fetch metrics response successfully")
	}
}

func TestStartMetricsFailure(t *testing.T) {
	freePort, err := getFreePort()
	if err != nil {
		t.Errorf("failed to get free port on host: %v", err)
		return
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", freePort))
	if err != nil {
		t.Errorf("failed to listen on port %s: %v", freePort, err)
		return
	}
	defer listener.Close()
	err = StartMetrics(metricsConfig{
		metricsPort: freePort,
		metricsPath: "/metrics",
	})
	if err == nil {
		t.Errorf("started metrics listener even when port is bound")
		return
	}
	if err.Error() != fmt.Sprintf("port %s is not free", freePort) {
		t.Errorf("unexpected error when starting metrics listener")
	}
}
