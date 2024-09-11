import 'dart:convert'; 
import 'client_messages.dart'; 

class ServerMessages {
  final ClientMessages clientMessages; 

  ServerMessages(this.clientMessages);

  void processMessage(String message) {
    Map<String, dynamic> jsonMessage = jsonDecode(message);

    // Check the "type" of the message
    switch (jsonMessage['type']) {
      case 'RESPONSE':
        _handleResponse(jsonMessage);
        break;

      case 'NEW_USER':
        _handleNewUser(jsonMessage); 
        break; 

      case 'NEW_STATUS':
        _handleNewStatus(jsonMessage);
        break;

       case 'USER_LIST':
        _handleUserList(jsonMessage);
        break; 

      case 'TEXT_FROM':
        _handleTextFrom(jsonMessage);
        break;  

      case 'PUBLIC_TEXT_FROM':
        _handlePublicTextFrom(jsonMessage);
        break;

      case 'INVITATION':
        _handleInvitation(jsonMessage);
        break;

      case 'JOINED_ROOM':
        _handleJoinedRoom(jsonMessage);
        break; 

      case 'ROOM_USER_LIST':
        _handleRoomUserList(jsonMessage);
        break;        

      case 'ROOM_TEXT_FROM':
        _handleRoomTextFrom(jsonMessage);
        break;  

      case 'LEFT_ROOM':
        _handleLeftRoom(jsonMessage);
        break;  

      case 'DISCONNECTED':
        _handleDisconnected(jsonMessage);
        break;  

      default:
        print("Unknown message type: ${jsonMessage['type']}");
    }
  }

  void _handleResponse(Map<String, dynamic> message) {
    String operation = message['operation'];
    String result = message['result'];

    // Handle specific operations and results
    switch (operation) {
      case 'IDENTIFY':
        if (result == 'USER_ALREADY_EXISTS') {
          String existingUsername = message['extra'];
          print('nombre de usuario "$existingUsername" ya existe. \n');
          _promptForValidUsername();
        } else if (result == 'SUCCESS') {
          String newUsername = message['extra'];
          client.updateUsername(newUsername); 
        }
        break;

      case 'NEW_USER':
        _handleNewUser(jsonMessage);
        break;

      case 'TEXT':
        if (result == 'NO_SUCH_USER') {
          _handleNoSuchUser(message);
        }
        break; 

      case 'NEW_ROOM':
        if (result == 'SUCCESS') {
          _handleNewRoomSuccess(message);
        } else if (result == 'ROOM_ALREADY_EXISTS') {
          _handleRoomAlreadyExists(message);
        }
        break;

      case 'INVITE':
        if (result == 'NO_SUCH_ROOM') {
          _handleNoSuchRoom(message);
        } else if (result == 'NO_SUCH_USER') {
          _handleNoSuchUser(message);
        }
        break;

      case 'JOIN_ROOM':
        if (result == 'SUCCESS') {
          _handleJoinRoomSuccess(message);
        } else if (result == 'NO_SUCH_ROOM') {
          _handleNoSuchRoom(message);
        } else if (result == 'NOT_INVITED') {
          _handleNotInvited(message);
        }
        break;

        case 'ROOM_USERS':
        if (result == 'NO_SUCH_ROOM') {
          _handleNoSuchRoom(message);
        } else if (result == 'NOT_JOINED') {
          _handleNotJoinedRoom(message);
        }
        break;

      case 'ROOM_TEXT':
        if (result == 'NO_SUCH_ROOM') {
          _handleNoSuchRoom(message); 
        } else if (result == 'NOT_JOINED') {
          _handleNotJoinedRoom(message);
        }
        break;

      case 'LEAVE_ROOM':
        if (result == 'NO_SUCH_ROOM') {
          _handleNoSuchRoom(message); 
        } else if (result == 'NOT_JOINED') {
          _handleNotJoined(message); 
        }
        break;


      default:
        print("Unknown operation: $operation");
    }
  }

    void _promptForValidUsername() {
    String? newUsername;
    // Keep asking until a valid username is provided
    do {
      print('Ingresa un nuevo nombre de usuario (no debe ser vacio ni mayor a 8 caracteres)) :');
      newUsername = stdin.readLineSync();
    } while (newUsername == null || newUsername.trim().isEmpty);

    // Once a valid username is entered, call set_username
    clientMessages.set_username(newUsername.trim());
  }

