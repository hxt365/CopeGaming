package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"coordinator/app/api/app"
	"coordinator/app/api/provider"
	"coordinator/app/client"
	"coordinator/app/ws"
	"coordinator/settings"

	"github.com/rs/cors"
)

var port = flag.Int("port", 8080, "port")

func main() {
	flag.Parse()

	hub := client.NewHub()

	mux := http.NewServeMux()
	mux.HandleFunc("/apps", app.GetAppList)
	mux.HandleFunc("/providers", func(w http.ResponseWriter, r *http.Request) {
		provider.GetProviderList(hub, w, r)
	})
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(hub, w, r)
	})

	c := cors.New(cors.Options{
		AllowedOrigins: settings.AllowedOrigins,
	})
	handler := c.Handler(mux)

	log.Println("Start listening on port", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), handler))
}
