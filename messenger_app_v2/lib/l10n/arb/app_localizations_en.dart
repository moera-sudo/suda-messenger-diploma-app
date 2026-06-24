// ignore: unused_import
import 'package:intl/intl.dart' as intl;
import 'app_localizations.dart';

// ignore_for_file: type=lint

/// The translations for English (`en`).
class AppLocalizationsEn extends AppLocalizations {
  AppLocalizationsEn([String locale = 'en']) : super(locale);

  @override
  String get appTitle => 'Crypto Messenger';

  @override
  String get cancel => 'Cancel';

  @override
  String get confirm => 'Confirm';

  @override
  String get settingsTitle => 'Settings';

  @override
  String get themeTitle => 'Theme';

  @override
  String get languageTitle => 'Language';

  @override
  String get languageEnglish => 'English';

  @override
  String get languageRussian => 'Russian';

  @override
  String get languageKazakh => 'Kazakh';

  @override
  String get authWelcomeTitle => 'Suda Messenger';

  @override
  String get authWelcomeSubtitle =>
      'The decentralized hub for your conversations and crypto assets. Connect, trade, and explore the new web3 ecosystem.';

  @override
  String get authGetStarted => 'Get Started';

  @override
  String get authLogin => 'Log In';

  @override
  String get authCreateAccount => 'Create Account';

  @override
  String get authWelcomeBack => 'Welcome Back!';

  @override
  String get authExcited => 'We\'re so excited to see you again!';

  @override
  String get authEmailOrUsername => 'Email or Username';

  @override
  String get authEmail => 'Email';

  @override
  String get authPassword => 'Password';

  @override
  String get authUsername => 'Username';

  @override
  String get authDisplayName => 'Display Name';

  @override
  String get authUsernameHint => '3–30 chars · letters and digits';

  @override
  String get authForgotPassword => 'Forgot Password?';

  @override
  String get authNeedAccount => 'Need an account?';

  @override
  String get authRegister => 'Register';

  @override
  String get authAlreadyHaveAccount => 'Already have an account?';

  @override
  String get authOrLoginWith => 'Or login with';

  @override
  String get authJoinFuture => 'Join the decentralized future';

  @override
  String get authVerifyTitle => 'Verify Identity';

  @override
  String get authVerifySubtitle =>
      'Enter the 6-digit code sent to your email to access your crypto wallet.';

  @override
  String get authVerifyButton => 'Verify & Continue';

  @override
  String get authResendCode => 'Resend Code';

  @override
  String get authCodeExpires => 'Code expires in';

  @override
  String get errorGeneric => 'Something went wrong';

  @override
  String get errorNetwork => 'No internet connection';

  @override
  String get placeholderErrorTitle => 'Error';

  @override
  String get placeholderErrorMessage =>
      'Something went wrong.\nWe are already fixing it.';

  @override
  String get placeholderInProgressTitle => 'In Development';

  @override
  String get placeholderInProgressMessage =>
      'This feature is not ready yet.\nCheck back later!';

  @override
  String get placeholderTestTitle => 'Test';

  @override
  String get placeholderTestMessage =>
      'This is a test screen to check navigation.';

  @override
  String get placeholderNoContentTitle => 'Empty';

  @override
  String get placeholderNoContentMessage => 'There is nothing here yet.';

  @override
  String get buttonRetry => 'Retry';

  @override
  String get authForgotPasswordTitle => 'Forgot Password';

  @override
  String get authForgotPasswordSubtitle =>
      'Enter your email and we will send a verification code.';

  @override
  String get authSendVerificationCode => 'Send Verification Code';

  @override
  String get authBackToLogin => 'Back to Login';

  @override
  String get authResetPasswordTitle => 'Reset Password';

  @override
  String authEnterCodeSentTo(String email) {
    return 'Enter the 6-digit code sent to\n$email';
  }

  @override
  String get authNewPassword => 'New Password';

  @override
  String get authConfirmPassword => 'Confirm Password';

  @override
  String get authPasswordChangedSuccess => 'Password changed successfully';

  @override
  String get authVerificationCodeSent =>
      'If email exists, verification code for password reset sent';

  @override
  String get authEmailRequired => 'Email is required';

  @override
  String get authEnterFullCode => 'Enter full verification code';

  @override
  String get authPasswordMinLength => 'Password must be at least 6 characters';

  @override
  String get authPasswordsDoNotMatch => 'Passwords do not match';

