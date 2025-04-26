
package zabbix

import (
	"encoding/json"
	"time"
)

type Item struct {
	ID              string `json:"itemid"`
	Nome            string `json:"name"`
	Status          string `json:"status"`
	Estado          string `json:"state"`
	UltimoValor     string `json:"lastvalue"`
	UltimaAlteracao string `json:"lastchange"`
}

type Trigger struct {
	ID              string `json:"triggerid"`
	Nome            string `json:"description"`
	Status          string `json:"status"`
	Valor           string `json:"value"`
	Prioridade      string `json:"priority"`
	UltimaAlteracao string `json:"lastchange"`
}

type Host struct {
	ID         string      `json:"hostid"`
	Nome       string      `json:"host"`
	Status     string      `json:"status"`
	Items      []Item      `json:"items"`
	Triggers   []Trigger   `json:"triggers"`
	Interfaces []Interface `json:"interfaces"`
}

type Interface struct {
	ID    string `json:"interfaceid"`
	Tipo  string `json:"type"`
	IP    string `json:"ip"`
	DNS   string `json:"dns"`
	Porta string `json:"port"`
}

type Problema struct {
	ID          string    `json:"eventid"`
	Nome        string    `json:"name"`
	Severidade  string    `json:"severity"`
	DataInicio  time.Time `json:"clock"`
	DataFim     time.Time `json:"r_clock"`
	Duracao     string    `json:"duration"`
	HostID      string    `json:"hostid"`
	TriggerID   string    `json:"triggerid"`
	Valor       string    `json:"value"`
	Hosts       []Host    `json:"hosts"`
}

type Evento struct {
	ID             string          `json:"eventid"`
	Nome           string          `json:"name"`
	Clock          string          `json:"clock"`
	Valor          string          `json:"value"`
	Severidade     string          `json:"severity"`
	Reconhecido    string          `json:"acknowledged"`
	HostID         string          `json:"hostid"`
	ObjetoID       string          `json:"objectid"`
	TipoObjeto     string          `json:"object"`
	ObjetoRelativo json.RawMessage `json:"relatedObject"`
}

type AnaliseMensal struct {
	HostID              string
	HostNome            string
	TotalProblemas      int
	ProblemasPorTrigger map[string]int
	TempoIndisponivel   time.Duration
	LimitesExcedidos    int
	PicoTrigger         struct {
		Nome      string
		DataPico  time.Time
		Contagem  int
		Gravidade string
	}
}
