package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"

	"../zabbix"
)

// GerenciadorUI gerencia a interface do usuário da aplicação
type GerenciadorUI struct {
	app             fyne.App
	janelaPrincipal fyne.Window
	janelaLogin     fyne.Window
	clienteZabbix   *zabbix.ClienteAPI
	config          *Config
}

// Config representa as configurações da aplicação
type Config struct {
	URLZabbix   string
	TempoLimite int
	Token       string
}

// NovoGerenciadorUI cria uma nova instância do gerenciador de UI
func NovoGerenciadorUI(app fyne.App) *GerenciadorUI {
	// Configuração padrão inicial
	config := &Config{
		URLZabbix:   "http://localhost/zabbix/api_jsonrpc.php",
		TempoLimite: 30,
	}

	// Criar o cliente Zabbix
	clienteZabbix := zabbix.NovoClienteAPI(zabbix.ConfigAPI{
		URL:         config.URLZabbix,
		TempoLimite: config.TempoLimite,
	})

	return &GerenciadorUI{
		app:           app,
		clienteZabbix: clienteZabbix,
		config:        config,
	}
}

// Iniciar inicia a aplicação
func (g *GerenciadorUI) Iniciar() error {
	// Aplicar o tema personalizado
	g.app.Settings().SetTheme(NovoTemaZabbix())

	// Iniciar a tela de login
	g.MostrarTelaLogin()

	return nil
}

// MostrarTelaLogin exibe a tela de login
func (g *GerenciadorUI) MostrarTelaLogin() {
	g.janelaLogin = g.app.NewWindow("Zabbix Manager - Login")
	g.janelaLogin.SetMaster()
	g.janelaLogin.Resize(fyne.NewSize(400, 300))
	g.janelaLogin.CenterOnScreen()

	g.janelaLogin.SetContent(criarTelaLogin(g.janelaLogin, g.config, g.Autenticar))
	g.janelaLogin.Show()
}

// MostrarTelaPrincipal exibe a tela principal
func (g *GerenciadorUI) MostrarTelaPrincipal() {
	if g.janelaPrincipal != nil {
		g.janelaPrincipal.Show()
		return
	}

	g.janelaPrincipal = g.app.NewWindow("Zabbix Manager")
	g.janelaPrincipal.SetMaster()
	g.janelaPrincipal.Resize(fyne.NewSize(900, 600))
	g.janelaPrincipal.CenterOnScreen()

	// Criar a tela principal
	conteudo, atualizarDados := criarTelaPrincipal(g.janelaPrincipal, g.clienteZabbix, g.MostrarTelaConfig)
	g.janelaPrincipal.SetContent(conteudo)
	
	// Carregar os dados iniciais
	atualizarDados()

	// Fechar a janela de login se estiver aberta
	if g.janelaLogin != nil {
		g.janelaLogin.Hide()
	}

	g.janelaPrincipal.Show()
}

// MostrarTelaConfig exibe a tela de configuração
func (g *GerenciadorUI) MostrarTelaConfig() {
	janelaConfig := g.app.NewWindow("Configurações")
	janelaConfig.Resize(fyne.NewSize(500, 300))
	janelaConfig.CenterOnScreen()

	janelaConfig.SetContent(criarTelaConfig(janelaConfig, g.config, g.AtualizarConfig))
	janelaConfig.Show()
}

// Autenticar realiza a autenticação do usuário
func (g *GerenciadorUI) Autenticar(usuario, senha string) {
	// Atualizar a URL da API com a configurada
	g.clienteZabbix.Config.URL = g.config.URLZabbix

	// Tentar autenticar
	token, err := g.clienteZabbix.Autenticar(usuario, senha)
	if err != nil {
		dialog.ShowError(fmt.Errorf("Erro na autenticação: %v", err), g.janelaLogin)
		return
	}

	// Salvar o token na configuração
	g.config.Token = token
	g.clienteZabbix.Config.Token = token

	// Mostrar a tela principal
	g.MostrarTelaPrincipal()
}

// AtualizarConfig atualiza as configurações da aplicação
func (g *GerenciadorUI) AtualizarConfig(novaConfig *Config) {
	// Atualizar as configurações
	g.config.URLZabbix = novaConfig.URLZabbix
	g.config.TempoLimite = novaConfig.TempoLimite

	// Atualizar o cliente Zabbix
	g.clienteZabbix.Config.URL = novaConfig.URLZabbix
	g.clienteZabbix.Config.TempoLimite = novaConfig.TempoLimite
}