  @override
  String get authEmailMissing => 'Email is missing. Request reset again.';

  @override
  String get chatTitle => 'Chats';

  @override
  String get chatFindOrStart => 'Find or start a conversation';

  @override
  String get chatFilterAll => 'All';

  @override
  String get chatFilterUnread => 'Unread';

  @override
  String get chatFilterPersonal => 'Personal';

  @override
  String get chatFilterGroups => 'Groups';

  @override
  String get chatFilterChannels => 'Channels';

  @override
  String get chatNewGroup => 'New group';

  @override
  String get chatNewChannel => 'New channel';

  @override
  String get chatFindContact => 'Find a contact';

  @override
  String get chatSavedMessages => 'Saved Messages';

  @override
  String get chatUnknownUser => 'Unknown User';

  @override
  String get chatNoMessagesYet => 'No messages yet';

  @override
  String get chatStatusTyping => 'Typing...';

  @override
  String get chatStatusOnline => 'Online';

  @override
  String get chatStatusOffline => 'Offline';

  @override
  String get searchTitle => 'Search';

  @override
  String get searchHint => 'Search users, chats, messages…';

  @override
  String get searchNoUsersFound => 'No users found';

  @override
  String get searchEnterUsername => 'Enter username to search';

  @override
  String get searchRecent => 'Recent';

  @override
  String get searchUsers => 'Users';

  @override
  String get searchChats => 'Chats';

  @override
  String get searchMessages => 'Messages';

  @override
  String get searchClearQuery => 'Clear';

  @override
  String searchNoResults(String query) {
    return 'No results for \"$query\"';
  }

  @override
  String get chatSearchTitle => 'Search in chat';

  @override
  String get chatSearchHint => 'Search messages…';

  @override
  String get profileUpdated => 'Profile updated';

  @override
  String get avatarUpdated => 'Avatar updated';

  @override
  String get loggedOut => 'Logged out';

  @override
  String get navChat => 'Chat';

  @override
  String get navHub => 'Hub';

  @override
  String get navFriends => 'Friends';

  @override
  String get navProfile => 'Profile';

  @override
  String get themeSelectTitle => 'Select Theme';

  @override
  String get chatLastSeenPrefix => 'last seen';

  @override
  String get chatLastSeenJustNow => 'just now';

  @override
  String get chatLastSeenMinutes => 'min ago';

  @override
  String get chatLastSeenHours => 'h ago';

  @override
  String get messageReadersTitle => 'Read by';

  @override
  String get messageReadersEmpty => 'No one has read this message yet';

  @override
  String get messageActionEdit => 'Edit';

  @override
  String get messageActionCopy => 'Copy';

  @override
  String get messageActionDeleteForMe => 'Delete for me';

  @override
  String get messageActionDeleteForEveryone => 'Delete for everyone';

  @override
  String get messageEditingLabel => 'Editing';

  @override
  String get messageEditCancel => 'Cancel editing';

  @override
  String get messageReply => 'Reply';

  @override
  String get messageForward => 'Forward';

  @override
  String get forwardTo => 'Forward to';

  @override
  String get chatInfo => 'Chat Info';

  @override
  String get chatMembers => 'Members';

  @override
  String get chatAddMember => 'Add member';

  @override
  String get chatLeave => 'Leave chat';

  @override
  String get chatDelete => 'Delete chat';

  @override
  String get chatLeaveConfirm => 'Are you sure you want to leave?';

  @override
  String get chatDeleteConfirm => 'Are you sure you want to delete this chat?';

  @override
  String get chatMute => 'Mute';

  @override
  String get chatUnmute => 'Unmute';

  @override
  String get chatDeleteForMe => 'Delete for me';

  @override
  String get chatDeleteForEveryone => 'Delete for everyone';

  @override
  String get chatDeleteGroup => 'Delete group';

  @override
  String get chatDeleteChannel => 'Delete channel';

  @override
  String get chatDeleteConfirmForAll =>
      'This will delete the chat for all members. Continue?';

  @override
  String get chatDeleteConfirmForMe =>
      'This will hide the chat and clear your history. Continue?';

  @override
  String get chatEditGroup => 'Edit group';

  @override
  String get chatGroupDescription => 'Description';

  @override
  String get chatGroupDescriptionHint => 'About this group…';

  @override
  String get chatSaveChanges => 'Save changes';

  @override
  String get chatAddMembers => 'Add members';

  @override
  String get chatMakeAdmin => 'Make admin';

