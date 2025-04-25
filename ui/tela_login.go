package ui

import (
        "strings"

        "fyne.io/fyne/v2"
        "fyne.io/fyne/v2/container"
        "fyne.io/fyne/v2/dialog"
        "fyne.io/fyne/v2/layout"
        "fyne.io/fyne/v2/theme"
        "fyne.io/fyne/v2/widget"

        "zabbix-manager/config"
)

// criarTelaLogin cria a tela de login
func criarTelaLogin(janela fyne.Window, config *config.Configuração, perfilAtual config.Perfil,
        fnAutenticar func(string, string, bool, string),
        fnAutenticarToken func(string, bool, string),
        fnSelecionarPerfil func(string),
        fnGerenciarPerfis func()) fyne.CanvasObject {

        // Guias para os métodos de autenticação
        tabLogin := container.NewTabItem("Login com Usuário/Senha", nil)
        tabToken := container.NewTabItem("Login com Token", nil)

        // Campo para selecionar o perfil
        selecionadorPerfil := widget.NewSelect(obterNomesPerfis(config), func(nomePerfil string) {
                if nomePerfil != perfilAtual.Nome {
                        fnSelecionarPerfil(nomePerfil)
                }
        })
        selecionadorPerfil.Selected = perfilAtual.Nome

        // Botão para gerenciar perfis
        botaoGerenciarPerfis := widget.NewButtonWithIcon("Gerenciar Perfis", theme.SettingsIcon(), fnGerenciarPerfis)
        
        // Nome do novo perfil
        entradaNomePerfil := widget.NewEntry()
        entradaNomePerfil.SetPlaceHolder("Nome do perfil (opcional)")

        // Campo de URL da API
        entradaURL := widget.NewEntry()
        entradaURL.Text = perfilAtual.URL
        entradaURL.SetPlaceHolder("URL da API do Zabbix")

        // ===== Tab de Login com Usuário/Senha =====
        // Campo de usuário
        entradaUsuario := widget.NewEntry()
        entradaUsuario.SetPlaceHolder("Nome de usuário")
        
        // Preencher com o usuário salvo, se existir
        if perfilAtual.Salvar && perfilAtual.Usuário != "" {
                entradaUsuario.Text = perfilAtual.Usuário
        }

        // Campo de senha
        entradaSenha := widget.NewPasswordEntry()
        entradaSenha.SetPlaceHolder("Senha")
        
        // Preencher com a senha salva, se existir
        if perfilAtual.Salvar && perfilAtual.Senha != "" {
                entradaSenha.Text = perfilAtual.Senha
        }

        // Checkbox para salvar credenciais
        checkboxSalvarCredenciais := widget.NewCheck("Salvar credenciais", nil)
        checkboxSalvarCredenciais.Checked = perfilAtual.Salvar

        // Botão de login com usuário/senha
        botaoLoginUsuario := widget.NewButton("Entrar", func() {
                // Validar campos
                if entradaUsuario.Text == "" || entradaSenha.Text == "" {
                        dialog.ShowError(fyne.NewError(1, "Por favor, informe o usuário e senha"), janela)
                        return
                }
                
                if entradaURL.Text == "" {
                        dialog.ShowError(fyne.NewError(1, "Por favor, informe a URL da API"), janela)
                        return
                }
                
                // Atualizar a URL na configuração
                perfil := perfilAtual
                perfil.URL = entradaURL.Text
                
                // Determinar o nome do perfil
                nomePerfil := selecionadorPerfil.Selected
                if entradaNomePerfil.Text != "" {
                        nomePerfil = entradaNomePerfil.Text
                }
                
                // Chamar a função de autenticação
                fnAutenticar(
                        entradaUsuario.Text, 
                        entradaSenha.Text, 
                        checkboxSalvarCredenciais.Checked,
                        nomePerfil,
                )
        })

        // Layout da tab de login com usuário/senha
        tabLogin.Content = container.NewVBox(
                container.NewGridWithColumns(2,
                        widget.NewLabel("Usuário:"),
                        entradaUsuario,
                ),
                container.NewGridWithColumns(2,
                        widget.NewLabel("Senha:"),
                        entradaSenha,
                ),
                checkboxSalvarCredenciais,
                container.NewHBox(
                        layout.NewSpacer(),
                        botaoLoginUsuario,
                ),
        )

        // ===== Tab de Login com Token =====
        // Campo de token
        entradaToken := widget.NewPasswordEntry()
        entradaToken.SetPlaceHolder("Token de autenticação da API")

        // Checkbox para salvar token
        checkboxSalvarToken := widget.NewCheck("Salvar token", nil)
        
        // Botão de login com token
        botaoLoginToken := widget.NewButton("Entrar com Token", func() {
                // Validar campos
                if entradaToken.Text == "" {
                        dialog.ShowError(fyne.NewError(1, "Por favor, informe o token"), janela)
                        return
                }
                
                if entradaURL.Text == "" {
                        dialog.ShowError(fyne.NewError(1, "Por favor, informe a URL da API"), janela)
                        return
                }
                
                // Atualizar a URL na configuração
                perfil := perfilAtual
                perfil.URL = entradaURL.Text
                
                // Determinar o nome do perfil
                nomePerfil := selecionadorPerfil.Selected
                if entradaNomePerfil.Text != "" {
                        nomePerfil = entradaNomePerfil.Text
                }
                
                // Chamar a função de autenticação com token
                fnAutenticarToken(
                        entradaToken.Text,
                        checkboxSalvarToken.Checked,
                        nomePerfil,
                )
        })

        // Layout da tab de login com token
        tabToken.Content = container.NewVBox(
                container.NewGridWithColumns(2,
                        widget.NewLabel("Token:"),
                        entradaToken,
                ),
                checkboxSalvarToken,
                container.NewHBox(
                        layout.NewSpacer(),
                        botaoLoginToken,
                ),
        )

        // Configurar as tabs
        tabs := container.NewAppTabs(
                tabLogin,
                tabToken,
        )

        // Logo
        logoZabbix := widget.NewIcon(CarregarLogoZabbix())

        // Seção do perfil (URL e seleção de perfil)
        secaoPerfil := container.NewVBox(
                container.NewGridWithColumns(2,
                        widget.NewLabel("Perfil:"),
                        container.NewBorder(nil, nil, nil, botaoGerenciarPerfis,
                                selecionadorPerfil,
                        ),
                ),
                container.NewGridWithColumns(2,
                        widget.NewLabel("Nome do perfil:"),
                        entradaNomePerfil,
                ),
                container.NewGridWithColumns(2,
                        widget.NewLabel("URL API:"),
                        entradaURL,
                ),
        )

        // Organizar a interface
        titulo := widget.NewLabelWithStyle("Zabbix Manager", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
        titulo.TextSize = 20

        conteudo := container.NewVBox(
                container.NewCenter(logoZabbix),
                titulo,
                widget.NewSeparator(),
                container.NewPadded(
                        container.NewVBox(
                                widget.NewLabel("Selecione o perfil ou crie um novo"),
                                secaoPerfil,
                                widget.NewSeparator(),
                                widget.NewLabel("Informe seus dados de acesso"),
                                tabs,
                        ),
                ),
        )

        // Permitir enviar o formulário pressionando Enter
        janela.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
                if key.Name == fyne.KeyReturn {
                        if tabs.Selected() == 0 {
                                botaoLoginUsuario.OnTapped()
                        } else {
                                botaoLoginToken.OnTapped()
                        }
                }
        })

        return conteudo
}

