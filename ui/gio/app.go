package gio

import (
        "image/color"
        "log"
        "time"

        "gioui.org/app"
        "gioui.org/font/gofont"
        "gioui.org/io/system"
        "gioui.org/layout"
        "gioui.org/op"
        "gioui.org/text"
        "gioui.org/unit"
        "gioui.org/widget"
        "gioui.org/widget/material"

        "zabbix-manager/config"
        "zabbix-manager/zabbix"
)

// Aplicacao representa a aplicação Gio para o Zabbix Manager
type Aplicacao struct {
        janela            *app.Window
        tema              *material.Theme
        config            *config.Configuração
        clienteAPI        *zabbix.ClienteAPI
        perfilAtual       *config.ConfiguracaoPerfil
        arquivoConfig     string
        estado            Estado
        estadoAnterior    Estado
        hosts             []zabbix.Host
        hostsFiltrados    []zabbix.Host
        listaPerfis       []config.ConfiguracaoPerfil
        indicePerfilAtual int
        
        // Widgets comuns
        txtBusca          *widget.Editor
        btnBusca          *widget.Clickable
        btnVoltar         *widget.Clickable
        btnSalvar         *widget.Clickable
        
        // Widgets tela de login
        txtNomePerfil     *widget.Editor
        txtURLAPI         *widget.Editor
        txtToken          *widget.Editor
        btnEntrar         *widget.Clickable
        btnAdicionar      *widget.Clickable
        btnEditar         *widget.Clickable
        btnRemover        *widget.Clickable
        lstPerfis         *widget.List
        
        // Widgets tela principal
        btnAtualizar      *widget.Clickable
        btnExportar       *widget.Clickable
        btnConfiguracoes  *widget.Clickable
        btnSair           *widget.Clickable
        lstHosts          *widget.List
        
        // Estado de mensagens
        mensagem          string
        tipoMensagem      TipoMensagem
        tempoMensagem     time.Time
        
        // Dimensões
        tamanhoPadrao     unit.Dp
}

// Estado representa o estado atual da aplicação
type Estado int

// Estados possíveis da aplicação
const (
        EstadoLogin Estado = iota
        EstadoPrincipal
        EstadoConfiguracao
)

// TipoMensagem define o tipo de mensagem a ser exibida
type TipoMensagem int

// Tipos de mensagem
const (
        MensagemInfo TipoMensagem = iota
        MensagemErro
        MensagemSucesso
)

// NovaAplicacao cria uma nova aplicação Gio
func NovaAplicacao() *Aplicacao {
        // Criar janela
        janela := app.NewWindow(
                app.Title("Zabbix Manager"),
                app.Size(unit.Dp(1000), unit.Dp(600)),
        )

        // Obter caminho do arquivo de configuração
        arquivoConfig := config.ObterCaminhoConfiguracao()

        // Carregar configuração
        configuracao, err := config.Carregar(arquivoConfig)
        if err != nil {
                log.Printf("Erro ao carregar configuração: %v", err)
                configuracao = config.NovaConfiguração()
        }

        // Verificar se há algum perfil ativo
        var perfilAtivo *config.ConfiguracaoPerfil
        if configuracao.PerfilAtual >= 0 && configuracao.PerfilAtual < len(configuracao.Perfis) {
                perfilAtivo, _ = configuracao.PerfilAtivo()
        }

        // Criar aplicação
        app := &Aplicacao{
                janela:            janela,
                tema:              material.NewTheme(gofont.Collection()),
                config:            configuracao,
                perfilAtual:       perfilAtivo,
                arquivoConfig:     arquivoConfig,
                estado:            EstadoLogin,
                listaPerfis:       configuracao.Perfis,
                indicePerfilAtual: configuracao.PerfilAtual,
                tamanhoPadrao:     unit.Dp(16),
                
                // Inicializar widgets
                txtBusca:          &widget.Editor{SingleLine: true},
                btnBusca:          &widget.Clickable{},
                btnVoltar:         &widget.Clickable{},
                btnSalvar:         &widget.Clickable{},
                
                txtNomePerfil:     &widget.Editor{SingleLine: true, Hint: "Nome do perfil"},
                txtURLAPI:         &widget.Editor{SingleLine: true, Hint: "URL da API (ex: http://zabbix.exemplo.com/api_jsonrpc.php)"},
                txtToken:          &widget.Editor{SingleLine: true, Hint: "Token da API"},
                btnEntrar:         &widget.Clickable{},
                btnAdicionar:      &widget.Clickable{},
                btnEditar:         &widget.Clickable{},
                btnRemover:        &widget.Clickable{},
                lstPerfis:         &widget.List{List: layout.List{Axis: layout.Vertical}},
                
                btnAtualizar:      &widget.Clickable{},
                btnExportar:       &widget.Clickable{},
                btnConfiguracoes:  &widget.Clickable{},
                btnSair:           &widget.Clickable{},
                lstHosts:          &widget.List{List: layout.List{Axis: layout.Vertical}},
        }

        // Configurar cliente API se houver perfil ativo
        if perfilAtivo != nil {
                app.ConfigurarClienteAPI(perfilAtivo)
        }

        return app
}

