# Suda Messenger — Mobile Client

Technical reference for the **client application only**. The backend, infrastructure, and
deployment are out of scope and intentionally not described here.

---

## 1. Overview

Suda Messenger is a cross-platform **Flutter** messaging client (Android-first, portrait
orientation locked at startup). It provides direct and group chats, broadcast channels with
comments, friends/contacts, search, media and voice messages, push notifications, and an
embedded crypto wallet (SUDA token) exposed as a WebView mini-app with channel treasuries,
donations, and token-gated subscriptions.

- **Framework:** Flutter (Material 3)
- **Language:** Dart, SDK constraint `^3.9.2`
- **Realtime:** WebSocket event stream alongside a REST API
- **Architecture:** feature-first, layered (`data` / `domain` / `presentation`)

---

## 2. Architecture

The client follows a **feature-first** organization with a clean **three-layer** convention
applied per feature:

| Layer | Responsibility | Contents |
| --- | --- | --- |
| `presentation` | UI and state | Pages, widgets, Blocs/Cubits (events + states) |
| `domain` | Contracts | Abstract repository interfaces, domain models |
| `data` | Implementation | Repository implementations, DTO models (+ generated `*.g.dart`), remote/local sources |

**Dependency direction:** `presentation → domain ← data`. The presentation layer depends
only on domain repository *interfaces*; concrete `*_repository_impl` classes in `data`
implement those interfaces and are wired at runtime through dependency injection. This keeps
Blocs free of HTTP/serialization details.

Cross-cutting concerns live outside features:

- `lib/app/` — application bootstrap, DI container, theming, navigation, configuration.
- `lib/shared/` — shared infrastructure (API client, socket client, storage, logger,
  notifications, feedback) and the reusable design-system widgets.

---

## 3. Project Structure

```
lib/
├── main.dart                      # Entry point → bootstrap()
├── firebase_options.dart          # FCM platform config
├── app/
│   ├── bootstrap.dart             # Async init: Firebase, Hive, DI, notifications, health check
│   ├── app.dart                   # MaterialApp.router + root MultiBlocProvider
│   ├── lifespan.dart              # WidgetsBindingObserver (app lifecycle)
│   ├── DI/
│   │   ├── get_it.dart            # @InjectableInit, configureDependencies()
│   │   ├── register_module.dart   # @module third-party registrations
│   │   └── get_it.config.dart     # generated DI wiring
│   ├── config/
│   │   ├── app_config.dart        # base URL, WS URL, timeouts
│   │   ├── const.dart
│   │   └── theme/                 # app_theme.dart, app_colors.dart
│   └── navigation/
│       ├── app_router.dart        # GoRouter definition
│       └── app_routes.dart        # route name constants
├── l10n/arb/                      # ARB sources + generated AppLocalizations (en/ru/kk)
├── shared/
│   ├── data/
│   │   ├── api/                   # api_client.dart, socket_client.dart, server_exception.dart
│   │   ├── dio/data/              # dio_module.dart, auth_interceptor.dart
│   │   ├── auth/                  # current_user.dart, auth_event_bus.dart
│   │   ├── storage/              # secure_storage_client.dart, hive_storage.dart
│   │   ├── locale/                # shared_storage.dart (SharedPreferences)
│   │   ├── notifications/         # notification_service.dart
│   │   ├── contacts/              # device_contacts_service.dart
│   │   └── logger/                # app_logger.dart
│   ├── domain/models/             # auth_token_pair, app_failure, api_error_response, socket_event, error_dictionary
│   └── presentation/
│       ├── bloc/                  # settings_cubit.dart (theme + locale)
│       ├── widgets/               # Suda* design-system widgets
│       ├── navbar/                # main_layout, custom_bottom_nav, sidebar_drawer
│       └── feedback/              # app_feedback.dart (snackbars)
└── features/
    ├── auth/        # login, register, email verify, forgot/reset password
    ├── chat/        # direct & group chats, messages, channels info/settings/comments, socket bloc
    ├── channels/    # broadcast channels, gating rules, treasury, donations
    ├── friends/     # friend list & requests
    ├── profile/     # own profile, edit, sessions, logout
    ├── user/        # other users, contacts, blocked users
    ├── search/      # global search + in-chat search
    ├── preferences/ # server-synced user preferences
    ├── media/       # upload/download pipeline, media URL resolution
    ├── wallet/      # WebView wallet mini-app host
    └── settings/    # settings hub (theme, password, sessions, invitations)
```

