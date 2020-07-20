package main

import (
	"encoding/json"
	"log"
	"net/http"

	libs "./libs"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/quotes/getall", libs.BasicAuthMiddleware(http.HandlerFunc(getAll))).Methods("GET")
	r.HandleFunc("/api/quotes/liked/{id}", libs.BasicAuthMiddleware(http.HandlerFunc(addLikesOnAQuote))).Methods("PUT")
	r.HandleFunc("/api/quotes/undolike/{id}", libs.BasicAuthMiddleware(http.HandlerFunc(removeLikesOnAQuote))).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8000", r))

}

func getAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(libs.GetAllQuotes(libs.DbConnection()))
}

func addLikesOnAQuote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id := params["id"]
	json.NewEncoder(w).Encode(libs.UpdateLikes(libs.DbConnection(), 1, id))
}

func removeLikesOnAQuote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id := params["id"]
	json.NewEncoder(w).Encode(libs.UpdateLikes(libs.DbConnection(), -1, id))
}
