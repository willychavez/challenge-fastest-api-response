package main

import (
	"fmt"
	"time"

	"github.com/willychavez/challenge-fastest-api-response/internal/api"
)

func main() {
	cep := "01153000"
	ch := make(chan api.APIResponse)

	go api.FetchBrasilAPI(cep, ch)
	go api.FetchViaCEP(cep, ch)

	select {
	case res := <-ch:
		if res.Error != nil {
			fmt.Printf("Error %s: %v\n", res.API, res.Error)
		} else {
			fmt.Printf("Result from %s: %+v\n", res.API, res.Address)
		}
	case <-time.After(1 * time.Second):
		fmt.Println("Error: Timeout exceeded.")
	}
}
