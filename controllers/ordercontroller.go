package controllers

import (
	"fmt"
	"chat-ecomm/database"
	"chat-ecomm/entities"
	"html/template"
	"net/http"
	"strconv"
)

func AddCart(w http.ResponseWriter, r *http.Request) {
	usercookie, _ := r.Cookie("cookieforUserID")
	P := r.FormValue("egProd1")
	UId, _ := strconv.ParseUint(usercookie.Value, 10, 64)
	PId, _ := strconv.ParseUint(P, 10, 64)
	
	var orderId uint64
	
	//orderuser table
	var findUserOrder entities.OrderUser
	database.Instance.Where("Userid = ? AND Status = ?", UId, "Cart").First(&findUserOrder) //have to use variable name in table and not in struct
	//if declared in struct as UserID then user_id created in table, if declared as Userid then userid created in table
	// fmt.Println(findUserOrder)

	if findUserOrder.ID == 0 { //userid does not exist
		userorder := entities.OrderUser{Userid: UId, Status: "Cart"} 
		database.Instance.Create(&userorder)
		database.Instance.Where("Userid = ? AND Status = ?", UId, "Cart").First(&findUserOrder)
		orderId = uint64(findUserOrder.ID) //retrives orderid
	} else {
		orderId = uint64(findUserOrder.ID) //retrieves orderid
	}
	
	//orderlistitems table
	var findProductOrder entities.OrderListItems
	database.Instance.Where("Orderid = ? AND Productid = ?", orderId, PId).First(&findProductOrder) //if product exists in that order id
	if findProductOrder.ID == 0 { //product does not exist in order
		orderlist := entities.OrderListItems{Orderid: orderId, Productid: PId, Quantity: 1}
		database.Instance.Create(&orderlist)
	} else{
		//updating orderlistitems table product quantity
		findProductOrder.Quantity = findProductOrder.Quantity + 1
		database.Instance.Save(&findProductOrder)
	}

	//updating inventory
	var product entities.Product
	database.Instance.First(&product, PId)
	product.Quantity = product.Quantity - 1
	database.Instance.Save(&product)

	oid := strconv.FormatUint(orderId, 10)
	ordercookie := http.Cookie{
		Name:  "cookieforOrderID",
		Value: oid,
	}
	
	http.SetCookie(w, usercookie)
	http.SetCookie(w, &ordercookie)
	http.Redirect(w, r, "http://localhost:8080/shoppingcart", http.StatusSeeOther)
}

func ShoppingCart(w http.ResponseWriter, r *http.Request) {
	usercookie, _ := r.Cookie("cookieforUserID")
	ordercookie, _ := r.Cookie("cookieforOrderID")

	tmpl := template.Must(template.ParseFiles("shoppingcart.html"))
	oid, _ := strconv.ParseUint(ordercookie.Value, 10, 64)
	//add html code to Display the added products and their sum
	
	var orderitems []entities.OrderListItems
	database.Instance.Where("Orderid = ?", oid).Find(&orderitems) //maps all products from the orderid orderitems

	var prdList entities.Productlist

	var prds []entities.ProductVO

	prdList.Total = 0

	for _, item := range orderitems {
		var prd entities.ProductVO
		prd.ID = uint(item.Productid)
		
		var product entities.Product
		database.Instance.First(&product, item.Productid) //querying Product table with productid and saving in product

		prd.Name = product.Name
		prd.Description = product.Description
		prd.Price = product.Price*float64(item.Quantity)
		prd.Quantity = item.Quantity
		prds = append(prds, prd)

		prdList.Total += prd.Price
	}

	prdList.Productdetails = prds
	// fmt.Println(prdList)
	http.SetCookie(w, usercookie)
	http.SetCookie(w, ordercookie)

	tmpl.Execute(w, prdList)
}

func Checkout(w http.ResponseWriter, r *http.Request) {
	ordercookie, _ := r.Cookie("cookieforOrderID")
	oid, _ := strconv.ParseUint(ordercookie.Value, 10, 64)
	var order entities.OrderUser
	database.Instance.First(&order, oid)
	order.Status = "Success" //make status in orderuser for that orderid to Success
	database.Instance.Save(&order)
	fmt.Fprintf(w, "Congratulations! Order #"+ordercookie.Value+" placed sucessfully.")
	// http.SetCookie(w, usercookie)
	// http.SetCookie(w, ordercookie)
}
