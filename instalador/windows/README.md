# Instruções para Criar o Instalador do Zabbix Manager para Windows

Este diretório contém os arquivos necessários para criar um instalador do Windows para o Zabbix Manager.

## Requisitos

- [NSIS (Nullsoft Scriptable Install System)](https://nsis.sourceforge.io/Download)
- Executável compilado do Zabbix Manager (`zabbix-manager.exe`)

## Passos para Criar o Instalador

1. Instale o NSIS em seu computador Windows.
2. Compile o aplicativo Zabbix Manager:
   ```
   go build -o zabbix-manager.exe
   ```
3. Execute o script `build-installer.bat` neste diretório.
4. Se tudo correr bem, o instalador será criado como `ZabbixManager-Setup.exe`.

## Conteúdo dos Arquivos

- `zabbix-manager-installer.nsi` - Script NSIS para criar o instalador
- `license.txt` - Texto da licença que será exibido durante a instalação
- `build-installer.bat` - Script batch para automatizar a compilação do instalador

## Personalização

Você pode personalizar o instalador modificando o arquivo `zabbix-manager-installer.nsi`:

- Altere as definições no início do arquivo para mudar o nome do aplicativo, versão, etc.
- Adicione arquivos adicionais à seção de instalação se necessário.
- Modifique as páginas e recursos do instalador conforme necessário.

## Problemas Comuns

- **Erro "NSIS não encontrado"**: Certifique-se de que o NSIS está instalado e o diretório de instalação está no PATH do sistema.
- **Erro "zabbix-manager.exe não encontrado"**: Certifique-se de compilar o aplicativo antes de criar o instalador.
- **Erros de compilação NSIS**: Verifique a sintaxe do script NSIS e certifique-se de que todos os caminhos de arquivo estão corretos.