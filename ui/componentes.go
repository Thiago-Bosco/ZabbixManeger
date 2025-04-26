package ui

// Componente representa um componente da interface
type Componente interface {
	// Renderizar renderiza o componente
	Renderizar() error
}

// Botao representa um botão na interface
type Botao struct {
	Texto     string
	OnClick   func() error
	Habilitado bool
}

// Campo representa um campo de entrada na interface
type Campo struct {
	Rotulo    string
	Valor     string
	Dica      string
	Senha     bool
	MultiLinha bool
	OnChange  func(string) error
	Habilitado bool
}

// Lista representa uma lista na interface
type Lista struct {
	Itens     []string
	Selecionado int
	OnSelect  func(int) error
	Habilitado bool
}

// Mensagem representa uma mensagem na interface
type Mensagem struct {
	Texto     string
	Tipo      TipoMensagem
}

// TipoMensagem representa o tipo de uma mensagem
type TipoMensagem int

// Tipos de mensagem
const (
	MensagemInfo TipoMensagem = iota
	MensagemErro
	MensagemSucesso
	MensagemAviso
)