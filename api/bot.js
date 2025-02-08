import { Telegraf } from "telegraf";
import axios from "axios";
import express from "express";
import dotenv from "dotenv";
import ytdl from "ytdl-core";

// Загружаем переменные окружения
dotenv.config();

const bot = new Telegraf(process.env.BOT_TOKEN); // Токен из переменных окружения
const app = express();
const PORT = process.env.PORT || 3000;

// Словарь с переводами
const LANGUAGES = {
  ru: {
    start_message: "Привет! Отправь мне ссылку на видео из Facebook.",
    downloading: "Загружаю видео, подождите...",
    error: "Произошла ошибка при скачивании видео.",
    not_facebook: "Это не ссылка на видео с Facebook!",
    video_ready: "Вот ваше видео!",
    language_changed: "Язык был изменен на русский.",
  },
  en: {
    start_message: "Hi! Send me a link to a Facebook video.",
    downloading: "Downloading video, please wait...",
    error: "An error occurred while downloading the video.",
    not_facebook: "This is not a Facebook video link!",
    video_ready: "Here is your video!",
    language_changed: "Language has been changed to English.",
  },
};

const userLanguages = {};

// Получаем язык пользователя (по умолчанию английский)
function getLanguage(userId) {
  return userLanguages[userId] || "en";
}

// Команда /start
bot.start((ctx) => {
  const language = getLanguage(ctx.from.id);
  ctx.reply(LANGUAGES[language].start_message);
});

// Обработчик текстовых сообщений (ожидаем ссылку)
bot.on("text", async (ctx) => {
  const url = ctx.message.text;
  const language = getLanguage(ctx.from.id);

  if (!url.includes("facebook.com")) {
    return ctx.reply(LANGUAGES[language].not_facebook);
  }

  ctx.reply(LANGUAGES[language].downloading);

  try {
    // Используем ytdl для получения информации о видео
    const info = await ytdl.getInfo(url);
    const videoUrl = info.formats.find(format => format.itag === 22); // Загружаем видео в хорошем качестве

    if (videoUrl) {
      ctx.replyWithVideo(videoUrl.url, { caption: LANGUAGES[language].video_ready });
    } else {
      ctx.reply(LANGUAGES[language].error);
    }
  } catch (error) {
    console.error(error);
    ctx.reply(LANGUAGES[language].error);
  }
});

// Запускаем бота
bot.launch();

// Express сервер для Vercel
app.get("/", (req, res) => {
  res.send("Bot is running!");
});

// Запуск сервера
app.listen(PORT, () => {
  console.log(`Server running on port ${PORT}`);
});
