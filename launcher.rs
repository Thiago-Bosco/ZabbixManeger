use std::fs::{self, File};
use std::io::Write;
use std::path::Path;
use std::process::Command;

fn main() {
    println!("=== Zabbix Manager Launcher ===");
    
    // Verificar se o executável existe
    let executable = "ZabbixManager-Console.exe";
    if !Path::new(executable).exists() {
        println!("Erro: {} não encontrado!", executable);
        println!("Verifique se o executável está na mesma pasta deste launcher.");
        pause();
        return;
    }
    
    // Garantir que a configuração existe
    ensure_config_exists();
    
    // Iniciar o aplicativo
    println!("Iniciando {}...", executable);
    
    match Command::new(executable).spawn() {
        Ok(_) => {
            println!("Aplicativo iniciado com sucesso!");
        },
        Err(e) => {
            println!("Erro ao iniciar o aplicativo: {}", e);
            pause();
        }
    }
}

fn ensure_config_exists() {
    // Criar diretório de configuração se não existir
    let config_dir = "config";
    if !Path::new(config_dir).exists() {
        match fs::create_dir(config_dir) {
            Ok(_) => println!("Diretório de configuração '{}' criado.", config_dir),
            Err(e) => println!("Erro ao criar diretório de configuração: {}", e)
        }
    }
    
    // Criar arquivo de configuração se não existir
    let config_file = format!("{}/config.json", config_dir);
    if !Path::new(&config_file).exists() {
        let initial_config = r#"{
  "servidores": []
}"#;
        
        match File::create(&config_file).and_then(|mut file| {
            file.write_all(initial_config.as_bytes())
        }) {
            Ok(_) => println!("Arquivo de configuração inicial '{}' criado.", config_file),
            Err(e) => println!("Erro ao criar arquivo de configuração: {}", e)
        }
    }
}

fn pause() {
    println!("Pressione ENTER para continuar...");
    let mut input = String::new();
    std::io::stdin().read_line(&mut input).unwrap();
}

// Para compilar este código Rust:
// 1. Instale o Rust: https://www.rust-lang.org/tools/install
// 2. Execute: rustc launcher.rs
// 3. Isso gerará um executável launcher.exe