import 'package:coworking_app/core/widgets/text_field.dart';
import 'package:coworking_app/core/l10n/localization_helper.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../../bloc/auth_bloc.dart';
import '../../bloc/auth_event.dart';
import '../../bloc/auth_state.dart';
import '../../../../core/di/service_locator.dart';

class RegisterScreen extends StatefulWidget {
  const RegisterScreen({super.key});

  @override
  State<RegisterScreen> createState() => _RegisterScreenState();
}

class _RegisterScreenState extends State<RegisterScreen> {
  final nameController = TextEditingController();
  final lastNameController = TextEditingController();
  final emailController = TextEditingController();
  final passwordController = TextEditingController();
  final confirmPasswordController = TextEditingController();
  final adminSecretController = TextEditingController();
  final _formKey = GlobalKey<FormState>();

  // Role constant as per requirements (usually students register themselves)
  final String _roleCode = 'student';

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (_) => sl<AuthBloc>(),
      child: Scaffold(
        body: BlocConsumer<AuthBloc, AuthState>(
          listener: (context, state) {
            if (state.isAuthenticated) {
              ScaffoldMessenger.of(context).showSnackBar(
                SnackBar(content: Text(context.l10n.registerSuccess)),
              );
              Navigator.pushReplacementNamed(context, '/main');
            } else if (state.status == AuthStatus.failure) {
              ScaffoldMessenger.of(
                context,
              ).showSnackBar(SnackBar(content: Text(state.error!)));
            }
          },
          builder: (context, state) {
            if (state.status == AuthStatus.loading) {
              return const Center(child: CircularProgressIndicator());
            }

            return Padding(
              padding: const EdgeInsets.only(left: 40, right: 40),
              child: Center(
                child: SizedBox(
                  width: 350.0,
                  child: SingleChildScrollView(
                    child: Form(
                      key: _formKey,
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        crossAxisAlignment: CrossAxisAlignment.center,
                        children: [
                          const SizedBox(height: 50),
                          Text(
                            context.l10n.registerTitle,
                            textAlign: TextAlign.center,
                            style: Theme.of(context).textTheme.headlineMedium,
                          ),
                          const SizedBox(height: 20),
                          CustomTextField(
                            label: context.l10n.firstNameFieldLabel,
                            hint: context.l10n.firstNameFieldHint,
                            prefixIcon: Icons.person,
                            controller: nameController,
                            validator: (value) {
                              if (value == null || value.isEmpty) {
                                return context.l10n.validationEmptyFirstName;
                              }
                              return null;
                            },
                          ),
                          const SizedBox(height: 16),
                          CustomTextField(
                            label: context.l10n.lastNameFieldLabel,
                            hint: context.l10n.lastNameFieldHint,
                            prefixIcon: Icons.person,
                            controller: lastNameController,
                            validator: (value) {
                              if (value == null || value.isEmpty) {
                                return context.l10n.validationEmptyLastName;
                              }
                              return null;
                            },
                          ),
                          const SizedBox(height: 16),
                          CustomTextField(
                            label: context.l10n.emailFieldLabel,
                            hint: context.l10n.emailFieldHint,
                            prefixIcon: Icons.email,
                            controller: emailController,
                            keyboardType: TextInputType.emailAddress,
                            validator: (value) {
                              if (value == null || value.isEmpty) {
                                return context.l10n.validationEmptyEmail;
                              }
                              if (!value.contains('@')) {
                                return context.l10n.validationInvalidEmail;
                              }
                              return null;
                            },
                          ),
                          const SizedBox(height: 16),
                          CustomTextField(
                            label: context.l10n.passwordFieldLabel,
                            hint: context.l10n.passwordFieldHint,
                            prefixIcon: Icons.lock,
                            isPassword: true,
                            controller: passwordController,
                            validator: (value) {
                              if (value == null || value.length < 8) {
                                return context.l10n.validationMinimumPassword;
                              }
                              return null;
                            },
                          ),
                          const SizedBox(height: 16),
                          CustomTextField(
                            label: context.l10n.confirmPasswordFieldLabel,
                            hint: context.l10n.confirmPasswordFieldHint,
                            isPassword: true,
                            controller: confirmPasswordController,
                            validator: (value) {
                              if (value != passwordController.text) {
                                return context.l10n.validationPasswordsNotMatch;
                              }
                              return null;
                            },
                          ),
                          const SizedBox(height: 20),
                          ElevatedButton(
                            onPressed: () {
                              if (_formKey.currentState!.validate()) {
                                context.read<AuthBloc>().add(
                                  AuthRegister(
                                    name: nameController.text.trim(),
                                    lastName: lastNameController.text.trim(),
                                    email: emailController.text.trim(),
                                    password: passwordController.text.trim(),
                                    role: _roleCode,
                                  ),
                                );
                              } else {
                                ScaffoldMessenger.of(context).showSnackBar(
                                  SnackBar(
                                    content: Text(
                                      context.l10n.validationInvalidData,
                                    ),
                                  ),
                                );
                              }
                            },
                            child: Text(context.l10n.registerButtonLabel),
                          ),
                          const SizedBox(height: 16),
                          Row(
                            mainAxisAlignment: MainAxisAlignment.center,
                            children: [
                              Text(context.l10n.alreadyHaveAccount),
                              TextButton(
                                onPressed: () {
                                  Navigator.pop(context);
                                },
                                child: Text(context.l10n.loginButtonLabel),
                              ),
                            ],
                          ),
                          const SizedBox(height: 50),
                        ],
                      ),
                    ),
                  ),
                ),
              ),
            );
          },
        ),
      ),
    );
  }
}
