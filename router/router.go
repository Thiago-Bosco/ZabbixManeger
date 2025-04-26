
package router

import (
	"net/http"
	"zabbix-manager/handlers"
)

func ConfigureRoutes(h *handlers.Handler) {
	http.HandleFunc("/", h.Home)
	http.HandleFunc("/login", h.Login)
	http.HandleFunc("/config", h.Config)
	http.HandleFunc("/analise", h.Analise)
	http.HandleFunc("/hosts", h.Hosts)
	http.HandleFunc("/hosts/buscar", h.BuscarHosts)
	http.HandleFunc("/exportar", h.ExportarCSV)
	
	// Servir arquivos estáticos
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
}
