@echo off
echo Compilando Zabbix Manager...

REM Verificar se Go está instalado
where go >nul 2>nul
if %ERRORLEVEL% neq 0 (
    echo Erro: Go nao esta instalado ou nao esta no PATH.
    echo Por favor, instale Go de https://golang.org/dl/
    pause
    exit /b 1
)

REM Modo headless (console) - mais simples e sem dependências externas
echo Criando versao headless (console)...
go build -tags=headless -o ZabbixManager-Console.exe

if %ERRORLEVEL% neq 0 (
    echo Erro ao compilar o aplicativo.
    pause
    exit /b 1
)

echo.
echo Compilacao concluida com sucesso!
echo.
echo O arquivo ZabbixManager-Console.exe foi criado.
echo Clique duas vezes nele para executar o aplicativo.
echo.
echo Pressione qualquer tecla para sair...
pause > nul