package zabbix

import (
        "bytes"
        "encoding/json"
        "errors"
        "fmt"
        "io/ioutil"
        "net/http"
        "time"
)

// ConfigAPI armazena a configuração da API do Zabbix
type ConfigAPI struct {
        URL         string
        Token       string
        TempoLimite int // em segundos
}

// ClienteAPI representa um cliente para API do Zabbix
type ClienteAPI struct {
        URL         string
        Token       string
        Cliente     *http.Client
        TempoLimite int
}

// NovoClienteAPI cria um novo cliente de API do Zabbix
func NovoClienteAPI(config ConfigAPI) *ClienteAPI {
        // Definir tempo limite padrão se não especificado
        tempoLimite := config.TempoLimite
        if tempoLimite <= 0 {
                tempoLimite = 30 // 30 segundos padrão
        }

        return &ClienteAPI{
                URL:         config.URL,
                Token:       config.Token,
                TempoLimite: tempoLimite,
                Cliente: &http.Client{
                        Timeout: time.Duration(tempoLimite) * time.Second,
                },
        }
}

// TestarConexao verifica se a conexão com a API está funcionando
func (c *ClienteAPI) TestarConexao() error {
        // Definir dados da requisição
        requisicao := map[string]interface{}{
                "jsonrpc": "2.0",
                "method":  "apiinfo.version",
                "params":  map[string]interface{}{},
                "id":      1,
        }

        // Enviar requisição
        resposta, err := c.enviarRequisicao(requisicao)
        if err != nil {
                return err
        }

        // Verificar se há erro na resposta
        if erro, ok := resposta["error"].(map[string]interface{}); ok {
                return fmt.Errorf("erro API: %v", erro["data"])
        }

        // Verificar se há resultado
        if _, ok := resposta["result"]; !ok {
                return errors.New("resposta sem resultado")
        }

        return nil
}

// ObterHosts recupera a lista de hosts do Zabbix
func (c *ClienteAPI) ObterHosts() ([]Host, error) {
        // Definir dados da requisição
        requisicao := map[string]interface{}{
                "jsonrpc": "2.0",
                "method":  "host.get",
                "params": map[string]interface{}{
                        "output":                []string{"hostid", "host", "name", "status"},
                        "selectItems":           []string{"itemid", "name", "key_", "status"},
                        "selectTriggers":        []string{"triggerid", "description", "status", "priority"},
                        "selectInterfaces":      []string{"interfaceid", "ip", "dns", "port", "type"},
                        "selectInventory":       []string{"os", "os_full", "hardware", "serialno_a"},
                        "selectGroups":          []string{"groupid", "name"},
                        "selectParentTemplates": []string{"templateid", "name"},
                },
                "auth": c.Token,
                "id":   2,
        }

        // Enviar requisição
        resposta, err := c.enviarRequisicao(requisicao)
        if err != nil {
                return nil, err
        }

        // Verificar se há erro na resposta
        if erro, ok := resposta["error"].(map[string]interface{}); ok {
                return nil, fmt.Errorf("erro API: %v", erro["data"])
        }

        // Verificar se há resultado
        resultado, ok := resposta["result"].([]interface{})
        if !ok {
                return nil, errors.New("resposta sem resultado válido")
        }

        // Converter resultado para lista de hosts
        hosts := []Host{}
        for _, item := range resultado {
                hostData := item.(map[string]interface{})
                
                // Construir objeto Host
                host := Host{
                        ID:     hostData["hostid"].(string),
                        Nome:   hostData["name"].(string),
                        Status: 0,
                }
                
                // Converter status
                if statusStr, ok := hostData["status"].(string); ok {
                        host.Status = parseInt(statusStr)
                }
                
                // Processar itens
                if itemsData, ok := hostData["items"].([]interface{}); ok {
                        for _, itemData := range itemsData {
                                item := Item{}
                                itemMap := itemData.(map[string]interface{})
                                
                                item.ID = getString(itemMap, "itemid")
                                item.Nome = getString(itemMap, "name")
                                item.Chave = getString(itemMap, "key_")
                                item.Status = parseInt(getString(itemMap, "status"))
                                
                                host.Items = append(host.Items, item)
                        }
                }
                
                // Processar triggers
                if triggersData, ok := hostData["triggers"].([]interface{}); ok {
                        for _, triggerData := range triggersData {
                                trigger := Trigger{}
                                triggerMap := triggerData.(map[string]interface{})
                                
                                trigger.ID = getString(triggerMap, "triggerid")
                                trigger.Descricao = getString(triggerMap, "description")
                                trigger.Status = parseInt(getString(triggerMap, "status"))
                                trigger.Prioridade = parseInt(getString(triggerMap, "priority"))
                                
                                host.Triggers = append(host.Triggers, trigger)
                        }
                }
                
                hosts = append(hosts, host)
        }

        return hosts, nil
}

// enviarRequisicao envia uma requisição para a API do Zabbix
func (c *ClienteAPI) enviarRequisicao(dados map[string]interface{}) (map[string]interface{}, error) {
        // Converter dados para JSON
        jsonData, err := json.Marshal(dados)
        if err != nil {
                return nil, fmt.Errorf("erro ao converter para JSON: %v", err)
        }

        // Criar requisição HTTP
        req, err := http.NewRequest("POST", c.URL, bytes.NewBuffer(jsonData))
        if err != nil {
                return nil, fmt.Errorf("erro ao criar requisição: %v", err)
        }

        // Definir cabeçalhos
        req.Header.Set("Content-Type", "application/json")

        // Enviar requisição
        resp, err := c.Cliente.Do(req)
        if err != nil {
                return nil, fmt.Errorf("erro ao enviar requisição: %v", err)
        }
        defer resp.Body.Close()

        // Ler corpo da resposta
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
                return nil, fmt.Errorf("erro ao ler resposta: %v", err)
        }

        // Verificar código de status HTTP
        if resp.StatusCode != http.StatusOK {
                return nil, fmt.Errorf("código de status HTTP inesperado: %d - %s", resp.StatusCode, string(body))
        }

        // Converter resposta para mapa
        var resultado map[string]interface{}
        err = json.Unmarshal(body, &resultado)
        if err != nil {
                return nil, fmt.Errorf("erro ao converter resposta: %v", err)
        }

        return resultado, nil
}

// Funções auxiliares para conversão de tipos
func getString(m map[string]interface{}, key string) string {
        if v, ok := m[key]; ok {
                if s, ok := v.(string); ok {
                        return s
                }
        }
        return ""
}

func parseInt(s string) int {
        var i int
        _, err := fmt.Sscanf(s, "%d", &i)
        if err != nil {
                return 0
        }
        return i
}