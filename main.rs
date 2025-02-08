use std::process::Command;
use std::env;
use reqwest::{Client, Error};
use std::fs::File;
use std::io::Write;
use serde::{Deserialize, Serialize};

#[derive(Deserialize)]
struct TelegramMessage {
    chat_id: String,
    text: String,
}

#[derive(Serialize)]
struct SendMessage {
    chat_id: String,
    text: String,
}

async fn send_message(chat_id: &str, text: &str, token: &str) -> Result<(), Error> {
    let client = Client::new();
    let url = format!("https://api.telegram.org/bot{}/sendMessage", token);
    
    let params = SendMessage {
        chat_id: chat_id.to_string(),
        text: text.to_string(),
    };
    
    client.post(url)
        .json(&params)
        .send()
        .await?;

    Ok(())
}

async fn download_video(url: &str, chat_id: &str, token: &str) -> Result<(), std::io::Error> {
    // Отправляем сообщение о начале загрузки
    send_message(chat_id, "Downloading video, please wait...", token).await.unwrap();

    // Используем yt-dlp для скачивания видео
    let output = Command::new("yt-dlp")
        .arg(url)
        .arg("-o")
        .arg("downloaded_video.%(ext)s")
        .output()
        .expect("Failed to execute command");

    if !output.status.success() {
        send_message(chat_id, "An error occurred while downloading the video.", token)
            .await.unwrap();
        return Err(std::io::Error::new(std::io::ErrorKind::Other, "Download failed"));
    }

    // Отправляем видео обратно
    let video_path = "downloaded_video.mp4";
    let mut video_file = File::open(video_path)?;
    let mut video_data = Vec::new();
    video_file.read_to_end(&mut video_data)?;

    // Отправляем файл видео обратно пользователю
    let upload_url = format!("https://api.telegram.org/bot{}/sendVideo", token);
    let client = Client::new();
    client.post(upload_url)
        .form(&[("chat_id", chat_id), ("video", video_data)])
        .send()
        .await?;

    send_message(chat_id, "Here is your video!", token).await.unwrap();
    Ok(())
}

#[tokio::main]
async fn main() {
    dotenv::dotenv().ok();
    let token = env::var("TELEGRAM_BOT_TOKEN").expect("TELEGRAM_BOT_TOKEN must be set");
    
    // URL видео и chat_id могут быть переданы через запрос или API
    let chat_id = "your_chat_id"; // это нужно будет заменить на id чата
    let video_url = "https://www.facebook.com/video_url"; // это пример URL на видео

    match download_video(video_url, chat_id, &token).await {
        Ok(_) => println!("Video downloaded and sent successfully!"),
        Err(err) => eprintln!("Error: {}", err),
    }
}
