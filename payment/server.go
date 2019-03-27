package payment

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Server - HTTP server for handling payment requests
type Server struct {
	*Handler
	Router *mux.Router
}

// NewServer - creates and initializes a new HTTP server
func NewServer(handler *Handler) *Server {
	server := &Server{Handler: handler, Router: mux.NewRouter()}
	server.configRouter()

	return server
}

// Run - runs the server on its router
func (s *Server) Run(host string) {
	loggedRouter := handlers.LoggingHandler(os.Stdout, s.Router)
	log.Fatal(http.ListenAndServe(host, loggedRouter))
}

func (s *Server) configRouter() {
	s.get("/v1/payments", s.GetAllPayments)
	s.get("/v1/payments/{paymentID}", s.GetPaymentByID)
	s.post("/v1/payments", s.AddPayment)
	s.delete("/v1/payments/{paymentID}", s.DeletePayment)
	s.put("/v1/payments/{paymentID}", s.UpdatePayment)
}

func (s *Server) get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	s.Router.HandleFunc(path, f).Methods(http.MethodGet)
}

func (s *Server) post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	s.Router.HandleFunc(path, f).Methods(http.MethodPost)
}

func (s *Server) put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	s.Router.HandleFunc(path, f).Methods(http.MethodPut)
}

func (s *Server) delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	s.Router.HandleFunc(path, f).Methods(http.MethodDelete)
}
