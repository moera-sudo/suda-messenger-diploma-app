// ignore: unused_import
import 'package:intl/intl.dart' as intl;
import 'app_localizations.dart';

// ignore_for_file: type=lint

/// The translations for Russian (`ru`).
class AppLocalizationsRu extends AppLocalizations {
  AppLocalizationsRu([String locale = 'ru']) : super(locale);

  @override
  String get appTitle => 'Крипто Мессенджер';

  @override
  String get cancel => 'Отмена';

  @override
  String get confirm => 'Подтвердить';

  @override
  String get settingsTitle => 'Настройки';

  @override
  String get themeTitle => 'Тема';

  @override
  String get languageTitle => 'Язык';

  @override
  String get languageEnglish => 'Английский';

  @override
  String get languageRussian => 'Русский';

  @override
  String get languageKazakh => 'Казахский';

  @override
  String get authWelcomeTitle => 'Suda Messenger';

  @override
  String get authWelcomeSubtitle =>
      'Децентрализованный хаб для ваших переписок и крипто-активов. Подключайтесь, торгуйте и исследуйте новую экосистему web3.';

  @override
  String get authGetStarted => 'Начать';

  @override
  String get authLogin => 'Войти';

  @override
  String get authCreateAccount => 'Создать аккаунт';

  @override
  String get authWelcomeBack => 'С возвращением!';

  @override
  String get authExcited => 'Мы рады видеть вас снова!';

  @override
  String get authEmailOrUsername => 'Email или Имя пользователя';

  @override
  String get authEmail => 'Email';

  @override
  String get authPassword => 'Пароль';

  @override
  String get authUsername => 'Имя пользователя';

  @override
  String get authDisplayName => 'Отображаемое имя';

  @override
  String get authUsernameHint => '3–30 символов · буквы и цифры';

  @override
  String get authForgotPassword => 'Забыли пароль?';

  @override
  String get authNeedAccount => 'Нужен аккаунт?';

  @override
  String get authRegister => 'Регистрация';

  @override
  String get authAlreadyHaveAccount => 'Уже есть аккаунт?';

  @override
  String get authOrLoginWith => 'Или войти через';

  @override
  String get authJoinFuture => 'Присоединяйся к децентрализованному будущему';

  @override
  String get authVerifyTitle => 'Подтверждение личности';

  @override
  String get authVerifySubtitle =>
      'Введите 6-значный код, отправленный на ваш email, чтобы получить доступ к кошельку.';

  @override
  String get authVerifyButton => 'Подтвердить и продолжить';

  @override
  String get authResendCode => 'Отправить код снова';

  @override
  String get authCodeExpires => 'Код истекает через';

  @override
  String get errorGeneric => 'Что-то пошло не так';

  @override
  String get errorNetwork => 'Нет соединения с интернетом';

  @override
  String get placeholderErrorTitle => 'Ошибка';

  @override
  String get placeholderErrorMessage =>
      'Что-то пошло не так.\nМы уже разбираемся с этим.';

  @override
  String get placeholderInProgressTitle => 'В разработке';

  @override
  String get placeholderInProgressMessage =>
      'Этот функционал еще не готов.\nЗагляните сюда позже!';

  @override
  String get placeholderTestTitle => 'Тест';

  @override
  String get placeholderTestMessage =>
      'Это тестовый экран для проверки навигации.';

  @override
  String get placeholderNoContentTitle => 'Пусто';

  @override
  String get placeholderNoContentMessage => 'Здесь пока ничего нет.';

  @override
  String get buttonRetry => 'Повторить';

  @override
  String get authForgotPasswordTitle => 'Забыли пароль';

  @override
  String get authForgotPasswordSubtitle =>
      'Введите email и мы отправим код подтверждения.';

  @override
  String get authSendVerificationCode => 'Отправить код';

  @override
  String get authBackToLogin => 'Назад к входу';

  @override
  String get authResetPasswordTitle => 'Сброс пароля';

  @override
  String authEnterCodeSentTo(String email) {
    return 'Введите 6-значный код, отправленный на\n$email';
  }

  @override
  String get authNewPassword => 'Новый пароль';

