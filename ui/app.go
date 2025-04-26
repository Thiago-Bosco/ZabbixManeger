package ui

import (
        "fmt"
        "log"
        "os"
        "path/filepath"

        "fyne.io/fyne/v2"
        "fyne.io/fyne/v2/app"
        "fyne.io/fyne/v2/dialog"
        "fyne.io/fyne/v2/theme"
        "zabbix-manager/config"
        "zabbix-manager/zabbix"
)

// AplicacaoZabbix representa a aplicação principal do Zabbix Manager
type AplicacaoZabbix struct {
        App          fyne.App
        Janela       fyne.Window
        Config       *config.Configuração
        Cliente      *zabbix.ClienteAPI
        ArquivoConfig string
        PerfilAtual  *config.ConfiguracaoPerfil
}

// NovaAplicacao cria uma nova instância da aplicação Zabbix Manager
func NovaAplicacao() *AplicacaoZabbix {
        // Criar aplicação Fyne
        aplicacaoFyne := app.New()
        aplicacaoFyne.Settings().SetTheme(theme.DarkTheme())

        // Configurar caminho do arquivo de configuração
        arquivoConfig := config.ObterCaminhoConfiguracao()

        // Carregar a configuração
        configuracao, err := config.Carregar(arquivoConfig)
        if err != nil {
                log.Printf("Erro ao carregar configuração: %v", err)
                configuracao = config.NovaConfiguração()
        }

        // Verificar se há algum perfil ativo
        perfilAtivo, _ := configuracao.PerfilAtivo()

        // Criar janela principal
        janela := aplicacaoFyne.NewWindow("Zabbix Manager")
        janela.Resize(fyne.NewSize(1000, 600))
        janela.CenterOnScreen()

        // Criar a aplicação Zabbix
        app := &AplicacaoZabbix{
                App:          aplicacaoFyne,
                Janela:       janela,
                Config:       configuracao,
                ArquivoConfig: arquivoConfig,
                PerfilAtual:  perfilAtivo,
        }

        // Configurar cliente API se houver perfil ativo
        if perfilAtivo != nil {
                app.ConfigurarClienteAPI(perfilAtivo)
        }

        return app
}

// ConfigurarClienteAPI configura o cliente da API do Zabbix com os dados do perfil
func (a *AplicacaoZabbix) ConfigurarClienteAPI(perfil *config.ConfiguracaoPerfil) {
        configAPI := zabbix.ConfigAPI{
                URL:         perfil.URL,
                Token:       perfil.Token,
                TempoLimite: a.Config.TempoLimite,
        }
        a.Cliente = zabbix.NovoClienteAPI(configAPI)
        a.PerfilAtual = perfil
}

// SalvarConfiguracao salva a configuração atual no arquivo
func (a *AplicacaoZabbix) SalvarConfiguracao() error {
        return a.Config.Salvar(a.ArquivoConfig)
}

// AdicionarPerfil adiciona um novo perfil à configuração
func (a *AplicacaoZabbix) AdicionarPerfil(perfil config.ConfiguracaoPerfil) error {
        // Verificar se já existe um perfil com o mesmo nome
        for _, p := range a.Config.Perfis {
                if p.Nome == perfil.Nome {
                        return fmt.Errorf("já existe um perfil com o nome '%s'", perfil.Nome)
                }
        }

        // Adicionar o perfil
        a.Config.AdicionarPerfil(perfil)

        // Salvar a configuração
        return a.SalvarConfiguracao()
}

// AtualizarPerfil atualiza um perfil existente
func (a *AplicacaoZabbix) AtualizarPerfil(indice int, perfil config.ConfiguracaoPerfil) error {
        // Verificar se existe outro perfil com o mesmo nome
        for i, p := range a.Config.Perfis {
                if p.Nome == perfil.Nome && i != indice {
                        return fmt.Errorf("já existe um perfil com o nome '%s'", perfil.Nome)
                }
        }

        // Atualizar o perfil
        err := a.Config.AtualizarPerfil(indice, perfil)
        if err != nil {
                return err
        }

        // Salvar a configuração
        return a.SalvarConfiguracao()
}

// RemoverPerfil remove um perfil existente
func (a *AplicacaoZabbix) RemoverPerfil(indice int) error {
        // Remover o perfil
        err := a.Config.RemoverPerfil(indice)
        if err != nil {
                return err
        }

        // Salvar a configuração
        return a.SalvarConfiguracao()
}

