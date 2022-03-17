package provider

import (
	"coordinator/app/api/response"
	"encoding/json"
	"log"
	"net/http"

	"coordinator/app/client"
)

type Provider struct {
	ID string `json:"id"`
}

type GetProviderListResp struct {
	Providers []*Provider `json:"providers"`
}

func GetProviderList(hub *client.Hub, w http.ResponseWriter, r *http.Request) {
	var providers []*Provider

	for _, p := range hub.GetProviders() {
		providers = append(providers, &Provider{ID: p.ID})
	}

	resp := response.Response{
		Data: GetProviderListResp{Providers: providers},
	}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Couldn't marshall get provider list response to JSON", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}
