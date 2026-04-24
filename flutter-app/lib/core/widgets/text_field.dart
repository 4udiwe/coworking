import 'dart:async';
import 'package:flutter/material.dart';

class CustomTextField extends StatefulWidget {
  final String? label;
  final String? hint;
  final IconData? prefixIcon;
  final bool isPassword;
  final TextEditingController? controller;
  final String? Function(String?)? validator;
  final TextInputType keyboardType;

  const CustomTextField({
    super.key,
    this.label,
    this.hint,
    this.prefixIcon,
    this.isPassword = false,
    this.controller,
    this.validator,
    this.keyboardType = TextInputType.text,
  });

  @override
  State<CustomTextField> createState() => _CustomTextFieldState();
}

class _CustomTextFieldState extends State<CustomTextField> {
  late FocusNode _focusNode;
  Timer? _debounce;
  bool _obscure = true;

  @override
  void initState() {
    super.initState();
    _focusNode = FocusNode();

    // 👉 Валидация при потере фокуса
    _focusNode.addListener(() {
      if (!_focusNode.hasFocus) {
        Form.of(context).validate();
      }
    });
  }

  @override
  void dispose() {
    _focusNode.dispose();
    _debounce?.cancel();
    super.dispose();
  }

  void _onChanged(String value) {
    if (widget.validator == null) return;

    _debounce?.cancel();
    _debounce = Timer(const Duration(milliseconds: 4000), () {
      Form.of(context).validate();
    });
  }

  @override
  Widget build(BuildContext context) {
    final isPassword = widget.isPassword;

    return TextFormField(
      controller: widget.controller,
      focusNode: _focusNode,
      validator: widget.validator,
      keyboardType: widget.keyboardType,
      obscureText: isPassword ? _obscure : false,
      autovalidateMode: AutovalidateMode.onUnfocus,
      onChanged: _onChanged,

      decoration: InputDecoration(
        labelText: widget.label,
        hintText: widget.hint,

        prefixIcon: widget.prefixIcon != null
            ? Icon(widget.prefixIcon)
            : null,

        // 👁️ показать/скрыть пароль
        suffixIcon: isPassword
            ? IconButton(
          icon: Icon(
            _obscure ? Icons.visibility_off : Icons.visibility,
          ),
          onPressed: () {
            setState(() => _obscure = !_obscure);
          },
        )
            : null,

        filled: true,
        fillColor: Colors.grey.shade100,

        contentPadding: const EdgeInsets.symmetric(
          vertical: 18,
          horizontal: 20,
        ),

        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(16),
          borderSide: BorderSide.none,
        ),

        enabledBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(16),
          borderSide: BorderSide(color: Colors.grey.shade300),
        ),

        focusedBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(16),
          borderSide: const BorderSide(color: Colors.blue, width: 2),
        ),

        errorBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(16),
          borderSide: const BorderSide(color: Colors.red),
        ),

        focusedErrorBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(16),
          borderSide: const BorderSide(color: Colors.red, width: 2),
        ),
      ),
    );
  }
}