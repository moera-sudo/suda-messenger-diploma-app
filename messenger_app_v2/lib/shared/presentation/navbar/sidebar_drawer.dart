import 'package:flutter/material.dart';

class SidebarDrawer extends StatelessWidget {
  const SidebarDrawer({super.key});

  @override
  Widget build(BuildContext context) {
    return Drawer(
      width: 80,
      backgroundColor: const Color(0xFF0C1114),
      child: SafeArea(
        child: Column(
          children: [
            const SizedBox(height: 16),
            // Home / DM Icon
            _SidebarItem(
              icon: Icons.chat_bubble, 
              color: const Color(0xFF29DBF2),
              isActive: true,
            ),
            const SizedBox(height: 16),
            // Mock Servers / Communities
            const _SidebarItem(imageUrl: 'https://via.placeholder.com/50', badgeCount: 5),
            const SizedBox(height: 12),
            const _SidebarItem(imageUrl: 'https://via.placeholder.com/50', badgeCount: 0),
            
            const Spacer(),
            // Add Button
            Container(
              width: 48, height: 48,
              decoration: BoxDecoration(
                color: Colors.white.withValues(alpha:0.05),
                shape: BoxShape.circle,
              ),
              child: const Icon(Icons.add, color: Colors.green),
            ),
            const SizedBox(height: 16),
          ],
        ),
      ),
    );
  }
}

class _SidebarItem extends StatelessWidget {
  final IconData? icon;
  final String? imageUrl;
  final Color? color;
  final bool isActive;
  final int badgeCount;

  const _SidebarItem({
    this.icon, 
    this.imageUrl, 
    this.color, 
    this.isActive = false, 
    this.badgeCount = 0
  });

  @override
  Widget build(BuildContext context) {
    return Stack(
      clipBehavior: Clip.none,
      children: [
        Container(
          width: 48, height: 48,
          decoration: BoxDecoration(
            color: color ?? Colors.grey.shade800,
            shape: BoxShape.circle,
            image: imageUrl != null ? DecorationImage(image: NetworkImage(imageUrl!)) : null,
            border: isActive ? Border.all(color: Colors.white, width: 2) : null,
          ),
          child: icon != null ? Icon(icon, color: Colors.white) : null,
        ),
        if (isActive)
          Positioned(
            left: -12, top: 12,
            child: Container(
              width: 4, height: 24,
              decoration: const BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.only(topRight: Radius.circular(4), bottomRight: Radius.circular(4))
              ),
            ),
          ),
        if (badgeCount > 0)
          Positioned(
            right: -4, bottom: -4,
            child: Container(
              padding: const EdgeInsets.all(4),
              decoration: const BoxDecoration(color: Colors.red, shape: BoxShape.circle),
              child: Text(
                badgeCount > 9 ? '9+' : '$badgeCount',
                style: const TextStyle(color: Colors.white, fontSize: 10, fontWeight: FontWeight.bold),
              ),
            ),
          )
      ],
    );
  }
}