# API Documentation

## Base URL
```
http://localhost:8080
```

## 1. Пользователи (Users)

### Регистрация пользователя
```http
POST /register
Content-Type: application/json

{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
}
```

### Вход в систему
```http
POST /login
Content-Type: application/json

{
    "email": "test@example.com",
    "password": "password123"
}
```

### Получение всех пользователей
```http
GET /users
```

### Получение пользователя по ID
```http
GET /users/{id}
Authorization: Bearer {token}
```

### Обновление пользователя
```http
PUT /users/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
    "username": "updated_username",
    "email": "updated@example.com"
}
```

### Удаление пользователя
```http
DELETE /users/{id}
Authorization: Bearer {token}
```

## 2. Продукты (Products)

### Создание продукта
```http
POST /products
Content-Type: application/json

{
    "name": "Test Product",
    "description": "Test Description",
    "price": 99.99,
    "stock": 100
}
```

### Получение всех продуктов
```http
GET /products
```

### Получение продукта по ID
```http
GET /products/{id}
```

### Обновление продукта
```http
PUT /products/{id}
Content-Type: application/json

{
    "name": "Updated Product",
    "description": "Updated Description",
    "price": 149.99,
    "stock": 150
}
```

### Удаление продукта
```http
DELETE /products/{id}
```

## 3. Заказы (Orders)

### Создание заказа
```http
POST /orders
Authorization: Bearer {token}
Content-Type: application/json

{
    "user_id": "user_id_here",
    "total_price": 199.99
}
```

### Получение заказа по ID
```http
GET /orders/{id}
Authorization: Bearer {token}
```

### Получение всех заказов
```http
GET /orders
Authorization: Bearer {token}
```

### Обновление заказа
```http
PUT /orders/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
    "status": "completed",
    "total_price": 199.99
}
```

### Удаление заказа
```http
DELETE /orders/{id}
Authorization: Bearer {token}
```

## 4. Корзина (Cart)

### Добавление товара в корзину
```http
POST /cart/{userID}
Content-Type: application/json

{
    "product_id": "product_id_here",
    "quantity": 2
}
```

### Просмотр корзины
```http
GET /cart/{userID}
```

### Удаление товара из корзины
```http
DELETE /cart/{userID}/{productID}
```

### Оформление заказа из корзины
```http
POST /checkout/{userID}
```

## Тестовые данные

### 1. Пользователи
```json
{
    "username": "testuser1",
    "email": "test1@example.com",
    "password": "password123"
}
```

```json
{
    "username": "testuser2",
    "email": "test2@example.com",
    "password": "password123"
}
```

### 2. Продукты
```json
{
    "name": "Laptop",
    "description": "High-performance laptop",
    "price": 999.99,
    "stock": 50
}
```

```json
{
    "name": "Smartphone",
    "description": "Latest model smartphone",
    "price": 699.99,
    "stock": 100
}
```

### 3. Заказы
```json
{
    "user_id": "user_id_here",
    "total_price": 1699.98
}
```

## Порядок тестирования

1. Сначала создайте пользователя через `/register`
2. Получите токен через `/login`
3. Используйте полученный токен в заголовке `Authorization: Bearer {token}` для защищенных эндпоинтов
4. Создайте несколько продуктов
5. Добавьте продукты в корзину
6. Оформите заказ из корзины
7. Проверьте получение заказов и статистики

## Ожидаемые ответы

- Успешные ответы будут иметь статус 200 (GET), 201 (POST), 204 (DELETE)
- Ошибки будут иметь статус 400 (Bad Request), 401 (Unauthorized), 404 (Not Found), 500 (Internal Server Error)
- Все ответы будут в формате JSON

## Примечания

1. Замените `{token}` на реальный JWT токен, полученный после входа
2. Замените `{id}` на реальные ID объектов
3. Все запросы к защищенным эндпоинтам должны содержать валидный JWT токен
4. При тестировании учитывайте, что некоторые операции могут быть кэшированы в Redis

## Примеры ответов

### Успешная регистрация
```json
{
    "id": "user_id",
    "username": "testuser",
    "email": "test@example.com"
}
```

### Успешный вход
```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Ошибка аутентификации
```json
{
    "error": "invalid credentials"
}
```

### Ошибка валидации
```json
{
    "error": "user with email test@example.com already exists"
}
``` 