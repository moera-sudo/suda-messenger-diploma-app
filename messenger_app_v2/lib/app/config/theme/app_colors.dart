import 'package:flutter/material.dart';

class AppPalette {
  final Color primary;
  final Color secondary;
  final Color accent;
  final Color background;
  final Color surface;
  final Color surface2;
  final Color textPrimary;
  final Color textAccent;
  final Color surfaceContainer;
  final Color textSecondary;
  final Color divider;
  final Color messageMeBg;
  final Color messageOtherBg;
  final Color danger;
  final Color success;
  final Color glow;

  const AppPalette({
    required this.primary,
    required this.secondary,
    required this.accent,
    required this.background,
    required this.surface,
    required this.surface2,
    required this.textPrimary,
    required this.textAccent,
    required this.surfaceContainer,
    required this.textSecondary,
    required this.divider,
    required this.messageMeBg,
    required this.messageOtherBg,
    required this.danger,
    required this.success,
    required this.glow,
  });
}

// ─── Suda Dark (default) ────────────────────────────────────
const sudaPalette = AppPalette(
  primary: Color(0xFF644599),
  secondary: Color(0xFF495599),
  accent: Color(0xFF29DBF2),
  background: Color(0xFF0F1020),
  surface: Color(0xFF1C1F3A),
  surface2: Color(0xFF262A4D),
  textPrimary: Color(0xFFE0E0E0),
  textAccent: Color(0xFF80DEEA),
  surfaceContainer: Color(0xFF11152D),
  textSecondary: Color(0xFF8A8FB5),
  divider: Color(0x0FFFFFFF),
  messageMeBg: Color(0xFF644599),
  messageOtherBg: Color(0xFF1C1F3A),
  danger: Color(0xFFFF6B7A),
  success: Color(0xFF3DDC97),
  glow: Color(0x5929DBF2),
);

// ─── Suda Enlightened (light) ───────────────────────────────
const sudaEnlightenedPalette = AppPalette(
  primary: Color(0xFF9575CD),
  secondary: Color(0xFF8E99F3),
  accent: Color(0xFF4DD0E1),
  background: Color(0xFFF6F7FB),
  surface: Color(0xFFF0E7F6),
  surface2: Color(0xFFE6DCEE),
  textPrimary: Color(0xFF1A1A1A),
  textAccent: Color(0xFF603CAB),
  surfaceContainer: Color(0xFFE8DFF0),
  textSecondary: Color(0xFF6A6A78),
  divider: Color(0x0F000000),
  messageMeBg: Color(0xFF9575CD),
  messageOtherBg: Color(0xFFFFFFFF),
  danger: Color(0xFFD9485B),
  success: Color(0xFF1FA567),
  glow: Color(0x594DD0E1),
);

// ─── TeaChats Light ────────────────────────────────────────
const teaChatsLightPalette = AppPalette(
  primary: Color(0xFF809F73),
  secondary: Color(0xFFA5C895),
  accent: Color(0xFF66BB6A),
  background: Color(0xFFF6F8F6),
  surface: Color(0xFFE5ECE2),
  surface2: Color(0xFFD8E1D4),
  textPrimary: Color(0xFF1E1E1E),
  textAccent: Color(0xFF3F6623),
  surfaceContainer: Color(0xFFD9E4D5),
  textSecondary: Color(0xFF5C6A56),
  divider: Color(0x0F000000),
  messageMeBg: Color(0xFF809F73),
  messageOtherBg: Color(0xFFFFFFFF),
  danger: Color(0xFFC7453B),
  success: Color(0xFF3F6623),
  glow: Color(0x5966BB6A),
);

// ─── TeaChats Dark ─────────────────────────────────────────
const teaChatsDarkPalette = AppPalette(
  primary: Color(0xFF607C5B),
  secondary: Color(0xFF4B5D43),
  accent: Color(0xFF81C784),
  background: Color(0xFF111512),
  surface: Color(0xFF1C221E),
  surface2: Color(0xFF252C26),
  textPrimary: Color(0xFFE0E0E0),
  textAccent: Color(0xFFA5D6A7),
  surfaceContainer: Color(0xFF151A17),
  textSecondary: Color(0xFF8E9A8B),
  divider: Color(0x0FFFFFFF),
  messageMeBg: Color(0xFF607C5B),
  messageOtherBg: Color(0xFF1C221E),
  danger: Color(0xFFFF6B7A),
  success: Color(0xFFA5D6A7),
  glow: Color(0x5981C784),
);

// ─── Ethereal Sky (light) ──────────────────────────────────
const etherealSkyPalette = AppPalette(
  primary: Color(0xFF5EC8D8),
  secondary: Color(0xFF90E0EF),
  accent: Color(0xFF00BCD4),
  background: Color(0xFFF4FBFD),
  surface: Color(0xFFE1F5FE),
  surface2: Color(0xFFD2EDF8),
  textPrimary: Color(0xFF1A1A1A),
  textAccent: Color(0xFF0274A8),
  surfaceContainer: Color(0xFFD4EEF7),
  textSecondary: Color(0xFF5A6B72),
  divider: Color(0x0F000000),
  messageMeBg: Color(0xFF5EC8D8),
  messageOtherBg: Color(0xFFFFFFFF),
  danger: Color(0xFFD9485B),
  success: Color(0xFF1FA567),
  glow: Color(0x5900BCD4),
);

// ─── Ethereal Abyss (dark) ─────────────────────────────────
const etherealAbyssPalette = AppPalette(
  primary: Color(0xFF3A9EB1),
  secondary: Color(0xFF2E7A8A),
  accent: Color(0xFF4DD0E1),
  background: Color(0xFF0C1114),
  surface: Color(0xFF152126),
  surface2: Color(0xFF1E2D33),
  textPrimary: Color(0xFFE0E0E0),
  textAccent: Color(0xFF80DEEA),
  surfaceContainer: Color(0xFF0F1A1E),
  textSecondary: Color(0xFF8AA0A8),
  divider: Color(0x0FFFFFFF),
  messageMeBg: Color(0xFF3A9EB1),
  messageOtherBg: Color(0xFF152126),
  danger: Color(0xFFFF6B7A),
  success: Color(0xFF3DDC97),
  glow: Color(0x594DD0E1),
);
