# 📋 НАСТРОЙКА ПЕРЕМЕННЫХ ОКРУЖЕНИЯ (.env)

**Статус:** ✅ ПОЛНОЕ РУКОВОДСТВО  
**Дата:** 2025-12-08

---

## ⚡ БЫСТРЫЙ СТАРТ (5 МИНУТ)

### 1️⃣ Создать .env файл

```bash
cd ~/send-to-kindle-telegram-bot
cp .env.example .env
nano .env  # или vim, или ваш редактор
```

### 2️⃣ Заполнить обязательные переменные

```bash
# Telegram бот токен (от @BotFather)
UBOT_TELEGRAM_TOKEN=123456789:ABCdefGHIjklMNOpqrsTUVwxyz

UBOT_EMAIL_FROM=your-email@gmail.com
UBOT_PASSWORD=xxxx-xxxx-xxxx-xxxx
UBOT_SMTP_HOST=smtp.gmail.com
UBOT_SMTP_PORT=587
UBOT_EMAIL_TO=your-kindle@kindle.com
```

### 3️⃣ Перезагрузить контейнер

```bash
docker compose down
docker compose up -d
docker compose logs sendtokindle --tail 20
```

**Готово! 🎉**

---

## 📖 ПОЛНОЕ РУКОВОДСТВО

### UBOT_TELEGRAM_TOKEN

**Что это?** Токен доступа к твоему Telegram боту

**Как получить?**
1. Откройте Telegram
2. Напишите @BotFather
3. Выберите /newbot
4. Следуйте инструкциям
5. Получите токен вида: `123456789:ABCdefGHIjklMNOpqrsTUVwxyz`

**Пример:**
```bash
UBOT_TELEGRAM_TOKEN=8283610744:AAHD5Ml9rAiuz3FcCbBrDqRYTlciLuVk4ws
```

**⚠️ ВАЖНО:** Это твой секретный ключ! Никому не давай!

---

### UBOT_EMAIL_FROM

**Что это?** Email адрес с которого отправляются книги на Kindle

**Какой email использовать?**
- ✅ Gmail (рекомендуется)
- ✅ Yandex.Почта
- ✅ Outlook / Microsoft 365
- ✅ Любой IMAP/SMTP email
- ✅ Корпоративная почта (если поддерживает SMTP)

**Пример:**
```bash
UBOT_EMAIL_FROM=book@nodkeys.com
```

**⚠️ ВАЖНО:** Этот email ДОЛЖЕН быть добавлен в Amazon белый лист!

Как добавить в Amazon:
1. Перейти https://www.amazon.com/gp/digital/fiona/
2. Найти "Approved Personal Document E-mail List"
3. Нажать "Add a new approved e-mail address"
4. Добавить твой `UBOT_EMAIL_FROM`
5. Убедиться что включено "Receive e-mail-based Kindle documents"

---

### UBOT_PASSWORD

**Что это?** Пароль для SMTP авторизации

**⚠️ КРИТИЧЕСКИ ВАЖНО:** Это НЕ твой основной пароль почты!

Всегда создавай отдельный пароль приложения:

**Для Gmail:**
1. Перейти https://myaccount.google.com/apppasswords
2. Выбрать "Mail" и "Windows Computer" (или твое устройство)
3. Скопировать сгенерированный пароль (16 символов с пробелами)
4. Убрать пробелы или оставить как есть

**Для Yandex.Почта:**
1. Перейти https://passport.yandex.ru/
2. "Безопасность" → "Пароли приложений"
3. Выбрать тип "Почта"
4. Создать пароль
5. Скопировать

**Для Outlook:**
1. Перейти https://account.live.com/
2. "Безопасность" → "Пароли приложений"
3. Создать пароль
4. Скопировать

**Для корпоративной почты:**
1. Обратиться к администратору
2. Попросить пароль приложения для SMTP

**Пример:**
```bash
UBOT_PASSWORD=pegboz-sibkus-4cYsjy
```

---

### UBOT_SMTP_HOST

**Что это?** Адрес SMTP сервера почтовой службы

**Популярные значения:**

| Провайдер | SMTP Host | Порт | Тип |
|-----------|-----------|------|-----|
| Gmail | `smtp.gmail.com` | 587 | TLS |
| Yandex | `smtp.yandex.com` | 587 | TLS |
| Yandex (alt) | `smtp.yandex.ru` | 587 | TLS |
| Outlook | `smtp.live.com` | 587 | TLS |
| Outlook (alt) | `smtp-mail.outlook.com` | 587 | TLS |
| Microsoft 365 | `smtp.office365.com` | 587 | TLS |
| Nodkeys | `mail.nodkeys.com` | 587 | TLS |

**Пример:**
```bash
UBOT_SMTP_HOST=mail.nodkeys.com
```

**Если не знаешь значение:** Гугли "SMTP settings [твой провайдер]"

---

### UBOT_SMTP_PORT

**Что это?** Порт SMTP сервера

