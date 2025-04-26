package zabbix

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// StatusHost mapeia códigos de status para descrições
var StatusHost = map[string]string{
	"0": "Ativo",
	"1": "Inativo",
}

// GerarRelatorioCSV gera um relatório CSV com informações dos hosts
func GerarRelatorioCSV(hosts []Host, caminhoArquivo string) error {
	// Garantir que o diretório existe
	dir := filepath.Dir(caminhoArquivo)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório para relatório: %w", err)
	}

	// Criar arquivo
	arquivo, err := os.Create(caminhoArquivo)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo CSV: %w", err)
	}
	defer arquivo.Close()

	// Gerar conteúdo CSV
	return gerarCSV(hosts, arquivo)
}

// GerarRelatorioCSVStream gera um relatório CSV diretamente para um io.Writer (usado na versão web)
func GerarRelatorioCSVStream(hosts []Host, writer io.Writer) error {
	return gerarCSV(hosts, writer)
}

// gerarCSV realiza a formatação e escrita do conteúdo CSV
func gerarCSV(hosts []Host, writer io.Writer) error {
	csvWriter := csv.NewWriter(writer)
	csvWriter.Comma = ';'
	defer csvWriter.Flush()

	// Escrever cabeçalhos
	cabecalhos := []string{
		"Host ID", "Host Nome", "Status", "Item 1 ID", "Item 1 Nome",
		"Item 2 ID", "Item 2 Nome", "Trigger 1 ID", "Trigger 1 Descrição",
		"Trigger 2 ID", "Trigger 2 Descrição",
	}
	if err := csvWriter.Write(cabecalhos); err != nil {
		return fmt.Errorf("erro ao escrever cabeçalhos no arquivo CSV: %w", err)
	}

	// Escrever dados
	for _, host := range hosts {
		status, ok := StatusHost[host.Status]
		if !ok {
			status = "Desconhecido"
		}

		// Processar itens (limitados a 2)
		itemCols := []string{}
		for i, item := range host.Items {
			itemCols = append(itemCols, item.ID, item.Nome)
			if i == 1 {
				break
			}
		}

		// Preencher células vazias caso necessário
		for len(itemCols) < 4 {
			itemCols = append(itemCols, "", "")
		}

		// Processar triggers (limitadas a 2)
		triggerCols := []string{}
		for i, trigger := range host.Triggers {
			triggerCols = append(triggerCols, trigger.ID, trigger.Nome)
			if i == 1 {
				break
			}
		}

		// Preencher células vazias caso necessário
		for len(triggerCols) < 4 {
			triggerCols = append(triggerCols, "", "")
		}

		// Montar linha
		linha := []string{
			host.ID,
			host.Nome,
			status,
		}
		linha = append(linha, itemCols...)
		linha = append(linha, triggerCols...)

		// Escrever linha
		if err := csvWriter.Write(linha); err != nil {
			return fmt.Errorf("erro ao escrever linha no arquivo CSV: %w", err)
		}
	}

	return nil
}