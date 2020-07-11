package main

// https://medium.com/the-andela-way/build-a-restful-json-api-with-golang-85a83420c9da
import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/bill", createBill).Methods("POST")
	router.HandleFunc("/bills", getAllBills).Methods("GET")
	router.HandleFunc("/bills/{id}", getOneBill).Methods("GET")
	router.HandleFunc("/bills/{id}", updateBill).Methods("PATCH")
	router.HandleFunc("/bills/{id}", deleteBill).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))

}
