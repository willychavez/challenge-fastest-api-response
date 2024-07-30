package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AddressBrasilAPI struct {
	CEP          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
}

type APIResponse struct {
	Address interface{}
	API     string
	Error   error
}

var BrasilAPIURL = "https://brasilapi.com.br/api/cep/v1/"

func FetchBrasilAPI(cep string, ch chan<- APIResponse) {
	client := &http.Client{Timeout: time.Second * 1}
	resp, err := client.Get(fmt.Sprintf("%s%s", BrasilAPIURL, cep))
	if err != nil {
		ch <- APIResponse{API: "BrasilAPI", Error: err}
		return
	}
	defer resp.Body.Close()

	var address AddressBrasilAPI
	if err := json.NewDecoder(resp.Body).Decode(&address); err != nil {
		ch <- APIResponse{API: "BrasilAPI", Error: err}
		return
	}
	ch <- APIResponse{API: "BrasilAPI", Address: address}
}
