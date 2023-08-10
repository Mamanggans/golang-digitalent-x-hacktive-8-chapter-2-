package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Product struct {
	Id int `json:"id"`

	Title string `json:"title"`

	Price int `json:"price"`
}

type ProductRequest struct {
	Title string `json:"nama"`
	Price int    `json:"harga"`
}

var Products = []Product{}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		var message []byte = []byte("Hello")

		w.Write(message)
	})

	http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {
			getProducts(w, r)
			return
		}

		if r.Method == "POST" {
			createProduct(w, r)
			return
		}

		if r.Method == "PUT" {
			updateProduct(w, r)
			return
		}

		if r.Method == "DELETE" {
			deleteProduct(w, r)
			return
		}
	})

	fmt.Println("Listening on PORT 8080")
	http.ListenAndServe(":8080", nil)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	idResult, err := strconv.Atoi(id)

	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, "id has to be an integer value")
		return
	}

	var productIndex = 0

	for index, value := range Products {
		if value.Id == idResult {
			productIndex = index
		}
	}

	copy(Products[productIndex:], Products[productIndex+1:])

	Products[len(Products)-1] = Product{}

	Products = Products[:len(Products)-1]

	fmt.Fprint(w, "product has been successfully deleted")
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	idResult, err := strconv.Atoi(id)

	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, "id has to be an integer value")
		return
	}

	var updatedProduct Product

	var productIndex = 0

	for index, value := range Products {
		if value.Id == idResult {
			updatedProduct = value
			productIndex = index
		}
	}

	if updatedProduct.Id == 0 {
		w.WriteHeader(404)
		fmt.Fprintf(w, "product with id %d does not exist", idResult)
		return
	}

	title := r.FormValue("title")

	price := r.FormValue("price")

	priceResult, err := strconv.Atoi(price)

	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, "price has to be an integer value")
		return
	}

	updatedProduct.Title = title

	updatedProduct.Price = priceResult

	Products[productIndex] = updatedProduct

	fmt.Fprintf(w, "product with id %d has been successfully updated", updatedProduct.Id)
}

func createProduct(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(400)

		fmt.Fprint(w, "invalid content-type")
		return
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(422)
		fmt.Fprintf(w, "invalid request body")
		return
	}

	var request ProductRequest

	err = json.Unmarshal(body, &request)

	// title := r.FormValue("title")

	// price := r.FormValue("price")

	// priceResult, err := strconv.Atoi(price)

	// if err != nil {
	// 	w.WriteHeader(400)
	// 	fmt.Fprint(w, "invalid price")
	// 	return
	// }

	productId := len(Products) + 1

	newProduct := Product{
		Id:    productId,
		Title: request.Title,
		Price: request.Price,
	}

	Products = append(Products, newProduct)

	var response map[string]string = map[string]string{
		"message": "New product has been successfully created",
	}

	bs, err := json.Marshal(response)

	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, "internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(201)

	w.Write(bs)
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	bs, err := json.Marshal(Products)

	if err != nil {
		w.WriteHeader(500)

		w.Write([]byte("something went wrong"))

		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(bs))
	return
}
