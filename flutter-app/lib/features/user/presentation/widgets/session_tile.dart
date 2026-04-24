import 'package:coworking_app/features/user/bloc/user_event.dart';
import 'package:coworking_app/core/l10n/localization_helper.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:intl/intl.dart';

import '../../../../core/models/user_session.dart';
import '../../bloc/user_bloc.dart';

final dateTimeFormat = DateFormat('dd MMM yyyy, HH:mm', 'RU');
final dateFormat = DateFormat('dd MMM yyyy', 'RU');

class SessionTile extends StatefulWidget {
  final UserSession session;

  const SessionTile({super.key, required this.session});

  @override
  State<SessionTile> createState() => _SessionTileState();
}

class _SessionTileState extends State<SessionTile> {
  bool expanded = false;

  Color _statusColor() {
    if (widget.session.revoked) return Colors.grey;
    if (widget.session.current) return Colors.blue;
    return Colors.green;
  }

  @override
  Widget build(BuildContext context) {
    final s = widget.session;

    return Padding(
      padding: EdgeInsets.symmetric(horizontal: 12, vertical: 6),
      child: Material(
        elevation: 4.0,
        surfaceTintColor: widget.session.current ? Colors.blue : null,
        borderRadius: BorderRadius.circular(24),
        child: InkWell(
          borderRadius: BorderRadius.circular(24),
          onTap: () => setState(() => expanded = !expanded),
          child: AnimatedContainer(
            duration: const Duration(milliseconds: 200),
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                /// HEADER
                Row(
                  mainAxisAlignment: MainAxisAlignment.end,
                  children: [
                    /// DEVICE
                    Expanded(
                      flex: 4,
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            s.device,
                            style: const TextStyle(fontWeight: FontWeight.w600),
                            maxLines: 1,
                            overflow: TextOverflow.ellipsis,
                          ),
                          const SizedBox(height: 4),
                          Text(
                            s.ipAddress,
                            style: TextStyle(color: Colors.grey[600]),
                          ),
                        ],
                      ),
                    ),

                    Expanded(
                      flex: 3,
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.end,
                        children: [
                          Container(
                            padding: const EdgeInsets.symmetric(
                              horizontal: 10,
                              vertical: 6,
                            ),
                            decoration: BoxDecoration(
                              color: _statusColor().withOpacity(0.15),
                              borderRadius: BorderRadius.circular(12),
                            ),
                            child: Text(
                              s.current
                                  ? context.l10n.currentSession
                                  : s.revoked
                                  ? context.l10n.revokedSession
                                  : context.l10n.activeSession,
                              style: TextStyle(
                                color: _statusColor(),
                                fontWeight: FontWeight.w500,
                              ),
                            ),
                          ),
                          const SizedBox(width: 8),
                          Icon(
                            expanded ? Icons.expand_less : Icons.expand_more,
                          ),
                        ],
                      ),
                    ),

                    /// STATUS
                  ],
                ),

                /// EXPANDED
                if (expanded) ...[
                  const SizedBox(height: 16),
                  const Divider(),

                  const SizedBox(height: 12),

                  Text('User Agent: ${s.userAgent}'),

                  const SizedBox(height: 4),

                  Text(
                    context.l10n.sessionCreated +
                        dateTimeFormat.format(s.createdAt.toLocal()).toString(),
                  ),

                  const SizedBox(height: 4),

                  Text(
                    context.l10n.sessionLastUsedAt +
                        dateTimeFormat
                            .format(s.lastUsedAt.toLocal())
                            .toString(),
                  ),

                  if (!s.revoked) ...[
                    const SizedBox(height: 4),

                    Text(
                      context.l10n.sessionExpiresAt +
                          dateFormat.format(s.expiresAt.toLocal()).toString(),
                    ),
                  ],

                  if (!s.revoked && !s.current) ...[
                    const SizedBox(height: 16),

                    SizedBox(
                      width: double.infinity,
                      child: ElevatedButton(
                        onPressed: () {
                          context.read<UserBloc>().add(
                            RevokeSessionEvent(s.id),
                          );
                          setState(() => expanded = false);
                        },
                        child: Text(context.l10n.revokeSession),
                      ),
                    ),
                  ],
                ],
              ],
            ),
          ),
        ),
      ),
    );
  }
}
