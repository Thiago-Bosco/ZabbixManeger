#!/usr/bin/env python3
"""
Launcher para Zabbix Manager
Este script inicia o aplicativo Zabbix Manager e cria arquivos de configuração iniciais se necessário.
"""

import os
import json
import subprocess
import sys
import time

# Configurações
EXECUTABLE = "ZabbixManager-Console.exe"
CONFIG_DIR = "config"
CONFIG_FILE = os.path.join(CONFIG_DIR, "config.json")

def ensure_config_exists():
    """Verifica se o diretório e arquivo de configuração existem, e cria se necessário."""
    # Criar diretório de configuração se não existir
    if not os.path.exists(CONFIG_DIR):
        os.makedirs(CONFIG_DIR)
        print(f"Diretório de configuração '{CONFIG_DIR}' criado.")
    
    # Criar arquivo de configuração se não existir
    if not os.path.exists(CONFIG_FILE):
        initial_config = {
            "servidores": []
        }
        
        with open(CONFIG_FILE, 'w', encoding='utf-8') as f:
            json.dump(initial_config, f, indent=2)
        
        print(f"Arquivo de configuração inicial '{CONFIG_FILE}' criado.")

def launch_application():
    """Inicia o aplicativo Zabbix Manager."""
    if not os.path.exists(EXECUTABLE):
        print(f"Erro: Executável '{EXECUTABLE}' não encontrado!")
        print("Verifique se o executável está na mesma pasta deste launcher.")
        input("Pressione ENTER para sair...")
        return False
    
    try:
        print(f"Iniciando {EXECUTABLE}...")
        
        # Em Windows
        if sys.platform.startswith('win'):
            subprocess.Popen([EXECUTABLE], creationflags=subprocess.CREATE_NEW_CONSOLE)
        # Em Linux/macOS
        else:
            subprocess.Popen([f"./{EXECUTABLE}"])
            
        return True
    except Exception as e:
        print(f"Erro ao iniciar o aplicativo: {e}")
        input("Pressione ENTER para sair...")
        return False

def main():
    """Função principal."""
    print("=== Zabbix Manager Launcher ===")
    
    # Garantir que a configuração existe
    ensure_config_exists()
    
    # Iniciar o aplicativo
    success = launch_application()
    
    if success:
        print("Aplicativo iniciado com sucesso!")
        time.sleep(2)  # Pequena pausa para mostrar a mensagem

if __name__ == "__main__":
    main()