Each feature directory follows the same `data/ · domain/ · presentation/` split where
applicable (`settings` and `media` are thinner: `settings` is presentation-only, `media`
has no UI).

---

## 4. State Management

State is managed with **`flutter_bloc`** (both `Bloc` and `Cubit`). Blocs/Cubits are
resolved from the DI container and provided at the root (`AuthBloc`, `SocketBloc`,
`SettingsCubit`) or per-screen.

| Bloc / Cubit | Type | Owns |
| --- | --- | --- |
| `AuthBloc` | Bloc | Login, register, email verification, forgot/reset password, auth status |
| `SocketBloc` | Bloc | WebSocket connection lifecycle (connect/disconnect) |
| `ChatListBloc` | Bloc | Chat list, pin/mute/delete, mark-read, applies socket updates |
| `ChatDetailBloc` | Bloc | Messages (load/send/edit/delete), replies, media, members, typing, socket updates |
| `SearchBloc` | Bloc | Global search (users/chats/messages), 400 ms debounce |
| `ChatSearchBloc` | Bloc | In-chat message search, 400 ms debounce |
| `TreasuryCubit` | Cubit | Channel treasury stats, donations, withdrawal |
| `PreferencesCubit` | Cubit | User preferences, synced via `PREFERENCES_UPDATED` socket events |
| `SettingsCubit` | Cubit | Active theme (`AppThemeType`) and locale |
| `FriendsCubit` | Cubit | Friends and requests, synced via socket events |
| `ProfileBloc` | Bloc | Own profile load/save, avatar upload, wallet load, logout |
| `UserProfileBloc` | Bloc | Foreign user profile, block/unblock |
| `ContactsBloc` | Bloc | Device contacts and blocked users |

- **`stream_transform`** provides a custom `debounce()` event transformer used by the search
  Blocs (`events.debounce(400ms).switchMap(mapper)`).
- **`bloc_concurrency`** is available for explicit event-transformer concurrency control.
- **`equatable`** backs value equality for events and states.

---

## 5. Dependency Injection

DI uses **`get_it`** as the service locator with **`injectable`** code generation.

- `lib/app/DI/get_it.dart` exposes the global `sl = GetIt.instance` and
  `configureDependencies()` annotated with `@InjectableInit`.
- `lib/app/DI/register_module.dart` is a `@module` registering third-party singletons
  (`FlutterSecureStorage`) and `@preResolve` async dependencies (`SharedPreferences`).
- `lib/app/DI/get_it.config.dart` is the generated wiring (`sl.init()`), run via
  `build_runner`.

Registration conventions: `@lazySingleton` for services and repositories (`ApiClient`,
`SocketClient`, `AppLogger`, `*RepositoryImpl`), `@injectable`/`@factory` for Blocs and
Cubits. `configureDependencies()` is awaited in `bootstrap()` before the app runs.

---

## 6. Networking

REST access is built on **`dio`**.

- **`DioModule`** (`@module`) builds a lazy-singleton `Dio` with `baseUrl` from `AppConfig`,
  15 s connect/receive timeouts, and JSON headers.
- **Interceptor chain:** `AuthInterceptor` → `pretty_dio_logger` (`PrettyDioLogger`).
- **`ApiClient`** wraps `get/post/put/delete`, catching `DioException` and converting it to a
  typed **`ServerException`** (extracting backend error code + status). Connection/timeout
  failures are mapped to `NETWORK_ERROR`.
- **`AuthInterceptor`** (`lib/shared/data/dio/data/auth_interceptor.dart`):
  - Injects `Authorization: Bearer <accessToken>` on non-auth requests.
  - On **401**, runs a **single-flight refresh**: concurrent 401s coalesce onto one
    in-flight refresh future (preventing refresh-token rotation races), then the original
    request is retried with the new token.
  - If refresh fails (or no refresh token exists), it triggers a forced logout via
    **`AuthEventBus`** (`onForcedLogout` broadcast stream), decoupling the network layer from
    navigation. The root `AuthBloc`/`app.dart` listens and redirects to Welcome.

