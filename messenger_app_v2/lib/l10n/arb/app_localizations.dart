import 'dart:async';

import 'package:flutter/foundation.dart';
import 'package:flutter/widgets.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:intl/intl.dart' as intl;

import 'app_localizations_en.dart';
import 'app_localizations_kk.dart';
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
/// import 'arb/app_localizations.dart';
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
    Locale('kk'),
    Locale('ru'),
  ];

  /// No description provided for @appTitle.
  ///
  /// In en, this message translates to:
  /// **'Crypto Messenger'**
  String get appTitle;

  /// No description provided for @cancel.
  ///
  /// In en, this message translates to:
  /// **'Cancel'**
  String get cancel;

  /// No description provided for @confirm.
  ///
  /// In en, this message translates to:
  /// **'Confirm'**
  String get confirm;

  /// No description provided for @settingsTitle.
  ///
  /// In en, this message translates to:
  /// **'Settings'**
  String get settingsTitle;

  /// No description provided for @themeTitle.
  ///
  /// In en, this message translates to:
  /// **'Theme'**
  String get themeTitle;

  /// No description provided for @languageTitle.
  ///
  /// In en, this message translates to:
  /// **'Language'**
  String get languageTitle;

  /// No description provided for @languageEnglish.
  ///
  /// In en, this message translates to:
  /// **'English'**
  String get languageEnglish;

  /// No description provided for @languageRussian.
  ///
  /// In en, this message translates to:
  /// **'Russian'**
  String get languageRussian;

  /// No description provided for @languageKazakh.
  ///
  /// In en, this message translates to:
  /// **'Kazakh'**
  String get languageKazakh;

  /// No description provided for @authWelcomeTitle.
  ///
  /// In en, this message translates to:
  /// **'Suda Messenger'**
  String get authWelcomeTitle;

  /// No description provided for @authWelcomeSubtitle.
  ///
  /// In en, this message translates to:
  /// **'The decentralized hub for your conversations and crypto assets. Connect, trade, and explore the new web3 ecosystem.'**
  String get authWelcomeSubtitle;

  /// No description provided for @authGetStarted.
  ///
  /// In en, this message translates to:
  /// **'Get Started'**
  String get authGetStarted;

  /// No description provided for @authLogin.
  ///
  /// In en, this message translates to:
  /// **'Log In'**
  String get authLogin;

  /// No description provided for @authCreateAccount.
  ///
  /// In en, this message translates to:
  /// **'Create Account'**
  String get authCreateAccount;

  /// No description provided for @authWelcomeBack.
  ///
  /// In en, this message translates to:
  /// **'Welcome Back!'**
  String get authWelcomeBack;

  /// No description provided for @authExcited.
  ///
  /// In en, this message translates to:
  /// **'We\'re so excited to see you again!'**
  String get authExcited;

  /// No description provided for @authEmailOrUsername.
  ///
  /// In en, this message translates to:
  /// **'Email or Username'**
  String get authEmailOrUsername;

  /// No description provided for @authEmail.
  ///
  /// In en, this message translates to:
  /// **'Email'**
  String get authEmail;

  /// No description provided for @authPassword.
  ///
  /// In en, this message translates to:
  /// **'Password'**
  String get authPassword;

  /// No description provided for @authUsername.
  ///
  /// In en, this message translates to:
  /// **'Username'**
  String get authUsername;

  /// No description provided for @authDisplayName.
  ///
  /// In en, this message translates to:
  /// **'Display Name'**
  String get authDisplayName;

  /// No description provided for @authUsernameHint.
  ///
  /// In en, this message translates to:
  /// **'3–30 chars · letters and digits'**
  String get authUsernameHint;

  /// No description provided for @authForgotPassword.
  ///
  /// In en, this message translates to:
  /// **'Forgot Password?'**
  String get authForgotPassword;

  /// No description provided for @authNeedAccount.
  ///
  /// In en, this message translates to:
  /// **'Need an account?'**
  String get authNeedAccount;

  /// No description provided for @authRegister.
  ///
  /// In en, this message translates to:
  /// **'Register'**
  String get authRegister;

  /// No description provided for @authAlreadyHaveAccount.
  ///
  /// In en, this message translates to:
  /// **'Already have an account?'**
  String get authAlreadyHaveAccount;

  /// No description provided for @authOrLoginWith.
  ///
  /// In en, this message translates to:
  /// **'Or login with'**
  String get authOrLoginWith;

  /// No description provided for @authJoinFuture.
  ///
  /// In en, this message translates to:
  /// **'Join the decentralized future'**
  String get authJoinFuture;

  /// No description provided for @authVerifyTitle.
  ///
  /// In en, this message translates to:
  /// **'Verify Identity'**
  String get authVerifyTitle;

  /// No description provided for @authVerifySubtitle.
  ///
  /// In en, this message translates to:
  /// **'Enter the 6-digit code sent to your email to access your crypto wallet.'**
  String get authVerifySubtitle;

  /// No description provided for @authVerifyButton.
  ///
  /// In en, this message translates to:
  /// **'Verify & Continue'**
  String get authVerifyButton;

  /// No description provided for @authResendCode.
  ///
  /// In en, this message translates to:
  /// **'Resend Code'**
  String get authResendCode;

  /// No description provided for @authCodeExpires.
  ///
  /// In en, this message translates to:
  /// **'Code expires in'**
  String get authCodeExpires;

  /// No description provided for @errorGeneric.
  ///
  /// In en, this message translates to:
  /// **'Something went wrong'**
  String get errorGeneric;

  /// No description provided for @errorNetwork.
  ///
  /// In en, this message translates to:
  /// **'No internet connection'**
  String get errorNetwork;

  /// No description provided for @placeholderErrorTitle.
  ///
  /// In en, this message translates to:
  /// **'Error'**
  String get placeholderErrorTitle;

  /// No description provided for @placeholderErrorMessage.
  ///
  /// In en, this message translates to:
  /// **'Something went wrong.\nWe are already fixing it.'**
  String get placeholderErrorMessage;

  /// No description provided for @placeholderInProgressTitle.
  ///
  /// In en, this message translates to:
  /// **'In Development'**
  String get placeholderInProgressTitle;

  /// No description provided for @placeholderInProgressMessage.
  ///
  /// In en, this message translates to:
  /// **'This feature is not ready yet.\nCheck back later!'**
  String get placeholderInProgressMessage;

  /// No description provided for @placeholderTestTitle.
  ///
  /// In en, this message translates to:
  /// **'Test'**
  String get placeholderTestTitle;

  /// No description provided for @placeholderTestMessage.
  ///
  /// In en, this message translates to:
  /// **'This is a test screen to check navigation.'**
  String get placeholderTestMessage;

  /// No description provided for @placeholderNoContentTitle.
  ///
  /// In en, this message translates to:
  /// **'Empty'**
  String get placeholderNoContentTitle;

  /// No description provided for @placeholderNoContentMessage.
  ///
  /// In en, this message translates to:
  /// **'There is nothing here yet.'**
  String get placeholderNoContentMessage;

  /// No description provided for @buttonRetry.
  ///
  /// In en, this message translates to:
  /// **'Retry'**
  String get buttonRetry;

  /// No description provided for @authForgotPasswordTitle.
  ///
  /// In en, this message translates to:
  /// **'Forgot Password'**
  String get authForgotPasswordTitle;

  /// No description provided for @authForgotPasswordSubtitle.
  ///
  /// In en, this message translates to:
  /// **'Enter your email and we will send a verification code.'**
  String get authForgotPasswordSubtitle;

  /// No description provided for @authSendVerificationCode.
  ///
  /// In en, this message translates to:
  /// **'Send Verification Code'**
  String get authSendVerificationCode;

  /// No description provided for @authBackToLogin.
  ///
  /// In en, this message translates to:
  /// **'Back to Login'**
  String get authBackToLogin;

  /// No description provided for @authResetPasswordTitle.
  ///
  /// In en, this message translates to:
  /// **'Reset Password'**
  String get authResetPasswordTitle;

  /// No description provided for @authEnterCodeSentTo.
  ///
  /// In en, this message translates to:
  /// **'Enter the 6-digit code sent to\n{email}'**
  String authEnterCodeSentTo(String email);

  /// No description provided for @authNewPassword.
  ///
  /// In en, this message translates to:
  /// **'New Password'**
  String get authNewPassword;

  /// No description provided for @authConfirmPassword.
  ///
  /// In en, this message translates to:
  /// **'Confirm Password'**
  String get authConfirmPassword;

  /// No description provided for @authPasswordChangedSuccess.
  ///
  /// In en, this message translates to:
  /// **'Password changed successfully'**
  String get authPasswordChangedSuccess;

  /// No description provided for @authVerificationCodeSent.
  ///
  /// In en, this message translates to:
  /// **'If email exists, verification code for password reset sent'**
  String get authVerificationCodeSent;

  /// No description provided for @authEmailRequired.
  ///
  /// In en, this message translates to:
  /// **'Email is required'**
  String get authEmailRequired;

  /// No description provided for @authEnterFullCode.
  ///
  /// In en, this message translates to:
  /// **'Enter full verification code'**
  String get authEnterFullCode;

  /// No description provided for @authPasswordMinLength.
  ///
  /// In en, this message translates to:
  /// **'Password must be at least 6 characters'**
  String get authPasswordMinLength;

  /// No description provided for @authPasswordsDoNotMatch.
  ///
  /// In en, this message translates to:
  /// **'Passwords do not match'**
  String get authPasswordsDoNotMatch;

  /// No description provided for @authEmailMissing.
  ///
  /// In en, this message translates to:
  /// **'Email is missing. Request reset again.'**
  String get authEmailMissing;

  /// No description provided for @chatTitle.
  ///
  /// In en, this message translates to:
  /// **'Chats'**
  String get chatTitle;

  /// No description provided for @chatFindOrStart.
  ///
  /// In en, this message translates to:
  /// **'Find or start a conversation'**
  String get chatFindOrStart;

  /// No description provided for @chatFilterAll.
  ///
  /// In en, this message translates to:
  /// **'All'**
  String get chatFilterAll;

  /// No description provided for @chatFilterUnread.
  ///
  /// In en, this message translates to:
  /// **'Unread'**
  String get chatFilterUnread;

  /// No description provided for @chatFilterPersonal.
  ///
  /// In en, this message translates to:
  /// **'Personal'**
  String get chatFilterPersonal;

  /// No description provided for @chatFilterGroups.
  ///
  /// In en, this message translates to:
  /// **'Groups'**
  String get chatFilterGroups;

  /// No description provided for @chatFilterChannels.
  ///
  /// In en, this message translates to:
  /// **'Channels'**
  String get chatFilterChannels;

  /// No description provided for @chatNewGroup.
  ///
  /// In en, this message translates to:
  /// **'New group'**
  String get chatNewGroup;

  /// No description provided for @chatNewChannel.
  ///
  /// In en, this message translates to:
  /// **'New channel'**
  String get chatNewChannel;

  /// No description provided for @chatFindContact.
  ///
  /// In en, this message translates to:
  /// **'Find a contact'**
  String get chatFindContact;

  /// No description provided for @chatSavedMessages.
  ///
  /// In en, this message translates to:
  /// **'Saved Messages'**
  String get chatSavedMessages;

  /// No description provided for @chatUnknownUser.
  ///
  /// In en, this message translates to:
  /// **'Unknown User'**
  String get chatUnknownUser;

  /// No description provided for @chatNoMessagesYet.
  ///
  /// In en, this message translates to:
  /// **'No messages yet'**
  String get chatNoMessagesYet;

  /// No description provided for @chatStatusTyping.
  ///
  /// In en, this message translates to:
  /// **'Typing...'**
  String get chatStatusTyping;

  /// No description provided for @chatStatusOnline.
  ///
  /// In en, this message translates to:
  /// **'Online'**
  String get chatStatusOnline;

  /// No description provided for @chatStatusOffline.
  ///
  /// In en, this message translates to:
  /// **'Offline'**
  String get chatStatusOffline;

  /// No description provided for @searchTitle.
  ///
  /// In en, this message translates to:
  /// **'Search'**
  String get searchTitle;

  /// No description provided for @searchHint.
  ///
  /// In en, this message translates to:
  /// **'Search users, chats, messages…'**
  String get searchHint;

  /// No description provided for @searchNoUsersFound.
  ///
  /// In en, this message translates to:
  /// **'No users found'**
  String get searchNoUsersFound;

  /// No description provided for @searchEnterUsername.
  ///
  /// In en, this message translates to:
  /// **'Enter username to search'**
  String get searchEnterUsername;

  /// No description provided for @searchRecent.
  ///
  /// In en, this message translates to:
  /// **'Recent'**
  String get searchRecent;

  /// No description provided for @searchUsers.
  ///
  /// In en, this message translates to:
  /// **'Users'**
  String get searchUsers;

  /// No description provided for @searchChats.
  ///
  /// In en, this message translates to:
  /// **'Chats'**
  String get searchChats;

  /// No description provided for @searchMessages.
  ///
  /// In en, this message translates to:
  /// **'Messages'**
  String get searchMessages;

  /// No description provided for @searchClearQuery.
  ///
  /// In en, this message translates to:
  /// **'Clear'**
  String get searchClearQuery;

  /// No description provided for @searchNoResults.
  ///
  /// In en, this message translates to:
  /// **'No results for \"{query}\"'**
  String searchNoResults(String query);

  /// No description provided for @chatSearchTitle.
  ///
  /// In en, this message translates to:
  /// **'Search in chat'**
  String get chatSearchTitle;

  /// No description provided for @chatSearchHint.
  ///
  /// In en, this message translates to:
  /// **'Search messages…'**
  String get chatSearchHint;

  /// No description provided for @profileUpdated.
  ///
  /// In en, this message translates to:
  /// **'Profile updated'**
  String get profileUpdated;

  /// No description provided for @avatarUpdated.
  ///
  /// In en, this message translates to:
  /// **'Avatar updated'**
  String get avatarUpdated;

  /// No description provided for @loggedOut.
  ///
  /// In en, this message translates to:
  /// **'Logged out'**
  String get loggedOut;

  /// No description provided for @navChat.
  ///
  /// In en, this message translates to:
  /// **'Chat'**
  String get navChat;

  /// No description provided for @navHub.
  ///
  /// In en, this message translates to:
  /// **'Hub'**
  String get navHub;

  /// No description provided for @navFriends.
  ///
  /// In en, this message translates to:
  /// **'Friends'**
  String get navFriends;

  /// No description provided for @navProfile.
  ///
  /// In en, this message translates to:
  /// **'Profile'**
  String get navProfile;

  /// No description provided for @themeSelectTitle.
  ///
  /// In en, this message translates to:
  /// **'Select Theme'**
  String get themeSelectTitle;

  /// No description provided for @chatLastSeenPrefix.
  ///
  /// In en, this message translates to:
  /// **'last seen'**
  String get chatLastSeenPrefix;

  /// No description provided for @chatLastSeenJustNow.
  ///
  /// In en, this message translates to:
  /// **'just now'**
  String get chatLastSeenJustNow;

  /// No description provided for @chatLastSeenMinutes.
  ///
  /// In en, this message translates to:
  /// **'min ago'**
  String get chatLastSeenMinutes;

  /// No description provided for @chatLastSeenHours.
  ///
  /// In en, this message translates to:
  /// **'h ago'**
  String get chatLastSeenHours;

  /// No description provided for @messageReadersTitle.
  ///
  /// In en, this message translates to:
  /// **'Read by'**
  String get messageReadersTitle;

  /// No description provided for @messageReadersEmpty.
  ///
  /// In en, this message translates to:
  /// **'No one has read this message yet'**
  String get messageReadersEmpty;

  /// No description provided for @messageActionEdit.
  ///
  /// In en, this message translates to:
  /// **'Edit'**
  String get messageActionEdit;

  /// No description provided for @messageActionCopy.
  ///
  /// In en, this message translates to:
  /// **'Copy'**
  String get messageActionCopy;

  /// No description provided for @messageActionDeleteForMe.
  ///
  /// In en, this message translates to:
  /// **'Delete for me'**
  String get messageActionDeleteForMe;

  /// No description provided for @messageActionDeleteForEveryone.
  ///
  /// In en, this message translates to:
  /// **'Delete for everyone'**
  String get messageActionDeleteForEveryone;

  /// No description provided for @messageEditingLabel.
  ///
  /// In en, this message translates to:
  /// **'Editing'**
  String get messageEditingLabel;

  /// No description provided for @messageEditCancel.
  ///
  /// In en, this message translates to:
  /// **'Cancel editing'**
  String get messageEditCancel;

  /// No description provided for @messageReply.
  ///
  /// In en, this message translates to:
  /// **'Reply'**
  String get messageReply;

  /// No description provided for @messageForward.
  ///
  /// In en, this message translates to:
  /// **'Forward'**
  String get messageForward;

  /// No description provided for @forwardTo.
  ///
  /// In en, this message translates to:
  /// **'Forward to'**
  String get forwardTo;

  /// No description provided for @chatInfo.
  ///
  /// In en, this message translates to:
  /// **'Chat Info'**
  String get chatInfo;

  /// No description provided for @chatMembers.
  ///
  /// In en, this message translates to:
  /// **'Members'**
  String get chatMembers;

  /// No description provided for @chatAddMember.
  ///
  /// In en, this message translates to:
  /// **'Add member'**
  String get chatAddMember;

  /// No description provided for @chatLeave.
  ///
  /// In en, this message translates to:
  /// **'Leave chat'**
  String get chatLeave;

  /// No description provided for @chatDelete.
  ///
  /// In en, this message translates to:
  /// **'Delete chat'**
  String get chatDelete;

  /// No description provided for @chatLeaveConfirm.
  ///
  /// In en, this message translates to:
  /// **'Are you sure you want to leave?'**
  String get chatLeaveConfirm;

  /// No description provided for @chatDeleteConfirm.
  ///
  /// In en, this message translates to:
  /// **'Are you sure you want to delete this chat?'**
  String get chatDeleteConfirm;

  /// No description provided for @chatMute.
  ///
  /// In en, this message translates to:
  /// **'Mute'**
  String get chatMute;

  /// No description provided for @chatUnmute.
  ///
  /// In en, this message translates to:
  /// **'Unmute'**
  String get chatUnmute;

  /// No description provided for @chatDeleteForMe.
  ///
  /// In en, this message translates to:
  /// **'Delete for me'**
  String get chatDeleteForMe;

  /// No description provided for @chatDeleteForEveryone.
  ///
  /// In en, this message translates to:
  /// **'Delete for everyone'**
  String get chatDeleteForEveryone;

  /// No description provided for @chatDeleteGroup.
  ///
  /// In en, this message translates to:
  /// **'Delete group'**
  String get chatDeleteGroup;

  /// No description provided for @chatDeleteChannel.
  ///
  /// In en, this message translates to:
  /// **'Delete channel'**
  String get chatDeleteChannel;

  /// No description provided for @chatDeleteConfirmForAll.
  ///
  /// In en, this message translates to:
  /// **'This will delete the chat for all members. Continue?'**
  String get chatDeleteConfirmForAll;

  /// No description provided for @chatDeleteConfirmForMe.
  ///
  /// In en, this message translates to:
  /// **'This will hide the chat and clear your history. Continue?'**
  String get chatDeleteConfirmForMe;

  /// No description provided for @chatEditGroup.
  ///
  /// In en, this message translates to:
  /// **'Edit group'**
  String get chatEditGroup;

  /// No description provided for @chatGroupDescription.
  ///
  /// In en, this message translates to:
  /// **'Description'**
  String get chatGroupDescription;

  /// No description provided for @chatGroupDescriptionHint.
  ///
  /// In en, this message translates to:
  /// **'About this group…'**
  String get chatGroupDescriptionHint;

  /// No description provided for @chatSaveChanges.
  ///
  /// In en, this message translates to:
  /// **'Save changes'**
  String get chatSaveChanges;

  /// No description provided for @chatAddMembers.
  ///
  /// In en, this message translates to:
  /// **'Add members'**
  String get chatAddMembers;

  /// No description provided for @chatMakeAdmin.
  ///
  /// In en, this message translates to:
  /// **'Make admin'**
  String get chatMakeAdmin;

  /// No description provided for @chatRemoveMember.
  ///
  /// In en, this message translates to:
  /// **'Remove'**
  String get chatRemoveMember;

  /// No description provided for @newChat.
  ///
  /// In en, this message translates to:
  /// **'New conversation'**
  String get newChat;

  /// No description provided for @newGroup.
  ///
  /// In en, this message translates to:
  /// **'New group'**
  String get newGroup;

  /// No description provided for @newChannel.
  ///
  /// In en, this message translates to:
  /// **'New channel'**
  String get newChannel;

  /// No description provided for @groupName.
  ///
  /// In en, this message translates to:
  /// **'Group name'**
  String get groupName;

  /// No description provided for @groupNameHint.
  ///
  /// In en, this message translates to:
  /// **'Enter group name'**
  String get groupNameHint;

  /// No description provided for @createGroup.
  ///
  /// In en, this message translates to:
  /// **'Create group'**
  String get createGroup;

  /// No description provided for @selectMembers.
  ///
  /// In en, this message translates to:
  /// **'Select members'**
  String get selectMembers;

  /// No description provided for @roleOwner.
  ///
  /// In en, this message translates to:
  /// **'Owner'**
  String get roleOwner;

  /// No description provided for @roleAdmin.
  ///
  /// In en, this message translates to:
  /// **'Admin'**
  String get roleAdmin;

  /// No description provided for @roleMember.
  ///
  /// In en, this message translates to:
  /// **'Member'**
  String get roleMember;

  /// No description provided for @profileTitle.
  ///
  /// In en, this message translates to:
  /// **'Profile'**
  String get profileTitle;

  /// No description provided for @profileSudaBalance.
  ///
  /// In en, this message translates to:
  /// **'SUDA Balance'**
  String get profileSudaBalance;

  /// No description provided for @profileOpenWallet.
  ///
  /// In en, this message translates to:
  /// **'Open wallet'**
  String get profileOpenWallet;

  /// No description provided for @profileChats.
  ///
  /// In en, this message translates to:
  /// **'Chats'**
  String get profileChats;

  /// No description provided for @profileChannels.
  ///
  /// In en, this message translates to:
  /// **'Channels'**
  String get profileChannels;

  /// No description provided for @profileContacts.
  ///
  /// In en, this message translates to:
  /// **'Contacts'**
  String get profileContacts;

  /// No description provided for @profileEditProfile.
  ///
  /// In en, this message translates to:
  /// **'Edit profile'**
  String get profileEditProfile;

  /// No description provided for @profileNotifications.
  ///
  /// In en, this message translates to:
  /// **'Notifications'**
  String get profileNotifications;

  /// No description provided for @profilePrivacy.
  ///
  /// In en, this message translates to:
  /// **'Privacy & security'**
  String get profilePrivacy;

  /// No description provided for @profileLanguage.
  ///
  /// In en, this message translates to:
  /// **'Language'**
  String get profileLanguage;

  /// No description provided for @profileBlocked.
  ///
  /// In en, this message translates to:
  /// **'Blocked users'**
  String get profileBlocked;

  /// No description provided for @profileAbout.
  ///
  /// In en, this message translates to:
  /// **'About Suda'**
  String get profileAbout;

  /// No description provided for @profileSignOut.
  ///
  /// In en, this message translates to:
  /// **'Sign out'**
  String get profileSignOut;

  /// No description provided for @settingsAccount.
  ///
  /// In en, this message translates to:
  /// **'Account'**
  String get settingsAccount;

  /// No description provided for @settingsPrivacy.
  ///
  /// In en, this message translates to:
  /// **'Privacy'**
  String get settingsPrivacy;

  /// No description provided for @settingsNotifData.
  ///
  /// In en, this message translates to:
  /// **'Notifications & data'**
  String get settingsNotifData;

  /// No description provided for @settingsAppearance.
  ///
  /// In en, this message translates to:
  /// **'Appearance'**
  String get settingsAppearance;

  /// No description provided for @settingsSuda.
  ///
  /// In en, this message translates to:
  /// **'Suda'**
  String get settingsSuda;

  /// No description provided for @settingsEmail.
  ///
  /// In en, this message translates to:
  /// **'Email'**
  String get settingsEmail;

  /// No description provided for @settingsChangePassword.
  ///
  /// In en, this message translates to:
  /// **'Change password'**
  String get settingsChangePassword;

  /// No description provided for @settingsActiveSessions.
  ///
  /// In en, this message translates to:
  /// **'Active sessions'**
  String get settingsActiveSessions;

  /// No description provided for @settingsShowOnline.
  ///
  /// In en, this message translates to:
  /// **'Show online status'**
  String get settingsShowOnline;

  /// No description provided for @settingsReadReceipts.
  ///
  /// In en, this message translates to:
  /// **'Read receipts'**
  String get settingsReadReceipts;

  /// No description provided for @settingsLastSeen.
  ///
  /// In en, this message translates to:
  /// **'Last seen visibility'**
  String get settingsLastSeen;

  /// No description provided for @settingsLastSeenEveryone.
  ///
  /// In en, this message translates to:
  /// **'Everyone'**
  String get settingsLastSeenEveryone;

  /// No description provided for @settingsLastSeenContacts.
  ///
  /// In en, this message translates to:
  /// **'Contacts'**
  String get settingsLastSeenContacts;

  /// No description provided for @settingsLastSeenNobody.
  ///
  /// In en, this message translates to:
  /// **'Nobody'**
  String get settingsLastSeenNobody;

  /// No description provided for @settingsBlocked.
  ///
  /// In en, this message translates to:
  /// **'Blocked users'**
  String get settingsBlocked;

  /// No description provided for @settingsPushNotifications.
  ///
  /// In en, this message translates to:
  /// **'Push notifications'**
  String get settingsPushNotifications;

  /// No description provided for @settingsAutoDownload.
  ///
  /// In en, this message translates to:
  /// **'Auto-download media'**
  String get settingsAutoDownload;

  /// No description provided for @settingsAutoAlways.
  ///
  /// In en, this message translates to:
  /// **'Always'**
  String get settingsAutoAlways;

  /// No description provided for @settingsAutoWifiOnly.
  ///
  /// In en, this message translates to:
  /// **'Wi-Fi only'**
  String get settingsAutoWifiOnly;

  /// No description provided for @settingsAutoNever.
  ///
  /// In en, this message translates to:
  /// **'Never'**
  String get settingsAutoNever;

  /// No description provided for @settingsTheme.
  ///
  /// In en, this message translates to:
  /// **'Theme'**
  String get settingsTheme;

  /// No description provided for @settingsWallet.
  ///
  /// In en, this message translates to:
  /// **'Wallet'**
  String get settingsWallet;

  /// No description provided for @settingsMarketplace.
  ///
  /// In en, this message translates to:
  /// **'Marketplace'**
  String get settingsMarketplace;

  /// No description provided for @settingsAboutSuda.
  ///
  /// In en, this message translates to:
  /// **'About Suda'**
  String get settingsAboutSuda;

  /// No description provided for @settingsPrivacyPolicy.
  ///
  /// In en, this message translates to:
  /// **'Privacy policy'**
  String get settingsPrivacyPolicy;

  /// No description provided for @settingsSignOut.
  ///
  /// In en, this message translates to:
  /// **'Sign out'**
  String get settingsSignOut;

  /// No description provided for @settingsVersion.
  ///
  /// In en, this message translates to:
  /// **'Suda v{version}'**
  String settingsVersion(String version);

  /// No description provided for @settingsAboutDesc.
  ///
  /// In en, this message translates to:
  /// **'The decentralized hub for your conversations and crypto assets.'**
  String get settingsAboutDesc;

  /// No description provided for @themePickerSubtitle.
  ///
  /// In en, this message translates to:
  /// **'Live preview applies instantly'**
  String get themePickerSubtitle;

  /// No description provided for @editProfileTitle.
  ///
  /// In en, this message translates to:
  /// **'Edit profile'**
  String get editProfileTitle;

  /// No description provided for @editProfileSave.
  ///
  /// In en, this message translates to:
  /// **'Save'**
  String get editProfileSave;

  /// No description provided for @editProfileDisplayName.
  ///
  /// In en, this message translates to:
  /// **'Display name'**
  String get editProfileDisplayName;

  /// No description provided for @editProfileFirstName.
  ///
  /// In en, this message translates to:
  /// **'First name'**
  String get editProfileFirstName;

  /// No description provided for @editProfileLastName.
  ///
  /// In en, this message translates to:
  /// **'Last name'**
  String get editProfileLastName;

  /// No description provided for @editProfileBio.
  ///
  /// In en, this message translates to:
  /// **'Bio'**
  String get editProfileBio;

  /// No description provided for @editProfileBioHint.
  ///
  /// In en, this message translates to:
  /// **'Tell people about yourself'**
  String get editProfileBioHint;

  /// No description provided for @editProfileOptional.
  ///
  /// In en, this message translates to:
  /// **'Optional'**
  String get editProfileOptional;

  /// No description provided for @editProfileUsername.
  ///
  /// In en, this message translates to:
  /// **'Username'**
  String get editProfileUsername;

  /// No description provided for @editProfileUsernameHint.
  ///
  /// In en, this message translates to:
  /// **'Username changes are limited to once per week.'**
  String get editProfileUsernameHint;

  /// No description provided for @walletOpenWallet.
  ///
  /// In en, this message translates to:
  /// **'Open wallet'**
  String get walletOpenWallet;

  /// No description provided for @walletMarketplace.
  ///
  /// In en, this message translates to:
  /// **'Marketplace'**
  String get walletMarketplace;

  /// No description provided for @walletSuda.
  ///
  /// In en, this message translates to:
  /// **'SUDA'**
  String get walletSuda;

  /// No description provided for @messageEdited.
  ///
  /// In en, this message translates to:
  /// **'edited'**
  String get messageEdited;

  /// No description provided for @messageForwardedLabel.
  ///
  /// In en, this message translates to:
  /// **'Forwarded'**
  String get messageForwardedLabel;

  /// No description provided for @userProfileOnline.
  ///
  /// In en, this message translates to:
  /// **'online'**
  String get userProfileOnline;

  /// No description provided for @userProfileLastSeen.
  ///
  /// In en, this message translates to:
  /// **'last seen {time}'**
  String userProfileLastSeen(String time);

  /// No description provided for @userProfileMessage.
  ///
  /// In en, this message translates to:
  /// **'Message'**
  String get userProfileMessage;

  /// No description provided for @userProfileSendSuda.
  ///
  /// In en, this message translates to:
  /// **'Send SUDA'**
  String get userProfileSendSuda;

  /// No description provided for @userProfileDonate.
  ///
  /// In en, this message translates to:
  /// **'Donate'**
  String get userProfileDonate;

  /// No description provided for @userProfileMute.
  ///
  /// In en, this message translates to:
  /// **'Mute'**
  String get userProfileMute;

  /// No description provided for @userProfileUnmute.
  ///
  /// In en, this message translates to:
  /// **'Unmute'**
  String get userProfileUnmute;

  /// No description provided for @userProfileBioSection.
  ///
  /// In en, this message translates to:
  /// **'BIO'**
  String get userProfileBioSection;

  /// No description provided for @userProfileWalletAddress.
  ///
  /// In en, this message translates to:
  /// **'Wallet address'**
  String get userProfileWalletAddress;

  /// No description provided for @userProfileSharedMedia.
  ///
  /// In en, this message translates to:
  /// **'Shared media'**
  String get userProfileSharedMedia;

  /// No description provided for @userProfileNotifications.
  ///
  /// In en, this message translates to:
  /// **'Notifications'**
  String get userProfileNotifications;

  /// No description provided for @userProfileBlockUser.
  ///
  /// In en, this message translates to:
  /// **'Block user'**
  String get userProfileBlockUser;

  /// No description provided for @userProfileUnblockUser.
  ///
  /// In en, this message translates to:
  /// **'Unblock'**
  String get userProfileUnblockUser;

  /// No description provided for @userProfileBlockConfirm.
  ///
  /// In en, this message translates to:
  /// **'Block {name}?'**
  String userProfileBlockConfirm(String name);

  /// No description provided for @userProfileBlockConfirmMsg.
  ///
  /// In en, this message translates to:
  /// **'They won\'t be able to message you.'**
  String get userProfileBlockConfirmMsg;

  /// No description provided for @userProfileBlockSuccess.
  ///
  /// In en, this message translates to:
  /// **'User blocked'**
  String get userProfileBlockSuccess;

  /// No description provided for @userProfileUnblockSuccess.
  ///
  /// In en, this message translates to:
  /// **'User unblocked'**
  String get userProfileUnblockSuccess;

  /// No description provided for @contactsTitle.
  ///
  /// In en, this message translates to:
  /// **'Contacts'**
  String get contactsTitle;

  /// No description provided for @contactsEmpty.
  ///
  /// In en, this message translates to:
  /// **'No contacts yet'**
  String get contactsEmpty;

  /// No description provided for @contactsAddContact.
  ///
  /// In en, this message translates to:
  /// **'Add contact'**
  String get contactsAddContact;

  /// No description provided for @contactsRemove.
  ///
  /// In en, this message translates to:
  /// **'Remove'**
  String get contactsRemove;

  /// No description provided for @blockedUsersTitle.
  ///
  /// In en, this message translates to:
  /// **'Blocked users'**
  String get blockedUsersTitle;

  /// No description provided for @blockedUsersEmpty.
  ///
  /// In en, this message translates to:
  /// **'No blocked users'**
  String get blockedUsersEmpty;

  /// No description provided for @blockedUsersUnblock.
  ///
  /// In en, this message translates to:
  /// **'Unblock'**
  String get blockedUsersUnblock;

  /// No description provided for @chatPin.
  ///
  /// In en, this message translates to:
  /// **'Pin'**
  String get chatPin;

  /// No description provided for @chatUnpin.
  ///
  /// In en, this message translates to:
  /// **'Unpin'**
  String get chatUnpin;

  /// No description provided for @comingSoon.
  ///
  /// In en, this message translates to:
  /// **'Coming soon'**
  String get comingSoon;

  /// No description provided for @attachPhoto.
  ///
  /// In en, this message translates to:
  /// **'Photo / Video'**
  String get attachPhoto;

  /// No description provided for @attachFile.
  ///
  /// In en, this message translates to:
  /// **'File'**
  String get attachFile;

  /// No description provided for @attachVoice.
  ///
  /// In en, this message translates to:
  /// **'Voice'**
  String get attachVoice;

  /// No description provided for @voiceRecording.
  ///
  /// In en, this message translates to:
  /// **'Recording…'**
  String get voiceRecording;

  /// No description provided for @voiceSlideCancelHint.
  ///
  /// In en, this message translates to:
  /// **'Tap ✕ to cancel'**
  String get voiceSlideCancelHint;

  /// No description provided for @messageInputHint.
  ///
  /// In en, this message translates to:
  /// **'Message…'**
  String get messageInputHint;

  /// No description provided for @messageTypeImage.
  ///
  /// In en, this message translates to:
  /// **'Photo'**
  String get messageTypeImage;

  /// No description provided for @messageTypeFile.
  ///
  /// In en, this message translates to:
  /// **'File'**
  String get messageTypeFile;

  /// No description provided for @messageTypeVoice.
  ///
  /// In en, this message translates to:
  /// **'Voice message'**
  String get messageTypeVoice;

  /// No description provided for @messageTypeVideo.
  ///
  /// In en, this message translates to:
  /// **'Video'**
  String get messageTypeVideo;

  /// No description provided for @openFile.
  ///
  /// In en, this message translates to:
  /// **'Open'**
  String get openFile;

  /// No description provided for @fileOpenError.
  ///
  /// In en, this message translates to:
  /// **'Couldn\'t open the file'**
  String get fileOpenError;

  /// No description provided for @micPermissionDenied.
  ///
  /// In en, this message translates to:
  /// **'Microphone access is required to record voice messages'**
  String get micPermissionDenied;

  /// No description provided for @voiceRecordError.
  ///
  /// In en, this message translates to:
  /// **'Couldn\'t record the voice message'**
  String get voiceRecordError;

  /// No description provided for @saveToGallery.
  ///
  /// In en, this message translates to:
  /// **'Save to gallery'**
  String get saveToGallery;

  /// No description provided for @saveToGallerySuccess.
  ///
  /// In en, this message translates to:
  /// **'Saved to gallery'**
  String get saveToGallerySuccess;

  /// No description provided for @saveToGalleryError.
  ///
  /// In en, this message translates to:
  /// **'Couldn\'t save to gallery'**
  String get saveToGalleryError;

  /// No description provided for @commentReplyDeleted.
  ///
  /// In en, this message translates to:
  /// **'Deleted comment'**
  String get commentReplyDeleted;

  /// No description provided for @sharedMediaTitle.
  ///
  /// In en, this message translates to:
  /// **'Shared media'**
  String get sharedMediaTitle;

  /// No description provided for @sharedMediaTabMedia.
  ///
  /// In en, this message translates to:
  /// **'Media'**
  String get sharedMediaTabMedia;

  /// No description provided for @sharedMediaPhotos.
  ///
  /// In en, this message translates to:
  /// **'Photos'**
  String get sharedMediaPhotos;

  /// No description provided for @sharedMediaVideos.
  ///
  /// In en, this message translates to:
  /// **'Videos'**
  String get sharedMediaVideos;

  /// No description provided for @sharedMediaFiles.
  ///
  /// In en, this message translates to:
  /// **'Files'**
  String get sharedMediaFiles;

  /// No description provided for @sharedMediaAudio.
  ///
  /// In en, this message translates to:
  /// **'Audio'**
  String get sharedMediaAudio;

  /// No description provided for @sharedMediaEmpty.
  ///
  /// In en, this message translates to:
  /// **'Nothing here yet'**
  String get sharedMediaEmpty;

  /// No description provided for @uploadFailed.
  ///
  /// In en, this message translates to:
  /// **'Upload failed. Try again.'**
  String get uploadFailed;

  /// No description provided for @uploadingMedia.
  ///
  /// In en, this message translates to:
  /// **'Uploading…'**
  String get uploadingMedia;

  /// No description provided for @channelNewTitle.
  ///
  /// In en, this message translates to:
  /// **'New channel'**
  String get channelNewTitle;

  /// No description provided for @channelName.
  ///
  /// In en, this message translates to:
  /// **'Channel name'**
  String get channelName;

  /// No description provided for @channelHandle.
  ///
  /// In en, this message translates to:
  /// **'Handle'**
  String get channelHandle;

  /// No description provided for @channelHandleHint.
  ///
  /// In en, this message translates to:
  /// **'@username'**
  String get channelHandleHint;

  /// No description provided for @channelDescription.
  ///
  /// In en, this message translates to:
  /// **'Description'**
  String get channelDescription;

  /// No description provided for @channelDescriptionHint.
  ///
  /// In en, this message translates to:
  /// **'Tell people about this channel'**
  String get channelDescriptionHint;

  /// No description provided for @channelVisibilityLabel.
  ///
  /// In en, this message translates to:
  /// **'Visibility'**
  String get channelVisibilityLabel;

  /// No description provided for @channelPublic.
  ///
  /// In en, this message translates to:
  /// **'Public'**
  String get channelPublic;

  /// No description provided for @channelPrivate.
  ///
  /// In en, this message translates to:
  /// **'Private'**
  String get channelPrivate;

  /// No description provided for @channelCreate.
  ///
  /// In en, this message translates to:
  /// **'Create channel'**
  String get channelCreate;

  /// No description provided for @channelSubscribersCount.
  ///
  /// In en, this message translates to:
  /// **'{count} subscribers'**
  String channelSubscribersCount(int count);

  /// No description provided for @channelSubscribe.
  ///
  /// In en, this message translates to:
  /// **'Subscribe'**
  String get channelSubscribe;

  /// No description provided for @channelUnsubscribe.
  ///
  /// In en, this message translates to:
  /// **'Unsubscribe'**
  String get channelUnsubscribe;

  /// No description provided for @channelTokenGated.
  ///
  /// In en, this message translates to:
  /// **'Paid subscription'**
  String get channelTokenGated;

  /// No description provided for @channelTokenGatedDesc.
  ///
  /// In en, this message translates to:
  /// **'Subscribe for {amount} SUDA to access this channel.'**
  String channelTokenGatedDesc(String amount);

  /// No description provided for @channelSubscribeForPrice.
  ///
  /// In en, this message translates to:
  /// **'Subscribe for {price} SUDA'**
  String channelSubscribeForPrice(String price);

  /// No description provided for @channelBalanceLabel.
  ///
  /// In en, this message translates to:
  /// **'Your balance'**
  String get channelBalanceLabel;

  /// No description provided for @channelUnlock.
  ///
  /// In en, this message translates to:
  /// **'Unlock & subscribe'**
  String get channelUnlock;

  /// No description provided for @channelTopUp.
  ///
  /// In en, this message translates to:
  /// **'Top up & subscribe'**
  String get channelTopUp;

  /// No description provided for @channelConfirmSubscribe.
  ///
  /// In en, this message translates to:
  /// **'Confirm & Pay'**
  String get channelConfirmSubscribe;

  /// No description provided for @channelSubscribeSuccess.
  ///
  /// In en, this message translates to:
  /// **'Subscribed!'**
  String get channelSubscribeSuccess;

  /// No description provided for @channelInsufficientBalance.
  ///
  /// In en, this message translates to:
  /// **'Not enough SUDA. Top up your wallet and try again.'**
  String get channelInsufficientBalance;

  /// No description provided for @channelNoWallet.
  ///
  /// In en, this message translates to:
  /// **'No wallet connected. Create one in the Wallet section.'**
  String get channelNoWallet;

  /// No description provided for @channelSubscribeFailed.
  ///
  /// In en, this message translates to:
  /// **'Payment failed. Please try again.'**
  String get channelSubscribeFailed;

  /// No description provided for @channelUnsubscribeSuccess.
  ///
  /// In en, this message translates to:
  /// **'Unsubscribed'**
  String get channelUnsubscribeSuccess;

  /// No description provided for @channelLeave.
  ///
  /// In en, this message translates to:
  /// **'Leave channel'**
  String get channelLeave;

  /// No description provided for @channelLeaveConfirm.
  ///
  /// In en, this message translates to:
  /// **'Leave this channel?'**
  String get channelLeaveConfirm;

  /// No description provided for @channelDelete.
  ///
  /// In en, this message translates to:
  /// **'Delete channel'**
  String get channelDelete;

  /// No description provided for @channelDeleteConfirm.
  ///
  /// In en, this message translates to:
  /// **'Delete this channel? This cannot be undone.'**
  String get channelDeleteConfirm;

  /// No description provided for @navContacts.
  ///
  /// In en, this message translates to:
  /// **'Contacts'**
  String get navContacts;

  /// No description provided for @contactsInSuda.
  ///
  /// In en, this message translates to:
  /// **'In Suda'**
  String get contactsInSuda;

  /// No description provided for @contactsOnPhone.
  ///
  /// In en, this message translates to:
  /// **'From phone'**
  String get contactsOnPhone;

  /// No description provided for @contactsPermissionDenied.
  ///
  /// In en, this message translates to:
  /// **'Contact access denied. Allow in settings.'**
  String get contactsPermissionDenied;

  /// No description provided for @channelPinnedPosts.
  ///
  /// In en, this message translates to:
  /// **'Pinned posts'**
  String get channelPinnedPosts;

  /// No description provided for @channelSearchIn.
  ///
  /// In en, this message translates to:
  /// **'Search in channel'**
  String get channelSearchIn;

  /// No description provided for @channelPost.
  ///
  /// In en, this message translates to:
  /// **'Post'**
  String get channelPost;

  /// No description provided for @channelInvite.
  ///
  /// In en, this message translates to:
  /// **'Invite'**
  String get channelInvite;

  /// No description provided for @channelInviteUsernameHint.
  ///
  /// In en, this message translates to:
  /// **'@username'**
  String get channelInviteUsernameHint;

  /// No description provided for @channelInviteSent.
  ///
  /// In en, this message translates to:
  /// **'Invitation sent'**
  String get channelInviteSent;

  /// No description provided for @channelTreasury.
  ///
  /// In en, this message translates to:
  /// **'Treasury'**
  String get channelTreasury;

  /// No description provided for @treasuryBalance.
  ///
  /// In en, this message translates to:
  /// **'Treasury Balance'**
  String get treasuryBalance;

  /// No description provided for @treasuryTotalDonations.
  ///
  /// In en, this message translates to:
  /// **'Total Donations'**
  String get treasuryTotalDonations;

  /// No description provided for @treasuryTopDonors.
  ///
  /// In en, this message translates to:
  /// **'Top Donors'**
  String get treasuryTopDonors;

  /// No description provided for @treasuryRecentDonations.
  ///
  /// In en, this message translates to:
  /// **'Recent Donations'**
  String get treasuryRecentDonations;

  /// No description provided for @treasuryEmpty.
  ///
  /// In en, this message translates to:
  /// **'No donations yet'**
  String get treasuryEmpty;

  /// No description provided for @treasuryWithdraw.
  ///
  /// In en, this message translates to:
  /// **'Withdraw'**
  String get treasuryWithdraw;

  /// No description provided for @treasuryWithdrawHint.
  ///
  /// In en, this message translates to:
  /// **'Amount (SUDA)'**
  String get treasuryWithdrawHint;

  /// No description provided for @treasuryWithdrawSuccess.
  ///
  /// In en, this message translates to:
  /// **'Withdrawal submitted successfully'**
  String get treasuryWithdrawSuccess;

  /// No description provided for @treasuryInsufficientFunds.
  ///
  /// In en, this message translates to:
  /// **'Insufficient funds in treasury'**
  String get treasuryInsufficientFunds;

  /// No description provided for @treasuryWithdrawFailed.
  ///
  /// In en, this message translates to:
  /// **'Withdrawal failed. Please try again.'**
  String get treasuryWithdrawFailed;

  /// No description provided for @channelEdit.
  ///
  /// In en, this message translates to:
  /// **'Edit'**
  String get channelEdit;

  /// No description provided for @channelMiniApps.
  ///
  /// In en, this message translates to:
  /// **'Mini-apps'**
  String get channelMiniApps;

  /// No description provided for @channelSettings.
  ///
  /// In en, this message translates to:
  /// **'Settings'**
  String get channelSettings;

  /// No description provided for @channelSettingsTitle.
  ///
  /// In en, this message translates to:
  /// **'Channel settings'**
  String get channelSettingsTitle;

  /// No description provided for @channelSettingsSave.
  ///
  /// In en, this message translates to:
  /// **'Save'**
  String get channelSettingsSave;

  /// No description provided for @channelSettingsProfile.
  ///
  /// In en, this message translates to:
  /// **'Profile'**
  String get channelSettingsProfile;

  /// No description provided for @channelSettingsSaved.
  ///
  /// In en, this message translates to:
  /// **'Settings saved'**
  String get channelSettingsSaved;

  /// No description provided for @gatingSettingsTitle.
  ///
  /// In en, this message translates to:
  /// **'Token gating'**
  String get gatingSettingsTitle;

  /// No description provided for @gatingSettingsEnable.
  ///
  /// In en, this message translates to:
  /// **'Require SUDA to join'**
  String get gatingSettingsEnable;

  /// No description provided for @gatingSettingsMinBalance.
  ///
  /// In en, this message translates to:
  /// **'Subscription price (SUDA)'**
  String get gatingSettingsMinBalance;

  /// No description provided for @channelCommentsEnabledLabel.
  ///
  /// In en, this message translates to:
  /// **'Comments'**
  String get channelCommentsEnabledLabel;

  /// No description provided for @channelComments.
  ///
  /// In en, this message translates to:
  /// **'Comments'**
  String get channelComments;

  /// No description provided for @channelCommentsCount.
  ///
  /// In en, this message translates to:
  /// **'{count, plural, =0{No comments} =1{1 comment} other{{count} comments}}'**
  String channelCommentsCount(int count);

  /// No description provided for @channelCommentHint.
  ///
  /// In en, this message translates to:
  /// **'Write a comment...'**
  String get channelCommentHint;

  /// No description provided for @channelCommentsEmpty.
  ///
  /// In en, this message translates to:
  /// **'No comments yet'**
  String get channelCommentsEmpty;

  /// No description provided for @channelCommentSubscribePrompt.
  ///
  /// In en, this message translates to:
  /// **'Subscribe to comment'**
  String get channelCommentSubscribePrompt;

  /// No description provided for @channelCommentsDisabled.
  ///
  /// In en, this message translates to:
  /// **'Comments are disabled'**
  String get channelCommentsDisabled;

  /// No description provided for @commentEdited.
  ///
  /// In en, this message translates to:
  /// **'edited'**
  String get commentEdited;

  /// No description provided for @commentActionEdit.
  ///
  /// In en, this message translates to:
  /// **'Edit'**
  String get commentActionEdit;

  /// No description provided for @commentActionDelete.
  ///
  /// In en, this message translates to:
  /// **'Delete'**
  String get commentActionDelete;

  /// No description provided for @commentActionReply.
  ///
  /// In en, this message translates to:
  /// **'Reply'**
  String get commentActionReply;

  /// No description provided for @friendsTabFriends.
  ///
  /// In en, this message translates to:
  /// **'Friends'**
  String get friendsTabFriends;

  /// No description provided for @friendsTabRequests.
  ///
  /// In en, this message translates to:
  /// **'Requests'**
  String get friendsTabRequests;

  /// No description provided for @friendsAddFriend.
  ///
  /// In en, this message translates to:
  /// **'Add friend'**
  String get friendsAddFriend;

  /// No description provided for @friendsCancelRequest.
  ///
  /// In en, this message translates to:
  /// **'Cancel request'**
  String get friendsCancelRequest;

  /// No description provided for @friendsAccept.
  ///
  /// In en, this message translates to:
  /// **'Accept'**
  String get friendsAccept;

  /// No description provided for @friendsReject.
  ///
  /// In en, this message translates to:
  /// **'Reject'**
  String get friendsReject;

  /// No description provided for @friendsUnfriend.
  ///
  /// In en, this message translates to:
  /// **'Unfriend'**
  String get friendsUnfriend;

  /// No description provided for @friendsRequestSent.
  ///
  /// In en, this message translates to:
  /// **'Request sent'**
  String get friendsRequestSent;

  /// No description provided for @friendsYouAreFriends.
  ///
  /// In en, this message translates to:
  /// **'Friends'**
  String get friendsYouAreFriends;

  /// No description provided for @friendsIncoming.
  ///
  /// In en, this message translates to:
  /// **'Incoming'**
  String get friendsIncoming;

  /// No description provided for @friendsOutgoing.
  ///
  /// In en, this message translates to:
  /// **'Outgoing'**
  String get friendsOutgoing;

  /// No description provided for @friendsEmptyFriends.
  ///
  /// In en, this message translates to:
  /// **'No friends yet. Search for users to add.'**
  String get friendsEmptyFriends;

  /// No description provided for @friendsEmptyRequests.
  ///
  /// In en, this message translates to:
  /// **'No pending requests.'**
  String get friendsEmptyRequests;

  /// No description provided for @friendsSince.
  ///
  /// In en, this message translates to:
  /// **'Friends since {date}'**
  String friendsSince(String date);

  /// No description provided for @walletOpenTitle.
  ///
  /// In en, this message translates to:
  /// **'Wallet'**
  String get walletOpenTitle;

  /// No description provided for @walletLoadingError.
  ///
  /// In en, this message translates to:
  /// **'Could not open wallet. Tap to retry.'**
  String get walletLoadingError;

  /// No description provided for @changePasswordOld.
  ///
  /// In en, this message translates to:
  /// **'Current password'**
  String get changePasswordOld;

  /// No description provided for @changePasswordNew.
  ///
  /// In en, this message translates to:
  /// **'New password'**
  String get changePasswordNew;

  /// No description provided for @changePasswordConfirm.
  ///
  /// In en, this message translates to:
  /// **'Confirm new password'**
  String get changePasswordConfirm;

  /// No description provided for @changePasswordSuccess.
  ///
  /// In en, this message translates to:
  /// **'Password changed. Please log in again.'**
  String get changePasswordSuccess;

  /// No description provided for @changePasswordMismatch.
  ///
  /// In en, this message translates to:
  /// **'Passwords do not match.'**
  String get changePasswordMismatch;

  /// No description provided for @changePasswordMinLength.
  ///
  /// In en, this message translates to:
  /// **'Password must be at least 6 characters.'**
  String get changePasswordMinLength;

  /// No description provided for @changePasswordInvalid.
  ///
  /// In en, this message translates to:
  /// **'Current password is incorrect.'**
  String get changePasswordInvalid;

  /// No description provided for @activeSessionsTitle.
  ///
  /// In en, this message translates to:
  /// **'Active sessions'**
  String get activeSessionsTitle;

  /// No description provided for @activeSessionsSignOutAll.
  ///
  /// In en, this message translates to:
  /// **'Sign out from all other devices'**
  String get activeSessionsSignOutAll;

  /// No description provided for @activeSessionsSignOutAllConfirm.
  ///
  /// In en, this message translates to:
  /// **'Sign out all other sessions?'**
  String get activeSessionsSignOutAllConfirm;

  /// No description provided for @activeSessionsTerminate.
  ///
  /// In en, this message translates to:
  /// **'Terminate'**
  String get activeSessionsTerminate;

  /// No description provided for @activeSessionsCurrent.
  ///
  /// In en, this message translates to:
  /// **'This device'**
  String get activeSessionsCurrent;

  /// No description provided for @msgSudaTransfer.
  ///
  /// In en, this message translates to:
  /// **'{from} → {to}: {amount} SUDA'**
  String msgSudaTransfer(String from, String to, String amount);

  /// No description provided for @msgDonation.
  ///
  /// In en, this message translates to:
  /// **'{from} donated {amount} SUDA'**
  String msgDonation(String from, String amount);

  /// No description provided for @sysGroupCreated.
  ///
  /// In en, this message translates to:
  /// **'Group created'**
  String get sysGroupCreated;

  /// No description provided for @sysGenericEvent.
  ///
  /// In en, this message translates to:
  /// **'Updated'**
  String get sysGenericEvent;

  /// No description provided for @transferPending.
  ///
  /// In en, this message translates to:
  /// **'Transfer sent! Awaiting confirmation.'**
  String get transferPending;

  /// No description provided for @transferAmountHint.
  ///
  /// In en, this message translates to:
  /// **'Amount in SUDA'**
  String get transferAmountHint;

  /// No description provided for @transferNoteHint.
  ///
  /// In en, this message translates to:
  /// **'Note (optional)'**
  String get transferNoteHint;

  /// No description provided for @transferSelectRecipient.
  ///
  /// In en, this message translates to:
  /// **'Select recipient'**
  String get transferSelectRecipient;

  /// No description provided for @transferSendingTo.
  ///
  /// In en, this message translates to:
  /// **'Sending to @{username}'**
  String transferSendingTo(String username);

  /// No description provided for @donateSent.
  ///
  /// In en, this message translates to:
  /// **'Donation sent!'**
  String get donateSent;

  /// No description provided for @donateAmountHint.
  ///
  /// In en, this message translates to:
  /// **'Amount in SUDA'**
  String get donateAmountHint;

  /// No description provided for @donateMessageHint.
  ///
  /// In en, this message translates to:
  /// **'Message (optional, max 200)'**
  String get donateMessageHint;

  /// No description provided for @gatingAlertTitle.
  ///
  /// In en, this message translates to:
  /// **'Token-gated channel'**
  String get gatingAlertTitle;

  /// No description provided for @gatingAlertBody.
  ///
  /// In en, this message translates to:
  /// **'Subscribe for {amount} SUDA to read this channel.'**
  String gatingAlertBody(String amount);

  /// No description provided for @gatingOpenWallet.
  ///
  /// In en, this message translates to:
  /// **'Open Wallet'**
  String get gatingOpenWallet;

  /// No description provided for @gatingBack.
  ///
  /// In en, this message translates to:
  /// **'Back'**
  String get gatingBack;

  /// No description provided for @gatingPay.
  ///
  /// In en, this message translates to:
  /// **'Pay'**
  String get gatingPay;

  /// No description provided for @channelUsernameTaken.
  ///
  /// In en, this message translates to:
  /// **'Username already taken'**
  String get channelUsernameTaken;

  /// No description provided for @commentDeleteConfirm.
  ///
  /// In en, this message translates to:
  /// **'Delete this comment?'**
  String get commentDeleteConfirm;

  /// No description provided for @chatBlockedBanner.
  ///
  /// In en, this message translates to:
  /// **'This chat is unavailable'**
  String get chatBlockedBanner;

  /// No description provided for @blockedUserFallbackName.
  ///
  /// In en, this message translates to:
  /// **'Blocked User'**
  String get blockedUserFallbackName;

  /// No description provided for @profileSetPhoto.
  ///
  /// In en, this message translates to:
  /// **'Set Photo'**
  String get profileSetPhoto;

  /// No description provided for @channelNonMember.
  ///
  /// In en, this message translates to:
  /// **'Subscribe to view this channel\'s content'**
  String get channelNonMember;

  /// No description provided for @attachTransfer.
  ///
  /// In en, this message translates to:
  /// **'Transfer'**
  String get attachTransfer;

  /// No description provided for @channelJoinRequest.
  ///
  /// In en, this message translates to:
  /// **'Request to Join'**
  String get channelJoinRequest;

  /// No description provided for @channelCancelRequest.
  ///
  /// In en, this message translates to:
  /// **'Cancel Request'**
  String get channelCancelRequest;

  /// No description provided for @channelRequestSent.
  ///
  /// In en, this message translates to:
  /// **'Request Sent'**
  String get channelRequestSent;

  /// No description provided for @channelAcceptInvite.
  ///
  /// In en, this message translates to:
  /// **'Accept Invite'**
  String get channelAcceptInvite;

  /// No description provided for @channelDeclineInvite.
  ///
  /// In en, this message translates to:
  /// **'Decline'**
  String get channelDeclineInvite;

  /// No description provided for @channelSubscribeBanner.
  ///
  /// In en, this message translates to:
  /// **'Subscribe to participate in this channel'**
  String get channelSubscribeBanner;

  /// No description provided for @channelPrivateBanner.
  ///
  /// In en, this message translates to:
  /// **'This is a private channel'**
  String get channelPrivateBanner;

  /// No description provided for @settingsChannelInvitations.
  ///
  /// In en, this message translates to:
  /// **'Channel Invitations'**
  String get settingsChannelInvitations;

  /// No description provided for @channelInvitationsTitle.
  ///
  /// In en, this message translates to:
  /// **'Channel Invitations'**
  String get channelInvitationsTitle;

  /// No description provided for @channelInvitationsEmpty.
  ///
  /// In en, this message translates to:
  /// **'No pending invitations'**
  String get channelInvitationsEmpty;

  /// No description provided for @messageSentByMe.
  ///
  /// In en, this message translates to:
  /// **'You'**
  String get messageSentByMe;

  /// No description provided for @forwardNoTargets.
  ///
  /// In en, this message translates to:
  /// **'No chats to forward to'**
  String get forwardNoTargets;

  /// No description provided for @replyDeleted.
  ///
  /// In en, this message translates to:
  /// **'Original message was deleted'**
  String get replyDeleted;

  /// No description provided for @messageForwardedFrom.
  ///
  /// In en, this message translates to:
  /// **'Forwarded from {name}'**
  String messageForwardedFrom(String name);

  /// No description provided for @replyInvalidTarget.
  ///
  /// In en, this message translates to:
  /// **'Cannot reply: message was deleted or belongs to another chat'**
  String get replyInvalidTarget;

  /// No description provided for @replyToUnsentMessage.
  ///
  /// In en, this message translates to:
  /// **'Cannot reply to a message that hasn\'t been sent yet'**
  String get replyToUnsentMessage;

  /// No description provided for @chatPreviewTransfer.
  ///
  /// In en, this message translates to:
  /// **'{from} → {to}: {amount} SUDA'**
  String chatPreviewTransfer(String from, String to, String amount);

  /// No description provided for @chatPreviewDonation.
  ///
  /// In en, this message translates to:
  /// **'{from} donated {amount} SUDA'**
  String chatPreviewDonation(String from, String amount);

  /// No description provided for @chatBlockedMeBanner.
  ///
  /// In en, this message translates to:
  /// **'You cannot write to this user'**
  String get chatBlockedMeBanner;

  /// No description provided for @chatBlockedByMeBanner.
  ///
  /// In en, this message translates to:
  /// **'You have blocked this user'**
  String get chatBlockedByMeBanner;

  /// No description provided for @unblockAction.
  ///
  /// In en, this message translates to:
  /// **'Unblock'**
  String get unblockAction;
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
      <String>['en', 'kk', 'ru'].contains(locale.languageCode);

  @override
  bool shouldReload(_AppLocalizationsDelegate old) => false;
}

AppLocalizations lookupAppLocalizations(Locale locale) {
  // Lookup logic when only language code is specified.
  switch (locale.languageCode) {
    case 'en':
      return AppLocalizationsEn();
    case 'kk':
      return AppLocalizationsKk();
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
