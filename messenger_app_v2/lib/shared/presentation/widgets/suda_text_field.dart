import 'package:flutter/material.dart';

import '../../../app/config/theme/app_theme.dart';

class SudaTextField extends StatefulWidget {
  final String label;
  final String? hint;
  final TextEditingController? controller;
  final IconData? leftIcon;
  final IconData? rightIcon;
  final VoidCallback? onRightTap;
  final TextInputType? keyboardType;
  final bool obscureText;
  final String? error;
  final String? hint2;
  final ValueChanged<String>? onChanged;
  final TextInputAction? textInputAction;
  final bool enabled;
  final int? maxLines;

  const SudaTextField({
    super.key,
    required this.label,
    this.hint,
    this.controller,
    this.leftIcon,
    this.rightIcon,
    this.onRightTap,
    this.keyboardType,
    this.obscureText = false,
    this.error,
    this.hint2,
    this.onChanged,
    this.textInputAction,
    this.enabled = true,
    this.maxLines = 1,
  });

  @override
  State<SudaTextField> createState() => _SudaTextFieldState();
}

class _SudaTextFieldState extends State<SudaTextField> {
  final FocusNode _focus = FocusNode();
  bool _focused = false;

  @override
  void initState() {
    super.initState();
    _focus.addListener(() => setState(() => _focused = _focus.hasFocus));
  }

  @override
  void dispose() {
    _focus.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final palette = Theme.of(context).extension<AppPaletteTheme>()!;
    final theme = Theme.of(context);
    final hasError = widget.error != null && widget.error!.isNotEmpty;

    final borderColor = hasError
        ? palette.danger
        : _focused
            ? theme.colorScheme.tertiary
            : palette.divider;

    final shadowColor = hasError
        ? palette.danger.withValues(alpha: 0.12)
        : _focused
            ? palette.glow
            : Colors.transparent;

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        // UPPERCASE label
        Text(
          widget.label.toUpperCase(),
          style: TextStyle(
            fontFamily: 'Manrope',
            fontSize: 11,
            fontWeight: FontWeight.w700,
            color: palette.textAccent,
            letterSpacing: 0.12,
          ),
        ),
        const SizedBox(height: 8),
        // Input wrap
        AnimatedContainer(
          duration: const Duration(milliseconds: 150),
          decoration: BoxDecoration(
            color: theme.colorScheme.surface,
            borderRadius: BorderRadius.circular(14),
            border: Border.all(color: borderColor, width: hasError ? 1.5 : 1),
            boxShadow: [
              BoxShadow(
                color: shadowColor,
                blurRadius: 0,
                spreadRadius: 4,
              ),
            ],
          ),
          constraints: const BoxConstraints(minHeight: 52),
          child: Row(
            children: [
              if (widget.leftIcon != null)
                Padding(
                  padding: const EdgeInsets.only(left: 14),
                  child: Icon(widget.leftIcon, size: 18, color: palette.textSecondary),
                ),
              Expanded(
                child: TextField(
                  focusNode: _focus,
                  controller: widget.controller,
                  keyboardType: widget.keyboardType,
                  obscureText: widget.obscureText,
                  onChanged: widget.onChanged,
                  textInputAction: widget.textInputAction,
                  enabled: widget.enabled,
                  maxLines: widget.obscureText ? 1 : widget.maxLines,
                  style: TextStyle(
                    fontFamily: 'Manrope',
                    fontSize: 15,
                    color: theme.colorScheme.onSurface,
                  ),
                  decoration: InputDecoration(
                    hintText: widget.hint,
                    hintStyle: TextStyle(
                      color: palette.textSecondary,
                      fontSize: 15,
                    ),
                    border: InputBorder.none,
                    enabledBorder: InputBorder.none,
                    focusedBorder: InputBorder.none,
                    contentPadding: EdgeInsets.symmetric(
                      horizontal: widget.leftIcon != null ? 10 : 14,
                      vertical: 14,
                    ),
                    isDense: true,
                  ),
                ),
              ),
              if (widget.rightIcon != null)
                GestureDetector(
                  onTap: widget.onRightTap,
                  child: Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 12),
                    child: Icon(widget.rightIcon, size: 20, color: palette.textSecondary),
                  ),
                ),
            ],
          ),
        ),
        if (widget.hint2 != null || hasError)
          Padding(
            padding: const EdgeInsets.only(top: 6, left: 4),
            child: Text(
              hasError ? widget.error! : widget.hint2!,
              style: TextStyle(
                fontSize: 12,
                color: hasError ? palette.danger : palette.textSecondary,
              ),
            ),
          ),
      ],
    );
  }
}
