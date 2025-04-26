===== ZABBIX MANAGER =====

Uma aplicação desktop para Windows desenvolvida em Go para gerenciamento e monitoramento 
de infraestrutura Zabbix, com foco em visualização de hosts, itens e triggers, 
e geração de relatórios em CSV.

== COMO USAR ==

1. Execute o arquivo "ZabbixManager-Console.exe" (versão console)
   - Basta clicar duas vezes no executável

2. Na primeira execução, você será solicitado a adicionar um servidor Zabbix:
   - Informe a URL do servidor (ex: https://seu-zabbix.exemplo.com)
   - Informe o token de API (gerado no frontend do Zabbix)
   - Opcionalmente, dê um nome para este perfil

3. Após conectar, você poderá:
   - Listar e buscar hosts no Zabbix
   - Ver itens e triggers de cada host
   - Exportar relatórios para arquivos CSV
   - Gerenciar múltiplos servidores Zabbix

== COMPILAÇÃO (OPCIONAL) ==

Se você quiser compilar o aplicativo a partir do código-fonte:

1. Execute o arquivo "build-windows.bat"
   - Um novo executável será gerado

== REQUISITOS ==

- Windows 7/8/10/11 (64-bit)
- Conexão com internet para acessar o servidor Zabbix

== CONTATO ==

Para suporte ou sugestões, entre em contato pelo email: suporte@zabbixmanager.com