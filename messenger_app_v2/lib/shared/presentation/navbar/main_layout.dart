import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import 'custom_bottom_nav.dart';
import 'sidebar_drawer.dart';

class MainLayout extends StatelessWidget {
  final Widget child;

  const MainLayout({super.key, required this.child});

  @override
  Widget build(BuildContext context) {
    final String location = GoRouterState.of(context).matchedLocation;

    int currentIndex = 0;
    if (location.startsWith('/friends')) currentIndex = 1;
    if (location.startsWith('/profile')) currentIndex = 2;

    return Scaffold(
      backgroundColor: Theme.of(context).scaffoldBackgroundColor,
      drawer: const SidebarDrawer(),
      body: child,
      bottomNavigationBar: CustomBottomNav(
        currentIndex: currentIndex,
        onTap: (index) {
          switch (index) {
            case 0:
              context.go('/chats');
            case 1:
              context.go('/friends');
            case 2:
              context.go('/profile');
          }
        },
      ),
    );
  }
}
