package ui

import (
        "fmt"

        "fyne.io/fyne/v2"
        "fyne.io/fyne/v2/dialog"

        "zabbix-manager/config"
        "zabbix-manager/zabbix"
)

// GerenciadorUI gerencia a interface do usuário da aplicação
type GerenciadorUI struct {
        app             fyne.App
        janelaPrincipal fyne.Window
        janelaLogin     fyne.Window
        clienteZabbix   *zabbix.ClienteAPI
        config          *config.Configuração
        perfilAtual     config.Perfil
}

// NovoGerenciadorUI cria uma nova instância do gerenciador de UI
func NovoGerenciadorUI(app fyne.App, configuracao *config.Configuração) *GerenciadorUI {
        // Obter o perfil ativo ou o primeiro da lista
        var perfilAtivo config.Perfil
        temPerfil := false

        if perfil, encontrado := configuracao.ObterPerfilAtivo(); encontrado {
                perfilAtivo = perfil
                temPerfil = true
        } else if len(configuracao.PerfisZabbix) > 0 {
                perfilAtivo = configuracao.PerfisZabbix[0]
                configuracao.DefinirPerfilAtivo(perfilAtivo.Nome)
                temPerfil = true
        } else {
                // Criar um perfil padrão se não existir nenhum
                perfilAtivo = config.Perfil{
                        Nome:     "Padrão",
                        URL:      "http://localhost/zabbix/api_jsonrpc.php",
                        Token:    "",
                        Usuário:  "",
                        Senha:    "",
                        Salvar:   false,
                }
                configuracao.AdicionarPerfil(perfilAtivo)
                configuracao.DefinirPerfilAtivo(perfilAtivo.Nome)
                configuracao.Salvar()
        }

        // Criar o cliente Zabbix com o perfil selecionado
        clienteZabbix := zabbix.NovoClienteAPI(zabbix.ConfigAPI{
                URL:         perfilAtivo.URL,
                Token:       perfilAtivo.Token,
                TempoLimite: configuracao.TempoLimite,
        })

        return &GerenciadorUI{
                app:           app,
                clienteZabbix: clienteZabbix,
                config:        configuracao,
                perfilAtual:   perfilAtivo,
        }
}

// Iniciar inicia a aplicação
func (g *GerenciadorUI) Iniciar() error {
        // Aplicar o tema personalizado
        g.app.Settings().SetTheme(NovoTemaZabbix())

        // Verificar se o token já é válido
        if g.clienteZabbix.Config.Token != "" {
                if g.clienteZabbix.VerificarAutenticação() {
                        // Token válido, mostrar a tela principal diretamente
                        g.MostrarTelaPrincipal()
                        return nil
                }
        }

        // Se chegou aqui, precisa de login
        g.MostrarTelaLogin()
        return nil
}

// MostrarTelaLogin exibe a tela de login
func (g *GerenciadorUI) MostrarTelaLogin() {
        g.janelaLogin = g.app.NewWindow("Zabbix Manager - Login")
        g.janelaLogin.SetMaster()
        g.janelaLogin.Resize(fyne.NewSize(500, 400))
        g.janelaLogin.CenterOnScreen()

        g.janelaLogin.SetContent(criarTelaLogin(g.janelaLogin, g.config, g.perfilAtual, g.Autenticar, g.AutenticarComToken, g.SelecionarPerfil, g.GerenciarPerfis))
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
        conteudo, atualizarDados := criarTelaPrincipal(g.janelaPrincipal, g.clienteZabbix, g.MostrarTelaConfig, g.MostrarTelaLogin)
        g.janelaPrincipal.SetContent(conteudo)
        
        // Carregar os dados iniciais
        atualizarDados()

        // Fechar a janela de login se estiver aberta
        if g.janelaLogin != nil {
                g.janelaLogin.Hide()
        }

        g.janelaPrincipal.Show()

        // Configurar evento de fechamento
        g.janelaPrincipal.SetOnClosed(func() {
                // Salvar configurações antes de fechar
                g.config.Salvar()
        })
}

