@echo off
echo Compilando o instalador para Zabbix Manager...

REM Verificar se o NSIS está instalado
where makensis >nul 2>nul
if %ERRORLEVEL% neq 0 (
    echo ERRO: NSIS não encontrado. Por favor, instale o NSIS antes de continuar.
    echo Você pode baixá-lo em: https://nsis.sourceforge.io/Download
    pause
    exit /b 1
)

REM Verificar se o executável existe
if not exist "..\..\zabbix-manager.exe" (
    echo ERRO: zabbix-manager.exe não encontrado.
    echo Por favor, compile o aplicativo antes de criar o instalador.
    pause
    exit /b 1
)

REM Verificar se o diretório assets existe
if not exist "..\..\assets" (
    echo AVISO: Diretório assets não encontrado. Criando diretório vazio...
    mkdir "..\..\assets"
)

REM Compilar o instalador
echo Compilando o instalador...
makensis zabbix-manager-installer.nsi

if %ERRORLEVEL% equ 0 (
    echo Instalador compilado com sucesso!
    echo O instalador está disponível em: %CD%\ZabbixManager-Setup.exe
) else (
    echo ERRO ao compilar o instalador.
)

pause