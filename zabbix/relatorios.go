package zabbix

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// StatusHost mapeia códigos de status para descrições
var StatusHost = map[string]string{
	"0": "Ativo",
	"1": "Inativo",
}

// DadosRelatorio estrutura para armazenar dados completos do relatório
type DadosRelatorio struct {
	Host              Host
	ProblemasRecentes []Problema
	Disponibilidade   float64
	UltimaColeta      time.Time
}

// GerarRelatorioCSV gera um relatório CSV com informações detalhadas dos hosts
func GerarRelatorioCSV(hosts []Host, caminhoArquivo string) error {
	dir := filepath.Dir(caminhoArquivo)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório: %w", err)
	}

	arquivo, err := os.Create(caminhoArquivo)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo CSV: %w", err)
	}
	defer arquivo.Close()

	return gerarCSV(hosts, arquivo)
}

// GerarRelatorioCSVStream gera relatório CSV para um io.Writer
func GerarRelatorioCSVStream(hosts []Host, writer io.Writer) error {
	return gerarCSV(hosts, writer)
}

func gerarCSV(hosts []Host, writer io.Writer) error {
	csvWriter := csv.NewWriter(writer)
	csvWriter.Comma = ';'
	defer csvWriter.Flush()

	// Cabeçalhos expandidos
	cabecalhos := []string{
		"Host ID", "Nome", "Status",
		"Disponibilidade (%)", "Última Coleta",
		"Total Items", "Items Ativos", "Items com Problema",
		"Total Triggers", "Triggers Ativas", "Triggers com Problema",
		"Problemas Últimas 24h", "Tempo Médio de Resolução",
		"Performance CPU (%)", "Performance Memória (%)",
		"Interface Principal", "Tráfego Entrada (avg)", "Tráfego Saída (avg)",
	}

	if err := csvWriter.Write(cabecalhos); err != nil {
		return fmt.Errorf("erro ao escrever cabeçalhos: %w", err)
	}

	for _, host := range hosts {
		status := StatusHost[host.Status]
		if status == "" {
			status = "Desconhecido"
		}

		// Calcular métricas
		itemsAtivos := contarItemsAtivos(host.Items)
		itemsProblema := contarItemsComProblema(host.Items)
		triggersAtivas := contarTriggersAtivas(host.Triggers)
		triggersProblema := contarTriggersComProblema(host.Triggers)

		linha := []string{
			host.ID,
			host.Nome,
			status,
			fmt.Sprintf("%.2f", calcularDisponibilidade(host)),
			obterUltimaColeta(host).Format("2006-01-02 15:04:05"),
			fmt.Sprintf("%d", len(host.Items)),
			fmt.Sprintf("%d", itemsAtivos),
			fmt.Sprintf("%d", itemsProblema),
			fmt.Sprintf("%d", len(host.Triggers)),
			fmt.Sprintf("%d", triggersAtivas),
			fmt.Sprintf("%d", triggersProblema),
			fmt.Sprintf("%d", contarProblemasRecentes(host)),
			calcularTempoMedioResolucao(host),
			fmt.Sprintf("%.2f", obterPerformanceCPU(host)),
			fmt.Sprintf("%.2f", obterPerformanceMemoria(host)),
			obterInterfacePrincipal(host),
			formatarTrafego(obterTrafegoDados(host, "in")),
			formatarTrafego(obterTrafegoDados(host, "out")),
		}

		if err := csvWriter.Write(linha); err != nil {
			return fmt.Errorf("erro ao escrever linha: %w", err)
		}
	}

	return nil
}

// Funções auxiliares
func calcularDisponibilidade(host Host) float64 {
	problemasRecentes := 0
	for _, trigger := range host.Triggers {
		if trigger.Valor == "1" { // problema ativo
			problemasRecentes++
		}
	}
	total := float64(len(host.Triggers))
	if total == 0 {
		return 100.0
	}
	return 100.0 * (1.0 - float64(problemasRecentes)/total)
}

