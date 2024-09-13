import 'dart:io';
import 'writer.dart';
import 'client_messages.dart';

class Client {
  final int port;
  final Socket socket;
  late ClientMessages clientMessages;
  late Writer writter;
  String? username;

  Client(this.port, this.socket) {
    this.writter = Writer(socket);
  }
      
    // Function to update the username after a successful identification
  void updateUsername(String newUsername) {
    username = newUsername;
    print('Te has unido al servidor como "$newUsername". Ya puedes empezar a chatear.\nEscribe /help para obtener información sobre como usar el chat');
  }

    set_username(String username) {
    Map<String, dynamic> identifyMessage = {
      'type': 'IDENTIFY',
      'username': username
    };
    writter.sendJsonMessage(identifyMessage);
  }

  void disconnect() {
    print('Desconectando del servidor...');
    socket.close();
  }
}