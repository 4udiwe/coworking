import 'dart:async';

import 'package:flutter/foundation.dart';
import 'package:flutter/widgets.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:intl/intl.dart' as intl;

import 'app_localizations_en.dart';
import 'app_localizations_ru.dart';

// ignore_for_file: type=lint

/// Callers can lookup localized strings with an instance of AppLocalizations
/// returned by `AppLocalizations.of(context)`.
///
/// Applications need to include `AppLocalizations.delegate()` in their app's
/// `localizationDelegates` list, and the locales they support in the app's
/// `supportedLocales` list. For example:
///
/// ```dart
/// import 'generated_l10n/app_localizations.dart';
///
/// return MaterialApp(
///   localizationsDelegates: AppLocalizations.localizationsDelegates,
///   supportedLocales: AppLocalizations.supportedLocales,
///   home: MyApplicationHome(),
/// );
/// ```
///
/// ## Update pubspec.yaml
///
/// Please make sure to update your pubspec.yaml to include the following
/// packages:
///
/// ```yaml
/// dependencies:
///   # Internationalization support.
///   flutter_localizations:
///     sdk: flutter
///   intl: any # Use the pinned version from flutter_localizations
///
///   # Rest of dependencies
/// ```
///
/// ## iOS Applications
///
/// iOS applications define key application metadata, including supported
/// locales, in an Info.plist file that is built into the application bundle.
/// To configure the locales supported by your app, you’ll need to edit this
/// file.
///
/// First, open your project’s ios/Runner.xcworkspace Xcode workspace file.
/// Then, in the Project Navigator, open the Info.plist file under the Runner
/// project’s Runner folder.
///
/// Next, select the Information Property List item, select Add Item from the
/// Editor menu, then select Localizations from the pop-up menu.
///
/// Select and expand the newly-created Localizations item then, for each
/// locale your application supports, add a new item and select the locale
/// you wish to add from the pop-up menu in the Value field. This list should
/// be consistent with the languages listed in the AppLocalizations.supportedLocales
/// property.
abstract class AppLocalizations {
  AppLocalizations(String locale)
    : localeName = intl.Intl.canonicalizedLocale(locale.toString());

  final String localeName;

  static AppLocalizations? of(BuildContext context) {
    return Localizations.of<AppLocalizations>(context, AppLocalizations);
  }

  static const LocalizationsDelegate<AppLocalizations> delegate =
      _AppLocalizationsDelegate();

  /// A list of this localizations delegate along with the default localizations
  /// delegates.
  ///
  /// Returns a list of localizations delegates containing this delegate along with
  /// GlobalMaterialLocalizations.delegate, GlobalCupertinoLocalizations.delegate,
  /// and GlobalWidgetsLocalizations.delegate.
  ///
  /// Additional delegates can be added by appending to this list in
  /// MaterialApp. This list does not have to be used at all if a custom list
  /// of delegates is preferred or required.
  static const List<LocalizationsDelegate<dynamic>> localizationsDelegates =
      <LocalizationsDelegate<dynamic>>[
        delegate,
        GlobalMaterialLocalizations.delegate,
        GlobalCupertinoLocalizations.delegate,
        GlobalWidgetsLocalizations.delegate,
      ];

  /// A list of this localizations delegate's supported locales.
  static const List<Locale> supportedLocales = <Locale>[
    Locale('en'),
    Locale('ru'),
  ];

  /// No description provided for @appTitle.
  ///
  /// In en, this message translates to:
  /// **'Coworking App'**
  String get appTitle;

  /// No description provided for @loginTitle.
  ///
  /// In en, this message translates to:
  /// **'Login'**
  String get loginTitle;

  /// No description provided for @loginButtonLabel.
  ///
  /// In en, this message translates to:
  /// **'Login'**
  String get loginButtonLabel;

  /// No description provided for @loginSuccess.
  ///
  /// In en, this message translates to:
  /// **'Login success'**
  String get loginSuccess;

  /// No description provided for @emailFieldLabel.
  ///
  /// In en, this message translates to:
  /// **'Email'**
  String get emailFieldLabel;

  /// No description provided for @emailFieldHint.
  ///
  /// In en, this message translates to:
  /// **'Enter email'**
  String get emailFieldHint;

  /// No description provided for @passwordFieldLabel.
  ///
  /// In en, this message translates to:
  /// **'Password'**
  String get passwordFieldLabel;

  /// No description provided for @passwordFieldHint.
  ///
  /// In en, this message translates to:
  /// **'Enter password'**
  String get passwordFieldHint;

  /// No description provided for @registerTitle.
  ///
  /// In en, this message translates to:
  /// **'Register'**
  String get registerTitle;

