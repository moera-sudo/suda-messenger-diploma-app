import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import '../../../app/config/theme/app_theme.dart';

class SudaOtpCell extends StatefulWidget {
  final TextEditingController controller;
  final FocusNode focusNode;
  final bool filled;
  final ValueChanged<String> onChanged;
  final VoidCallback onBackspace;

  const SudaOtpCell({
    super.key,
    required this.controller,
    required this.focusNode,
    required this.filled,
    required this.onChanged,
    required this.onBackspace,
  });

  @override
  State<SudaOtpCell> createState() => _SudaOtpCellState();
}

class _SudaOtpCellState extends State<SudaOtpCell> {
  // KeyboardListener needs a stable FocusNode — create it once, not on every build.
  final FocusNode _keyboardListenerNode = FocusNode();
  bool _isFocused = false;

  @override
  void initState() {
    super.initState();
    widget.focusNode.addListener(_onFocusChange);
  }

  @override
  void didUpdateWidget(SudaOtpCell old) {
    super.didUpdateWidget(old);
    if (old.focusNode != widget.focusNode) {
      old.focusNode.removeListener(_onFocusChange);
      widget.focusNode.addListener(_onFocusChange);
    }
  }

  void _onFocusChange() {
    setState(() => _isFocused = widget.focusNode.hasFocus);

    // When the cell gains focus, always move cursor to end.
    // This prevents right-to-left typing that occurs when a tap on a centred
    // character places the cursor at offset 0 instead of offset 1.
    if (widget.focusNode.hasFocus && widget.controller.text.isNotEmpty) {
      // Select the existing digit so typing a new one replaces it.
      // Setting cursor to end would block input (maxLength=1 already full).
      widget.controller.selection = TextSelection(
        baseOffset: 0,
        extentOffset: widget.controller.text.length,
      );
    }
  }

  @override
  void dispose() {
    widget.focusNode.removeListener(_onFocusChange);
    _keyboardListenerNode.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final palette = Theme.of(context).extension<AppPaletteTheme>()!;
    final theme = Theme.of(context);

    return AnimatedContainer(
      duration: const Duration(milliseconds: 150),
      width: 44,
      height: 52,
      decoration: BoxDecoration(
        color: theme.colorScheme.surface,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(
          color: widget.filled || _isFocused
              ? theme.colorScheme.tertiary
              : palette.divider,
          width: 1.5,
        ),
        boxShadow: widget.filled || _isFocused
            ? [BoxShadow(color: palette.glow, blurRadius: 0, spreadRadius: 3)]
            : null,
      ),
      child: KeyboardListener(
        focusNode: _keyboardListenerNode,
        onKeyEvent: (event) {
          if (event is KeyDownEvent &&
              event.logicalKey == LogicalKeyboardKey.backspace &&
              widget.controller.text.isEmpty) {
            widget.onBackspace();
          }
        },
        child: TextField(
          controller: widget.controller,
          focusNode: widget.focusNode,
          textAlign: TextAlign.center,
          keyboardType: TextInputType.number,
          maxLength: 1,
          style: TextStyle(
            fontFamily: 'JetBrainsMono',
            fontSize: 22,
            fontWeight: FontWeight.w700,
            color: theme.colorScheme.onSurface,
          ),
          decoration: const InputDecoration(
            counterText: '',
            border: InputBorder.none,
            enabledBorder: InputBorder.none,
            focusedBorder: InputBorder.none,
            isDense: true,
            contentPadding: EdgeInsets.zero,
          ),
          onChanged: widget.onChanged,
        ),
      ),
    );
  }
}
