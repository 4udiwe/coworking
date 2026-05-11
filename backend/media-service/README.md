# Media Service

Сервис хранения и обработки изображений. Отвечает за загрузку, хранение, асинхронный ресайз и выдачу изображений коворкингов.

Сервис полностью изолирован от бизнес-логики коворкингов и работает как отдельное media-хранилище.

## Архитектура

### MongoDB
MongoDB хранит metadata изображения:

- статус обработки
- mime type
- размеры variants
- retry count
- timestamps

Пример документа:
```json
{
  "_id": "69fee84b6b7056213f41d693",
  "status": "ready",
  "mime_type": "image/webp",
  "variants": [
    {
      "size": "thumbnail",
      "width": 300
    },
    {
      "size": "medium",
      "width": 900
    },
    {
      "size": "large",
      "width": 1600
    }
  ],
  "created_at": "...",
  "updated_at": "..."
}
```

### MinIO

Файлы изображений хранятся в MinIO.

Структура ключей:
```
media/{media_id}/original.webp
media/{media_id}/thumbnail.webp
media/{media_id}/medium.webp
media/{media_id}/large.webp
```
MinIO bucket используется как публичное CDN-хранилище.

## Pipeline обработки

### Upload
При загрузке изображения:

- создаётся запись в MongoDB
- original загружается в MinIO
- синхронно создаётся thumbnail
- thumbnail загружается в MinIO
- статус переводится в processing
- запускается async resize остальных размеров

Ответ upload endpoint:
```json
{
  "id": "69fee84b6b7056213f41d693",
  "status": "processing",
  "urls": {
    "thumbnail": "/media/69fee84b6b7056213f41d693/thumbnail.webp"
  }
}
```

### Async processing
Фоновая обработка:

- скачивает original
- параллельно генерирует размеры
- загружает variants в MinIO
- обновляет MongoDB

Используется errgroup для конкурентной обработки.

### Stale Processing Recovery

Если обработка зависла:

- stale worker находит медиа в статусе processing
- увеличивает retry_count
- повторно запускает resize

Если число retry превышено:

```
status = failed
```

### Статусы Media
- **pending** - Файл создан, upload ещё не завершён.
- **processing** - Идёт async генерация размеров.
- **ready** - Все размеры успешно созданы.
- **failed** - Обработка завершилась ошибкой.

### Размеры изображений

Поддерживаются variants:

| Размер    | Назначение              |
| --------- | ----------------------- |
| thumbnail | карточки, preview       |
| medium    | списки, галереи         |
| large     | full-screen отображение |
| original  | исходное изображение    |

## API

### Upload 
```json
POST /api/media/upload
```
Multipart upload изображения.

### Delete
```
DELETE /api/media/{id}
```
Удаляет metadata и файлы из storage.

Подробнее в [swagger](../docs/swagger.yaml).

## Публичная выдача файлов
Изображения отдаются напрямую через gateway → MinIO:
```
GET /media/{media_id}/thumbnail.webp
GET /media/{media_id}/medium.webp
GET /media/{media_id}/large.webp
```
Media-service не участвует в выдаче файлов.

## Интеграция
Использует публичный RSA-ключ для валидации access token входящих HTTP запросов.
Путь к файлу с ключем (pubilc.pem) указывается в [config.yaml](config/config.yaml)


## Конфигурация

Через `.env`. Для запуска можно скопировать `.env.example`