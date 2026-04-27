import 'package:flutter/material.dart';

class AppToast extends StatelessWidget {
  final String title;
  final String body;
  final VoidCallback? onTap;
  final VoidCallback? onMarkRead;

  const AppToast({
    super.key,
    required this.title,
    required this.body,
    this.onTap,
    this.onMarkRead,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Material(
      color: Colors.transparent,
      child: Container(
        margin: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
        padding: const EdgeInsets.only(top: 12, right: 12, left: 12, bottom: 4),
        decoration: BoxDecoration(
          //border: BoxBorder.all(color: theme.colorScheme.primary, width: 2),
          color: theme.colorScheme.onPrimaryContainer.withAlpha(235),
          borderRadius: BorderRadius.circular(12),
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            /// 🔔 Контент
            GestureDetector(
              onTap: onTap,
              child: Row(
                children: [
                  Icon(Icons.notifications, color: theme.colorScheme.onPrimary),
                  const SizedBox(width: 12),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          title,
                          style: TextStyle(
                            color: theme.colorScheme.onPrimary,
                            fontWeight: FontWeight.bold,
                            fontSize: 16,
                          ),
                        ),
                        const SizedBox(height: 4),
                        Text(
                          body,
                          style: TextStyle(
                            color: theme.colorScheme.onSecondary,
                          ),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            ),

            // const SizedBox(height: 2),

            /// ✅ Кнопка
            Align(
              alignment: Alignment.centerRight,
              child: TextButton(
                onPressed: onMarkRead,
                child: Text(
                  'Прочитано',
                  style: TextStyle(color: theme.colorScheme.primaryFixed),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
