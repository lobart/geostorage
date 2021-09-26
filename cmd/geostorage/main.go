package main

import (
	"github.com/lobart/go_geoserver.git/internal/http"
)

func main() {
	s := http.ServerGeo{}
	s.InitDBConnection()
	defer s.CloseDBConnection()
	s.StartServer()
}