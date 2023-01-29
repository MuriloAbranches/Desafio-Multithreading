package main

import (
	"io/ioutil"
	"net/http"
	"time"
)

const (
	API_CEP_URL = "https://cdn.apicep.com/file/apicep/"
	VIA_CEP_URL = "http://viacep.com.br/ws/"
)

func main() {
	channelApiCep := make(chan string)
	channelViaCep := make(chan string)

	cep := "01310-000"

	go func() {
		url := API_CEP_URL + cep + ".json"

		response := request(url)
		if response != nil {
			channelApiCep <- string(response)
		}
	}()

	go func() {
		url := VIA_CEP_URL + cep + "/json/"

		response := request(url)
		if response != nil {
			channelViaCep <- string(response)
		}
	}()

	select {
	case response := <-channelApiCep:
		println(response)

	case response := <-channelViaCep:
		println(response)

	case <-time.After(time.Second):
		println("Request timeout")
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
