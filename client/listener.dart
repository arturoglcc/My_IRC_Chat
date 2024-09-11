import 'dart:io';
import 'client.dart';
import 'server_messages.dart';

class Listener {
  final Socket socket;
  late Client client;

  Listener(this.socket, this.client);

  void startListening() {
    ServerMessages serverMessages = ServerMessages(client);
    socket.listen((List<int> data) {
      String message = String.fromCharCodes(data).trim();
      serverMessages.processMessage(message); 
    });
  }
}
