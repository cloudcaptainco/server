package server

// TODO complete wait function.

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gorilla/mux"
)

// Config provides Server configuration
type Config struct {
	Address string
	Port    string
}

// Server type contains http server and controller methods.
type Server struct {
	router      *mux.Router
	config      Config
	server      *http.Server
	listener    *net.Listener
	exitChannel chan os.Signal
}

// HTTPController web controller object, for server.New init.
type HTTPController interface {
	ServeHTTP(writer http.ResponseWriter, request *http.Request)
}

// New returns an initialized server struct.
func New(routes *map[string]HTTPController, config Config) *Server {

	router := mux.NewRouter()
	for path, handler := range *routes {
		router.PathPrefix(path).Handler(handler)
	}

	s := Server{}
	s.router = router
	s.config = config

	return &s
}

// Serve start the webserver in a separate go routine.
func (s *Server) Serve(c chan int) {

	if s.router == nil {
		panic("router is nil, has setup been called?")
	}
	srv := &http.Server{
		Handler: s.router,
	}
	s.server = srv

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", s.config.Port))

	s.listener = &listener

	if err != nil {
		panic(err)
	}
	// Writing to channel to Wait func knows to continue.
	// See func (s *Server) Wait(c chan int) int
	c <- 1

	log.Fatal(srv.Serve(listener))
}

// Wait waits for server to recieve a Operating System call
func (s *Server) Wait(c chan int) int {
	// This is why Serve writes to the channel with a 1.
	_ = <-c

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// deref server
	server := *s.server

	exitChannel := make(chan int)
	go func() {
		for {
			sig := <-signalChannel
			log.Print("os signal received processing.")
			switch sig {
			case syscall.SIGTERM, syscall.SIGINT:
				err := server.Shutdown(context.Background())
				// Register server level errors, return code.
				if err != nil {
					log.Println(err)
					exitChannel <- 1
				} else {
					exitChannel <- 0
				}
			}
		}
	}()

	code := <-exitChannel
	return code

}

// GetPort returns the TCP port the server is using.
func (s *Server) GetPort() string {
	if s.server == nil {
		return s.config.Port
	} else {
		l := *s.listener
		return strconv.Itoa(l.Addr().(*net.TCPAddr).Port)
	}
}
