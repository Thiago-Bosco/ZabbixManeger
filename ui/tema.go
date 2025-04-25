package ui

import (
        "image/color"

        "fyne.io/fyne/v2"
        "fyne.io/fyne/v2/theme"
)

// TemaZabbix é um tema personalizado para a aplicação
type TemaZabbix struct{}

// NovoTemaZabbix cria uma nova instância do tema Zabbix
func NovoTemaZabbix() fyne.Theme {
        return &TemaZabbix{}
}

// Color retorna a cor para o nome especificado
func (t *TemaZabbix) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
        switch name {
        case theme.ColorNamePrimary:
                return CorPrincipal
        case theme.ColorNameForeground:
                return CorTexto
        case theme.ColorNameBackground:
                return CorFundo
        case theme.ColorNameButton:
                return CorSecundaria
        case theme.ColorNameDisabled:
                return color.NRGBA{R: 180, G: 180, B: 180, A: 128}
        case theme.ColorNamePlaceHolder:
                return color.NRGBA{R: 150, G: 150, B: 150, A: 255}
        default:
                return theme.DefaultTheme().Color(name, variant)
        }
}

// Icon retorna o recurso do ícone para o nome especificado
func (t *TemaZabbix) Icon(name fyne.ThemeIconName) fyne.Resource {
        return theme.DefaultTheme().Icon(name)
}

// Font retorna o nome da fonte para o estilo especificado
func (t *TemaZabbix) Font(style fyne.TextStyle) fyne.Resource {
        return theme.DefaultTheme().Font(style)
}

// Size retorna o tamanho para o nome especificado
func (t *TemaZabbix) Size(name fyne.ThemeSizeName) float32 {
        switch name {
        case theme.SizeNamePadding:
                return 4
        case theme.SizeNameInlineIcon:
                return 20
        case theme.SizeNameScrollBar:
                return 10
        case theme.SizeNameScrollBarSmall:
                return 5
        case theme.SizeNameText:
                return 12
        case theme.SizeNameHeadingText:
                return 16
        case theme.SizeNameSubHeadingText:
                return 14
        case theme.SizeNameInputBorder:
                return 1
        default:
                return theme.DefaultTheme().Size(name)
        }
}