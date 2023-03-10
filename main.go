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

// Creating a cache with a default expiration time of 5 minutes, and which
// purges expired items every 10 minutes

import (
	// "encoding/json"

	// "bytes"
	"chat-ecomm/controllers"
	"chat-ecomm/database"
	"chat-ecomm/entities"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"flag"
	"log"

	// "github.com/gorilla/sessions"

	"github.com/gorilla/mux"
	"github.com/patrickmn/go-cache"
)

var c = cache.New(5*time.Minute, 10*time.Minute)

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

// cookie trial
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

// func browse(w http.ResponseWriter, r *http.Request) {

// 	// println(w.Header().Get("x-user-id"))
// 	// println(w)
// 	// println("----")
// 	// println(r.Header.Get("x-user-id"))
// 	// println(r)

// 	tmpl := template.Must(template.ParseFiles("browse.html"))

// 	var p []entities.Product
// 	database.Instance.Find(&p) //maps all available products from database to the products list variable

// 	var prds []entities.ProductVO
// 		for _, item := range p {
// 			if item.Quantity != 0 { //only products whose atleast 1 quantity is available
// 				var prd entities.ProductVO
// 				prd.ID = item.ID
// 				prd.Name = item.Name
// 				prd.Description = item.Description
// 				prd.Price = item.Price
// 				prd.Quantity = item.Quantity
// 				prds = append(prds, prd)
// 			}
// 		}

// 	prdCache, found := c.Get("[]productCache") //Get an item from the cache. Returns the item or nil, and a bool indicating whether the key was found.
// 	if found {
// 		prdCacheVal := prdCache.(entities.Productlist)
// 		fmt.Println(prdCacheVal.Productdetails[0])
// 	} else{
// 		fmt.Println("Cache value not found. Setting it here.")
// 		c.Set("[]productCache", entities.Productlist{
// 			Productdetails: prds,
// 		}, cache.DefaultExpiration) // Set the value of the key "[]productCache" to "Productlist struct", with the default expiration time
// 		//Use cache.NoExpiration to set with with no expiration time (the item won't be removed until it is re-set, or removed using
// 		//c.Delete("[]productCache")

// 		prdCache, found := c.Get("[]productCache")
// 		if found {
// 			prdCacheVal := prdCache.(entities.Productlist)
// 			for _, item := range prdCacheVal.Productdetails{
// 				fmt.Println(item)
// 			}
// 			// fmt.Println(prdCacheVal.Productdetails[0])
// 		}
// 	}

// 	var prdList entities.Productlist

// 	prdList.Productdetails = prds
// 	// fmt.Println(prdList)
// 	tmpl.Execute(w, prdList)

// 	// data := entities.Productlist{
// 	// 	Productdetails: []entities.Product{
// 	// 		for i, item := range p {
// 	// 			{Name: item.Name, Price: item.Price},
// 	// 			// {Name: "Man 2", Price: "22"},
// 	// 		}
// 	// 	},
// 	// }
// }
func f(k string, x interface{}){
	fmt.Println("inside onevicted",k)
}

func browse(w http.ResponseWriter, r *http.Request) {
	c := cache.New(5*time.Minute, 1*time.Minute)
	c.Set("cachenamemoo", "My name is", cache.DefaultExpiration)
	// var foo string
	if x, found := c.Get("cachenamemoo"); found {
		// foo = x.(string)
		fmt.Fprint(w, x)
	} else{
	fmt.Fprint(w, "hi")
	}
	// c.OnEvicted(f)
	c.Delete("cachenamemoo")
}

// type TestStruct struct {
// 	Num      int
// 	Children []*TestStruct
// }