  @override
  String get authConfirmPassword => 'Подтвердите пароль';

  @override
  String get authPasswordChangedSuccess => 'Пароль успешно изменён';

  @override
  String get authVerificationCodeSent =>
      'Если email существует, код для сброса пароля отправлен';

  @override
  String get authEmailRequired => 'Введите email';

  @override
  String get authEnterFullCode => 'Введите полный код подтверждения';

  @override
  String get authPasswordMinLength => 'Пароль должен быть не менее 6 символов';

  @override
  String get authPasswordsDoNotMatch => 'Пароли не совпадают';

  @override
  String get authEmailMissing => 'Email отсутствует. Запросите сброс снова.';

  @override
  String get chatTitle => 'Чаты';

  @override
  String get chatFindOrStart => 'Найти или начать беседу';

  @override
  String get chatFilterAll => 'Все';

  @override
  String get chatFilterUnread => 'Непрочитанные';

  @override
  String get chatFilterPersonal => 'Личные';

  @override
  String get chatFilterGroups => 'Группы';

  @override
  String get chatFilterChannels => 'Каналы';

  @override
  String get chatNewGroup => 'Новая группа';

  @override
  String get chatNewChannel => 'Новый канал';

  @override
  String get chatFindContact => 'Найти контакт';

  @override
  String get chatSavedMessages => 'Сохранённые сообщения';

  @override
  String get chatUnknownUser => 'Неизвестный пользователь';

  @override
  String get chatNoMessagesYet => 'Сообщений пока нет';

  @override
  String get chatStatusTyping => 'Печатает...';

  @override
  String get chatStatusOnline => 'В сети';

  @override
  String get chatStatusOffline => 'Не в сети';

  @override
  String get searchTitle => 'Поиск';

  @override
  String get searchHint => 'Поиск пользователей, чатов, сообщений…';

  @override
  String get searchNoUsersFound => 'Пользователи не найдены';

  @override
  String get searchEnterUsername => 'Введите имя для поиска';

  @override
  String get searchRecent => 'Недавние';

  @override
  String get searchUsers => 'Пользователи';

  @override
  String get searchChats => 'Чаты';

  @override
  String get searchMessages => 'Сообщения';

  @override
  String get searchClearQuery => 'Очистить';

  @override
  String searchNoResults(String query) {
    return 'Ничего не найдено по \"$query\"';
  }

  @override
  String get chatSearchTitle => 'Поиск в чате';

  @override
  String get chatSearchHint => 'Поиск сообщений…';

  @override
  String get profileUpdated => 'Профиль обновлён';

  @override
  String get avatarUpdated => 'Аватар обновлён';

  @override
  String get loggedOut => 'Вы вышли из аккаунта';

  @override
  String get navChat => 'Чаты';

  @override
  String get navHub => 'Хаб';

  @override
  String get navFriends => 'Друзья';

  @override
  String get navProfile => 'Профиль';

  @override
  String get themeSelectTitle => 'Выбор темы';

  @override
  String get chatLastSeenPrefix => 'был(а)';

  @override
  String get chatLastSeenJustNow => 'только что';

  @override
  String get chatLastSeenMinutes => 'мин назад';

  @override
  String get chatLastSeenHours => 'ч назад';

  @override
  String get messageReadersTitle => 'Прочитали';

  @override
  String get messageReadersEmpty => 'Никто ещё не прочитал это сообщение';

  @override
  String get messageActionEdit => 'Редактировать';

  @override
  String get messageActionCopy => 'Копировать';

  @override
  String get messageActionDeleteForMe => 'Удалить у меня';

  @override
  String get messageActionDeleteForEveryone => 'Удалить у всех';

  @override
  String get messageEditingLabel => 'Редактирование';

  @override
  String get messageEditCancel => 'Отменить редактирование';

  @override
  String get messageReply => 'Ответить';

  @override
  String get messageForward => 'Переслать';

  @override
  String get forwardTo => 'Переслать в';

  @override
  String get chatInfo => 'Информация о чате';

  @override
  String get chatMembers => 'Участники';

  @override
  String get chatAddMember => 'Добавить участника';

