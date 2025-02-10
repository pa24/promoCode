# Сервер для управления промокодами

Это серверное приложение для управления промокодами в игре, написанное на Go. Оно использует Docker для контейнеризации и PostgreSQL для хранения данных. 

1. Склонируйте проект на вашу машину: https://github.com/pa24/promoCode.git
2. Перейдите в рабочую директорию
```bash
 cd promoCode
```
4. Для того чтобы запустить приложение в Docker, используйте команду:
```bash
docker-compose up --build
```
5. Если контейнеры были успешно запущены, сервер будет доступен по адресу http://localhost:8080
6.   Интерфейс для создания промокода доступен по адресу http://localhost:8080/admin
7. Миграции базы данных будут применяться автоматически при старте приложения

| Метод | URL           | Описание                                             |
|-------|---------------|------------------------------------------------------|
| POST  | /admin/create | Создание нового промокода                            |
| GET   | /admin        | Страница администратора для создания промокодов      |
| POST  | /api/apply    | Применение промокода     

## Примеры запросов

### 1. Создание промокода

**Параметры запроса:**

| Параметр  | Тип     | Описание                              |
|-----------|---------|---------------------------------------|
| `code`    | string  | Название промокода                    |
| `reward`  | number  | Вознаграждение                        |
| `max_uses`| number  | Максимальное количество использований |
| `player_id`| number | Id игрока |

**Пример запроса:**
1. Запрос на создание промокода

```bash
curl -X POST http://localhost:8080/admin/create \
     -H "Content-Type: application/x-www-form-urlencoded" \
     -d "code=FREE100&reward=100&max_uses=10"
```

Ответ (успешное создание):
```json
{
  "message": "Promo code created successfully"
}
```

2. Применение промокода
```bash
curl -X POST http://localhost:8080/api/apply \
     -H "Content-Type: application/json" \
     -d '{
           "player_id": 123,
           "code": "FREE100"
         }'
```

Ответ (Успешное использование промокода):
```json
{
  "message": "Promocode applied successfully"
}
```
