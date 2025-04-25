package zabbix

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// GerarRelatorioCSV gera um arquivo CSV com os dados dos hosts
func GerarRelatorioCSV(hosts []Host, caminho string) error {
	// Se o caminho não foi especificado, usar um nome padrão
	if caminho == "" {
		dataAtual := time.Now().Format("2006-01-02_15-04-05")
		caminho = fmt.Sprintf("Relatorio_Zabbix_%s.csv", dataAtual)
	}

	// Garantir que o diretório exista
	dir := filepath.Dir(caminho)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("erro ao criar diretório para o relatório: %w", err)
		}
	}

	arquivo, err := os.Create(caminho)
	if err != nil {
		return fmt.Errorf("erro ao criar o arquivo CSV: %w", err)
	}
	defer arquivo.Close()

	writer := csv.NewWriter(arquivo)
	writer.Comma = ';'

	cabecalhos := []string{
		"Host ID", "Host Nome", "Status", "Item 1 ID", "Item 1 Nome",
		"Item 2 ID", "Item 2 Nome", "Trigger 1 ID", "Trigger 1 Descrição",
		"Trigger 2 ID", "Trigger 2 Descrição",
	}
	
	if err := writer.Write(cabecalhos); err != nil {
		return fmt.Errorf("erro ao escrever os cabeçalhos no arquivo CSV: %w", err)
	}

	for _, host := range hosts {
		status, ok := StatusHost[host.Status]
		if !ok {
			status = "Desconhecido"
		}

		itemCols := []string{}
		for i, item := range host.Items {
			itemCols = append(itemCols, item.ID, item.Nome)
			if i == 1 {
				break
			}
		}

		// Garantir que temos pelo menos 4 colunas para itens (2 itens)
		for len(itemCols) < 4 {
			itemCols = append(itemCols, "", "")
		}

		triggerCols := []string{}
		for i, trigger := range host.Triggers {
			triggerCols = append(triggerCols, trigger.ID, trigger.Nome)
			if i == 1 {
				break
			}
		}

		// Garantir que temos pelo menos 4 colunas para triggers (2 triggers)
		for len(triggerCols) < 4 {
			triggerCols = append(triggerCols, "", "")
		}

		linha := []string{
			host.ID,
			host.Nome,
			status,
		}
		linha = append(linha, itemCols...)
		linha = append(linha, triggerCols...)

		if err := writer.Write(linha); err != nil {
			return fmt.Errorf("erro ao escrever os dados no arquivo CSV: %w", err)
		}
	}

	writer.Flush()
	
	if err := writer.Error(); err != nil {
		return fmt.Errorf("erro ao finalizar o arquivo CSV: %w", err)
	}

	return nil
}