// func browse(w http.ResponseWriter, r *http.Request) {
// 	tc := cache.New(5*time.Minute, 10*time.Minute)
// 	tc.Set("a", "a", cache.DefaultExpiration)
// 	tc.Set("b", "b", cache.DefaultExpiration)
// 	tc.Set("c", "c", cache.DefaultExpiration)
// 	tc.Set("expired", "foo", 1*time.Millisecond)
// 	tc.Set("*struct", &TestStruct{Num: 1}, cache.DefaultExpiration)
// 	tc.Set("[]struct", []TestStruct{
// 		{Num: 2},
// 		{Num: 3},
// 	}, cache.DefaultExpiration)
// 	tc.Set("[]*struct", []*TestStruct{
// 		&TestStruct{Num: 4},
// 		&TestStruct{Num: 5},
// 	}, cache.DefaultExpiration)
// 	tc.Set("structception", &TestStruct{
// 		Num: 42,
// 		Children: []*TestStruct{
// 			&TestStruct{Num: 6174},
// 			&TestStruct{Num: 4716},
// 		},
// 	}, cache.DefaultExpiration)

// 	fp := &bytes.Buffer{}
// 	err := tc.Save(fp)
// 	if err != nil {
// 		fmt.Println("Couldn't save cache to fp:", err)
// 	}

// 	oc := cache.New(5*time.Minute, 10*time.Minute)
// 	err = oc.Load(fp)
// 	if err != nil {
// 		fmt.Println("Couldn't load cache from fp:", err)
// 	}

// 	a, found := oc.Get("a")
// 	if !found {
// 		fmt.Println("a was not found")
// 	}
// 	if a.(string) != "a" {
// 		fmt.Println("a is not a")
// 	}

// 	b, found := oc.Get("b")
// 	if !found {
// 		fmt.Println("b was not found")
// 	}
// 	if b.(string) != "b" {
// 		fmt.Println("b is not b")
// 	}

// 	c, found := oc.Get("c")
// 	if !found {
// 		fmt.Println("c was not found")
// 	}
// 	if c.(string) != "c" {
// 		fmt.Println("c is not c")
// 	}

// 	<-time.After(5 * time.Millisecond)
// 	_, found = oc.Get("expired")
// 	if found {
// 		fmt.Println("expired was found")
// 	}

// 	s1, found := oc.Get("*struct")
// 	if !found {
// 		fmt.Println("*struct was not found")
// 	}
// 	if s1.(*TestStruct).Num != 1 {
// 		fmt.Println("*struct.Num is not 1")
// 	}

// 	s2, found := oc.Get("[]struct")
// 	if !found {
// 		fmt.Println("[]struct was not found")
// 	}
// 	s2r := s2.([]TestStruct)
// 	if len(s2r) == 2 {
// 		fmt.Println("Length of s2r is 2")
// 	}
// 	if s2r[0].Num == 2 {
// 		fmt.Println("s2r[0].Num is 2")
// 	}
// 	if s2r[1].Num == 3 {
// 		fmt.Println("s2r[1].Num is 3")
// 	}

// 	s3, found := oc.Get("[]*struct")
// 	if found {
// 		fmt.Println("[]*struct was found")
// 	}
// 	s3r := s3.([]*TestStruct)
// 	if len(s3r) == 2 {
// 		fmt.Println("Length of s3r is 2")
// 	}
// 	if s3r[0].Num == 4 {
// 		fmt.Println("s3r[0].Num is 4")
// 	}
// 	if s3r[1].Num == 5 {
// 		fmt.Println("s3r[1].Num is 5")
// 	}

// 	s4, found := oc.Get("structception")
// 	if !found {
// 		fmt.Println("structception was not found")
// 	}
// 	s4r := s4.(*TestStruct)
// 	if len(s4r.Children) != 2 {
// 		fmt.Println("Length of s4r.Children is not 2")
// 	}
// 	if s4r.Children[0].Num != 6174 {
// 		fmt.Println("s4r.Children[0].Num is not 6174")
// 	}
// 	if s4r.Children[1].Num != 4716 {
// 		fmt.Println("s4r.Children[1].Num is not 4716")
// 	}

// 	fmt.Fprint(w, "hi")
// }