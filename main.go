package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"zabbix-manager/config"
	"zabbix-manager/zabbix"
)

// HTML Template Structures
type PaginaLogin struct {
	Erro         string
	Sucesso      string
	ListaPerfis  []config.ConfiguracaoPerfil
	PerfilAtivo  int
	ModoEdicao   bool
	PerfilEditar *config.ConfiguracaoPerfil
}

type PaginaPrincipal struct {
	NomeServidor    string
	URLServidor     string
	Hosts           []zabbix.Host
	TermoBusca      string
	MensagemErro    string
	MensagemSucesso string
}

var (
	cfg            *config.Configuração
	arquivoConfig  string
	clienteAPI     *zabbix.ClienteAPI
	templatesCache map[string]*template.Template
	funcMap        template.FuncMap
)

func init() {
	// Configure logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Initialize template cache
	templatesCache = make(map[string]*template.Template)

	// Custom template functions
	funcMap = template.FuncMap{
		"subtract": func(a, b int) int {
			return a - b
		},
	}

	// Load templates
	templates := []string{"login", "principal", "config", "analise"}
	for _, nome := range templates {
		carregarTemplate(nome)
	}
}

func carregarTemplate(nome string) {
	basePath := "templates"
	templatePath := filepath.Join(basePath, nome+".html")
	layoutPath := filepath.Join(basePath, "layout.html")

	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		log.Printf("Template %s not found: %v", nome, err)
		return
	}

	if _, err := os.Stat(layoutPath); os.IsNotExist(err) {
		tmpl, err := template.New(nome).Funcs(funcMap).ParseFiles(templatePath)
		if err != nil {
			log.Printf("Error processing template %s: %v", nome, err)
			return
		}
		templatesCache[nome] = tmpl
		return
	}

	tmpl, err := template.New(nome).Funcs(funcMap).ParseFiles(layoutPath, templatePath)
	if err != nil {
		log.Printf("Error processing template %s: %v", nome, err)
		return
	}

	templatesCache[nome] = tmpl
}

func renderizarTemplate(w http.ResponseWriter, nome string, dados interface{}) {
	tmpl, ok := templatesCache[nome]
	if !ok {
		log.Printf("Template %s not found in cache", nome)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err := tmpl.ExecuteTemplate(w, "layout", dados)
	if err != nil {
		log.Printf("Error rendering template %s: %v", nome, err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}

func inicializarClienteAPI() {
	perfilAtivo, err := cfg.PerfilAtivo()
	if err != nil {
		log.Printf("Warning: %v", err)
		clienteAPI = nil
		return
	}

	configAPI := zabbix.ConfigAPI{
		URL:         perfilAtivo.URL,
		Token:       perfilAtivo.Token,
		TempoLimite: cfg.TempoLimite,
	}
	clienteAPI = zabbix.NovoClienteAPI(configAPI)

	// Optional: Add monthly analysis
	ano := time.Now().Year()
	mes := int(time.Now().Month())
	log.Printf("Loading analysis for %d/%d", ano, mes)
	if analises, err := clienteAPI.AnalisarProblemasMensais(ano, mes); err != nil {
		log.Printf("Error analyzing problems: %v", err)
	} else {
		log.Printf("Found %d analysis records", len(analises))
	}

	if err := clienteAPI.TestarConexao(); err != nil {
		log.Printf("Error testing Zabbix server connection: %v", err)
		clienteAPI = nil
	}
}

// Handler Functions
func manipuladorHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if clienteAPI == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

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
		ModoEdicao:  false,
	}
	renderizarTemplate(w, "config", pagina)
}

func manipuladorAdicionarPerfil(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/config", http.StatusFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Erro ao processar formulário", http.StatusBadRequest)
		return
	}

	nome := r.Form.Get("nome")
	url := r.Form.Get("url")
	token := r.Form.Get("token")

	if nome == "" || url == "" || token == "" {
		pagina := PaginaLogin{
			ListaPerfis: cfg.Perfis,
			PerfilAtivo: cfg.PerfilAtual,
			Erro:        "Todos os campos são obrigatórios",
		}
		renderizarTemplate(w, "config", pagina)
		return
	}

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

	perfil := config.ConfiguracaoPerfil{
		Nome:  nome,
		URL:   url,
		Token: token,
	}
	cfg.AdicionarPerfil(perfil)

	if err := cfg.Salvar(arquivoConfig); err != nil {
		pagina := PaginaLogin{
			ListaPerfis: cfg.Perfis,
			PerfilAtivo: cfg.PerfilAtual,
			Erro:        fmt.Sprintf("Erro ao salvar configuração: %v", err),
		}
		renderizarTemplate(w, "config", pagina)
		return
	}

	http.Redirect(w, r, "/config?sucesso=Perfil adicionado com sucesso", http.StatusFound)
}