// ConfigurarClienteAPI configura o cliente da API do Zabbix com os dados do perfil
func (a *Aplicacao) ConfigurarClienteAPI(perfil *config.ConfiguracaoPerfil) {
        configAPI := zabbix.ConfigAPI{
                URL:         perfil.URL,
                Token:       perfil.Token,
                TempoLimite: a.config.TempoLimite,
        }
        a.clienteAPI = zabbix.NovoClienteAPI(configAPI)
        a.perfilAtual = perfil
}

// MostrarErro exibe uma mensagem de erro
func (a *Aplicacao) MostrarErro(msg string) {
        a.mensagem = msg
        a.tipoMensagem = MensagemErro
        a.tempoMensagem = time.Now()
}

// MostrarInfo exibe uma mensagem informativa
func (a *Aplicacao) MostrarInfo(msg string) {
        a.mensagem = msg
        a.tipoMensagem = MensagemInfo
        a.tempoMensagem = time.Now()
}

// MostrarSucesso exibe uma mensagem de sucesso
func (a *Aplicacao) MostrarSucesso(msg string) {
        a.mensagem = msg
        a.tipoMensagem = MensagemSucesso
        a.tempoMensagem = time.Now()
}

// Executar inicia a execução da aplicação
func (a *Aplicacao) Executar() error {
        var ops op.Ops

        // Definir estado inicial
        if len(a.config.Perfis) == 0 {
                a.estado = EstadoLogin
        } else if a.config.PerfilAtual < 0 || a.config.PerfilAtual >= len(a.config.Perfis) {
                a.estado = EstadoLogin
        } else {
                // Verificar conexão com o servidor
                if a.clienteAPI != nil {
                        err := a.clienteAPI.TestarConexao()
                        if err == nil {
                                a.estado = EstadoPrincipal
                                // Carregar hosts
                                a.CarregarHosts()
                        } else {
                                a.estado = EstadoLogin
                                a.MostrarErro("Erro de conexão: " + err.Error())
                        }
                } else {
                        a.estado = EstadoLogin
                }
        }

        // Loop principal
        for {
                e := <-a.janela.Events()
                
                switch e := e.(type) {
                case system.DestroyEvent:
                        return e.Err
                        
                case system.FrameEvent:
                        // Iniciar operações da GUI
                        gtx := layout.NewContext(&ops, e)
                        
                        // Renderizar interface de acordo com o estado
                        switch a.estado {
                        case EstadoLogin:
                                a.RenderizarTelaLogin(gtx)
                        case EstadoPrincipal:
                                a.RenderizarTelaPrincipal(gtx)
                        case EstadoConfiguracao:
                                a.RenderizarTelaConfiguracao(gtx)
                        }
                        
                        // Desenhar
                        e.Frame(gtx.Ops)
                }
        }
}

