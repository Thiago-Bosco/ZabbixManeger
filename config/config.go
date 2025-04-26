package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// ConfiguracaoPerfil representa um perfil de configuração para um servidor Zabbix
type ConfiguracaoPerfil struct {
	Nome  string `json:"nome"`
	URL   string `json:"url"`
	Token string `json:"token"`
}

// Configuração representa a configuração do aplicativo
type Configuração struct {
	Perfis      []ConfiguracaoPerfil `json:"perfis"`
	PerfilAtual int                  `json:"perfil_atual"`
	TempoLimite int                  `json:"tempo_limite"`
}

// NovaConfiguração cria uma nova configuração com valores padrão
func NovaConfiguração() *Configuração {
	return &Configuração{
		Perfis:      []ConfiguracaoPerfil{},
		PerfilAtual: -1,
		TempoLimite: 30,
	}
}

// AdicionarPerfil adiciona um novo perfil à configuração
func (c *Configuração) AdicionarPerfil(perfil ConfiguracaoPerfil) {
	c.Perfis = append(c.Perfis, perfil)
	// Se for o primeiro perfil, selecionar automaticamente
	if len(c.Perfis) == 1 {
		c.PerfilAtual = 0
	}
}

// AtualizarPerfil atualiza um perfil existente
func (c *Configuração) AtualizarPerfil(indice int, perfil ConfiguracaoPerfil) error {
	if indice < 0 || indice >= len(c.Perfis) {
		return fmt.Errorf("índice de perfil inválido: %d", indice)
	}

	c.Perfis[indice] = perfil
	return nil
}

// RemoverPerfil remove um perfil existente
func (c *Configuração) RemoverPerfil(indice int) error {
	if indice < 0 || indice >= len(c.Perfis) {
		return fmt.Errorf("índice de perfil inválido: %d", indice)
	}

	// Remover o perfil do slice
	c.Perfis = append(c.Perfis[:indice], c.Perfis[indice+1:]...)

	// Ajustar o índice do perfil atual se necessário
	if c.PerfilAtual == indice {
		// Se removemos o perfil atual, selecionar outro
		if len(c.Perfis) > 0 {
			c.PerfilAtual = 0
		} else {
			c.PerfilAtual = -1
		}
	} else if c.PerfilAtual > indice {
		// Se removemos um perfil antes do atual, atualizar o índice
		c.PerfilAtual--
	}

	return nil
}

// PerfilAtivo retorna o perfil atualmente selecionado
func (c *Configuração) PerfilAtivo() (*ConfiguracaoPerfil, error) {
	if c.PerfilAtual < 0 || c.PerfilAtual >= len(c.Perfis) {
		return nil, fmt.Errorf("nenhum perfil selecionado")
	}
	return &c.Perfis[c.PerfilAtual], nil
}

// SelecionarPerfil seleciona um perfil pelo índice
func (c *Configuração) SelecionarPerfil(indice int) error {
	if indice < 0 || indice >= len(c.Perfis) {
		return fmt.Errorf("índice de perfil inválido: %d", indice)
	}
	c.PerfilAtual = indice
	return nil
}

// Carregar carrega a configuração do arquivo
func Carregar(caminhoArquivo string) (*Configuração, error) {
	// Verificar se o arquivo existe
	if _, err := os.Stat(caminhoArquivo); os.IsNotExist(err) {
		// Se não existir, criar uma nova configuração
		config := NovaConfiguração()
		return config, nil
	}

	// Ler o conteúdo do arquivo
	dados, err := ioutil.ReadFile(caminhoArquivo)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo de configuração: %w", err)
	}

	// Desserializar o JSON
	var config Configuração
	err = json.Unmarshal(dados, &config)
	if err != nil {
		return nil, fmt.Errorf("erro ao analisar configuração: %w", err)
	}

	return &config, nil
}

// Salvar salva a configuração no arquivo
func (c *Configuração) Salvar(caminhoArquivo string) error {
	// Criar o diretório se não existir
	diretorio := filepath.Dir(caminhoArquivo)
	err := os.MkdirAll(diretorio, 0755)
	if err != nil {
		return fmt.Errorf("erro ao criar diretório de configuração: %w", err)
	}

	// Serializar a configuração para JSON
	dados, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("erro ao serializar configuração: %w", err)
	}

	// Escrever no arquivo
	err = ioutil.WriteFile(caminhoArquivo, dados, 0644)
	if err != nil {
		return fmt.Errorf("erro ao escrever arquivo de configuração: %w", err)
	}

	return nil
}

// ObterCaminhoConfiguracao retorna o caminho para o arquivo de configuração
func ObterCaminhoConfiguracao() string {
	// Obter o diretório home do usuário
	diretorioHome, err := os.UserHomeDir()
	if err != nil {
		// Fallback para o diretório atual se não conseguir obter o home
		diretorioHome, _ = os.Getwd()
	}

	// Criar o caminho para o arquivo de configuração
	return filepath.Join(diretorioHome, ".zabbix-manager", "config.json")
}