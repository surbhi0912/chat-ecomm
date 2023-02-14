package controllers

import (
	"fmt"
	"log"
	// "encoding/json"
	"chat-ecomm/database"
	"chat-ecomm/entities"
	"net/http"
)

// func CreateUser(w http.ResponseWriter, r *http.Request) {//for GET and POST via postman
// 	w.Header().Set("Content-Type", "application/json")
// 	var user entities.User                //new user variable
// 	json.NewDecoder(r.Body).Decode(&user) //decodes body of incoming json request and maps it to newly created variable
// 	database.Instance.Create(&user)       //Using GORM we create a new user by passing in the parsed user, which creates a new record in the user table
// 	json.NewEncoder(w).Encode(user) //Returns the newly created user data back to the client
// }

func CreateUser(w http.ResponseWriter, r *http.Request) {
	usernameValue := r.FormValue("exampleInputEmail1")
	passwordValue := r.FormValue("exampleInputPassword1")
	nameValue := r.FormValue("exampleName1")
	addressValue := r.FormValue("exampleAddress1")
	dobValue := r.FormValue("exampleDOB1")
	var user1 entities.User
	database.Instance.Where("Email = ?", usernameValue).First(&user1) //finds records in db with matching usernameFromForm
	if user1.ID == 0 { //Email doesn't exist in database
		user := entities.User{Name: nameValue, Email: usernameValue, Password: passwordValue, Address: addressValue, DOB: dobValue}
		result := database.Instance.Create(&user)
		if result == nil {
			log.Fatal("Error inserting " + usernameValue + " into database.")
			fmt.Fprintf(w, "Oops! try again later.")
		} else {
			fmt.Fprintf(w, "Congratulations " + usernameValue + " You are successfully registered.")
		}
	} else { //Email exists in database
		fmt.Fprintf(w, "An account with " + usernameValue + " is already registered, please signup with a new username.")
	}
}

// func GetUsers(w http.ResponseWriter, r *http.Request) {
// 	var users []entities.User
// 	database.Instance.Find(&users) //maps all available products from database to the products list variable
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(users)
// }