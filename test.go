package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

// структура для записей
type Article struct {
	Id                           uint16
	Title, Description, Textarea string
}

// структура для регистрации
type ContactDetails struct {
	Login         string
	Password      string
	Success       bool
	StorageAccess string
}

var posts = []Article{}
var showPost = Article{}

func home(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("template/home_page.html", "template/footer.html", "template/header.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/go")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	res, err := db.Query("SELECT * FROM `article`")
	if err != nil {
		panic(err)
	}

	posts = []Article{}
	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Description, &post.Textarea)
		if err != nil {
			panic(err)
		}

		posts = append(posts, post)
	}

	t.ExecuteTemplate(w, "home", posts)
}

func create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("template/create.html", "template/footer.html", "template/header.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "create", nil)
}

func save_article(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	description := r.FormValue("description")
	textarea := r.FormValue("textarea")

	if title == "" || description == "" || textarea == "" {
		fmt.Fprintf(w, "nonono")
	} else {
		db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/go")
		if err != nil {
			panic(err)
		}

		defer db.Close()

		insert, err := db.Query(fmt.Sprintf("INSERT INTO `article` (`title`,`description`,`textarea`) VALUES ('%s','%s','%s')", title, description, textarea))
		if err != nil {
			panic(err)
		}

		defer insert.Close()

		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
}

func show_post(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	t, err := template.ParseFiles("template/show.html", "template/footer.html", "template/header.html")

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/go")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	res, err := db.Query(fmt.Sprintf("SELECT * FROM `article` WHERE `id` ='%s'", vars["id"]))
	if err != nil {
		panic(err)
	}

	showPost = Article{}
	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Description, &post.Textarea)
		if err != nil {
			panic(err)
		}

		showPost = post
	}

	t.ExecuteTemplate(w, "show", showPost)
}

func register(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("template/register.html", "template/footer.html", "template/header.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "register", nil)
}

func handleFunc() {
	rtr := mux.NewRouter()

	rtr.HandleFunc("/", home).Methods("GET")
	rtr.HandleFunc("/create", create).Methods("GET")
	rtr.HandleFunc("/save_article", save_article).Methods("POST")
	rtr.HandleFunc("/post/{id:[0-9]+}", show_post).Methods("GET")
	rtr.HandleFunc("/register", register).Methods("GET")

	http.Handle("/", rtr)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.ListenAndServe(":8080", nil)
}

func main() {
	handleFunc()
}
