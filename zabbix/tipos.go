package zabbix

// Host representa um host no Zabbix
type Host struct {
	ID       string    `json:"hostid"`
	Nome     string    `json:"name"`
	Status   string    `json:"status"`
	Items    []Item    `json:"items"`
	Triggers []Trigger `json:"triggers"`
}

// Item representa um item de monitoramento no Zabbix
type Item struct {
	ID     string `json:"itemid"`
	Nome   string `json:"name"`
	Chave  string `json:"key_"`
	Status string `json:"status"`
}

// Trigger representa um trigger (alerta) no Zabbix
type Trigger struct {
	ID         string `json:"triggerid"`
	Descricao  string `json:"description"`
	Status     string `json:"status"`
	Prioridade string `json:"priority"`
}

// StatusHost mapeia os códigos de status dos hosts para textos
var StatusHost = map[string]string{
	"0": "Ativo",
	"1": "Inativo",
}

// StatusItem mapeia os códigos de status dos itens para textos
var StatusItem = map[string]string{
	"0": "Ativo",
	"1": "Inativo",
	"2": "Não suportado",
}

// StatusTrigger mapeia os códigos de status dos triggers para textos
var StatusTrigger = map[string]string{
	"0": "Ativo",
	"1": "Inativo",
}

// PrioridadeTrigger mapeia os códigos de prioridade dos triggers para textos
var PrioridadeTrigger = map[string]string{
	"0": "Não classificada",
	"1": "Informação",
	"2": "Atenção",
	"3": "Média",
	"4": "Alta",
	"5": "Desastre",
}