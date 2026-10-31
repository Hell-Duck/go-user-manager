# Go User CRUD API

Учебный REST API на Go для управления пользователями.

Реализованы создание, получение, обновление и удаление пользователей. Данные хранятся в оперативной памяти через `map`.

## Запуск

Из корневой папки проекта:

```bash
go run .
```

Сервер запускается по адресу `http://localhost:4040`.

## HTTP API

| Метод | Адрес | Описание |
|---|---|---|
| GET | `/user/1` | Получить пользователя |
| POST | `/users` | Создать пользователя |
| PUT | `/user/1` | Обновить пользователя |
| DELETE | `/user/1` | Удалить пользователя |

## Примеры запросов

Команды для PowerShell.

**Создание пользователя**

```powershell
'{"Name":"Alex","Age":25}' | curl.exe -i -X POST http://localhost:4040/users -H "Content-Type: application/json" --data-binary "@-"
```

**Получение пользователя**

```powershell
curl.exe -i http://localhost:4040/user/1
```

**Обновление пользователя**

```powershell
'{"Name":"Bob","Age":30}' | curl.exe -i -X PUT http://localhost:4040/user/1 -H "Content-Type: application/json" --data-binary "@-"
```

**Удаление пользователя**

```powershell
curl.exe -i -X DELETE http://localhost:4040/user/1
```

## Тестирование

```bash
go test ./...
```

## Ограничения

- Данные не сохраняются после перезапуска сервера.
- Хранилище пока не защищено от одновременных запросов.
- База данных и авторизация пока не подключены.