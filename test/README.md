# Интеграционный тест

Интеграционный тест охватывает полный цикл пути пользователя - от регистрации, до создания, просмотра и отмены бронирования.

## Сценарий

Интеграционный тест проводится по следующему сценарию
1. **User Registration** - Создание пользователя в auth-service
2. **List Coworkings** - Получение доступных коворкингов
3. **Get Available Places** - Получение доступных мест для бронирования в определенное время
4. **Create Booking** - Создание бронирования длительностью в 2 часа
5. **View Active Bookings** - Получение активных бронирований пользователя
6. **Cancel Booking** - Отмена созданного бронирования
7. **Verify Cancellation** - Подтверждение перехода созданного бронирования в статус "отменено"

## Тестовые данные

### Coworking
- **ID**: `550e8400-e29b-41d4-a716-446655440000`
- **Name**: Test Coworking Space
- **Address**: 123 Test Street, Test City, State 12345

### Places
- **Open Desk A**: ID `550e8400-e29b-41d4-a716-446655441001`
- **Meeting Room B**: ID `550e8400-e29b-41d4-a716-446655441002`
- **Private Office C**: ID `550e8400-e29b-41d4-a716-446655441003`

### Test User
- **Email**: test@example.com
- **Password**: password1234
- **Name**: Test User
- **Role**: student

### Время бронирования
- **Start**: Now + 2 hours
- **End**: Now + 4 hours (2-hour duration)


### Запуск теста

```bash
cd test
make test
```

### Сервисы
- **Auth Service**: Порт 8081
- **Booking Service**: Порт 8082 

### Databases
- **postgres_booking_test**: 
  тестовые данные загружаются скриптом `init/booking-seed.sql`

