/// Maps HTTP status codes and backend error keys to user-facing messages.
/// Format: "<code> <Title>. <Description>" — per CLAUDE.md coding standards.
class ErrorDictionary {
  const ErrorDictionary._();

  static String humanMessage({int? statusCode, String? errorKey}) {
    // Backend-specific error key takes priority for precise messages
    if (errorKey != null) {
      final byKey = _byKey[errorKey.toUpperCase()];
      if (byKey != null) return byKey;
    }

    // Fall back to HTTP status code
    return _byCode[statusCode] ?? 'Something went wrong. Please try again.';
  }

  static const _byCode = {
    400: '400 Bad Request. Check your input and try again.',
    401: '401 Unauthorized. Please log in again.',
    402: '402 Insufficient Balance. Top up your SUDA balance.',
    403: '403 Forbidden. You don\'t have permission for this action.',
    404: '404 Not Found.',
    409: '409 Conflict. Already exists or already handled.',
    422: '422 Validation Error. Check your input fields.',
    424: '424 Wallet Not Ready. Please verify your email first.',
    429: '429 Too Many Requests. Please wait a moment.',
    500: '500 Internal Server Error. Please try later.',
    503: '503 Service Unavailable. Please try later.',
  };

  static const _byKey = {
    'INSUFFICIENT_BALANCE': '402 Insufficient SUDA balance. Top up to continue.',
    'NO_WALLET': '424 Wallet not ready. Please verify your email first.',
    'NFT_REQUIRED': '403 NFT required to subscribe to this channel.',
    'GATING_REQUIREMENT_NOT_MET': '403 Requirements not met to join this channel.',
    'RATE_LIMITED': '429 Too many requests. Wait a moment and try again.',
    'INVALID_REQUEST': '400 Invalid request. Check your input.',
    'VALIDATION_ERROR': '400 Validation error. Check your input fields.',
    'ALREADY_TAKEN': '409 Username already taken. Choose another one.',
    'ALREADY_HANDLED': '409 Already processed.',
    'RECIPIENT_NOT_FOUND': '404 Recipient not found.',
    'CHAT_NOT_FOUND': '404 Chat not found.',
    'CANNOT_TRANSFER_TO_SELF': '400 Cannot transfer to yourself.',
    'RECIPIENT_HAS_NO_WALLET': '400 Recipient doesn\'t have a wallet yet.',
    'SENDER_NOT_IN_CHAT': '403 You are not a member of this chat.',
    'RECIPIENT_NOT_IN_CHAT': '403 Recipient is not a member of this chat.',
    'WALLET_NOT_FOUND': '424 Wallet not found. Verify your email to create one.',
    'CHANNEL_WALLET_NOT_FOUND': '424 Channel treasury not set up yet.',
    'NOT_OWNER': '403 Only the owner can perform this action.',
    'BLOCKED': '403 You cannot write to this user.',
    'FORBIDDEN': '403 You don\'t have permission for this action.',
    'UNAUTHORIZED': '401 Unauthorized. Please log in again.',
    'UNKNOWN_PACKAGE': '400 Unknown purchase package.',
    'PURCHASE_NOT_FOUND': '404 Purchase not found.',
    'INTERNAL_ERROR': '500 Internal Server Error. Please try later.',
  };
}
