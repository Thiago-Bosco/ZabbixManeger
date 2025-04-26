
package ui

import (
	"zabbix-manager/config"
	"zabbix-manager/zabbix"
)

// Aplicacao representa a interface da aplicação web
type Aplicacao interface {
	// Iniciar inicia a aplicação web
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
}
