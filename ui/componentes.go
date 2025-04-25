package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/storage"
	"image/color"
)

// Cores da aplicação
var (
	CorPrincipal = color.NRGBA{R: 215, G: 43, B: 43, A: 255}  // Vermelho Zabbix
	CorSecundaria = color.NRGBA{R: 68, G: 68, B: 68, A: 255}  // Cinza escuro
	CorFundo = color.NRGBA{R: 250, G: 250, B: 250, A: 255}    // Quase branco
	CorTexto = color.NRGBA{R: 33, G: 33, B: 33, A: 255}       // Quase preto
)

// CarregarLogoZabbix carrega o logo do Zabbix
func CarregarLogoZabbix() fyne.Resource {
	// Retornar o recurso SVG embutido
	recurso, _ := fyne.LoadResourceFromPath("./assets/zabbix_logo.svg")
	if recurso == nil {
		// Fallback para um recurso estático
		recurso = logoZabbixResource
	}
	return recurso
}

// CriarElementoComCor cria um elemento com a cor especificada
func CriarElementoComCor(cor color.Color) *canvas.Rectangle {
	retangulo := canvas.NewRectangle(cor)
	retangulo.SetMinSize(fyne.NewSize(20, 20))
	return retangulo
}

// Uma variável para evitar ter que carregar o recurso várias vezes
var logoZabbixResource = fyne.NewStaticResource("zabbix_logo", []byte(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 202 60" width="64" height="64">
  <path fill="#D72A1F" d="M126.333,18.624l-28.49,27.032L121.147,40.4L93.761,62.669l28.891-25.343l-23.493,4.061Z"/>
  <path fill="#222222" d="M53.237,35.668V18.624h21.592v4.822H58.483v3.381h14.705v4.823H58.483v4.018Z"/>
  <path fill="#222222" d="M137.961,35.668V18.624h5.247V35.668Z"/>
  <path fill="#222222" d="M6.433,35.668,0,18.624H5.82l3.833,11.025l3.787-11.025h5.867L13.181,35.668Z"/>
  <path fill="#222222" d="M21.692,35.668V18.624h21.592v4.822H26.939v3.381h14.705v4.823H26.939v4.018Z"/>
  <path fill="#222222" d="M94.688,35.668V18.624h21.592v4.822H99.935v3.381h14.705v4.823H99.935v4.018Z"/>
  <path fill="#222222" d="M148.186,35.668V18.624h12.812a11.024,11.024,0,0,1,3.69.532,5.88,5.88,0,0,1,2.462,1.549,3.648,3.648,0,0,1,.867,2.461,3.993,3.993,0,0,1-1.409,3.137,7.826,7.826,0,0,1-4.205,1.6,5.906,5.906,0,0,1,1.77.867,10.6,10.6,0,0,1,1.725,1.8l2.744,4.963h-6.289l-3.271-5.291a4.831,4.831,0,0,0-1.454-1.6,3.607,3.607,0,0,0-1.707-.327h-.673v7.354Zm5.246-11.464h4.018a6.769,6.769,0,0,0,1.435-.187,1.557,1.557,0,0,0,.962-.61,1.761,1.761,0,0,0,.337-1.1,1.571,1.571,0,0,0-.542-1.268,2.932,2.932,0,0,0-1.931-.458h-4.279Z"/>
  <path fill="#222222" d="M173.649,35.668V18.624h5.246V30.844h13.159v4.824Z"/>
  <path fill="#222222" d="M77.509,35.668V18.624h5.246V35.668Z"/>
</svg>
`))

// Estrutura para implementar fyne.URIWriteCloser para salvamento de arquivos
type storage struct {
	fyne.URI
}

func (s *storage) Write(p []byte) (n int, err error) {
	return 0, nil
}

func (s *storage) Close() error {
	return nil
}

// NewExtensionFileFilter cria um filtro de arquivo por extensão
func NewExtensionFileFilter(extensions []string) storage.FileFilter {
	return &extensionFileFilter{extensions}
}

type extensionFileFilter struct {
	extensions []string
}

func (e *extensionFileFilter) Matches(uri fyne.URI) bool {
	if uri == nil || uri.String() == "" {
		return true
	}

	path := uri.Path()
	if path == "" {
		return true
	}

	for _, ext := range e.extensions {
		if filepath.Ext(path) == ext {
			return true
		}
	}
	return false
}

func (e *extensionFileFilter) String() string {
	var builder strings.Builder
	for i, ext := range e.extensions {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString("*")
		builder.WriteString(ext)
	}
	return builder.String()
}
