package provider

import (
	"encoding/json"
	"log"
	"net/http"

	"coordinator/app/api/response"
	"coordinator/app/client"
)

type Provider struct {
	ID         string  `json:"id"`
	HostName   string  `json:"hostName"`
	Platform   string  `json:"platform"`
	CpuName    string  `json:"cpuName"`
	CpuNum     int     `json:"cpuNum"`
	MemSize    float64 `json:"memSize"`
	CpuPercent float64 `json:"cpuPercent"`
	MemPercent float64 `json:"memPercent"`
}

type GetProviderListResp struct {
	Providers []*Provider `json:"providers"`
}

func GetProviderList(hub *client.Hub, w http.ResponseWriter, r *http.Request) {
	hasOwnerIDParam := r.URL.Query().Has("owner")
	ownerID := r.URL.Query().Get("owner")

	providers := make([]*Provider, 0)

	for _, p := range hub.GetProviders() {
		if !hasOwnerIDParam || p.Provider.OwnerID == ownerID {
			providers = append(providers, &Provider{
				ID:         p.ID,
				HostName:   p.Provider.HostName,
				Platform:   p.Provider.Platform,
				CpuName:    p.Provider.CpuName,
				CpuNum:     p.Provider.CpuNum,
				MemSize:    p.Provider.MemSize,
				CpuPercent: p.Provider.CpuPercent,
				MemPercent: p.Provider.MemPercent,
			})
		}
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
