// +build headless

package main

import (
        "bufio"
        "fmt"
        "log"
        "os"
        "strings"

        "zabbix-manager/config"
        "zabbix-manager/zabbix"
)

func main() {
        // Configurar diretório de logs
        configurarLogs()

        // Iniciar a versão console
        log.Println("Iniciando Zabbix Manager (Modo Console)...")
        
        // Carregar configuração
        arquivoConfig := config.ObterCaminhoConfiguracao()
        configuracao, err := config.Carregar(arquivoConfig)
        if err != nil {
                log.Printf("Erro ao carregar configuração: %v", err)
                configuracao = config.NovaConfiguração()
        }

        // Verificar se há perfis
        if len(configuracao.Perfis) == 0 {
                fmt.Println("Não há perfis de servidor Zabbix configurados.")
                adicionarPerfil(configuracao)
        }

        // Loop principal
        scanner := bufio.NewScanner(os.Stdin)
        for {
                exibirMenu()
                fmt.Print("Escolha uma opção: ")
                scanner.Scan()
                opcao := scanner.Text()

                switch opcao {
                case "1":
                        listarPerfis(configuracao)
                case "2":
                        adicionarPerfil(configuracao)
                case "3":
                        removerPerfil(configuracao, scanner)
                case "4":
                        listarHosts(configuracao, scanner)
                case "5":
                        exportarRelatorio(configuracao, scanner)
                case "0":
                        fmt.Println("Saindo do programa.")
                        return
                default:
                        fmt.Println("Opção inválida. Por favor, escolha novamente.")
                }

                // Salvar configuração após cada operação
                configuracao.Salvar(arquivoConfig)
        }
}

func exibirMenu() {
        fmt.Println("\n===== ZABBIX MANAGER (MODO CONSOLE) =====")
        fmt.Println("1. Listar perfis de servidores")
        fmt.Println("2. Adicionar perfil de servidor")
        fmt.Println("3. Remover perfil de servidor")
        fmt.Println("4. Listar hosts de um servidor")
        fmt.Println("5. Exportar relatório de hosts")
        fmt.Println("0. Sair")
        fmt.Println("========================================")
}

