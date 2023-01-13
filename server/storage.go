package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Storage interface {
	GetAllUsers() ([]*User, error)
	GetAllStashes() ([]*Stash, error)
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

func (h *Handler) createUsersTable() error {
	query := `CREATE TABLE if not exists user (
		id int PRIMARY KEY AUTO_INCREMENT,
		name varchar(255) NOT NULL,
		image varchar(255) NOT NULL,
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

func scanIntoUser(rows *sql.Rows) (*User, error) {
	user := new(User)
	err := rows.Scan(
		&user.ID,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	return user, err
}

func (h *Handler) GetAllStashes() ([]*Stash, error) {
	rows, err := h.db.Query("select * from stash")
	if err != nil {
		return nil, err
	}

	stashes := []*Stash{}
	for rows.Next() {
		var stash *Stash
		err := rows.Scan(
			&stash.ID,
			&stash.Location,
			&stash.Name,
			&stash.CreatedAt,
			&stash.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		stashes = append(stashes, stash)
	}

	return stashes, nil
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

	// fmt.Println("dropping all tables")
	// h.db.Exec("DROP TABLE user")
	// h.db.Exec("DROP TABLE stash")
	// h.db.Exec("DROP TABLE product")
	// fmt.Println("dropped all tables")

	// fmt.Println("creating tables")
	// h.createUsersTable()
	// h.createStashesTable()
	// h.createProductsTable()
	// fmt.Println("just created tables")

	// _, err := h.db.Exec("INSERT INTO `stash` (name, location) VALUES ('Mendes store', 'Sweden');")
	// h.db.Exec("INSERT INTO `product` (name, price, stashId) VALUES ('Red hoodie', 49.99, 1);")
	// h.db.Exec("INSERT INTO `product` (name, price, stashId) VALUES ('Blue hoodie', 19.99, 1);")
	// h.db.Exec("INSERT INTO `product` (name, price, stashId) VALUES ('Pink shoe', 49.99, 2);")
	// h.db.Exec("INSERT INTO `product` (name, price, stashId) VALUES ('Orange hat', 49.99, 2);")

	// if err != nil {
	// 	return err
	// }

	return nil
}
