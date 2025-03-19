package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func getHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}
func getHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", getHome)
	r.HandleFunc("/hello", getHello)

	fmt.Printf("server listening at port 8080\n")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
