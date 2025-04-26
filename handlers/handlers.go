
package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"zabbix-manager/config"
	"zabbix-manager/zabbix"
)

type Handler struct {
	Config      *config.Configuração
	ClienteAPI  *zabbix.ClienteAPI
	RenderTemplate func(w http.ResponseWriter, nome string, dados interface{})
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if h.ClienteAPI == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	http.Redirect(w, r, "/hosts", http.StatusFound)
}

func (h *Handler) Analise(w http.ResponseWriter, r *http.Request) {
	if h.ClienteAPI == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	ano := time.Now().Year()
	mes := int(time.Now().Month())

	if anoStr := r.URL.Query().Get("ano"); anoStr != "" {
		if anoInt, err := strconv.Atoi(anoStr); err == nil {
			ano = anoInt
		}
	}
	if mesStr := r.URL.Query().Get("mes"); mesStr != "" {
		if mesInt, err := strconv.Atoi(mesStr); err == nil {
			mes = mesInt
		}
	}

	analises, err := h.ClienteAPI.AnalisarProblemasMensais(ano, mes)
	if err != nil {
		h.RenderTemplate(w, "analise", map[string]interface{}{
			"Erro": fmt.Sprintf("Erro ao analisar problemas: %v", err),
		})
		return
	}

	dados := map[string]interface{}{
		"Analises":       analises,
		"AnoSelecionado": ano,
		"MesSelecionado": mes,
		"Anos":          []int{ano - 1, ano, ano + 1},
		"Meses":         []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
		"NomesMeses":    []string{"Janeiro", "Fevereiro", "Março", "Abril", "Maio", "Junho", "Julho", "Agosto", "Setembro", "Outubro", "Novembro", "Dezembro"},
		"TipoFiltro":    "mensal",
	}

	h.RenderTemplate(w, "analise", dados)
}
