package main

import (
	"github.com/lobart/go_geoserver.git/internal/http"
)

func main() {
	s, err := http.ServerGeo{}.New()
	if err != nil {
		panic(err)
	}
	s.InitDBConnection()
	defer s.CloseDBConnection()
	err = s.StartServer()
	if err != nil {
		panic(err)
	}
}