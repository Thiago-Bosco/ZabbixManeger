package main

import (
        "fmt"
        "html/template"
        "log"
        "net/http"
        "os"
        "path/filepath"
        "strings"
        "time"

        "zabbix-manager/config"
        "zabbix-manager/zabbix"
)

// Estruturas para os modelos HTML
type PaginaLogin struct {
        Erro         string
        Sucesso      string
        ListaPerfis  []config.ConfiguracaoPerfil
        PerfilAtivo  int
        ModoEdicao   bool
        PerfilEditar *config.ConfiguracaoPerfil
}

type PaginaPrincipal struct {
        NomeServidor   string
        URLServidor    string
        Hosts          []zabbix.Host
        TermoBusca     string
        MensagemErro   string
        MensagemSucesso string
}

var (
        cfg            *config.Configuração
        arquivoConfig  string
        clienteAPI     *zabbix.ClienteAPI
        templatesCache map[string]*template.Template
)

func init() {
        // Configurar log
        log.SetFlags(log.LstdFlags | log.Lshortfile)
        
        // Inicializar cache de templates
        templatesCache = make(map[string]*template.Template)
        
        // Carregar templates
        templates := []string{"login", "principal", "config"}
        for _, nome := range templates {
                carregarTemplate(nome)
        }
}

func carregarTemplate(nome string) {
        caminhoBase := "templates"
        caminhoTemplate := filepath.Join(caminhoBase, nome+".html")
        caminhoLayout := filepath.Join(caminhoBase, "layout.html")
        
        // Verificar se os arquivos existem
        if _, err := os.Stat(caminhoTemplate); os.IsNotExist(err) {
                log.Printf("Template %s não encontrado: %v", nome, err)
                return
        }
        
        if _, err := os.Stat(caminhoLayout); os.IsNotExist(err) {
                // Tentar criar apenas com o template principal
                tmpl, err := template.New(nome).ParseFiles(caminhoTemplate)
                if err != nil {
                        log.Printf("Erro ao processar template %s: %v", nome, err)
                        return
                }
                templatesCache[nome] = tmpl
                return
        }
        
        // Carregar template com layout
        tmpl, err := template.New(nome).ParseFiles(caminhoLayout, caminhoTemplate)
        if err != nil {
                log.Printf("Erro ao processar template %s: %v", nome, err)
                return
        }
        
        templatesCache[nome] = tmpl
}

func renderizarTemplate(w http.ResponseWriter, nome string, dados interface{}) {
        tmpl, ok := templatesCache[nome]
        if !ok {
                log.Printf("Template %s não encontrado no cache", nome)
                http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
                return
        }
        
        err := tmpl.ExecuteTemplate(w, "layout", dados)
        if err != nil {
                log.Printf("Erro ao renderizar template %s: %v", nome, err)
                http.Error(w, "Erro ao renderizar página", http.StatusInternalServerError)
        }
}

func main() {
        // Carregar configuração
        arquivoConfig = config.ObterCaminhoConfiguracao()
        var err error
        cfg, err = config.Carregar(arquivoConfig)
        if err != nil {
                log.Printf("Erro ao carregar configuração: %v", err)
                // Criar configuração padrão
                cfg = config.NovaPadrao()
                if err := cfg.Salvar(arquivoConfig); err != nil {
                        log.Fatalf("Erro ao salvar configuração padrão: %v", err)
                }
        }
        
        // Inicializar API se houver perfil ativo
        inicializarClienteAPI()
        
        // Configurar rotas
        http.HandleFunc("/", manipuladorHome)
        http.HandleFunc("/login", manipuladorLogin)
        http.HandleFunc("/config", manipuladorConfig)
        http.HandleFunc("/perfil/adicionar", manipuladorAdicionarPerfil)
        http.HandleFunc("/perfil/editar", manipuladorEditarPerfil)
        http.HandleFunc("/perfil/remover", manipuladorRemoverPerfil)
        http.HandleFunc("/perfil/selecionar", manipuladorSelecionarPerfil)
        http.HandleFunc("/hosts", manipuladorHosts)
        http.HandleFunc("/hosts/buscar", manipuladorBuscarHosts)
        http.HandleFunc("/exportar", manipuladorExportarCSV)
        
        // Servir arquivos estáticos
        fs := http.FileServer(http.Dir("static"))
        http.Handle("/static/", http.StripPrefix("/static/", fs))
        
        // Iniciar servidor
        porta := "5000"
        log.Printf("Zabbix Manager Web iniciando na porta %s...", porta)
        addr := fmt.Sprintf("0.0.0.0:%s", porta)
        if err := http.ListenAndServe(addr, nil); err != nil {
                log.Fatalf("Erro ao iniciar servidor: %v", err)
        }
}

