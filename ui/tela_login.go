package ui

import (
	"zabbix-manager/config"
	"zabbix-manager/zabbix"
)

// TelaLogin define uma interface para a tela de login
type TelaLogin interface {
	// Exibir exibe a tela de login
	Exibir()
	
	// ProcessarLogin processa o login com os dados fornecidos
	ProcessarLogin(url, token string) error
	
	// AdicionarPerfil adiciona um novo perfil de servidor
	AdicionarPerfil(nome, url, token string) error
	
	// EditarPerfil edita um perfil existente
	EditarPerfil(indice int, nome, url, token string) error
	
	// RemoverPerfil remove um perfil existente
	RemoverPerfil(indice int) error
	
	// ObterPerfis retorna a lista de perfis configurados
	ObterPerfis() []config.ConfiguracaoPerfil
}

// TestarConexao testa a conexão com o servidor Zabbix
func TestarConexao(url, token string) error {
	// Configurar cliente API temporário
	configAPI := zabbix.ConfigAPI{
		URL:         url,
		Token:       token,
		TempoLimite: 30, // 30 segundos
	}
	clienteAPI := zabbix.NovoClienteAPI(configAPI)
	
	// Testar conexão
	return clienteAPI.TestarConexao()
}