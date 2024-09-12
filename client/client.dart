import 'dart:io';
import 'writer.dart';
import 'server_listener.dart';
import 'client_messages.dart';
import 'server_messages.dart';

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
    print('Te has unido al servidor como "$newUsername"., ya puedes empezar a chatear. escribe /help para obtener información sobre como usar el chat');
  }

    set_username(String username) {
    Map<String, dynamic> identifyMessage = {
      'type': 'IDENTIFY',
      'username': username
    };
    writter.sendJsonMessage(identifyMessage);
  }

  
}