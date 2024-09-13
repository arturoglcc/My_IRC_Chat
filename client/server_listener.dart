import 'dart:io';
import 'server_messages.dart';
import 'dart:convert';
import 'client.dart';

class ServerListener {
  final Socket socket;
  final Client client;
  late ServerMessages serverMessages;

  ServerListener(this.socket,  this.client) {
     this.serverMessages = ServerMessages(client);
  }

 void startListening() {
    socket.listen(
      (data) {
        // Decode the received data using UTF-8 to handle special characters
        String message = utf8.decode(data).trim();
        print('Message from server: $message');
        serverMessages.processMessage(message); 
      },
      onError: (error) {
        print('Server error: $error');
        socket.destroy();
      },
      onDone: () {
        print('Server closed the connection');
        socket.destroy();
      },
    );
  }
}
