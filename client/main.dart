import 'client.dart';
import 'server_listener.dart';
import 'dart:io';
import 'user_input_listener.dart';


void main(List<String> arguments) async {
  int port = 1234; // Default port
  // Check if a port is provided, and attempt to parse it
  if (arguments.isNotEmpty) {
    try {
      port = int.parse(arguments[0]);
    } catch (e) {
      print('Invalid port provided. Using default port 1234.');
    }
  }

  Socket socket = await Socket.connect('localhost', port);

  Client client = Client(port, socket);

  ServerListener serverListener = ServerListener(socket, client);
  serverListener.startListening(); 
    
  UserInputListener userInputListener = UserInputListener(socket, client);
  userInputListener.startListening(); 
}
