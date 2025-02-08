use reqwest::Client;
use teloxide::prelude::*;
use tokio;
use dotenv::dotenv;
use std::env;

#[tokio::main]
async fn main() {
    dotenv().ok();
    let bot_token = env::var("TELEGRAM_BOT_TOKEN").expect("TELEGRAM_BOT_TOKEN not set");
    let bot = Bot::new(bot_token);

    teloxide::repl(bot, |message: Message, bot: Bot| async move {
        if let Some(text) = message.text() {
            if text.contains("facebook.com") {
                bot.send_message(message.chat.id, "Downloading video...").await?;
                
                // Получаем ссылку для скачивания (функция на основе внешнего API)
                let video_url = get_video_url(text).await;
                
                match video_url {
                    Some(url) => {
                        bot.send_video(message.chat.id, url)
                            .caption("Here is your Facebook video!")
                            .await?;
                    }
                    None => {
                        bot.send_message(message.chat.id, "Failed to download video.").await?;
                    }
                }
            } else {
                bot.send_message(message.chat.id, "This is not a Facebook link!").await?;
            }
        }
        Ok(())
    })
    .await;
}

async fn get_video_url(url: &str) -> Option<String> {
    // Подключение к внешнему API для скачивания видео (например, fdown.net)
    let client = Client::new();
    let res = client
        .get(format!("https://api.fdown.net/api/download?url={}", url))
        .send()
        .await;

    match res {
        Ok(response) => {
            if let Ok(json) = response.json::<serde_json::Value>().await {
                json["downloadUrl"].as_str().map(String::from)
            } else {
                None
            }
        }
        Err(_) => None,
    }
}
