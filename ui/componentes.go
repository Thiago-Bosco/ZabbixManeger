package ui

import (
	"fmt"
	"image/color"
	"io/ioutil"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

// CriarImagemDaURI carrega uma imagem a partir de uma URI
func CriarImagemDaURI(uri fyne.URI) (*canvas.Image, error) {
	// Abrir o arquivo de imagem
	r, err := storage.Reader(uri)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir imagem: %w", err)
	}
	defer r.Close()

	// Ler o conteúdo do arquivo
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler imagem: %w", err)
	}

	// Criar um arquivo temporário
	tempFile, err := ioutil.TempFile("", "zabbix-manager-*.png")
	if err != nil {
		return nil, fmt.Errorf("erro ao criar arquivo temporário: %w", err)
	}
	defer tempFile.Close()

	// Escrever os dados no arquivo temporário
	if _, err := tempFile.Write(data); err != nil {
		return nil, fmt.Errorf("erro ao escrever arquivo temporário: %w", err)
	}

	// Criar a imagem
	img := canvas.NewImageFromFile(tempFile.Name())
	img.FillMode = canvas.ImageFillContain

	return img, nil
}

// CriarImagemLocal carrega uma imagem a partir de um caminho local
func CriarImagemLocal(caminho string) (*canvas.Image, error) {
	// Verificar se o arquivo existe
	if _, err := os.Stat(caminho); os.IsNotExist(err) {
		return nil, fmt.Errorf("arquivo de imagem não encontrado: %s", caminho)
	}

	// Criar a imagem
	img := canvas.NewImageFromFile(caminho)
	img.FillMode = canvas.ImageFillContain

	return img, nil
}

// CriarCorDoTema retorna uma cor do tema Zabbix
func CriarCorDoTema(nome string) color.Color {
	// Cores do tema Zabbix (baseadas nas cores da interface web)
	cores := map[string]color.RGBA{
		"vermelho":    {R: 218, G: 68, B: 83, A: 255},   // #DA4453
		"vermelho-escuro": {R: 165, G: 42, B: 42, A: 255},  // #A52A2A
		"cinza":       {R: 170, G: 170, B: 170, A: 255}, // #AAAAAA
		"cinza-escuro": {R: 85, G: 85, B: 85, A: 255},   // #555555
		"branco":      {R: 255, G: 255, B: 255, A: 255}, // #FFFFFF
		"preto":       {R: 0, G: 0, B: 0, A: 255},       // #000000
	}

	// Retornar a cor escolhida ou preto como padrão
	if cor, ok := cores[nome]; ok {
		return cor
	}
	return cores["preto"]
}

// CriarBotaoPrimario cria um botão com estilo primário
func CriarBotaoPrimario(texto string, acao func()) *widget.Button {
	botao := widget.NewButton(texto, acao)
	botao.Importance = widget.HighImportance
	return botao
}

// CriarBotaoSecundario cria um botão com estilo secundário
func CriarBotaoSecundario(texto string, acao func()) *widget.Button {
	botao := widget.NewButton(texto, acao)
	botao.Importance = widget.MediumImportance
	return botao
}

// CriarDiretorioSeNecessario cria um diretório se ele não existir
func CriarDiretorioSeNecessario(caminho string) error {
	return os.MkdirAll(caminho, 0755)
}

// CriarLogoDaAPI gera uma imagem do logo do Zabbix a partir da API
func CriarLogoDaAPI(urlAPI string) (*canvas.Image, error) {
	// URL do logo do Zabbix
	caminhoLogo := "assets/logo.png"

	// Verificar se o arquivo existe
	if _, err := os.Stat(caminhoLogo); os.IsNotExist(err) {
		// Criar diretório de assets se não existir
		err := CriarDiretorioSeNecessario(filepath.Dir(caminhoLogo))
		if err != nil {
			return nil, fmt.Errorf("erro ao criar diretório de assets: %w", err)
		}

		// Criar um logo padrão
		logoData := []byte(`<svg width="200" height="50" xmlns="http://www.w3.org/2000/svg">
			<rect width="200" height="50" fill="#DA4453"/>
			<text x="50%" y="50%" font-family="Arial" font-size="20" fill="white" text-anchor="middle" dominant-baseline="middle">ZABBIX</text>
		</svg>`)

		// Salvar o logo
		err = ioutil.WriteFile(caminhoLogo, logoData, 0644)
		if err != nil {
			return nil, fmt.Errorf("erro ao salvar logo: %w", err)
		}
	}

	// Carregar a imagem
	return CriarImagemLocal(caminhoLogo)
}