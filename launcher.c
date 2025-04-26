#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <windows.h>

// Função para verificar se um arquivo existe
int file_exists(const char *filename) {
    FILE *file = fopen(filename, "r");
    if (file) {
        fclose(file);
        return 1;
    }
    return 0;
}

// Função para criar um arquivo de configuração se não existir
void create_config_if_not_exists() {
    const char *config_dir = ".\\config";
    const char *config_file = ".\\config\\config.json";
    
    // Verificar se o diretório config existe, se não, criar
    if (GetFileAttributes(config_dir) == INVALID_FILE_ATTRIBUTES) {
        CreateDirectory(config_dir, NULL);
        printf("Diretório de configuração criado.\n");
    }
    
    // Verificar se o arquivo de configuração existe
    if (!file_exists(config_file)) {
        FILE *file = fopen(config_file, "w");
        if (file) {
            // Estrutura básica JSON para configuração
            fprintf(file, "{\n");
            fprintf(file, "  \"servidores\": []\n");
            fprintf(file, "}\n");
            fclose(file);
            printf("Arquivo de configuração inicial criado.\n");
        } else {
            printf("Erro ao criar arquivo de configuração.\n");
        }
    }
}

int main() {
    char cmd[256];
    const char *executable = "ZabbixManager-Console.exe";
    
    // Verificar se o executável principal existe
    if (!file_exists(executable)) {
        printf("Erro: %s não encontrado!\n", executable);
        printf("Verifique se o executável está na mesma pasta deste launcher.\n");
        printf("Pressione qualquer tecla para sair...");
        getchar();
        return 1;
    }
    
    // Criar configuração se necessário
    create_config_if_not_exists();
    
    // Construir comando para iniciar o programa principal
    sprintf(cmd, "%s", executable);
    
    // Iniciar o programa
    printf("Iniciando Zabbix Manager...\n");
    
    STARTUPINFO si;
    PROCESS_INFORMATION pi;
    
    ZeroMemory(&si, sizeof(si));
    si.cb = sizeof(si);
    ZeroMemory(&pi, sizeof(pi));
    
    // Iniciar o processo
    if (CreateProcess(NULL, cmd, NULL, NULL, FALSE, 0, NULL, NULL, &si, &pi)) {
        // Fechar os handles
        CloseHandle(pi.hProcess);
        CloseHandle(pi.hThread);
        return 0;
    } else {
        printf("Erro ao iniciar o processo: %d\n", GetLastError());
        printf("Pressione qualquer tecla para sair...");
        getchar();
        return 1;
    }
}