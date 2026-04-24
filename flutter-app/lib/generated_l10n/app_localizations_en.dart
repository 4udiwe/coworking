// ignore: unused_import
import 'package:intl/intl.dart' as intl;
import 'app_localizations.dart';

// ignore_for_file: type=lint

/// The translations for English (`en`).
class AppLocalizationsEn extends AppLocalizations {
  AppLocalizationsEn([String locale = 'en']) : super(locale);

  @override
  String get appTitle => 'Coworking App';

  @override
  String get loginTitle => 'Login';

  @override
  String get loginButtonLabel => 'Login';

  @override
  String get loginSuccess => 'Login success';

  @override
  String get emailFieldLabel => 'Email';

  @override
  String get emailFieldHint => 'Enter email';

  @override
  String get passwordFieldLabel => 'Password';

  @override
  String get passwordFieldHint => 'Enter password';

  @override
  String get registerTitle => 'Register';

  @override
  String get registerButtonLabel => 'Register';

  @override
  String get registerSuccess => 'Registration success';

  @override
  String get firstNameFieldLabel => 'First name';

  @override
  String get firstNameFieldHint => 'Enter first name';

  @override
  String get lastNameFieldLabel => 'Last name';

  @override
  String get lastNameFieldHint => 'Enter last name';

  @override
  String get confirmPasswordFieldLabel => 'Confirm password';

  @override
  String get confirmPasswordFieldHint => 'Confirm password';

  @override
  String get validationEmptyEmail => 'Enter email';

  @override
  String get validationInvalidEmail => 'Invalid email';

  @override
  String get validationMinimumPassword => 'Minimum 8 characters';

  @override
  String get validationEmptyFirstName => 'Enter first name';

  @override
  String get validationEmptyLastName => 'Enter last name';

  @override
  String get validationPasswordsNotMatch => 'Passwords don\'t match';

  @override
  String get validationInvalidData => 'Invalid data';

  @override
  String get haveNoAccount => 'Have no account?';

  @override
  String get toRegister => 'To register';

  @override
  String get alreadyHaveAccount => 'Already have an account?';

  @override
  String get incorrectData => 'Incorrect data';

  @override
  String get navCoworkings => 'Coworkings';

  @override
  String get navBookings => 'Bookings';

  @override
  String get navNotifications => 'Notifications';

  @override
  String get navProfile => 'Profile';

  @override
  String get navAdmin => 'Admin';

  @override
  String get profileScreenTitle => 'Profile';

  @override
  String get bookingsScreenTitle => 'Bookings';

  @override
  String get bookingsTabActive => 'Active';

  @override
  String get bookingsTabHistory => 'History';

  @override
  String get bookingLabelCoworking => 'Coworking:';

  @override
  String get bookingLabelPlace => 'Place:';

  @override
  String get bookingLabelCancelReason => 'Cancel reason:';

  @override
  String get bookingButtonCancel => 'Cancel booking';

  @override
  String get cancelReasonByUser => 'Cancelled by user';

  @override
  String get failedFetchActiveBookings => 'Failed to fetch active bookings';

  @override
  String get failedFetchHistoryBookings => 'Failed to fetch history bookings';

  @override
  String get bookingCancelled => 'Booking cancelled';

  @override
  String get failedCancelBooking => 'Failed to cancel booking';

  @override
  String get coworkingsScreenTitle => 'Coworkings';

  @override
  String get coworkingDefaultName => 'Coworking';

  @override
  String get coworkingUnavailable => 'Coworking is unavailable right now';

  @override
  String get bookingCreatedSuccess => 'Booking created successfully';

  @override
  String get adminPanelTitle => 'Admin Panel';

  @override
  String get adminNoCoworkings => 'No coworkings';

  @override
  String get adminFailedLoadCoworkings => 'Failed to load coworkings';

  @override
  String get errorUserAlreadyExists => 'User already exists';

  @override
  String get errorRegistrationFailed => 'Registration failed';

  @override
  String get errorInvalidCredentials => 'Invalid credentials';

  @override
  String get errorLoginFailed => 'Login failed';

  @override
  String get errorRefreshTokenFailed => 'Refresh token failed';

  @override
  String get errorLogoutFailed => 'Logout failed';

  @override
  String get errorUnauthorized => 'Unauthorized';

  @override
  String get errorFetchProfile => 'Failed to fetch profile';

  @override
  String get errorFetchActiveSessions => 'Failed to fetch active sessions';

  @override
  String get errorFetchAllSessions => 'Failed to fetch all sessions';

  @override
  String get errorRevokeSession => 'Failed to revoke session';

  @override
  String get errorFetchCoworkings => 'Failed to fetch coworkings';

  @override
  String get errorFetchCoworkingDetails => 'Failed to fetch coworking details';

  @override
  String get errorFetchAvailablePlaces => 'Failed to fetch available places';

  @override
  String get errorFetchLayout => 'Failed to fetch layout for coworking';

  @override
  String get errorCreateBooking => 'Failed to create booking';

  @override
  String get errorCancelBooking => 'Failed to cancel booking';

  @override
  String get settingsLanguage => 'Language';

  @override
  String get settingsSystemDefault => 'System default';

  @override
  String get settingsEnglish => 'English';

  @override
  String get settingsRussian => 'Русский';

  @override
  String get currentSession => 'current';

  @override
  String get activeSession => 'active';

  @override
  String get revokedSession => 'revoked';

  @override
  String get sessionCreated => 'Created: ';

  @override
  String get sessionLastUsedAt => 'Last used at: ';

  @override
  String get sessionExpiresAt => 'Expires at: ';

  @override
  String get revokeSession => 'Revoke session';

  @override
  String get sessions => 'Sessions';

  @override
  String get notifications => 'Notifications';

  @override
  String get markAllAsRead => 'Mark all as read';

  @override
  String get notificationReadAt => 'Read';

  @override
  String get bookingStatusActive => 'Active';

  @override
  String get bookingStatusCompleted => 'Completed';

  @override
  String get bookingStatusCancelled => 'Cancelled';

  @override
  String get placeTypeOpenDesk => 'Open desk';

  @override
  String get noActiveBookings => 'No active bookings';

  @override
  String get noHistoryBookings => 'No booking history';
}