  /// No description provided for @registerButtonLabel.
  ///
  /// In en, this message translates to:
  /// **'Register'**
  String get registerButtonLabel;

  /// No description provided for @registerSuccess.
  ///
  /// In en, this message translates to:
  /// **'Registration success'**
  String get registerSuccess;

  /// No description provided for @firstNameFieldLabel.
  ///
  /// In en, this message translates to:
  /// **'First name'**
  String get firstNameFieldLabel;

  /// No description provided for @firstNameFieldHint.
  ///
  /// In en, this message translates to:
  /// **'Enter first name'**
  String get firstNameFieldHint;

  /// No description provided for @lastNameFieldLabel.
  ///
  /// In en, this message translates to:
  /// **'Last name'**
  String get lastNameFieldLabel;

  /// No description provided for @lastNameFieldHint.
  ///
  /// In en, this message translates to:
  /// **'Enter last name'**
  String get lastNameFieldHint;

  /// No description provided for @confirmPasswordFieldLabel.
  ///
  /// In en, this message translates to:
  /// **'Confirm password'**
  String get confirmPasswordFieldLabel;

  /// No description provided for @confirmPasswordFieldHint.
  ///
  /// In en, this message translates to:
  /// **'Confirm password'**
  String get confirmPasswordFieldHint;

  /// No description provided for @validationEmptyEmail.
  ///
  /// In en, this message translates to:
  /// **'Enter email'**
  String get validationEmptyEmail;

  /// No description provided for @validationInvalidEmail.
  ///
  /// In en, this message translates to:
  /// **'Invalid email'**
  String get validationInvalidEmail;

  /// No description provided for @validationMinimumPassword.
  ///
  /// In en, this message translates to:
  /// **'Minimum 8 characters'**
  String get validationMinimumPassword;

  /// No description provided for @validationEmptyFirstName.
  ///
  /// In en, this message translates to:
  /// **'Enter first name'**
  String get validationEmptyFirstName;

  /// No description provided for @validationEmptyLastName.
  ///
  /// In en, this message translates to:
  /// **'Enter last name'**
  String get validationEmptyLastName;

  /// No description provided for @validationPasswordsNotMatch.
  ///
  /// In en, this message translates to:
  /// **'Passwords don\'t match'**
  String get validationPasswordsNotMatch;

  /// No description provided for @validationInvalidData.
  ///
  /// In en, this message translates to:
  /// **'Invalid data'**
  String get validationInvalidData;

  /// No description provided for @haveNoAccount.
  ///
  /// In en, this message translates to:
  /// **'Have no account?'**
  String get haveNoAccount;

  /// No description provided for @toRegister.
  ///
  /// In en, this message translates to:
  /// **'To register'**
  String get toRegister;

  /// No description provided for @alreadyHaveAccount.
  ///
  /// In en, this message translates to:
  /// **'Already have an account?'**
  String get alreadyHaveAccount;

  /// No description provided for @incorrectData.
  ///
  /// In en, this message translates to:
  /// **'Incorrect data'**
  String get incorrectData;

  /// No description provided for @navCoworkings.
  ///
  /// In en, this message translates to:
  /// **'Coworkings'**
  String get navCoworkings;

  /// No description provided for @navBookings.
  ///
  /// In en, this message translates to:
  /// **'Bookings'**
  String get navBookings;

  /// No description provided for @navNotifications.
  ///
  /// In en, this message translates to:
  /// **'Notifications'**
  String get navNotifications;

  /// No description provided for @navProfile.
  ///
  /// In en, this message translates to:
  /// **'Profile'**
  String get navProfile;

  /// No description provided for @navAdmin.
  ///
  /// In en, this message translates to:
  /// **'Admin'**
  String get navAdmin;

  /// No description provided for @profileScreenTitle.
  ///
  /// In en, this message translates to:
  /// **'Profile'**
  String get profileScreenTitle;

  /// No description provided for @bookingsScreenTitle.
  ///
  /// In en, this message translates to:
  /// **'Bookings'**
  String get bookingsScreenTitle;

  /// No description provided for @bookingsTabActive.
  ///
  /// In en, this message translates to:
  /// **'Active'**
  String get bookingsTabActive;

  /// No description provided for @bookingsTabHistory.
  ///
  /// In en, this message translates to:
  /// **'History'**
  String get bookingsTabHistory;

  /// No description provided for @bookingLabelCoworking.
  ///
  /// In en, this message translates to:
  /// **'Coworking:'**
  String get bookingLabelCoworking;

  /// No description provided for @bookingLabelPlace.
  ///
  /// In en, this message translates to:
  /// **'Place:'**
  String get bookingLabelPlace;

