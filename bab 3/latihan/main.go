package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = "5432"
	username = "postgres"
	password = "jakarta2017"
	dbname   = "challengeTujuh"
	dialect  = "postgres"
)

type Book struct {
	Id        int       `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type BookRequest struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}

var (
	db  *sql.DB
	err error
)

const PORT = ":4000"

func init() {
	psglinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, username, password, dbname)

	db, err = sql.Open(dialect, psglinfo)

	if err != nil {
		log.Panicf("error connecting to db: %s", err.Error())
	}
	_ = db

	err = db.Ping()
	if err != nil {
		log.Panicf("erorr while verify connection to database books: %s", err.Error())
	}

	createBooksTableQuery := `
	CREATE TABLE IF NOT EXISTS "books" (
		id SERIAL PRIMARY KEY, 
		title varchar(255) NOT NULL, 
		author  varchar(255) NOT NULL, 
		createdAt timestamptz DEFAULT now(),
		updatedAt timestamptz DEFAULT now()
	);`

	if err != nil {
		log.Panicf("cannot creating database: %s", err.Error())
	}
	_, err = db.Exec(createBooksTableQuery)
	fmt.Println("sukes")
}

func main() {
	http.HandleFunc("/books", booksEndpoint)
	fmt.Println("Listening on PORT:", PORT)
	http.ListenAndServe(PORT, nil)
}

func booksEndpoint(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if r.URL.Query().Get("id") != "" {
			getOneBook(w, r)
			return
		}
		getBook(w, r)
		return
	}
	if r.Method == http.MethodPut {
		updatedBooksById(w, r)
		return
	}
	if r.Method == http.MethodPost {
		createBook(w, r)
		return
	}
	if r.Method == http.MethodDelete {
		deleteBook(w, r)
		return
	}

}

func getBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	getBooksQuery := `
	SELECT id, title, author, createdAt, updatedAt from "books"
	`

	rows, err := db.Query(getBooksQuery)

	if err != nil {
		log.Panicf("error has been occured when you trying books data please check your selected book id: %s \n", err.Error())
	}

	defer rows.Close()

	books := []Book{}

	for rows.Next() {
		var book Book

		err = rows.Scan(&book.Id, &book.Title, &book.Author, &book.CreatedAt, &book.UpdatedAt)

		if err != nil {
			log.Panicf("error while scanning books data : %s \n", err.Error())
		}

		books = append(books, book)
	}

	bs, err := json.Marshal(books)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("something went wrong please contact admin "))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bs)

}

func getOneBook(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	book := Book{}
	err = db.QueryRow("SELECT id, title, author FROM books WHERE id=$1", id).Scan(&book.Id, &book.Title, &book.Author)
	if err == sql.ErrNoRows {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func errorHandler(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}

func responseHandler(w http.ResponseWriter, statusCode int, data interface{}) {

	bs, err := json.Marshal(data)

	if err != nil {
		errorHandler(w, http.StatusInternalServerError, "something when wrong please contact admin")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	w.Write(bs)

}

func updatedBooksById(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	idResult, err := strconv.Atoi(id)

	if err != nil {
		errorHandler(w, http.StatusBadRequest, "your book id is not found please enter your id with number (int)")
		return
	}

	var bookRequest BookRequest

	body, _ := ioutil.ReadAll(r.Body)

	err = json.Unmarshal(body, &bookRequest)

	if err != nil {
		errorHandler(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	updateBookQuery :=
		`
	UPDATE "books" 
	SET title = $2, 
	author = $3, 
	updatedAt = $4
	WHERE id = $1
	RETURNING id
	`

	row := db.QueryRow(updateBookQuery, idResult, bookRequest.Title, bookRequest.Author, time.Now())

	var bookId int

	err = row.Scan(&bookId)

	if err != nil {
		if err == sql.ErrNoRows {
			errorHandler(w, http.StatusNotFound, "your data hasn't successfully updated, because the book id is not found")
			return
		}
		errorHandler(w, http.StatusInternalServerError, "your data hasn't successfully updated please try again")
		return
	}

	message := fmt.Sprintf("your book with id %d has been succefully updated", bookId)

	responseHandler(w, http.StatusOK, message)
}

func createBook(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-Type") != "text/plain" {
		w.WriteHeader(400)

		fmt.Fprintf(w, "invalid content type please use JSON (raw) for post your books data")
		return
	}

	var bookRequest BookRequest
	body, _ := ioutil.ReadAll(r.Body)

	err = json.Unmarshal(body, &bookRequest)

	if err != nil {
		errorHandler(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	createBooksQuery := `
	INSERT INTO "books"
	(
		title, 
		author
	)
	VALUES($1, $2)
	RETURNING id, title, author, createdAt, updatedAt
	`
	row := db.QueryRow(createBooksQuery, bookRequest.Title, bookRequest.Author)

	var response map[string]string = map[string]string{
		"message": "your books detail has been succesfully added",
	}

	bs, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, "please fill the data correctly or internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(bs)
	var newBook Book
	err = row.Scan(&newBook.Id, &newBook.Title, &newBook.Author, &newBook.CreatedAt, &newBook.UpdatedAt)

	if err != nil {
		log.Panicf("error while user creating new book data :  %s\n", err.Error())
	}

}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	_, err := db.Exec("DELETE FROM books WHERE id=$1", id)
	if err != nil {
		log.Fatal(err)
		fmt.Fprint(w, "id has to be an number (integer) value")
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "your chosen books has been succesfully deleted")
}