func listarPerfis(cfg *config.Configuração) {
        fmt.Println("\n--- PERFIS DE SERVIDORES ---")
        if len(cfg.Perfis) == 0 {
                fmt.Println("Não há perfis cadastrados.")
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

func adicionarPerfil(cfg *config.Configuração) {
        scanner := bufio.NewScanner(os.Stdin)
        
        fmt.Println("\n--- ADICIONAR PERFIL DE SERVIDOR ---")
        
        fmt.Print("Nome para identificar o servidor: ")
        scanner.Scan()
        nome := scanner.Text()
        
        fmt.Print("URL da API (ex: http://zabbix.example.com/api_jsonrpc.php): ")
        scanner.Scan()
        url := scanner.Text()
        
        fmt.Print("Token de API: ")
        scanner.Scan()
        token := scanner.Text()
        
        // Validar campos
        if nome == "" || url == "" || token == "" {
                fmt.Println("Erro: Todos os campos são obrigatórios.")
                return
        }
        
        // Testar conexão
        configAPI := zabbix.ConfigAPI{
                URL:         url,
                Token:       token,
                TempoLimite: cfg.TempoLimite,
        }
        cliente := zabbix.NovoClienteAPI(configAPI)
        err := cliente.TestarConexao()
        if err != nil {
                fmt.Printf("Erro ao conectar ao servidor: %v\n", err)
                return
        }
        
        // Adicionar perfil
        perfil := config.ConfiguracaoPerfil{
                Nome:  nome,
                URL:   url,
                Token: token,
        }
        
        cfg.AdicionarPerfil(perfil)
        fmt.Println("Perfil adicionado com sucesso!")
}

func removerPerfil(cfg *config.Configuração, scanner *bufio.Scanner) {
        if len(cfg.Perfis) == 0 {
                fmt.Println("Não há perfis para remover.")
                return
        }
        
        listarPerfis(cfg)
        
        fmt.Print("Digite o número do perfil que deseja remover: ")
        scanner.Scan()
        numStr := scanner.Text()
        
        // Converter para número
        var num int
        _, err := fmt.Sscanf(numStr, "%d", &num)
        if err != nil || num < 1 || num > len(cfg.Perfis) {
                fmt.Println("Número de perfil inválido.")
                return
        }
        
        // Índice real (0-based)
        indice := num - 1
        
        // Confirmar
        fmt.Printf("Tem certeza que deseja remover o perfil '%s'? (s/n): ", cfg.Perfis[indice].Nome)
        scanner.Scan()
        confirmacao := strings.ToLower(scanner.Text())
        
        if confirmacao == "s" || confirmacao == "sim" {
                err := cfg.RemoverPerfil(indice)
                if err != nil {
                        fmt.Printf("Erro ao remover perfil: %v\n", err)
                        return
                }
                fmt.Println("Perfil removido com sucesso!")
        } else {
                fmt.Println("Operação cancelada.")
        }
}

func listarHosts(cfg *config.Configuração, scanner *bufio.Scanner) {
        if len(cfg.Perfis) == 0 {
                fmt.Println("Não há perfis configurados.")
                return
        }
        
        // Selecionar perfil
        listarPerfis(cfg)
        
        fmt.Print("Digite o número do perfil para listar hosts: ")
        scanner.Scan()
        numStr := scanner.Text()
        
        // Converter para número
        var num int
        _, err := fmt.Sscanf(numStr, "%d", &num)
        if err != nil || num < 1 || num > len(cfg.Perfis) {
                fmt.Println("Número de perfil inválido.")
                return
        }
        
        // Índice real (0-based)
        indice := num - 1
        
        // Configurar cliente
        perfil := cfg.Perfis[indice]
        configAPI := zabbix.ConfigAPI{
                URL:         perfil.URL,
                Token:       perfil.Token,
                TempoLimite: cfg.TempoLimite,
        }
        cliente := zabbix.NovoClienteAPI(configAPI)
        
        // Buscar hosts
        fmt.Println("Buscando hosts...")
        hosts, err := cliente.ObterHosts()
        if err != nil {
                fmt.Printf("Erro ao buscar hosts: %v\n", err)
                return
        }
        
        // Exibir hosts
        fmt.Printf("\n--- HOSTS DO SERVIDOR %s ---\n", perfil.Nome)
        fmt.Printf("Total de hosts: %d\n\n", len(hosts))
        
        fmt.Println("ID | NOME | STATUS | ITENS | TRIGGERS")
        fmt.Println("-------------------------------------------")
        for _, host := range hosts {
                status := zabbix.StatusHost[host.Status]
                if status == "" {
                        status = "Desconhecido"
                }
                fmt.Printf("%s | %s | %s | %d | %d\n", 
                        host.ID, host.Nome, status, len(host.Items), len(host.Triggers))
        }
}

func configurarLogs() {
        // Obter diretório home do usuário
        diretorioHome, err := os.UserHomeDir()
        if err != nil {
                // Usar diretório atual se não conseguir obter o home
                diretorioHome, _ = os.Getwd()
        }

        // Criar diretório de logs
        diretorioLogs := fmt.Sprintf("%s/.zabbix-manager/logs", diretorioHome)
        err = os.MkdirAll(diretorioLogs, 0755)
        if err != nil {
                log.Printf("Erro ao criar diretório de logs: %v", err)
                return
        }

        // Configurar log para escrever na saída padrão
        log.SetOutput(os.Stdout)
}

func exportarRelatorio(cfg *config.Configuração, scanner *bufio.Scanner) {
        if len(cfg.Perfis) == 0 {
                fmt.Println("Não há perfis configurados.")
                return
        }
        
        // Selecionar perfil
        listarPerfis(cfg)
        
        fmt.Print("Digite o número do perfil para exportar relatório: ")
        scanner.Scan()
        numStr := scanner.Text()
        
        // Converter para número
        var num int
        _, err := fmt.Sscanf(numStr, "%d", &num)
        if err != nil || num < 1 || num > len(cfg.Perfis) {
                fmt.Println("Número de perfil inválido.")
                return
        }
        
        // Índice real (0-based)
        indice := num - 1
        
        // Configurar cliente
        perfil := cfg.Perfis[indice]
        configAPI := zabbix.ConfigAPI{
                URL:         perfil.URL,
                Token:       perfil.Token,
                TempoLimite: cfg.TempoLimite,
        }
        cliente := zabbix.NovoClienteAPI(configAPI)
        
        // Buscar hosts
        fmt.Println("Buscando hosts...")
        hosts, err := cliente.ObterHosts()
        if err != nil {
                fmt.Printf("Erro ao buscar hosts: %v\n", err)
                return
        }
        
        // Criar diretório para relatórios
        diretorioHome, _ := os.UserHomeDir()
        diretorioRelatorios := fmt.Sprintf("%s/Relatórios Zabbix", diretorioHome)
        err = os.MkdirAll(diretorioRelatorios, 0755)
        if err != nil {
                fmt.Printf("Erro ao criar diretório de relatórios: %v\n", err)
                return
        }
        
        // Nome do arquivo
        caminhoArquivo := fmt.Sprintf("%s/relatorio_%s.csv", diretorioRelatorios, perfil.Nome)
        
        // Gerar relatório
        err = zabbix.GerarRelatorioCSV(hosts, caminhoArquivo)
        if err != nil {
                fmt.Printf("Erro ao gerar relatório: %v\n", err)
                return
        }
        
        fmt.Printf("Relatório exportado com sucesso para: %s\n", caminhoArquivo)
}