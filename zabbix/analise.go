package zabbix

import "time"

type AnaliseProblema struct {
	HostNome        string
	TotalProblemas  int
	LimitesExcedidos int
	PicoTrigger     struct {
		Nome      string
		DataPico  time.Time
		Contagem  int
		Gravidade string
	}
	DuracaoMedia string
	TempoTotal   string
}