// obterNomesPerfis retorna a lista de nomes dos perfis
func obterNomesPerfis(config *config.Configuração) []string {
        nomes := make([]string, len(config.PerfisZabbix))
        for i, p := range config.PerfisZabbix {
                nomes[i] = p.Nome
        }
        return nomes
}

// criarTelaPerfis cria a tela de gerenciamento de perfis
func criarTelaPerfis(janela fyne.Window, config *config.Configuração, 
        fnAtualizarPerfil func(config.Perfil), 
        fnExcluirPerfil func(string)) fyne.CanvasObject {
        
        var perfilSelecionado config.Perfil
        var indiceSelecionado int = -1
        
        // Tabela de perfis
        tabelaPerfis := widget.NewTable(
                func() (int, int) {
                        return len(config.PerfisZabbix) + 1, 4 // Cabeçalho + linhas, 4 colunas
                },
                func() fyne.CanvasObject {
                        return widget.NewLabel("Carregando...")
                },
                func(id widget.TableCellID, obj fyne.CanvasObject) {
                        label := obj.(*widget.Label)
                        label.Alignment = fyne.TextAlignLeading
                        
                        // Cabeçalho
                        if id.Row == 0 {
                                label.TextStyle = fyne.TextStyle{Bold: true}
                                switch id.Col {
                                case 0:
                                        label.SetText("Nome")
                                case 1:
                                        label.SetText("URL")
                                case 2:
                                        label.SetText("Tem token?")
                                case 3:
                                        label.SetText("Lembrar?")
                                }
                                return
                        }
                        
                        // Dados
                        if id.Row-1 < len(config.PerfisZabbix) {
                                perfil := config.PerfisZabbix[id.Row-1]
                                switch id.Col {
                                case 0:
                                        label.SetText(perfil.Nome)
                                        if perfil.Nome == config.PerfilAtivo {
                                                label.TextStyle = fyne.TextStyle{Bold: true}
                                        }
                                case 1:
                                        label.SetText(perfil.URL)
                                case 2:
                                        if perfil.Token != "" {
                                                label.SetText("Sim")
                                        } else {
                                                label.SetText("Não")
                                        }
                                case 3:
                                        if perfil.Salvar {
                                                label.SetText("Sim")
                                        } else {
                                                label.SetText("Não")
                                        }
                                }
                        }
                },
        )
        
        // Ajustar tamanho das colunas
        tabelaPerfis.SetColumnWidth(0, 150)
        tabelaPerfis.SetColumnWidth(1, 250)
        tabelaPerfis.SetColumnWidth(2, 100)
        tabelaPerfis.SetColumnWidth(3, 100)
        
        // Selecionar um perfil na tabela
        tabelaPerfis.OnSelected = func(id widget.TableCellID) {
                if id.Row == 0 || id.Row-1 >= len(config.PerfisZabbix) {
                        return
                }
                
                indiceSelecionado = id.Row - 1
                perfilSelecionado = config.PerfisZabbix[indiceSelecionado]
        }
        
        // Formulário para editar/criar perfil
        nomePerfil := widget.NewEntry()
        nomePerfil.SetPlaceHolder("Nome do perfil")
        
        urlPerfil := widget.NewEntry()
        urlPerfil.SetPlaceHolder("URL da API do Zabbix")
        
        salvarPerfil := widget.NewCheck("Salvar credenciais", nil)
        
        // Botões de ação
        botaoNovo := widget.NewButton("Novo Perfil", func() {
                nomePerfil.SetText("")
                urlPerfil.SetText("http://localhost/zabbix/api_jsonrpc.php")
                salvarPerfil.SetChecked(false)
                indiceSelecionado = -1
        })
        
        botaoExcluir := widget.NewButton("Excluir", func() {
                if indiceSelecionado >= 0 && indiceSelecionado < len(config.PerfisZabbix) {
                        // Confirmar exclusão
                        dialog.ShowConfirm(
                                "Excluir perfil",
                                "Tem certeza que deseja excluir o perfil '" + perfilSelecionado.Nome + "'?",
                                func(confirma bool) {
                                        if confirma {
                                                fnExcluirPerfil(perfilSelecionado.Nome)
                                                tabelaPerfis.Refresh()
                                                indiceSelecionado = -1
                                        }
                                },
                                janela,
                        )
                }
        })
        
        botaoSalvar := widget.NewButton("Salvar", func() {
                if strings.TrimSpace(nomePerfil.Text) == "" || strings.TrimSpace(urlPerfil.Text) == "" {
                        dialog.ShowError(fyne.NewError(1, "Nome e URL são obrigatórios"), janela)
                        return
                }
                
                // Criar ou atualizar o perfil
                var perfil config.Perfil
                
                if indiceSelecionado >= 0 && indiceSelecionado < len(config.PerfisZabbix) {
                        // Atualizar perfil existente
                        perfil = perfilSelecionado
                        perfil.Nome = nomePerfil.Text
                        perfil.URL = urlPerfil.Text
                        perfil.Salvar = salvarPerfil.Checked
                } else {
                        // Novo perfil
                        perfil = config.Perfil{
                                Nome:     nomePerfil.Text,
                                URL:      urlPerfil.Text,
                                Token:    "",
                                Usuário:  "",
                                Senha:    "",
                                Salvar:   salvarPerfil.Checked,
                        }
                }
                
                fnAtualizarPerfil(perfil)
                tabelaPerfis.Refresh()
                
                dialog.ShowInformation("Sucesso", "Perfil salvo com sucesso!", janela)
        })
        
        botaoFechar := widget.NewButton("Fechar", func() {
                janela.Close()
        })
        
        // Atualizar os campos quando selecionar um perfil
        tabelaPerfis.OnSelected = func(id widget.TableCellID) {
                if id.Row == 0 || id.Row-1 >= len(config.PerfisZabbix) {
                        return
                }
                
                indiceSelecionado = id.Row - 1
                perfilSelecionado = config.PerfisZabbix[indiceSelecionado]
                
                nomePerfil.SetText(perfilSelecionado.Nome)
                urlPerfil.SetText(perfilSelecionado.URL)
                salvarPerfil.SetChecked(perfilSelecionado.Salvar)
        }
        
        // Layout
        formulario := widget.NewForm(
                widget.NewFormItem("Nome", nomePerfil),
                widget.NewFormItem("URL da API", urlPerfil),
                widget.NewFormItem("", salvarPerfil),
        )
        
        botoesAcao := container.NewHBox(
                botaoNovo,
                layout.NewSpacer(),
                botaoExcluir,
                botaoSalvar,
                botaoFechar,
        )
        
        conteudo := container.NewBorder(
                container.NewVBox(
                        widget.NewLabelWithStyle("Gerenciador de Perfis", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
                        widget.NewSeparator(),
                ),
                container.NewVBox(
                        widget.NewSeparator(),
                        botoesAcao,
                ),
                nil, nil,
                container.NewVSplit(
                        container.NewBorder(
                                nil, 
                                widget.NewLabel("Selecione um perfil na tabela para editar"),
                                nil, nil, 
                                tabelaPerfis,
                        ),
                        container.NewPadded(
                                container.NewVBox(
                                        widget.NewLabelWithStyle("Editar Perfil", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
                                        formulario,
                                ),
                        ),
                ),
        )
        
        return conteudo
}