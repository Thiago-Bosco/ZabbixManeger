package zabbix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// ConfigAPI define as configurações para a API do Zabbix
type ConfigAPI struct {
	URL         string
	Token       string
	TempoLimite int
}

// ClienteAPI representa um cliente para a API do Zabbix
type ClienteAPI struct {
	URL         string
	Token       string
	Cliente     *http.Client
	TempoLimite int
}

// NovoClienteAPI cria um novo cliente para a API do Zabbix
func NovoClienteAPI(config ConfigAPI) *ClienteAPI {
	// Usar tempo limite padrão de 30 segundos se não for especificado
	timeout := config.TempoLimite
	if timeout <= 0 {
		timeout = 30
	}

	return &ClienteAPI{
		URL:         config.URL,
		Token:       config.Token,
		Cliente:     &http.Client{Timeout: time.Duration(timeout) * time.Second},
		TempoLimite: timeout,
	}
}

// TestarConexao testa a conexão com o servidor Zabbix
func (c *ClienteAPI) TestarConexao() error {
	// Criar uma requisição simples para testar a conexão
	requisicao := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "apiinfo.version",
		"params":  []string{},
		"id":      1,
	}

	// Enviar a requisição
	resposta, err := c.enviarRequisicao(requisicao)
	if err != nil {
		return fmt.Errorf("erro ao testar conexão: %w", err)
	}

	// Verificar se há erro na resposta
	if resposta["error"] != nil {
		erro := resposta["error"].(map[string]interface{})
		return fmt.Errorf("erro da API: %s", erro["data"])
	}

	return nil
}

// ObterHosts retorna todos os hosts do servidor Zabbix
func (c *ClienteAPI) ObterHosts() ([]Host, error) {
	log.Println("Obtendo hosts do servidor Zabbix...")

	// Montar a requisição para obter os hosts
	requisicao := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "host.get",
		"params": map[string]interface{}{
			"output":    "extend",
			"selectItems": []string{"itemid", "name", "key_", "status"},
			"selectTriggers": []string{"triggerid", "description", "status", "priority"},
		},
		"auth": c.Token,
		"id":   1,
	}

	// Enviar a requisição
	resposta, err := c.enviarRequisicao(requisicao)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter hosts: %w", err)
	}

	// Verificar se há erro na resposta
	if resposta["error"] != nil {
		erro := resposta["error"].(map[string]interface{})
		return nil, fmt.Errorf("erro da API: %s", erro["data"])
	}

	// Processar a resposta
	hostsRaw, ok := resposta["result"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("formato de resposta inválido")
	}

	// Converter para o tipo Host
	hosts := []Host{}
	for _, hostRaw := range hostsRaw {
		hostMap := hostRaw.(map[string]interface{})
		
		// Processar itens
		var items []Item
		itemsRaw, ok := hostMap["items"].([]interface{})
		if ok {
			for _, itemRaw := range itemsRaw {
				itemMap := itemRaw.(map[string]interface{})
				status := "0" // Padrão: ativo
				if statusRaw, ok := itemMap["status"].(string); ok {
					status = statusRaw
				}

				item := Item{
					ID:     itemMap["itemid"].(string),
					Nome:   itemMap["name"].(string),
					Chave:  itemMap["key_"].(string),
					Status: status,
				}
				items = append(items, item)
			}
		}

		// Processar triggers
		var triggers []Trigger
		triggersRaw, ok := hostMap["triggers"].([]interface{})
		if ok {
			for _, triggerRaw := range triggersRaw {
				triggerMap := triggerRaw.(map[string]interface{})
				
				status := "0" // Padrão: ativo
				if statusRaw, ok := triggerMap["status"].(string); ok {
					status = statusRaw
				}
				
				prioridade := "0" // Padrão: não classificada
				if prioridadeRaw, ok := triggerMap["priority"].(string); ok {
					prioridade = prioridadeRaw
				}

				trigger := Trigger{
					ID:          triggerMap["triggerid"].(string),
					Descricao:   triggerMap["description"].(string),
					Status:      status,
					Prioridade:  prioridade,
				}
				triggers = append(triggers, trigger)
			}
		}

		status := "0" // Padrão: ativo
		if statusRaw, ok := hostMap["status"].(string); ok {
			status = statusRaw
		}

		host := Host{
			ID:       hostMap["hostid"].(string),
			Nome:     hostMap["name"].(string),
			Status:   status,
			Items:    items,
			Triggers: triggers,
		}
		hosts = append(hosts, host)
	}

	log.Printf("Obtidos %d hosts do servidor Zabbix", len(hosts))
	return hosts, nil
}

// enviarRequisicao envia uma requisição para a API do Zabbix
func (c *ClienteAPI) enviarRequisicao(dados map[string]interface{}) (map[string]interface{}, error) {
	// Converter a requisição para JSON
	jsonData, err := json.Marshal(dados)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar requisição: %w", err)
	}

	// Criar a requisição HTTP
	req, err := http.NewRequest("POST", c.URL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição HTTP: %w", err)
	}

	// Configurar headers
	req.Header.Set("Content-Type", "application/json")

	// Enviar a requisição
	resp, err := c.Cliente.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao enviar requisição: %w", err)
	}
	defer resp.Body.Close()

	// Ler a resposta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta: %w", err)
	}

	// Verificar código de status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status inválido: %d - %s", resp.StatusCode, string(body))
	}

	// Converter a resposta de JSON para map
	var resultado map[string]interface{}
	err = json.Unmarshal(body, &resultado)
	if err != nil {
		return nil, fmt.Errorf("erro ao desserializar resposta: %w", err)
	}

	return resultado, nil
}