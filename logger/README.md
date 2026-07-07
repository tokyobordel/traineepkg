# ContextLogger - Логгер с обязательным контекстом

## Описание

`ContextLogger` - это обертка над `logrus.Logger`, которая требует обязательный параметр `context.Context` для всех методов логирования. Логгер автоматически извлекает `traceSpread` из контекста и добавляет его в каждое сообщение лога.

## Особенности

- **Обязательный контекст**: Все методы логирования требуют `context.Context` как первый параметр
- **Автоматическое извлечение spread**: Логгер автоматически извлекает `traceSpread` из контекста
- **Совместимость с logrus**: Поддерживает все уровни логирования logrus
- **Дополнительные поля**: Поддержка `WithField` и `WithFields` с контекстом

## Использование

### Создание логгера

```go
// Создание нового ContextLogger
logger := logger.NewContextLogger()

// Или создание из существующего logrus.Logger
existingLogger := logrus.New()
contextLogger := logger.NewContextLoggerFromLogger(existingLogger)
```

### Добавление spread в контекст

```go
import "traineepkg/context/trace"

// Создание контекста со spread
ctx := context.Background()
ctx = trace.WithSpread(ctx, "custom-spread-id")

// Автоматическая генерация spread
ctx = trace.WithSpread(ctx, "") // пустая строка приведет к автогенерации
```

### Методы логирования

```go
// Основные методы
logger.Info(ctx, "Информационное сообщение")
logger.Debug(ctx, "Отладочное сообщение")
logger.Warn(ctx, "Предупреждение")
logger.Error(ctx, "Ошибка")
logger.Fatal(ctx, "Критическая ошибка")
logger.Panic(ctx, "Паника")

// Форматированные методы
logger.Infof(ctx, "Пользователь %s создан с ID %d", "john", 123)
logger.Errorf(ctx, "Ошибка: %s", err.Error())
```

### Дополнительные поля

```go
// Одно поле
entry := logger.WithField(ctx, "operation", "user_creation")
entry.Info("Операция начата")

// Несколько полей
entry := logger.WithFields(ctx, logrus.Fields{
    "user_id": 123,
    "action": "login",
})
entry.Info("Пользователь вошел в систему")
```

## Формат вывода

Логгер выводит сообщения в следующем формате:
```
>> [2024-01-01 15:04:05] [INFO   ] [file.go:42 package.function] [spr-your-spread-id] msg:> Ваше сообщение
```

Где:
- `spr-your-spread-id` - это spread из контекста
- Если spread не найден, будет выведено `spr-Spread not found`

## Миграция с обычного logrus.Logger

### Было:
```go
type Repository struct {
    logger *logrus.Logger
}

func (r *Repository) SomeMethod() {
    r.logger.Info("Сообщение")
    r.logger.Errorf("Ошибка: %s", err.Error())
}
```

### Стало:
```go
type Repository struct {
    logger *logger.ContextLogger
}

func (r *Repository) SomeMethod(ctx context.Context) {
    r.logger.Info(ctx, "Сообщение")
    r.logger.Errorf(ctx, "Ошибка: %s", err.Error())
}
```

## Требования

- Все методы, использующие логгер, должны принимать `context.Context`
- Рекомендуется добавлять spread в контекст на входе в приложение
- При отсутствии spread в контексте, логгер все равно будет работать

## Примеры

См. файл `example_usage.go` для подробных примеров использования. 