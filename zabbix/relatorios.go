package zabbix

import (
        "encoding/csv"
        "fmt"
        "os"
        "time"
)

// GerarRelatorioCSV gera um relatório CSV contendo informações dos hosts
func GerarRelatorioCSV(hosts []Host, caminhoArquivo string) error {
        // Criar o arquivo
        arquivo, err := os.Create(caminhoArquivo)
        if err != nil {
                return fmt.Errorf("erro ao criar arquivo: %w", err)
        }
        defer arquivo.Close()

        // Criar o escritor CSV
        escritor := csv.NewWriter(arquivo)
        defer escritor.Flush()

        // Escrever cabeçalho
        cabecalho := []string{
                "ID do Host",
                "Nome do Host",
                "Status",
                "Quantidade de Itens",
                "Quantidade de Triggers",
                "Data de Exportação",
        }
        if err := escritor.Write(cabecalho); err != nil {
                return fmt.Errorf("erro ao escrever cabeçalho: %w", err)
        }

        // Obter a data e hora atual
        dataAtual := time.Now().Format("2006-01-02 15:04:05")

        // Escrever dados de cada host
        for _, host := range hosts {
                status := StatusHost[host.Status]
                if status == "" {
                        status = "Desconhecido"
                }

                linha := []string{
                        host.ID,
                        host.Nome,
                        status,
                        fmt.Sprintf("%d", len(host.Items)),
                        fmt.Sprintf("%d", len(host.Triggers)),
                        dataAtual,
                }

                if err := escritor.Write(linha); err != nil {
                        return fmt.Errorf("erro ao escrever linha: %w", err)
                }
        }

        return nil
}