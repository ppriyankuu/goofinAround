package main

import (
	"fmt"
	"net/http"
)

func main() {
	port := "8081"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Response from server %s", port)
	})
	http.ListenAndServe(":"+port, nil)
}
