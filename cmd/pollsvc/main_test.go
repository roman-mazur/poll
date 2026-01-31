package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestApiHTTP(t *testing.T) {
	const testSecret = "test"
	_, addr := chooseAvailablePortTCP()
	cmd := exec.Command("go", "run", ".", "-addr", addr, "-admin-secret", testSecret)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = exec.Command("killall", "pollsvc").Run()
		t.Log("shutting down", cmd.Process.Pid)
		err := cmd.Process.Signal(os.Interrupt)
		if err != nil {
			t.Log("failed to send interrupt signal:", err)
		}
		t.Log("waiting for server shutdown")
		_, _ = cmd.Process.Wait()
		t.Log("server exited")
	})

	waitForReadiness(t, addr)
	baseUrl := "http://" + addr

	req, _ := http.NewRequest("POST", baseUrl+"/config/new?key=test-prefix", nil)
	req.Header.Set("authorization", testSecret)
	token := expectHTTPSuccess(t)(http.DefaultClient.Do(req))
	if !strings.HasPrefix(token, "test-prefix-") {
		t.Errorf("invalid token: %s", token)
	}

	t.Log(expectHTTPSuccess(t)(http.Get(baseUrl + "/config/current")))

	t.Log(expectHTTPSuccess(t)(http.Get(baseUrl + "/ping")))
}

func expectHTTPSuccess(t *testing.T) func(resp *http.Response, err error) string {
	return func(resp *http.Response, err error) string {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatal("expected 200 OK, got: ", resp.StatusCode)
		}
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		return string(data)
	}
}

func chooseAvailablePortTCP() (int, string) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	port := listener.Addr().(*net.TCPAddr).Port
	return port, fmt.Sprintf("localhost:%d", port)
}

func waitForReadiness(t *testing.T, addr string) {
	t.Helper()
	var err error
	for i := range 10 {
		var conn net.Conn
		conn, err = net.Dial("tcp", addr)
		if err == nil {
			_ = conn.Close()
			return
		}
		t.Log("failed", i+1, err)
		time.Sleep(100 * time.Millisecond * (1 << i))
	}
	t.Fatal(err)
}
