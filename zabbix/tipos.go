package zabbix

// Item representa um item monitorado no Zabbix
type Item struct {
	ID   string `json:"itemid"`
	Nome string `json:"name"`
}

// Trigger representa um gatilho/alerta no Zabbix
type Trigger struct {
	ID   string `json:"triggerid"`
	Nome string `json:"description"`
}

// Host representa um host no Zabbix
type Host struct {
	ID       string    `json:"hostid"`
	Nome     string    `json:"host"`
	Status   string    `json:"status"`
	Items    []Item    `json:"items"`
	Triggers []Trigger `json:"triggers"`
}

// Resposta representa a estrutura da resposta da API do Zabbix
type Resposta struct {
	Resultados []Host `json:"result"`
}

// ConfigAPI armazena as configurações de conexão com a API Zabbix
type ConfigAPI struct {
	URL         string
	Token       string
	TempoLimite int // em segundos
}

// StatusHost mapeia os valores de status para textos em português
var StatusHost = map[string]string{
	"0": "Ativo",
	"1": "Inativo",
}