func obterUltimaColeta(host Host) time.Time {
	ultimaColeta := time.Time{}
	for _, item := range host.Items {
		if item.UltimoValor != "" {
			timestamp, err := strconv.ParseInt(item.UltimaAlteracao, 10, 64)
			if err == nil {
				itemColeta := time.Unix(timestamp, 0)
				if itemColeta.After(ultimaColeta) {
					ultimaColeta = itemColeta
				}
			}
		}
	}
	return ultimaColeta
}

func contarItemsAtivos(items []Item) int {
	ativos := 0
	for _, item := range items {
		if item.Status == "0" { // 0 = ativo
			ativos++
		}
	}
	return ativos
}

func contarItemsComProblema(items []Item) int {
	problemas := 0
	for _, item := range items {
		if item.UltimoValor != "" && item.Estado == "1" { // 1 = problema
			problemas++
		}
	}
	return problemas
}

func contarTriggersAtivas(triggers []Trigger) int {
	ativas := 0
	for _, trigger := range triggers {
		if trigger.Status == "0" { // 0 = habilitada
			ativas++
		}
	}
	return ativas
}

func contarTriggersComProblema(triggers []Trigger) int {
	problemas := 0
	for _, trigger := range triggers {
		if trigger.Valor == "1" { // 1 = problema
			problemas++
		}
	}
	return problemas
}

func contarProblemasRecentes(host Host) int {
	count := 0
	for _, trigger := range host.Triggers {
		if trigger.Valor == "1" {
			count++
		}
	}
	return count
}

func calcularTempoMedioResolucao(host Host) string {
	var total time.Duration
	count := 0
	for _, trigger := range host.Triggers {
		if trigger.UltimaAlteracao != "" {
			timestamp, err := strconv.ParseInt(trigger.UltimaAlteracao, 10, 64)
			if err == nil {
				duracao := time.Since(time.Unix(timestamp, 0))
				total += duracao
				count++
			}
		}
	}
	if count == 0 {
		return "N/A"
	}
	media := total / time.Duration(count)
	return media.Round(time.Minute).String()
}

func obterPerformanceCPU(host Host) float64 {
	for _, item := range host.Items {
		if strings.Contains(strings.ToLower(item.Nome), "cpu") {
			if valor, err := strconv.ParseFloat(item.UltimoValor, 64); err == nil {
				return valor
			}
		}
	}
	return 0
}

func obterPerformanceMemoria(host Host) float64 {
	for _, item := range host.Items {
		if strings.Contains(strings.ToLower(item.Nome), "memory") {
			if valor, err := strconv.ParseFloat(item.UltimoValor, 64); err == nil {
				return valor
			}
		}
	}
	return 0
}

func obterInterfacePrincipal(host Host) string {
	if len(host.Interfaces) > 0 {
		return fmt.Sprintf("%s (%s)", host.Interfaces[0].IP, host.Interfaces[0].DNS)
	}
	return "N/A"
}

func obterTrafegoDados(host Host, direcao string) float64 {
	for _, item := range host.Items {
		if strings.Contains(strings.ToLower(item.Nome), fmt.Sprintf("network %s", direcao)) {
			if valor, err := strconv.ParseFloat(item.UltimoValor, 64); err == nil {
				return valor
			}
		}
	}
	return 0
}

func formatarTrafego(bytes float64) string {
	// Implementar formatação de tráfego (B, KB, MB, GB)
	return "1.02 MB/s" // Placeholder -  Needs implementation for proper formatting.
}


//Necessary structs moved to tipos.go
type Host struct {
	ID          string
	Nome        string
	Status      string
	Items       []Item
	Triggers    []Trigger
	Interfaces []Interface
}

type Item struct {
	ID             string
	Nome           string
	Status         string
	UltimoValor    string
	UltimaAlteracao string
	Estado         string
}

type Trigger struct {
	ID             string
	Nome           string
	Status         string
	Valor          string
	UltimaAlteracao string
}

type Interface struct {
	IP  string
	DNS string
}

type Problema struct {
	// ... fields ...
}