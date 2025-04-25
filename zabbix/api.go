package zabbix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ClienteAPI representa o cliente de API do Zabbix
type ClienteAPI struct {
	Config     ConfigAPI
	HTTPClient *http.Client
}

// NovoClienteAPI cria uma nova instância do cliente de API do Zabbix
func NovoClienteAPI(config ConfigAPI) *ClienteAPI {
	// Se o tempo limite não estiver configurado, usar 30 segundos como padrão
	if config.TempoLimite <= 0 {
		config.TempoLimite = 30
	}

	return &ClienteAPI{
		Config: config,
		HTTPClient: &http.Client{
			Timeout: time.Duration(config.TempoLimite) * time.Second,
		},
	}
}

// ObterHosts retorna todos os hosts cadastrados no Zabbix
func (c *ClienteAPI) ObterHosts() ([]Host, error) {
	cabecalhos := map[string]string{
		"Content-Type": "application/json-rpc",
	}

	pedido := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "host.get",
		"params": map[string]interface{}{
			"output":         []string{"hostid", "host", "status"},
			"selectItems":    []string{"itemid", "name"},
			"selectTriggers": []string{"triggerid", "description"},
		},
		"auth": c.Config.Token,
		"id":   1,
	}

	var resposta Resposta
	erro := c.enviarRequisicao(c.Config.URL, cabecalhos, pedido, &resposta)
	if erro != nil {
		return nil, erro
	}

	return resposta.Resultados, nil
}

// TestarConexao testa a conexão com a API do Zabbix
func (c *ClienteAPI) TestarConexao() error {
	cabecalhos := map[string]string{
		"Content-Type": "application/json-rpc",
	}

	pedido := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "apiinfo.version",
		"params":  map[string]interface{}{},
		"id":      1,
	}

	var resposta interface{}
	return c.enviarRequisicao(c.Config.URL, cabecalhos, pedido, &resposta)
}

// ObterHostsFiltrados retorna hosts filtrados por termo de busca
func (c *ClienteAPI) ObterHostsFiltrados(termoBusca string) ([]Host, error) {
	cabecalhos := map[string]string{
		"Content-Type": "application/json-rpc",
	}

	pedido := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "host.get",
		"params": map[string]interface{}{
			"output":         []string{"hostid", "host", "status"},
			"selectItems":    []string{"itemid", "name"},
			"selectTriggers": []string{"triggerid", "description"},
			"search": map[string]string{
				"host": termoBusca,
			},
			"searchByAny": true,
		},
		"auth": c.Config.Token,
		"id":   1,
	}

	var resposta Resposta
	erro := c.enviarRequisicao(c.Config.URL, cabecalhos, pedido, &resposta)
	if erro != nil {
		return nil, erro
	}

	return resposta.Resultados, nil
}

// Autenticar realiza a autenticação na API do Zabbix
func (c *ClienteAPI) Autenticar(usuario, senha string) (string, error) {
	cabecalhos := map[string]string{
		"Content-Type": "application/json-rpc",
	}

	pedido := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "user.login",
		"params": map[string]string{
			"user":     usuario,
			"password": senha,
		},
		"id": 1,
	}

	type RespostaAuth struct {
		Jsonrpc string `json:"jsonrpc"`
		Result  string `json:"result"`
		ID      int    `json:"id"`
	}

	var resposta RespostaAuth
	erro := c.enviarRequisicao(c.Config.URL, cabecalhos, pedido, &resposta)
	if erro != nil {
		return "", erro
	}

	if resposta.Result == "" {
		return "", fmt.Errorf("falha na autenticação: token vazio")
	}

	// Atualizar o token no cliente
	c.Config.Token = resposta.Result
	return resposta.Result, nil
}

// enviarRequisicao envia uma requisição para a API do Zabbix
func (c *ClienteAPI) enviarRequisicao(url string, cabecalhos map[string]string, pedido map[string]interface{}, resposta interface{}) error {
	pedidoBytes, erro := json.Marshal(pedido)
	if erro != nil {
		return fmt.Errorf("erro ao criar o pedido JSON: %w", erro)
	}

	req, erro := http.NewRequest("POST", url, bytes.NewBuffer(pedidoBytes))
	if erro != nil {
		return fmt.Errorf("erro ao criar a requisição: %w", erro)
	}

	for chave, valor := range cabecalhos {
		req.Header.Set(chave, valor)
	}

	resp, erro := c.HTTPClient.Do(req)
	if erro != nil {
		return fmt.Errorf("erro na requisição: %w", erro)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("erro na requisição, código de status: %d", resp.StatusCode)
	}

	erro = json.NewDecoder(resp.Body).Decode(resposta)
	if erro != nil {
		return fmt.Errorf("erro ao decodificar a resposta JSON: %w", erro)
	}

	return nil
}