**Стандартные значения:**
- `587` - TLS (рекомендуется, безопасно)
- `465` - SSL (альтернатива, тоже безопасно)
- `25` - обычный (редко работает, не рекомендуется)

**Пример:**
```bash
UBOT_SMTP_PORT=587
```

**Совет:** Если 587 не работает, попробуй 465

---

### UBOT_EMAIL_TO (или UBOT_KINDLE_DEVICES)

**Что это?** Email адрес твоего Kindle

#### Вариант А: ОДНО устройство

```bash
UBOT_EMAIL_TO=vera_muhamedova_abyH2D@kindle.com
```

**Как найти свой Kindle email?**
1. Перейти https://www.amazon.com/gp/digital/fiona/
2. Раздел "Devices"
3. Найти своё устройство
4. Email вида: `username_randomchars@kindle.com`

#### Вариант Б: НЕСКОЛЬКО устройств (НОВОЕ!)

```bash
UBOT_KINDLE_DEVICES=Kindle1:vera_muhamedova_abyH2D@kindle.com
```

**Формат:**
```
Имя:email@kindle.com|Имя2:email2@kindle.com|...
```

**Примеры:**
```bash
# Два устройства
UBOT_KINDLE_DEVICES=My Paperwhite:personal@kindle.com|Work Tablet:work@kindle.com

# Три устройства
UBOT_KINDLE_DEVICES=Спальня:bed@kindle.com|Офис:office@kindle.com|Путешествия:travel@kindle.com
```

**⚠️ ВАЖНО:** Если используешь `UBOT_KINDLE_DEVICES`, оставь `UBOT_EMAIL_TO` закомментированным или пустым!

---

### UBOT_SMTP_INSECURE (опционально)

**Что это?** Флаг для отключения проверки SSL сертификата

**Когда использовать?** Только если:
- Используешь самоподписанный сертификат
- SMTP сервер с проблемами SSL

**Значения:**
```bash
# Включено (небезопасно, но работает с некоторыми сервами)
UBOT_SMTP_INSECURE=true

# Отключено (безопасно, рекомендуется)
UBOT_SMTP_INSECURE=false
# или просто закомментируй эту строку
```

**По умолчанию:** false (безопасно)

---

## ✅ ПОЛНЫЙ ПРИМЕР .env

### Для одного Kindle

```bash
# Telegram BOT Configuration
UBOT_TELEGRAM_TOKEN=8283610744:AAHD5Ml9rAiuz3FcCbBrDqRYTlciLuVk4ws

# Email Configuration (sender)
UBOT_EMAIL_FROM=book@nodkeys.com
UBOT_PASSWORD=pegboz-sibkus-4cYsjy

# Kindle Email Configuration - Single Device
UBOT_EMAIL_TO=vera_muhamedova_abyH2D@kindle.com

# SMTP Configuration
UBOT_SMTP_HOST=mail.nodkeys.com
UBOT_SMTP_PORT=587
UBOT_SMTP_INSECURE=true
```

### Для нескольких Kindle

```bash
# Telegram BOT Configuration
UBOT_TELEGRAM_TOKEN=8283610744:AAHD5Ml9rAiuz3FcCbBrDqRYTlciLuVk4ws

# Email Configuration (sender)
UBOT_EMAIL_FROM=book@nodkeys.com
UBOT_PASSWORD=pegboz-sibkus-4cYsjy

# Kindle Email Configuration - Multiple Devices
UBOT_KINDLE_DEVICES=Kindle1:vera_muhamedova_abyH2D@kindle.com|Kindle2:email2@kindle.com

# SMTP Configuration
UBOT_SMTP_HOST=mail.nodkeys.com
UBOT_SMTP_PORT=587
UBOT_SMTP_INSECURE=true
```

---

## 🆘 ЧАСТЫЕ ОШИБКИ И ИСПРАВЛЕНИЯ

### ❌ Ошибка: Синтаксис в .env

**НЕПРАВИЛЬНО:**
```bash
# Комментарий в середине строки (Java-style)
UBOT_EMAIL_FROM=book@nodkeys.com  # Это email

# Пробелы вокруг = (Python-style)
UBOT_PASSWORD = pegboz-sibkus-4cYsjy

# URL с особыми символами без кавычек
UBOT_SMTP_HOST=mail.nodkeys.com?ssl=true

# Закомментированная переменная в меню (Markdown-style)
#UBOT_EMAIL_TO=vera@kindle.com
```

**ПРАВИЛЬНО:**
```bash
# Комментарий на отдельной строке
UBOT_EMAIL_FROM=book@nodkeys.com

# Без пробелов вокруг =
UBOT_PASSWORD=pegboz-sibkus-4cYsjy

# Простые значения без спецсимволов
UBOT_SMTP_HOST=mail.nodkeys.com

# Закомментировано правильно
# UBOT_EMAIL_TO=vera@kindle.com
```

### ❌ Ошибка: "could not start telegram bot: emailto not set"

**Причина:** `UBOT_EMAIL_TO` и `UBOT_KINDLE_DEVICES` оба не установлены

