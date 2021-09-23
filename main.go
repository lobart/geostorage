package go_geoserver

import (
	"net/http"
	"github.com/lobart/go_geoserver.git/pkg/server"
)

func main() {
	http.HandleFunc("/", server.StartPage)
	http.HandleFunc("/hello", server.Hello)
	http.HandleFunc("/kick", server.Kick)
	http.ListenAndServe(":8090", nil)
}