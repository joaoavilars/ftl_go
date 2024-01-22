package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const baseURL = "https://localhost"

func main() {
	// Argumento (arquivo CSV)
	csvFile := os.Args[1]

	// Abrir o arquivo CSV
	file, err := os.Open(csvFile)
	if err != nil {
		fmt.Printf("Erro ao abrir arquivo CSV: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Criar um leitor CSV
	reader := csv.NewReader(file)

	// Desabilitar a verificação do certificado TLS
	http.DefaultTransport.(*http.Transport).TLSClientConfig.InsecureSkipVerify = true

	// Obter a data e hora atual para o nome do arquivo de log
	currentTime := time.Now()
	logFileName := fmt.Sprintf("log_%s.txt", currentTime.Format("20060102150405"))

	// Criar ou abrir o arquivo de log
	logFile, err := os.Create(logFileName)
	if err != nil {
		fmt.Printf("Erro ao criar ou abrir arquivo de log: %v\n", err)
		os.Exit(1)
	}
	defer logFile.Close()

	// Criar um escritor para o arquivo de log
	logWriter := io.MultiWriter(os.Stdout, logFile)

	// Loop pelos registros do CSV
	for {
		// Ler o próximo registro do CSV
		record, err := reader.Read()

		// Verificar se atingiu o final do arquivo
		if err == io.EOF {
			break
		}

		// Verificar se ocorreu algum outro erro
		if err != nil {
			fmt.Printf("Erro ao ler registro do CSV: %v\n", err)
			os.Exit(1)
		}

		// Extrair as informações do registro
		cnpj := record[0]
		dataIni := record[1]
		dataFim := record[2]

		// Formatar a data inicial para o formato 'ddMMyyyy'
		dataIniFormatted, err := time.Parse("02012006", dataIni)
		if err != nil {
			fmt.Printf("Erro ao formatar a data inicial: %v\n", err)
			continue
		}
		dataIniFormattedStr := dataIniFormatted.Format("02012006")

		// Formatar a data final para o formato 'ddMMyyyy'
		dataFimFormatted, err := time.Parse("02012006", dataFim)
		if err != nil {
			fmt.Printf("Erro ao formatar a data final: %v\n", err)
			continue
		}
		dataFimFormattedStr := dataFimFormatted.Format("02012006")

		// URL do serviço web
		url := fmt.Sprintf("%s/GenerateZipXml.ashx?dataini=%s&datafim=%s&cnpj=%s", baseURL, dataIniFormattedStr, dataFimFormattedStr, cnpj)

		// Imprimir a URL no arquivo de log
		fmt.Fprintln(logWriter, url)

		// Fazer a requisição HTTP
		response, err := http.Get(url)
		if err != nil {
			fmt.Printf("Erro ao fazer a requisição para CNPJ %s: %v\n", cnpj, err)
			fmt.Fprintln(logWriter, fmt.Sprintf("Erro para CNPJ %s: %v", cnpj, err))
			continue
		}
		defer response.Body.Close()

		// Ler o conteúdo da resposta
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("Erro ao ler o conteúdo da resposta para CNPJ %s: %v\n", cnpj, err)
			fmt.Fprintln(logWriter, fmt.Sprintf("Erro ao ler resposta para CNPJ %s: %v", cnpj, err))
			continue
		}

		// Imprimir o resultado no arquivo de log
		fmt.Fprintln(logWriter, fmt.Sprintf("CNPJ: %s - de %s a %s : %s", cnpj, dataIniFormattedStr, dataFimFormattedStr, body))
		fmt.Fprintln(logWriter, "")
	}
}
