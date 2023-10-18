package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
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
	cepParam := "93125-480"
	cepParam = strings.TrimSpace(cepParam)
	fmt.Printf("Your postal code is %s\n", cepParam)

	var responseString string

	// OPTION 1
	viaCep, error := FindPostalCodeViaCEP(cepParam)
	if error != nil {
		log.Println(error)
		fmt.Fprintf(os.Stderr, "Error while getting the response: %v\n", error)
	}
	if viaCep.Erro {
		responseString = "Invalid postal code (CEP)"
	} else {
		responseString = viaCep.Logradouro + ", " + viaCep.Bairro + ", " + viaCep.Localidade + ", " + viaCep.Uf + ", " + viaCep.Cep
	}
	fmt.Println(responseString)

	// OPTION 2
	apiCep, error := FindPostalCodeBrasilApiCEP(cepParam)
	if error != nil {
		log.Println(error)
		fmt.Fprintf(os.Stderr, "Error while getting the response: %v\n", error)
	}
	if apiCep.Message != "" {
		responseString = "Invalid postal code (CEP)"
	} else {
		responseString = apiCep.Street + ", " + apiCep.Neighborhood + ", " + apiCep.City + ", " + apiCep.State + ", " + apiCep.Cep
	}
	fmt.Println(responseString)
}

func FindPostalCodeViaCEP(cep string) (*ViaCEP, error) {
	resp, error := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
	if error != nil {
		return nil, error
	}
	defer resp.Body.Close()
	body, error := io.ReadAll(resp.Body)
	if error != nil {
		log.Println(error)
		return nil, error
	}
	var c ViaCEP
	error = json.Unmarshal(body, &c)
	if error != nil {
		log.Println(error)
		return nil, error
	}
	return &c, nil
}

func FindPostalCodeBrasilApiCEP(cep string) (*BrasilApiCep, error) {
	resp, error := http.Get("https://brasilapi.com.br/api/cep/v1/" + cep)
	if error != nil {
		return nil, error
	}
	defer resp.Body.Close()
	body, error := io.ReadAll(resp.Body)
	if error != nil {
		log.Println(error)
		return nil, error
	}
	var c BrasilApiCep
	error = json.Unmarshal(body, &c)
	if error != nil {
		log.Println(error)
		return nil, error
	}
	return &c, nil
}
