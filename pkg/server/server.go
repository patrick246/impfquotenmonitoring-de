package server

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/patrick246/impfquotenmonitoring-de/pkg/persistence"
	"log"
	"net/http"
)

type Server struct {
	addr    string
	router  *httprouter.Router
	storage persistence.StorageService
}

var endRequest = struct{}{}

func NewServer(addr string, storage persistence.StorageService) (*Server, error) {
	router := httprouter.New()
	server := &Server{
		addr:    addr,
		router:  router,
		storage: storage,
	}
	server.registerRoutes()
	return server, nil
}

func (s *Server) registerRoutes() {
	s.router.GET("/api/vaccinestats/:year/:month", s.HandleMonthRequest)
	s.router.GET("/readyz", s.readiness)
}

func (s *Server) ListenAndServe() error {
	log.Printf("starting listener on %s", s.addr)
	return http.ListenAndServe(s.addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != endRequest && r != nil {
				panic(r)
			}
		}()

		s.router.ServeHTTP(w, r)
	}))
}

func (s *Server) error(w http.ResponseWriter, r *http.Request, message string, code int) {
	log.Printf("error during request, uri=%s, message=%s, code=%d", r.RequestURI, message, code)

	bytes, err := json.Marshal(struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	}{Message: message, Code: code})

	if err != nil {
		log.Printf("error sending request error, message=%s, code=%d", message, code)
	}

	_, _ = w.Write(bytes)
	panic(endRequest)
}

func (s *Server) readiness(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	_, _ = w.Write([]byte("Ready\n"))
}
