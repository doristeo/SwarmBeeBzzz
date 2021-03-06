package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

//---------------------------------------------------------------------------------------------------------------------------------
func CreateDatabase() (*sql.DB, error) {
	serverName := "10.1.1.109:3306"
	user := "home"
	password := "111111"
	dbName := "swarms"

	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true&multiStatements=true", user, password, serverName, dbName)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	return db, nil
}

//---------------------------------------------------------------------------------------------------------------------------------
type Logdata struct {
	Name      string `json:"name"`
	Peers     int    `json:"peers"`
	Diskavail int    `json:"diskavail"`
	Diskfree  int    `json:"diskfree"`
}

//---------------------------------------------------------------------------------------------------------------------------------
func postFunction(w http.ResponseWriter, r *http.Request) {

	database, err := CreateDatabase()
	if err != nil {
		log.Fatal("Database connection failed")
	}

	var logdata Logdata
	json.NewDecoder(r.Body).Decode(&logdata)

	s := "INSERT nodes(name, peers, diskavail, diskfree) VALUES ('" + logdata.Name + "', " + strconv.Itoa(logdata.Peers) + ", " +
		strconv.Itoa(logdata.Diskavail) + ", " + strconv.Itoa(logdata.Diskfree) + ")"
	fmt.Println(s)
	_, err = database.Exec(s)
	if err != nil {
		log.Println("Database INSERT failed")
	}
}

//---------------------------------------------------------------------------------------------------------------------------------
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", postFunction).Methods("POST")

	srv := &http.Server{
		Handler: router,
		Addr:    ":8080",
	}
	log.Fatal(srv.ListenAndServe())

	//setupRouter(router)

	//log.Fatal(http.ListenAndServe(":8080", router))
}
