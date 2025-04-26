
package zabbix

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
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
	UltimaColeta     time.Time
}

// Problema representa um problema/evento do Zabbix
type Problema struct {
	ID          string
	Nome        string
	Severidade  string
	DataInicio  time.Time
	DataFim     time.Time
	Duracao     string
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
	// Implementar cálculo baseado no histórico de problemas
	return 99.99 // Placeholder
}

func obterUltimaColeta(host Host) time.Time {
	return time.Now() // Implementar busca real
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
	// Implementar contagem de problemas nas últimas 24h
	return 0 // Placeholder
}

func calcularTempoMedioResolucao(host Host) string {
	// Implementar cálculo do tempo médio de resolução
	return "1h 30m" // Placeholder
}

func obterPerformanceCPU(host Host) float64 {
	// Implementar busca de performance CPU
	return 45.5 // Placeholder
}

func obterPerformanceMemoria(host Host) float64 {
	// Implementar busca de performance memória
	return 67.8 // Placeholder
}

func obterInterfacePrincipal(host Host) string {
	// Implementar busca da interface principal
	return "eth0" // Placeholder
}

func obterTrafegoDados(host Host, direcao string) float64 {
	// Implementar busca de dados de tráfego
	return 1024.5 // Placeholder
}

func formatarTrafego(bytes float64) string {
	// Implementar formatação de tráfego (B, KB, MB, GB)
	return "1.02 MB/s" // Placeholder
}
