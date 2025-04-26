
package metrics

import (
	"sync"
	"time"
)

type Metrics struct {
	contadores map[string]int64
	timings    map[string][]time.Duration
	mu         sync.RWMutex
}

func Novo() *Metrics {
	return &Metrics{
		contadores: make(map[string]int64),
		timings:    make(map[string][]time.Duration),
	}
}

func (m *Metrics) IncrementarContador(nome string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.contadores[nome]++
}

func (m *Metrics) RegistrarTempo(nome string, duracao time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.timings[nome] = append(m.timings[nome], duracao)
}
