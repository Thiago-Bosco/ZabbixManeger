package zabbix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// ConfigAPI armazena configurações para conexão com a API Zabbix
type ConfigAPI struct {
	URL         string        // URL do servidor (Ex: http://zabbix.example.com)
	Token       string        // Token de autenticação da API
	TempoLimite time.Duration // Tempo limite para requisições (em segundos)
}

// ClienteAPI encapsula funcionalidades para interagir com a API do Zabbix
type ClienteAPI struct {
	config ConfigAPI
	client *http.Client
}



// RespostaAPI encapsula a resposta da API do Zabbix
type RespostaAPI struct {
	Jsonrpc string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    string `json:"data"`
	} `json:"error"`
	ID int `json:"id"`
}

// NovoClienteAPI cria uma nova instância de ClienteAPI
func NovoClienteAPI(config ConfigAPI) *ClienteAPI {
	// Definir tempo limite padrão se não especificado
	if config.TempoLimite <= 0 {
		config.TempoLimite = 30 * time.Second
	}

	// Criar cliente HTTP com timeout configurado
	client := &http.Client{
		Timeout: config.TempoLimite,
	}

	return &ClienteAPI{
		config: config,
		client: client,
	}
}

// TestarConexao verifica se a conexão com a API do Zabbix está funcionando
func (c *ClienteAPI) TestarConexao() error {
	// Requisição para obter a versão da API (método simples para testar conexão)
	pedido := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "apiinfo.version",
		"params":  map[string]interface{}{},
		"id":      1,
	}

	var resposta RespostaAPI
	err := c.realizarRequisicao(pedido, &resposta)
	if err != nil {
		return err
	}

	// Verificar se houve erro na resposta
	if resposta.Error != nil {
		return fmt.Errorf("erro na API: %s - %s", resposta.Error.Message, resposta.Error.Data)
	}

	return nil
}

// ObterHosts retorna a lista de hosts do Zabbix com seus itens e triggers
// ObterProblemasPeriodo obtém problemas de um período específico
func (c *ClienteAPI) ObterProblemasPeriodo(inicio, fim time.Time) ([]Problema, error) {
	pedido := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "problem.get",
		"params": map[string]interface{}{
			"output":      "extend",
			"time_from":   inicio.Unix(),
			"time_till":   fim.Unix(),
			"sortfield":   []string{"eventid"},
			"selectHosts": []string{"hostid", "host"},
		},
		"auth": c.config.Token,
		"id":   1,
	}

	var resposta RespostaAPI
	err := c.realizarRequisicao(pedido, &resposta)
	if err != nil {
		return nil, err
	}

	var problemas []Problema
	err = json.Unmarshal(resposta.Result, &problemas)
	return problemas, err
}

// AnalisarProblemasMensais analisa problemas de um mês específico
func (c *ClienteAPI) AnalisarProblemasMensais(ano int, mes int) ([]AnaliseMensal, error) {
	inicio := time.Date(ano, time.Month(mes), 1, 0, 0, 0, 0, time.UTC)
	fim := inicio.AddDate(0, 1, 0).Add(-time.Second)
	
	problemas, err := c.ObterProblemasPeriodo(inicio, fim)
	if err != nil {
		return nil, err
	}

	analises := make(map[string]*AnaliseMensal)
	for _, p := range problemas {
		if _, existe := analises[p.HostID]; !existe {
			analises[p.HostID] = &AnaliseMensal{
				HostID: p.HostID,
				ProblemasPorTrigger: make(map[string]int),
			}
		}
		
		analise := analises[p.HostID]
		analise.TotalProblemas++
		analise.ProblemasPorTrigger[p.TriggerID]++
		
		if p.Valor == "1" { // Problema ativo
			analise.LimitesExcedidos++
		}
	}

	resultado := make([]AnaliseMensal, 0, len(analises))
	for _, a := range analises {
		resultado = append(resultado, *a)
	}
	
	return resultado, nil
}

func (c *ClienteAPI) ObterHosts() ([]Host, error) {
	// Preparar requisição para obter hosts com itens e triggers
	pedido := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "host.get",
		"params": map[string]interface{}{
			"output":         []string{"hostid", "host", "status"},
			"selectItems":    []string{"itemid", "name"},
			"selectTriggers": []string{"triggerid", "description"},
		},
		"auth": c.config.Token,
		"id":   1,
	}

	var resposta RespostaAPI
	err := c.realizarRequisicao(pedido, &resposta)
	if err != nil {
		return nil, err
	}

	// Verificar se houve erro na resposta
	if resposta.Error != nil {
		return nil, fmt.Errorf("erro na API: %s - %s", resposta.Error.Message, resposta.Error.Data)
	}

	// Decodificar resultado para slice de hosts
	var hosts []Host
	err = json.Unmarshal(resposta.Result, &hosts)
	if err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	return hosts, nil
}

// realizarRequisicao envia uma requisição para a API do Zabbix
func (c *ClienteAPI) realizarRequisicao(pedido map[string]interface{}, resposta *RespostaAPI) error {
	// Converter pedido para JSON
	pedidoBytes, err := json.Marshal(pedido)
	if err != nil {
		return fmt.Errorf("erro ao criar pedido JSON: %w", err)
	}

	// Criar requisição HTTP
	apiURL := strings.TrimRight(c.config.URL, "/") + "/api_jsonrpc.php"
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(pedidoBytes))
	if err != nil {
		return fmt.Errorf("erro ao criar requisição: %w", err)
	}

	// Definir cabeçalhos
	req.Header.Set("Content-Type", "application/json-rpc")

	// Enviar requisição
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("erro na requisição: %w", err)
	}
	defer resp.Body.Close()

	// Verificar código de status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("erro na API, código de status: %d", resp.StatusCode)
	}

	// Decodificar resposta
	err = json.NewDecoder(resp.Body).Decode(resposta)
	if err != nil {
		return fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	return nil
}