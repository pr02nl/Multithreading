package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	var cep string
	_, err := fmt.Scan(&cep)
	if err != nil {
		fmt.Printf("Erro ao ler o CEP, %s", err)
	}
	apicep := make(chan string)
	viacep := make(chan string)
	go cepByApi("https://cdn.apicep.com/file/apicep/%s.json", cep, apicep)
	go cepByApi("https://viacep.com.br/ws/%s/json/", cep, viacep)
	select {
	case result := <-apicep:
		fmt.Printf("Resultado via cdn.apicep.com: %s\n", result)
	case result := <-viacep:
		fmt.Printf("Resultado via viacep.com.br: %s\n", result)
	case <-time.After(time.Second):
		fmt.Println("Timeout")
	}
}

func cepByApi(api, cep string, result chan<- string) {
	url := fmt.Sprintf(api, cep)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Erro ao buscar o CEP, %s", err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Erro ao ler o CEP, %s", err)
		return
	}
	result <- string(body)
}
