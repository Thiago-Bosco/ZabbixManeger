package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Perfil representa um perfil de conexão do Zabbix
type Perfil struct {
	Nome     string `json:"nome"`
	URL      string `json:"url"`
	Token    string `json:"token"`
	Usuário  string `json:"usuario"`
	Senha    string `json:"senha"`
	Salvar   bool   `json:"salvar"`
}

// Configuração da aplicação
type Configuração struct {
	TempoLimite    int               `json:"tempo_limite"`
	PerfisZabbix   []Perfil          `json:"perfis_zabbix"`
	PerfilAtivo    string            `json:"perfil_ativo"`
	DiretórioCSV   string            `json:"diretorio_csv"`
	mu             sync.Mutex
	arquivoConfig  string
}

// Nova cria uma nova instância de configuração
func Nova(arquivoConfig string) *Configuração {
	// Configuração padrão
	config := &Configuração{
		TempoLimite:   30,
		PerfisZabbix:  []Perfil{},
		PerfilAtivo:   "",
		DiretórioCSV:  "",
		arquivoConfig: arquivoConfig,
	}

	// Adicionar perfil padrão se não existir arquivo de configuração
	if _, err := os.Stat(arquivoConfig); os.IsNotExist(err) {
		config.PerfisZabbix = append(config.PerfisZabbix, Perfil{
			Nome:     "Padrão",
			URL:      "http://localhost/zabbix/api_jsonrpc.php",
			Token:    "",
			Usuário:  "",
			Senha:    "",
			Salvar:   false,
		})
		config.PerfilAtivo = "Padrão"
		
		// Criar o diretório se não existir
		dir := filepath.Dir(arquivoConfig)
		if err := os.MkdirAll(dir, 0755); err == nil {
			config.Salvar()
		}
	} else {
		config.Carregar()
	}

	return config
}

// Carregar carrega a configuração do arquivo
func (c *Configuração) Carregar() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Verificar se o arquivo existe
	if _, err := os.Stat(c.arquivoConfig); os.IsNotExist(err) {
		return nil // Não há erro, apenas não existe o arquivo ainda
	}

	// Ler o arquivo
	dados, err := os.ReadFile(c.arquivoConfig)
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo de configuração: %w", err)
	}

	// Decodificar o JSON
	err = json.Unmarshal(dados, c)
	if err != nil {
		return fmt.Errorf("erro ao decodificar configuração: %w", err)
	}

	return nil
}

// Salvar salva a configuração no arquivo
func (c *Configuração) Salvar() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Criar o diretório se não existir
	dir := filepath.Dir(c.arquivoConfig)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório de configuração: %w", err)
	}

	// Serializar para JSON
	dados, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("erro ao serializar configuração: %w", err)
	}

	// Escrever no arquivo
	err = os.WriteFile(c.arquivoConfig, dados, 0644)
	if err != nil {
		return fmt.Errorf("erro ao salvar configuração: %w", err)
	}

	return nil
}

// AdicionarPerfil adiciona ou atualiza um perfil de conexão
func (c *Configuração) AdicionarPerfil(perfil Perfil) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Verificar se o perfil já existe
	for i, p := range c.PerfisZabbix {
		if p.Nome == perfil.Nome {
			// Atualizar o perfil existente
			c.PerfisZabbix[i] = perfil
			return
		}
	}

	// Adicionar novo perfil
	c.PerfisZabbix = append(c.PerfisZabbix, perfil)
}

// RemoverPerfil remove um perfil pelo nome
func (c *Configuração) RemoverPerfil(nome string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Filtrar os perfis, removendo o que tem o nome especificado
	novoPerfis := []Perfil{}
	for _, p := range c.PerfisZabbix {
		if p.Nome != nome {
			novoPerfis = append(novoPerfis, p)
		}
	}

	c.PerfisZabbix = novoPerfis

	// Se removeu o perfil ativo, selecionar outro perfil se houver
	if c.PerfilAtivo == nome && len(c.PerfisZabbix) > 0 {
		c.PerfilAtivo = c.PerfisZabbix[0].Nome
	} else if len(c.PerfisZabbix) == 0 {
		c.PerfilAtivo = ""
	}
}

// ObterPerfil obtém um perfil pelo nome
func (c *Configuração) ObterPerfil(nome string) (Perfil, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, p := range c.PerfisZabbix {
		if p.Nome == nome {
			return p, true
		}
	}

	return Perfil{}, false
}

// ObterPerfilAtivo obtém o perfil ativo
func (c *Configuração) ObterPerfilAtivo() (Perfil, bool) {
	if c.PerfilAtivo == "" {
		return Perfil{}, false
	}
	return c.ObterPerfil(c.PerfilAtivo)
}

// DefinirPerfilAtivo define o perfil ativo pelo nome
func (c *Configuração) DefinirPerfilAtivo(nome string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Verificar se o perfil existe
	for _, p := range c.PerfisZabbix {
		if p.Nome == nome {
			c.PerfilAtivo = nome
			return true
		}
	}

	return false
}

// AtualizarPerfilAtivo atualiza os dados do perfil ativo
func (c *Configuração) AtualizarPerfilAtivo(perfil Perfil) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Verificar se o perfil ativo existe
	for i, p := range c.PerfisZabbix {
		if p.Nome == c.PerfilAtivo {
			// Manter o nome original
			perfil.Nome = p.Nome
			
			// Atualizar o perfil
			c.PerfisZabbix[i] = perfil
			return true
		}
	}

	return false
}