package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/cloudcaptainco/server/pkg/server"
)

type testController struct{}

func (tc testController) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	m := map[string]string{
		"test": "key",
	}
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(m)
}

func main() {

	testController := testController{}

	routes := &map[string]server.HTTPController{
		"/test": testController,
	}
	config := server.Config{"172.0.0.1", "8081"}

	s := server.New(routes, config)
	channel := make(chan int)
	go s.Serve(channel)

	os.Exit(s.Wait(channel))
}
