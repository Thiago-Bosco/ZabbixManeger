package ui

import (
        "fmt"
        "path/filepath"
        "strconv"
        "strings"
        "time"

        "fyne.io/fyne/v2"
        "fyne.io/fyne/v2/container"
        "fyne.io/fyne/v2/dialog"
        "fyne.io/fyne/v2/layout"
        "fyne.io/fyne/v2/theme"
        "fyne.io/fyne/v2/widget"

        "zabbix-manager/zabbix"
)

// criarTelaPrincipal cria a tela principal da aplicação
func criarTelaPrincipal(janela fyne.Window, cliente *zabbix.ClienteAPI, fnAbrirConfig func(), fnAbrirLogin func()) (fyne.CanvasObject, func()) {
        var hosts []zabbix.Host
        var hostsSelecionados []zabbix.Host

        // Tabela de hosts
        tabelaHosts := widget.NewTable(
                func() (int, int) {
                        return len(hosts) + 1, 3 // Cabeçalho + linhas, 3 colunas
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
                                        label.SetText("ID")
                                case 1:
                                        label.SetText("Nome do Host")
                                case 2:
                                        label.SetText("Status")
                                }
                                return
                        }

                        // Dados
                        if id.Row-1 < len(hosts) {
                                host := hosts[id.Row-1]
                                switch id.Col {
                                case 0:
                                        label.SetText(host.ID)
                                case 1:
                                        label.SetText(host.Nome)
                                case 2:
                                        status := zabbix.StatusHost[host.Status]
                                        if status == "" {
                                                status = "Desconhecido"
                                        }
                                        label.SetText(status)
                                        
                                        // Cor de acordo com o status
                                        if host.Status == "0" {
                                                label.TextStyle = fyne.TextStyle{Bold: true}
                                        } else {
                                                label.TextStyle = fyne.TextStyle{Italic: true}
                                        }
                                }
                        }
                },
        )

        // Ajustar tamanho das colunas
        tabelaHosts.SetColumnWidth(0, 80)
        tabelaHosts.SetColumnWidth(1, 220)
        tabelaHosts.SetColumnWidth(2, 100)

        // Campo de busca
        campoBusca := widget.NewEntry()
        campoBusca.SetPlaceHolder("Buscar hosts...")
        
        botaoBuscar := widget.NewButtonWithIcon("Buscar", theme.SearchIcon(), func() {
                atualizarListaHosts(cliente, campoBusca.Text, func(h []zabbix.Host, err error) {
                        if err != nil {
                                dialog.ShowError(err, janela)
                                return
                        }
                        hosts = h
                        tabelaHosts.Refresh()
                })
        })

        // Detalhes do host
        labelHostSelecionado := widget.NewLabel("Selecione um host para ver os detalhes")
        labelHostSelecionado.Wrapping = fyne.TextWrapWord

        // Tabs para items e triggers
        listaItems := widget.NewList(
                func() int {
                        if len(hostsSelecionados) == 0 {
                                return 0
                        }
                        return len(hostsSelecionados[0].Items)
                },
                func() fyne.CanvasObject {
                        return container.NewHBox(
                                widget.NewIcon(theme.DocumentIcon()),
                                widget.NewLabel("Item"),
                        )
                },
                func(id widget.ListItemID, obj fyne.CanvasObject) {
                        if len(hostsSelecionados) == 0 || id >= len(hostsSelecionados[0].Items) {
                                return
                        }
                        
                        container := obj.(*fyne.Container)
                        label := container.Objects[1].(*widget.Label)
                        
                        item := hostsSelecionados[0].Items[id]
                        texto := fmt.Sprintf("%s: %s", item.ID, item.Nome)
                        label.SetText(texto)
                },
        )

        listaTriggers := widget.NewList(
                func() int {
                        if len(hostsSelecionados) == 0 {
                                return 0
                        }
                        return len(hostsSelecionados[0].Triggers)
                },
                func() fyne.CanvasObject {
                        return container.NewHBox(
                                widget.NewIcon(theme.WarningIcon()),
                                widget.NewLabel("Trigger"),
                        )
                },
                func(id widget.ListItemID, obj fyne.CanvasObject) {
                        if len(hostsSelecionados) == 0 || id >= len(hostsSelecionados[0].Triggers) {
                                return
                        }
                        
                        container := obj.(*fyne.Container)
                        label := container.Objects[1].(*widget.Label)
                        
                        trigger := hostsSelecionados[0].Triggers[id]
                        texto := fmt.Sprintf("%s: %s", trigger.ID, trigger.Nome)
                        label.SetText(texto)
                },
        )

        tabs := container.NewAppTabs(
                container.NewTabItem("Items", listaItems),
                container.NewTabItem("Triggers", listaTriggers),
        )

        painelDetalhes := container.NewVBox(
                labelHostSelecionado,
                widget.NewSeparator(),
                tabs,
        )

        // Evento ao selecionar um host na tabela
        tabelaHosts.OnSelected = func(id widget.TableCellID) {
                if id.Row == 0 || id.Row-1 >= len(hosts) {
                        return
                }

                host := hosts[id.Row-1]
                hostsSelecionados = []zabbix.Host{host}
                labelHostSelecionado.SetText(fmt.Sprintf("Host: %s (ID: %s)\nStatus: %s", 
                        host.Nome, 
                        host.ID, 
                        zabbix.StatusHost[host.Status]))
                
                listaItems.Refresh()
                listaTriggers.Refresh()
        }

        // Botão para exportar relatório
        botaoExportar := widget.NewButtonWithIcon("Exportar Relatório", theme.DocumentSaveIcon(), func() {
                if len(hosts) == 0 {
                        dialog.ShowInformation("Informação", "Não há dados para exportar.", janela)
                        return
                }

                // Diálogo para salvar arquivo
                dialogoSalvar := dialog.NewFileSave(
                        func(escritor fyne.URIWriteCloser, err error) {
                                if err != nil {
                                        dialog.ShowError(err, janela)
                                        return
                                }
                                if escritor == nil {
                                        return // Usuário cancelou
                                }
                                defer escritor.Close()

                                // Obter o caminho do arquivo
                                caminho := escritor.URI().Path()
                                
                                // Garantir que tenha a extensão .csv
                                if !strings.HasSuffix(strings.ToLower(caminho), ".csv") {
                                        caminho += ".csv"
                                }
                                
                                // Gerar o relatório
                                err = zabbix.GerarRelatorioCSV(hosts, caminho)
                                if err != nil {
                                        dialog.ShowError(fmt.Errorf("Erro ao gerar relatório: %v", err), janela)
                                        return
                                }

                                dialog.ShowInformation("Sucesso", "Relatório gerado com sucesso!", janela)
                        },
                        janela,
                )
                
                // Sugerir um nome de arquivo padrão
                dataAtual := time.Now().Format("2006-01-02")
                dialogoSalvar.SetFileName(fmt.Sprintf("Relatorio_Zabbix_%s.csv", dataAtual))
                dialogoSalvar.SetFilter(storage.NewExtensionFileFilter([]string{".csv"}))
                dialogoSalvar.Show()
        })

        // Botão de atualizar dados
        botaoAtualizar := widget.NewButtonWithIcon("Atualizar", theme.ViewRefreshIcon(), func() {
                atualizarListaHosts(cliente, campoBusca.Text, func(h []zabbix.Host, err error) {
                        if err != nil {
                                dialog.ShowError(err, janela)
                                return
                        }
                        hosts = h
                        tabelaHosts.Refresh()
                        
                        // Limpar seleção
                        hostsSelecionados = []zabbix.Host{}
                        labelHostSelecionado.SetText("Selecione um host para ver os detalhes")
                        listaItems.Refresh()
                        listaTriggers.Refresh()
                })
        })

        // Botão de configurações
        botaoConfig := widget.NewButtonWithIcon("Configurações", theme.SettingsIcon(), fnAbrirConfig)
        
        // Botão para alterar conta/perfil
        botaoAlterarConta := widget.NewButtonWithIcon("Alterar Conta", theme.AccountIcon(), fnAbrirLogin)

        // Barra de status
        barraStatus := widget.NewLabel("Pronto")

        // Layout da tela principal
        split := container.NewHSplit(
                container.NewBorder(
                        container.NewVBox(
                                container.NewHBox(
                                        widget.NewLabel("Filtro:"),
                                        campoBusca,
                                        botaoBuscar,
                                ),
                        ),
                        nil, nil, nil,
                        tabelaHosts,
                ),
                painelDetalhes,
        )
        split.SetOffset(0.6)

        conteudo := container.NewBorder(
                nil,
                container.NewBorder(
                        nil, nil, nil, nil,
                        container.NewHBox(
                                barraStatus,
                                layout.NewSpacer(),
                                botaoAtualizar,
                                botaoExportar,
                                botaoAlterarConta,
                                botaoConfig,
                        ),
                ),
                nil, nil,
                split,
        )

        // Função para atualizar os dados
        fnAtualizar := func() {
                barraStatus.SetText("Carregando hosts...")
                atualizarListaHosts(cliente, "", func(h []zabbix.Host, err error) {
                        if err != nil {
                                barraStatus.SetText(fmt.Sprintf("Erro: %v", err))
                                dialog.ShowError(err, janela)
                                return
                        }
                        hosts = h
                        tabelaHosts.Refresh()
                        barraStatus.SetText(fmt.Sprintf("Carregados %d hosts", len(hosts)))
                })
        }

        // Configurar menu da aplicação
        janela.SetMainMenu(fyne.NewMainMenu(
                fyne.NewMenu("Arquivo",
                        fyne.NewMenuItem("Atualizar", func() {
                                botaoAtualizar.OnTapped()
                        }),
                        fyne.NewMenuItem("Exportar Relatório", func() {
                                botaoExportar.OnTapped()
                        }),
                        fyne.NewMenuItemSeparator(),
                        fyne.NewMenuItem("Sair", func() {
                                janela.Close()
                        }),
                ),
                fyne.NewMenu("Conta",
                        fyne.NewMenuItem("Alterar Conta/Perfil", func() {
                                fnAbrirLogin()
                        }),
                        fyne.NewMenuItem("Configurações", func() {
                                fnAbrirConfig()
                        }),
                ),
                fyne.NewMenu("Ajuda",
                        fyne.NewMenuItem("Sobre", func() {
                                dialog.ShowInformation("Sobre", "Zabbix Manager\nVersão 1.0\n\nUma aplicação para gerenciamento do Zabbix.", janela)
                        }),
                ),
        ))

        return conteudo, fnAtualizar
}

// atualizarListaHosts atualiza a lista de hosts
func atualizarListaHosts(cliente *zabbix.ClienteAPI, filtro string, callback func([]zabbix.Host, error)) {
        go func() {
                var hosts []zabbix.Host
                var err error

                if filtro == "" {
                        hosts, err = cliente.ObterHosts()
                } else {
                        hosts, err = cliente.ObterHostsFiltrados(filtro)
                }

                // Chamada de retorno na thread principal da UI
                fyne.CurrentApp().Driver().RunOnMain(func() {
                        callback(hosts, err)
                })
        }()
}