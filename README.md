# Zabbix Manager

Uma aplicação desktop multiplataforma (Windows e Linux) que utiliza o código Go para conectar à API do Zabbix, visualizar hosts, itens e triggers, e gerar relatórios em CSV.

## Funcionalidades

- Interface desktop amigável e nativa
- Gerenciamento de múltiplos servidores Zabbix
- Visualização de hosts e seus itens de monitoramento
- Busca e filtragem de hosts
- Exportação de relatórios em formato CSV
- Suporte a autenticação via token da API Zabbix

## Requisitos

- Windows 7 ou superior / Linux
- Acesso a um servidor Zabbix (URL e token da API)

## Instalação

1. Faça o download do executável mais recente para seu sistema operacional
2. Execute o arquivo `zabbix-manager` (Linux) ou `zabbix-manager.exe` (Windows)

## Compilação

### Ambiente de Desenvolvimento

Para compilar a aplicação, você precisa ter o Go instalado (versão 1.18 ou superior).

```bash
# Clone o repositório
git clone https://github.com/seu-usuario/zabbix-manager.git
cd zabbix-manager

# Instale as dependências
go mod tidy

# Compile para seu sistema atual
go build
```

### Compilação para Windows (a partir de qualquer sistema)

```bash
GOOS=windows GOARCH=amd64 go build -o zabbix-manager.exe
```

### Compilação para Linux (a partir de qualquer sistema)

```bash
GOOS=linux GOARCH=amd64 go build -o zabbix-manager
```

### Interface Gráfica vs. Modo Headless

O projeto suporta dois modos de operação através de tags de compilação:

1. **Modo GUI (padrão)** - Interface gráfica nativa com Gio
   ```bash
   go build
   # ou explicitamente
   go build -tags=""
   ```

2. **Modo Headless** - Interface de linha de comando
   ```bash
   go build -tags=headless
   # ou para executar diretamente
   go run -tags=headless .
   ```

## Estrutura do Projeto

- `/zabbix`: Código de integração com a API do Zabbix
- `/ui`: Componentes da interface gráfica
  - `/ui/gio`: Interface gráfica nativa com Gio
- `/config`: Gerenciamento de configurações
- `/assets`: Recursos visuais (imagens, ícones)

## Tecnologias Utilizadas

- **Go**: Linguagem principal do projeto
- **Gio**: Framework de interface gráfica nativa
- **Go modules**: Gerenciamento de dependências

## Configuração

A aplicação armazena suas configurações em:

- Windows: `%USERPROFILE%\.zabbix-manager\config.json`
- Linux: `$HOME/.zabbix-manager/config.json`

Os relatórios são exportados para:
- Windows: `%USERPROFILE%\Relatórios Zabbix\`
- Linux: `$HOME/Relatórios Zabbix/`

## Benefícios da Abordagem Multiplataforma

- **Código Único**: A mesma base de código funciona em Windows e Linux
- **Interface Nativa**: Interface gráfica nativa em ambos os sistemas
- **Modo Headless**: Opção de usar em servidores sem interface gráfica
- **Desempenho**: Aplicação compilada com excelente desempenho

## Licença

Este projeto está licenciado sob a licença MIT - veja o arquivo [LICENSE](LICENSE) para detalhes.