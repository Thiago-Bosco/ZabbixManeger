
package zabbix

// Item representa um item de monitoramento do Zabbix
type Item struct {
	ID          string `json:"itemid"`
	Nome        string `json:"name"`
	Status      string `json:"status"`
	Estado      string `json:"state"`
	UltimoValor string `json:"lastvalue"`
}

// Trigger representa uma trigger (alarme) do Zabbix
type Trigger struct {
	ID              string `json:"triggerid"`
	Nome            string `json:"description"`
	Status          string `json:"status"`
	Valor           string `json:"value"`
	Prioridade      string `json:"priority"`
	UltimaAlteracao string `json:"lastchange"`
}

// Host representa um host do Zabbix com seus itens e triggers
type Host struct {
	ID         string      `json:"hostid"`
	Nome       string      `json:"host"`
	Status     string      `json:"status"`
	Items      []Item      `json:"items"`
	Triggers   []Trigger   `json:"triggers"`
	Interfaces []Interface `json:"interfaces"`
}

// Interface representa uma interface de rede do host
type Interface struct {
	ID    string `json:"interfaceid"`
	Tipo  string `json:"type"`
	IP    string `json:"ip"`
	DNS   string `json:"dns"`
	Porta string `json:"port"`
}
