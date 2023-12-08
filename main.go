package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Dados representa a estrutura dos dados a serem ordenados
type Dados struct {
	Nome      string
	Idade     int
	Pontuacao int
}

// PorNomeIdade implementa a interface Sort para ordenar por Nome e Idade
type PorNomeIdade []Dados

func (a PorNomeIdade) Len() int      { return len(a) }
func (a PorNomeIdade) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a PorNomeIdade) Less(i, j int) bool {
	// Ordenar por Nome usando a tabela ASCII (case-insensitive)
	if strings.ToLower(a[i].Nome) != strings.ToLower(a[j].Nome) {
		return strings.ToLower(a[i].Nome) < strings.ToLower(a[j].Nome)
	}
	// Se os nomes são iguais, ordenar por Idade
	return a[i].Idade < a[j].Idade
}

// LerArquivo lê dados de um arquivo CSV
func LerArquivo(nomeArquivo string) ([]Dados, error) {
	arquivo, err := os.Open(nomeArquivo)
	if err != nil {
		return nil, err
	}
	defer arquivo.Close()

	leitorCSV := csv.NewReader(arquivo)
	leitorCSV.Comma = ','

	registros, err := leitorCSV.ReadAll()
	if err != nil {
		return nil, err
	}

	// Converter os dados do CSV para a estrutura definida
	var dados []Dados
	for _, registro := range registros[1:] {
		nome := registro[0]
		idade, err := strconv.Atoi(registro[1])
		if err != nil {
			return nil, err
		}
		pontuacao, err := strconv.Atoi(registro[2])
		if err != nil {
			return nil, err
		}

		dados = append(dados, Dados{nome, idade, pontuacao})
	}

	return dados, nil
}

// EscreverArquivo escreve dados em um arquivo CSV
func EscreverArquivo(nomeArquivo string, dados []Dados) error {
	arquivo, err := os.Create(nomeArquivo)
	if err != nil {
		return err
	}
	defer arquivo.Close()

	escritorCSV := csv.NewWriter(arquivo)
	escritorCSV.Comma = ','

	// Escrever o cabeçalho
	cabecalho := []string{"Nome", "Idade", "Pontuacao"}
	if err := escritorCSV.Write(cabecalho); err != nil {
		return err
	}

	// Escrever os dados
	for _, dado := range dados {
		if err := escritorCSV.Write([]string{dado.Nome, strconv.Itoa(dado.Idade), strconv.Itoa(dado.Pontuacao)}); err != nil {
			return err
		}
	}

	escritorCSV.Flush()
	return escritorCSV.Error()
}

// OrdenarDados ordena os dados por nome e por idade
func OrdenarDados(dados []Dados) {
	sort.Sort(PorNomeIdade(dados))
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Uso: go run main.go <arquivo-de-entrada.csv> <arquivo-de-saida.csv>")
		os.Exit(1)
	}

	arquivoEntrada := os.Args[1]
	arquivoSaida := os.Args[2]

	// Ler dados do arquivo de entrada
	dados, err := LerArquivo(arquivoEntrada)
	if err != nil {
		fmt.Printf("Erro ao ler o arquivo de entrada: %v\n", err)
		os.Exit(1)
	}

	// Ordenar os dados
	OrdenarDados(dados)

	// Escrever dados ordenados no arquivo de saída
	err = EscreverArquivo(arquivoSaida, dados)
	if err != nil {
		fmt.Printf("Erro ao escrever no arquivo de saída: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Ordenação concluída com sucesso.")
}
