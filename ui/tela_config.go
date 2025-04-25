package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"strconv"
)

// criarTelaConfig cria a tela de configurações
func criarTelaConfig(janela fyne.Window, config *Config, fnSalvarConfig func(*Config)) fyne.CanvasObject {
	// Entradas para as configurações
	entradaURL := widget.NewEntry()
	entradaURL.Text = config.URLZabbix
	entradaURL.SetPlaceHolder("URL da API do Zabbix")

	entradaTempoLimite := widget.NewEntry()
	entradaTempoLimite.Text = formatarNumero(config.TempoLimite)
	entradaTempoLimite.SetPlaceHolder("Tempo limite (segundos)")

	// Botões
	botaoSalvar := widget.NewButton("Salvar", func() {
		// Validar URL
		if entradaURL.Text == "" {
			dialog.ShowError(fyne.NewError(1, "A URL da API é obrigatória"), janela)
			return
		}

		// Validar tempo limite
		tempoLimite, err := parseNumero(entradaTempoLimite.Text)
		if err != nil || tempoLimite <= 0 {
			dialog.ShowError(fyne.NewError(1, "Tempo limite inválido. Use um número maior que zero."), janela)
			return
		}

		// Criar nova configuração
		novaConfig := &Config{
			URLZabbix:   entradaURL.Text,
			TempoLimite: tempoLimite,
			Token:       config.Token, // Manter o token atual
		}

		// Salvar configuração
		fnSalvarConfig(novaConfig)

		// Fechar janela
		dialog.ShowInformation("Sucesso", "Configurações salvas com sucesso!", janela)
		janela.Close()
	})

	botaoCancelar := widget.NewButton("Cancelar", func() {
		janela.Close()
	})

	// Formulário
	formulario := widget.NewForm(
		widget.NewFormItem("URL da API", entradaURL),
		widget.NewFormItem("Tempo Limite (segundos)", entradaTempoLimite),
	)

	// Layout
	conteudo := container.NewVBox(
		widget.NewLabelWithStyle("Configurações", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		formulario,
		container.NewHBox(
			layout.NewSpacer(),
			botaoCancelar,
			botaoSalvar,
		),
	)

	return container.NewPadded(conteudo)
}

// formatarNumero formata um número como string
func formatarNumero(numero int) string {
	return strconv.Itoa(numero)
}

// parseNumero converte uma string para número
func parseNumero(texto string) (int, error) {
	return strconv.Atoi(texto)
}
