package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
)

type FormData struct {
	URL string
}

func main() {
	app := Homepage("Ip Map")

	http.Handle("/", templ.Handler(app))

	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "API response")
	})

	http.HandleFunc("POST /handle-form", func(w http.ResponseWriter, r *http.Request) {
		details := FormData{
			URL: r.FormValue("url"),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(details)
	})

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", nil)
}
