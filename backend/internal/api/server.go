package api


import (
	"log"
	"net/http"

	"github.com/ram291/opamp-control-pane/backend/internal/opamp"
)



type Server struct {

	opamp *opamp.Server

}



func NewServer(
	opampServer *opamp.Server,
) *Server {


	return &Server{
		opamp: opampServer,
	}

}



func (s *Server) Start(address string){


	http.HandleFunc(
		"/health",
		func(w http.ResponseWriter,r *http.Request){

			w.Write([]byte("OK"))

		},
	)


	log.Println(
		"API listening on",
		address,
	)


	http.ListenAndServe(
		address,
		nil,
	)

}

