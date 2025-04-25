package main

import (
	"log"

	"fyne.io/fyne/v2/app"

	"./ui"
)

func main() {
	// Criar a aplicação Fyne
	aplicacao := app.New()
	
	// Configurar o título da aplicação
	aplicacao.SetIcon(ui.CarregarLogoZabbix())
	
	// Iniciar a interface da aplicação
	gerenciadorUI := ui.NovoGerenciadorUI(aplicacao)
	if err := gerenciadorUI.Iniciar(); err != nil {
		log.Fatalf("Erro ao iniciar a aplicação: %v", err)
	}
}
