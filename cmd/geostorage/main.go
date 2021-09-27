package main

import (
	"github.com/lobart/go_geoserver.git/internal/http"
	"github.com/lobart/go_geoserver.git/internal/pubsub"
)

func main() {
	pS := pubsub.Pubsub{}.New()
	s := http.ServerGeo{Ps : pS}
	s.InitDBConnection()
	defer s.CloseDBConnection()
	s.StartServer()
}