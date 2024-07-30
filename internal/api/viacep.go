package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AddressViaCEP struct {
	CEP          string `json:"cep"`
	State        string `json:"uf"`
	City         string `json:"localidade"`
	Neighborhood string `json:"bairro"`
	Street       string `json:"logradouro"`
}

var ViaCEPURL = "https://viacep.com.br/ws/"

func FetchViaCEP(cep string, ch chan<- APIResponse) {
	client := &http.Client{Timeout: time.Second * 1}
	resp, err := client.Get(fmt.Sprintf("%s%s/json", ViaCEPURL, cep))
	if err != nil {
		ch <- APIResponse{API: "ViaCEP", Error: err}
		return
	}
	defer resp.Body.Close()

	var address AddressViaCEP
	if err := json.NewDecoder(resp.Body).Decode(&address); err != nil {
		ch <- APIResponse{API: "ViaCEP", Error: err}
		return
	}
	ch <- APIResponse{API: "ViaCEP", Address: address}
}
