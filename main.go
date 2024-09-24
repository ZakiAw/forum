package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

var temp *template.Template

type User struct {
	ID       int
	Email    string
	Username string
	Password string
}

func main() {
	os.Mkdir("SQL", 0o755)

	db, err := sql.Open("sqlite3", "./SQL/data.db")
	if err != nil {
		fmt.Println("Error Opening database:", err)
		os.Exit(0)
	}

	tables(db)
	handle()
}

func tables(db *sql.DB) {
	tab := `
	CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	email TEXT NOT NULL UNIQUE,
	username TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL,
	);
	CREATE TABLE IF NOT EXISTS comments (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	post_id INTEGER NOT NULL,
	username TEXT NOT NULL,
	content TEXT NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS likes (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	post_id INTEGER NOT NULL,
	username TEXT NOT NULL,
	UNIQUE(post_id,username)
	);
	CREATE TABLE IF NOT EXISTS commentsLikes(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	comment_id INTEGER NOT NULL,
	username TEXT NOT NULL,
	UNIQUE(comment_id,username)
	);`
	_, err := db.Exec(tab)
	if err != nil {
		log.Fatal(err)
	}
}

func handle() {
	fs := http.FileServer(http.Dir("style"))
	http.Handle("/style/", http.StripPrefix("/style/", fs))
	temp, _ = template.ParseGlob("html/*.html")
	http.HandleFunc("/", mainHandler)
	http.ListenAndServe(":8080", nil)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	temp.ExecuteTemplate(w, "register.html", nil)
}