// RenderizarTelaLogin renderiza a tela de login
func (a *Aplicacao) RenderizarTelaLogin(gtx layout.Context) {
        // Verificar eventos dos botões
        if a.btnEntrar.Clicked() {
                a.acaoEntrar()
        }
        if a.btnAdicionar.Clicked() {
                a.acaoAdicionarPerfil()
        }
        if a.btnEditar.Clicked() {
                a.acaoEditarPerfil()
        }
        if a.btnRemover.Clicked() {
                a.acaoRemoverPerfil()
        }

        // Layout principal
        layout.Flex{Axis: layout.Vertical}.Layout(gtx,
                layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                        return layout.Inset{Top: a.tamanhoPadrao, Bottom: a.tamanhoPadrao}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                                titulo := material.H4(a.tema, "Zabbix Manager")
                                titulo.Alignment = text.Middle
                                return titulo.Layout(gtx)
                        })
                }),
                layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                        return layout.Inset{
                                Top:    a.tamanhoPadrao,
                                Bottom: a.tamanhoPadrao,
                                Left:   a.tamanhoPadrao,
                                Right:  a.tamanhoPadrao,
                        }.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                                return a.renderizarFormularioLogin(gtx)
                        })
                }),
                layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                        return layout.Inset{
                                Bottom: a.tamanhoPadrao,
                                Left:   a.tamanhoPadrao,
                                Right:  a.tamanhoPadrao,
                        }.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                                return a.renderizarMensagem(gtx)
                        })
                }),
        )
}

// renderizarFormularioLogin renderiza o formulário de login
func (a *Aplicacao) renderizarFormularioLogin(gtx layout.Context) layout.Dimensions {
        if len(a.listaPerfis) > 0 {
                // Se existem perfis, mostrar a lista
                return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
                        layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                                return material.H6(a.tema, "Selecione um perfil:").Layout(gtx)
                        }),
                        layout.Rigid(layout.Spacer{Height: a.tamanhoPadrao}.Layout),
                        layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                                return a.renderizarListaPerfis(gtx)
                        }),
                        layout.Rigid(layout.Spacer{Height: a.tamanhoPadrao}.Layout),
                        layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                                return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
                                        layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
                                                return material.Button(a.tema, a.btnEntrar, "Entrar").Layout(gtx)
                                        }),
                                )
                        }),
                        layout.Rigid(layout.Spacer{Height: a.tamanhoPadrao}.Layout),
                        layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                                return material.H6(a.tema, "Adicionar ou editar perfil:").Layout(gtx)
                        }),
                        layout.Rigid(layout.Spacer{Height: a.tamanhoPadrao}.Layout),
                        layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                                return material.Body1(a.tema, "Nome:").Layout(gtx)
                        }),
                        layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                                return material.Editor(a.tema, a.txtNomePerfil, "Nome para identificar este servidor").Layout(gtx)
                        }),
                        layout.Rigid(layout.Spacer{Height: a.tamanhoPadrao/2}.Layout),
                        layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                                return material.Body1(a.tema, "URL da API:").Layout(gtx)
                        }),
                        layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                                return material.Editor(a.tema, a.txtURLAPI, "http://seu-servidor-zabbix/api_jsonrpc.php").Layout(gtx)
                        }),
                        layout.Rigid(layout.Spacer{Height: a.tamanhoPadrao/2}.Layout),
                        layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                                return material.Body1(a.tema, "Token da API:").Layout(gtx)
                        }),
                        layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                                return material.Editor(a.tema, a.txtToken, "Token de API do Zabbix").Layout(gtx)
                        }),
                        layout.Rigid(layout.Spacer{Height: a.tamanhoPadrao}.Layout),
                        layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                                return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
                                        layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
                                                return material.Button(a.tema, a.btnAdicionar, "Adicionar").Layout(gtx)
                                        }),
                                        layout.Rigid(layout.Spacer{Width: a.tamanhoPadrao}.Layout),
                                        layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
                                                return material.Button(a.tema, a.btnEditar, "Editar").Layout(gtx)
                                        }),
                                        layout.Rigid(layout.Spacer{Width: a.tamanhoPadrao}.Layout),
                                        layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
                                                return material.Button(a.tema, a.btnRemover, "Remover").Layout(gtx)
                                        }),
                                )
                        }),
                )
        } else {
                // Se não existem perfis, mostrar apenas o formulário de cadastro
                return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
                        layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                                return material.H6(a.tema, "Não há perfis cadastrados. Adicione um novo:").Layout(gtx)
                        }),
                        layout.Rigid(layout.Spacer{Height: a.tamanhoPadrao}.Layout),
                        layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                                return material.Body1(a.tema, "Nome:").Layout(gtx)
                        }),
                        layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                                return material.Editor(a.tema, a.txtNomePerfil, "Nome para identificar este servidor").Layout(gtx)
                        }),
                        layout.Rigid(layout.Spacer{Height: a.tamanhoPadrao/2}.Layout),
                        layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                                return material.Body1(a.tema, "URL da API:").Layout(gtx)
                        }),
                        layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                                return material.Editor(a.tema, a.txtURLAPI, "http://seu-servidor-zabbix/api_jsonrpc.php").Layout(gtx)
                        }),
                        layout.Rigid(layout.Spacer{Height: a.tamanhoPadrao/2}.Layout),
                        layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                                return material.Body1(a.tema, "Token da API:").Layout(gtx)
                        }),
                        layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                                return material.Editor(a.tema, a.txtToken, "Token de API do Zabbix").Layout(gtx)
                        }),
                        layout.Rigid(layout.Spacer{Height: a.tamanhoPadrao}.Layout),
                        layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                                return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
                                        layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
                                                return material.Button(a.tema, a.btnAdicionar, "Adicionar").Layout(gtx)
                                        }),
                                )
                        }),
                )
        }
}

