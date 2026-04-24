import 'package:coworking_app/core/utils/bloc_load_state.dart';
import '../../../core/models/user.dart';
import '../../../core/models/user_session.dart';

class UserState {
  final LoadState<User> profile;
  final LoadState<List<UserSession>> activeSessions;
  final LoadState<List<UserSession>> allSessions;

  final String? messageId;
  final String? actionMessage;
  final bool isError;


  const UserState({
    this.profile = const LoadState(),
    this.activeSessions = const LoadState(),
    this.allSessions = const LoadState(),
    this.messageId,
    this.actionMessage,
    this.isError = false,
  });

  UserState copyWith({
    LoadState<User>? profile,
    LoadState<List<UserSession>>? activeSessions,
    LoadState<List<UserSession>>? allSessions,
    String? messageId,
    String? actionMessage,
    bool? isError,
  }) {
    return UserState(
      profile: profile ?? this.profile,
      activeSessions: activeSessions ?? this.activeSessions,
      allSessions: allSessions ?? this.allSessions,
      messageId: messageId ?? this.messageId,
      actionMessage: actionMessage ?? this.actionMessage,
    );
  }
}