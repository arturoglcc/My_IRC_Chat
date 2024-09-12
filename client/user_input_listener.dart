import 'dart:io';
import 'client.dart';
import 'client_messages.dart';
import 'writer.dart';

class UserInputListener {
  final Socket socket;
  bool isFirstMessage = true;
  final Client client;
  late Writer writer;
  late ClientMessages clientMessages;

  UserInputListener(this.socket, this.client) {
    this.writer = Writer(socket);
    this.clientMessages = ClientMessages(writer);
  }

  void startListening() {
    print('Bienvenido al servidor, escribe tu nombre de usuario: ');
    stdin.listen((data) {
      String message = String.fromCharCodes(data).trim();
      if (message.isNotEmpty) {
        if (isFirstMessage == true) {
          client.set_username(message);
          isFirstMessage = false;
        } else {
          try {
            clientMessages.processMessage(message);
        } catch (e) {
          print('Failed to send message: $e');
          }
        }
      }
    });
  }
}
