package zabbix

// Host representa um host do Zabbix
type Host struct {
	ID       string
	Nome     string
	Status   int
	Items    []Item
	Triggers []Trigger
}

// Item representa um item de monitoramento do Zabbix
type Item struct {
	ID     string
	Nome   string
	Chave  string
	Status int
}

// Trigger representa um gatilho do Zabbix
type Trigger struct {
	ID         string
	Descricao  string
	Status     int
	Prioridade int
}

// StatusHost define os possíveis status de um host
var StatusHost = map[int]string{
	0: "Ativo",
	1: "Desativado",
}

// StatusItem define os possíveis status de um item
var StatusItem = map[int]string{
	0: "Ativo",
	1: "Desativado",
}

// StatusTrigger define os possíveis status de um trigger
var StatusTrigger = map[int]string{
	0: "Ativo",
	1: "Desativado",
}

// PrioridadeTrigger define as possíveis prioridades de um trigger
var PrioridadeTrigger = map[int]string{
	0: "Não classificado",
	1: "Informação",
	2: "Aviso",
	3: "Médio",
	4: "Alto",
	5: "Crítico",
}