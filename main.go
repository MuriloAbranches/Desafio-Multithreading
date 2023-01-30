package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	API_CEP_URL = "https://cdn.apicep.com/file/apicep/"
	VIA_CEP_URL = "http://viacep.com.br/ws/"
)

type ApiCepResponse struct {
	Code       string `json:"code"`
	State      string `json:"state"`
	City       string `json:"city"`
	District   string `json:"district"`
	Address    string `json:"address"`
	Status     int    `json:"status"`
	Ok         bool   `json:"ok"`
	StatusText string `json:"statusText"`
}

type ViaCepResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func main() {
	channelApiCep := make(chan ApiCepResponse)
	channelViaCep := make(chan ViaCepResponse)

	cep := "01310-000"

	go func() {
		url := API_CEP_URL + cep + ".json"
		
		response := request(url)
		if response != nil {
			var apiCepResponse ApiCepResponse

			if err := json.Unmarshal(response, &apiCepResponse); err != nil {
				return
			}

			channelApiCep <- apiCepResponse
		}
	}()

	go func() {
		url := VIA_CEP_URL + cep + "/json/"

		response := request(url)
		if response != nil {
			var viaCepResponse ViaCepResponse

			if err := json.Unmarshal(response, &viaCepResponse); err != nil {
				return
			}

			channelViaCep <- viaCepResponse
		}
	}()

	select {
	case response := <-channelApiCep:
		fmt.Printf("API CEP Response: %+v\n", response)

	case response := <-channelViaCep:
		fmt.Printf("Via CEP Response: %+v \n", response)

	case <-time.After(time.Second):
		fmt.Println("Request timeout")
	}
}

func request(url string) []byte {
	response, err := http.Get(url)

	if err != nil {
		return nil
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil
	}

	return body
}
