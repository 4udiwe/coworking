# Система локализации Coworking App

## 📝 Описание
Приложение поддерживает локализацию на русском и английском языках с возможностью:
- **Авто-определения** языка на основе системной локали устройства
- **Ручного выбора** языка через настройки приложения
- **Сохранение выбора** пользователя в защищенном хранилище
- **Fallback на английский** если система не определена (браузер, неизвестная локаль)

## 📁 Структура файлов

```
lib/
├── l10n/
│   ├── app_en.arb         # Английские тексты
│   └── app_ru.arb         # Русские тексты
├── generator_l10n/
│   ├── app_localizations.dart       # Основной класс
│   ├── app_localizations_en.dart    # Реализация англ.
│   └── app_localizations_ru.dart    # Реализация русс.
├── core/l10n/
│   ├── locale_provider.dart         # Управление локалью
│   └── localization_helper.dart     # Helper для доступа к строкам
└── main.dart                        # Конфигурация приложения

l10n.yaml                            # Конфигурация локализации
```

## 🔧 Как использовать локализованные строки

### В widget's методах build:

```dart
import 'package:coworking_app/core/l10n/localization_helper.dart';

// В методе build:
Text(context.l10n.loginButtonLabel)
Text(context.l10n.emailFieldHint)
Text(context.l10n.errorUserAlreadyExists)
```

### Изменение языка программно:

```dart
final localeProvider = context.read<LocaleProvider>();

// Переключить на русский
await localeProvider.setLocale(const Locale('ru'));

// Переключить на английский
await localeProvider.setLocale(const Locale('en'));

// Использовать системную локаль
await localeProvider.setLocale(null);

// Переключить на следующий
await localeProvider.toggleLocale();
```

## 🌐 Как добавить новый язык

1. Создать новый файл `lib/l10n/app_xx.arb` (где xx - код языка)
2. Скопировать структуру из `app_en.arb`
3. Перевести все значения
4. Обновить `l10n.yaml`:
   ```yaml
   preferred-supported-locales:
     - en
     - ru
     - xx  # Добавить новый код
   ```
5. Запустить: `flutter gen-l10n`
6. Обновить `LocaleProvider.supportedLocales`

## ⚙️ Настройка системной локали

Если нужно, чтобы приложение всегда определяло язык по системе:

```dart
// В LocaleProvider.initialize()
// Будет автоматически использоваться языковой код системы
```

Система поддерживает любой языковой код `Locale`, но покажет текст только для поддерживаемых локалей (en, ru).

## 🎯 Инициализация в main.dart

```dart
final localeProvider = LocaleProvider();
await localeProvider.initialize();

runApp(
  ChangeNotifierProvider(
    create: (_) => localeProvider,
    child: MultiBlocProvider(
      // ... providers
      child: const MyApp(),
    ),
  ),
);
```

## 💾 Сохранение выбора
Выбор языка пользователя сохраняется в `flutter_secure_storage` и восстанавливается при загрузке приложения.