func manipuladorEditarPerfil(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		indiceStr := r.URL.Query().Get("indice")
		if indiceStr == "" {
			http.Redirect(w, r, "/config", http.StatusFound)
			return
		}

		var indice int
		fmt.Sscanf(indiceStr, "%d", &indice)

		if indice < 0 || indice >= len(cfg.Perfis) {
			http.Redirect(w, r, "/config", http.StatusFound)
			return
		}

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
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Erro ao processar formulário", http.StatusBadRequest)
			return
		}

		indiceStr := r.Form.Get("indice")
		nome := r.Form.Get("nome")
		url := r.Form.Get("url")
		token := r.Form.Get("token")

		var indice int
		fmt.Sscanf(indiceStr, "%d", &indice)

		if indice < 0 || indice >= len(cfg.Perfis) || nome == "" || url == "" || token == "" {
			http.Redirect(w, r, "/config?erro=Dados inválidos", http.StatusFound)
			return
		}

		cfg.Perfis[indice].Nome = nome
		cfg.Perfis[indice].URL = url
		cfg.Perfis[indice].Token = token

		if err := cfg.Salvar(arquivoConfig); err != nil {
			http.Redirect(w, r, fmt.Sprintf("/config?erro=%s", err), http.StatusFound)
			return
		}

		if indice == cfg.PerfilAtual {
			inicializarClienteAPI()
		}

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

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Erro ao processar formulário", http.StatusBadRequest)
		return
	}

	indiceStr := r.Form.Get("indice")
	if indiceStr == "" {
		http.Redirect(w, r, "/config", http.StatusFound)
		return
	}

	var indice int
	fmt.Sscanf(indiceStr, "%d", &indice)

	if err := cfg.RemoverPerfil(indice); err != nil {
		http.Redirect(w, r, fmt.Sprintf("/config?erro=%s", err), http.StatusFound)
		return
	}

	if err := cfg.Salvar(arquivoConfig); err != nil {
		http.Redirect(w, r, fmt.Sprintf("/config?erro=%s", err), http.StatusFound)
		return
	}

	inicializarClienteAPI()
	http.Redirect(w, r, "/config?sucesso=Perfil removido com sucesso", http.StatusFound)
}

func manipuladorSelecionarPerfil(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/config", http.StatusFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Erro ao processar formulário", http.StatusBadRequest)
		return
	}

	indiceStr := r.Form.Get("indice")
	if indiceStr == "" {
		http.Redirect(w, r, "/config", http.StatusFound)
		return
	}

	var indice int
	fmt.Sscanf(indiceStr, "%d", &indice)

	if err := cfg.SelecionarPerfil(indice); err != nil {
		http.Redirect(w, r, fmt.Sprintf("/config?erro=%s", err), http.StatusFound)
		return
	}

	if err := cfg.Salvar(arquivoConfig); err != nil {
		http.Redirect(w, r, fmt.Sprintf("/config?erro=%s", err), http.StatusFound)
		return
	}

	inicializarClienteAPI()
	http.Redirect(w, r, "/?sucesso=Perfil selecionado com sucesso", http.StatusFound)
}

func manipuladorHosts(w http.ResponseWriter, r *http.Request) {
	if clienteAPI == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	perfilAtivo, err := cfg.PerfilAtivo()
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

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

	pagina := PaginaPrincipal{
		NomeServidor:    perfilAtivo.Nome,
		URLServidor:     perfilAtivo.URL,
		Hosts:           hosts,
		MensagemSucesso: r.URL.Query().Get("sucesso"),
		MensagemErro:    r.URL.Query().Get("erro"),
	}
	renderizarTemplate(w, "principal", pagina)
}

func manipuladorBuscarHosts(w http.ResponseWriter, r *http.Request) {
	if clienteAPI == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	perfilAtivo, err := cfg.PerfilAtivo()
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	termo := r.URL.Query().Get("termo")
	if termo == "" {
		http.Redirect(w, r, "/hosts", http.StatusFound)
		return
	}

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

	var hostsFiltrados []zabbix.Host
	if termo != "" {
		termoLower := strings.ToLower(termo)
		for _, host := range hosts {
			nomeLower := strings.ToLower(host.Nome)
			if strings.Contains(nomeLower, termoLower) {
				hostsFiltrados = append(hostsFiltrados, host)
			}
		}
	} else {
		hostsFiltrados = hosts
	}

	pagina := PaginaPrincipal{
		NomeServidor: perfilAtivo.Nome,
		URLServidor:  perfilAtivo.URL,
		TermoBusca:   termo,
		Hosts:        hostsFiltrados,
	}
	renderizarTemplate(w, "principal", pagina)
}

