package ui

import (
        "fmt"
        "strconv"
        "strings"

        "fyne.io/fyne/v2"
        "fyne.io/fyne/v2/container"
        "fyne.io/fyne/v2/layout"
        "fyne.io/fyne/v2/theme"
        "fyne.io/fyne/v2/widget"
        "zabbix-manager/zabbix"
)

// TelaPrincipal representa a tela principal do aplicativo
type TelaPrincipal struct {
        App                *AplicacaoZabbix
        Container          *fyne.Container
        TabelaHosts        *widget.Table
        CampoBusca         *widget.Entry
        BotaoAtualizar     *widget.Button
        BotaoExportar      *widget.Button
        BotaoSair          *widget.Button
        BotaoConfigurar    *widget.Button
        LabelStatus        *widget.Label
        BarraStatus        *fyne.Container
        Hosts              []zabbix.Host
        HostsFiltrados     []zabbix.Host
}

// MostrarTelaPrincipal exibe a tela principal
func (a *AplicacaoZabbix) MostrarTelaPrincipal() {
        tela := &TelaPrincipal{
                App: a,
        }

        tela.CriarInterface()
        a.Janela.SetContent(tela.Container)

        // Carregar dados iniciais
        tela.CarregarDados()
}

// CriarInterface cria a interface da tela principal
func (t *TelaPrincipal) CriarInterface() {
        // Criar widgets da tela
        t.CriarWidgets()

        // Barra de ferramentas
        barraFerramentas := container.NewHBox(
                t.CampoBusca,
                t.BotaoAtualizar,
                t.BotaoExportar,
                t.BotaoConfigurar,
                layout.NewSpacer(),
                t.BotaoSair,
        )

        // Barra de status
        t.BarraStatus = container.NewHBox(t.LabelStatus)

        // Container principal
        t.Container = container.NewBorder(barraFerramentas, t.BarraStatus, nil, nil, t.TabelaHosts)
}

// CriarWidgets cria os widgets da tela
func (t *TelaPrincipal) CriarWidgets() {
        // Campo de busca
        t.CampoBusca = widget.NewEntry()
        t.CampoBusca.SetPlaceHolder("Buscar hosts...")
        t.CampoBusca.OnChanged = t.AplicarFiltro

        // Botões
        t.BotaoAtualizar = widget.NewButtonWithIcon("Atualizar", theme.ViewRefreshIcon(), t.CarregarDados)
        t.BotaoExportar = widget.NewButtonWithIcon("Exportar CSV", theme.DocumentSaveIcon(), t.ExportarDados)
        t.BotaoConfigurar = widget.NewButtonWithIcon("Configurar", theme.SettingsIcon(), t.AbrirConfiguracoes)
        t.BotaoSair = widget.NewButtonWithIcon("Sair", theme.LogoutIcon(), t.Sair)

        // Label de status
        t.LabelStatus = widget.NewLabel("Pronto")

        // Tabela de hosts
        t.TabelaHosts = widget.NewTable(
                func() (int, int) {
                        return len(t.HostsFiltrados) + 1, 5 // +1 para o cabeçalho
                },
                func() fyne.CanvasObject {
                        return widget.NewLabel("Texto longo para prever o tamanho")
                },
                func(i widget.TableCellID, o fyne.CanvasObject) {
                        label := o.(*widget.Label)
                        label.TextStyle = fyne.TextStyle{Bold: i.Row == 0}
                        label.Alignment = fyne.TextAlignLeading

                        if i.Row == 0 {
                                // Cabeçalho
                                switch i.Col {
                                case 0:
                                        label.SetText("ID do Host")
                                case 1:
                                        label.SetText("Nome do Host")
                                case 2:
                                        label.SetText("Status")
                                case 3:
                                        label.SetText("Itens")
                                case 4:
                                        label.SetText("Triggers")
                                }
                        } else if i.Row-1 < len(t.HostsFiltrados) {
                                // Dados
                                host := t.HostsFiltrados[i.Row-1]
                                switch i.Col {
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
                                case 3:
                                        label.SetText(strconv.Itoa(len(host.Items)))
                                case 4:
                                        label.SetText(strconv.Itoa(len(host.Triggers)))
                                }
                        }
                },
        )

        // Configurar tamanho das colunas
        t.TabelaHosts.SetColumnWidth(0, 100)
        t.TabelaHosts.SetColumnWidth(1, 300)
        t.TabelaHosts.SetColumnWidth(2, 100)
        t.TabelaHosts.SetColumnWidth(3, 80)
        t.TabelaHosts.SetColumnWidth(4, 80)
}

