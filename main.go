// +build !headless

package main

import (
        "fmt"
        "log"
        "os"

        "zabbix-manager/config"
        "zabbix-manager/ui"
)

func main() {
        // Configurar diretório de logs
        configurarLogs()

        // Iniciar a aplicação
        log.Println("Iniciando Zabbix Manager...")
        
        // Esta versão usa GUI e será compilada apenas quando a tag 'headless' não estiver presente
        app := ui.NovaAplicacao()
        app.Iniciar()
}

func configurarLogs() {
        // Obter diretório home do usuário
        diretorioHome, err := os.UserHomeDir()
        if err != nil {
                // Usar diretório atual se não conseguir obter o home
                diretorioHome, _ = os.Getwd()
        }

        // Criar diretório de logs
        diretorioLogs := fmt.Sprintf("%s/.zabbix-manager/logs", diretorioHome)
        err = os.MkdirAll(diretorioLogs, 0755)
        if err != nil {
                log.Printf("Erro ao criar diretório de logs: %v", err)
                return
        }

        // Criar arquivo de log
        arquivoLog, err := os.OpenFile(
                fmt.Sprintf("%s/zabbix-manager.log", diretorioLogs),
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