func inicializarClienteAPI() {
        // Verificar se há perfil configurado e ativo
        perfilAtivo, err := cfg.PerfilAtivo()
        if err != nil {
                log.Printf("Aviso: %v", err)
                clienteAPI = nil
                return
        }
        
        // Criar cliente API
        configAPI := zabbix.ConfigAPI{
                URL:         perfilAtivo.URL,
                Token:       perfilAtivo.Token,
                TempoLimite: cfg.TempoLimite,
        }
        clienteAPI = zabbix.NovoClienteAPI(configAPI)
        
        // Testar conexão
        if err := clienteAPI.TestarConexao(); err != nil {
                log.Printf("Erro ao testar conexão com o servidor Zabbix: %v", err)
                clienteAPI = nil
        }
}

func manipuladorHome(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/" {
                http.NotFound(w, r)
                return
        }
        
        // Redirecionar para a página de login se não houver perfil ativo
        if clienteAPI == nil {
                http.Redirect(w, r, "/login", http.StatusFound)
                return
        }
        
        // Redirecionar para a página de hosts
        http.Redirect(w, r, "/hosts", http.StatusFound)
}

func manipuladorLogin(w http.ResponseWriter, r *http.Request) {
        pagina := PaginaLogin{
                ListaPerfis: cfg.Perfis,
                PerfilAtivo: cfg.PerfilAtual,
        }
        
        renderizarTemplate(w, "login", pagina)
}

func manipuladorConfig(w http.ResponseWriter, r *http.Request) {
        pagina := PaginaLogin{
                ListaPerfis: cfg.Perfis,
                PerfilAtivo: cfg.PerfilAtual,
                ModoEdicao: false,
        }
        
        renderizarTemplate(w, "config", pagina)
}

func manipuladorAdicionarPerfil(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
                http.Redirect(w, r, "/config", http.StatusFound)
                return
        }
        
        // Obter dados do formulário
        err := r.ParseForm()
        if err != nil {
                http.Error(w, "Erro ao processar formulário", http.StatusBadRequest)
                return
        }
        
        nome := r.Form.Get("nome")
        url := r.Form.Get("url")
        token := r.Form.Get("token")
        
        // Validar dados
        if nome == "" || url == "" || token == "" {
                pagina := PaginaLogin{
                        ListaPerfis: cfg.Perfis,
                        PerfilAtivo: cfg.PerfilAtual,
                        Erro:        "Todos os campos são obrigatórios",
                }
                renderizarTemplate(w, "config", pagina)
                return
        }
        
        // Testar conexão
        configAPI := zabbix.ConfigAPI{
                URL:         url,
                Token:       token,
                TempoLimite: cfg.TempoLimite,
        }
        clienteTemporario := zabbix.NovoClienteAPI(configAPI)
        if err := clienteTemporario.TestarConexao(); err != nil {
                pagina := PaginaLogin{
                        ListaPerfis: cfg.Perfis,
                        PerfilAtivo: cfg.PerfilAtual,
                        Erro:        fmt.Sprintf("Erro ao conectar: %v", err),
                }
                renderizarTemplate(w, "config", pagina)
                return
        }
        
        // Adicionar perfil
        perfil := config.ConfiguracaoPerfil{
                Nome:  nome,
                URL:   url,
                Token: token,
        }
        cfg.AdicionarPerfil(perfil)
        
        // Salvar configuração
        if err := cfg.Salvar(arquivoConfig); err != nil {
                pagina := PaginaLogin{
                        ListaPerfis: cfg.Perfis,
                        PerfilAtivo: cfg.PerfilAtual,
                        Erro:        fmt.Sprintf("Erro ao salvar configuração: %v", err),
                }
                renderizarTemplate(w, "config", pagina)
                return
        }
        
        // Redirecionar para a página de configuração
        http.Redirect(w, r, "/config?sucesso=Perfil adicionado com sucesso", http.StatusFound)
}

