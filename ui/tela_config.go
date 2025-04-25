package ui

import (
        "fmt"
        "path/filepath"
        "strconv"

        "fyne.io/fyne/v2"
        "fyne.io/fyne/v2/container"
        "fyne.io/fyne/v2/dialog"
        "fyne.io/fyne/v2/layout"
        "fyne.io/fyne/v2/widget"

        "zabbix-manager/config"
)

// criarTelaConfig cria a tela de configurações
func criarTelaConfig(janela fyne.Window, config *config.Configuração, perfilAtual config.Perfil, 
        fnSalvarConfig func(int, string, config.Perfil)) fyne.CanvasObject {
        
        // Perfil a ser editado (cópia do atual)
        perfilEditado := perfilAtual
        
        // Seção de configurações gerais
        entradaTempoLimite := widget.NewEntry()
        entradaTempoLimite.Text = formatarNumero(config.TempoLimite)
        entradaTempoLimite.SetPlaceHolder("Tempo limite (segundos)")

        // Diretório padrão para arquivos CSV
        entradaDiretorioCSV := widget.NewEntry()
        entradaDiretorioCSV.Text = config.DiretórioCSV
        entradaDiretorioCSV.SetPlaceHolder("Diretório para arquivos CSV")
        
        botaoSelecionarDiretorio := widget.NewButton("Selecionar...", func() {
                // Abrir diálogo de selecionar diretório
                dialogoDiretorio := dialog.NewFolderOpen(
                        func(uri fyne.ListableURI, err error) {
                                if err != nil {
                                        dialog.ShowError(err, janela)
                                        return
                                }
                                if uri == nil {
                                        return // Usuário cancelou
                                }
                                
                                // Obter o caminho do diretório
                                caminho := uri.Path()
                                entradaDiretorioCSV.SetText(caminho)
                        },
                        janela,
                )
                dialogoDiretorio.Show()
        })

        // Seção de configurações do perfil ativo
        entradaURL := widget.NewEntry()
        entradaURL.Text = perfilAtual.URL
        entradaURL.SetPlaceHolder("URL da API do Zabbix")
        
        entradaNome := widget.NewLabel(perfilAtual.Nome)
        
        infoToken := widget.NewLabel("Não há token configurado")
        if perfilAtual.Token != "" {
                infoToken.SetText("Token configurado")
        }
        
        infoSalvar := widget.NewLabel("Não salvar credenciais")
        if perfilAtual.Salvar {
                infoSalvar.SetText("Salvar credenciais")
        }
        
        // Botão para limpar o token
        botaoLimparToken := widget.NewButton("Limpar Token", func() {
                dialog.ShowConfirm(
                        "Limpar Token",
                        "Tem certeza que deseja limpar o token? Será necessário fazer login novamente.",
                        func(confirma bool) {
                                if confirma {
                                        perfilEditado.Token = ""
                                        infoToken.SetText("Não há token configurado")
                                }
                        },
                        janela,
                )
        })
        
        // Botão para limpar credenciais salvas
        botaoLimparCredenciais := widget.NewButton("Limpar Credenciais", func() {
                dialog.ShowConfirm(
                        "Limpar Credenciais",
                        "Tem certeza que deseja limpar as credenciais salvas?",
                        func(confirma bool) {
                                if confirma {
                                        perfilEditado.Usuário = ""
                                        perfilEditado.Senha = ""
                                        perfilEditado.Salvar = false
                                        infoSalvar.SetText("Não salvar credenciais")
                                }
                        },
                        janela,
                )
        })
        
        // Atualizar as informações ao editar a URL
        entradaURL.OnChanged = func(texto string) {
                perfilEditado.URL = texto
        }

        // Botões de ação
        botaoSalvar := widget.NewButton("Salvar", func() {
                // Validar o tempo limite
                tempoLimite, err := parseNumero(entradaTempoLimite.Text)
                if err != nil || tempoLimite <= 0 {
                        dialog.ShowError(fyne.NewError(1, "Tempo limite inválido. Use um número maior que zero."), janela)
                        return
                }
                
                // Validar a URL
                if entradaURL.Text == "" {
                        dialog.ShowError(fyne.NewError(1, "A URL da API é obrigatória"), janela)
                        return
                }
                
                // Atualizar o perfil com a URL
                perfilEditado.URL = entradaURL.Text
                
                // Salvar as configurações
                fnSalvarConfig(tempoLimite, entradaDiretorioCSV.Text, perfilEditado)
                
                // Fechar a janela
                dialog.ShowInformation("Sucesso", "Configurações salvas com sucesso!", janela)
                janela.Close()
        })

        botaoCancelar := widget.NewButton("Cancelar", func() {
                janela.Close()
        })

        // Formulário de configurações gerais
        formularioGeral := widget.NewForm(
                widget.NewFormItem("Tempo Limite (segundos)", entradaTempoLimite),
                widget.NewFormItem("Diretório CSV", container.NewBorder(nil, nil, nil, botaoSelecionarDiretorio, entradaDiretorioCSV)),
        )
        
        // Formulário de configurações do perfil
        formularioPerfil := widget.NewForm(
                widget.NewFormItem("Nome do Perfil", entradaNome),
                widget.NewFormItem("URL da API", entradaURL),
                widget.NewFormItem("Token", container.NewHBox(infoToken, layout.NewSpacer(), botaoLimparToken)),
                widget.NewFormItem("Credenciais", container.NewHBox(infoSalvar, layout.NewSpacer(), botaoLimparCredenciais)),
        )

        // Layout
        conteudo := container.NewVBox(
                widget.NewLabelWithStyle("Configurações", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
                widget.NewSeparator(),
                widget.NewLabelWithStyle("Configurações Gerais", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
                formularioGeral,
                widget.NewSeparator(),
                widget.NewLabelWithStyle(fmt.Sprintf("Configurações do Perfil: %s", perfilAtual.Nome), 
                        fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
                formularioPerfil,
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