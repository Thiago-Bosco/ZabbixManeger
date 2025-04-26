package ui

import (
	"image/color"
)

// Tema representa o tema visual da aplicação
type Tema struct {
	// Cores
	CorPrimaria       color.RGBA
	CorSecundaria     color.RGBA
	CorTerciaria      color.RGBA
	CorFundo          color.RGBA
	CorTexto          color.RGBA
	CorTextoPrimaria  color.RGBA
	CorTextoSecundaria color.RGBA
	CorBorda          color.RGBA
	CorErro           color.RGBA
	CorSucesso        color.RGBA
	CorAviso          color.RGBA
	CorInfo           color.RGBA
	
	// Tamanhos
	TamanhoFonteNormal  float64
	TamanhoFontePequena float64
	TamanhoFonteGrande  float64
	TamanhoFonteTitulo  float64
	EspacamentoPadrao   float64
	BorderaPadrao       float64
}

// TemaPadrao retorna o tema padrão da aplicação
func TemaPadrao() *Tema {
	return &Tema{
		// Cores inspiradas no Zabbix (vermelho e cinza)
		CorPrimaria:       color.RGBA{R: 215, G: 0, B: 0, A: 255},       // Vermelho Zabbix
		CorSecundaria:     color.RGBA{R: 100, G: 100, B: 100, A: 255},   // Cinza escuro
		CorTerciaria:      color.RGBA{R: 50, G: 50, B: 50, A: 255},      // Cinza mais escuro
		CorFundo:          color.RGBA{R: 240, G: 240, B: 240, A: 255},   // Cinza claro para fundo
		CorTexto:          color.RGBA{R: 33, G: 33, B: 33, A: 255},      // Quase preto
		CorTextoPrimaria:  color.RGBA{R: 215, G: 0, B: 0, A: 255},       // Vermelho Zabbix
		CorTextoSecundaria: color.RGBA{R: 100, G: 100, B: 100, A: 255},  // Cinza escuro
		CorBorda:          color.RGBA{R: 200, G: 200, B: 200, A: 255},   // Cinza médio
		CorErro:           color.RGBA{R: 200, G: 0, B: 0, A: 255},       // Vermelho
		CorSucesso:        color.RGBA{R: 0, G: 150, B: 0, A: 255},       // Verde
		CorAviso:          color.RGBA{R: 255, G: 150, B: 0, A: 255},     // Laranja
		CorInfo:           color.RGBA{R: 0, G: 100, B: 200, A: 255},     // Azul
		
		// Tamanhos
		TamanhoFonteNormal:  16.0,
		TamanhoFontePequena: 12.0,
		TamanhoFonteGrande:  20.0,
		TamanhoFonteTitulo:  24.0,
		EspacamentoPadrao:   8.0,
		BorderaPadrao:       1.0,
	}
}

// TemaEscuro retorna um tema escuro para a aplicação
func TemaEscuro() *Tema {
	return &Tema{
		// Cores inspiradas no Zabbix (vermelho e cinza) com fundo escuro
		CorPrimaria:       color.RGBA{R: 215, G: 30, B: 30, A: 255},     // Vermelho Zabbix
		CorSecundaria:     color.RGBA{R: 150, G: 150, B: 150, A: 255},   // Cinza claro
		CorTerciaria:      color.RGBA{R: 100, G: 100, B: 100, A: 255},   // Cinza médio
		CorFundo:          color.RGBA{R: 30, G: 30, B: 30, A: 255},      // Cinza muito escuro para fundo
		CorTexto:          color.RGBA{R: 220, G: 220, B: 220, A: 255},   // Quase branco
		CorTextoPrimaria:  color.RGBA{R: 255, G: 80, B: 80, A: 255},     // Vermelho claro
		CorTextoSecundaria: color.RGBA{R: 180, G: 180, B: 180, A: 255},  // Cinza claro
		CorBorda:          color.RGBA{R: 70, G: 70, B: 70, A: 255},      // Cinza escuro
		CorErro:           color.RGBA{R: 255, G: 70, B: 70, A: 255},     // Vermelho claro
		CorSucesso:        color.RGBA{R: 70, G: 200, B: 70, A: 255},     // Verde claro
		CorAviso:          color.RGBA{R: 255, G: 180, B: 50, A: 255},    // Laranja claro
		CorInfo:           color.RGBA{R: 70, G: 150, B: 255, A: 255},    // Azul claro
		
		// Tamanhos
		TamanhoFonteNormal:  16.0,
		TamanhoFontePequena: 12.0,
		TamanhoFonteGrande:  20.0,
		TamanhoFonteTitulo:  24.0,
		EspacamentoPadrao:   8.0,
		BorderaPadrao:       1.0,
	}
}