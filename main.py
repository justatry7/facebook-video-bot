import logging
import yt_dlp
from aiogram import Bot, Dispatcher, types
from aiogram.types import InlineKeyboardMarkup, InlineKeyboardButton
from aiogram.utils import executor
from flask import Flask, request, jsonify

# Настройка логирования
logging.basicConfig(level=logging.DEBUG)
logger = logging.getLogger()

# Токен API (вставлен прямо в коде)
API_TOKEN = "7798675393:AAEltxpXHGY6uJ920eyrp_CR4XrQ79W1odQ"

bot = Bot(token=API_TOKEN)
dp = Dispatcher(bot)

# Словарь с переводами
LANGUAGES = {
    "ru": {
        "start_message": "Привет! Отправь мне ссылку на видео из Facebook.",
        "downloading": "Загружаю видео, подождите...",
        "error": "Произошла ошибка при скачивании видео.",
        "not_facebook": "Это не ссылка на видео с Facebook!",
        "video_ready": "Вот ваше видео!",
        "language_changed": "Язык был изменен на русский.",
    },
    "en": {
        "start_message": "Hi! Send me a link to a Facebook video.",
        "downloading": "Downloading video, please wait...",
        "error": "An error occurred while downloading the video.",
        "not_facebook": "This is not a Facebook video link!",
        "video_ready": "Here is your video!",
        "language_changed": "Language has been changed to English.",
    },
}

# Словарь для хранения языка пользователя
user_languages = {}

# Функция для получения языка пользователя (по умолчанию английский)
def get_language(user_id):
    return user_languages.get(user_id, "en")

# Функция для создания кнопки выбора языка
def language_keyboard():
    keyboard = InlineKeyboardMarkup(row_width=1)
    btn_en = InlineKeyboardButton("English", callback_data="set_language_en")
    btn_ru = InlineKeyboardButton("Русский", callback_data="set_language_ru")
    keyboard.add(btn_en, btn_ru)
    return keyboard

# Обработчик команды /start
@dp.message_handler(commands=["start"])
async def start_command(message: types.Message):
    language = get_language(message.from_user.id)
    await message.reply(LANGUAGES[language]["start_message"], reply_markup=language_keyboard())

# Обработчик callback запросов для смены языка
@dp.callback_query_handler(lambda c: c.data.startswith("set_language_"))
async def set_language(callback_query: types.CallbackQuery):
    language = callback_query.data.split('_')[-1]

    if language in LANGUAGES:
        user_languages[callback_query.from_user.id] = language  # Сохраняем язык для пользователя
        await bot.answer_callback_query(callback_query.id, text=LANGUAGES[language]["language_changed"])

        # Отправляем сообщение о смене языка
        await bot.send_message(callback_query.from_user.id, LANGUAGES[language]["language_changed"])

        # Запросить ссылку после смены языка
        await bot.send_message(callback_query.from_user.id, LANGUAGES[language]["start_message"])

# Обработчик для получения ссылки и скачивания видео
@dp.message_handler()
async def download_video(message: types.Message):
    url = message.text
    language = get_language(message.from_user.id)

    if "facebook.com" in url:
        await message.reply(LANGUAGES[language]["downloading"])

        try:
            ydl_opts = {
                'quiet': True,
                'format': 'best',
                'outtmpl': '%(id)s.%(ext)s'
            }

            with yt_dlp.YoutubeDL(ydl_opts) as ydl:
                result = ydl.extract_info(url, download=True)

                video_file = f"{result['id']}.mp4"
                await message.reply_video(open(video_file, 'rb'), caption=LANGUAGES[language]["video_ready"])

        except Exception as e:
            logger.error(f"Error downloading video: {e}")
            await message.reply(LANGUAGES[language]["error"])
    else:
        await message.reply(LANGUAGES[language]["not_facebook"])

# Создание веб-сервера для Vercel
app = Flask(__name__)

@app.route('/webhook', methods=['POST'])
def webhook():
    if request.method == 'POST':
        data = request.json
        logger.debug(f"Received webhook data: {data}")
        return jsonify({"status": "ok"}), 200
