
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

func (c *ClienteAPI) AnalisarProblemasMensais(ano, mes int) ([]AnaliseProblema, error) {
	// Implementação básica para teste
	analises := []AnaliseProblema{
		{
			HostNome:       "Servidor 1",
			TotalProblemas: 10,
			LimitesExcedidos: 2,
			PicoTrigger: struct {
				Nome      string
				DataPico  time.Time
				Contagem  int
				Gravidade string
			}{
				Nome:      "CPU Alto",
				DataPico:  time.Now(),
				Contagem:  5,
				Gravidade: "4",
			},
			DuracaoMedia: "2h 30min",
			TempoTotal:   "25h",
		},
	}
	return analises, nil
}
