# Cinema Hall

Микросервисная система для управления кинотеатром с возможностью бронирования билетов, управления сеансами, фильмами и залами. Система построена на архитектуре микросервисов с использованием Apache Kafka для асинхронной обработки событий и API Gateway в качестве единой точки входа.

## Функционал приложения

- Регистрация и авторизация пользователей, управление профилем
- Просмотр фильмов и сеансов
- Бронирование билетов, с автоматической отменой через 15 минут, при отсутствии оплаты
- Управление залами и местами
- Асинхронная обработка событий через Kafka

## Технологический стек

![Go](https://img.shields.io/badge/Go-1.25.4-00ADD8?logo=go)
![Gin](https://img.shields.io/badge/Gin-1.9+-00D8FF?logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-336791?logo=postgresql)
![GORM](https://img.shields.io/badge/GORM-1.25+-FF6B6B?logo=go)
![Kafka](https://img.shields.io/badge/Apache%20Kafka-3.5+-231F20?logo=apache-kafka)
![JWT](https://img.shields.io/badge/JWT-Auth-000000?logo=json-web-tokens)
![Docker](https://img.shields.io/badge/Docker-20.10+-2496ED?logo=docker)
![Docker Compose](https://img.shields.io/badge/Docker%20Compose-2.0+-2496ED?logo=docker)

**Архитектура**: Микросервисная (5 сервисов: User, Movie, Cinema, Booking, Gateway)

```
             HTTP      ┌───────────────────┐
┌────────┐ ──────────> │  Gateway Service  │
│ Client │             └─────────┬─────────┘
└────────┘                       │
           ┌─────────────────────┼─────────────────────┐
           │                     │                     │
           v                     v                     v
    ┌──────────────┐      ┌──────────────┐      ┌──────────────┐
    │    Movie     │      │   Booking    │      │    User      │
    │   Service    │      │   Service    │      │   Service    │
    └──────┬───────┘      └──────┬───────┘      └──────┬───────┘
           │                     │                     │
      PostgreSQL                 │ HTTP           PostgreSQL
                                 │ (get seats,
                                 │  check availability)
                                 v
                          ┌──────────────┐
                          │   Cinema     │
                          │   Service    │
                          └──────┬───────┘
                                 │
                            PostgreSQL
                                 │
                          Kafka Topics
                     ┌─────────────────────┐
                     │ booking.confirmed   │
                     │ booking.cancelled   │
                     └─────────────────────┘
```


## Участники разработки

- [Али Умаров](https://github.com/AliUmarov)
- [Усман Дзакаев](https://github.com/dzakaev)
- [Бекхан Хаджимагомадов](https://github.com/Bekkhanbs)
- [Вис Магомадов](https://github.com/magadov)
