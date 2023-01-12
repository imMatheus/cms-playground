package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	store      Storage
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/", makeHTTPHandleFunc(s.handleBase)).Methods(http.MethodGet)
	router.HandleFunc("/health", makeHTTPHandleFunc(s.handleHealth)).Methods(http.MethodGet)
	router.HandleFunc("/users", makeHTTPHandleFunc(s.handleGetUsers)).Methods(http.MethodGet)
	router.HandleFunc("/stashes", makeHTTPHandleFunc(s.handleGetStashes)).Methods(http.MethodGet)

	log.Println("Starting APIServer on port ", s.listenAddr)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	http.ListenAndServe(s.listenAddr, handlers.CORS(originsOk, headersOk, methodsOk)(router))
}

func (s *APIServer) handleBase(w http.ResponseWriter, r *http.Request) error {
	return WriteJSON(w, http.StatusOK, "server is running")
}

func (s *APIServer) handleHealth(w http.ResponseWriter, r *http.Request) error {
	return WriteJSON(w, http.StatusOK, "ok")
}

func (s *APIServer) handleGetUsers(w http.ResponseWriter, r *http.Request) error {
	users, err := s.store.GetAllUsers()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, users)
}
func (s *APIServer) handleGetStashes(w http.ResponseWriter, r *http.Request) error {
	stashes, err := s.store.GetAllStashes()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, stashes)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}