**Решение:** Убедись что хотя бы один из них есть:
```bash
# Вариант 1: UBOT_EMAIL_TO
UBOT_EMAIL_TO=vera_muhamedova_abyH2D@kindle.com

# Вариант 2: UBOT_KINDLE_DEVICES
UBOT_KINDLE_DEVICES=Kindle1:vera_muhamedova_abyH2D@kindle.com
```

### ❌ Ошибка: "Authentication failed"

**Причины:**
1. Неправильный `UBOT_PASSWORD`
2. Email не в белом списке Amazon
3. SMTP сервер требует TLS/SSL

**Решение:**
```bash
# 1. Проверить пароль приложения (не основной пароль!)
UBOT_PASSWORD=xxxx-xxxx-xxxx-xxxx

# 2. Добавить email в Amazon
# https://www.amazon.com/gp/digital/fiona/

# 3. Убедиться что SMTP порт правильный
UBOT_SMTP_PORT=587  # TLS
# или
UBOT_SMTP_PORT=465  # SSL
```

### ❌ Ошибка: "Connection refused" или "Connection timeout"

**Причины:**
1. Неправильный SMTP сервер
2. Неправильный порт
3. Сетевые проблемы

**Решение:**
```bash
# Проверить SMTP сервер и порт
UBOT_SMTP_HOST=mail.nodkeys.com
UBOT_SMTP_PORT=587

# Проверить подключение к интернету
ping mail.nodkeys.com

# Если нужен SMTP без SSL проверки (self-signed certificates)
UBOT_SMTP_INSECURE=true
```

---

## 🧪 ПРОВЕРКА КОНФИГУРАЦИИ

### 1️⃣ Проверить что .env существует

```bash
ls -la ~/send-to-kindle-telegram-bot/.env
```

**Должно быть:**
```
-rw-r--r-- 1 root root 500 Dec  8 10:00 .env
```

### 2️⃣ Проверить что переменные заполнены

```bash
grep "^UBOT_" ~/send-to-kindle-telegram-bot/.env
```

**Должно быть:**
```
UBOT_TELEGRAM_TOKEN=8283610744:AAHD5Ml9rAiuz3FcCbBrDqRYTlciLuVk4ws
UBOT_EMAIL_FROM=book@nodkeys.com
UBOT_PASSWORD=pegboz-sibkus-4cYsjy
UBOT_SMTP_HOST=mail.nodkeys.com
UBOT_SMTP_PORT=587
UBOT_EMAIL_TO=vera_muhamedova_abyH2D@kindle.com
```

### 3️⃣ Проверить что нет пустых значений

```bash
grep "=$" ~/send-to-kindle-telegram-bot/.env
```

**Не должно быть вывода!** Если есть - значит переменная пуста.

### 4️⃣ Проверить синтаксис в контейнере

```bash
# Перезагрузить контейнер
docker compose down
docker compose up -d

# Проверить логи
sleep 3
docker compose logs sendtokindle --tail 20
```

**Должно быть:**
```
✅ Нормальный запуск БЕЗ ошибок "emailto not set"
```

---

## 🔐 БЕЗОПАСНОСТЬ

### Защита .env файла

```bash
# Установить права доступа (только владелец может читать)
chmod 600 ~/.env

# Проверить
ls -la ~/.env
# Должно быть: -rw------- 1 root root ...
```

### Резервная копия

```bash
# Сохранить backup
cp ~/.env ~/.env.backup

# Восстановить если что-то сломается
cp ~/.env.backup ~/.env
```

### Git - НЕ коммитить .env!

```bash
# Добавить в .gitignore
echo ".env" >> .gitignore
echo ".env.backup" >> .gitignore

# Проверить что .env не в git
git status | grep env
# Не должно быть вывода
```

---

## ✨ ИТОГОВЫЙ ЧЕКЛИСТ

Перез запуском контейнера:
- [ ] Файл `.env` создан
- [ ] Токен от @BotFather установлен в `UBOT_TELEGRAM_TOKEN`
- [ ] Email отправителя установлен в `UBOT_EMAIL_FROM`
- [ ] Пароль приложения установлен в `UBOT_PASSWORD`
- [ ] SMTP хост установлен в `UBOT_SMTP_HOST`
- [ ] SMTP порт установлен в `UBOT_SMTP_PORT`
- [ ] Либо `UBOT_EMAIL_TO`, либо `UBOT_KINDLE_DEVICES` установлен
- [ ] Email отправителя добавлен в Amazon белый лист
- [ ] Все переменные имеют значения (не пусто)
- [ ] Нет синтаксических ошибок
- [ ] Права доступа установлены: `chmod 600 .env`

После запуска контейнера:
- [ ] `docker compose logs` не содержит ошибку "emailto not set"
- [ ] Контейнер запущен и healthy
- [ ] Бот отвечает на сообщения в Telegram
- [ ] Файлы отправляются на Kindle

**Если все чекбоксы отмечены - всё работает! 🎉**

---

**Дата создания:** 2025-12-08  
**Статус:** ✅ ПОЛНОЕ РУКОВОДСТВО  
**Версия:** 1.0
