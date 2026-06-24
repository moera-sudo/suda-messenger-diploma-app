import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import '../../../../app/navigation/app_routes.dart';
import '../bloc/auth_bloc.dart';
import '../../../chat/presentation/bloc/socket_bloc.dart';

class SplashPage extends StatefulWidget {
  const SplashPage({super.key});

  @override
  State<SplashPage> createState() => _SplashPageState();
}

class _SplashPageState extends State<SplashPage> {
  @override
  void initState() {
    super.initState();
    // Проверяем авторизацию при старте
    context.read<AuthBloc>().add(AuthCheckRequested());
  }

  @override
  Widget build(BuildContext context) {
    return BlocListener<AuthBloc, AuthState>(
      listener: (context, state) {
        if (state.status == AuthStatus.authenticated) {
          // ЕСЛИ УЖЕ АВТОРИЗОВАН - ЗАПУСКАЕМ СОКЕТ!
          context.read<SocketBloc>().add(SocketConnect());
          context.goNamed(AppRoutes.chats);
        } else if (state.status == AuthStatus.unauthenticated) {
          context.goNamed(AppRoutes.welcome);
        }
      },
      child: const Scaffold(
        backgroundColor: Color(0xFF0F1020),
        body: Center(
          child: CircularProgressIndicator(color: Color(0xFF29DBF2)),
        ),
      ),
    );
  }
}
