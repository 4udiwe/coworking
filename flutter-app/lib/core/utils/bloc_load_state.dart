enum LoadStatus {
  initial,
  loading,
  success,
  error,
}

class LoadState<T> {
  final T? data;
  final LoadStatus status;
  final String? error;

  const LoadState({
    this.data,
    this.status = LoadStatus.initial,
    this.error,
  });

  LoadState<T> copyWith({
    T? data,
    LoadStatus? status,
    String? error,
  }) {
    return LoadState<T>(
      data: data ?? this.data,
      status: status ?? this.status,
      error: error,
    );
  }
}