// MostrarTelaConfig exibe a tela de configuração
func (g *GerenciadorUI) MostrarTelaConfig() {
        janelaConfig := g.app.NewWindow("Configurações")
        janelaConfig.Resize(fyne.NewSize(500, 350))
        janelaConfig.CenterOnScreen()

        janelaConfig.SetContent(criarTelaConfig(janelaConfig, g.config, g.perfilAtual, g.AtualizarConfig))
        janelaConfig.Show()
}

// MostrarGerenciadorPerfis exibe a tela de gerenciamento de perfis
func (g *GerenciadorUI) GerenciarPerfis() {
        janelaPerfis := g.app.NewWindow("Gerenciador de Perfis")
        janelaPerfis.Resize(fyne.NewSize(600, 400))
        janelaPerfis.CenterOnScreen()

        janelaPerfis.SetContent(criarTelaPerfis(janelaPerfis, g.config, g.AtualizarPerfil, g.ExcluirPerfil))
        janelaPerfis.Show()
}

// SelecionarPerfil seleciona um perfil pelo nome
func (g *GerenciadorUI) SelecionarPerfil(nomePerfil string) {
        // Buscar o perfil pelo nome
        perfil, encontrado := g.config.ObterPerfil(nomePerfil)
        if encontrado {
                // Atualizar perfil ativo na configuração
                g.config.DefinirPerfilAtivo(nomePerfil)
                
                // Atualizar o perfil atual no gerenciador
                g.perfilAtual = perfil
                
                // Atualizar o cliente Zabbix
                g.clienteZabbix.Config.URL = perfil.URL
                g.clienteZabbix.Config.Token = perfil.Token
                
                // Salvar as alterações
                g.config.Salvar()
                
                // Atualizar a interface se necessário
                if g.janelaLogin != nil {
                        g.janelaLogin.SetContent(criarTelaLogin(g.janelaLogin, g.config, g.perfilAtual, g.Autenticar, g.AutenticarComToken, g.SelecionarPerfil, g.GerenciarPerfis))
                }
        }
}

// Autenticar realiza a autenticação do usuário
func (g *GerenciadorUI) Autenticar(usuario, senha string, salvarCredenciais bool, nomePerfil string) {
        // Obter o perfil selecionado ou criar um novo se o nome for diferente
        var perfil config.Perfil
        encontrado := false
        
        if nomePerfil != g.perfilAtual.Nome {
                // Verificar se o perfil já existe
                perfil, encontrado = g.config.ObterPerfil(nomePerfil)
                if !encontrado {
                        // Criar novo perfil
                        perfil = config.Perfil{
                                Nome:     nomePerfil,
                                URL:      g.perfilAtual.URL,
                                Usuário:  usuario,
                                Senha:    senha,
                                Token:    "",
                                Salvar:   salvarCredenciais,
                        }
                        g.config.AdicionarPerfil(perfil)
                }
        } else {
                perfil = g.perfilAtual
                perfil.Usuário = usuario
                perfil.Senha = senha
                perfil.Salvar = salvarCredenciais
        }
        
        // Atualizar o cliente Zabbix com os dados do perfil
        g.clienteZabbix.Config.URL = perfil.URL

        // Tentar autenticar
        token, err := g.clienteZabbix.Autenticar(usuario, senha)
        if err != nil {
                dialog.ShowError(fmt.Errorf("Erro na autenticação: %v", err), g.janelaLogin)
                return
        }

        // Atualizar o token no perfil
        perfil.Token = token
        
        // Limpar senha se não for para salvar
        if !salvarCredenciais {
                perfil.Senha = ""
        }
        
        // Atualizar o perfil na configuração
        g.config.AdicionarPerfil(perfil)
        g.config.DefinirPerfilAtivo(perfil.Nome)
        g.perfilAtual = perfil
        
        // Salvar configurações
        g.config.Salvar()

        // Mostrar a tela principal
        g.MostrarTelaPrincipal()
}

