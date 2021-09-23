package server

import (
	"fmt"
	"github.com/gorilla/schema"
	"log"
	"net/http"
	"github.com/lobart/go_geoserver.git/pkg/models"
	"github.com/lobart/go_geoserver.git/pkg/db_driver_fabric"
	"encoding/json"
)




var decoder  = schema.NewDecoder()

func StartPage(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello!\n")
	fmt.Fprintf(w, "Your query should be like: \n .../kick?companyName='blabla'&kickName='blabla'&longitude=10.0&latitude=10.0&speed=100.0&status='BUSY'\n")
}

func Hello(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "hello\n")
}

func Kick(w http.ResponseWriter, req *http.Request) {
	var kick models.KickConfig

	err := decoder.Decode(&kick, req.URL.Query())

	answer, _ := json.Marshal(kick)
	fmt.Fprintf(w, "Your kick parameters is:\n")
	fmt.Fprintf(w, string(answer) + "\n")
	if err != nil {
		log.Println("Error in GET parameters : ", err)
	} else {
		log.Println("GET parameters : ", kick)

	}
	 dbCreator := db_driver_fabric.DriverCreator{}
	 driver := dbCreator.CreateDriver()
	 driver.Push(kick)
	fmt.Fprintf(w, "Your kick is pushed!\n")
}

