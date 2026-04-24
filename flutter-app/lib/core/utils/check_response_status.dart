void checkStatus(dynamic response, {List<int> validCodes = const [200]}) {
  if (!validCodes.contains(response.statusCode)) {
    throw Exception(
      'Request failed: ${response.statusCode}, body: ${response.body}',
    );
  }
}