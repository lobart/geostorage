package http

import (
	"encoding/json"
	"fmt"
	"github.com/lobart/go_geoserver.git/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
)


func TestServerGeo_Hello(t *testing.T) {
	req, err := http.NewRequest("GET", "/hello", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	Hello(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := "hello"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestServerGeo_StartPage(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	StartPage(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := "Hello!\nYour query should be like: \n .../kick?companyName='blabla'&kickName='blabla'&longitude=10.0&latitude=10.0&speed=100.0&status='BUSY'\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestServerGeo_New(t *testing.T) {
	s, err := ServerGeo{}.New()
	if (s==nil) || (err!=nil) {
		t.Errorf("Constructor of server returned nil pointer")
	}
}

func TestServerGeo_InitDBConnection(t *testing.T) {
	s, err := ServerGeo{}.New()
	_, err = s.InitDBConnection()
	if err!=nil {
		t.Errorf("Init DB connection error: %v", err.Error())
	}
}

func TestServerGeo_CloseDBConnection(t *testing.T) {
	s, err := ServerGeo{}.New()
	s.InitDBConnection()
	fmt.Print("Connected")
	s, err = s.CloseDBConnection()
	if err!=nil {
		t.Errorf("Close DB connection error: %v", err.Error())
	}
}

func TestServerGeo_Kick(t *testing.T) {
	s, err := ServerGeo{}.New()
	if err != nil {
		panic(err)
	}
	s.InitDBConnection()
	defer s.CloseDBConnection()


	req, err := http.NewRequest("GET", "/kick?companyName=%27blabla%27&kickName=%27blabla%27&longitude=10.5&latitude=10.5&speed=100.0&status=%27BUSY%27", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	s.Kick(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	kick:= models.KickConfig{"'blabla'","'blabla'",10.5,10.5,100,"'BUSY'"}
	ans, _ := json.Marshal(kick)
	expected := "Your kick parameters is:\n"+ string(ans) +"\nYour kick is pushed!\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}


}

/*func TestServerGeo_StartServer(t *testing.T) {
	s, _ := ServerGeo{}.New()
	err := s.StartServer()
	if err != nil {
		t.Errorf("Server not started with error: %v",err.Error())
	}
}*/