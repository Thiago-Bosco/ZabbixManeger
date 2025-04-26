package ui

import (
	"zabbix-manager/zabbix"
)

// TelaPrincipal define uma interface para a tela principal
type TelaPrincipal interface {
	// Exibir exibe a tela principal
	Exibir()
	
	// CarregarHosts carrega os hosts do servidor
	CarregarHosts() ([]zabbix.Host, error)
	
	// FiltrarHosts filtra os hosts baseado em um termo de busca
	FiltrarHosts(hosts []zabbix.Host, termo string) []zabbix.Host
	
	// ExportarRelatorio exporta os dados dos hosts para um arquivo CSV
	ExportarRelatorio(hosts []zabbix.Host, nomeArquivo string) error
}

// FiltrarHostsPorTermo filtra uma lista de hosts baseado em um termo de busca
func FiltrarHostsPorTermo(hosts []zabbix.Host, termo string) []zabbix.Host {
	if termo == "" {
		return hosts
	}
	
	resultado := []zabbix.Host{}
	for _, host := range hosts {
		if contem(host.Nome, termo) || contem(host.ID, termo) {
			resultado = append(resultado, host)
		}
	}
	
	return resultado
}

// contem verifica se uma string contém outra, ignorando maiúsculas/minúsculas
func contem(s, substr string) bool {
	// Converter ambas as strings para minúsculas
	s1 := []rune(s)
	s2 := []rune(substr)
	
	// Conversão manual para minúsculas (sem usar strings.ToLower)
	for i := 0; i < len(s1); i++ {
		if s1[i] >= 'A' && s1[i] <= 'Z' {
			s1[i] = s1[i] + ('a' - 'A')
		}
	}
	
	for i := 0; i < len(s2); i++ {
		if s2[i] >= 'A' && s2[i] <= 'Z' {
			s2[i] = s2[i] + ('a' - 'A')
		}
	}
	
	// Verificar se s1 contém s2
	n1, n2 := len(s1), len(s2)
	if n2 > n1 {
		return false
	}
	
	for i := 0; i <= n1-n2; i++ {
		match := true
		for j := 0; j < n2; j++ {
			if s1[i+j] != s2[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	
	return false
}