  /// No description provided for @bookingLabelCancelReason.
  ///
  /// In en, this message translates to:
  /// **'Cancel reason:'**
  String get bookingLabelCancelReason;

  /// No description provided for @bookingButtonCancel.
  ///
  /// In en, this message translates to:
  /// **'Cancel booking'**
  String get bookingButtonCancel;

  /// No description provided for @cancelReasonByUser.
  ///
  /// In en, this message translates to:
  /// **'Cancelled by user'**
  String get cancelReasonByUser;

  /// No description provided for @failedFetchActiveBookings.
  ///
  /// In en, this message translates to:
  /// **'Failed to fetch active bookings'**
  String get failedFetchActiveBookings;

  /// No description provided for @failedFetchHistoryBookings.
  ///
  /// In en, this message translates to:
  /// **'Failed to fetch history bookings'**
  String get failedFetchHistoryBookings;

  /// No description provided for @bookingCancelled.
  ///
  /// In en, this message translates to:
  /// **'Booking cancelled'**
  String get bookingCancelled;

  /// No description provided for @failedCancelBooking.
  ///
  /// In en, this message translates to:
  /// **'Failed to cancel booking'**
  String get failedCancelBooking;

  /// No description provided for @coworkingsScreenTitle.
  ///
  /// In en, this message translates to:
  /// **'Coworkings'**
  String get coworkingsScreenTitle;

  /// No description provided for @coworkingDefaultName.
  ///
  /// In en, this message translates to:
  /// **'Coworking'**
  String get coworkingDefaultName;

  /// No description provided for @coworkingUnavailable.
  ///
  /// In en, this message translates to:
  /// **'Coworking is unavailable right now'**
  String get coworkingUnavailable;

  /// No description provided for @bookingCreatedSuccess.
  ///
  /// In en, this message translates to:
  /// **'Booking created successfully'**
  String get bookingCreatedSuccess;

  /// No description provided for @adminPanelTitle.
  ///
  /// In en, this message translates to:
  /// **'Admin Panel'**
  String get adminPanelTitle;

  /// No description provided for @adminNoCoworkings.
  ///
  /// In en, this message translates to:
  /// **'No coworkings'**
  String get adminNoCoworkings;

  /// No description provided for @adminFailedLoadCoworkings.
  ///
  /// In en, this message translates to:
  /// **'Failed to load coworkings'**
  String get adminFailedLoadCoworkings;

  /// No description provided for @errorUserAlreadyExists.
  ///
  /// In en, this message translates to:
  /// **'User already exists'**
  String get errorUserAlreadyExists;

  /// No description provided for @errorRegistrationFailed.
  ///
  /// In en, this message translates to:
  /// **'Registration failed'**
  String get errorRegistrationFailed;

  /// No description provided for @errorInvalidCredentials.
  ///
  /// In en, this message translates to:
  /// **'Invalid credentials'**
  String get errorInvalidCredentials;

  /// No description provided for @errorLoginFailed.
  ///
  /// In en, this message translates to:
  /// **'Login failed'**
  String get errorLoginFailed;

  /// No description provided for @errorRefreshTokenFailed.
  ///
  /// In en, this message translates to:
  /// **'Refresh token failed'**
  String get errorRefreshTokenFailed;

  /// No description provided for @errorLogoutFailed.
  ///
  /// In en, this message translates to:
  /// **'Logout failed'**
  String get errorLogoutFailed;

  /// No description provided for @errorUnauthorized.
  ///
  /// In en, this message translates to:
  /// **'Unauthorized'**
  String get errorUnauthorized;

  /// No description provided for @errorFetchProfile.
  ///
  /// In en, this message translates to:
  /// **'Failed to fetch profile'**
  String get errorFetchProfile;

  /// No description provided for @errorFetchActiveSessions.
  ///
  /// In en, this message translates to:
  /// **'Failed to fetch active sessions'**
  String get errorFetchActiveSessions;

  /// No description provided for @errorFetchAllSessions.
  ///
  /// In en, this message translates to:
  /// **'Failed to fetch all sessions'**
  String get errorFetchAllSessions;

  /// No description provided for @errorRevokeSession.
  ///
  /// In en, this message translates to:
  /// **'Failed to revoke session'**
  String get errorRevokeSession;

  /// No description provided for @errorFetchCoworkings.
  ///
  /// In en, this message translates to:
  /// **'Failed to fetch coworkings'**
  String get errorFetchCoworkings;

  /// No description provided for @errorFetchCoworkingDetails.
  ///
  /// In en, this message translates to:
  /// **'Failed to fetch coworking details'**
  String get errorFetchCoworkingDetails;

  /// No description provided for @errorFetchAvailablePlaces.
  ///
  /// In en, this message translates to:
  /// **'Failed to fetch available places'**
  String get errorFetchAvailablePlaces;

