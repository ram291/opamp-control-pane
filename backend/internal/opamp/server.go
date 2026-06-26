package opamp


import "log"


type Server struct {

}


func NewServer() *Server {

	log.Println("Initializing OPAMP server")

	return &Server{}

}


func (s *Server) Start() {

	log.Println("OPAMP server started")

}