**Error model:** `DioException → ServerException → AppFailure`, where `AppFailure` has
subtypes `ServerFailure` (backend code), `NetworkFailure` (`NETWORK_ERROR`), and
`UnknownFailure`. **`ErrorDictionary`** maps status codes and backend error keys to
human-readable messages surfaced through `AppFeedback`.

---

## 7. Realtime / WebSocket

Realtime updates use **`web_socket_channel`** (`IOWebSocketChannel`).

- **`SocketClient`** (`lib/shared/data/api/socket_client.dart`) connects to `AppConfig.wsUrl`
  with the access token passed as a query parameter, awaits `channel.ready`, and exposes a
  **broadcast `Stream<SocketEvent>`** (`events`).
- **`SocketEvent`** (`lib/shared/domain/models/socket_event.dart`) is a typed
  `{type, payload}` DTO. Event types include `NEW_MESSAGE`, `MESSAGE_EDITED`,
  `MESSAGES_READ`, `USER_ONLINE`/`USER_OFFLINE`, `MEMBER_ADDED`/`MEMBER_REMOVED`,
  `REACTION_ADDED`, `MESSAGE_PINNED`, `CHAT_UPDATED`, `DONATION_RECEIVED`,
  `PREFERENCES_UPDATED`, etc.
- **Reconnection** uses exponential backoff (1 s → 2 s → 5 s → 10 s → 30 s cap). Intentional
  `disconnect()` suppresses auto-reconnect; stream errors/closes trigger it.
- **`SocketBloc`** drives the connection lifecycle; feature Blocs (e.g. `ChatListBloc`,
  `ChatDetailBloc`, `FriendsCubit`, `PreferencesCubit`) subscribe to the broadcast stream and
  fold relevant events into their own state.

---

## 8. Data & Persistence

Repositories return functional results using **`dartz`** `Either<AppFailure, T>` — `Left`
for failures, `Right` for success — so call sites handle errors via `fold(...)` instead of
exceptions.

Storage is split by purpose:

| Store | Library | Used for |
| --- | --- | --- |
| Secure storage | `flutter_secure_storage` | JWT access/refresh tokens (`SecureStorageClient`) |
| Key-value cache | `hive` / `hive_flutter` | Cached business objects (boxes via `HiveClient`) |
| Lightweight prefs | `shared_preferences` | Theme + locale (`PreferencesLocalSource`) |

The current user identity is derived from the access token via **`jwt_decoder`**
(`CurrentUser.id()` reads the `user_id`/`sub` claim).

---

## 9. Serialization & Models

- DTOs use **`json_annotation`** + **`json_serializable`** with generated `*.g.dart` files
  (e.g. `chat_models.g.dart`, `auth_request_models.g.dart`, `user_models.g.dart`,
  `friend_models.g.dart`, `media_model.g.dart`, `user_preferences.g.dart`, `gating_rule.g.dart`).
- **`equatable`** provides value equality for models, events, and states.
- **`build_runner`** + `json_serializable` + `injectable_generator` drive all code
  generation.

---

## 10. Navigation

Routing uses **`go_router`**.

- **`AppRouter`** is a lazy singleton holding a configured `GoRouter` (`lib/app/navigation/`),
  with route names centralized in `app_routes.dart`.
- A **`ShellRoute`** wraps the primary tabs (`/chats`, `/friends`, `/profile`) in
  `MainLayout`, which renders the persistent `CustomBottomNav`.
- Routes use path parameters (`:id`, `:postId`) and pass complex arguments through
  `state.extra`.
- Authentication screens are wrapped in a fixed Suda theme regardless of the user's selected
  theme.
- The root `BlocListener<AuthBloc>` performs imperative redirects on auth status changes
  (authenticated → apply preferences + register FCM; unauthenticated → tear down socket and
  go to Welcome).

---

## 11. UI / Design System

- **Material 3** (`useMaterial3: true`) with a custom `AppPaletteTheme` `ThemeExtension`
  exposing palette colors not covered by `ColorScheme` (`surfaceContainer`, `messageMeBg`,
  `messageOtherBg`, `textAccent`, `danger`, `success`, `glow`, …).
- **`AppThemeType`** offers **6 variants** (light/dark pairs): `suda`, `sudaEnlightened`,
  `teaChatsLight`, `teaChatsDark`, `etherealSky`, `etherealAbyss`. `AppTheme.getTheme()`
  builds the `ThemeData`; `SettingsCubit` selects and persists the active theme.