  @override
  String get chatLeave => 'Покинуть чат';

  @override
  String get chatDelete => 'Удалить чат';

  @override
  String get chatLeaveConfirm => 'Вы уверены, что хотите покинуть чат?';

  @override
  String get chatDeleteConfirm => 'Вы уверены, что хотите удалить чат?';

  @override
  String get chatMute => 'Заглушить';

  @override
  String get chatUnmute => 'Включить уведомления';

  @override
  String get chatDeleteForMe => 'Удалить у меня';

  @override
  String get chatDeleteForEveryone => 'Удалить у всех';

  @override
  String get chatDeleteGroup => 'Удалить группу';

  @override
  String get chatDeleteChannel => 'Удалить канал';

  @override
  String get chatDeleteConfirmForAll =>
      'Это удалит чат для всех участников. Продолжить?';

  @override
  String get chatDeleteConfirmForMe =>
      'Это скроет чат и очистит вашу историю. Продолжить?';

  @override
  String get chatEditGroup => 'Редактировать группу';

  @override
  String get chatGroupDescription => 'Описание';

  @override
  String get chatGroupDescriptionHint => 'О группе…';

  @override
  String get chatSaveChanges => 'Сохранить';

  @override
  String get chatAddMembers => 'Добавить участников';

  @override
  String get chatMakeAdmin => 'Назначить администратором';

  @override
  String get chatRemoveMember => 'Удалить';

  @override
  String get newChat => 'Новый диалог';

  @override
  String get newGroup => 'Новая группа';

  @override
  String get newChannel => 'Новый канал';

  @override
  String get groupName => 'Название группы';

  @override
  String get groupNameHint => 'Введите название группы';

  @override
  String get createGroup => 'Создать группу';

  @override
  String get selectMembers => 'Выбрать участников';

  @override
  String get roleOwner => 'Владелец';

  @override
  String get roleAdmin => 'Администратор';

  @override
  String get roleMember => 'Участник';

  @override
  String get profileTitle => 'Профиль';

  @override
  String get profileSudaBalance => 'Баланс SUDA';

  @override
  String get profileOpenWallet => 'Открыть кошелёк';

  @override
  String get profileChats => 'Чаты';

  @override
  String get profileChannels => 'Каналы';

  @override
  String get profileContacts => 'Контакты';

  @override
  String get profileEditProfile => 'Редактировать профиль';

  @override
  String get profileNotifications => 'Уведомления';

  @override
  String get profilePrivacy => 'Конфиденциальность';

  @override
  String get profileLanguage => 'Язык';

  @override
  String get profileBlocked => 'Заблокированные';

  @override
  String get profileAbout => 'О Suda';

  @override
  String get profileSignOut => 'Выйти';

  @override
  String get settingsAccount => 'Аккаунт';

  @override
  String get settingsPrivacy => 'Приватность';

  @override
  String get settingsNotifData => 'Уведомления и данные';

  @override
  String get settingsAppearance => 'Оформление';

  @override
  String get settingsSuda => 'Suda';

  @override
  String get settingsEmail => 'Почта';

  @override
  String get settingsChangePassword => 'Сменить пароль';

  @override
  String get settingsActiveSessions => 'Активные сессии';

  @override
  String get settingsShowOnline => 'Показывать статус «онлайн»';

  @override
  String get settingsReadReceipts => 'Отчёты о прочтении';

  @override
  String get settingsLastSeen => 'Видимость «был(а) в сети»';

  @override
  String get settingsLastSeenEveryone => 'Все';

  @override
  String get settingsLastSeenContacts => 'Контакты';

  @override
  String get settingsLastSeenNobody => 'Никто';

  @override
  String get settingsBlocked => 'Заблокированные';

  @override
  String get settingsPushNotifications => 'Push-уведомления';

  @override
  String get settingsAutoDownload => 'Автозагрузка медиа';

  @override
  String get settingsAutoAlways => 'Всегда';

  @override
  String get settingsAutoWifiOnly => 'Только Wi-Fi';

  @override
  String get settingsAutoNever => 'Никогда';

  @override
  String get settingsTheme => 'Тема';

