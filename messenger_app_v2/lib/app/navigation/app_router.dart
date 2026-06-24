import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:injectable/injectable.dart';

import 'package:messenger_app_v2/app/config/theme/app_theme.dart';
import 'package:messenger_app_v2/features/auth/auth_view.dart';
import 'package:messenger_app_v2/features/chat/presentation/pages/channel_comments_page.dart';
import 'package:messenger_app_v2/features/chat/presentation/pages/shared_media_page.dart';
import 'package:messenger_app_v2/features/chat/presentation/pages/channel_info_page.dart';
import 'package:messenger_app_v2/features/chat/presentation/pages/channel_settings_page.dart';
import 'package:messenger_app_v2/features/chat/presentation/pages/chat_detail_page.dart';
import 'package:messenger_app_v2/features/chat/data/models/chat_models.dart';
import 'package:messenger_app_v2/features/chat/presentation/pages/chat_info_page.dart';
import 'package:messenger_app_v2/features/chat/presentation/pages/chats_page.dart';
import 'package:messenger_app_v2/features/chat/presentation/pages/create_channel_page.dart';
import 'package:messenger_app_v2/features/chat/presentation/pages/create_group_page.dart';
import 'package:messenger_app_v2/features/profile/presentation/pages/edit_profile_page.dart';
import 'package:messenger_app_v2/features/profile/presentation/pages/profile_page.dart';
import 'package:messenger_app_v2/features/search/presentation/pages/chat_search_page.dart';
import 'package:messenger_app_v2/features/search/presentation/pages/search_page.dart';
import 'package:messenger_app_v2/features/friends/presentation/pages/friends_page.dart';
import 'package:messenger_app_v2/features/settings/presentation/pages/active_sessions_page.dart';
import 'package:messenger_app_v2/features/settings/presentation/pages/change_password_page.dart';
import 'package:messenger_app_v2/features/settings/presentation/pages/settings_page.dart';
import 'package:messenger_app_v2/features/settings/presentation/pages/theme_picker_page.dart';
import 'package:messenger_app_v2/features/channels/presentation/pages/channel_invitations_page.dart';
import 'package:messenger_app_v2/features/user/presentation/pages/blocked_users_page.dart';
import 'package:messenger_app_v2/features/user/presentation/pages/user_profile_page.dart';
import 'package:messenger_app_v2/features/channels/presentation/pages/treasury_page.dart';
import 'package:messenger_app_v2/features/wallet/presentation/pages/wallet_webview_screen.dart';
import '../../shared/presentation/navbar/main_layout.dart';
import 'app_routes.dart';

// Always render auth screens with Suda Dark regardless of the user-chosen theme.
Widget _suda(Widget child) => Theme(
      data: AppTheme.getTheme(AppThemeType.suda),
      child: child,
    );