  @override
  String get chatRemoveMember => 'Remove';

  @override
  String get newChat => 'New conversation';

  @override
  String get newGroup => 'New group';

  @override
  String get newChannel => 'New channel';

  @override
  String get groupName => 'Group name';

  @override
  String get groupNameHint => 'Enter group name';

  @override
  String get createGroup => 'Create group';

  @override
  String get selectMembers => 'Select members';

  @override
  String get roleOwner => 'Owner';

  @override
  String get roleAdmin => 'Admin';

  @override
  String get roleMember => 'Member';

  @override
  String get profileTitle => 'Profile';

  @override
  String get profileSudaBalance => 'SUDA Balance';

  @override
  String get profileOpenWallet => 'Open wallet';

  @override
  String get profileChats => 'Chats';

  @override
  String get profileChannels => 'Channels';

  @override
  String get profileContacts => 'Contacts';

  @override
  String get profileEditProfile => 'Edit profile';

  @override
  String get profileNotifications => 'Notifications';

  @override
  String get profilePrivacy => 'Privacy & security';

  @override
  String get profileLanguage => 'Language';

  @override
  String get profileBlocked => 'Blocked users';

  @override
  String get profileAbout => 'About Suda';

  @override
  String get profileSignOut => 'Sign out';

  @override
  String get settingsAccount => 'Account';

  @override
  String get settingsPrivacy => 'Privacy';

  @override
  String get settingsNotifData => 'Notifications & data';

  @override
  String get settingsAppearance => 'Appearance';

  @override
  String get settingsSuda => 'Suda';

  @override
  String get settingsEmail => 'Email';

  @override
  String get settingsChangePassword => 'Change password';

  @override
  String get settingsActiveSessions => 'Active sessions';

  @override
  String get settingsShowOnline => 'Show online status';

  @override
  String get settingsReadReceipts => 'Read receipts';

  @override
  String get settingsLastSeen => 'Last seen visibility';

  @override
  String get settingsLastSeenEveryone => 'Everyone';

  @override
  String get settingsLastSeenContacts => 'Contacts';

  @override
  String get settingsLastSeenNobody => 'Nobody';

  @override
  String get settingsBlocked => 'Blocked users';

  @override
  String get settingsPushNotifications => 'Push notifications';

  @override
  String get settingsAutoDownload => 'Auto-download media';

  @override
  String get settingsAutoAlways => 'Always';

  @override
  String get settingsAutoWifiOnly => 'Wi-Fi only';

  @override
  String get settingsAutoNever => 'Never';

  @override
  String get settingsTheme => 'Theme';

  @override
  String get settingsWallet => 'Wallet';

  @override
  String get settingsMarketplace => 'Marketplace';

  @override
  String get settingsAboutSuda => 'About Suda';

  @override
  String get settingsPrivacyPolicy => 'Privacy policy';

  @override
  String get settingsSignOut => 'Sign out';

  @override
  String settingsVersion(String version) {
    return 'Suda v$version';
  }

  @override
  String get settingsAboutDesc =>
      'The decentralized hub for your conversations and crypto assets.';

  @override
  String get themePickerSubtitle => 'Live preview applies instantly';

  @override
  String get editProfileTitle => 'Edit profile';

  @override
  String get editProfileSave => 'Save';

  @override
  String get editProfileDisplayName => 'Display name';

  @override
  String get editProfileFirstName => 'First name';

  @override
  String get editProfileLastName => 'Last name';

  @override
  String get editProfileBio => 'Bio';

  @override
  String get editProfileBioHint => 'Tell people about yourself';

  @override
  String get editProfileOptional => 'Optional';

  @override
  String get editProfileUsername => 'Username';

  @override
  String get editProfileUsernameHint =>
      'Username changes are limited to once per week.';

  @override
  String get walletOpenWallet => 'Open wallet';

  @override
  String get walletMarketplace => 'Marketplace';

  @override
  String get walletSuda => 'SUDA';

  @override
  String get messageEdited => 'edited';

  @override
  String get messageForwardedLabel => 'Forwarded';

  @override
  String get userProfileOnline => 'online';

  @override
  String userProfileLastSeen(String time) {
    return 'last seen $time';
  }

  @override
  String get userProfileMessage => 'Message';

  @override
  String get userProfileSendSuda => 'Send SUDA';

  @override
  String get userProfileDonate => 'Donate';

  @override
  String get userProfileMute => 'Mute';

  @override
  String get userProfileUnmute => 'Unmute';

