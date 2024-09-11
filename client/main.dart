import 'dart:io';
import 'client.dart';

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

  // Initialize the client
  Client client = Client(port);
  await client.connect();
    
  // Wait for user input
  while (true) {
    String? input = stdin.readLineSync();
    if (input != null && input.isNotEmpty) {
      client.writer.sendMessage(input);
    }
  }
}