// SelecionarPerfil seleciona um perfil pelo índice
func (a *AplicacaoZabbix) SelecionarPerfil(indice int) error {
        // Selecionar o perfil
        err := a.Config.SelecionarPerfil(indice)
        if err != nil {
                return err
        }

        // Atualizar o cliente API
        perfil, err := a.Config.PerfilAtivo()
        if err != nil {
                return err
        }

        a.ConfigurarClienteAPI(perfil)

        // Salvar a configuração
        return a.SalvarConfiguracao()
}

// TestarConexao testa a conexão com o servidor Zabbix
func (a *AplicacaoZabbix) TestarConexao(url, token string) error {
        configTemporaria := zabbix.ConfigAPI{
                URL:         url,
                Token:       token,
                TempoLimite: a.Config.TempoLimite,
        }
        clienteTemporario := zabbix.NovoClienteAPI(configTemporaria)
        return clienteTemporario.TestarConexao()
}

// ExportarRelatorio exporta um relatório CSV com os dados dos hosts
func (a *AplicacaoZabbix) ExportarRelatorio() error {
        // Verificar se há cliente configurado
        if a.Cliente == nil {
                return fmt.Errorf("nenhum servidor Zabbix configurado")
        }

        // Obter os hosts
        hosts, err := a.Cliente.ObterHosts()
        if err != nil {
                return fmt.Errorf("erro ao obter hosts: %w", err)
        }

        // Criar diretório para relatórios se não existir
        diretorioHome, err := os.UserHomeDir()
        if err != nil {
                diretorioHome, _ = os.Getwd()
        }

        diretorioRelatorios := filepath.Join(diretorioHome, "Relatórios Zabbix")
        err = os.MkdirAll(diretorioRelatorios, 0755)
        if err != nil {
                return fmt.Errorf("erro ao criar diretório de relatórios: %w", err)
        }

        // Gerar nome do arquivo
        nomeServidor := "zabbix"
        if a.PerfilAtual != nil {
                nomeServidor = a.PerfilAtual.Nome
        }
        caminhoArquivo := filepath.Join(diretorioRelatorios, fmt.Sprintf("relatorio_%s.csv", nomeServidor))

        // Gerar o relatório
        return zabbix.GerarRelatorioCSV(hosts, caminhoArquivo)
}

// MostrarErro exibe uma caixa de diálogo de erro
func (a *AplicacaoZabbix) MostrarErro(titulo, mensagem string) {
        dialog.ShowError(fmt.Errorf(mensagem), a.Janela)
}

// MostrarInfo exibe uma caixa de diálogo de informação
func (a *AplicacaoZabbix) MostrarInfo(titulo, mensagem string) {
        dialog.ShowInformation(titulo, mensagem, a.Janela)
}

// MostrarConfirmacao exibe uma caixa de diálogo de confirmação
func (a *AplicacaoZabbix) MostrarConfirmacao(titulo, mensagem string, callback func(bool)) {
        dialog.ShowConfirm(titulo, mensagem, callback, a.Janela)
}

// ObterDiretorioHome retorna o diretório home do usuário
func (a *AplicacaoZabbix) ObterDiretorioHome() (string, error) {
        diretorioHome, err := os.UserHomeDir()
        if err != nil {
                diretorioHome, _ = os.Getwd()
        }
        return diretorioHome, nil
}

// Iniciar inicia a aplicação
func (a *AplicacaoZabbix) Iniciar() {
        // Verificar se há algum perfil configurado
        if len(a.Config.Perfis) == 0 {
                // Se não houver perfis, mostrar a tela de configuração
                a.MostrarTelaLogin()
        } else {
                // Se houver perfis, verificar se há um perfil ativo
                if a.Config.PerfilAtual < 0 || a.Config.PerfilAtual >= len(a.Config.Perfis) {
                        // Se não houver perfil ativo, mostrar a tela de login
                        a.MostrarTelaLogin()
                } else {
                        // Se houver perfil ativo, verificar a conexão com o servidor
                        perfil, _ := a.Config.PerfilAtivo()
                        a.ConfigurarClienteAPI(perfil)
                        
                        err := a.Cliente.TestarConexao()
                        if err != nil {
                                // Se a conexão falhar, mostrar a tela de login
                                a.MostrarTelaLogin()
                        } else {
                                // Se a conexão for bem-sucedida, mostrar a tela principal
                                a.MostrarTelaPrincipal()
                        }
                }
        }

        // Executar a aplicação
        a.Janela.ShowAndRun()
}