  @override
  String get settingsWallet => 'Кошелёк';

  @override
  String get settingsMarketplace => 'Маркетплейс';

  @override
  String get settingsAboutSuda => 'О Suda';

  @override
  String get settingsPrivacyPolicy => 'Политика конфиденциальности';

  @override
  String get settingsSignOut => 'Выйти';

  @override
  String settingsVersion(String version) {
    return 'Suda v$version';
  }

  @override
  String get settingsAboutDesc =>
      'Децентрализованный хаб для ваших переписок и криптоактивов.';

  @override
  String get themePickerSubtitle => 'Предпросмотр применяется мгновенно';

  @override
  String get editProfileTitle => 'Редактировать профиль';

  @override
  String get editProfileSave => 'Сохранить';

  @override
  String get editProfileDisplayName => 'Отображаемое имя';

  @override
  String get editProfileFirstName => 'Имя';

  @override
  String get editProfileLastName => 'Фамилия';

  @override
  String get editProfileBio => 'О себе';

  @override
  String get editProfileBioHint => 'Расскажите о себе';

  @override
  String get editProfileOptional => 'Необязательно';

  @override
  String get editProfileUsername => 'Имя пользователя';

  @override
  String get editProfileUsernameHint =>
      'Имя пользователя можно менять не чаще раза в неделю.';

  @override
  String get walletOpenWallet => 'Открыть кошелёк';

  @override
  String get walletMarketplace => 'Маркетплейс';

  @override
  String get walletSuda => 'SUDA';

  @override
  String get messageEdited => 'изменено';

  @override
  String get messageForwardedLabel => 'Переслано';

  @override
  String get userProfileOnline => 'в сети';

  @override
  String userProfileLastSeen(String time) {
    return 'был(а) $time';
  }

  @override
  String get userProfileMessage => 'Написать';

  @override
  String get userProfileSendSuda => 'Отправить';

  @override
  String get userProfileDonate => 'Донат';

  @override
  String get userProfileMute => 'Заглушить';

  @override
  String get userProfileUnmute => 'Включить';

  @override
  String get userProfileBioSection => 'О СЕБЕ';

  @override
  String get userProfileWalletAddress => 'Адрес кошелька';

  @override
  String get userProfileSharedMedia => 'Общие медиа';

  @override
  String get userProfileNotifications => 'Уведомления';

  @override
  String get userProfileBlockUser => 'Заблокировать';

  @override
  String get userProfileUnblockUser => 'Разблокировать';

  @override
  String userProfileBlockConfirm(String name) {
    return 'Заблокировать $name?';
  }

  @override
  String get userProfileBlockConfirmMsg =>
      'Пользователь не сможет написать вам.';

  @override
  String get userProfileBlockSuccess => 'Пользователь заблокирован';

  @override
  String get userProfileUnblockSuccess => 'Пользователь разблокирован';

  @override
  String get contactsTitle => 'Контакты';

  @override
  String get contactsEmpty => 'Пока нет контактов';

  @override
  String get contactsAddContact => 'Добавить контакт';

  @override
  String get contactsRemove => 'Удалить';

  @override
  String get blockedUsersTitle => 'Заблокированные';

  @override
  String get blockedUsersEmpty => 'Нет заблокированных пользователей';

  @override
  String get blockedUsersUnblock => 'Разблокировать';

  @override
  String get chatPin => 'Закрепить';

  @override
  String get chatUnpin => 'Открепить';

  @override
  String get comingSoon => 'Скоро';

  @override
  String get attachPhoto => 'Фото / Видео';

  @override
  String get attachFile => 'Файл';

  @override
  String get attachVoice => 'Голос';

  @override
  String get voiceRecording => 'Запись…';

  @override
  String get voiceSlideCancelHint => 'Нажмите ✕ для отмены';

  @override
  String get messageInputHint => 'Сообщение…';

  @override
  String get messageTypeImage => 'Фото';

  @override
  String get messageTypeFile => 'Файл';

  @override
  String get messageTypeVoice => 'Голосовое сообщение';

  @override
  String get messageTypeVideo => 'Видео';

  @override
  String get openFile => 'Открыть';