  void _handleNewUser(Map<String, dynamic> message) {
    String username = message['username'];
    print('Usuario "$username" se ha unido al servidor.');
  }

  // Handle "NEW_STATUS" type messages
  void _handleNewStatus(Map<String, dynamic> message) {
    String username = message['username'];
    String status = message['status'];
    print('Usuario "$username" cambió su estado a: "$status".');
  }

  // Handle "USER_LIST" type messages
  void _handleUserList(Map<String, dynamic> message) {
    Map<String, dynamic> users = message['users'];
    StringBuffer userList = StringBuffer();

    users.forEach((username, status) {
      userList.writeln('$username: $status');
    });

    print('User list:\n$userList');
  }

  // Handle "TEXT_FROM" type messages
  void _handleTextFrom(Map<String, dynamic> message) {
    String sender = message['username'];
    String text = message['text'];
    print('$sender: $text');
  }

  // Handle "NO_SUCH_USER" result
  void _handleNoSuchUser(Map<String, dynamic> message) {
    String username = message['extra'];
    print('Error: El usuario "$username" no existe.');
  }

    // Handle "PUBLIC_TEXT_FROM" type messages
  void _handlePublicTextFrom(Map<String, dynamic> message) {
    String username = message['username'];
    String text = message['text'];

    // Notify the user with the formatted message
    print('$username [general]: $text');
  }

    // Handle success for "NEW_ROOM" operation
  void _handleNewRoomSuccess(Map<String, dynamic> message) {
    String roomName = message['extra'];
    print('Cuarto "$roomName" Se ha creado exitosamente.');
  }

  // Handle room already exists case
  void _handleRoomAlreadyExists(Map<String, dynamic> message) {
    String roomName = message['extra'];
    print('El cuarto "$roomName" ya existe.');
  }

    // Handle "INVITATION" type messages
  void _handleInvitation(Map<String, dynamic> message) {
    String username = message['username'];
    String roomName = message['roomname'];
    print('Has sido invitado a $roomName por $username.');
  }

    // Handle "NO_SUCH_ROOM" result for INVITE operation
  void _handleNoSuchRoom(Map<String, dynamic> message) {
    String roomName = message['extra'];
    print('No existe un cuarto llamado $roomName.');
  }

  // Handle "JOIN_ROOM" success result
  void _handleJoinRoomSuccess(Map<String, dynamic> message) {
    String roomName = message['extra'];
    print('Te has unido a $roomName.');
  }

    // Handle "JOINED_ROOM" type messages
  void _handleJoinedRoom(Map<String, dynamic> message) {
    String username = message['username'];
    String roomName = message['roomname'];
    print('$username se ha unido a la $roomName.');
  }

  // Handle "NOT_JOINED" result for ROOM_TEXT operation
  void _handleNotJoinedRoom(Map<String, dynamic> message) {
    String roomName = message['extra'];
    print('No eres parte de $roomName.');
  }

    // Handle "NOT_INVITED" result for JOIN_ROOM operation
  void _handleNotInvited(Map<String, dynamic> message) {
    String roomName = message['extra'];
    print('No has sido invitado a $roomName.');
  }

  // Handle "ROOM_USER_LIST" type messages
  void _handleRoomUserList(Map<String, dynamic> message) {
    String roomName = message['roomname'];
    Map<String, dynamic> users = message['users'];

    StringBuffer userList = StringBuffer();
    userList.writeln('Room: $roomName\n');
    
    users.forEach((username, status) {
      userList.writeln('$username: $status');
    });

    print(userList.toString());
  }

   // Handle "ROOM_TEXT_FROM" type messages
  void _handleRoomTextFrom(Map<String, dynamic> message) {
    String username = message['username'];
    String roomName = message['roomname'];
    String text = message['text'];

    print('$username [$roomName]: $text');
  }

  // Handle "LEFT_ROOM" type messages
  void _handleLeftRoom(Map<String, dynamic> message) {
    String username = message['username'];
    String roomName = message['roomname'];

    print('$username se ha ido de $roomName.');
  }

  // Handle "DISCONNECTED" type messages
  void _handleDisconnected(Map<String, dynamic> message) {
    String username = message['username'];
    print('Se ha desconectado $username.');
  }


}