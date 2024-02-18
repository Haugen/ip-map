package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
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
		ctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()

		details := FormData{
			URL: r.FormValue("url"),
		}

		now := time.Now()
		chromedp.ListenTarget(ctx, func(ev interface{}) {
			if ev, ok := ev.(*network.EventResponseReceived); ok {
				if ev.Response.RemoteIPAddress != "" {
					fmt.Println(ev.Response.RemoteIPAddress, ev.Response.ResponseTime.Time().Sub(now).Milliseconds())
				}
			}
		})

		if err := chromedp.Run(ctx, chromedp.Navigate(details.URL)); err != nil {
			log.Fatal(err)
		}

		fmt.Println("")
		fmt.Println("#########")
		fmt.Println("")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(details)
	})

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", nil)
}