  @override
  String get fileOpenError => 'Не удалось открыть файл';

  @override
  String get micPermissionDenied =>
      'Для записи голосовых нужен доступ к микрофону';

  @override
  String get voiceRecordError => 'Не удалось записать голосовое сообщение';

  @override
  String get saveToGallery => 'Сохранить в галерею';

  @override
  String get saveToGallerySuccess => 'Сохранено в галерею';

  @override
  String get saveToGalleryError => 'Не удалось сохранить в галерею';

  @override
  String get commentReplyDeleted => 'Удалённый комментарий';

  @override
  String get sharedMediaTitle => 'Общие медиа';

  @override
  String get sharedMediaTabMedia => 'Медиа';

  @override
  String get sharedMediaPhotos => 'Фото';

  @override
  String get sharedMediaVideos => 'Видео';

  @override
  String get sharedMediaFiles => 'Файлы';

  @override
  String get sharedMediaAudio => 'Аудио';

  @override
  String get sharedMediaEmpty => 'Пока здесь пусто';

  @override
  String get uploadFailed => 'Ошибка загрузки. Повторите.';

  @override
  String get uploadingMedia => 'Загрузка…';

  @override
  String get channelNewTitle => 'Новый канал';

  @override
  String get channelName => 'Название канала';

  @override
  String get channelHandle => 'Хэндл';

  @override
  String get channelHandleHint => '@username';

  @override
  String get channelDescription => 'Описание';

  @override
  String get channelDescriptionHint => 'Расскажите о канале';

  @override
  String get channelVisibilityLabel => 'Видимость';

  @override
  String get channelPublic => 'Публичный';

  @override
  String get channelPrivate => 'Приватный';

  @override
  String get channelCreate => 'Создать канал';

  @override
  String channelSubscribersCount(int count) {
    return '$count подписчиков';
  }

  @override
  String get channelSubscribe => 'Подписаться';

  @override
  String get channelUnsubscribe => 'Отписаться';

  @override
  String get channelTokenGated => 'Платная подписка';

  @override
  String channelTokenGatedDesc(String amount) {
    return 'Подпишитесь за $amount SUDA, чтобы получить доступ.';
  }

  @override
  String channelSubscribeForPrice(String price) {
    return 'Подписаться за $price SUDA';
  }

  @override
  String get channelBalanceLabel => 'Ваш баланс';

  @override
  String get channelUnlock => 'Разблокировать и подписаться';

  @override
  String get channelTopUp => 'Пополнить и подписаться';

  @override
  String get channelConfirmSubscribe => 'Подтвердить и оплатить';

  @override
  String get channelSubscribeSuccess => 'Подписка оформлена!';

  @override
  String get channelInsufficientBalance =>
      'Недостаточно SUDA. Пополните кошелёк и повторите попытку.';

  @override
  String get channelNoWallet =>
      'Кошелёк не подключён. Создайте его в разделе Кошелёк.';

  @override
  String get channelSubscribeFailed => 'Ошибка оплаты. Попробуйте ещё раз.';

  @override
  String get channelUnsubscribeSuccess => 'Подписка отменена';

  @override
  String get channelLeave => 'Покинуть канал';

  @override
  String get channelLeaveConfirm => 'Покинуть этот канал?';

  @override
  String get channelDelete => 'Удалить канал';

  @override
  String get channelDeleteConfirm => 'Удалить этот канал? Действие необратимо.';

  @override
  String get navContacts => 'Контакты';

  @override
  String get contactsInSuda => 'В Suda';

  @override
  String get contactsOnPhone => 'С телефона';

  @override
  String get contactsPermissionDenied =>
      'Доступ к контактам запрещён. Разрешите в настройках.';

  @override
  String get channelPinnedPosts => 'Закреплённые посты';

  @override
  String get channelSearchIn => 'Поиск в канале';

  @override
  String get channelPost => 'Пост';

  @override
  String get channelInvite => 'Пригласить';

  @override
  String get channelInviteUsernameHint => '@username';

  @override
  String get channelInviteSent => 'Приглашение отправлено';

  @override
  String get channelTreasury => 'Казна';

