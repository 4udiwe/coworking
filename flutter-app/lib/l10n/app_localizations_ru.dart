// ignore: unused_import
import 'package:intl/intl.dart' as intl;
import 'app_localizations.dart';

// ignore_for_file: type=lint

/// The translations for Russian (`ru`).
class AppLocalizationsRu extends AppLocalizations {
  AppLocalizationsRu([String locale = 'ru']) : super(locale);

  @override
  String get appTitle => 'Приложение коворкинга';

  @override
  String get loginTitle => 'Вход';

  @override
  String get loginButtonLabel => 'Вход';

  @override
  String get loginSuccess => 'Вход успешен';

  @override
  String get emailFieldLabel => 'Email';

  @override
  String get emailFieldHint => 'Введите email';

  @override
  String get passwordFieldLabel => 'Пароль';

  @override
  String get passwordFieldHint => 'Введите пароль';

  @override
  String get registerTitle => 'Регистрация';

  @override
  String get registerButtonLabel => 'Регистрация';

  @override
  String get registerSuccess => 'Регистрация успешна';

  @override
  String get firstNameFieldLabel => 'Имя';

  @override
  String get firstNameFieldHint => 'Введите имя';

  @override
  String get lastNameFieldLabel => 'Фамилия';

  @override
  String get lastNameFieldHint => 'Введите фамилию';

  @override
  String get confirmPasswordFieldLabel => 'Подтвердите пароль';

  @override
  String get confirmPasswordFieldHint => 'Пароль';

  @override
  String get validationEmptyEmail => 'Введите email';

  @override
  String get validationInvalidEmail => 'Некорректный email';

  @override
  String get validationMinimumPassword => 'Минимум 8 символов';

  @override
  String get validationEmptyFirstName => 'Введите имя';

  @override
  String get validationEmptyLastName => 'Введите фамилию';

  @override
  String get validationPasswordsNotMatch => 'Пароли не совпадают';

  @override
  String get validationInvalidData => 'Некорректные данные';

  @override
  String get haveNoAccount => 'Нет аккаунта?';

  @override
  String get toRegister => 'Зарегистрироваться';

  @override
  String get alreadyHaveAccount => 'Уже есть аккаунт?';

  @override
  String get incorrectData => 'Некорректные данные';

  @override
  String get navCoworkings => 'Коворкинги';

  @override
  String get navBookings => 'Бронирования';

  @override
  String get navNotifications => 'Уведомления';

  @override
  String get navProfile => 'Профиль';

  @override
  String get navAdmin => 'Админ';

  @override
  String get profileScreenTitle => 'Профиль';

  @override
  String get bookingsScreenTitle => 'Бронирования';

  @override
  String get bookingsTabActive => 'Активные';

  @override
  String get bookingsTabHistory => 'История';

  @override
  String get bookingLabelCoworking => 'Коворкинг:';

  @override
  String get bookingLabelPlace => 'Место:';

  @override
  String get bookingLabelCancelReason => 'Причина отмены:';

  @override
  String get bookingButtonCancel => 'Отменить бронирование';

  @override
  String get cancelReasonByUser => 'Отменено пользователем';

  @override
  String get failedFetchActiveBookings => 'Не удалось получить активные бронирования';

  @override
  String get failedFetchHistoryBookings => 'Не удалось получить историю бронирований';

  @override
  String get bookingCancelled => 'Бронирование отменено';

  @override
  String get failedCancelBooking => 'Не удалось отменить бронирование';

  @override
  String get coworkingsScreenTitle => 'Коворкинги';

  @override
  String get coworkingDefaultName => 'Коворкинг';

  @override
  String get coworkingUnavailable => 'Коворкинг недоступен сейчас';

  @override
  String get bookingCreatedSuccess => 'Бронирование успешно создано';

  @override
  String get adminPanelTitle => 'Панель администратора';

  @override
  String get adminNoCoworkings => 'Нет коворкингов';

  @override
  String get adminFailedLoadCoworkings => 'Не удалось загрузить коворкинги';

  @override
  String get errorUserAlreadyExists => 'Пользователь уже существует';

  @override
  String get errorRegistrationFailed => 'Ошибка регистрации';

  @override
  String get errorInvalidCredentials => 'Неверные учетные данные';

  @override
  String get errorLoginFailed => 'Ошибка входа';

  @override
  String get errorRefreshTokenFailed => 'Ошибка обновления токена';

  @override
  String get errorLogoutFailed => 'Ошибка выхода';

  @override
  String get errorUnauthorized => 'Не авторизовано';

  @override
  String get errorFetchProfile => 'Не удалось загрузить профиль';

  @override
  String get errorFetchActiveSessions => 'Не удалось загрузить активные сессии';

  @override
  String get errorFetchAllSessions => 'Не удалось загрузить все сессии';

  @override
  String get errorRevokeSession => 'Не удалось отозвать сессию';

  @override
  String get errorFetchCoworkings => 'Не удалось загрузить коворкинги';

  @override
  String get errorFetchCoworkingDetails => 'Не удалось загрузить информацию о коворкинге';

  @override
  String get errorFetchAvailablePlaces => 'Не удалось загрузить доступные места';

  @override
  String get errorFetchLayout => 'Не удалось загрузить макет коворкинга';

  @override
  String get errorCreateBooking => 'Не удалось создать бронирование';

  @override
  String get errorCancelBooking => 'Не удалось отменить бронирование';

  @override
  String get settingsLanguage => 'Язык';

  @override
  String get settingsSystemDefault => 'По умолчанию системы';

  @override
  String get settingsEnglish => 'English';

  @override
  String get settingsRussian => 'Русский';
}