  @override
  String get userProfileBioSection => 'BIO';

  @override
  String get userProfileWalletAddress => 'Wallet address';

  @override
  String get userProfileSharedMedia => 'Shared media';

  @override
  String get userProfileNotifications => 'Notifications';

  @override
  String get userProfileBlockUser => 'Block user';

  @override
  String get userProfileUnblockUser => 'Unblock';

  @override
  String userProfileBlockConfirm(String name) {
    return 'Block $name?';
  }

  @override
  String get userProfileBlockConfirmMsg =>
      'They won\'t be able to message you.';

  @override
  String get userProfileBlockSuccess => 'User blocked';

  @override
  String get userProfileUnblockSuccess => 'User unblocked';

  @override
  String get contactsTitle => 'Contacts';

  @override
  String get contactsEmpty => 'No contacts yet';

  @override
  String get contactsAddContact => 'Add contact';

  @override
  String get contactsRemove => 'Remove';

  @override
  String get blockedUsersTitle => 'Blocked users';

  @override
  String get blockedUsersEmpty => 'No blocked users';

  @override
  String get blockedUsersUnblock => 'Unblock';

  @override
  String get chatPin => 'Pin';

  @override
  String get chatUnpin => 'Unpin';

  @override
  String get comingSoon => 'Coming soon';

  @override
  String get attachPhoto => 'Photo / Video';

  @override
  String get attachFile => 'File';

  @override
  String get attachVoice => 'Voice';

  @override
  String get voiceRecording => 'Recording…';

  @override
  String get voiceSlideCancelHint => 'Tap ✕ to cancel';

  @override
  String get messageInputHint => 'Message…';

  @override
  String get messageTypeImage => 'Photo';

  @override
  String get messageTypeFile => 'File';

  @override
  String get messageTypeVoice => 'Voice message';

  @override
  String get messageTypeVideo => 'Video';

  @override
  String get openFile => 'Open';

  @override
  String get fileOpenError => 'Couldn\'t open the file';

  @override
  String get micPermissionDenied =>
      'Microphone access is required to record voice messages';

  @override
  String get voiceRecordError => 'Couldn\'t record the voice message';

  @override
  String get saveToGallery => 'Save to gallery';

  @override
  String get saveToGallerySuccess => 'Saved to gallery';

  @override
  String get saveToGalleryError => 'Couldn\'t save to gallery';

  @override
  String get commentReplyDeleted => 'Deleted comment';

  @override
  String get sharedMediaTitle => 'Shared media';

  @override
  String get sharedMediaTabMedia => 'Media';

  @override
  String get sharedMediaPhotos => 'Photos';

  @override
  String get sharedMediaVideos => 'Videos';

  @override
  String get sharedMediaFiles => 'Files';

  @override
  String get sharedMediaAudio => 'Audio';

  @override
  String get sharedMediaEmpty => 'Nothing here yet';

  @override
  String get uploadFailed => 'Upload failed. Try again.';

  @override
  String get uploadingMedia => 'Uploading…';

  @override
  String get channelNewTitle => 'New channel';

  @override
  String get channelName => 'Channel name';

  @override
  String get channelHandle => 'Handle';

  @override
  String get channelHandleHint => '@username';

  @override
  String get channelDescription => 'Description';

  @override
  String get channelDescriptionHint => 'Tell people about this channel';

  @override
  String get channelVisibilityLabel => 'Visibility';

  @override
  String get channelPublic => 'Public';

  @override
  String get channelPrivate => 'Private';

  @override
  String get channelCreate => 'Create channel';

  @override
  String channelSubscribersCount(int count) {
    return '$count subscribers';
  }

  @override
  String get channelSubscribe => 'Subscribe';

  @override
  String get channelUnsubscribe => 'Unsubscribe';

  @override
  String get channelTokenGated => 'Paid subscription';

  @override
  String channelTokenGatedDesc(String amount) {
    return 'Subscribe for $amount SUDA to access this channel.';
  }

  @override
  String channelSubscribeForPrice(String price) {
    return 'Subscribe for $price SUDA';
  }

  @override
  String get channelBalanceLabel => 'Your balance';

  @override
  String get channelUnlock => 'Unlock & subscribe';

  @override
  String get channelTopUp => 'Top up & subscribe';

  @override
  String get channelConfirmSubscribe => 'Confirm & Pay';

  @override
  String get channelSubscribeSuccess => 'Subscribed!';