  @override
  String get treasuryBalance => 'Баланс казны';

  @override
  String get treasuryTotalDonations => 'Всего донатов';

  @override
  String get treasuryTopDonors => 'Топ доноры';

  @override
  String get treasuryRecentDonations => 'Последние донаты';

  @override
  String get treasuryEmpty => 'Пока нет донатов';

  @override
  String get treasuryWithdraw => 'Вывести';

  @override
  String get treasuryWithdrawHint => 'Сумма (SUDA)';

  @override
  String get treasuryWithdrawSuccess => 'Вывод успешно отправлен';

  @override
  String get treasuryInsufficientFunds => 'Недостаточно средств в казне';

  @override
  String get treasuryWithdrawFailed => 'Ошибка вывода. Попробуйте ещё раз.';

  @override
  String get channelEdit => 'Изменить';

  @override
  String get channelMiniApps => 'Мини-приложения';

  @override
  String get channelSettings => 'Настройки';

  @override
  String get channelSettingsTitle => 'Настройки канала';

  @override
  String get channelSettingsSave => 'Сохранить';

  @override
  String get channelSettingsProfile => 'Профиль';

  @override
  String get channelSettingsSaved => 'Настройки сохранены';

  @override
  String get gatingSettingsTitle => 'Токен-доступ';

  @override
  String get gatingSettingsEnable => 'Требовать SUDA для входа';

  @override
  String get gatingSettingsMinBalance => 'Цена подписки (SUDA)';

  @override
  String get channelCommentsEnabledLabel => 'Комментарии';

  @override
  String get channelComments => 'Комментарии';

  @override
  String channelCommentsCount(int count) {
    String _temp0 = intl.Intl.pluralLogic(
      count,
      locale: localeName,
      other: '$count комментариев',
      few: '$count комментария',
      one: '1 комментарий',
      zero: 'Нет комментариев',
    );
    return '$_temp0';
  }

  @override
  String get channelCommentHint => 'Написать комментарий...';

  @override
  String get channelCommentsEmpty => 'Пока нет комментариев';

  @override
  String get channelCommentSubscribePrompt =>
      'Подпишитесь, чтобы комментировать';

  @override
  String get channelCommentsDisabled => 'Комментарии отключены';

  @override
  String get commentEdited => 'изменено';

  @override
  String get commentActionEdit => 'Изменить';

  @override
  String get commentActionDelete => 'Удалить';

  @override
  String get commentActionReply => 'Ответить';

  @override
  String get friendsTabFriends => 'Друзья';

  @override
  String get friendsTabRequests => 'Запросы';

  @override
  String get friendsAddFriend => 'Добавить в друзья';

  @override
  String get friendsCancelRequest => 'Отменить запрос';

  @override
  String get friendsAccept => 'Принять';

  @override
  String get friendsReject => 'Отклонить';

  @override
  String get friendsUnfriend => 'Удалить из друзей';

  @override
  String get friendsRequestSent => 'Запрос отправлен';

  @override
  String get friendsYouAreFriends => 'Вы друзья';

  @override
  String get friendsIncoming => 'Входящие';

  @override
  String get friendsOutgoing => 'Исходящие';

  @override
  String get friendsEmptyFriends =>
      'Нет друзей. Найдите пользователей через поиск.';

  @override
  String get friendsEmptyRequests => 'Нет входящих запросов.';

  @override
  String friendsSince(String date) {
    return 'Друзья с $date';
  }

  @override
  String get walletOpenTitle => 'Кошелёк';

  @override
  String get walletLoadingError =>
      'Не удалось открыть кошелёк. Нажмите для повтора.';

  @override
  String get changePasswordOld => 'Текущий пароль';

  @override
  String get changePasswordNew => 'Новый пароль';

  @override
  String get changePasswordConfirm => 'Подтвердите новый пароль';

  @override
  String get changePasswordSuccess =>
      'Пароль изменён. Пожалуйста, войдите заново.';

  @override
  String get changePasswordMismatch => 'Пароли не совпадают.';

  @override
  String get changePasswordMinLength =>
      'Пароль должен содержать минимум 6 символов.';

