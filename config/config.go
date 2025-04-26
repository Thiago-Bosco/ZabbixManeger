package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Configuração armazena as configurações da aplicação
type Configuração struct {
	TempoLimite int // em segundos
	PerfilAtual int
	Perfis      []ConfiguracaoPerfil
}

// ConfiguracaoPerfil armazena a configuração de um servidor Zabbix
type ConfiguracaoPerfil struct {
	Nome  string
	URL   string
	Token string
}

// NovaConfiguração cria uma nova configuração com valores padrão
func NovaConfiguração() *Configuração {
	return &Configuração{
		TempoLimite: 30, // 30 segundos padrão
		PerfilAtual: -1,
		Perfis:      []ConfiguracaoPerfil{},
	}
}

// Carregar carrega a configuração do arquivo
func Carregar(caminho string) (*Configuração, error) {
	// Verificar se o arquivo existe
	if _, err := os.Stat(caminho); os.IsNotExist(err) {
		// Criar configuração padrão
		cfg := NovaConfiguração()
		// Salvar configuração
		if err := cfg.Salvar(caminho); err != nil {
			return nil, fmt.Errorf("erro ao salvar configuração padrão: %v", err)
		}
		return cfg, nil
	}

	// Abrir arquivo
	arquivo, err := os.Open(caminho)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir arquivo: %v", err)
	}
	defer arquivo.Close()

	// Decodificar JSON
	var cfg Configuração
	decodificador := json.NewDecoder(arquivo)
	if err := decodificador.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("erro ao decodificar JSON: %v", err)
	}

	return &cfg, nil
}

// Salvar salva a configuração no arquivo
func (c *Configuração) Salvar(caminho string) error {
	// Criar diretório pai, se necessário
	diretorioPai := filepath.Dir(caminho)
	if err := os.MkdirAll(diretorioPai, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório: %v", err)
	}

	// Criar arquivo
	arquivo, err := os.Create(caminho)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo: %v", err)
	}
	defer arquivo.Close()

	// Codificar JSON
	codificador := json.NewEncoder(arquivo)
	codificador.SetIndent("", "  ")
	if err := codificador.Encode(c); err != nil {
		return fmt.Errorf("erro ao codificar JSON: %v", err)
	}

	return nil
}

// AdicionarPerfil adiciona um perfil de servidor à configuração
func (c *Configuração) AdicionarPerfil(perfil ConfiguracaoPerfil) {
	c.Perfis = append(c.Perfis, perfil)
	// Se for o primeiro perfil, selecioná-lo como atual
	if len(c.Perfis) == 1 {
		c.PerfilAtual = 0
	}
}

// RemoverPerfil remove um perfil de servidor da configuração
func (c *Configuração) RemoverPerfil(indice int) error {
	if indice < 0 || indice >= len(c.Perfis) {
		return errors.New("índice de perfil inválido")
	}

	// Remover perfil
	c.Perfis = append(c.Perfis[:indice], c.Perfis[indice+1:]...)

	// Atualizar perfil atual, se necessário
	if c.PerfilAtual == indice {
		if len(c.Perfis) > 0 {
			c.PerfilAtual = 0
		} else {
			c.PerfilAtual = -1
		}
	} else if c.PerfilAtual > indice {
		c.PerfilAtual--
	}

	return nil
}

// AtualizarPerfil atualiza um perfil existente
func (c *Configuração) AtualizarPerfil(indice int, perfil ConfiguracaoPerfil) error {
	if indice < 0 || indice >= len(c.Perfis) {
		return errors.New("índice de perfil inválido")
	}

	c.Perfis[indice] = perfil
	return nil
}

// SelecionarPerfil seleciona um perfil como atual
func (c *Configuração) SelecionarPerfil(indice int) error {
	if indice < 0 || indice >= len(c.Perfis) {
		return errors.New("índice de perfil inválido")
	}

	c.PerfilAtual = indice
	return nil
}

// PerfilAtivo retorna o perfil ativo ou erro, se não houver
func (c *Configuração) PerfilAtivo() (*ConfiguracaoPerfil, error) {
	if c.PerfilAtual < 0 || c.PerfilAtual >= len(c.Perfis) {
		return nil, errors.New("nenhum perfil ativo")
	}

	return &c.Perfis[c.PerfilAtual], nil
}

// ObterCaminhoConfiguracao retorna o caminho para o arquivo de configuração
func ObterCaminhoConfiguracao() string {
	// Obter diretório home do usuário
	diretorioHome, err := os.UserHomeDir()
	if err != nil {
		// Usar diretório atual se não conseguir obter o home
		diretorioHome, _ = os.Getwd()
	}

	// Caminho para o arquivo de configuração
	return filepath.Join(diretorioHome, ".zabbix-manager", "config.json")
}