package zabbix

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// GerarRelatorioCSV exporta os dados de hosts para um arquivo CSV
func GerarRelatorioCSV(hosts []Host, caminhoArquivo string) error {
	// Criar diretório pai, se necessário
	diretorioPai := filepath.Dir(caminhoArquivo)
	if err := os.MkdirAll(diretorioPai, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório: %v", err)
	}

	// Criar o arquivo
	arquivo, err := os.Create(caminhoArquivo)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo: %v", err)
	}
	defer arquivo.Close()

	// Configurar escritor CSV
	escritor := csv.NewWriter(arquivo)
	defer escritor.Flush()

	// Escrever cabeçalho
	cabecalho := []string{
		"ID", "Nome", "Status", "Quantidade de Itens", "Quantidade de Triggers",
		"Data de Exportação",
	}
	if err := escritor.Write(cabecalho); err != nil {
		return fmt.Errorf("erro ao escrever cabeçalho: %v", err)
	}

	// Escrever dados para cada host
	for _, host := range hosts {
		// Obter status como texto
		status := StatusHost[host.Status]
		if status == "" {
			status = "Desconhecido"
		}

		// Criar linha
		linha := []string{
			host.ID,
			host.Nome,
			status,
			fmt.Sprintf("%d", len(host.Items)),
			fmt.Sprintf("%d", len(host.Triggers)),
			time.Now().Format("2006-01-02 15:04:05"),
		}

		// Escrever linha
		if err := escritor.Write(linha); err != nil {
			return fmt.Errorf("erro ao escrever linha: %v", err)
		}
	}

	return nil
}