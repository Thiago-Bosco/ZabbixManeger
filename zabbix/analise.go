
package zabbix

import "time"

type AnaliseProblema struct {
	HostID           string
	HostNome         string
	TotalProblemas   int
	LimitesExcedidos int
	PicoTrigger      struct {
		Nome      string
		DataPico  time.Time
		Contagem  int
		Gravidade string
	}
	ProblemasPorTrigger map[string]int
	TempoIndisponivel   time.Duration
}
