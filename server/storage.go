package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type Storage interface {
	GetAllUsers() ([]*User, error)
	GetUserById(id int64) (*User, error)
	CreateUser(user CreateUserRequest) (*User, error)
	CreateSession(userId int64, sessionId []byte) error
	GetSessionById(sessionId []byte) (int64, error)
	GetAllStashes() ([]*Stash, error)
	CreateStash(stash CreateStashRequest) (*CreateStashResponse, error)
	GetAllProducts() ([]*Product, error)
	GetProductById(id int) (*Product, error)
}

// A Handler is an HTTP API server handler.
type Handler struct {
	db *sql.DB
}

func NewStore() (*Handler, error) {
	// Load environment variables from file.
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load environment variables: %v", err)
	}

	// Connect to PlanetScale database using DSN environment variable.
	db, err := sql.Open("mysql", os.Getenv("DSN"))

	if err != nil {
		log.Fatalf("failed to connect to PlanetScale: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to connect to PlanetScale: %v", err)
		return nil, err
	}

	log.Println("Successfully connected to PlanetScale!")

	return &Handler{
		db: db,
	}, nil
}

func (h *Handler) createStashesTable() error {
	query := `CREATE TABLE if not exists stash (
		id int PRIMARY KEY AUTO_INCREMENT,
		name varchar(255) NOT NULL,
		location varchar(255) NOT NULL,
		createdAt datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
		updatedAt datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3)
	);`

	_, err := h.db.Exec(query)
	return err
}
func (h *Handler) createSessionsTable() error {
	query := `CREATE TABLE if not exists session (
		id int PRIMARY KEY AUTO_INCREMENT,
		userId int NOT NULL,
		sessionId BLOB NOT NULL,
		createdAt datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
		updatedAt datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3)
	);`

	_, err := h.db.Exec(query)
	return err
}
func (h *Handler) createUsersTable() error {
	query := `CREATE TABLE if not exists user (
		id int PRIMARY KEY AUTO_INCREMENT,
		name varchar(255) NOT NULL,
		image varchar(255) NOT NULL,
		email varchar(255) NOT NULL UNIQUE,
		password varchar(255) NOT NULL,
		createdAt datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
		updatedAt datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3)
	);`

	_, err := h.db.Exec(query)

	return err
}

func (h *Handler) createProductsTable() error {
	query := `CREATE TABLE if not exists product (
		id int PRIMARY KEY AUTO_INCREMENT,
		name varchar(255) NOT NULL,
		price double NOT NULL,
		stashId int NOT NULL,
		createdAt datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
		updatedAt datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3)
	);`

	_, err := h.db.Exec(query)
	return err
}

func (h *Handler) GetAllUsers() ([]*User, error) {
	rows, err := h.db.Query("select * from user")
	if err != nil {
		return nil, err
	}

	users := []*User{}
	for rows.Next() {
		user, err := scanIntoUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (h *Handler) GetUserById(id int64) (*User, error) {
	rows, err := h.db.Query("select id, name, email, image, createdAt, updatedAt from user where id = ?", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoUser(rows)
	}

	return nil, fmt.Errorf("account %d not found", id)
}

func (h *Handler) CreateUser(user CreateUserRequest) (*User, error) {
	query := `INSERT INTO user
	(name, image, email, password)
	values (?, ?, ?, ?)`

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err.Error())
	}

	rows, err := h.db.Exec(
		query,
		user.Name,
		"https://avatars.githubusercontent.com/u/77362975?v=4",
		user.Email,
		hashedPassword,
	)

	if err != nil {
		return nil, err
	}

	userId, err := rows.LastInsertId()

	if err != nil {
		return nil, err
	}

	finalUser, err := h.GetUserById(userId)

	if err != nil {
		return nil, err
	}

	return finalUser, nil
}

func scanIntoUser(rows *sql.Rows) (*User, error) {
	user := new(User)
	err := rows.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Image,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	return user, err
}