// renderizarListaPerfis renderiza a lista de perfis
func (a *Aplicacao) renderizarListaPerfis(gtx layout.Context) layout.Dimensions {
        return layout.List{
                Axis: layout.Vertical,
        }.Layout(gtx, len(a.listaPerfis), func(gtx layout.Context, i int) layout.Dimensions {
                perfil := a.listaPerfis[i]
                selected := i == a.indicePerfilAtual
                
                var colorBg color.NRGBA
                if selected {
                        colorBg = color.NRGBA{R: 200, G: 200, B: 200, A: 255}
                } else {
                        colorBg = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
                }
                
                return layout.Background{
                        Color: colorBg,
                }.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                        return layout.UniformInset(a.tamanhoPadrao/2).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                                return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
                                        layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                                                return material.Body1(a.tema, perfil.Nome).Layout(gtx)
                                        }),
                                        layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                                                return material.Caption(a.tema, perfil.URL).Layout(gtx)
                                        }),
                                )
                        })
                })
        })
}

// renderizarMensagem renderiza a mensagem de feedback
func (a *Aplicacao) renderizarMensagem(gtx layout.Context) layout.Dimensions {
        // Verificar se há mensagem para exibir
        if a.mensagem == "" || time.Since(a.tempoMensagem) > 5*time.Second {
                return layout.Dimensions{}
        }
        
        var cor color.NRGBA
        switch a.tipoMensagem {
        case MensagemErro:
                cor = color.NRGBA{R: 255, G: 0, B: 0, A: 255}
        case MensagemSucesso:
                cor = color.NRGBA{R: 0, G: 200, B: 0, A: 255}
        case MensagemInfo:
                cor = color.NRGBA{R: 0, G: 0, B: 200, A: 255}
        }
        
        return layout.Background{
                Color: color.NRGBA{R: 240, G: 240, B: 240, A: 255},
        }.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                return layout.UniformInset(a.tamanhoPadrao).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                        msg := material.Body1(a.tema, a.mensagem)
                        msg.Color = cor
                        return msg.Layout(gtx)
                })
        })
}

