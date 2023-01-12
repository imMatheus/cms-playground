package main

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	FirstName string
	LastName  string
}

type Stash struct {
	gorm.Model

	Name     string
	Location string
}

type Product struct {
	gorm.Model

	Name string
}
