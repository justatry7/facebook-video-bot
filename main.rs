use reqwest::Client;
use serde::{Deserialize, Serialize};
use std::env;
use std::error::Error;
use std::fs;
use tokio;

#[derive(Serialize, Deserialize)]
struct TelegramResponse {
    ok: bool,
    result: Option<Vec<String>>,
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn Error>> {
    dotenv::dotenv().ok();  // Загружаем переменные окружения из .env

    let token = env::var("TELEGRAM_BOT_TOKEN")?;  // Получаем токен из переменной окружения

    // Пример, как вызвать API Telegram
    let url = format!("https://api.telegram.org/bot{}/getMe", token);
    let client = Client::new();
    let res = client.get(&url).send().await?;

    let body = res.text().await?;
    println!("Telegram API response: {}", body);

    // Загружаем видео по ссылке, передаваемой в сообщении
    // Для этого тебе нужно будет использовать библиотеку для скачивания видео из Facebook (например, через fdown)

    Ok(())
}