func manipuladorEditarPerfil(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodGet {
                // Obter índice do perfil a ser editado
                indiceStr := r.URL.Query().Get("indice")
                if indiceStr == "" {
                        http.Redirect(w, r, "/config", http.StatusFound)
                        return
                }
                
                // Converter índice para inteiro
                indice := 0
                fmt.Sscanf(indiceStr, "%d", &indice)
                
                // Validar índice
                if indice < 0 || indice >= len(cfg.Perfis) {
                        http.Redirect(w, r, "/config", http.StatusFound)
                        return
                }
                
                // Renderizar página de configuração com perfil a ser editado
                pagina := PaginaLogin{
                        ListaPerfis:  cfg.Perfis,
                        PerfilAtivo:  cfg.PerfilAtual,
                        ModoEdicao:   true,
                        PerfilEditar: &cfg.Perfis[indice],
                }
                renderizarTemplate(w, "config", pagina)
                return
        }
        
        if r.Method == http.MethodPost {
                // Obter dados do formulário
                err := r.ParseForm()
                if err != nil {
                        http.Error(w, "Erro ao processar formulário", http.StatusBadRequest)
                        return
                }
                
                indiceStr := r.Form.Get("indice")
                nome := r.Form.Get("nome")
                url := r.Form.Get("url")
                token := r.Form.Get("token")
                
                // Validar dados
                indice := 0
                fmt.Sscanf(indiceStr, "%d", &indice)
                
                if indice < 0 || indice >= len(cfg.Perfis) || nome == "" || url == "" || token == "" {
                        http.Redirect(w, r, "/config?erro=Dados inválidos", http.StatusFound)
                        return
                }
                
                // Atualizar perfil
                cfg.Perfis[indice].Nome = nome
                cfg.Perfis[indice].URL = url
                cfg.Perfis[indice].Token = token
                
                // Salvar configuração
                if err := cfg.Salvar(arquivoConfig); err != nil {
                        http.Redirect(w, r, fmt.Sprintf("/config?erro=%s", err), http.StatusFound)
                        return
                }
                
                // Reinicializar cliente API se o perfil editado for o ativo
                if indice == cfg.PerfilAtual {
                        inicializarClienteAPI()
                }
                
                // Redirecionar para a página de configuração
                http.Redirect(w, r, "/config?sucesso=Perfil atualizado com sucesso", http.StatusFound)
                return
        }
        
        http.Redirect(w, r, "/config", http.StatusFound)
}

func manipuladorRemoverPerfil(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
                http.Redirect(w, r, "/config", http.StatusFound)
                return
        }
        
        // Obter índice do perfil a ser removido
        err := r.ParseForm()
        if err != nil {
                http.Error(w, "Erro ao processar formulário", http.StatusBadRequest)
                return
        }
        
        indiceStr := r.Form.Get("indice")
        if indiceStr == "" {
                http.Redirect(w, r, "/config", http.StatusFound)
                return
        }
        
        // Converter índice para inteiro
        indice := 0
        fmt.Sscanf(indiceStr, "%d", &indice)
        
        // Remover perfil
        err = cfg.RemoverPerfil(indice)
        if err != nil {
                http.Redirect(w, r, fmt.Sprintf("/config?erro=%s", err), http.StatusFound)
                return
        }
        
        // Salvar configuração
        if err := cfg.Salvar(arquivoConfig); err != nil {
                http.Redirect(w, r, fmt.Sprintf("/config?erro=%s", err), http.StatusFound)
                return
        }
        
        // Reinicializar cliente API
        inicializarClienteAPI()
        
        // Redirecionar para a página de configuração
        http.Redirect(w, r, "/config?sucesso=Perfil removido com sucesso", http.StatusFound)
}

func manipuladorSelecionarPerfil(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
                http.Redirect(w, r, "/config", http.StatusFound)
                return
        }
        
        // Obter índice do perfil a ser selecionado
        err := r.ParseForm()
        if err != nil {
                http.Error(w, "Erro ao processar formulário", http.StatusBadRequest)
                return
        }
        
        indiceStr := r.Form.Get("indice")
        if indiceStr == "" {
                http.Redirect(w, r, "/config", http.StatusFound)
                return
        }
        
        // Converter índice para inteiro
        indice := 0
        fmt.Sscanf(indiceStr, "%d", &indice)
        
        // Selecionar perfil
        err = cfg.SelecionarPerfil(indice)
        if err != nil {
                http.Redirect(w, r, fmt.Sprintf("/config?erro=%s", err), http.StatusFound)
                return
        }
        
        // Salvar configuração
        if err := cfg.Salvar(arquivoConfig); err != nil {
                http.Redirect(w, r, fmt.Sprintf("/config?erro=%s", err), http.StatusFound)
                return
        }
        
        // Reinicializar cliente API
        inicializarClienteAPI()
        
        // Redirecionar para a página inicial
        http.Redirect(w, r, "/?sucesso=Perfil selecionado com sucesso", http.StatusFound)
}

