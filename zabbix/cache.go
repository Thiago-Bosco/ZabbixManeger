
package zabbix

import (
	"sync"
	"time"
)

type Cache struct {
	dados map[string]interface{}
	exp   map[string]time.Time
	mu    sync.RWMutex
}

func NovoCache() *Cache {
	return &Cache{
		dados: make(map[string]interface{}),
		exp:   make(map[string]time.Time),
	}
}

func (c *Cache) Set(chave string, valor interface{}, duracao time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.dados[chave] = valor
	c.exp[chave] = time.Now().Add(duracao)
}

func (c *Cache) Get(chave string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	if exp, ok := c.exp[chave]; ok {
		if time.Now().After(exp) {
			delete(c.dados, chave)
			delete(c.exp, chave)
			return nil, false
		}
		return c.dados[chave], true
	}
	return nil, false
}
