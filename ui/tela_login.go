package ui

import (
        "fmt"
        "net/url"
        "strings"

        "fyne.io/fyne/v2"
        "fyne.io/fyne/v2/canvas"
        "fyne.io/fyne/v2/container"
        "fyne.io/fyne/v2/layout"
        "fyne.io/fyne/v2/theme"
        "fyne.io/fyne/v2/widget"
        "zabbix-manager/config"
)

// TelaLogin representa a tela de login do aplicativo
type TelaLogin struct {
        App                     *AplicacaoZabbix
        Container               *fyne.Container
        TabContainer           *container.AppTabs
        TabAtual                int
        ComboServidores         *widget.Select
        CampoURL                *widget.Entry
        CampoToken              *widget.Entry
        CampoPerfil             *widget.Entry
        BotaoAdicionar          *widget.Button
        BotaoEditar             *widget.Button
        BotaoRemover            *widget.Button
        BotaoEntrar             *widget.Button
        ListaPerfis             []config.ConfiguracaoPerfil
}

// MostrarTelaLogin exibe a tela de login
func (a *AplicacaoZabbix) MostrarTelaLogin() {
        tela := &TelaLogin{
                App:       a,
                TabAtual:  0,
                ListaPerfis: a.Config.Perfis,
        }

        tela.CriarInterface()
        a.Janela.SetContent(tela.Container)
}

