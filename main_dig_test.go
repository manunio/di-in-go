package main

import (
	"net/http"
	"testing"
)

// TODO: make it work
func TestMain(t *testing.T) {
	// wire up server
	container := BuildContainer()

	err := container.Invoke(func(server *Server) {
		server.Run()
	})

	if err != nil {
		panic(err)
	}

	//send get request
	resp, err := http.Get("http://localhost:8001/people")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.Body == nil {
		t.Fatal("empty response")
	}
}
