package entities

import "gorm.io/gorm"

// import "gorm.io/gorm"

type OrderUser struct {
	gorm.Model //creates ID of uint represents orderid
	Userid     uint64
	Status     string
}

type OrderListItems struct {
	gorm.Model //creates ID of uint represents orderlistitemsid
	Orderid    uint64
	Productid  uint64
	Quantity uint64
}