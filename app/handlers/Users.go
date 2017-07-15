package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"test/app/dbc"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int    `json:"Id"`
	Username string `json:"Username"`
	Password string `json:"Password"`
	Email    string `json:"Email"`
	Date     int    `json:"Date"`
}
type HttpResp struct {
	Status      int    `json:"status"`
	Description string `json:"description"`
	Body        string `json:"body"`
	Error       string `json:"error"`
	UUID        string `json:"uuid"`
}

type Users []User

func hashPassword(pass string) string {
	password := []byte(pass)
	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	}
	return (string(hashedPassword))
}

func compearPassword(hP string, p string) bool {
	hashedPassword := []byte(hP)
	password := []byte(p)

	err := bcrypt.CompareHashAndPassword(hashedPassword, password)
	if err != nil {
		fmt.Println(err) // nil means it is a match
		return false
	} else {
		return true
	}
}

func getUser(u User) User {
	db := dbc.Connect()
	fmt.Print(db)
	var user User
	err := db.QueryRow("SELECT idUsers, Username, Password, Email, Date FROM Users where Username = ?", u.Username).Scan(&user.Id, &user.Username, &user.Password, &user.Email, &user.Date)
	if err != nil {
		log.Print(err.Error()) // proper error handling instead of panic in your app
		log.Print("Failed to select tag from database")
		db.Close()
		return user
	} else {
		log.Print(user)
		db.Close()
		return user
	}
}

func setUser(u User) {
	db := dbc.Connect()
	defer db.Close()
	stmt, _ := db.Prepare("INSERT INTO Users (Username, Password, Email, Date) values (?,?,?,?) ")
	res, err := stmt.Exec(u.Username, hashPassword(u.Password), u.Email, u.Date)
	if err != nil {
		log.Print(err.Error()) // proper error handling instead of panic in your app
		log.Print("Failed to insert post into database")
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Print("Failed to get last insert id")
	}
	log.Print("Successfully Inserted Post %i Into the Database", id)

	//SetAudit(&Audit{Type: "stuff", Desc: "test", Date: int(time.Now().Unix())})
}

func Register(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var user User
	err := decoder.Decode(&user)
	if err != nil {
		json.NewEncoder(w).Encode(HttpResp{Status: 500, Description: "Failed to registered", Error: err.Error()})
	} else {
		setUser(user)
		setAudit(Audit{TITLE: "AUTH", TEXT: "A user registered", DATE: time.Now().Unix(), USERNAME: user.Username})
		json.NewEncoder(w).Encode(HttpResp{Status: 200, Description: "Sucssfully Registered"})
	}

}

func Login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var user User
	err := decoder.Decode(&user)
	if err != nil {
		json.NewEncoder(w).Encode(HttpResp{Status: 500, Description: "Failed to registered", Error: err.Error()})
	} else {
		u := getUser(user)
		if compearPassword(u.Password, user.Password) {
			a := Audit{TITLE: "AUTH", TEXT: "A user logged in", DATE: time.Now().Unix(), USERNAME: u.Username}
			setAudit(a)
			json.NewEncoder(w).Encode(HttpResp{Status: 200, Description: "Login Sucesss "})
		} else {
			setAudit(Audit{TITLE: "AUTH", TEXT: "A user logged in failed", DATE: time.Now().Unix(), USERNAME: u.Username})
			json.NewEncoder(w).Encode(HttpResp{Status: 500, Description: "Login failed"})
		}
	}
}
