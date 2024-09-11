import 'dart:io';
import 'writter.dart';
import 'listener.dart';
import 'client_messages.dart';

class Client {
  final int port;
  late Socket socket;
  late Writer writer;
  late Listener listener;
  String? username;
  bool isFirstMessage = true; // Flag to track if it's the first message

  Client(this.port);

  Future<void> connect() async {
    try {
      socket = await Socket.connect('localhost', port);
      print('Connected to server on port $port');

      writer = Writer(socket);
      listener = Listener(socket);

// Start receiving user input and sending to server
      _startUserInputListener();
    } catch (e) {
      print('Error: $e');
      exit(1);
    }
  }

  // This method will handle user input
  void _startUserInputListener() {
    print('Bienvenido al servidor! Escribe tu nombre de usuario: ');

    // Start listening for user input from the keyboard
    stdin.listen((List<int> data) {
      String message = String.fromCharCodes(data).trim();

    // Check if this is the first message
    if (isFirstMessage) {
      clientMessages.set_username(message);
      isFirstMessage = false; 
    }

      // Send the message to the ClientMessages object for processing
      clientMessages.processMessage(message);
    });
  }

    // Function to update the username after a successful identification
  void updateUsername(String newUsername) {
    username = newUsername;
    print('Te has unido al servidor como "$newUsername"., ya puedes empezar a chatear. escribe /help para obtener información sobre como usar el chat');
  }
}