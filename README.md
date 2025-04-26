# Zabbix Manager Web

Uma aplicação web desenvolvida em Go para gerenciamento e monitoramento de infraestrutura Zabbix, com foco em visualização de hosts, itens e triggers, e geração de relatórios em CSV.

## Características Principais

- Interface web responsiva e amigável
- Suporte para múltiplos perfis de servidor Zabbix
- Visualização e busca de hosts monitorados
- Visualização de itens e triggers de cada host
- Exportação de relatórios em formato CSV
- Implementação em Go para desempenho e eficiência
- Sem necessidade de banco de dados adicional

## Requisitos

- Go 1.18 ou superior
- Acesso a um servidor Zabbix com API ativa
- Token de API gerado no frontend do Zabbix

## Instalação

### Clonar o repositório

```bash
git clone https://github.com/seu-usuario/zabbix-manager.git
cd zabbix-manager
```

### Compilar e executar

```bash
# Compilar
go build

# Executar
./zabbix-manager
```

O servidor iniciará na porta 5000. Acesse no navegador: http://localhost:5000

## Configuração

Na primeira execução, o sistema solicitará a configuração de um servidor Zabbix:

1. Acesse a interface web em http://localhost:5000
2. Clique em "Adicionar Servidor"
3. Preencha:
   - Nome: Um identificador para o servidor (Ex: "Zabbix Produção")
   - URL da API: URL completa do endpoint da API (Ex: "https://zabbix.exemplo.com/api_jsonrpc.php")
   - Token da API: Token de autenticação gerado no frontend do Zabbix

## Como obter um token da API Zabbix

1. Faça login no frontend do Zabbix com um usuário que tenha permissões adequadas
2. Acesse: Administração > Usuários
3. Selecione seu usuário > guia API tokens
4. Clique em "Criar token da API"
5. Dê um nome ao token e defina uma data de expiração (opcional)
6. Copie o token gerado e use-o na configuração do Zabbix Manager

## Estrutura do Projeto

- `main.go`: Ponto de entrada da aplicação web
- `zabbix/`: Pacote com implementação da API do Zabbix
  - `api.go`: Cliente para API do Zabbix
  - `relatorios.go`: Geração de relatórios CSV
  - `tipos.go`: Definições de tipos utilizados
- `config/`: Configurações da aplicação
  - `config.go`: Gerenciamento de configurações
- `templates/`: Templates HTML
- `static/`: Arquivos estáticos (CSS, JS)

## Licença

Este projeto está licenciado sob a licença MIT - veja o arquivo LICENSE para mais detalhes.