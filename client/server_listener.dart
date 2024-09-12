import 'dart:io';
import 'server_messages.dart';
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
        String message = String.fromCharCodes(data).trim();
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
