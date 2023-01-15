package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const SESSION_NAME string = "matusSessionId"
const userIDKey string = "userID"

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
	router.HandleFunc("/me", makeHTTPHandleFunc(s.handleGetMe)).Methods(http.MethodGet)

	router.HandleFunc("/users", makeHTTPHandleFunc(s.handleGetUsers)).Methods(http.MethodGet)
	router.HandleFunc("/sign-up", makeHTTPHandleFunc(s.handleSignUp)).Methods(http.MethodPost)
	router.HandleFunc("/login", makeHTTPHandleFunc(s.handleLogin)).Methods(http.MethodPost)
	router.HandleFunc("/stash", makeHTTPHandleFunc(s.handleCreateStash)).Methods(http.MethodPost)
	router.HandleFunc("/stashes", makeHTTPHandleFunc(s.handleGetStashes)).Methods(http.MethodGet)
	router.HandleFunc("/products", makeHTTPHandleFunc(s.handleGetProducts)).Methods(http.MethodGet)
	router.HandleFunc("/products/{id}", makeHTTPHandleFunc(s.handleGetProductById)).Methods(http.MethodGet)

	log.Println("Starting APIServer on port ", s.listenAddr)

	// Create a new CORS handler
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:5173"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowCredentials(),
	)

	http.ListenAndServe(s.listenAddr, corsHandler(router))
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

func (s *APIServer) handleGetMe(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("jellooo")
	// Get the session ID from the cookies
	cookie, err := r.Cookie(SESSION_NAME)
	if err != nil {
		fmt.Println("error lol")
		fmt.Println(cookie)

		return permissionDenied(w)
	}
	fmt.Println(cookie)

	nonce, err := r.Cookie("session-nonce")
	if err != nil {
		fmt.Println("error lol")
		fmt.Println(cookie)

		return permissionDenied(w)
	}
	fmt.Println("::::::::::::.::::::::::")
	fmt.Println(nonce)
	nonceEncrypted, _ := base64.URLEncoding.DecodeString(nonce.Value)

	// nonceEncrypted := []byte(nonce.Value)
	fmt.Println(nonce.Value)
	fmt.Println(nonceEncrypted)
	// fmt.Println([]byte(nonceEncrypted))
	fmt.Println("::::::::::::.::::::::::")
	// Decrypt the session ID
	sessionIdEncrypted, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		fmt.Println("could not decrypt session id")
		return permissionDenied(w)
	}

	fmt.Println("session id909090+: ", sessionIdEncrypted)

	secretKey := []byte(os.Getenv("AUTH_SECRET"))
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		fmt.Println("made it to line 130")
		fmt.Println(err)
		return permissionDenied(w)
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println("made it to line 136")
		fmt.Println(err)
		return permissionDenied(w)
	}

	nonceSize := aesgcm.NonceSize()

	if len(sessionIdEncrypted) < nonceSize {
		return permissionDenied(w)
	}

	sessionIdEncrypted = sessionIdEncrypted[nonceSize:]
	fmt.Println("made it lenght of shit")
	fmt.Println(len(sessionIdEncrypted))
	fmt.Println(nonceSize)
	fmt.Println("Made it to line 145")
	fmt.Println(nonceEncrypted)
	fmt.Println(sessionIdEncrypted)

	sessionId, err := aesgcm.Open(nil, nonceEncrypted, sessionIdEncrypted, nil)

	fmt.Println("so here is the sessionId at last")
	fmt.Println(sessionId)
	fmt.Println(string(sessionId))

	if err != nil {
		fmt.Println("Made it to line 151")
		fmt.Println(err)
		return permissionDenied(w)
	}

	fmt.Println("sessionId fr this time: ", sessionId)

	fmt.Println("got to last line ffs")

	return WriteJSON(w, http.StatusOK, "cookie")
}

func (s *APIServer) handleSignUp(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("sigggnin up")
	createUserReq := new(CreateUserRequest)

	if err := json.NewDecoder(r.Body).Decode(createUserReq); err != nil {
		return err
	}

	fmt.Println("7000000")

	fmt.Println("%+v", createUserReq)

	// stash := NewAccount(createAccountReq.FirstName, createAccountReq.LastName)
	user, err := s.store.CreateUser(*createUserReq)
	if err != nil {
		return err
	}
	sessionId := make([]byte, 16)
	_, err = rand.Read(sessionId)
	if err != nil {
		panic(err)
	}

	fmt.Println("secret key and shiiit")
	secretKey := []byte(os.Getenv("AUTH_SECRET"))
	fmt.Println(secretKey)

	// Encrypt the session ID using AES
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		panic(err)
	}

	fmt.Println("seeewwwyyy")
	nonce := make([]byte, 12)
	_, err = rand.Read(nonce)
	if err != nil {
		panic(err)
	}

	fmt.Println("soooo, lest get a user init", nonce)

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}
	sessionIdEncrypted := aesgcm.Seal(nil, nonce, sessionId, nil)
	fmt.Println("??????????????")
	fmt.Println(sessionIdEncrypted)
	fmt.Println(nonce)
	err = s.store.CreateSession(int64(user.ID), sessionIdEncrypted)
	if err != nil {
		panic(err)
	}

	fmt.Println("user that was created %+v", user)
	fmt.Println(base64.URLEncoding.EncodeToString(sessionIdEncrypted))
	http.SetCookie(w, &http.Cookie{
		Name:     SESSION_NAME,
		Value:    base64.URLEncoding.EncodeToString(sessionIdEncrypted),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	fmt.Println("about to set cookie for nonce")
	fmt.Println(nonce)
	fmt.Println(base64.URLEncoding.EncodeToString(nonce))

	http.SetCookie(w, &http.Cookie{
		Name:     "session-nonce",
		Value:    base64.URLEncoding.EncodeToString(nonce),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	return WriteJSON(w, http.StatusOK, user)
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	products, err := s.store.GetAllProducts()

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, products)
}

func (s *APIServer) validateSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get the session ID from the cookies
		cookie, err := r.Cookie(SESSION_NAME)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Decrypt the session ID
		sessionIdEncrypted, err := base64.URLEncoding.DecodeString(cookie.Value)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		secretKey := []byte(os.Getenv("AUTH_SECRET"))
		block, err := aes.NewCipher(secretKey)
		if err != nil {
			panic(err)
		}

		aesgcm, err := cipher.NewGCM(block)
		if err != nil {
			panic(err)
		}

		nonceSize := aesgcm.NonceSize()
		if len(sessionIdEncrypted) < nonceSize {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		nonce, sessionIdEncrypted := sessionIdEncrypted[:nonceSize], sessionIdEncrypted[nonceSize:]
		sessionId, err := aesgcm.Open(nil, nonce, sessionIdEncrypted, nil)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Lookup the session in the database
		userId, err := s.store.GetSessionById(sessionId)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Add the user ID to the request context
		ctx := context.WithValue(r.Context(), userIDKey, userId)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
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

func permissionDenied(w http.ResponseWriter) error {
	return WriteJSON(w, http.StatusForbidden, ApiError{Error: "permission denied"})
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
