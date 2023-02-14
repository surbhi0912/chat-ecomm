package main

//Requirements:
//Product should have following attributes : name, product id, description and price
//User should have following attributes: name, user id, address, date of birth
//Admin should be able to add, update and delete products
//Logged in user should be able to browse products
//Logged in user should have a shopping cart where user should be able to add multiple products
//User should have ability to checkout and total payable should be displayed while checkout
//User and Product information should be persisted in database

//Learning:
//Designing classes to store the user and Product information
//Taking input from console and storing into models
//Persisting data in database
//Database table design
//Coding best practices like naming of variables, class names, designing helping and service classes

import (
	// "encoding/json"

	"chat-ecomm/controllers"
	"chat-ecomm/database"
	"chat-ecomm/entities"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"flag"
	"log"

	// "github.com/gorilla/sessions"

	"github.com/gorilla/mux"
)

var addr = flag.String("addr", ":8080", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/chat" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

// var (
//     // key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
//     key = []byte("super-secret-key")
//     store = sessions.NewCookieStore(key)
// )

func main() {
	LoadAppConfig() //loads configurations from config.json using viper

	//initialize database
	database.Connect(AppConfig.ConnectionString)
	database.Migrate()

	router := mux.NewRouter().StrictSlash(true) //initialise the router
	//strictslash when false, if the route path is "/path", accessing "/path/" will not match this route and vice versa

	flag.Parse()
	hub := newHub()
	go hub.run()
	// http.HandleFunc("/chat", serveHome)
	// http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
	// 	serveWs(hub, w, r)
	// })
	// err := http.ListenAndServe(*addr, nil)
	// if err != nil {
	// 	log.Fatal("ListenAndServe: ", err)
	// }
	router.HandleFunc("/chat", serveHome)
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	router.HandleFunc("/welcome", welcome)

	router.HandleFunc("/signup", signup)
	router.HandleFunc("/", login)
	router.HandleFunc("/signin", signin)
	router.HandleFunc("/createUser", controllers.CreateUser)

	// router.HandleFunc("/api/users", controllers.GetUsers).Methods("GET") //read
	// router.HandleFunc("/api/users", controllers.CreateUser).Methods("POST") //create

	router.HandleFunc("/admincontrol", admincontrol)
	router.HandleFunc("/browse", browse)

	router.HandleFunc("/addProduct", controllers.CreateProduct)

	router.HandleFunc("/api/products", controllers.GetProducts).Methods("GET")           //read
	router.HandleFunc("/api/products", controllers.CreateProduct).Methods("POST")        //create
	router.HandleFunc("/api/products/{id}", controllers.GetProductById).Methods("GET")   //read
	router.HandleFunc("/api/products/{id}", controllers.UpdateProduct).Methods("PUT")    //update
	router.HandleFunc("/api/products/{id}", controllers.DeleteProduct).Methods("DELETE") //delete

	router.HandleFunc("/addCart", controllers.AddCart)
	router.HandleFunc("/shoppingcart", controllers.ShoppingCart)
	router.HandleFunc("/checkout", controllers.Checkout)

	// http.ListenAndServe(fmt.Sprintf(":%v", AppConfig.Port), router)
	err := http.ListenAndServe(fmt.Sprintf(":%v", AppConfig.Port), router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome! You were expected!")
}

func login(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("loginPage.html"))
	tmpl.Execute(w, nil)
}

func signup(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("register.html"))
	tmpl.Execute(w, nil)
}

func signin(w http.ResponseWriter, r *http.Request) {
	// session, _ := store.Get(r, "cookieforUserID")

	usernameValue := r.FormValue("inputEmailValue1")
	passwordValue := r.FormValue("inputPasswordValue1")
	var user entities.User
	database.Instance.Where("Email = ?", usernameValue).First(&user)
	if user.ID == 0 { //Email doesn't exist in database
		http.Redirect(w, r, "http://localhost:8080/signup", http.StatusSeeOther)
	} else { //Email exists in database
		if user.Password == passwordValue {
			// w.Header().Add("x-user-id", "22")
			// // w.Write()
			if user.Role == "admin" {
				http.Redirect(w, r, "http://localhost:8080/admincontrol", http.StatusSeeOther)
			} else {
				uid := strconv.FormatUint(uint64(user.ID), 10)
				// println("====")
				// println(w.Header().Get("x-user-id"))
				// println("===")
				usercookie := http.Cookie{
					Name:  "cookieforUserID",
					Value: uid,
				}
				http.SetCookie(w, &usercookie)
				// session.Save(r, w)
				http.Redirect(w, r, "http://localhost:8080/browse", http.StatusSeeOther)
			}
		} else {
			fmt.Fprintf(w, "Oops! Username and password did not match.")
		}
	}
}

func admincontrol(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("adminpanel.html"))
	tmpl.Execute(w, nil)
}

func browse(w http.ResponseWriter, r *http.Request) {
	// println(w.Header().Get("x-user-id"))
	// println(w)
	// println("----")
	// println(r.Header.Get("x-user-id"))
	// println(r)
	var p []entities.Product
	database.Instance.Find(&p) //maps all available products from database to the products list variable

	tmpl := template.Must(template.ParseFiles("browse.html"))

	var prdList entities.Productlist

	var prds []entities.ProductVO
	for _, item := range p {
		if item.Quantity != 0 { //only products whose atleast 1 quantity is available
			var prd entities.ProductVO
			prd.ID = item.ID
			prd.Name = item.Name
			prd.Description = item.Description
			prd.Price = item.Price
			prd.Quantity = item.Quantity
			prds = append(prds, prd)
		}
	}

	prdList.Productdetails = prds
	// fmt.Println(prdList)
	tmpl.Execute(w, prdList)

	// data := entities.Productlist{
	// 	Productdetails: []entities.Product{
	// 		for i, item := range p {
	// 			{Name: item.Name, Price: item.Price},
	// 			// {Name: "Man 2", Price: "22"},
	// 		}
	// 	},
	// }
}