- **Fonts (variable):** `Manrope` (body/UI default), `SpaceGrotesk` (display/headings),
  `JetBrainsMono` (wallet addresses/amounts).
- **`Suda*` widget kit** (`lib/shared/presentation/widgets/`): `SudaButton` (primary/ghost/
  outline/dangerOutline/link, three sizes, loading state), `SudaTextField` (focus glow,
  error states), `SudaChip`, `SudaInfoRow`/`SudaInfoList`, `SudaAvatar` (cascading
  mediaUrl→mediaId→userId resolution with a process-wide URL cache), `SudaOtpCell`.
- Supporting UI libs: **`flutter_animate`** (entrance animations), **`shimmer`** (loading
  skeletons), **`scrollable_positioned_list`** (jump-to-message in chat).

---

## 12. Media

Upload follows a three-step pipeline in `media_repository_impl.dart`: **init** (request
media id + presigned URL) → **presigned PUT** (raw bytes) → **confirm**, returning a
`MediaResult` (id, url, type, filename, size). Downloads are cached to the temp directory.

| Concern | Library |
| --- | --- |
| In-app gallery picker (themed) | `wechat_assets_picker` |
| Generic file picker | `file_picker` |
| Camera / single image (avatars) | `image_picker` |
| Photo-library metadata | `photo_manager` |
| Remote image caching | `cached_network_image` |
| Zoom/pan image viewer | `photo_view` |
| Video playback + controls | `video_player` + `chewie` |
| Voice record/playback + waveform | `audio_waveforms` |
| Open downloaded documents | `open_file` |
| Temp/cache directories | `path_provider` |

Message bubbles render `IMAGE`/`VIDEO`/`VOICE`/`FILE` types accordingly; voice recording
(`.m4a`/AAC) requires `Permission.microphone`.

---

## 13. Localization

- Flutter **gen-l10n** with **`intl`** and `flutter_localizations`.
- ARB sources in `lib/l10n/arb/` for **3 locales**: English (`app_en.arb`), Russian
  (`app_ru.arb`), Kazakh (`app_kk.arb`); generated into `app_localizations*.dart`.
- A `context.l10n` extension (`l10n_extensions.dart`) shortcuts `AppLocalizations.of(context)`.
- Locale is selected via `SettingsCubit.changeLocale()` and persisted; server preferences can
  override it on login.

---

## 14. Notifications

Push uses **`firebase_messaging`** (FCM) with **`flutter_local_notifications`** for display
(`notification_service.dart`).

- **Setup (`NotificationService.setup()`)** in bootstrap: creates the Android channel
  (`messenger_messages`, high importance), requests permissions, and registers
  foreground/background/cold-start handlers.
- **Token registration (`registerToken()`)** runs after authentication, posts the FCM token
  to the backend, and re-registers on `onTokenRefresh`.
- **Foreground** messages render a local notification; **background/terminated** taps and
  cold-start `getInitialMessage()` stash the target `chat_id`, which `app.dart` consumes
  after auth to deep-link into the chat.

---

## 15. Wallet / Web3 (client side)

The wallet is implemented as a **WebView-hosted mini-app SPA**, not native UI.

- **`webview_flutter`** hosts the SPA (`WalletWebViewScreen` / `MiniAppWebViewScreen`).
  Authentication to the mini-app is handed off via an **HMAC-signed `initData`** query
  parameter; a `close://` URL scheme signals the host to dismiss the WebView. On-chain
  transactions are backend-driven.
- SUDA amounts are handled with **BigInt** wei↔display conversion helpers (`sudaToWei`,
  `weiToSudaDisplay`) to avoid floating-point overflow.
- **Token-gating** (`GatingRule`) models paid channel subscriptions (subscription price in
  wei, legacy min-balance, optional NFT collection). **Treasury/donations** UI
  (`TreasuryCubit`, `TreasuryPage`, `DonateSheet`) shows balances, top donors, donation
  history, and owner withdrawals.

---

## 16. Miscellaneous / Platform

| Library | Purpose |
| --- | --- |
| `flutter_contacts` | Device contacts enumeration (`DeviceContactsService`, permission-gated, session-cached) |
| `device_info_plus` | Device/user-agent string sent at login |
| `permission_handler` | Runtime permissions (microphone for voice) |
| `uuid` | Client-side ID generation |
| `logger` | Structured app logging (`AppLogger`), distinct from `pretty_dio_logger` |

