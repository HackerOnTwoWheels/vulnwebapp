package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var GlobalDB sql.DB

func dbSampleUsers(db *sql.DB, id int, email string) {
	fmt.Printf("## Creating user :%v\n", email)
	insUserstat, err := db.Prepare("INSERT into users values (?,?)")
	if err != nil {
		fmt.Println(err)
		return
	}
	insUserstat.Exec(id, email)
	insUserstat.Close()
}
func dbCreate(db *sql.DB) {
	fmt.Println("## Creating table ##")
	crtTbStat, err := db.Prepare("CREATE TABLE if not exists users (id INT PRIMARY KEY, email varchar(250))")
	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := crtTbStat.Exec()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
	crtTbStat.Close()
	dbSampleUsers(db, 1, "bob@local.com")
	dbSampleUsers(db, 2, "smith@local.com")
	dbSampleUsers(db, 3822333, "admin@flag.com")

}
func getUserEmail(id string) string {
	var email string
	query := "select email from users where id = " + id
	fmt.Println("query = ", query)
	selQuery := GlobalDB.QueryRow(query) //sqli  here
	selQuery.Scan(&email)
	return email
}
func sqli_handler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("html/sqli.html")
	queryID, ok := r.URL.Query()["id"]
	if !ok || len(queryID) == 0 {
		fmt.Fprintf(w, "Missing the ID GET param, please use ?id= to get User email")
	} else {
		userEmail := getUserEmail(queryID[0])
		t.Execute(w, userEmail)
	}

}

func index_handler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("html/index.html")
	t.Execute(w, nil)
}

func testquery() {
	selusers := GlobalDB.QueryRow("select email from users where id = 2")
	var email string
	selusers.Scan(&email)
	fmt.Println(email)
}
func main() {
	db, err := sql.Open("sqlite3", "file:vulndb.db?cache=shared")
	if err != nil {
		fmt.Println(err)
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	err = db.Ping()
	if err != nil {
		fmt.Printf("%v", err)
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	GlobalDB = *db
	dbCreate(db)

	testquery()
	fmt.Println("~~~~ Starting Web Server ~~~~")
	http.HandleFunc("/", index_handler)
	http.HandleFunc("/sqli", sqli_handler)
	http.ListenAndServe(":9001", nil)

}
