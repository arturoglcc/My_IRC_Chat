import 'client.dart';
import 'server_listener.dart';
import 'dart:io';
import 'user_input_listener.dart';


void main(List<String> arguments) async {
  String host = 'localhost'; // Default host
  int port = 1234;           // Default port

  // Check if the host is provided as the first argument
  if (arguments.isNotEmpty) {
    host = arguments[0]; // First argument is the host

    // Check if the port is provided as the second argument
    if (arguments.length >= 2) {
      try {
        port = int.parse(arguments[1]); // Second argument is the port
      } catch (e) {
        print('Invalid port provided. Using default port 1234.');
      }
    }
  }

  try {
    // Connect to the specified host and port
    Socket socket = await Socket.connect(host, port);
    print('Connected to $host on port $port');

    Client client = Client(port, socket);

    // Start server and user input listeners
    ServerListener serverListener = ServerListener(socket, client);
    serverListener.startListening();

    UserInputListener userInputListener = UserInputListener(socket, client);
    userInputListener.startListening();
  } catch (e) {
    print('Failed to connect: $e');
  }
}