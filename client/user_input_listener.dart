import 'dart:io';
import 'client.dart';

class UserInputListener {
  final Socket socket;
  bool isFirstMessage = true;
  final Client client;

  UserInputListener(this.socket, this.client);

  void startListening() {
    print('Bienvenido al servidor, escribe tu nombre de usuario: ');
    stdin.listen((data) {
      String message = String.fromCharCodes(data).trim();

      if (message.isNotEmpty) {
        if (isFirstMessage == true) {
          client.set_username(message);
        }
        try {
          socket.write(message);
        } catch (e) {
          print('Failed to send message: $e');
        }
      }
    });
  }
}
