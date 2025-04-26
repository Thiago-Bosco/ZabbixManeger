# PowerShell Script para criar um launcher para o Zabbix Manager
# Este script cria um atalho (shortcut) .lnk para o executável do Zabbix Manager

# Configurações
$executableName = "ZabbixManager-Console.exe"
$shortcutName = "Zabbix Manager.lnk"
$iconFile = "zabbix-manager-icon.ico" # Opcional - se você tiver um arquivo de ícone

# Verifica se o executável existe
if (-not (Test-Path $executableName)) {
    Write-Host "Erro: $executableName não encontrado na pasta atual."
    Write-Host "Execute este script na mesma pasta onde o executável está localizado."
    Read-Host "Pressione ENTER para sair"
    exit
}

# Cria um atalho para o executável
$WshShell = New-Object -ComObject WScript.Shell
$Shortcut = $WshShell.CreateShortcut($shortcutName)
$Shortcut.TargetPath = (Get-Item $executableName).FullName
$Shortcut.WorkingDirectory = (Get-Location).Path

# Define um ícone se o arquivo existir
if (Test-Path $iconFile) {
    $Shortcut.IconLocation = (Get-Item $iconFile).FullName
}

# Adiciona uma descrição
$Shortcut.Description = "Gerenciador de infraestrutura Zabbix"

# Salva o atalho
$Shortcut.Save()

Write-Host "Atalho '$shortcutName' criado com sucesso!"
Write-Host "Você pode agora clicar duas vezes neste atalho para iniciar o Zabbix Manager."

# Opcionalmente, cria um atalho na área de trabalho
$desktopPath = [System.Environment]::GetFolderPath("Desktop")
$desktopShortcut = "$desktopPath\$shortcutName"

$createDesktopShortcut = Read-Host "Deseja criar um atalho na área de trabalho? (S/N)"
if ($createDesktopShortcut -eq "S" -or $createDesktopShortcut -eq "s") {
    Copy-Item $shortcutName $desktopShortcut
    Write-Host "Atalho criado na área de trabalho."
}

Write-Host "Concluído!"
Read-Host "Pressione ENTER para sair"