func manipuladorHosts(w http.ResponseWriter, r *http.Request) {
        // Verificar se há perfil ativo
        if clienteAPI == nil {
                http.Redirect(w, r, "/login", http.StatusFound)
                return
        }
        
        // Obter perfil ativo
        perfilAtivo, err := cfg.PerfilAtivo()
        if err != nil {
                http.Redirect(w, r, "/login", http.StatusFound)
                return
        }
        
        // Obter hosts
        hosts, err := clienteAPI.ObterHosts()
        if err != nil {
                pagina := PaginaPrincipal{
                        NomeServidor: perfilAtivo.Nome,
                        URLServidor:  perfilAtivo.URL,
                        Hosts:        []zabbix.Host{},
                        MensagemErro: fmt.Sprintf("Erro ao obter hosts: %v", err),
                }
                renderizarTemplate(w, "principal", pagina)
                return
        }
        
        // Renderizar página
        pagina := PaginaPrincipal{
                NomeServidor: perfilAtivo.Nome,
                URLServidor:  perfilAtivo.URL,
                Hosts:        hosts,
                MensagemSucesso: r.URL.Query().Get("sucesso"),
                MensagemErro: r.URL.Query().Get("erro"),
        }
        renderizarTemplate(w, "principal", pagina)
}

func manipuladorBuscarHosts(w http.ResponseWriter, r *http.Request) {
        // Verificar se há perfil ativo
        if clienteAPI == nil {
                http.Redirect(w, r, "/login", http.StatusFound)
                return
        }
        
        // Obter perfil ativo
        perfilAtivo, err := cfg.PerfilAtivo()
        if err != nil {
                http.Redirect(w, r, "/login", http.StatusFound)
                return
        }
        
        // Obter termo de busca
        termo := r.URL.Query().Get("termo")
        if termo == "" {
                http.Redirect(w, r, "/hosts", http.StatusFound)
                return
        }
        
        // Obter hosts
        hosts, err := clienteAPI.ObterHosts()
        if err != nil {
                pagina := PaginaPrincipal{
                        NomeServidor: perfilAtivo.Nome,
                        URLServidor:  perfilAtivo.URL,
                        TermoBusca:   termo,
                        Hosts:        []zabbix.Host{},
                        MensagemErro: fmt.Sprintf("Erro ao obter hosts: %v", err),
                }
                renderizarTemplate(w, "principal", pagina)
                return
        }
        
        // Filtrar hosts
        termoLower := strings.ToLower(termo)
        hostsFiltrados := []zabbix.Host{}
        for _, host := range hosts {
                if strings.Contains(strings.ToLower(host.Nome), termoLower) ||
                        strings.Contains(strings.ToLower(host.ID), termoLower) {
                        hostsFiltrados = append(hostsFiltrados, host)
                }
        }
        
        // Renderizar página
        pagina := PaginaPrincipal{
                NomeServidor: perfilAtivo.Nome,
                URLServidor:  perfilAtivo.URL,
                TermoBusca:   termo,
                Hosts:        hostsFiltrados,
        }
        renderizarTemplate(w, "principal", pagina)
}

func manipuladorExportarCSV(w http.ResponseWriter, r *http.Request) {
        // Verificar se há perfil ativo
        if clienteAPI == nil {
                http.Redirect(w, r, "/login", http.StatusFound)
                return
        }
        
        // Obter perfil ativo
        perfilAtivo, err := cfg.PerfilAtivo()
        if err != nil {
                http.Redirect(w, r, "/login", http.StatusFound)
                return
        }
        
        // Obter hosts
        hosts, err := clienteAPI.ObterHosts()
        if err != nil {
                http.Redirect(w, r, fmt.Sprintf("/hosts?erro=%s", err), http.StatusFound)
                return
        }
        
        // Configurar cabeçalhos para download
        nomeArquivo := fmt.Sprintf("relatorio_%s_%s.csv", 
                perfilAtivo.Nome, 
                time.Now().Format("2006-01-02_15-04-05"))
        w.Header().Set("Content-Type", "text/csv")
        w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", nomeArquivo))
        
        // Gerar CSV diretamente para o response writer
        err = zabbix.GerarRelatorioCSVStream(hosts, w)
        if err != nil {
                http.Redirect(w, r, fmt.Sprintf("/hosts?erro=%s", err), http.StatusFound)
                return
        }
}