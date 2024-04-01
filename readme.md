[![GitHub Repository](https://img.shields.io/badge/GitHub-Repository-blue?logo=github)](https://github.com/chazari-x/hmtpk_zammad_vk_bot/tree/master) [![Docker Hub Container](https://img.shields.io/badge/Docker%20Hub-Container-blue?logo=docker)](https://hub.docker.com/r/chazari/zammad-vk-bot)

# Zammad VK Bot

Zammad VK Bot - это интеграция между платформой обслуживания клиентов Zammad и социальной сетью VK (ВКонтакте), написанная на языке программирования Golang. Этот бот предназначен для упрощения взаимодействия пользователей с системой Zammad через VK, обеспечивая более удобный способ отправки запросов и получения поддержки.

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

Пользователи могут взаимодействовать с ботом, отправляя сообщения с запросами на поддержку. Бот автоматически создаст тикет в Zammad на основе полученного сообщения. В некоторых случаях бот может попросить вас предоставить логин и пароль от вашей учетной записи Zammad, чтобы создать тикет.

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