// RenderizarTelaPrincipal renderiza a tela principal
func (a *Aplicacao) RenderizarTelaPrincipal(gtx layout.Context) {
        // Verificar eventos dos botões
        if a.btnAtualizar.Clicked() {
                a.CarregarHosts()
        }
        if a.btnExportar.Clicked() {
                a.ExportarRelatorio()
        }
        if a.btnConfiguracoes.Clicked() {
                a.estado = EstadoConfiguracao
        }
        if a.btnSair.Clicked() {
                a.estado = EstadoLogin
        }
        if a.btnBusca.Clicked() || a.txtBusca.Submitted() {
                a.AplicarFiltro(a.txtBusca.Text())
        }

        // Layout principal
        layout.Flex{Axis: layout.Vertical}.Layout(gtx,
                layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                        return layout.Inset{
                                Top:    a.tamanhoPadrao,
                                Bottom: a.tamanhoPadrao,
                                Left:   a.tamanhoPadrao,
                                Right:  a.tamanhoPadrao,
                        }.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                                return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
                                        layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
                                                return material.Editor(a.tema, a.txtBusca, "Buscar hosts...").Layout(gtx)
                                        }),
                                        layout.Rigid(layout.Spacer{Width: a.tamanhoPadrao/2}.Layout),
                                        layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                                                return material.Button(a.tema, a.btnBusca, "Buscar").Layout(gtx)
                                        }),
                                        layout.Rigid(layout.Spacer{Width: a.tamanhoPadrao}.Layout),
                                        layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                                                return material.Button(a.tema, a.btnAtualizar, "Atualizar").Layout(gtx)
                                        }),
                                        layout.Rigid(layout.Spacer{Width: a.tamanhoPadrao/2}.Layout),
                                        layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                                                return material.Button(a.tema, a.btnExportar, "Exportar CSV").Layout(gtx)
                                        }),
                                        layout.Rigid(layout.Spacer{Width: a.tamanhoPadrao/2}.Layout),
                                        layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                                                return material.Button(a.tema, a.btnConfiguracoes, "Configurar").Layout(gtx)
                                        }),
                                        layout.Rigid(layout.Spacer{Width: a.tamanhoPadrao/2}.Layout),
                                        layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                                                return material.Button(a.tema, a.btnSair, "Sair").Layout(gtx)
                                        }),
                                )
                        })
                }),
                layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
                        return layout.Inset{
                                Bottom: a.tamanhoPadrao,
                                Left:   a.tamanhoPadrao,
                                Right:  a.tamanhoPadrao,
                        }.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                                return a.renderizarTabelaHosts(gtx)
                        })
                }),
                layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                        return layout.Inset{
                                Bottom: a.tamanhoPadrao,
                                Left:   a.tamanhoPadrao,
                                Right:  a.tamanhoPadrao,
                        }.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                                return a.renderizarMensagem(gtx)
                        })
                }),
        )
}

// renderizarTabelaHosts renderiza a tabela de hosts
func (a *Aplicacao) renderizarTabelaHosts(gtx layout.Context) layout.Dimensions {
        // Primeiro, renderizar o cabeçalho
        header := layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
                layout.Flexed(0.2, func(gtx layout.Context) layout.Dimensions {
                        text := material.Body1(a.tema, "ID")
                        text.Font.Weight = text.Font.Weight + 300 // Tornar negrito
                        return text.Layout(gtx)
                }),
                layout.Flexed(0.4, func(gtx layout.Context) layout.Dimensions {
                        text := material.Body1(a.tema, "Nome")
                        text.Font.Weight = text.Font.Weight + 300
                        return text.Layout(gtx)
                }),
                layout.Flexed(0.2, func(gtx layout.Context) layout.Dimensions {
                        text := material.Body1(a.tema, "Status")
                        text.Font.Weight = text.Font.Weight + 300
                        return text.Layout(gtx)
                }),
                layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
                        text := material.Body1(a.tema, "Itens")
                        text.Font.Weight = text.Font.Weight + 300
                        return text.Layout(gtx)
                }),
                layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
                        text := material.Body1(a.tema, "Triggers")
                        text.Font.Weight = text.Font.Weight + 300
                        return text.Layout(gtx)
                }),
        )

        // Em seguida, a lista de hosts
        list := layout.List{
                Axis: layout.Vertical,
        }.Layout(gtx, len(a.hostsFiltrados), func(gtx layout.Context, i int) layout.Dimensions {
                host := a.hostsFiltrados[i]
                
                // Alternar cores de fundo para linhas pares/ímpares
                var colorBg color.NRGBA
                if i%2 == 0 {
                        colorBg = color.NRGBA{R: 240, G: 240, B: 245, A: 255}
                } else {
                        colorBg = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
                }
                
                return layout.Background{
                        Color: colorBg,
                }.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                        return layout.UniformInset(a.tamanhoPadrao/2).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                                // Obter status como texto
                                status := zabbix.StatusHost[host.Status]
                                if status == "" {
                                        status = "Desconhecido"
                                }
                                
                                return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
                                        layout.Flexed(0.2, func(gtx layout.Context) layout.Dimensions {
                                                return material.Body2(a.tema, host.ID).Layout(gtx)
                                        }),
                                        layout.Flexed(0.4, func(gtx layout.Context) layout.Dimensions {
                                                return material.Body2(a.tema, host.Nome).Layout(gtx)
                                        }),
                                        layout.Flexed(0.2, func(gtx layout.Context) layout.Dimensions {
                                                return material.Body2(a.tema, status).Layout(gtx)
                                        }),
                                        layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
                                                return material.Body2(a.tema, fmt.Sprintf("%d", len(host.Items))).Layout(gtx)
                                        }),
                                        layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
                                                return material.Body2(a.tema, fmt.Sprintf("%d", len(host.Triggers))).Layout(gtx)
                                        }),
                                )
                        })
                })
        })

        // Combinar cabeçalho e lista em um container vertical
        return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
                layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                        return layout.Background{
                                Color: color.NRGBA{R: 220, G: 220, B: 220, A: 255},
                        }.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                                return layout.UniformInset(a.tamanhoPadrao/2).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                                        return header
                                })
                        })
                }),
                layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
                        return list
                }),
        )
}