  @override
  String get channelInsufficientBalance =>
      'Not enough SUDA. Top up your wallet and try again.';

  @override
  String get channelNoWallet =>
      'No wallet connected. Create one in the Wallet section.';

  @override
  String get channelSubscribeFailed => 'Payment failed. Please try again.';

  @override
  String get channelUnsubscribeSuccess => 'Unsubscribed';

  @override
  String get channelLeave => 'Leave channel';

  @override
  String get channelLeaveConfirm => 'Leave this channel?';

  @override
  String get channelDelete => 'Delete channel';

  @override
  String get channelDeleteConfirm =>
      'Delete this channel? This cannot be undone.';

  @override
  String get navContacts => 'Contacts';

  @override
  String get contactsInSuda => 'In Suda';

  @override
  String get contactsOnPhone => 'From phone';

  @override
  String get contactsPermissionDenied =>
      'Contact access denied. Allow in settings.';

  @override
  String get channelPinnedPosts => 'Pinned posts';

  @override
  String get channelSearchIn => 'Search in channel';

  @override
  String get channelPost => 'Post';

  @override
  String get channelInvite => 'Invite';

  @override
  String get channelInviteUsernameHint => '@username';

  @override
  String get channelInviteSent => 'Invitation sent';

  @override
  String get channelTreasury => 'Treasury';

  @override
  String get treasuryBalance => 'Treasury Balance';

  @override
  String get treasuryTotalDonations => 'Total Donations';

  @override
  String get treasuryTopDonors => 'Top Donors';

  @override
  String get treasuryRecentDonations => 'Recent Donations';

  @override
  String get treasuryEmpty => 'No donations yet';

  @override
  String get treasuryWithdraw => 'Withdraw';

  @override
  String get treasuryWithdrawHint => 'Amount (SUDA)';

  @override
  String get treasuryWithdrawSuccess => 'Withdrawal submitted successfully';

  @override
  String get treasuryInsufficientFunds => 'Insufficient funds in treasury';

  @override
  String get treasuryWithdrawFailed => 'Withdrawal failed. Please try again.';

  @override
  String get channelEdit => 'Edit';

  @override
  String get channelMiniApps => 'Mini-apps';

  @override
  String get channelSettings => 'Settings';

  @override
  String get channelSettingsTitle => 'Channel settings';

  @override
  String get channelSettingsSave => 'Save';

  @override
  String get channelSettingsProfile => 'Profile';

  @override
  String get channelSettingsSaved => 'Settings saved';

  @override
  String get gatingSettingsTitle => 'Token gating';

  @override
  String get gatingSettingsEnable => 'Require SUDA to join';

  @override
  String get gatingSettingsMinBalance => 'Subscription price (SUDA)';

  @override
  String get channelCommentsEnabledLabel => 'Comments';

  @override
  String get channelComments => 'Comments';

  @override
  String channelCommentsCount(int count) {
    String _temp0 = intl.Intl.pluralLogic(
      count,
      locale: localeName,
      other: '$count comments',
      one: '1 comment',
      zero: 'No comments',
    );
    return '$_temp0';
  }

  @override
  String get channelCommentHint => 'Write a comment...';

  @override
  String get channelCommentsEmpty => 'No comments yet';

  @override
  String get channelCommentSubscribePrompt => 'Subscribe to comment';

  @override
  String get channelCommentsDisabled => 'Comments are disabled';

  @override
  String get commentEdited => 'edited';

  @override
  String get commentActionEdit => 'Edit';

  @override
  String get commentActionDelete => 'Delete';

  @override
  String get commentActionReply => 'Reply';

  @override
  String get friendsTabFriends => 'Friends';

  @override
  String get friendsTabRequests => 'Requests';

  @override
  String get friendsAddFriend => 'Add friend';

  @override
  String get friendsCancelRequest => 'Cancel request';

  @override
  String get friendsAccept => 'Accept';

  @override
  String get friendsReject => 'Reject';

  @override
  String get friendsUnfriend => 'Unfriend';

  @override
  String get friendsRequestSent => 'Request sent';

  @override
  String get friendsYouAreFriends => 'Friends';

  @override
  String get friendsIncoming => 'Incoming';

  @override
  String get friendsOutgoing => 'Outgoing';

  @override
  String get friendsEmptyFriends => 'No friends yet. Search for users to add.';

  @override
  String get friendsEmptyRequests => 'No pending requests.';

