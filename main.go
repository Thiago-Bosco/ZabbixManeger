package main

import (
        "log"
        "path/filepath"

        "fyne.io/fyne/v2/app"

        "zabbix-manager/config"
        "zabbix-manager/ui"
)

func main() {
        // Inicializar a configuração
        caminhoConfig := filepath.Join(".", "config.json")
        configuracao := config.Nova(caminhoConfig)

        // Criar a aplicação Fyne
        aplicacao := app.New()
        
        // Configurar o ícone da aplicação
        aplicacao.SetIcon(ui.CarregarLogoZabbix())
        
        // Iniciar a interface da aplicação
        gerenciadorUI := ui.NovoGerenciadorUI(aplicacao, configuracao)
        if err := gerenciadorUI.Iniciar(); err != nil {
                log.Fatalf("Erro ao iniciar a aplicação: %v", err)
        }
}
