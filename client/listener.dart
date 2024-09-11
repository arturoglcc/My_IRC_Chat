import 'dart:io';
import 'server_messages.dart';

class Listener {
  final Socket socket;

  Listener(this.socket);

  void startListening() {
    socket.listen((List<int> data) {
      String message = String.fromCharCodes(data).trim();
      ServerMessages.processMessage(message);
    });
  }
}
