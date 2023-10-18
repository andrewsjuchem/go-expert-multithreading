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
	cepParam := "93125-486"
	cepParam = strings.TrimSpace(cepParam)
	fmt.Printf("Your postal code is %s\n", cepParam)

	viaCep, error := FindPostalCodeViaCEP(cepParam)
	if error != nil {
		log.Println(error)
		fmt.Fprintf(os.Stderr, "Error while parsing the response: %v\n", error)
	}
	jsonDataIndent, error := json.MarshalIndent(viaCep, "", "    ")
	if error != nil {
		log.Println(error)
		return
	}
	fmt.Println(string(jsonDataIndent))

	apiCep, error := FindPostalCodeBrasilApiCEP(cepParam)
	if error != nil {
		fmt.Fprintf(os.Stderr, "Error while parsing the response: %v\n", error)
	}
	jsonDataIndent, error = json.MarshalIndent(apiCep, "", "    ")
	if error != nil {
		log.Println(error)
		return
	}
	fmt.Println(string(jsonDataIndent))
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
