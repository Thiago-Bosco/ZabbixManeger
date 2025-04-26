// +build headless

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"zabbix-manager/config"
	"zabbix-manager/zabbix"
)

func main() {
	// Configurar diretório de logs
	configurarLogs()

	// Iniciar aplicação no modo console
	log.Println("Iniciando Zabbix Manager (Modo Console)...")
	executarModoConsole()
}

func configurarLogs() {
	// Obter diretório home do usuário
	diretorioHome, err := os.UserHomeDir()
	if err != nil {
		// Usar diretório atual se não conseguir obter o home
		diretorioHome, _ = os.Getwd()
	}

	// Criar diretório de logs
	diretorioLogs := filepath.Join(diretorioHome, ".zabbix-manager", "logs")
	err = os.MkdirAll(diretorioLogs, 0755)
	if err != nil {
		log.Printf("Erro ao criar diretório de logs: %v", err)
		return
	}

	// Criar arquivo de log
	arquivoLog, err := os.OpenFile(
		filepath.Join(diretorioLogs, "zabbix-manager.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		log.Printf("Erro ao criar arquivo de log: %v", err)
		return
	}

	// Configurar log para escrever no arquivo e também na saída padrão
	log.SetOutput(os.Stdout)
}

func executarModoConsole() {
	// Carregar configuração
	arquivoConfig := config.ObterCaminhoConfiguracao()
	cfg, err := config.Carregar(arquivoConfig)
	if err != nil {
		log.Fatalf("Erro ao carregar configuração: %v", err)
	}

	// Verificar se há perfis configurados
	if len(cfg.Perfis) == 0 {
		fmt.Println("Não há perfis de servidor Zabbix configurados.")
		adicionarNovoPerfil(cfg, arquivoConfig)
	}

	// Se não há perfil ativo, selecionar o primeiro
	if cfg.PerfilAtual < 0 || cfg.PerfilAtual >= len(cfg.Perfis) {
		if len(cfg.Perfis) > 0 {
			cfg.PerfilAtual = 0
			cfg.Salvar(arquivoConfig)
		}
	}

	// Obter perfil ativo
	perfilAtivo, err := cfg.PerfilAtivo()
	if err != nil {
		log.Fatalf("Erro ao obter perfil ativo: %v", err)
	}

	// Criar cliente API
	configAPI := zabbix.ConfigAPI{
		URL:         perfilAtivo.URL,
		Token:       perfilAtivo.Token,
		TempoLimite: cfg.TempoLimite,
	}
	clienteAPI := zabbix.NovoClienteAPI(configAPI)

	// Testar conexão
	fmt.Printf("Conectando ao servidor Zabbix %s (%s)...\n", perfilAtivo.Nome, perfilAtivo.URL)
	if err := clienteAPI.TestarConexao(); err != nil {
		fmt.Printf("Erro ao conectar: %v\n", err)
		return
	}
	fmt.Println("Conexão estabelecida com sucesso!")

	// Menu principal
	exibir := true
	for exibir {
		exibir = exibirMenuPrincipal(clienteAPI, cfg, arquivoConfig, perfilAtivo)
	}
}

func exibirMenuPrincipal(clienteAPI *zabbix.ClienteAPI, cfg *config.Configuração, arquivoConfig string, perfilAtivo *config.ConfiguracaoPerfil) bool {
	fmt.Println("\n--- MENU PRINCIPAL ---")
	fmt.Println("1. Listar hosts")
	fmt.Println("2. Buscar hosts")
	fmt.Println("3. Exportar relatório CSV")
	fmt.Println("4. Gerenciar perfis")
	fmt.Println("5. Sair")
	fmt.Print("Escolha uma opção: ")

	opcao := lerEntrada()
	switch opcao {
	case "1":
		listarHosts(clienteAPI)
	case "2":
		buscarHosts(clienteAPI)
	case "3":
		exportarRelatorio(clienteAPI, perfilAtivo.Nome)
	case "4":
		gerenciarPerfis(cfg, arquivoConfig)
	case "5":
		return false
	default:
		fmt.Println("Opção inválida")
	}
	return true
}

func lerEntrada() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func listarHosts(clienteAPI *zabbix.ClienteAPI) {
	fmt.Println("\nCarregando hosts...")
	hosts, err := clienteAPI.ObterHosts()
	if err != nil {
		fmt.Printf("Erro ao obter hosts: %v\n", err)
		return
	}

	fmt.Printf("Total de hosts: %d\n\n", len(hosts))
	fmt.Println("ID\t\t\tNome\t\t\tStatus\t\tItens\tTriggers")
	fmt.Println("-------------------------------------------------------------------------------------")

	for _, host := range hosts {
		status := zabbix.StatusHost[host.Status]
		if status == "" {
			status = "Desconhecido"
		}
		fmt.Printf("%s\t%s\t\t%s\t\t%d\t%d\n", host.ID, host.Nome, status, len(host.Items), len(host.Triggers))
	}
}

func buscarHosts(clienteAPI *zabbix.ClienteAPI) {
	fmt.Print("\nDigite o termo de busca: ")
	termo := lerEntrada()

	fmt.Println("\nCarregando hosts...")
	hosts, err := clienteAPI.ObterHosts()
	if err != nil {
		fmt.Printf("Erro ao obter hosts: %v\n", err)
		return
	}

	termo = strings.ToLower(termo)
	hostsFiltrados := []zabbix.Host{}
	for _, host := range hosts {
		if strings.Contains(strings.ToLower(host.Nome), termo) ||
			strings.Contains(strings.ToLower(host.ID), termo) {
			hostsFiltrados = append(hostsFiltrados, host)
		}
	}

	fmt.Printf("Resultado da busca - Total: %d\n\n", len(hostsFiltrados))
	fmt.Println("ID\t\t\tNome\t\t\tStatus\t\tItens\tTriggers")
	fmt.Println("-------------------------------------------------------------------------------------")

	for _, host := range hostsFiltrados {
		status := zabbix.StatusHost[host.Status]
		if status == "" {
			status = "Desconhecido"
		}
		fmt.Printf("%s\t%s\t\t%s\t\t%d\t%d\n", host.ID, host.Nome, status, len(host.Items), len(host.Triggers))
	}
}

func exportarRelatorio(clienteAPI *zabbix.ClienteAPI, nomeServidor string) {
	fmt.Println("\nCarregando hosts para o relatório...")
	hosts, err := clienteAPI.ObterHosts()
	if err != nil {
		fmt.Printf("Erro ao obter hosts: %v\n", err)
		return
	}

	fmt.Printf("Total de hosts: %d\n", len(hosts))

	// Obter diretório home do usuário
	diretorioHome, err := os.UserHomeDir()
	if err != nil {
		diretorioHome, _ = os.Getwd()
	}

	// Criar diretório para relatórios
	diretorioRelatorios := filepath.Join(diretorioHome, "Relatórios Zabbix")
	err = os.MkdirAll(diretorioRelatorios, 0755)
	if err != nil {
		fmt.Printf("Erro ao criar diretório para relatórios: %v\n", err)
		return
	}

	// Caminho para o arquivo CSV
	caminhoArquivo := filepath.Join(diretorioRelatorios, fmt.Sprintf("relatorio_%s.csv", nomeServidor))

	// Exportar para CSV
	err = zabbix.GerarRelatorioCSV(hosts, caminhoArquivo)
	if err != nil {
		fmt.Printf("Erro ao gerar relatório: %v\n", err)
		return
	}

	fmt.Printf("Relatório exportado com sucesso para: %s\n", caminhoArquivo)
}

func gerenciarPerfis(cfg *config.Configuração, arquivoConfig string) {
	fmt.Println("\n--- GERENCIAR PERFIS ---")
	fmt.Println("1. Listar perfis")
	fmt.Println("2. Adicionar perfil")
	fmt.Println("3. Selecionar perfil")
	fmt.Println("4. Remover perfil")
	fmt.Println("5. Voltar")
	fmt.Print("Escolha uma opção: ")

	opcao := lerEntrada()
	switch opcao {
	case "1":
		listarPerfis(cfg)
	case "2":
		adicionarNovoPerfil(cfg, arquivoConfig)
	case "3":
		selecionarPerfil(cfg, arquivoConfig)
	case "4":
		removerPerfil(cfg, arquivoConfig)
	case "5":
		return
	default:
		fmt.Println("Opção inválida")
	}
}

func listarPerfis(cfg *config.Configuração) {
	fmt.Println("\n--- LISTA DE PERFIS ---")
	if len(cfg.Perfis) == 0 {
		fmt.Println("Não há perfis cadastrados")
		return
	}

	for i, perfil := range cfg.Perfis {
		ativo := ""
		if i == cfg.PerfilAtual {
			ativo = " (Ativo)"
		}
		fmt.Printf("%d. %s - %s%s\n", i+1, perfil.Nome, perfil.URL, ativo)
	}
}

func adicionarNovoPerfil(cfg *config.Configuração, arquivoConfig string) {
	fmt.Println("\n--- ADICIONAR PERFIL DE SERVIDOR ---")
	fmt.Print("Nome para identificar o servidor: ")
	nome := lerEntrada()

	fmt.Print("URL da API (ex: http://zabbix.example.com/api_jsonrpc.php): ")
	url := lerEntrada()

	fmt.Print("Token da API: ")
	token := lerEntrada()

	// Validar dados
	if nome == "" || url == "" || token == "" {
		fmt.Println("Todos os campos são obrigatórios")
		return
	}

	// Testar conexão
	configAPI := zabbix.ConfigAPI{
		URL:         url,
		Token:       token,
		TempoLimite: cfg.TempoLimite,
	}
	clienteTemporario := zabbix.NovoClienteAPI(configAPI)
	fmt.Println("Testando conexão...")
	if err := clienteTemporario.TestarConexao(); err != nil {
		fmt.Printf("Erro ao conectar: %v\n", err)
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
		fmt.Printf("Erro ao salvar configuração: %v\n", err)
		return
	}

	fmt.Println("Perfil adicionado com sucesso!")
}

func selecionarPerfil(cfg *config.Configuração, arquivoConfig string) {
	listarPerfis(cfg)
	if len(cfg.Perfis) == 0 {
		return
	}

	fmt.Print("\nDigite o número do perfil a ser selecionado: ")
	indiceStr := lerEntrada()
	indice, err := strconv.Atoi(indiceStr)
	if err != nil || indice < 1 || indice > len(cfg.Perfis) {
		fmt.Println("Número de perfil inválido")
		return
	}

	// Selecionar perfil
	err = cfg.SelecionarPerfil(indice - 1)
	if err != nil {
		fmt.Printf("Erro ao selecionar perfil: %v\n", err)
		return
	}

	// Salvar configuração
	if err := cfg.Salvar(arquivoConfig); err != nil {
		fmt.Printf("Erro ao salvar configuração: %v\n", err)
		return
	}

	fmt.Printf("Perfil '%s' selecionado com sucesso!\n", cfg.Perfis[indice-1].Nome)
	fmt.Println("Reinicie a aplicação para aplicar as alterações.")
}

func removerPerfil(cfg *config.Configuração, arquivoConfig string) {
	listarPerfis(cfg)
	if len(cfg.Perfis) == 0 {
		return
	}

	fmt.Print("\nDigite o número do perfil a ser removido: ")
	indiceStr := lerEntrada()
	indice, err := strconv.Atoi(indiceStr)
	if err != nil || indice < 1 || indice > len(cfg.Perfis) {
		fmt.Println("Número de perfil inválido")
		return
	}

	// Pedir confirmação
	nomePerfil := cfg.Perfis[indice-1].Nome
	fmt.Printf("Tem certeza que deseja remover o perfil '%s'? (s/n): ", nomePerfil)
	confirmacao := lerEntrada()
	if strings.ToLower(confirmacao) != "s" {
		fmt.Println("Operação cancelada")
		return
	}

	// Remover perfil
	err = cfg.RemoverPerfil(indice - 1)
	if err != nil {
		fmt.Printf("Erro ao remover perfil: %v\n", err)
		return
	}

	// Salvar configuração
	if err := cfg.Salvar(arquivoConfig); err != nil {
		fmt.Printf("Erro ao salvar configuração: %v\n", err)
		return
	}

	fmt.Printf("Perfil '%s' removido com sucesso!\n", nomePerfil)
	fmt.Println("Reinicie a aplicação para aplicar as alterações.")
}