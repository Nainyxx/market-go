# market-go

Мой учебный проект: маркетплейс на Go, который я делаю микросервисами. Цель простая: научиться писать бэкенд так, как его делают в командах, а не одним большим `main.go`.

Пока готовы два сервиса, остальные в планах.

## Что уже есть

**AuthService**: регистрация и вход.
- `POST /auth/register` создаёт пользователя и его профиль
- `POST /auth/login` проверяет пароль (bcrypt) и выдаёт JWT

**UserService**: профиль пользователя.
- имя, фамилия, аватар
- баланс (хранится в копейках, `int64`) и рейтинг
- списание баланса идёт одним SQL-запросом (`UPDATE ... WHERE balance >= $1`), поэтому два параллельных запроса не уведут баланс в минус

## Как это работает вместе?

При регистрации AuthService сначала создаёт пользователя у себя, потом вызывает UserService, чтобы тот создал профиль с тем же `id`. Общей транзакции между двумя базами нет, поэтому если профиль создать не вышло, AuthService откатывает своего пользователя. Внутренние вызовы между сервисами защищены токеном.

В каждом сервисе одни и те же слои: `handler → service → repo → models`. Зависимости передаются через интерфейсы, ошибки оборачиваются через `%w` и на границе превращаются в HTTP-статусы.

## Стек

Go, Gin, PostgreSQL, JWT, bcrypt.

## Как запустить

1. Поднять PostgreSQL и создать по базе на каждый сервис.
2. Применить миграции:
```bash
   psql -d <база_auth> -f services/AuthService/migration/001.sql
   psql -d <база_user> -f services/UserService/migration/001.sql
```
3. Скопировать `.env.example` в `.env` в каждом сервисе и заполнить.
4. Запустить:
```bash
   cd services/UserService && go run .
   cd services/AuthService && go run .
```

Пример:
```bash
curl -X POST localhost:8080/auth/register -H 'Content-Type: application/json' \
  -d '{"name":"Иван","surname":"Петров","phone":"+79991234567","email":"ivan@example.com","password":"correct-horse"}'
```

## Что дальше

- [ ] тесты на сервисный слой
- [ ] Product Service
- [ ] API-шлюз (Envoy) и docker-compose
