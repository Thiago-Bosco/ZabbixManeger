
package zabbix

type ErroAPI struct {
	Codigo    int
	Mensagem  string
	Detalhes  string
}

func (e *ErroAPI) Error() string {
	return e.Mensagem
}

func NovoErroAPI(codigo int, mensagem, detalhes string) *ErroAPI {
	return &ErroAPI{
		Codigo:   codigo,
		Mensagem: mensagem,
		Detalhes: detalhes,
	}
}
