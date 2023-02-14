package entities

import "gorm.io/gorm"

type Product struct {
	//ID          uint    `json:"id"`
	gorm.Model //creates ID of uint
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	Quantity    uint64  `json:"qty"`
}

type ProductVO struct { //making a new class that will not make changes to DB
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	Quantity    uint64  `json:"qty"`
}

type Productlist struct {
	Productdetails []ProductVO
	Total float64
}