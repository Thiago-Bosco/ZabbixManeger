// +build !headless

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"zabbix-manager/ui/gio"
)

func main() {
	// Configurar diretório de logs
	configurarLogs()

	// Iniciar a aplicação Gio
	log.Println("Iniciando Zabbix Manager (GUI)...")
	app := gio.NovaAplicacao()
	if err := app.Executar(); err != nil {
		fmt.Printf("Erro ao executar aplicação: %v\n", err)
		os.Exit(1)
	}
}

func configurarLogs() {
	// Obter diretório home do usuário
	diretorioHome, err := os.UserHomeDir()
	if err != nil {
		// Usar diretório atual se não conseguir obter o home
		diretorioHome, _ = os.Getwd()
	}

	// Criar diretório de logs
	diretorioLogs := filepath.Join(diretorioHome, ".zabbix-manager", "logs")
	err = os.MkdirAll(diretorioLogs, 0755)
	if err != nil {
		log.Printf("Erro ao criar diretório de logs: %v", err)
		return
	}

	// Criar arquivo de log
	arquivoLog, err := os.OpenFile(
		filepath.Join(diretorioLogs, "zabbix-manager.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		log.Printf("Erro ao criar arquivo de log: %v", err)
		return
	}

	// Configurar log para escrever no arquivo e também na saída padrão
	log.SetOutput(os.Stdout)
}