  @override
  String get changePasswordInvalid => 'Неверный текущий пароль.';

  @override
  String get activeSessionsTitle => 'Активные сессии';

  @override
  String get activeSessionsSignOutAll => 'Выйти на всех других устройствах';

  @override
  String get activeSessionsSignOutAllConfirm => 'Завершить все другие сессии?';

  @override
  String get activeSessionsTerminate => 'Завершить';

  @override
  String get activeSessionsCurrent => 'Это устройство';

  @override
  String msgSudaTransfer(String from, String to, String amount) {
    return '$from → $to: $amount SUDA';
  }

  @override
  String msgDonation(String from, String amount) {
    return '$from пожертвовал $amount SUDA';
  }

  @override
  String get sysGroupCreated => 'Группа создана';

  @override
  String get sysGenericEvent => 'Обновление';

  @override
  String get transferPending => 'Перевод отправлен! Ожидаем подтверждения.';

  @override
  String get transferAmountHint => 'Сумма в SUDA';

  @override
  String get transferNoteHint => 'Заметка (необязательно)';

  @override
  String get transferSelectRecipient => 'Выберите получателя';

  @override
  String transferSendingTo(String username) {
    return 'Отправка @$username';
  }

  @override
  String get donateSent => 'Донат отправлен!';

  @override
  String get donateAmountHint => 'Сумма в SUDA';

  @override
  String get donateMessageHint => 'Сообщение (необязательно, макс. 200)';

  @override
  String get gatingAlertTitle => 'Канал с токен-доступом';

  @override
  String gatingAlertBody(String amount) {
    return 'Подпишитесь за $amount SUDA, чтобы читать канал.';
  }

  @override
  String get gatingOpenWallet => 'Открыть кошелёк';

  @override
  String get gatingBack => 'Назад';

  @override
  String get gatingPay => 'Оплатить';

  @override
  String get channelUsernameTaken => 'Handle уже занят';

  @override
  String get commentDeleteConfirm => 'Удалить комментарий?';

  @override
  String get chatBlockedBanner => 'Этот чат недоступен';

  @override
  String get blockedUserFallbackName => 'Заблокированный пользователь';

  @override
  String get profileSetPhoto => 'Сменить фото';

  @override
  String get channelNonMember => 'Подпишитесь, чтобы увидеть контент канала';

  @override
  String get attachTransfer => 'Перевод';

  @override
  String get channelJoinRequest => 'Запросить вступление';

  @override
  String get channelCancelRequest => 'Отменить запрос';

  @override
  String get channelRequestSent => 'Запрос отправлен';

  @override
  String get channelAcceptInvite => 'Принять приглашение';

  @override
  String get channelDeclineInvite => 'Отклонить';

  @override
  String get channelSubscribeBanner =>
      'Подпишитесь, чтобы участвовать в канале';

  @override
  String get channelPrivateBanner => 'Это закрытый канал';

  @override
  String get settingsChannelInvitations => 'Приглашения в каналы';

  @override
  String get channelInvitationsTitle => 'Приглашения в каналы';

  @override
  String get channelInvitationsEmpty => 'Нет активных приглашений';

  @override
  String get messageSentByMe => 'Вы';

  @override
  String get forwardNoTargets => 'Нет чатов для пересылки';

  @override
  String get replyDeleted => 'Оригинальное сообщение удалено';

  @override
  String messageForwardedFrom(String name) {
    return 'Переслано от $name';
  }

  @override
  String get replyInvalidTarget =>
      'Не удалось ответить: сообщение удалено или из другого чата';

  @override
  String get replyToUnsentMessage =>
      'Нельзя отвечать на неотправленное сообщение';

  @override
  String chatPreviewTransfer(String from, String to, String amount) {
    return '$from → $to: $amount SUDA';
  }

  @override
  String chatPreviewDonation(String from, String amount) {
    return '$from пожертвовал $amount SUDA';
  }

  @override
  String get chatBlockedMeBanner => 'Вы не можете писать этому пользователю';

  @override
  String get chatBlockedByMeBanner => 'Вы заблокировали этого пользователя';

  @override
  String get unblockAction => 'Разблокировать';
}