  /// No description provided for @errorFetchLayout.
  ///
  /// In en, this message translates to:
  /// **'Failed to fetch layout for coworking'**
  String get errorFetchLayout;

  /// No description provided for @errorCreateBooking.
  ///
  /// In en, this message translates to:
  /// **'Failed to create booking'**
  String get errorCreateBooking;

  /// No description provided for @errorCancelBooking.
  ///
  /// In en, this message translates to:
  /// **'Failed to cancel booking'**
  String get errorCancelBooking;

  /// No description provided for @settingsLanguage.
  ///
  /// In en, this message translates to:
  /// **'Language'**
  String get settingsLanguage;

  /// No description provided for @settingsSystemDefault.
  ///
  /// In en, this message translates to:
  /// **'System default'**
  String get settingsSystemDefault;

  /// No description provided for @settingsEnglish.
  ///
  /// In en, this message translates to:
  /// **'English'**
  String get settingsEnglish;

  /// No description provided for @settingsRussian.
  ///
  /// In en, this message translates to:
  /// **'Русский'**
  String get settingsRussian;

  /// No description provided for @currentSession.
  ///
  /// In en, this message translates to:
  /// **'current'**
  String get currentSession;

  /// No description provided for @activeSession.
  ///
  /// In en, this message translates to:
  /// **'active'**
  String get activeSession;

  /// No description provided for @revokedSession.
  ///
  /// In en, this message translates to:
  /// **'revoked'**
  String get revokedSession;

  /// No description provided for @sessionCreated.
  ///
  /// In en, this message translates to:
  /// **'Created: '**
  String get sessionCreated;

  /// No description provided for @sessionLastUsedAt.
  ///
  /// In en, this message translates to:
  /// **'Last used at: '**
  String get sessionLastUsedAt;

  /// No description provided for @sessionExpiresAt.
  ///
  /// In en, this message translates to:
  /// **'Expires at: '**
  String get sessionExpiresAt;

  /// No description provided for @revokeSession.
  ///
  /// In en, this message translates to:
  /// **'Revoke session'**
  String get revokeSession;

  /// No description provided for @sessions.
  ///
  /// In en, this message translates to:
  /// **'Sessions'**
  String get sessions;

  /// No description provided for @notifications.
  ///
  /// In en, this message translates to:
  /// **'Notifications'**
  String get notifications;

  /// No description provided for @markAllAsRead.
  ///
  /// In en, this message translates to:
  /// **'Mark all as read'**
  String get markAllAsRead;

  /// No description provided for @notificationReadAt.
  ///
  /// In en, this message translates to:
  /// **'Read'**
  String get notificationReadAt;

  /// No description provided for @bookingStatusActive.
  ///
  /// In en, this message translates to:
  /// **'Active'**
  String get bookingStatusActive;

  /// No description provided for @bookingStatusCompleted.
  ///
  /// In en, this message translates to:
  /// **'Completed'**
  String get bookingStatusCompleted;

  /// No description provided for @bookingStatusCancelled.
  ///
  /// In en, this message translates to:
  /// **'Cancelled'**
  String get bookingStatusCancelled;

  /// No description provided for @placeTypeOpenDesk.
  ///
  /// In en, this message translates to:
  /// **'Open desk'**
  String get placeTypeOpenDesk;

  /// No description provided for @noActiveBookings.
  ///
  /// In en, this message translates to:
  /// **'No active bookings'**
  String get noActiveBookings;

  /// No description provided for @noHistoryBookings.
  ///
  /// In en, this message translates to:
  /// **'No booking history'**
  String get noHistoryBookings;
}

class _AppLocalizationsDelegate
    extends LocalizationsDelegate<AppLocalizations> {
  const _AppLocalizationsDelegate();

  @override
  Future<AppLocalizations> load(Locale locale) {
    return SynchronousFuture<AppLocalizations>(lookupAppLocalizations(locale));
  }

  @override
  bool isSupported(Locale locale) =>
      <String>['en', 'ru'].contains(locale.languageCode);

  @override
  bool shouldReload(_AppLocalizationsDelegate old) => false;
}

AppLocalizations lookupAppLocalizations(Locale locale) {
  // Lookup logic when only language code is specified.
  switch (locale.languageCode) {
    case 'en':
      return AppLocalizationsEn();
    case 'ru':
      return AppLocalizationsRu();
  }

  throw FlutterError(
    'AppLocalizations.delegate failed to load unsupported locale "$locale". This is likely '
    'an issue with the localizations generation tool. Please file an issue '
    'on GitHub with a reproducible sample app and the gen-l10n configuration '
    'that was used.',
  );
}
