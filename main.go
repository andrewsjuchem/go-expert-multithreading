package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type ViaCEP struct {
	Cep          string `json:"cep"`
	Street       string `json:"logradouro"`
	Street2      string `json:"complemento"`
	Neighborhood string `json:"bairro"`
	City         string `json:"localidade"`
	Uf           string `json:"uf"`
	Ibge         string `json:"ibge"`
	Gia          string `json:"gia"`
	Ddd          string `json:"ddd"`
	Siafi        string `json:"siafi"`
	Error        bool   `json:"erro"`
}

type BrasilApiCep struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
	Message      string `json:"message"`
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your Postal Code (CEP): ")
	cepParam, _ := reader.ReadString('\n')
	// cepParam := "93520-575"
	cepParam = strings.TrimSpace(cepParam)
	fmt.Printf("Your postal code is %s\n", cepParam)

	viaCepChannel := make(chan ViaCEP)
	brasilApiCepChannel := make(chan BrasilApiCep)
	var responseString string

	go FindPostalCodeViaCEP(cepParam, viaCepChannel)
	go FindPostalCodeBrasilApiCEP(cepParam, brasilApiCepChannel)

	select {
	case viaCep := <-viaCepChannel:
		if viaCep.Error {
			responseString = "Invalid postal code (CEP)"
		} else {
			responseString = viaCep.Street + ", " + viaCep.Neighborhood + ", " + viaCep.City + ", " + viaCep.Uf + ", " + viaCep.Cep
		}
		fmt.Println("Source: ViaCEP")
		fmt.Println("Address: " + responseString)

	case brasilApiCep := <-brasilApiCepChannel:
		if brasilApiCep.Message != "" {
			responseString = "Invalid postal code (CEP)"
		} else {
			responseString = brasilApiCep.Street + ", " + brasilApiCep.Neighborhood + ", " + brasilApiCep.City + ", " + brasilApiCep.State + ", " + brasilApiCep.Cep
		}
		fmt.Println("Source: BrasilApiCep")
		fmt.Println("Address: " + responseString)

	case <-time.After(time.Second * 1):
		fmt.Println("Timeout: the APIs took more than one second to respond")
	}
}

func FindPostalCodeViaCEP(cep string, viaCepChannel chan ViaCEP) {
	resp, error := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
	if error != nil {
		log.Println(error)
		panic(error)
	}
	defer resp.Body.Close()
	body, error := io.ReadAll(resp.Body)
	if error != nil {
		log.Println(error)
		panic(error)
	}
	var c ViaCEP
	error = json.Unmarshal(body, &c)
	if error != nil {
		log.Println(error)
		panic(error)
	}
	viaCepChannel <- c
}

func FindPostalCodeBrasilApiCEP(cep string, brasilApiCepChannel chan BrasilApiCep) {
	resp, error := http.Get("https://brasilapi.com.br/api/cep/v1/" + cep)
	if error != nil {
		log.Println(error)
		panic(error)
	}
	defer resp.Body.Close()
	body, error := io.ReadAll(resp.Body)
	if error != nil {
		log.Println(error)
		panic(error)
	}
	var c BrasilApiCep
	error = json.Unmarshal(body, &c)
	if error != nil {
		log.Println(error)
		panic(error)
	}
	brasilApiCepChannel <- c
}