// CriarInterface cria a interface da tela de login
func (t *TelaLogin) CriarInterface() {
        // Criação dos widgets
        t.CriarWidgets()

        // Layout principal da tela
        titulo := widget.NewLabelWithStyle("Zabbix Manager", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

        imagem := canvas.NewImageFromFile("assets/logo.png")
        imagem.FillMode = canvas.ImageFillContain
        imagem.SetMinSize(fyne.NewSize(200, 100))

        // Tabs de login
        t.TabContainer = container.NewAppTabs(
                container.NewTabItem("Login com Token", t.CriarFormularioLogin()),
                container.NewTabItem("Gerenciar Perfis", t.CriarFormularioPerfis()),
        )

        t.TabContainer.SetTabLocation(container.TabLocationTop)
        t.TabContainer.SelectIndex(t.TabAtual)

        // Container principal
        t.Container = container.NewVBox(
                titulo,
                container.NewHBox(layout.NewSpacer(), imagem, layout.NewSpacer()),
                widget.NewSeparator(),
                t.TabContainer,
        )
}

// CriarWidgets cria os widgets da tela
func (t *TelaLogin) CriarWidgets() {
        // Campo URL
        t.CampoURL = widget.NewEntry()
        t.CampoURL.SetPlaceHolder("http://seu-servidor-zabbix/api_jsonrpc.php")

        // Campo Token
        t.CampoToken = widget.NewEntry()
        t.CampoToken.Password = true
        t.CampoToken.SetPlaceHolder("Token de API do Zabbix")

        // Campo Nome do Perfil
        t.CampoPerfil = widget.NewEntry()
        t.CampoPerfil.SetPlaceHolder("Nome para identificar este servidor")

        // Combo de servidores
        t.AtualizarComboServidores()

        // Botões
        t.BotaoEntrar = widget.NewButton("Entrar", t.acaoEntrar)
        t.BotaoAdicionar = widget.NewButton("Adicionar Perfil", t.acaoAdicionarPerfil)
        t.BotaoEditar = widget.NewButton("Editar Perfil", t.acaoEditarPerfil)
        t.BotaoRemover = widget.NewButton("Remover", t.acaoRemoverPerfil)

        // Desabilitar botões de edição/remoção se não houver perfis
        t.atualizarEstadoBotoes()
}

// AtualizarComboServidores atualiza o combo de servidores com os perfis disponíveis
func (t *TelaLogin) AtualizarComboServidores() {
        nomesPerfis := []string{}
        for _, perfil := range t.ListaPerfis {
                nomesPerfis = append(nomesPerfis, perfil.Nome)
        }

        var valorSelecionado string
        if t.ComboServidores != nil && t.ComboServidores.Selected != "" {
                valorSelecionado = t.ComboServidores.Selected
        } else if t.App.Config.PerfilAtual >= 0 && t.App.Config.PerfilAtual < len(t.ListaPerfis) {
                valorSelecionado = t.ListaPerfis[t.App.Config.PerfilAtual].Nome
        } else if len(nomesPerfis) > 0 {
                valorSelecionado = nomesPerfis[0]
        }

        t.ComboServidores = widget.NewSelect(nomesPerfis, t.acaoSelecionarPerfil)
        if valorSelecionado != "" {
                t.ComboServidores.SetSelected(valorSelecionado)
        }
}

// CriarFormularioLogin cria o formulário de login
func (t *TelaLogin) CriarFormularioLogin() *fyne.Container {
        var formWidget fyne.CanvasObject

        if len(t.ListaPerfis) > 0 {
                // Se houver perfis, mostrar o combo de seleção
                formWidget = container.NewVBox(
                        widget.NewLabelWithStyle("Selecione o servidor Zabbix:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
                        t.ComboServidores,
                        container.NewHBox(layout.NewSpacer(), t.BotaoEntrar, layout.NewSpacer()),
                )
        } else {
                // Se não houver perfis, mostrar mensagem informativa
                formWidget = container.NewVBox(
                        widget.NewLabelWithStyle("Não há servidores Zabbix configurados", fyne.TextAlignCenter, fyne.TextStyle{}),
                        widget.NewLabelWithStyle("Adicione um perfil na aba 'Gerenciar Perfis'", fyne.TextAlignCenter, fyne.TextStyle{}),
                )
        }

        return container.NewVBox(
                formWidget,
        )
}

// CriarFormularioPerfis cria o formulário de gerenciamento de perfis
func (t *TelaLogin) CriarFormularioPerfis() *fyne.Container {
        return container.NewVBox(
                widget.NewLabelWithStyle("URL do servidor Zabbix:", fyne.TextAlignLeading, fyne.TextStyle{}),
                t.CampoURL,
                widget.NewLabelWithStyle("Token de API:", fyne.TextAlignLeading, fyne.TextStyle{}),
                t.CampoToken,
                widget.NewLabelWithStyle("Nome do Perfil:", fyne.TextAlignLeading, fyne.TextStyle{}),
                t.CampoPerfil,
                container.NewHBox(
                        t.BotaoAdicionar,
                        layout.NewSpacer(),
                        t.BotaoEditar,
                        layout.NewSpacer(),
                        t.BotaoRemover,
                ),
        )
}

// atualizarEstadoBotoes atualiza o estado dos botões baseado nos perfis disponíveis
func (t *TelaLogin) atualizarEstadoBotoes() {
        temPerfis := len(t.ListaPerfis) > 0
        t.BotaoEditar.Disabled = !temPerfis
        t.BotaoRemover.Disabled = !temPerfis
}

// acaoEntrar processa a ação de entrar no sistema
func (t *TelaLogin) acaoEntrar() {
        if len(t.ListaPerfis) == 0 {
                t.App.MostrarErro("Erro", "Não há perfis configurados")
                t.TabContainer.SelectIndex(1) // Mudar para a aba de gerenciamento de perfis
                return
        }

        // Obter o índice do perfil selecionado
        perfilSelecionado := t.ComboServidores.Selected
        var indice int = -1
        for i, p := range t.ListaPerfis {
                if p.Nome == perfilSelecionado {
                        indice = i
                        break
                }
        }

        if indice < 0 {
                t.App.MostrarErro("Erro", "Perfil não encontrado")
                return
        }

        // Selecionar o perfil
        err := t.App.SelecionarPerfil(indice)
        if err != nil {
                t.App.MostrarErro("Erro", fmt.Sprintf("Erro ao selecionar perfil: %v", err))
                return
        }

        // Testar a conexão
        perfil := t.ListaPerfis[indice]
        err = t.App.TestarConexao(perfil.URL, perfil.Token)
        if err != nil {
                t.App.MostrarErro("Erro de Conexão", fmt.Sprintf("Não foi possível conectar ao servidor Zabbix: %v", err))
                return
        }

        // Se chegou aqui, a conexão foi bem-sucedida
        t.App.MostrarTelaPrincipal()
}

// acaoSelecionarPerfil processa a ação de selecionar um perfil no combo
func (t *TelaLogin) acaoSelecionarPerfil(nomePerfil string) {
        // Preencher os campos com os dados do perfil selecionado
        for _, p := range t.ListaPerfis {
                if p.Nome == nomePerfil {
                        t.CampoURL.Text = p.URL
                        t.CampoToken.Text = p.Token
                        t.CampoPerfil.Text = p.Nome
                        t.CampoURL.Refresh()
                        t.CampoToken.Refresh()
                        t.CampoPerfil.Refresh()
                        return
                }
        }
}

// acaoAdicionarPerfil processa a ação de adicionar um novo perfil
func (t *TelaLogin) acaoAdicionarPerfil() {
        // Validar campos
        url := strings.TrimSpace(t.CampoURL.Text)
        token := strings.TrimSpace(t.CampoToken.Text)
        nome := strings.TrimSpace(t.CampoPerfil.Text)

        if url == "" || token == "" || nome == "" {
                t.App.MostrarErro("Erro", "Todos os campos são obrigatórios")
                return
        }

        // Validar URL
        _, err := url.Parse(url)
        if err != nil {
                t.App.MostrarErro("URL Inválida", "A URL informada não é válida")
                return
        }

        // Testar a conexão
        err = t.App.TestarConexao(url, token)
        if err != nil {
                t.App.MostrarErro("Erro de Conexão", fmt.Sprintf("Não foi possível conectar ao servidor Zabbix: %v", err))
                return
        }

        // Criar e adicionar o perfil
        perfil := config.ConfiguracaoPerfil{
                Nome:  nome,
                URL:   url,
                Token: token,
        }

        err = t.App.AdicionarPerfil(perfil)
        if err != nil {
                t.App.MostrarErro("Erro", fmt.Sprintf("Erro ao adicionar perfil: %v", err))
                return
        }

        // Atualizar a lista de perfis
        t.ListaPerfis = t.App.Config.Perfis
        t.AtualizarComboServidores()
        t.atualizarEstadoBotoes()

        // Limpar os campos
        t.CampoURL.Text = ""
        t.CampoToken.Text = ""
        t.CampoPerfil.Text = ""
        t.CampoURL.Refresh()
        t.CampoToken.Refresh()
        t.CampoPerfil.Refresh()

        // Mostrar mensagem de sucesso
        t.App.MostrarInfo("Sucesso", "Perfil adicionado com sucesso")

        // Mudar para a aba de login
        t.TabContainer.SelectIndex(0)
}

// acaoEditarPerfil processa a ação de editar um perfil existente
func (t *TelaLogin) acaoEditarPerfil() {
        if len(t.ListaPerfis) == 0 {
                t.App.MostrarErro("Erro", "Não há perfis para editar")
                return
        }

        // Validar campos
        url := strings.TrimSpace(t.CampoURL.Text)
        token := strings.TrimSpace(t.CampoToken.Text)
        nome := strings.TrimSpace(t.CampoPerfil.Text)

        if url == "" || token == "" || nome == "" {
                t.App.MostrarErro("Erro", "Todos os campos são obrigatórios")
                return
        }

        // Validar URL
        _, err := url.Parse(url)
        if err != nil {
                t.App.MostrarErro("URL Inválida", "A URL informada não é válida")
                return
        }

        // Testar a conexão
        err = t.App.TestarConexao(url, token)
        if err != nil {
                t.App.MostrarErro("Erro de Conexão", fmt.Sprintf("Não foi possível conectar ao servidor Zabbix: %v", err))
                return
        }

        // Obter o índice do perfil selecionado
        perfilSelecionado := t.ComboServidores.Selected
        var indice int = -1
        for i, p := range t.ListaPerfis {
                if p.Nome == perfilSelecionado {
                        indice = i
                        break
                }
        }

        if indice < 0 {
                t.App.MostrarErro("Erro", "Selecione um perfil para editar")
                return
        }

        // Criar e atualizar o perfil
        perfil := config.ConfiguracaoPerfil{
                Nome:  nome,
                URL:   url,
                Token: token,
        }

        err = t.App.AtualizarPerfil(indice, perfil)
        if err != nil {
                t.App.MostrarErro("Erro", fmt.Sprintf("Erro ao atualizar perfil: %v", err))
                return
        }

        // Atualizar a lista de perfis
        t.ListaPerfis = t.App.Config.Perfis
        t.AtualizarComboServidores()

        // Mostrar mensagem de sucesso
        t.App.MostrarInfo("Sucesso", "Perfil atualizado com sucesso")
}

// acaoRemoverPerfil processa a ação de remover um perfil existente
func (t *TelaLogin) acaoRemoverPerfil() {
        if len(t.ListaPerfis) == 0 {
                t.App.MostrarErro("Erro", "Não há perfis para remover")
                return
        }

        // Obter o índice do perfil selecionado
        perfilSelecionado := t.ComboServidores.Selected
        var indice int = -1
        for i, p := range t.ListaPerfis {
                if p.Nome == perfilSelecionado {
                        indice = i
                        break
                }
        }

        if indice < 0 {
                t.App.MostrarErro("Erro", "Selecione um perfil para remover")
                return
        }

        // Pedir confirmação
        t.App.MostrarConfirmacao(
                "Remover Perfil",
                fmt.Sprintf("Deseja realmente remover o perfil '%s'?", perfilSelecionado),
                func(confirmar bool) {
                        if !confirmar {
                                return
                        }

                        // Remover o perfil
                        err := t.App.RemoverPerfil(indice)
                        if err != nil {
                                t.App.MostrarErro("Erro", fmt.Sprintf("Erro ao remover perfil: %v", err))
                                return
                        }

                        // Atualizar a lista de perfis
                        t.ListaPerfis = t.App.Config.Perfis
                        t.AtualizarComboServidores()
                        t.atualizarEstadoBotoes()

                        // Limpar os campos
                        t.CampoURL.Text = ""
                        t.CampoToken.Text = ""
                        t.CampoPerfil.Text = ""
                        t.CampoURL.Refresh()
                        t.CampoToken.Refresh()
                        t.CampoPerfil.Refresh()

                        // Mostrar mensagem de sucesso
                        t.App.MostrarInfo("Sucesso", "Perfil removido com sucesso")
                },
        )
}