package main

import (
	"crypto/tls"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var (
	baseURLFlag = flag.String("u", "https://localhost", "URL base para o serviço web")
	helpFlag    = flag.Bool("h", false, "Mostra a mensagem de ajuda")
)

func main() {
	// Flags
	flag.Parse()

	// Verificar se a opção de ajuda foi fornecida
	if *helpFlag {
		printUsage()
		return
	}

	// Verificar se há argumentos suficientes
	if flag.NArg() < 1 {
		fmt.Println("Erro: Forneça o caminho do arquivo CSV como argumento. CSV separado por virgula e sem aspas.")
		fmt.Println("Formato CSV: CNPJ,DATA_INICIO,DATA_FIM - as datas devem estar no layout ddmmyyyy.")
		printUsage()
		os.Exit(1)
	}

	// Argumento não processado (arquivo CSV)
	csvFile := flag.Arg(0)

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
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

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

		// Construir a URL do serviço web
		url := fmt.Sprintf("%s/GenerateZipXml.ashx?dataini=%s&datafim=%s&cnpj=%s", *baseURLFlag, dataIniFormattedStr, dataFimFormattedStr, cnpj)

		// Fazer a requisição HTTP
		response, err := http.Get(url)
		if err != nil {
			fmt.Printf("Erro ao fazer a requisição para CNPJ %s: %v\n", cnpj, err)
			continue
		}
		defer response.Body.Close()

		// Ler o conteúdo da resposta
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("Erro ao ler o conteúdo da resposta para CNPJ %s: %v\n", cnpj, err)
			continue
		}

		fmt.Println(url)
		fmt.Printf("CNPJ: %s - de %s a %s : %s\n", cnpj, dataIniFormattedStr, dataFimFormattedStr, body)
		fmt.Println()
	}
}

func printUsage() {
	fmt.Println("Uso:")
	fmt.Printf("  %s [-u <URL>] arquivo_csv\n", os.Args[0])
	fmt.Println("Opções:")
	flag.PrintDefaults()
}
