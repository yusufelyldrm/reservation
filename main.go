package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", reservationHandler)

	_ = http.ListenAndServe(":8080", nil)

}

func reservationHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(fmt.Sprintf("Reservation request received: %s", r.URL.Path))
}
a
