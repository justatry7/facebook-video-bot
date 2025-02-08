import { Telegraf } from "telegraf";
import axios from "axios";
import dotenv from "dotenv";

dotenv.config();

const bot = new Telegraf(process.env.BOT_TOKEN); // Токен из переменных окружения

// Команда /start
bot.start((ctx) => {
  console.log("User started the bot.");
  ctx.reply("Send me a Facebook video link, and I'll download it for you!");
});

// Обработчик текстовых сообщений (ожидаем ссылку)
bot.on("text", async (ctx) => {
  const url = ctx.message.text;

  if (!url.includes("facebook.com")) {
    console.log("Received non-Facebook link: " + url);
    return ctx.reply("This is not a Facebook video link!");
  }

  console.log("Downloading video from: " + url);
  ctx.reply("Downloading video, please wait...");

  try {
    const response = await axios.get(`https://api.fdown.net/api/download?url=${encodeURIComponent(url)}`);
    console.log("API response:", response.data); // Логируем ответ от API

    if (response.data && response.data.downloadUrl) {
      ctx.replyWithVideo({ url: response.data.downloadUrl }, { caption: "Here is your video!" });
    } else {
      console.error("Error: No download URL returned");
      ctx.reply("Sorry, I couldn't download this video.");
    }
  } catch (error) {
    console.error("Error during video download:", error);
    ctx.reply("An error occurred while downloading the video.");
  }
});

// Запускаем бота
bot.launch().then(() => {
  console.log("Bot started successfully!");
}).catch((error) => {
  console.error("Error starting the bot:", error);
});
