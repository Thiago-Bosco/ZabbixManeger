package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// criarTelaLogin cria a tela de login
func criarTelaLogin(janela fyne.Window, config *Config, fnAutenticar func(string, string)) fyne.CanvasObject {
	// Campo de usuário
	entradaUsuario := widget.NewEntry()
	entradaUsuario.SetPlaceHolder("Nome de usuário")

	// Campo de senha
	entradaSenha := widget.NewPasswordEntry()
	entradaSenha.SetPlaceHolder("Senha")

	// Campo de URL da API
	entradaURL := widget.NewEntry()
	entradaURL.Text = config.URLZabbix
	entradaURL.SetPlaceHolder("URL da API do Zabbix")

	// Botão de login
	botaoLogin := widget.NewButton("Entrar", func() {
		// Validar campos
		if entradaUsuario.Text == "" || entradaSenha.Text == "" {
			dialog.ShowError(fyne.NewError(1, "Por favor, informe o usuário e senha"), janela)
			return
		}
		
		// Atualizar a URL na configuração
		config.URLZabbix = entradaURL.Text
		
		// Chamar a função de autenticação
		fnAutenticar(entradaUsuario.Text, entradaSenha.Text)
	})

	// Botão de configurações avançadas
	botaoConfig := widget.NewButton("Configurações Avançadas", func() {
		// Criar uma janela de diálogo para configurações
		entradaTempoLimite := widget.NewEntry()
		entradaTempoLimite.SetText(formatarNumero(config.TempoLimite))
		entradaTempoLimite.SetPlaceHolder("Tempo limite (segundos)")
		
		formConfig := widget.NewForm(
			widget.NewFormItem("Tempo Limite (segundos)", entradaTempoLimite),
		)
		
		dialogConfig := dialog.NewCustomConfirm(
			"Configurações Avançadas",
			"Salvar",
			"Cancelar",
			formConfig,
			func(confirmado bool) {
				if confirmado {
					// Tentar converter para inteiro
					tempoLimite, err := parseNumero(entradaTempoLimite.Text)
					if err != nil || tempoLimite <= 0 {
						dialog.ShowError(fyne.NewError(1, "Tempo limite inválido. Use um número maior que zero."), janela)
						return
					}
					
					// Atualizar a configuração
					config.TempoLimite = tempoLimite
				}
			},
			janela,
		)
		
		dialogConfig.Show()
	})

	// Logo
	logoZabbix := widget.NewIcon(CarregarLogoZabbix())

	// Organizar a interface
	titulo := widget.NewLabelWithStyle("Zabbix Manager", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	titulo.TextSize = 20

	conteudo := container.NewVBox(
		container.NewCenter(logoZabbix),
		titulo,
		widget.NewSeparator(),
		container.NewPadded(
			container.NewVBox(
				widget.NewLabel("Informe seus dados de acesso"),
				container.NewGridWithColumns(2,
					widget.NewLabel("Usuário:"),
					entradaUsuario,
				),
				container.NewGridWithColumns(2,
					widget.NewLabel("Senha:"),
					entradaSenha,
				),
				container.NewGridWithColumns(2,
					widget.NewLabel("URL API:"),
					entradaURL,
				),
				container.NewHBox(
					layout.NewSpacer(),
					botaoLogin,
				),
				container.NewHBox(
					layout.NewSpacer(),
					botaoConfig,
				),
			),
		),
	)

	// Permitir enviar o formulário pressionando Enter
	janela.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
		if key.Name == fyne.KeyReturn {
			botaoLogin.OnTapped()
		}
	})

	return conteudo
}
