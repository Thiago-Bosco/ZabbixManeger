package ui

import (
	"zabbix-manager/config"
	"zabbix-manager/zabbix"
)

// Aplicacao representa a interface da aplicação
type Aplicacao interface {
	// Iniciar inicia a aplicação
	Iniciar() error
	
	// Encerrar encerra a aplicação
	Encerrar() error
	
	// ObterTelaLogin retorna a tela de login
	ObterTelaLogin() TelaLogin
	
	// ObterTelaPrincipal retorna a tela principal
	ObterTelaPrincipal() TelaPrincipal
	
	// ObterTelaConfig retorna a tela de configuração
	ObterTelaConfig() TelaConfig
	
	// ObterClienteAPI retorna o cliente da API
	ObterClienteAPI() *zabbix.ClienteAPI
	
	// ObterConfiguracao retorna a configuração da aplicação
	ObterConfiguracao() *config.Configuração
	
	// ConfigurarAPI configura o cliente da API
	ConfigurarAPI(url, token string) error
}