---

## 17. Dependencies

### Runtime (`dependencies`)

| Package | Version | Purpose |
| --- | --- | --- |
| `flutter` (sdk) | — | Framework |
| `flutter_bloc` | ^9.1.1 | State management (Bloc/Cubit) |
| `bloc_concurrency` | ^0.3.0 | Event-transformer concurrency control |
| `stream_transform` | ^2.1.1 | Stream operators (search debounce) |
| `get_it` | ^9.2.0 | Service locator |
| `injectable` | ^2.7.1+4 | DI annotations / codegen |
| `go_router` | ^17.1.0 | Declarative navigation |
| `dio` | ^5.9.1 | HTTP client |
| `pretty_dio_logger` | ^1.4.0 | Dio request/response logging |
| `web_socket_channel` | ^3.0.3 | WebSocket realtime |
| `dartz` | ^0.10.1 | Functional `Either` error handling |
| `equatable` | ^2.0.8 | Value equality |
| `json_annotation` | ^4.10.0 | JSON serialization annotations |
| `jwt_decoder` | ^2.0.1 | JWT claim decoding |
| `flutter_secure_storage` | ^10.0.0 | Encrypted token storage |
| `hive` | ^2.2.3 | Local key-value DB |
| `hive_flutter` | ^1.1.0 | Hive Flutter bindings |
| `shared_preferences` | ^2.5.4 | Lightweight settings store |
| `firebase_core` | ^4.4.0 | Firebase init |
| `firebase_messaging` | ^16.1.1 | FCM push |
| `flutter_local_notifications` | ^20.0.0 | Local notification display |
| `intl` | ^0.20.2 | i18n formatting / gen-l10n |
| `cached_network_image` | ^3.4.1 | Remote image caching |
| `image_picker` | ^1.2.1 | Camera/gallery single picks |
| `file_picker` | ^10.3.10 | File browsing |
| `photo_manager` | ^3.9.0 | Photo-library access |
| `wechat_assets_picker` | ^9.5.0 | Themed in-app gallery picker |
| `photo_view` | ^0.15.0 | Image zoom/pan viewer |
| `video_player` | ^2.9.2 | Video playback |
| `chewie` | ^1.8.5 | Video player UI |
| `audio_waveforms` | ^1.2.0 | Voice record/playback + waveform |
| `open_file` | ^3.5.11 | Open downloaded files |
| `path_provider` | ^2.1.5 | App directories |
| `permission_handler` | ^12.0.1 | Runtime permissions |
| `flutter_contacts` | ^1.1.9+2 | Device contacts |
| `device_info_plus` | ^12.3.0 | Device info / user agent |
| `flutter_animate` | ^4.5.2 | UI animations |
| `shimmer` | ^3.0.0 | Loading skeletons |
| `scrollable_positioned_list` | ^0.3.8 | Index-addressable scroll list |
| `webview_flutter` | ^4.13.1 | WebView host (wallet mini-app) |
| `logger` | ^2.6.2 | Structured logging |
| `uuid` | ^4.5.2 | UUID generation |

### Development (`dev_dependencies`)

| Package | Version | Purpose |
| --- | --- | --- |
| `flutter_test` (sdk) | — | Unit/widget testing |
| `flutter_lints` | ^5.0.0 | Lint rules |
| `build_runner` | ^2.9.0 | Code-generation runner |
| `json_serializable` | ^6.12.0 | JSON (de)serialization codegen |
| `injectable_generator` | ^2.9.0 | DI codegen |
| `mocktail` | ^1.0.4 | Mocking in tests |

---

## 18. Build & Run (client only)

```bash
# 1. Fetch dependencies
flutter pub get

# 2. Generate code (DI, JSON models). Run after changing annotated classes.
dart run build_runner build --delete-conflicting-outputs

# 3. (If ARB strings changed) regenerate localizations
flutter gen-l10n

# 4. Run on a connected device / emulator
flutter run
```

**Assets:** images under `assets/images/`, variable fonts under `assets/fonts/`
(`Manrope`, `SpaceGrotesk`, `JetBrainsMono`), declared in `pubspec.yaml`. Firebase
configuration (`firebase_options.dart`) must be present for notifications.

> The client targets a running backend; the API base URL and WebSocket URL are configured in
> `lib/app/config/app_config.dart`.
