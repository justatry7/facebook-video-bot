import { Telegraf } from "telegraf";
import axios from "axios";
import dotenv from "dotenv";

dotenv.config();

const bot = new Telegraf(process.env.BOT_TOKEN);

bot.start((ctx) => {
  ctx.reply("Send me a Facebook video link, and I'll download it for you!");
});

bot.on("text", async (ctx) => {
  const url = ctx.message.text;

  // Проверка, является ли ссылка ссылкой на видео с Facebook
  if (!url.includes("facebook.com")) {
    return ctx.reply("This is not a Facebook video link!");
  }

  ctx.reply("Downloading video, please wait...");

  try {
    // Получаем ссылку на видео через fdown.net API
    const response = await axios.get(`https://api.fdown.net/api/download?url=${encodeURIComponent(url)}`);
    console.log("API Response: ", response.data);

    // Проверяем, есть ли в ответе ссылка на видео
    if (response.data && response.data.downloadUrl) {
      // Если ссылка есть, отправляем видео
      ctx.replyWithVideo({ url: response.data.downloadUrl }, { caption: "Here is your video!" });
    } else {
      // Если ссылка на видео не получена, выводим ошибку
      ctx.reply("Sorry, I couldn't download this video.");
    }
  } catch (error) {
    // Логирование ошибок
    console.error("Error while downloading video:", error);
    ctx.reply("An error occurred while downloading the video.");
  }
});

bot.launch();
