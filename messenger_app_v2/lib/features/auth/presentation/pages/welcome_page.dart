import 'package:flutter/material.dart';
import 'package:flutter_animate/flutter_animate.dart';
import 'package:go_router/go_router.dart';

import '../../../../app/config/theme/app_theme.dart';
import '../../../../app/navigation/app_routes.dart';
import '../../../../shared/presentation/l10n_extensions.dart';
import '../../../../shared/presentation/widgets/suda_button.dart';
import '../widgets/language_selector.dart';

class WelcomePage extends StatelessWidget {
  const WelcomePage({super.key});

  @override
  Widget build(BuildContext context) {
    final palette = Theme.of(context).extension<AppPaletteTheme>()!;
    final theme = Theme.of(context);

    return Scaffold(
      backgroundColor: theme.scaffoldBackgroundColor,
      body: Stack(
        children: [
          // Radial gradient background
          Positioned.fill(
            child: CustomPaint(painter: _SplashBgPainter(theme.colorScheme.primary, theme.colorScheme.tertiary)),
          ),
          SafeArea(
            child: Column(
              children: [
                // Language button top-right
                Align(
                  alignment: Alignment.topRight,
                  child: Padding(
                    padding: const EdgeInsets.all(8),
                    child: const LanguageSelector(),
                  ),
                ),

                // Logo hero
                Expanded(
                  child: Center(
                    child: _SudaLogo(
                      primaryColor: theme.colorScheme.primary,
                      accentColor: theme.colorScheme.tertiary,
                    ).animate().fadeIn(duration: 600.ms).scale(
                          begin: const Offset(0.8, 0.8),
                          end: const Offset(1.0, 1.0),
                          duration: 600.ms,
                          curve: Curves.easeOutBack,
                        ),
                  ),
                ),

                // Copy block
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 28),
                  child: Column(
                    children: [
                      // "Suda" gradient title
                      ShaderMask(
                        shaderCallback: (bounds) => LinearGradient(
                          begin: Alignment.topLeft,
                          end: Alignment.bottomRight,
                          colors: [Colors.white, palette.textAccent],
                        ).createShader(bounds),
                        child: Text(
                          'Suda',
                          style: TextStyle(
                            fontFamily: 'SpaceGrotesk',
                            fontSize: 38,
                            fontWeight: FontWeight.w800,
                            letterSpacing: -0.04,
                            color: Colors.white,
                          ),
                        ),
                      ).animate().fadeIn(delay: 200.ms).slideY(begin: 0.3, end: 0),

                      const SizedBox(height: 10),

                      Text(
                        context.l10n.authWelcomeSubtitle,
                        textAlign: TextAlign.center,
                        style: TextStyle(
                          fontFamily: 'Manrope',
                          fontSize: 15,
                          height: 1.6,
                          color: palette.textSecondary,
                        ),
                      ).animate().fadeIn(delay: 350.ms),
                    ],
                  ),
                ),

                const SizedBox(height: 32),

                // Action buttons
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20),
                  child: Column(
                    children: [
                      SudaButton(
                        label: context.l10n.authGetStarted,
                        variant: SudaButtonVariant.primary,
                        size: SudaButtonSize.lg,
                        fullWidth: true,
                        iconRight: Icons.arrow_forward_rounded,
                        onPressed: () => context.pushNamed(AppRoutes.register),
                      ),
                      const SizedBox(height: 12),
                      SudaButton(
                        label: context.l10n.authLogin,
                        variant: SudaButtonVariant.ghost,
                        size: SudaButtonSize.lg,
                        fullWidth: true,
                        onPressed: () => context.pushNamed(AppRoutes.login),
                      ),
                    ],
                  ),
                ).animate().fadeIn(delay: 500.ms).slideY(begin: 0.5, end: 0),

                const SizedBox(height: 28),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

// ── Compass-star SVG logo ──────────────────────────────────────
class _SudaLogo extends StatelessWidget {
  final Color primaryColor;
  final Color accentColor;

  const _SudaLogo({required this.primaryColor, required this.accentColor});

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      width: 180,
      height: 180,
      child: CustomPaint(
        painter: _LogoPainter(primaryColor: primaryColor, accentColor: accentColor),
      ),
    );
  }
}

class _LogoPainter extends CustomPainter {
  final Color primaryColor;
  final Color accentColor;

  const _LogoPainter({required this.primaryColor, required this.accentColor});

  @override
  void paint(Canvas canvas, Size size) {
    final center = Offset(size.width / 2, size.height / 2);
    final r = size.width / 2;

    // Outer glow circle
    final glowPaint = Paint()
      ..color = accentColor.withValues(alpha: 0.18)
      ..maskFilter = const MaskFilter.blur(BlurStyle.normal, 20);
    canvas.drawCircle(center, r * 0.96, glowPaint);

    // Orbit rings
    final ringPaint = Paint()
      ..color = accentColor.withValues(alpha: 0.18)
      ..style = PaintingStyle.stroke
      ..strokeWidth = 1;
    canvas.drawCircle(center, r * 0.68, ringPaint);
    canvas.drawCircle(center, r * 0.48, ringPaint);

    // 4-pointer diamond (compass star)
    final path = Path();
    final s = r * 0.42; // star half-size
    final m = r * 0.11; // inner notch

    path.moveTo(center.dx, center.dy - s);
    path.lineTo(center.dx + m, center.dy - m);
    path.lineTo(center.dx + s, center.dy);
    path.lineTo(center.dx + m, center.dy + m);
    path.lineTo(center.dx, center.dy + s);
    path.lineTo(center.dx - m, center.dy + m);
    path.lineTo(center.dx - s, center.dy);
    path.lineTo(center.dx - m, center.dy - m);
    path.close();

    final starPaint = Paint()
      ..shader = LinearGradient(
        begin: Alignment.topLeft,
        end: Alignment.bottomRight,
        colors: [accentColor, primaryColor],
      ).createShader(Rect.fromCircle(center: center, radius: s));
    canvas.drawPath(path, starPaint);

    // Top-right highlight shard
    final highlightPath = Path();
    highlightPath.moveTo(center.dx, center.dy - s);
    highlightPath.lineTo(center.dx + m, center.dy - m);
    highlightPath.lineTo(center.dx, center.dy);
    highlightPath.close();
    canvas.drawPath(
      highlightPath,
      Paint()..color = Colors.white.withValues(alpha: 0.55),
    );

    // Satellite dots
    final dotPaint = Paint()..color = accentColor;
    canvas.drawCircle(Offset(center.dx + r * 0.68, center.dy), 2.5, dotPaint);
    canvas.drawCircle(
      Offset(center.dx - r * 0.68, center.dy),
      2,
      Paint()..color = accentColor.withValues(alpha: 0.7),
    );
    canvas.drawCircle(
      Offset(center.dx, center.dy + r * 0.7),
      2,
      Paint()..color = accentColor.withValues(alpha: 0.6),
    );
    canvas.drawCircle(
      Offset(center.dx, center.dy - r * 0.7),
      1.8,
      Paint()..color = Colors.white.withValues(alpha: 0.5),
    );
  }

  @override
  bool shouldRepaint(covariant _LogoPainter old) =>
      old.primaryColor != primaryColor || old.accentColor != accentColor;
}

class _SplashBgPainter extends CustomPainter {
  final Color primaryColor;
  final Color accentColor;

  const _SplashBgPainter(this.primaryColor, this.accentColor);

  @override
  void paint(Canvas canvas, Size size) {
    // Top-center purple glow
    final p1 = Paint()
      ..shader = RadialGradient(
        colors: [
          primaryColor.withValues(alpha: 0.7),
          primaryColor.withValues(alpha: 0.3),
          Colors.transparent,
        ],
        stops: const [0.0, 0.5, 1.0],
      ).createShader(Rect.fromCircle(
        center: Offset(size.width * 0.5, size.height * 0.4),
        radius: size.width * 0.65,
      ));
    canvas.drawRect(Rect.fromLTWH(0, 0, size.width, size.height), p1);

    // Bottom-right cyan glow
    final p2 = Paint()
      ..shader = RadialGradient(
        colors: [
          accentColor.withValues(alpha: 0.5),
          Colors.transparent,
        ],
      ).createShader(Rect.fromCircle(
        center: Offset(size.width * 0.8, size.height * 0.85),
        radius: size.width * 0.42,
      ));
    canvas.drawRect(Rect.fromLTWH(0, 0, size.width, size.height), p2);
  }

  @override
  bool shouldRepaint(covariant _SplashBgPainter old) => false;
}
