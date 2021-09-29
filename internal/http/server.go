package http

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/schema"
	"github.com/lobart/go_geoserver.git/internal/db"
	"github.com/lobart/go_geoserver.git/internal/models"
	"github.com/lobart/go_geoserver.git/internal/pubsub"
	"log"
	"net/http"
)


type ServerGeo struct {
	Driver db.DriverDB
	Ps *pubsub.Pubs
}

var decoder  = schema.NewDecoder()

func (s ServerGeo) New() (*ServerGeo, error){
	var err error
	s.Ps, err = pubsub.Pubs{}.New()
	if err!=nil{
		return &s, err
	}
	return &s, nil
}

func StartPage(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello!\n")
	_, err := fmt.Fprintf(w, "Your query should be like: \n .../kick?companyName='blabla'&kickName='blabla'&longitude=10.0&latitude=10.0&speed=100.0&status='BUSY'\n")
	if err!=nil{
		fmt.Fprintf(w, err.Error())
	}

}

func Hello(w http.ResponseWriter, req *http.Request) {
	_, err := fmt.Fprintf(w, "hello")
	if err!=nil{
		fmt.Fprintf(w, err.Error())
	}
}

func  (s *ServerGeo) Kick(w http.ResponseWriter, req *http.Request) {
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
	fmt.Println("pubbub is ",s.Ps)
	go s.Ps.Publish("kick", kick)
	stErr := fmt.Sprintf("Error %s", err)
	if err!=nil{
		fmt.Fprintf(w, stErr)
	}
	fmt.Fprintf(w, "Your kick is pushed!\n")
}

func (s *ServerGeo) InitDBConnection() (*ServerGeo, error) {
	var err error
	s.Driver, err = db.New(s.Ps)
	return s, err
}


func (s *ServerGeo) CloseDBConnection() (*ServerGeo, error) {
	err := s.Driver.Close()
	return s, err
}



func (s *ServerGeo) StartServer() error {
	http.HandleFunc("/", StartPage)
	http.HandleFunc("/hello", Hello)
	http.HandleFunc("/kick", s.Kick)
	if err := http.ListenAndServe(":8090", nil); err!=nil {
		return err
	}
	return nil
}