// RenderizarTelaConfiguracao renderiza a tela de configuração
func (a *Aplicacao) RenderizarTelaConfiguracao(gtx layout.Context) {
        // Verificar eventos dos botões
        if a.btnSalvar.Clicked() {
                a.acaoSalvarConfiguracao()
        }
        if a.btnVoltar.Clicked() {
                a.estado = EstadoPrincipal
        }

        // Layout principal
        layout.Flex{Axis: layout.Vertical}.Layout(gtx,
                layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                        return layout.Inset{Top: a.tamanhoPadrao, Bottom: a.tamanhoPadrao}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                                titulo := material.H4(a.tema, "Configurações")
                                titulo.Alignment = text.Middle
                                return titulo.Layout(gtx)
                        })
                }),
                layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                        return layout.Inset{
                                Top:    a.tamanhoPadrao,
                                Bottom: a.tamanhoPadrao,
                                Left:   a.tamanhoPadrao,
                                Right:  a.tamanhoPadrao,
                        }.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                                return a.renderizarFormularioConfiguracao(gtx)
                        })
                }),
                layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                        return layout.Inset{
                                Bottom: a.tamanhoPadrao,
                                Left:   a.tamanhoPadrao,
                                Right:  a.tamanhoPadrao,
                        }.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                                return a.renderizarMensagem(gtx)
                        })
                }),
        )
}

// renderizarFormularioConfiguracao renderiza o formulário de configuração
func (a *Aplicacao) renderizarFormularioConfiguracao(gtx layout.Context) layout.Dimensions {
        return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
                layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                        return material.H6(a.tema, "Configurações Gerais").Layout(gtx)
                }),
                layout.Rigid(layout.Spacer{Height: a.tamanhoPadrao}.Layout),
                layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                        return material.Body1(a.tema, "Tempo limite de requisição (segundos):").Layout(gtx)
                }),
                layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                        // TODO: Implementar campo para tempo limite
                        return layout.Dimensions{Size: gtx.Constraints.Min}
                }),
                layout.Rigid(layout.Spacer{Height: a.tamanhoPadrao}.Layout),
                layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                        return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
                                layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
                                        return material.Button(a.tema, a.btnVoltar, "Voltar").Layout(gtx)
                                }),
                                layout.Rigid(layout.Spacer{Width: a.tamanhoPadrao}.Layout),
                                layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
                                        return material.Button(a.tema, a.btnSalvar, "Salvar").Layout(gtx)
                                }),
                        )
                }),
        )
}

// CarregarHosts carrega os hosts do servidor Zabbix
func (a *Aplicacao) CarregarHosts() {
        if a.clienteAPI == nil {
                a.MostrarErro("Nenhum servidor Zabbix configurado")
                return
        }

        // Buscar hosts
        a.MostrarInfo("Carregando hosts...")
        hosts, err := a.clienteAPI.ObterHosts()
        if err != nil {
                a.MostrarErro("Erro ao obter hosts: " + err.Error())
                return
        }

        a.hosts = hosts
        a.AplicarFiltro(a.txtBusca.Text())
        a.MostrarSucesso(fmt.Sprintf("Carregados %d hosts", len(hosts)))
}

// AplicarFiltro aplica filtro de busca aos hosts
func (a *Aplicacao) AplicarFiltro(filtro string) {
        if filtro == "" {
                a.hostsFiltrados = a.hosts
                return
        }

        filtrados := []zabbix.Host{}
        for _, host := range a.hosts {
                if strings.Contains(strings.ToLower(host.Nome), strings.ToLower(filtro)) ||
                        strings.Contains(strings.ToLower(host.ID), strings.ToLower(filtro)) {
                        filtrados = append(filtrados, host)
                }
        }

        a.hostsFiltrados = filtrados
        a.MostrarInfo(fmt.Sprintf("Mostrando %d de %d hosts", len(a.hostsFiltrados), len(a.hosts)))
}

