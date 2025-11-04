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
	t.Cleanup(func() {
		_ = cmd.Process.Signal(os.Interrupt)
	})
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
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
	for i := range 5 {
		var conn net.Conn
		conn, err = net.Dial("tcp", addr)
		if err == nil {
			_ = conn.Close()
			return
		}
		time.Sleep(100 * time.Millisecond * (1 << i))
	}
	t.Fatal(err)
}
