// UI formatting helpers shared across features.

/// Returns a username rendered as a single-`@` handle, e.g. `@john`.
///
/// Strips any number of leading `@` from [username] before prepending exactly
/// one, so values that already include `@` never render as `@@john`. Returns an
/// empty string for null/blank input.
String formatHandle(String? username) {
  final raw = username?.trim() ?? '';
  if (raw.isEmpty) return '';
  return '@${raw.replaceFirst(RegExp(r'^@+'), '')}';
}
