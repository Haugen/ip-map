package main

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
)

func main() {
	component := Homepage("Buddy!")

	http.Handle("/", templ.Handler(component))

	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "API response")
	})

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", nil)
}
