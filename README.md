# Wallet
REST API для тестового задания.

## Функционал данного приложения:
Приложение позволяет переводить денежные средства с одного кошешлька на другой. При переводе средств осуществляются проверки о том что указанные кошельки существуют и баланс необходимы для списания средств.
Если условия выполняются, то сервер отправляет в очередь задание для сущности Beaver для осуществления транзакции, посредством брокера сообщений RabbitMQ.
Данная сущность блокирует таблицу для изменений и осуществляет обновление записей в БД с нужными суммами.

## Как запустить приложение
1. Запустить файл docker-compose  командой:
```
docker-compose up -d
```
2. Запустить испольняемый файл main.

## API:
1. GET /api/v1/{walletid}

Статусы ответа сервера:
```
200 - запрос успешен
400 - ошибка запрос
404 - кошелек не существует
500 - внутренняя ошибка
```
2. POST /api/v1/{walletid}
Пример тела запроса:
```
{
	"id": "123456",
	"amount": "10.00"
} 
```
Статусы ответа сервера:
```
200 - запрос успешен
400 - ошибка запрос
406 - недостаточно средств на кошельке
404 - кошелек не существует
500 - внутренняя ошибка
```
