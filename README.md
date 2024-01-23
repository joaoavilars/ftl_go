# GeraZip

GeraZip é uma aplicação em Go (Golang) que realiza requisições HTTP a um serviço web para gerar arquivos ZIP com base em dados fornecidos em um arquivo CSV.

## Pré-requisitos

Certifique-se de ter o Go instalado no seu ambiente. Para mais informações, visite: https://golang.org/doc/install

## Como Usar

1. Baixe o código-fonte do projeto:

   ```bash
   git clone https://github.com/seu-usuario/GeraZip.git
   cd GeraZip
```
2. Compile o código:

- Windows
```GOARCH=amd64 GOOS=windows go build -o gerazip.exe```

3. Execute o programa, fornecendo o caminho para o arquivo CSV:

```./gerazip arquivo.csv```

# Opções

-u: Especifica a URL base para o serviço web. Padrão: https://localhost
-h: Exibe a mensagem de ajuda.

Exemplo de uso com URL personalizada:
```./gerazip -u https://localhost:442 arquivo.csv```

# Padrão do arquivo CSV
```CNPJ,DATA_INICIO,DATA_FIM```

Ex.:
```11111111111111,01012020,31012020```
