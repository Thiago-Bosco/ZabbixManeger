package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ConfiguracaoPerfil representa um perfil de configuração para um servidor Zabbix
type ConfiguracaoPerfil struct {
	Nome  string `json:"nome"`  // Nome do perfil
	URL   string `json:"url"`   // URL da API do Zabbix
	Token string `json:"token"` // Token de autenticação da API
}

// Configuração armazena as configurações gerais da aplicação
type Configuração struct {
	Perfis      []ConfiguracaoPerfil `json:"perfis"`       // Lista de perfis de servidores
	PerfilAtual int                  `json:"perfilAtual"`  // Índice do perfil ativo (-1 = nenhum)
	TempoLimite time.Duration        `json:"tempoLimite"`  // Tempo limite para requisições (em segundos)
}

// NovaPadrao cria uma configuração com valores padrão
func NovaPadrao() *Configuração {
	return &Configuração{
		Perfis:      []ConfiguracaoPerfil{},
		PerfilAtual: -1,
		TempoLimite: 30 * time.Second,
	}
}

// Carregar carrega a configuração a partir de um arquivo JSON
func Carregar(caminhoArquivo string) (*Configuração, error) {
	// Verificar se o arquivo existe
	if _, err := os.Stat(caminhoArquivo); os.IsNotExist(err) {
		// Criar diretório se não existir
		dir := filepath.Dir(caminhoArquivo)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("erro ao criar diretório de configuração: %w", err)
		}

		// Retornar configuração padrão
		cfg := NovaPadrao()
		return cfg, nil
	}

	// Abrir arquivo
	arquivo, err := os.Open(caminhoArquivo)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir arquivo de configuração: %w", err)
	}
	defer arquivo.Close()

	// Decodificar JSON
	var cfg Configuração
	err = json.NewDecoder(arquivo).Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("erro ao decodificar configuração: %w", err)
	}

	// Retornar configuração
	return &cfg, nil
}

// Salvar salva a configuração em um arquivo JSON
func (c *Configuração) Salvar(caminhoArquivo string) error {
	// Criar diretório se não existir
	dir := filepath.Dir(caminhoArquivo)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório de configuração: %w", err)
	}

	// Criar arquivo
	arquivo, err := os.Create(caminhoArquivo)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo de configuração: %w", err)
	}
	defer arquivo.Close()

	// Codificar JSON com indentação para facilitar leitura
	encoder := json.NewEncoder(arquivo)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(c)
	if err != nil {
		return fmt.Errorf("erro ao codificar configuração: %w", err)
	}

	return nil
}

// ObterCaminhoConfiguracao retorna o caminho para o arquivo de configuração
func ObterCaminhoConfiguracao() string {
	// Obter diretório home do usuário
	diretorioHome, err := os.UserHomeDir()
	if err != nil {
		// Usar diretório atual se não conseguir obter o home
		diretorioHome, _ = os.Getwd()
	}

	// Retornar caminho para o arquivo de configuração
	return filepath.Join(diretorioHome, ".zabbix-manager", "config.json")
}

// PerfilAtivo retorna o perfil ativo
func (c *Configuração) PerfilAtivo() (*ConfiguracaoPerfil, error) {
	// Verificar se há perfis cadastrados
	if len(c.Perfis) == 0 {
		return nil, fmt.Errorf("não há perfis de servidor cadastrados")
	}

	// Verificar se há perfil ativo
	if c.PerfilAtual < 0 || c.PerfilAtual >= len(c.Perfis) {
		return nil, fmt.Errorf("não há perfil ativo, selecione um perfil")
	}

	// Retornar perfil ativo
	return &c.Perfis[c.PerfilAtual], nil
}

// AdicionarPerfil adiciona um novo perfil de servidor
func (c *Configuração) AdicionarPerfil(perfil ConfiguracaoPerfil) {
	c.Perfis = append(c.Perfis, perfil)

	// Se este for o primeiro perfil, selecioná-lo automaticamente
	if len(c.Perfis) == 1 {
		c.PerfilAtual = 0
	}
}

// SelecionarPerfil seleciona um perfil de servidor
func (c *Configuração) SelecionarPerfil(indice int) error {
	// Verificar se o índice é válido
	if indice < 0 || indice >= len(c.Perfis) {
		return fmt.Errorf("índice de perfil inválido: %d", indice)
	}

	// Selecionar perfil
	c.PerfilAtual = indice
	return nil
}

// RemoverPerfil remove um perfil de servidor
func (c *Configuração) RemoverPerfil(indice int) error {
	// Verificar se o índice é válido
	if indice < 0 || indice >= len(c.Perfis) {
		return fmt.Errorf("índice de perfil inválido: %d", indice)
	}

	// Remover perfil
	c.Perfis = append(c.Perfis[:indice], c.Perfis[indice+1:]...)

	// Ajustar índice do perfil ativo se necessário
	if c.PerfilAtual == indice {
		// Se removeu o perfil ativo, selecionar o primeiro
		if len(c.Perfis) > 0 {
			c.PerfilAtual = 0
		} else {
			c.PerfilAtual = -1
		}
	} else if c.PerfilAtual > indice {
		// Se removeu um perfil antes do ativo, ajustar índice
		c.PerfilAtual--
	}

	return nil
}