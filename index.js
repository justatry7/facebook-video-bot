import { Telegraf } from "telegraf";
import express from "express";
import axios from "axios";
import dotenv from "dotenv";

dotenv.config();

const bot = new Telegraf(process.env.BOT_TOKEN); // Токен из переменных окружения

const app = express();
const PORT = process.env.PORT || 3000;

// Команда /start
bot.start((ctx) => {
  ctx.reply("Send me a Facebook video link, and I'll download it for you!");
});

// Обработчик текстовых сообщений (ожидаем ссылку)
bot.on("text", async (ctx) => {
  const url = ctx.message.text;
  
  if (!url.includes("facebook.com")) {
    return ctx.reply("This is not a Facebook video link!");
  }

  ctx.reply("Downloading video, please wait...");

  try {
    const response = await axios.get(`https://api.fdown.net/api/download?url=${encodeURIComponent(url)}`);
    if (response.data && response.data.downloadUrl) {
      ctx.replyWithVideo({ url: response.data.downloadUrl }, { caption: "Here is your video!" });
    } else {
      ctx.reply("Sorry, I couldn't download this video.");
    }
  } catch (error) {
    ctx.reply("An error occurred while downloading the video.");
  }
});

// Запускаем бота
bot.launch();

// Express для Vercel
app.get("/", (req, res) => {
  res.send("Bot is running!");
});

// Запуск сервера
app.listen(PORT, () => {
  console.log(`Server running on port ${PORT}`);
});
