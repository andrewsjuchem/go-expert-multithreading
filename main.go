package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type ViaCEP struct {
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
	Erro        bool   `json:"erro"`
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
	// reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your Postal Code (CEP): ")
	// cepParam, _ := reader.ReadString('\n')
	cepParam := "93520-575"
	cepParam = strings.TrimSpace(cepParam)
	fmt.Printf("Your postal code is %s\n", cepParam)

	viaCepChannel := make(chan ViaCEP)
	brasilApiCepChannel := make(chan BrasilApiCep)
	var responseString string

	go FindPostalCodeViaCEP(cepParam, viaCepChannel)
	go FindPostalCodeBrasilApiCEP(cepParam, brasilApiCepChannel)

	select {
	case viaCep := <-viaCepChannel:
		if viaCep.Erro {
			responseString = "Invalid postal code (CEP)"
		} else {
			responseString = viaCep.Logradouro + ", " + viaCep.Bairro + ", " + viaCep.Localidade + ", " + viaCep.Uf + ", " + viaCep.Cep
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
		println("timeout")
		// default:
		// 	println("default")
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
