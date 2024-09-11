import 'dart:io';
import 'dart:convert'; 

class Writer {
  final Socket socket;

  Writer(this.socket);

  void sendJsonMessage(Map<String, dynamic> message) {
    String jsonString = jsonEncode(message);
    sendMessage(jsonString); 
  }

  void sendMessage(String message) {
    socket.write(message);
  }
}
