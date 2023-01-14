package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	jwt "github.com/golang-jwt/jwt/v4"
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
	router.HandleFunc("/sign-up", makeHTTPHandleFunc(s.handleSignUp)).Methods(http.MethodPost)
	router.HandleFunc("/login", makeHTTPHandleFunc(s.handleLogin)).Methods(http.MethodPost)
	router.HandleFunc("/stash", makeHTTPHandleFunc(s.handleCreateStash)).Methods(http.MethodPost)
	router.HandleFunc("/stashes", makeHTTPHandleFunc(s.handleGetStashes)).Methods(http.MethodGet)
	router.HandleFunc("/products", makeHTTPHandleFunc(s.handleGetProducts)).Methods(http.MethodGet)
	router.HandleFunc("/products/{id}", makeHTTPHandleFunc(s.handleGetProductById)).Methods(http.MethodGet)

	log.Println("Starting APIServer on port ", s.listenAddr)

	// headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	http.ListenAndServe(s.listenAddr, handlers.CORS(originsOk, methodsOk)(router))
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

func (s *APIServer) handleGetProducts(w http.ResponseWriter, r *http.Request) error {
	products, err := s.store.GetAllProducts()

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, products)
}

func (s *APIServer) handleGetProductById(w http.ResponseWriter, r *http.Request) error {
	id, err := getURLParam(r, "id")
	if err != nil {
		return err
	}

	product, err := s.store.GetProductById(id)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, product)
}

func (s *APIServer) handleSignUp(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("sigggnin up")
	createUserReq := new(CreateUserRequest)

	if err := json.NewDecoder(r.Body).Decode(createUserReq); err != nil {
		return err
	}

	fmt.Println("maddddeee 22222")

	fmt.Println("%+v", createUserReq)

	// stash := NewAccount(createAccountReq.FirstName, createAccountReq.LastName)
	user, err := s.store.CreateUser(*createUserReq)
	if err != nil {
		return err
	}

	fmt.Println("stash %+v", user)

	return WriteJSON(w, http.StatusOK, user)
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	products, err := s.store.GetAllProducts()

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, products)
}

func (s *APIServer) handleGetStashes(w http.ResponseWriter, r *http.Request) error {
	stashes, err := s.store.GetAllStashes()

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, stashes)
}

func (s *APIServer) handleCreateStash(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("hellllooo")
	createStashReq := new(CreateStashRequest)
	if err := json.NewDecoder(r.Body).Decode(createStashReq); err != nil {
		return err
	}
	fmt.Println("maddddeee 22222")

	fmt.Println("%+v", createStashReq)

	// stash := NewAccount(createAccountReq.FirstName, createAccountReq.LastName)
	stash, err := s.store.CreateStash(*createStashReq)
	if err != nil {
		return err
	}

	fmt.Println("stash %+v", stash)

	return WriteJSON(w, http.StatusOK, stash)
}

func permissionDenied(w http.ResponseWriter) {
	WriteJSON(w, http.StatusForbidden, ApiError{Error: "permission denied"})
}

func createJWT(user *User) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt": 15000,
		"userId":    user.ID,
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fmt.Println("made i t h55555")
	fmt.Println(token)

	return token.SignedString([]byte(secret))
}

func withJWTAuth(handlerFunc http.HandlerFunc, s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("calling JWT auth middleware")

		tokenString := r.Header.Get("x-jwt-token")
		fmt.Println(tokenString)
		token, err := validateJWT(tokenString)
		fmt.Println(token)
		if err != nil {
			permissionDenied(w)
			return
		}

		if !token.Valid {
			permissionDenied(w)
			return
		}

		// userID, err := getID(r)
		// if err != nil {
		// 	permissionDenied(w)
		// 	return
		// }
		// account, err := s.GetAccountByID(userID)
		// if err != nil {
		// 	permissionDenied(w)
		// 	return
		// }

		claims := token.Claims.(jwt.MapClaims)
		fmt.Println(claims)
		// if account.Number != int64(claims["accountNumber"].(float64)) {
		// 	permissionDenied(w)
		// 	return
		// }

		if err != nil {
			WriteJSON(w, http.StatusForbidden, ApiError{Error: "invalid token"})
			return
		}

		handlerFunc(w, r)
	}
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})
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
	fmt.Println("888888888888888888")
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func getURLParam(r *http.Request, name string) (int, error) {
	idStr := mux.Vars(r)[name]
	param, err := strconv.Atoi(idStr)
	if err != nil {
		return param, fmt.Errorf("invalid id given %s", idStr)
	}
	return param, nil
}
