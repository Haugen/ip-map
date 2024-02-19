package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"slices"
	"time"

	"github.com/a-h/templ"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

var baseIpApiUrl = "http://ip-api.com/json/"
var myClient = &http.Client{Timeout: 10 * time.Second}

type FormData struct {
	URL string
}

type Location struct {
	Lat, Lon float64
}

type ResponseData struct {
	ip           string
	responseTime int64
	location     Location
}

type IpApiResponse struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
	Query       string  `json:"query"`
}

func main() {
	app := Homepage("Ip Map")

	http.Handle("/", templ.Handler(app))

	http.HandleFunc("POST /handle-form", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()

		details := FormData{
			URL: r.FormValue("url"),
		}

		now := time.Now()
		networkData := make([]ResponseData, 0)
		chromedp.ListenTarget(ctx, func(ev interface{}) {
			if ev, ok := ev.(*network.EventResponseReceived); ok {
				if ev.Response.RemoteIPAddress != "" {
					networkData = append(networkData, ResponseData{
						ev.Response.RemoteIPAddress,
						ev.Response.ResponseTime.Time().Sub(now).Milliseconds(),
						Location{},
					})
					fmt.Println(ev.Response.RemoteIPAddress, ev.Response.ResponseTime.Time().Sub(now).Milliseconds())
				}
			}
		})

		if err := chromedp.Run(ctx, chromedp.Navigate(details.URL)); err != nil {
			log.Fatal(err)
		}

		type CheckedAddresses struct {
			ip       string
			Location Location
		}

		checkedAddresses := make([]CheckedAddresses, 0)
		for i, data := range networkData {
			idx := slices.IndexFunc(checkedAddresses, func(c CheckedAddresses) bool { return c.ip == data.ip })

			if idx != -1 {
				networkData[i].location = checkedAddresses[idx].Location
				continue
			}

			resp, err := myClient.Get(baseIpApiUrl + data.ip)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()

			resBody, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}

			var ipApiResponse IpApiResponse
			if err := json.Unmarshal(resBody, &ipApiResponse); err != nil {
				log.Fatal(err)
			}

			address := CheckedAddresses{
				ip: data.ip,
				Location: Location{
					ipApiResponse.Lat,
					ipApiResponse.Lon,
				},
			}
			checkedAddresses = append(checkedAddresses, address)

			networkData[i].location = Location{
				ipApiResponse.Lat,
				ipApiResponse.Lon,
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(networkData)
	})

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", nil)
}