// ExportarRelatorio exporta os dados para CSV
func (a *Aplicacao) ExportarRelatorio() {
        if len(a.hosts) == 0 {
                a.MostrarErro("Não há dados para exportar")
                return
        }

        // Obter diretório home do usuário
        diretorioHome, err := os.UserHomeDir()
        if err != nil {
                diretorioHome, _ = os.Getwd()
        }

        // Criar diretório para relatórios
        diretorioRelatorios := filepath.Join(diretorioHome, "Relatórios Zabbix")
        err = os.MkdirAll(diretorioRelatorios, 0755)
        if err != nil {
                a.MostrarErro("Erro ao criar diretório para relatórios: " + err.Error())
                return
        }

        // Nome do arquivo
        nomeServidor := "zabbix"
        if a.perfilAtual != nil {
                nomeServidor = a.perfilAtual.Nome
        }
        caminhoArquivo := filepath.Join(diretorioRelatorios, fmt.Sprintf("relatorio_%s.csv", nomeServidor))

        // Gerar relatório
        err = zabbix.GerarRelatorioCSV(a.hosts, caminhoArquivo)
        if err != nil {
                a.MostrarErro("Erro ao gerar relatório: " + err.Error())
                return
        }

        a.MostrarSucesso(fmt.Sprintf("Relatório exportado para: %s", caminhoArquivo))
}

// acaoEntrar processa a ação de entrar no sistema
func (a *Aplicacao) acaoEntrar() {
        if len(a.listaPerfis) == 0 {
                a.MostrarErro("Não há perfis configurados")
                return
        }

        indice := a.indicePerfilAtual
        if indice < 0 || indice >= len(a.listaPerfis) {
                a.MostrarErro("Selecione um perfil válido")
                return
        }

        // Selecionar o perfil
        err := a.config.SelecionarPerfil(indice)
        if err != nil {
                a.MostrarErro("Erro ao selecionar perfil: " + err.Error())
                return
        }

        // Configurar cliente API
        perfil := a.listaPerfis[indice]
        a.ConfigurarClienteAPI(&perfil)

        // Testar conexão
        err = a.clienteAPI.TestarConexao()
        if err != nil {
                a.MostrarErro("Erro de conexão: " + err.Error())
                return
        }

        // Salvar configuração
        err = a.config.Salvar(a.arquivoConfig)
        if err != nil {
                a.MostrarErro("Erro ao salvar configuração: " + err.Error())
                return
        }

        // Ir para a tela principal
        a.estado = EstadoPrincipal
        a.CarregarHosts()
}

// acaoAdicionarPerfil processa a ação de adicionar um perfil
func (a *Aplicacao) acaoAdicionarPerfil() {
        // Obter dados dos campos
        nome := a.txtNomePerfil.Text()
        url := a.txtURLAPI.Text()
        token := a.txtToken.Text()

        // Validar campos
        if nome == "" || url == "" || token == "" {
                a.MostrarErro("Todos os campos são obrigatórios")
                return
        }

        // Testar conexão
        configAPI := zabbix.ConfigAPI{
                URL:         url,
                Token:       token,
                TempoLimite: a.config.TempoLimite,
        }
        clienteTemporario := zabbix.NovoClienteAPI(configAPI)
        err := clienteTemporario.TestarConexao()
        if err != nil {
                a.MostrarErro("Erro de conexão: " + err.Error())
                return
        }

        // Criar o novo perfil
        perfil := config.ConfiguracaoPerfil{
                Nome:  nome,
                URL:   url,
                Token: token,
        }

        // Adicionar o perfil à configuração
        a.config.AdicionarPerfil(perfil)

        // Atualizar a lista de perfis
        a.listaPerfis = a.config.Perfis
        a.indicePerfilAtual = a.config.PerfilAtual

        // Salvar a configuração
        err = a.config.Salvar(a.arquivoConfig)
        if err != nil {
                a.MostrarErro("Erro ao salvar configuração: " + err.Error())
                return
        }

        // Limpar campos
        a.txtNomePerfil.SetText("")
        a.txtURLAPI.SetText("")
        a.txtToken.SetText("")

        a.MostrarSucesso("Perfil adicionado com sucesso")
}

