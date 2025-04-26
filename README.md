# Zabbix Manager

Uma aplicação desktop para Windows que utiliza o código Go para conectar à API do Zabbix, visualizar hosts, itens e triggers, e gerar relatórios em CSV.

## Funcionalidades

- Interface desktop amigável
- Gerenciamento de múltiplos servidores Zabbix
- Visualização de hosts e seus itens de monitoramento
- Busca e filtragem de hosts
- Exportação de relatórios em formato CSV
- Suporte a autenticação via token da API Zabbix

## Requisitos

- Windows 7 ou superior
- Acesso a um servidor Zabbix (URL e token da API)

## Instalação

1. Faça o download do executável mais recente na seção de [Releases](https://github.com/seu-usuario/zabbix-manager/releases)
2. Execute o arquivo `zabbix-manager.exe`

## Compilação

### Ambiente de Desenvolvimento

Para compilar a aplicação, você precisa ter o Go instalado (versão 1.18 ou superior).

```
go get -u github.com/seu-usuario/zabbix-manager
cd zabbix-manager
go build
```

### Compilação para Windows

```
GOOS=windows GOARCH=amd64 go build -o zabbix-manager.exe
```

### Modo Headless (sem GUI)

Para executar em modo console (sem interface gráfica):

```
go run -tags=headless .
```

## Estrutura do Projeto

- `/zabbix`: Código de integração com a API do Zabbix
- `/ui`: Componentes da interface gráfica
- `/config`: Gerenciamento de configurações
- `/assets`: Recursos visuais (imagens, ícones)

## Configuração

A aplicação armazena suas configurações em:

- Windows: `%USERPROFILE%\.zabbix-manager\config.json`

## Licença

Este projeto está licenciado sob a licença MIT - veja o arquivo [LICENSE](LICENSE) para detalhes.

## Autor

Seu Nome