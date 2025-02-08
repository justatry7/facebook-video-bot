import { Telegraf } from 'telegraf';
import express from 'express';
import axios from 'axios';
import dotenv from 'dotenv';

dotenv.config();

// Создаем объект бота с токеном из переменных окружения
const bot = new Telegraf(process.env.BOT_TOKEN);

// Настроим Express сервер для Vercel
const app = express();
const PORT = process.env.PORT || 3000;

// Команда /start
bot.start((ctx) => {
  ctx.reply("Send me a Facebook video link, and I'll download it for you!");
});

// Обработчик текстовых сообщений
bot.on('text', async (ctx) => {
  const url = ctx.message.text;

  if (!url.includes("facebook.com")) {
    return ctx.reply("This is not a Facebook video link!");
  }

  ctx.reply("Downloading video, please wait...");

  try {
    // Запрос для скачивания видео через fdown API
    const response = await axios.get(`https://api.fdown.net/api/download?url=${encodeURIComponent(url)}`);
    
    if (response.data && response.data.downloadUrl) {
      // Отправляем видео, если ссылка на скачивание найдена
      ctx.replyWithVideo({ url: response.data.downloadUrl }, { caption: "Here is your video!" });
    } else {
      ctx.reply("Sorry, I couldn't download this video.");
    }
  } catch (error) {
    console.error(error);
    ctx.reply("An error occurred while downloading the video.");
  }
});

// Запускаем бот
bot.launch();

// Простая экспресс-страница, чтобы убедиться, что сервер работает
app.get("/", (req, res) => {
  res.send("Bot is running!");
});

// Запуск сервера Express
app.listen(PORT, () => {
  console.log(`Server running on port ${PORT}`);
});
