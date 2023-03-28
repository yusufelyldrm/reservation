package main

import (
	"fmt"
	"net/http"
)

const portNumber = ":8080"

// About is the about page handler
func About(w http.ResponseWriter, r *http.Request) {

}

// Home is the home page handler
func Home(w http.ResponseWriter, r *http.Request) {
}

func main() {

	http.HandleFunc("/", Home)
	http.HandleFunc("/about", About)

	fmt.Println(fmt.Sprintf("Starting application on port %s\n Press 'Ctrl + C' to stop", portNumber))
	err := http.ListenAndServe(portNumber, nil)
	if err != nil {
		fmt.Println(err)
	}
}
