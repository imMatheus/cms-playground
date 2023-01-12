package main

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Name string
}

type Stash struct {
	gorm.Model

	Name     string `json:"name"`
	Location string `json:"location"`
}

type Product struct {
	gorm.Model

	Name    string  `json:"name"`
	Price   float64 `json:"price"`
	StashId int     `json:"stash_id"`
	Stash   Stash   `json:"stash"`
}