  @override
  String friendsSince(String date) {
    return 'Friends since $date';
  }

  @override
  String get walletOpenTitle => 'Wallet';

  @override
  String get walletLoadingError => 'Could not open wallet. Tap to retry.';

  @override
  String get changePasswordOld => 'Current password';

  @override
  String get changePasswordNew => 'New password';

  @override
  String get changePasswordConfirm => 'Confirm new password';

  @override
  String get changePasswordSuccess => 'Password changed. Please log in again.';

  @override
  String get changePasswordMismatch => 'Passwords do not match.';

  @override
  String get changePasswordMinLength =>
      'Password must be at least 6 characters.';

  @override
  String get changePasswordInvalid => 'Current password is incorrect.';

  @override
  String get activeSessionsTitle => 'Active sessions';

  @override
  String get activeSessionsSignOutAll => 'Sign out from all other devices';

  @override
  String get activeSessionsSignOutAllConfirm => 'Sign out all other sessions?';

  @override
  String get activeSessionsTerminate => 'Terminate';

  @override
  String get activeSessionsCurrent => 'This device';

  @override
  String msgSudaTransfer(String from, String to, String amount) {
    return '$from → $to: $amount SUDA';
  }

  @override
  String msgDonation(String from, String amount) {
    return '$from donated $amount SUDA';
  }

  @override
  String get sysGroupCreated => 'Group created';

  @override
  String get sysGenericEvent => 'Updated';

  @override
  String get transferPending => 'Transfer sent! Awaiting confirmation.';

  @override
  String get transferAmountHint => 'Amount in SUDA';

  @override
  String get transferNoteHint => 'Note (optional)';

  @override
  String get transferSelectRecipient => 'Select recipient';

  @override
  String transferSendingTo(String username) {
    return 'Sending to @$username';
  }

  @override
  String get donateSent => 'Donation sent!';

  @override
  String get donateAmountHint => 'Amount in SUDA';

  @override
  String get donateMessageHint => 'Message (optional, max 200)';

  @override
  String get gatingAlertTitle => 'Token-gated channel';

  @override
  String gatingAlertBody(String amount) {
    return 'Subscribe for $amount SUDA to read this channel.';
  }

  @override
  String get gatingOpenWallet => 'Open Wallet';

  @override
  String get gatingBack => 'Back';

  @override
  String get gatingPay => 'Pay';

  @override
  String get channelUsernameTaken => 'Username already taken';

  @override
  String get commentDeleteConfirm => 'Delete this comment?';

  @override
  String get chatBlockedBanner => 'This chat is unavailable';

  @override
  String get blockedUserFallbackName => 'Blocked User';

  @override
  String get profileSetPhoto => 'Set Photo';

  @override
  String get channelNonMember => 'Subscribe to view this channel\'s content';

  @override
  String get attachTransfer => 'Transfer';

  @override
  String get channelJoinRequest => 'Request to Join';

  @override
  String get channelCancelRequest => 'Cancel Request';

  @override
  String get channelRequestSent => 'Request Sent';

  @override
  String get channelAcceptInvite => 'Accept Invite';

  @override
  String get channelDeclineInvite => 'Decline';

  @override
  String get channelSubscribeBanner =>
      'Subscribe to participate in this channel';

  @override
  String get channelPrivateBanner => 'This is a private channel';

  @override
  String get settingsChannelInvitations => 'Channel Invitations';

  @override
  String get channelInvitationsTitle => 'Channel Invitations';

  @override
  String get channelInvitationsEmpty => 'No pending invitations';

  @override
  String get messageSentByMe => 'You';

  @override
  String get forwardNoTargets => 'No chats to forward to';

  @override
  String get replyDeleted => 'Original message was deleted';

  @override
  String messageForwardedFrom(String name) {
    return 'Forwarded from $name';
  }

  @override
  String get replyInvalidTarget =>
      'Cannot reply: message was deleted or belongs to another chat';

  @override
  String get replyToUnsentMessage =>
      'Cannot reply to a message that hasn\'t been sent yet';

  @override
  String chatPreviewTransfer(String from, String to, String amount) {
    return '$from → $to: $amount SUDA';
  }

  @override
  String chatPreviewDonation(String from, String amount) {
    return '$from donated $amount SUDA';
  }

  @override
  String get chatBlockedMeBanner => 'You cannot write to this user';

  @override
  String get chatBlockedByMeBanner => 'You have blocked this user';

  @override
  String get unblockAction => 'Unblock';
}
