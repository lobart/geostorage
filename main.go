package main

import (
	"net/http"
	"github.com/lobart/go_geoserver.git/pkg/server"
)

func main() {
	var server server.ServerGeo
	server.InitDBConnection()
	defer server.CloseDBConnection()
	http.HandleFunc("/", server.StartPage)
	http.HandleFunc("/hello", server.Hello)
	http.HandleFunc("/kick", server.Kick)
	http.ListenAndServe(":8090", nil)
}