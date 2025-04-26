package ui

import (
	"zabbix-manager/config"
)

// TelaConfig define uma interface para a tela de configuração
type TelaConfig interface {
	// Exibir exibe a tela de configuração
	Exibir()
	
	// SalvarConfiguracoes salva as configurações
	SalvarConfiguracoes(tempoLimite int) error
	
	// ObterConfiguracoes retorna as configurações atuais
	ObterConfiguracoes() *config.Configuração
}

// AplicarConfiguracao aplica as configurações globais
func AplicarConfiguracao(cfg *config.Configuração) error {
	// Obter caminho do arquivo de configuração
	caminho := config.ObterCaminhoConfiguracao()
	
	// Salvar configuração
	return cfg.Salvar(caminho)
}