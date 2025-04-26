package ui

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"zabbix-manager/config"
)

// TelaConfig representa a tela de configurações do aplicativo
type TelaConfig struct {
	App           *AplicacaoZabbix
	Container     *fyne.Container
	CampoTimeout  *widget.Entry
	BotaoSalvar   *widget.Button
	BotaoVoltar   *widget.Button
	ConfigOriginal *config.Configuração
}

// MostrarTelaConfig exibe a tela de configurações
func (a *AplicacaoZabbix) MostrarTelaConfig() {
	// Fazer uma cópia da configuração atual
	configOriginal := *a.Config

	tela := &TelaConfig{
		App:            a,
		ConfigOriginal: &configOriginal,
	}

	tela.CriarInterface()
	a.Janela.SetContent(tela.Container)
}

// CriarInterface cria a interface da tela de configurações
func (t *TelaConfig) CriarInterface() {
	// Criar widgets
	t.CriarWidgets()

	// Container principal
	t.Container = container.NewVBox(
		widget.NewLabelWithStyle("Configurações", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		container.NewVBox(
			widget.NewLabel("Tempo limite de requisição (segundos):"),
			t.CampoTimeout,
		),
		widget.NewSeparator(),
		container.NewHBox(
			layout.NewSpacer(),
			t.BotaoVoltar,
			t.BotaoSalvar,
		),
	)
}

// CriarWidgets cria os widgets da tela
func (t *TelaConfig) CriarWidgets() {
	// Campo de tempo limite
	t.CampoTimeout = widget.NewEntry()
	t.CampoTimeout.Text = strconv.Itoa(t.App.Config.TempoLimite)
	// Validar que é um número inteiro positivo
	t.CampoTimeout.Validator = validation.NewRegexp(`^[1-9]\d*$`, "Deve ser um número inteiro positivo")

	// Botões
	t.BotaoSalvar = widget.NewButton("Salvar", t.Salvar)
	t.BotaoVoltar = widget.NewButton("Voltar", t.Voltar)
}

// Salvar salva as configurações
func (t *TelaConfig) Salvar() {
	// Validar campos
	if err := t.CampoTimeout.Validate(); err != nil {
		t.App.MostrarErro("Erro de Validação", "Tempo limite inválido. Deve ser um número inteiro positivo.")
		return
	}

	// Converter tempo limite para inteiro
	timeout, err := strconv.Atoi(t.CampoTimeout.Text)
	if err != nil {
		t.App.MostrarErro("Erro", "Tempo limite inválido")
		return
	}

	// Atualizar configurações
	t.App.Config.TempoLimite = timeout

	// Salvar configurações
	err = t.App.SalvarConfiguracao()
	if err != nil {
		t.App.MostrarErro("Erro", fmt.Sprintf("Erro ao salvar configurações: %v", err))
		return
	}

	// Atualizar cliente API se existir um perfil ativo
	if t.App.PerfilAtual != nil {
		t.App.ConfigurarClienteAPI(t.App.PerfilAtual)
	}

	// Mostrar mensagem de sucesso
	t.App.MostrarInfo("Sucesso", "Configurações salvas com sucesso")

	// Voltar para a tela anterior
	t.Voltar()
}

// Voltar volta para a tela anterior
func (t *TelaConfig) Voltar() {
	// Restaurar a configuração original se não foi salva
	if t.App.Config.TempoLimite != t.ConfigOriginal.TempoLimite {
		pergunta := "Deseja descartar as alterações?"
		t.App.MostrarConfirmacao("Alterações não salvas", pergunta, func(confirmar bool) {
			if confirmar {
				// Restaurar configuração original
				*t.App.Config = *t.ConfigOriginal

				// Voltar para a tela principal
				t.App.MostrarTelaPrincipal()
			}
		})
	} else {
		// Se não houve alterações, voltar direto
		t.App.MostrarTelaPrincipal()
	}
}