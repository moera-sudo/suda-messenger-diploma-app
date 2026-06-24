import 'package:flutter/widgets.dart';
import 'package:messenger_app_v2/l10n/arb/app_localizations.dart';

extension AppLocalizationsX on BuildContext {
  AppLocalizations get l10n => AppLocalizations.of(this)!;
}