@lazySingleton
class AppRouter {
  late final GoRouter router = GoRouter(
    initialLocation: '/',
    routes: [
      ShellRoute(
        builder: (context, state, child) => MainLayout(child: child),
        routes: [
          GoRoute(
            path: '/chats',
            name: AppRoutes.chats,
            builder: (context, state) => const ChatsPage(),
          ),
          GoRoute(
            path: '/profile',
            name: AppRoutes.profile,
            builder: (context, state) => const ProfilePage(),
          ),
          // Friends tab lives inside ShellRoute so the bottom-bar is visible.
          GoRoute(
            path: '/friends',
            name: AppRoutes.friends,
            builder: (context, state) => const FriendsPage(),
          ),
        ],
      ),

      // ── Auth screens — always Suda theme ────────────────────────────────
      GoRoute(
        path: '/',
        name: AppRoutes.splash,
        builder: (context, state) => _suda(const SplashPage()),
      ),
      GoRoute(
        path: '/welcome',
        name: AppRoutes.welcome,
        builder: (context, state) => _suda(const WelcomePage()),
      ),
      GoRoute(
        path: '/login',
        name: AppRoutes.login,
        builder: (context, state) => _suda(const LoginPage()),
      ),
      GoRoute(
        path: '/register',
        name: AppRoutes.register,
        builder: (context, state) => _suda(const RegisterPage()),
      ),
      GoRoute(
        path: '/verify',
        name: AppRoutes.verify,
        builder: (context, state) {
          final email = state.extra as String? ?? 'your email';
          return _suda(VerifyPage(email: email));
        },
      ),
      GoRoute(
        path: '/forgot-password',
        name: AppRoutes.forgotPassword,
        builder: (context, state) {
          final email = state.extra as String? ?? '';
          return _suda(ForgotPasswordPage(initialEmail: email));
        },
      ),
      GoRoute(
        path: '/reset-password',
        name: AppRoutes.resetPassword,
        builder: (context, state) {
          final email = state.extra as String? ?? '';
          return _suda(ResetPasswordPage(email: email));
        },
      ),

      // ── Search ──────────────────────────────────────────────────────────
      GoRoute(path: '/search', builder: (context, state) => const SearchPage()),

      // ── Chat detail — before /chats/:id to avoid /chats/new collision ──
      // Fix: /groups/new avoids matching against /chats/:id pattern.
      GoRoute(
        path: '/groups/new',
        name: AppRoutes.chatNew,
        builder: (context, state) => const CreateGroupPage(),
      ),

      GoRoute(
        path: '/chats/:id',
        name: AppRoutes.chatDetail,
        builder: (context, state) {
          final id = state.pathParameters['id'] ?? '';
          final extraMap = state.extra as Map<String, dynamic>? ?? {};
          final name = extraMap['name'] as String? ?? 'Chat';
          final interlocutorId = extraMap['interlocutorId'] as String? ?? '';
          final chatType = extraMap['chatType'] as String? ?? 'DIRECT';
          return ChatDetailPage(
            chatId: id,
            chatName: name,
            interlocutorId: interlocutorId,
            chatType: chatType,
          );
        },
      ),

      GoRoute(
        path: '/chats/:id/search',
        name: AppRoutes.chatSearch,
        builder: (context, state) {
          final chatId = state.pathParameters['id'] ?? '';
          final extra = state.extra as Map<String, dynamic>? ?? {};
          return ChatSearchPage(
            chatId: chatId,
            chatName: extra['chatName'] as String? ?? '',
          );
        },
      ),

      GoRoute(
        path: '/chats/:id/shared-media',
        name: AppRoutes.sharedMedia,
        builder: (context, state) {
          final chatId = state.pathParameters['id'] ?? '';
          final extra = state.extra as Map<String, dynamic>? ?? {};
          return SharedMediaPage(
            chatId: chatId,
            chatName: extra['chatName'] as String? ?? '',
          );
        },
      ),

      GoRoute(
        path: '/chats/:id/info',
        name: AppRoutes.chatInfo,
        builder: (context, state) {
          final chatId = state.pathParameters['id'] ?? '';
          final extra = state.extra as Map<String, dynamic>? ?? {};
          final chatType = extra['chatType'] as String? ?? 'GROUP';
          final myRole = extra['myRole'] as String? ?? 'MEMBER';
          if (chatType == 'CHANNEL') {
            final isMuted = extra['isMuted'] as bool? ?? false;
            return ChannelInfoPage(chatId: chatId, initialIsMuted: isMuted);
          }
          return ChatInfoPage(chatId: chatId, chatType: chatType, myRole: myRole);
        },
      ),

      GoRoute(
        path: '/channels/new',
        name: AppRoutes.channelNew,
        builder: (context, state) => const CreateChannelPage(),
      ),

      GoRoute(
        path: '/channels/:id/settings',
        name: AppRoutes.channelSettings,
        builder: (context, state) =>
            ChannelSettingsPage(chatId: state.pathParameters['id'] ?? ''),
      ),

      GoRoute(
        path: '/channels/:id/treasury',
        name: AppRoutes.channelTreasury,
        builder: (context, state) {
          final extra = state.extra as Map<String, dynamic>? ?? {};
          return TreasuryPage(
            channelId: state.pathParameters['id'] ?? '',
            isOwner: extra['isOwner'] as bool? ?? false,
          );
        },
      ),

      GoRoute(
        path: '/channels/:id/posts/:postId/comments',
        name: AppRoutes.channelComments,
        builder: (context, state) {
          final extra = state.extra as Map<String, dynamic>? ?? {};
          return ChannelCommentsPage(
            channelId: state.pathParameters['id'] ?? '',
            postId: int.tryParse(state.pathParameters['postId'] ?? '') ?? 0,
            post: extra['post'] as Message?,
            commentsEnabled: extra['commentsEnabled'] as bool? ?? false,
            isSubscriber: extra['isSubscriber'] as bool? ?? false,
            isAdmin: extra['isAdmin'] as bool? ?? false,
            currentUserId: extra['currentUserId'] as String? ?? '',
          );
        },
      ),

      GoRoute(
        path: '/user_profile/:id',
        name: AppRoutes.chatProfile,
        builder: (context, state) {
          final userId = state.pathParameters['id'] ?? '';
          final extra = state.extra as Map<String, dynamic>? ?? {};
          return UserProfilePage(
            userId: userId,
            chatId: extra['chatId'] as String?,
            chatName: extra['chatName'] as String?,
          );
        },
      ),

      GoRoute(
        path: '/blocked',
        name: AppRoutes.blocked,
        builder: (context, state) => const BlockedUsersPage(),
      ),

      // ── Settings ─────────────────────────────────────────────────────────
      GoRoute(
        path: '/settings',
        name: AppRoutes.settings,
        builder: (context, state) => const SettingsPage(),
      ),
      GoRoute(
        path: '/settings/edit-profile',
        name: AppRoutes.editProfile,
        builder: (context, state) => const EditProfilePage(),
      ),
      GoRoute(
        path: '/settings/theme',
        name: AppRoutes.themePicker,
        builder: (context, state) => const ThemePickerPage(),
      ),
      GoRoute(
        path: '/settings/change-password',
        name: AppRoutes.changePassword,
        builder: (context, state) => const ChangePasswordPage(),
      ),
      GoRoute(
        path: '/settings/active-sessions',
        name: AppRoutes.activeSessions,
        builder: (context, state) => const ActiveSessionsPage(),
      ),
      GoRoute(
        path: '/settings/channel-invitations',
        name: AppRoutes.channelInvitations,
        builder: (context, state) => const ChannelInvitationsPage(),
      ),

      // ── Wallet WebView ────────────────────────────────────────────────────
      GoRoute(
        path: '/wallet',
        name: AppRoutes.wallet,
        builder: (context, state) => const WalletWebViewScreen(),
      ),
    ],
  );
}
