package app

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"coordinator/app/api/response"
	"coordinator/utils"

	"gopkg.in/yaml.v3"
)

type App struct {
	ID        string   `yaml:"id" json:"id"`
	Name      string   `yaml:"name" json:"name"`
	Type      string   `yaml:"type" json:"type"`
	PosterURL string   `yaml:"poster_url" json:"posterURL"`
	Device    []string `yaml:"device" json:"device"`
}

var appList []*App

func getAppList() ([]*App, error) {
	ymlFile, err := ioutil.ReadFile("app/api/app/apps.yml")
	if err != nil {
		return nil, err
	}

	var apps []*App
	if err = yaml.Unmarshal(ymlFile, &apps); err != nil {
		return nil, err
	}

	return apps, nil
}

func init() {
	var err error

	appList, err = getAppList()
	if err != nil {
		log.Fatalln("Couldn't read app list", err)
	}
}

type GetAppListResponse struct {
	Apps []*App `json:"apps"`
}

func GetAppList(w http.ResponseWriter, r *http.Request) {
	resp := response.Response{
		Data: GetAppListResponse{Apps: appList},
	}

	deviceParams, ok := r.URL.Query()["device"]
	if ok && len(deviceParams[0]) > 0 {
		device := deviceParams[0]

		var filteredAppList []*App
		for _, app := range appList {
			if utils.InStringSlice(app.Device, device) {
				filteredAppList = append(filteredAppList, app)
			}
		}

		resp = response.Response{
			Data: GetAppListResponse{Apps: filteredAppList},
		}
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Couldn't marshall get app list response to JSON", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}