// CarregarDados carrega os dados dos hosts
func (t *TelaPrincipal) CarregarDados() {
        // Verificar se há cliente configurado
        if t.App.Cliente == nil {
                t.App.MostrarErro("Erro", "Nenhum servidor Zabbix configurado")
                t.App.MostrarTelaLogin()
                return
        }

        // Atualizar label de status
        t.LabelStatus.SetText("Carregando dados...")
        t.LabelStatus.Refresh()

        // Obter os hosts
        var erro error
        t.Hosts, erro = t.App.Cliente.ObterHosts()
        if erro != nil {
                t.App.MostrarErro("Erro", fmt.Sprintf("Erro ao obter dados dos hosts: %v", erro))
                t.LabelStatus.SetText("Erro ao carregar dados")
                t.LabelStatus.Refresh()
                return
        }

        // Aplicar filtro
        t.AplicarFiltro(t.CampoBusca.Text)

        // Atualizar status
        t.LabelStatus.SetText(fmt.Sprintf("Dados carregados. Total de hosts: %d", len(t.Hosts)))
        t.LabelStatus.Refresh()
}

// AplicarFiltro aplica o filtro de busca aos hosts
func (t *TelaPrincipal) AplicarFiltro(termo string) {
        // Se o termo estiver vazio, mostrar todos os hosts
        if termo == "" {
                t.HostsFiltrados = t.Hosts
                t.TabelaHosts.Refresh()
                return
        }

        // Filtrar hosts pelo nome
        t.HostsFiltrados = []zabbix.Host{}
        for _, host := range t.Hosts {
                if ContémTexto(host.Nome, termo) || ContémTexto(host.ID, termo) {
                        t.HostsFiltrados = append(t.HostsFiltrados, host)
                }
        }

        // Atualizar a tabela
        t.TabelaHosts.Refresh()

        // Atualizar status
        t.LabelStatus.SetText(fmt.Sprintf("Mostrando %d de %d hosts", len(t.HostsFiltrados), len(t.Hosts)))
        t.LabelStatus.Refresh()
}

// ExportarDados exporta os dados para CSV
func (t *TelaPrincipal) ExportarDados() {
        // Verificar se há hosts carregados
        if len(t.Hosts) == 0 {
                t.App.MostrarErro("Erro", "Não há dados para exportar")
                return
        }

        // Exportar relatório
        erro := t.App.ExportarRelatorio()
        if erro != nil {
                t.App.MostrarErro("Erro", fmt.Sprintf("Erro ao exportar relatório: %v", erro))
                return
        }

        // Mostrar mensagem de sucesso
        diretorioHome, _ := t.App.ObterDiretorioHome()
        caminho := fmt.Sprintf("%s/Relatórios Zabbix", diretorioHome)
        t.App.MostrarInfo("Relatório Exportado", fmt.Sprintf("Relatório exportado com sucesso para o diretório:\n%s", caminho))
}

// AbrirConfiguracoes abre a tela de configurações
func (t *TelaPrincipal) AbrirConfiguracoes() {
        t.App.MostrarTelaConfig()
}

// Sair volta para a tela de login
func (t *TelaPrincipal) Sair() {
        t.App.MostrarTelaLogin()
}

// ContémTexto verifica se um texto contém outro (case insensitive)
func ContémTexto(texto, busca string) bool {
        if len(busca) == 0 {
                return true
        }
        if len(texto) == 0 {
                return false
        }
        return strings.Contains(
                strings.ToLower(texto),
                strings.ToLower(busca),
        )
}