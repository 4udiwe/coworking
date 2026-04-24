import 'package:flutter/material.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class LocaleProvider extends ChangeNotifier {
  static const String _storageKey = 'app_locale';
  static const String _systemDefault = 'system';
  static const List<Locale> supportedLocales = [Locale('en'), Locale('ru')];

  final FlutterSecureStorage _storage = const FlutterSecureStorage();

  Locale? _currentLocale;
  bool _isInitialized = false;

  Locale? get currentLocale => _currentLocale;
  bool get isInitialized => _isInitialized;

  /// Инициализирует локаль из хранилища или из системы
  Future<void> initialize() async {
    try {
      final savedLocale = await _storage.read(key: _storageKey);

      if (savedLocale == null || savedLocale == _systemDefault) {
        // Используем системную локаль
        _currentLocale = null; // null позволит использовать системную локаль
      } else {
        // Используем сохраненную локаль
        final parts = savedLocale.split('_');
        if (parts.length == 2) {
          final locale = Locale(parts[0], parts[1]);
          // Проверяем, поддерживается ли эта локаль
          _currentLocale = supportedLocales.contains(locale) ? locale : null;
        } else if (parts.isNotEmpty) {
          final locale = Locale(parts[0]);
          // Проверяем, поддерживается ли эта локаль
          _currentLocale = supportedLocales.contains(locale) ? locale : null;
        }
      }
    } catch (e) {
      // В случае ошибки используем системную локаль
      _currentLocale = null;
    }

    _isInitialized = true;
    notifyListeners();
  }

  /// Устанавливает локаль приложения
  /// Если [locale] == null, используется системная локаль
  Future<void> setLocale(Locale? locale) async {
    _currentLocale = locale;

    try {
      if (locale == null) {
        await _storage.write(key: _storageKey, value: _systemDefault);
      } else {
        final localeString = locale.countryCode != null
            ? '${locale.languageCode}_${locale.countryCode}'
            : locale.languageCode;
        await _storage.write(key: _storageKey, value: localeString);
      }
    } catch (e) {
      // Ошибка сохранения, но локаль все равно изменится в памяти
      debugPrint('Error saving locale preference: $e');
    }

    notifyListeners();
  }

  /// Переключает между русским, английским и системной локалью
  Future<void> toggleLocale() async {
    if (_currentLocale == null) {
      // Система -> Русский
      await setLocale(const Locale('ru'));
    } else if (_currentLocale!.languageCode == 'ru') {
      // Русский -> Английский
      await setLocale(const Locale('en'));
    } else {
      // Английский -> Система
      await setLocale(null);
    }
  }

  /// Получить текущий языковой код
  String get currentLanguageCode {
    if (_currentLocale == null) {
      return 'system';
    }
    return _currentLocale!.languageCode;
  }
}
