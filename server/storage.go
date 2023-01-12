package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Storage interface {
	GetAllUsers() ([]*User, error)
	GetAllStashes() ([]*Stash, error)
	GetAllProducts() ([]*Product, error)
}

// A Handler is an HTTP API server handler.
type Handler struct {
	db *gorm.DB
}

func NewStore() (*Handler, error) {
	// Load environment variables from file.
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load environment variables: %v", err)
	}

	// Connect to PlanetScale database using DSN environment variable.
	db, err := gorm.Open(mysql.Open(os.Getenv("DSN")), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		log.Fatalf("failed to connect to PlanetScale: %v", err)
	}

	log.Println("Successfully connected to PlanetScale!")

	return &Handler{
		db: db,
	}, nil
}

func (h *Handler) GetAllUsers() ([]*User, error) {
	var users []*User
	result := h.db.Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

func (h *Handler) GetAllStashes() ([]*Stash, error) {
	var stashes []*Stash
	result := h.db.Find(&stashes)

	if result.Error != nil {
		return nil, result.Error
	}

	return stashes, nil
}

func (h *Handler) GetAllProducts() ([]*Product, error) {
	var products []*Product
	result := h.db.Find(&products)

	if result.Error != nil {
		return nil, result.Error
	}

	return products, nil
}

func (h *Handler) Init() error {
	fmt.Println("os: ", os.Getenv("environment"))
	if os.Getenv("environment") == "production" {
		return nil
	}

	fmt.Println("Drop all dbs")
	h.db.Migrator().DropTable(&User{})
	h.db.Migrator().DropTable(&Stash{})
	h.db.Migrator().DropTable(&Product{})
	h.db.Migrator().DropTable("projects")

	fmt.Println("Auto migrate user")
	if err := h.db.AutoMigrate(&User{}); err != nil {
		return err
	}

	fmt.Println("Auto migrate stashes")
	if err := h.db.AutoMigrate(&Stash{}); err != nil {
		return err
	}

	fmt.Println("Auto migrate products")
	if err := h.db.AutoMigrate(&Product{}); err != nil {
		return err
	}

	h.db.Create(&Stash{
		Name:     "Cool things 2",
		Location: "yo mama house",
	})

	h.db.Create(&Stash{
		Name:     "Cool things 2",
		Location: "yo mama house",
	})
	h.db.Create(&Stash{
		Name:     "Cool things 2",
		Location: "yo mama house",
	})
	h.db.Create(&Product{
		Name:  "Cool things 2",
		Price: 45.99,
	})
	h.db.Create(&Product{
		Name:  "Cool things 4",
		Price: 12.99,
	})
	h.db.Create(&Product{
		Name:  "Cool things 2",
		Price: 75.00,
	})
	h.db.Create(&Product{
		Name:  "Hat attack",
		Price: 99,
	})
	h.db.Create(&Stash{
		Name:     "Cool things 2",
		Location: "yo mama house",
	})

	h.db.Create(&User{
		Name: "MAth",
	})

	return nil
}
