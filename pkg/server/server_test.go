package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

var (
	ADDRESS string = "127.0.0.1"
)

type TestController struct{}

func (tc TestController) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	m := map[string]string{
		"test": "key",
	}
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(m)
}

func makeServer() *Server {
	port := "0"
	testController := TestController{}
	testMap := &map[string]HTTPController{
		"/test": testController,
	}
	config := Config{
		ADDRESS,
		port,
	}
	return New(testMap, config)
}

// Test creating a  new server object.
func TestNewServer(t *testing.T) {
	_ = makeServer()
}

func TestServerGetPort(t *testing.T) {
	s := makeServer()
	if s.GetPort() != "0" {
		t.Error("Port is not expected ':0'")
	}
	t.Logf("Server running on %s", s.GetPort())
}

// Test starting a new server, then shutting down.
func TestServerServing(t *testing.T) {
	s := makeServer()
	ch := make(chan int)
	go s.Serve(ch)
	// Read first item from channel to ensure
	<-ch
}

func TestServerResponding(t *testing.T) {
	isServerTests(t)
	s := makeServer()
	ch := make(chan int)
	go s.Serve(ch)

	<-ch

	address := fmt.Sprintf("%s:%s", ADDRESS, s.GetPort())

	t.Logf("address is %s", address)

	resp, err := http.Get(fmt.Sprintf("http://%s/test", address))

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	testMap := make(map[string]string)
	jsonErr := json.Unmarshal(body, &testMap)

	if jsonErr != nil {
		t.Error(jsonErr)
	}

	for key, val := range testMap {
		t.Logf("%s %s", key, val)
	}

	if val, ok := testMap["test"]; ok {
		if val != "key" {
			t.Error("value is not key")
		}
	}

}

func isServerTests(t *testing.T) {
	if os.Getenv("SERVERTEST") != "" {
		t.Skip("Skipping server tests")
	}
}