func (h *Handler) CreateSession(userId int64, sessionId []byte) error {
	_, err := h.db.Exec("INSERT INTO session (userId, sessionId) VALUES (?, ?)", userId, sessionId)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) GetSessionById(sessionId []byte) (int64, error) {
	rows, err := h.db.Query("SELECT userId FROM session WHERE sessionId = ?", sessionId)
	if err != nil {
		return 0, err
	}

	type UserId struct {
		UserId int64 `json:"userId"`
	}

	for rows.Next() {
		userId := new(UserId)
		err := rows.Scan(
			&userId.UserId,
		)

		if err == nil {
			return userId.UserId, nil
		}
	}

	return 0, fmt.Errorf("session %d not found", sessionId)
}

func scanIntoSession(rows *sql.Rows) (*Session, error) {
	session := new(Session)
	err := rows.Scan(
		&session.ID,
		&session.UserId,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	return session, err
}

func (h *Handler) GetAllStashes() ([]*Stash, error) {
	rows, err := h.db.Query("select * from stash")
	if err != nil {
		return nil, err
	}

	stashes := []*Stash{}
	for rows.Next() {
		stash, err := scanIntoStash(rows)
		if err != nil {
			return nil, err
		}
		stashes = append(stashes, stash)
	}

	return stashes, nil
}

func scanIntoStash(rows *sql.Rows) (*Stash, error) {
	stash := new(Stash)
	err := rows.Scan(
		&stash.ID,
		&stash.Name,
		&stash.Location,
		&stash.CreatedAt,
		&stash.UpdatedAt,
	)

	return stash, err
}

type CreateStashResponse struct {
	ID int64 `json:"id"`
}

func (h *Handler) CreateStash(stash CreateStashRequest) (*CreateStashResponse, error) {
	query := `INSERT INTO stash
	(name, location)
	values (?, ?)`

	resp, err := h.db.Exec(
		query,
		stash.Name, stash.Location,
	)

	if err != nil {
		return nil, err
	}

	id, _ := resp.LastInsertId()
	fmt.Printf("%+v\n", resp)

	res := &CreateStashResponse{
		ID: id,
	}

	return res, nil
}

func (h *Handler) GetAllProducts() ([]*Product, error) {
	rows, err := h.db.Query("select * from product")
	if err != nil {
		return nil, err
	}

	products := []*Product{}
	for rows.Next() {
		product, err := scanIntoProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (h *Handler) GetProductById(id int) (*Product, error) {
	rows, err := h.db.Query("select * from product where id = ?", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoProduct(rows)
	}

	return nil, fmt.Errorf("account %d not found", id)
}

func scanIntoProduct(rows *sql.Rows) (*Product, error) {
	product := new(Product)
	err := rows.Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.StashId,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	return product, err
}

func (h *Handler) Init() error {
	fmt.Println("os: ", os.Getenv("environment"))
	if os.Getenv("environment") == "production" {
		return nil
	}

	fmt.Println("dropping all tables")
	h.db.Exec("DROP TABLE user")
	h.db.Exec("DROP TABLE stash")
	h.db.Exec("DROP TABLE product")
	h.db.Exec("DROP TABLE session")
	fmt.Println("dropped all tables")

	fmt.Println("creating tables")
	h.createUsersTable()
	h.createStashesTable()
	h.createProductsTable()
	h.createSessionsTable()
	fmt.Println("just created tables")

	_, err := h.db.Exec("INSERT INTO `stash` (name, location) VALUES ('Mendes store', 'Sweden');")
	h.db.Exec("INSERT INTO `product` (name, price, stashId) VALUES ('Red hoodie', 49.99, 1);")
	h.db.Exec("INSERT INTO `product` (name, price, stashId) VALUES ('Blue hoodie', 19.99, 1);")
	h.db.Exec("INSERT INTO `product` (name, price, stashId) VALUES ('Pink shoe', 49.99, 2);")
	h.db.Exec("INSERT INTO `product` (name, price, stashId) VALUES ('Orange hat', 49.99, 2);")

	if err != nil {
		return err
	}

	return nil
}
