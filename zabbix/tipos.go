package zabbix

import (
        "encoding/json"
)

// ConfigAPI representa a configuração para a API do Zabbix
type ConfigAPI struct {
        URL         string
        Token       string
        TempoLimite int
}

// Resposta representa a resposta da API do Zabbix
type Resposta struct {
        Jsonrpc string  `json:"jsonrpc"`
        Result  []Host  `json:"result"`
        ID      int     `json:"id"`
        
        // Para lidar com erros da API
        Error   *struct {
                Code    int    `json:"code"`
                Message string `json:"message"`
                Data    string `json:"data"`
        } `json:"error"`
        
        // Para acesso fácil aos resultados
        Resultados []Host
}

// UnmarshalJSON é um método customizado para desserialização que popula o campo Resultados
func (r *Resposta) UnmarshalJSON(data []byte) error {
        // Criar uma estrutura temporária para evitar recursão infinita
        type TempResposta Resposta
        var temp TempResposta
        
        if err := json.Unmarshal(data, &temp); err != nil {
                return err
        }
        
        // Copiar os dados para a resposta original
        *r = Resposta(temp)
        
        // Atribuir o mesmo valor para resultados para facilitar o acesso
        r.Resultados = r.Result
        
        return nil
}

// Host representa um host do Zabbix
type Host struct {
        ID        string    `json:"hostid"`
        Nome      string    `json:"host"`
        Status    string    `json:"status"`
        Items     []Item    `json:"items"`
        Triggers  []Trigger `json:"triggers"`
}

// Item representa um item de monitoramento do Zabbix
type Item struct {
        ID        string `json:"itemid"`
        Nome      string `json:"name"`
}

// Trigger representa uma trigger do Zabbix
type Trigger struct {
        ID        string `json:"triggerid"`
        Nome      string `json:"description"`
}

// StatusHost mapeia o status do host para uma string legível
var StatusHost = map[string]string{
        "0": "Ativo",
        "1": "Inativo",
}