// acaoEditarPerfil processa a ação de editar um perfil
func (a *Aplicacao) acaoEditarPerfil() {
        if len(a.listaPerfis) == 0 {
                a.MostrarErro("Não há perfis para editar")
                return
        }

        indice := a.indicePerfilAtual
        if indice < 0 || indice >= len(a.listaPerfis) {
                a.MostrarErro("Selecione um perfil válido")
                return
        }

        // Obter dados dos campos
        nome := a.txtNomePerfil.Text()
        url := a.txtURLAPI.Text()
        token := a.txtToken.Text()

        // Validar campos
        if nome == "" || url == "" || token == "" {
                a.MostrarErro("Todos os campos são obrigatórios")
                return
        }

        // Testar conexão
        configAPI := zabbix.ConfigAPI{
                URL:         url,
                Token:       token,
                TempoLimite: a.config.TempoLimite,
        }
        clienteTemporario := zabbix.NovoClienteAPI(configAPI)
        err := clienteTemporario.TestarConexao()
        if err != nil {
                a.MostrarErro("Erro de conexão: " + err.Error())
                return
        }

        // Criar o perfil atualizado
        perfil := config.ConfiguracaoPerfil{
                Nome:  nome,
                URL:   url,
                Token: token,
        }

        // Atualizar o perfil
        err = a.config.AtualizarPerfil(indice, perfil)
        if err != nil {
                a.MostrarErro("Erro ao atualizar perfil: " + err.Error())
                return
        }

        // Atualizar a lista de perfis
        a.listaPerfis = a.config.Perfis

        // Salvar a configuração
        err = a.config.Salvar(a.arquivoConfig)
        if err != nil {
                a.MostrarErro("Erro ao salvar configuração: " + err.Error())
                return
        }

        // Se estiver editando o perfil ativo, atualizar o cliente API
        if indice == a.config.PerfilAtual {
                a.ConfigurarClienteAPI(&a.config.Perfis[indice])
        }

        a.MostrarSucesso("Perfil atualizado com sucesso")
}

// acaoRemoverPerfil processa a ação de remover um perfil
func (a *Aplicacao) acaoRemoverPerfil() {
        if len(a.listaPerfis) == 0 {
                a.MostrarErro("Não há perfis para remover")
                return
        }

        indice := a.indicePerfilAtual
        if indice < 0 || indice >= len(a.listaPerfis) {
                a.MostrarErro("Selecione um perfil válido")
                return
        }

        // Remover o perfil
        err := a.config.RemoverPerfil(indice)
        if err != nil {
                a.MostrarErro("Erro ao remover perfil: " + err.Error())
                return
        }

        // Atualizar a lista de perfis
        a.listaPerfis = a.config.Perfis
        a.indicePerfilAtual = a.config.PerfilAtual

        // Salvar a configuração
        err = a.config.Salvar(a.arquivoConfig)
        if err != nil {
                a.MostrarErro("Erro ao salvar configuração: " + err.Error())
                return
        }

        // Limpar campos
        a.txtNomePerfil.SetText("")
        a.txtURLAPI.SetText("")
        a.txtToken.SetText("")

        // Atualizar cliente API se necessário
        if len(a.config.Perfis) > 0 && a.config.PerfilAtual >= 0 {
                a.ConfigurarClienteAPI(&a.config.Perfis[a.config.PerfilAtual])
        } else {
                a.clienteAPI = nil
                a.perfilAtual = nil
        }

        a.MostrarSucesso("Perfil removido com sucesso")
}

// acaoSalvarConfiguracao processa a ação de salvar a configuração
func (a *Aplicacao) acaoSalvarConfiguracao() {
        // TODO: Implementar a leitura do campo de tempo limite
        //tempoLimite, _ := strconv.Atoi(a.txtTempoLimite.Text())
        //a.config.TempoLimite = tempoLimite

        // Salvar a configuração
        err := a.config.Salvar(a.arquivoConfig)
        if err != nil {
                a.MostrarErro("Erro ao salvar configuração: " + err.Error())
                return
        }

        // Atualizar cliente API se necessário
        if a.perfilAtual != nil {
                a.ConfigurarClienteAPI(a.perfilAtual)
        }

        a.MostrarSucesso("Configuração salva com sucesso")
}