// AutenticarComToken autentica usando um token diretamente
func (g *GerenciadorUI) AutenticarComToken(token string, salvarToken bool, nomePerfil string) {
        // Obter o perfil selecionado ou criar um novo se o nome for diferente
        var perfil config.Perfil
        encontrado := false
        
        if nomePerfil != g.perfilAtual.Nome {
                // Verificar se o perfil já existe
                perfil, encontrado = g.config.ObterPerfil(nomePerfil)
                if !encontrado {
                        // Criar novo perfil
                        perfil = config.Perfil{
                                Nome:     nomePerfil,
                                URL:      g.perfilAtual.URL,
                                Usuário:  "",
                                Senha:    "",
                                Token:    token,
                                Salvar:   salvarToken,
                        }
                        g.config.AdicionarPerfil(perfil)
                }
        } else {
                perfil = g.perfilAtual
                perfil.Token = token
                perfil.Salvar = salvarToken
        }
        
        // Configurar o cliente Zabbix com o token
        g.clienteZabbix.Config.URL = perfil.URL
        g.clienteZabbix.AutenticarComToken(token)
        
        // Verificar se o token é válido
        if !g.clienteZabbix.VerificarAutenticação() {
                dialog.ShowError(fmt.Errorf("Erro na autenticação: token inválido"), g.janelaLogin)
                return
        }
        
        // Não salvar o token se não for solicitado
        if !salvarToken {
                perfil.Token = ""
        }
        
        // Atualizar o perfil na configuração
        g.config.AdicionarPerfil(perfil)
        g.config.DefinirPerfilAtivo(perfil.Nome)
        g.perfilAtual = perfil
        
        // Salvar configurações
        g.config.Salvar()

        // Mostrar a tela principal
        g.MostrarTelaPrincipal()
}

// AtualizarConfig atualiza as configurações da aplicação
func (g *GerenciadorUI) AtualizarConfig(tempoLimite int, diretorioCSV string, perfilAtualizado config.Perfil) {
        // Atualizar as configurações gerais
        g.config.TempoLimite = tempoLimite
        g.config.DiretórioCSV = diretorioCSV
        
        // Atualizar o perfil atual
        g.config.AdicionarPerfil(perfilAtualizado)
        g.perfilAtual = perfilAtualizado
        
        // Atualizar o cliente Zabbix
        g.clienteZabbix.Config.URL = perfilAtualizado.URL
        g.clienteZabbix.Config.TempoLimite = tempoLimite
        
        // Salvar as alterações
        g.config.Salvar()
}

// AtualizarPerfil atualiza um perfil
func (g *GerenciadorUI) AtualizarPerfil(perfil config.Perfil) {
        g.config.AdicionarPerfil(perfil)
        
        // Se for o perfil ativo, atualizar também o cliente Zabbix
        if perfil.Nome == g.perfilAtual.Nome {
                g.perfilAtual = perfil
                g.clienteZabbix.Config.URL = perfil.URL
                g.clienteZabbix.Config.Token = perfil.Token
        }
        
        // Salvar as alterações
        g.config.Salvar()
}

// ExcluirPerfil exclui um perfil pelo nome
func (g *GerenciadorUI) ExcluirPerfil(nomePerfil string) {
        // Não permitir excluir se for o único perfil
        if len(g.config.PerfisZabbix) <= 1 {
                return
        }
        
        // Verificar se é o perfil ativo
        if nomePerfil == g.perfilAtual.Nome {
                // Buscar outro perfil para definir como ativo
                for _, p := range g.config.PerfisZabbix {
                        if p.Nome != nomePerfil {
                                g.config.DefinirPerfilAtivo(p.Nome)
                                g.perfilAtual = p
                                g.clienteZabbix.Config.URL = p.URL
                                g.clienteZabbix.Config.Token = p.Token
                                break
                        }
                }
        }
        
        // Remover o perfil
        g.config.RemoverPerfil(nomePerfil)
        
        // Salvar as alterações
        g.config.Salvar()
}
