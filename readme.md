[![GitHub Repository](https://img.shields.io/badge/GitHub-Repository-blue?logo=github)](https://github.com/chazari-x/hmtpk_zammad_vk_bot/tree/master) [![Docker Hub Container](https://img.shields.io/badge/Docker%20Hub-Container-blue?logo=docker)](https://hub.docker.com/r/chazari/zammad-vk-bot)

# Zammad VK Bot

Zammad VK Bot - это интеграция между платформой обслуживания клиентов Zammad и социальной сетью VK (ВКонтакте), написанная на языке программирования Golang. Этот бот предназначен для упрощения взаимодействия пользователей с системой Zammad через VK, обеспечивая более удобный способ отправки запросов и получения поддержки.

## Настройка Zammad

### Zammad пользователь для VK Bot

Для корректной работы бота необходимо создать специального пользователя в системе Zammad с правами администратора. Этот пользователь `должен использоваться исключительно ботом` для взаимодействия с системой.

<div style="border-left: 4px solid #007bff; padding-left: 10px;">
    <i>Создание отдельного пользователя для бота позволит избежать случайного срабатывания триггеров в системе Zammad из-за действий обычных пользователей.</i>
</div>

Пользователь бота `должен иметь доступ ко всем группам, под которыми создаются тикеты`.

<div style="border-left: 4px solid #007bff; padding-left: 10px;">
    <i>Это гарантирует, что бот сможет эффективно создавать и обрабатывать тикеты без каких-либо ограничений по доступу к необходимым данным в системе Zammad.</i>
</div>

### Доступ по ключу к API Zammad `/#profile/token_access`

Создайте личный ключ доступа для приложения с указанными `правами доступа`:

- admin.api, 
- admin.channel_chat, 
- admin.channel_sms, 
- admin.group, 
- admin.object, 
- admin.user, 
- admin.webhook, 
- chat.agent, 
- ticket.agent

### Веб-перехватчик (Webhook) `/#manage/webhook`

Вы можете использовать веб-перехватчики (webhooks) в Zammad для отправки данных о заявках, статьях и вложениях при каждом срабатывании триггера. Просто создайте и настройте веб-перехватчик (webhooks) с конечной точкой (endpoint) HTTP(S) и соответствующими параметрами безопасности, а затем настройте триггер для его выполнения.

#### Укажите `конечную точку (endpoint)`: 
``` 
http(s)://host:post/webhook
```

#### Укажите `ключ дподписи hnac sha1`, который указан в конфигурации `docker-compose.yml`:
```
your_webhook_secret_key
```

#### Включите `проверку уровня защищенных сокетов (ssl)`, если конечная точка расположена на `https`, или выключите, если конечная точка расположена на `http`.

### Триггеры `/#manage/trigger`

Создайте триггеры:

```json
[
    {
        "Имя": "botChangeGroup",
        "Активировано": "Действие",
        "Выполнение действия": "Selective (default)",
        "Условия отбора объектов": {
            "Заявка.Группа": "Изменился"
        },
        "Применить изменения к объектам": {
            "Веб-перехватчик (Webhook)": "установленный webhook"
        }
    },
    {
        "Имя": "botChangeOwner",
        "Активировано": "Действие",
        "Выполнение действия": "Selective (default)",
        "Условия отбора объектов": {
            "Заявка.Владелец": "Изменился"
        },
        "Применить изменения к объектам": {
            "Веб-перехватчик (Webhook)": "установленный webhook"
        }
    },
    {
        "Имя": "botChangePriority",
        "Активировано": "Действие",
        "Выполнение действия": "Selective (default)",
        "Условия отбора объектов": {
            "Заявка.Приоритет": "Изменился"
        },
        "Применить изменения к объектам": {
            "Веб-перехватчик (Webhook)": "установленный webhook"
        }
    },
    {
        "Имя": "botChangeStatus",
        "Активировано": "Действие",
        "Выполнение действия": "Selective (default)",
        "Условия отбора объектов": {
            "Заявка.Состояние": "Изменился"
        },
        "Применить изменения к объектам": {
            "Веб-перехватчик (Webhook)": "установленный webhook"
        }
    },
    {
        "Имя": "botChangeTitle",
        "Активировано": "Действие",
        "Выполнение действия": "Selective (default)",
        "Условия отбора объектов": {
            "Заявка.Заголовок": "Изменился"
        },
        "Применить изменения к объектам": {
            "Веб-перехватчик (Webhook)": "установленный webhook"
        }
    },
    {
        "Имя": "botNewMessage",
        "Активировано": "Действие",
        "Выполнение действия": "Selective (default)",
        "Условия отбора объектов": {
            "Заявка.Действие": "is Обновление",
            "Статья.Отправитель": "is Агент"
        },
        "Применить изменения к объектам": {
            "Веб-перехватчик (Webhook)": "установленный webhook"
        }
    },
    {
        "Имя": "botNewTicket",
        "Активировано": "Действие",
        "Выполнение действия": "Selective (default)",
        "Условия отбора объектов": {
            "Статья.Действие": "is Создано",
            "Заявка.Действие": "не Обновлен"
        },
        "Применить изменения к объектам": {
            "Веб-перехватчик (Webhook)": "установленный webhook"
        }
    }
]
```

## Запуск через Docker Compose

Для запуска Zammad VK Bot вы можете использовать Docker Compose. Ниже приведен пример файла `docker-compose.yml`, который вы можете использовать:

```yaml
version: '3'
services:
  vk-bot:
    container_name: vk-bot
    image: chazari/zammad-vk-bot:latest
    environment:
      VK_TOKEN: your_vk_token
      VK_API_HREF: https://api.vk.com
      ZAMMAD_TOKEN: your_zammad_token
      ZAMMAD_HREF: your_zammad_url
      WEBHOOK_SECRET_KEY: your_webhook_secret_key
      WEBHOOK_PORT: :8181
      POSTGRES_DB: postgres
      POSTGRES_USER: postgres
      POSTGRES_PORT: "5432"
      POSTGRES_HOST: postgres
      POSTGRES_PASS: your_postgres_password
    restart: always
    command: ["/app/main", "vk-bot"]
  postgres:
    container_name: postgres
    image: postgres:latest
    environment:
      POSTGRES_DB: postgres
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: your_postgres_password
    restart: on-failure
    ports:
      - "5432:5432"
```

После создания файла `docker-compose.yml` и замены параметров на свои, вы можете запустить Zammad VK Bot с помощью команды:

```bash
docker-compose up -d
```

## Взаимодействие с ботом

Пользователи могут взаимодействовать с ботом, отправляя сообщения с запросами на поддержку. 
Бот автоматически создаст тикет в Zammad на основе полученного сообщения. 

<div style="border-left: 4px solid #007bff; padding-left: 10px;">
    <i>В некоторых случаях бот может попросить вас предоставить логин и пароль от вашей учетной записи Zammad, чтобы создать тикет.</i>
</div>

### Пример создания тикета

Для создания тикета достаточно просто написать сообщение с описанием вашей проблемы. Например:

```
Привет, не могу войти в систему. Помогите, пожалуйста!
```

## Ограничения

- Бот может обрабатывать только определенные типы запросов, описанные в его функциональных возможностях.
- Необходима стабильная сетевая связь для взаимодействия с VK и Zammad.

## Лицензия

Этот проект лицензирован под MIT License. См. файл [LICENSE](LICENSE) для получения дополнительной информации.