func manipuladorExportarCSV(w http.ResponseWriter, r *http.Request) {
	if clienteAPI == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	perfilAtivo, err := cfg.PerfilAtivo()
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	hosts, err := clienteAPI.ObterHosts()
	if err != nil {
		http.Redirect(w, r, fmt.Sprintf("/hosts?erro=%s", err), http.StatusFound)
		return
	}

	nomeArquivo := fmt.Sprintf("relatorio_%s_%s.csv",
		perfilAtivo.Nome,
		time.Now().Format("2006-01-02_15-04-05"))
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", nomeArquivo))

	if err := zabbix.GerarRelatorioCSVStream(hosts, w); err != nil {
		http.Redirect(w, r, fmt.Sprintf("/hosts?erro=%s", err), http.StatusFound)
		return
	}
}

func main() {
	// Load configuration
	arquivoConfig = config.ObterCaminhoConfiguracao()
	var err error
	cfg, err = config.Carregar(arquivoConfig)
	if err != nil {
		log.Printf("Error loading configuration: %v", err)
		cfg = config.NovaPadrao()
		if err := cfg.Salvar(arquivoConfig); err != nil {
			log.Fatalf("Error saving default configuration: %v", err)
		}
	}

	// Initialize API if active profile exists
	inicializarClienteAPI()

	// Configure routes
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
	http.HandleFunc("/analise", manipuladorAnalise)

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Start server
	porta := "5000"
	log.Printf("Zabbix Manager Web starting on port %s...", porta)
	addr := fmt.Sprintf("0.0.0.0:%s", porta)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
func manipuladorAnalise(w http.ResponseWriter, r *http.Request) {
	if clienteAPI == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	ano := time.Now().Year()
	mes := int(time.Now().Month())
	tipoFiltro := r.URL.Query().Get("tipo_filtro")
	if tipoFiltro == "" {
		tipoFiltro = "mensal"
	}

	var analises []zabbix.AnaliseProblema
	var err error

	if tipoFiltro == "mensal" {
		if anoStr := r.URL.Query().Get("ano"); anoStr != "" {
			if anoInt, err := strconv.Atoi(anoStr); err == nil {
				ano = anoInt
			}
		}
		if mesStr := r.URL.Query().Get("mes"); mesStr != "" {
			if mesInt, err := strconv.Atoi(mesStr); err == nil {
				mes = mesInt
			}
		}
		
		log.Printf("Analisando problemas para %d/%d", mes, ano)
		analisesMensais, err := clienteAPI.AnalisarProblemasMensais(ano, mes)
		if err != nil {
			log.Printf("Erro ao obter análises mensais: %v", err)
			return
		}
		
		// Converter AnaliseMensal para AnaliseProblema
		analises = make([]zabbix.AnaliseProblema, len(analisesMensais))
		for i, am := range analisesMensais {
			ap := zabbix.AnaliseProblema{
				HostID:              am.HostID,
				HostNome:            am.HostNome,
				TotalProblemas:      am.TotalProblemas,
				LimitesExcedidos:    am.LimitesExcedidos,
				ProblemasPorTrigger: am.ProblemasPorTrigger,
				PicoTrigger: struct {
					Nome      string
					DataPico  time.Time
					Contagem  int
					Gravidade string
				}{
					Nome:      am.PicoTrigger.Nome,
					DataPico:  am.PicoTrigger.DataPico,
					Contagem:  am.PicoTrigger.Contagem,
					Gravidade: am.PicoTrigger.Gravidade,
				},
			}
			analises[i] = ap
			log.Printf("Análise convertida para host %s: %d problemas, %d limites excedidos", 
				ap.HostNome, ap.TotalProblemas, ap.LimitesExcedidos)
		}
	} else {
		dataInicial := r.URL.Query().Get("data_inicial")
		dataFinal := r.URL.Query().Get("data_final")
		
		if dataInicial != "" && dataFinal != "" {
			log.Printf("Analisando problemas entre %s e %s", dataInicial, dataFinal)
			// Implemente a análise por período específico aqui
		}
	}

	if err != nil {
		log.Printf("Erro ao analisar problemas: %v", err)
		renderizarTemplate(w, "analise", map[string]interface{}{
			"Erro": fmt.Sprintf("Erro ao analisar problemas: %v", err),
			"TipoFiltro": tipoFiltro,
		})
		return
	}

	log.Printf("Encontrados %d registros de análise", len(analises))

	// Preparar dados para o template
	dados := map[string]interface{}{
		"Analises":       analises,
		"AnoSelecionado": ano,
		"MesSelecionado": mes,
		"Anos":          []int{ano - 1, ano, ano + 1},
		"Meses":         []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
		"NomesMeses":    []string{"Janeiro", "Fevereiro", "Março", "Abril", "Maio", "Junho", "Julho", "Agosto", "Setembro", "Outubro", "Novembro", "Dezembro"},
		"TipoFiltro":    "mensal",
	}

	renderizarTemplate(w, "analise", dados)
}
