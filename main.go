package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
"github.com/julienschmidt/httprouter"
	_ "github.com/go-sql-driver/mysql"
)

// User is representation of a user
type User struct {
	ID       string `form:"id" json:"id" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Name     string `form:"name" json:"name" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required"`
	Hakbun   string `form:"hakbun" json:"hakbun" binding:"required"`
	Birth    string `form:"birth" json:"birth"`
}

//DuplicateMsg is
type DuplicateMsg struct {
	Status       int    `json:"status"`
	ID           string `json:"id"`
	IsDuplicated bool   `json:"is_duplicated"`
}

type SigninMsg struct {
	Status       int    `json:"status"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

var signinTemplate, _ = template.ParseFiles("html/signin.html")
var signupTemplate, _ = template.ParseFiles("html/signup.html")
var indexTemplate, _ = template.ParseFiles("html/index.html")

func signinPage(w http.ResponseWriter, r *http.Request) {
	signinTemplate.Execute(w, nil)
}
func signupPage(w http.ResponseWriter, r *http.Request) {
	signupTemplate.Execute(w, nil)
}
func indexPage(w http.ResponseWriter, r *http.Request) {
	indexTemplate.Execute(w, nil)
}

// SigninHandler is signin handler
func SigninHandler(w http.ResponseWriter, r *http.Request) {
	var person User
	if r.Header.Get("Content-Type") == "application/json" {
		len := r.ContentLength
		body := make([]byte, len)
		r.Body.Read(body)
		json.Unmarshal([]byte(body), &person)
	} else if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
		r.ParseForm()
		person.ID = r.FormValue("id")
		person.Password = r.FormValue("password")
	}
	fmt.Println(person)
	db, err := sql.Open("mysql", "root:111111@tcp(127.0.0.1:3306)/user")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	row, err := db.Query("SELECT id, password from user where id >= ?", 1)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var found bool = false
	var tempid string
	var temppassword string
	for row.Next() {
		err := row.Scan(&tempid, &temppassword)
		if err != nil {
			log.Fatal(err)
		}
		// duplicate
		if person.ID == tempid && person.Password == temppassword {
			found = true
			/*
				msg := SigninMsg{200, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NTg1MzcxNzcsImlzcyI6ImRhc29tLmlvIiwic3ViIjoiYWNjZXNzX3Rva2VuIiwidXNlcl9iaXJ0aCI6IjE5OTgtMDEtMDFUMDA6MDA6MDBaIiwidXNlcl9oYWtidW4iOjAsInVzZXJfaWQiOiIyMDE3MTAwMDAwIiwidXNlcl9qb2luX2RhdGUiOiIyMDE5LTA1LTIyVDA0OjQ3OjQ4WiIsInVzZXJfbmFtZSI6Iu2FjOyKpO2KuCJ9.h55HKpBmxMRHVT-wxsDTIqglV-GLKjWCwEvIuF3yY-s", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NTg1NzIyNzcsImlzcyI6ImRhc29tLmlvIiwic3ViIjoicmVmcmVzaF90b2tlbiIsInVzZXJfaWQiOiIyMDE3MTAwMDAwIn0.9vpOkQYFbb_MeocXhQdnTC3wTdTbeX-MaQnmptfdODs"}
				enc := json.NewEncoder(w)
				w.Header().Set("Content-Type", "application/json")
				enc.Encode(msg)
				//fmt.Fprintf(w, "%s", msg)
				fmt.Println("good")
			*/
			break
		}
		if !found {
			defer http.Redirect(w, r, "/signinpage", http.StatusFound)
		}
		if found {
			jwtToken, err := person.GetJwtToken()
		}
	}
}

// Signin is signin
func Signin(w http.ResponseWriter, r *http.Request) {

}

// SignupHandler is signup handler
func SignupHandler(w http.ResponseWriter, r *http.Request) {

	var person User

	if r.Header.Get("Content-Type") == "application/json" {
		len := r.ContentLength
		body := make([]byte, len)
		r.Body.Read(body)
		json.Unmarshal([]byte(body), &person)
	} else if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
		r.ParseForm()
		person.ID = r.FormValue("id")
		person.Password = r.FormValue("password")
		person.Name = r.FormValue("name")
		person.Email = r.FormValue("email")
		person.Hakbun = r.FormValue("hakbun")
		person.Birth = r.FormValue("birth")
	}
	fmt.Println(person)
	var id = person.ID
	var password = person.Password
	var name = person.Name
	var email = person.Email
	var hakbun = person.Hakbun
	var birth = person.Birth
	fmt.Println(id)
	url := "http://127.0.0.1:3000/duplicate?id=" + id
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err == nil {
		var msg DuplicateMsg
		json.Unmarshal([]byte(responseData), &msg)
		fmt.Println(msg.IsDuplicated)
		if msg.IsDuplicated == false {
			db, err := sql.Open("mysql", "root:111111@tcp(127.0.0.1:3306)/user")
			if err != nil {
				log.Fatal(err)
			}
			defer db.Close()
			_, err = db.Exec("INSERT user (id, password, name, email, hakbun, birth) values (?, ?, ?, ?, ?, ?)", id, password, name, email, hakbun, birth)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			defer http.Redirect(w, r, "/", http.StatusFound)
		}
	}

}

// DuplicateHandler is check if ID is duplicate
func DuplicateHandler(w http.ResponseWriter, r *http.Request) {
	key, _ := r.URL.Query()["id"]

	db, err := sql.Open("mysql", "root:111111@tcp(127.0.0.1:3306)/user")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	row, err := db.Query("SELECT id from user where UID >= ?", 1)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var found bool = false
	var tempid string
	for row.Next() {
		err := row.Scan(&tempid)
		if err != nil {
			log.Fatal(err)
		}
		// duplicate
		if key[0] == tempid {
			found = true
			msg := DuplicateMsg{200, key[0], true}
			enc := json.NewEncoder(w)
			w.Header().Set("Content-Type", "application/json")
			enc.Encode(msg)
			break
		}
	}
	// not duplicate
	if !found {
		msg := DuplicateMsg{200, key[0], false}
		enc := json.NewEncoder(w)
		w.Header().Set("Content-Type", "application/json")
		enc.Encode(msg)
	}
}

func main() {

	http.HandleFunc("/", signinPage)
	http.HandleFunc("/signuppage", signupPage)
	http.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			SigninHandler(w, r)
		}
	})
	http.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			SignupHandler(w, r)
		}
	})
	http.HandleFunc("/duplicate", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			DuplicateHandler(w, r)
		}
	})
	http.HandleFunc("/index", indexPage)
	http.ListenAndServe(":3000", nil)

}
