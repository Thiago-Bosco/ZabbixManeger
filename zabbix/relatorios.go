package zabbix

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

// GerarRelatorioCSV gera um relatório CSV com os dados dos hosts
func GerarRelatorioCSV(hosts []Host, caminhoArquivo string) error {
	log.Printf("Gerando relatório CSV em: %s", caminhoArquivo)

	// Criar o arquivo
	arquivo, err := os.Create(caminhoArquivo)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo: %w", err)
	}
	defer arquivo.Close()

	// Criar o escritor CSV
	escritor := csv.NewWriter(arquivo)
	defer escritor.Flush()

	// Escrever cabeçalhos
	cabecalhos := []string{
		"ID do Host", 
		"Nome do Host", 
		"Status", 
		"Quantidade de Itens", 
		"Quantidade de Triggers",
		"Data de Geração",
	}
	err = escritor.Write(cabecalhos)
	if err != nil {
		return fmt.Errorf("erro ao escrever cabeçalhos: %w", err)
	}

	// Obter data e hora atual
	dataHora := time.Now().Format("02/01/2006 15:04:05")

	// Escrever dados dos hosts
	for _, host := range hosts {
		// Obter o status do host como texto
		statusTexto := StatusHost[host.Status]
		if statusTexto == "" {
			statusTexto = "Desconhecido"
		}

		// Escrever uma linha para cada host
		linha := []string{
			host.ID,
			host.Nome,
			statusTexto,
			strconv.Itoa(len(host.Items)),
			strconv.Itoa(len(host.Triggers)),
			dataHora,
		}

		err = escritor.Write(linha)
		if err != nil {
			return fmt.Errorf("erro ao escrever linha: %w", err)
		}
	}

	log.Printf("Relatório gerado com sucesso: %d hosts", len(hosts))
	return nil
}

// GerarRelatorioDetalhado gera um relatório CSV detalhado com os dados dos hosts, incluindo itens e triggers
func GerarRelatorioDetalhado(hosts []Host, caminhoArquivo string) error {
	log.Printf("Gerando relatório detalhado em: %s", caminhoArquivo)

	// Criar o arquivo
	arquivo, err := os.Create(caminhoArquivo)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo: %w", err)
	}
	defer arquivo.Close()

	// Criar o escritor CSV
	escritor := csv.NewWriter(arquivo)
	defer escritor.Flush()

	// Escrever cabeçalhos
	cabecalhos := []string{
		"ID do Host", 
		"Nome do Host", 
		"Status do Host",
		"ID do Item",
		"Nome do Item",
		"Chave do Item",
		"Status do Item",
		"ID do Trigger",
		"Descrição do Trigger",
		"Status do Trigger",
		"Prioridade do Trigger",
		"Data de Geração",
	}
	err = escritor.Write(cabecalhos)
	if err != nil {
		return fmt.Errorf("erro ao escrever cabeçalhos: %w", err)
	}

	// Obter data e hora atual
	dataHora := time.Now().Format("02/01/2006 15:04:05")

	// Escrever dados detalhados
	for _, host := range hosts {
		// Obter o status do host como texto
		statusHostTexto := StatusHost[host.Status]
		if statusHostTexto == "" {
			statusHostTexto = "Desconhecido"
		}

		// Se o host não tem itens nem triggers, escrever uma linha só para o host
		if len(host.Items) == 0 && len(host.Triggers) == 0 {
			linha := []string{
				host.ID,
				host.Nome,
				statusHostTexto,
				"", // ID do Item
				"", // Nome do Item
				"", // Chave do Item
				"", // Status do Item
				"", // ID do Trigger
				"", // Descrição do Trigger
				"", // Status do Trigger
				"", // Prioridade do Trigger
				dataHora,
			}
			err = escritor.Write(linha)
			if err != nil {
				return fmt.Errorf("erro ao escrever linha: %w", err)
			}
			continue
		}

		// Processa itens
		for _, item := range host.Items {
			// Obter o status do item como texto
			statusItemTexto := StatusItem[item.Status]
			if statusItemTexto == "" {
				statusItemTexto = "Desconhecido"
			}

			// Escrever uma linha para cada item
			linha := []string{
				host.ID,
				host.Nome,
				statusHostTexto,
				item.ID,
				item.Nome,
				item.Chave,
				statusItemTexto,
				"", // ID do Trigger
				"", // Descrição do Trigger
				"", // Status do Trigger
				"", // Prioridade do Trigger
				dataHora,
			}
			err = escritor.Write(linha)
			if err != nil {
				return fmt.Errorf("erro ao escrever linha de item: %w", err)
			}
		}

		// Processa triggers
		for _, trigger := range host.Triggers {
			// Obter o status do trigger como texto
			statusTriggerTexto := StatusTrigger[trigger.Status]
			if statusTriggerTexto == "" {
				statusTriggerTexto = "Desconhecido"
			}

			// Obter a prioridade do trigger como texto
			prioridadeTriggerTexto := PrioridadeTrigger[trigger.Prioridade]
			if prioridadeTriggerTexto == "" {
				prioridadeTriggerTexto = "Desconhecida"
			}

			// Escrever uma linha para cada trigger
			linha := []string{
				host.ID,
				host.Nome,
				statusHostTexto,
				"", // ID do Item
				"", // Nome do Item
				"", // Chave do Item
				"", // Status do Item
				trigger.ID,
				trigger.Descricao,
				statusTriggerTexto,
				prioridadeTriggerTexto,
				dataHora,
			}
			err = escritor.Write(linha)
			if err != nil {
				return fmt.Errorf("erro ao escrever linha de trigger: %w", err)
			}
		}
	}

	log.Printf("Relatório detalhado gerado com sucesso: %d hosts", len(hosts))
	return nil
}