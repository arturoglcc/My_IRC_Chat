import 'dart:convert'; 
import 'client.dart'; 
import 'handle_server_messages.dart'; 
import 'dart:io';

class ServerMessages {
  final Client client;
  ServerMessages(this.client);

  // This function will process the JSON message from the server
   void processMessage(String message) {
    Map<String, dynamic> jsonMessage = jsonDecode(message);

    switch (jsonMessage['type']) {
      case 'RESPONSE':
        handleResponse(jsonMessage);
        break;
      case 'NEW_USER':
        HandleServerMessages.handleNewUser(jsonMessage);
        break;
      case 'NEW_STATUS':
        HandleServerMessages.handleNewStatus(jsonMessage);
        break;
      case 'USER_LIST':
        HandleServerMessages.handleUserList(jsonMessage);
        break;
      case 'TEXT_FROM':
        HandleServerMessages.handleTextFrom(jsonMessage);
        break;
      case 'PUBLIC_TEXT_FROM':
        HandleServerMessages.handlePublicTextFrom(jsonMessage);
        break;
      case 'DISCONNECTED':
        HandleServerMessages.handleDisconnected(jsonMessage);
        break;
       case 'INVITATION':
        HandleServerMessages.handleInvitation(jsonMessage);
        break;
      case 'JOINED_ROOM':
      HandleServerMessages.handleJoinedRoom(jsonMessage);
      case 'ROOM_USER_LIST':
        HandleServerMessages.handleRoomUserList(jsonMessage);
        break;
      case 'ROOM_TEXT_FROM':
        HandleServerMessages.handleRoomTextFrom(jsonMessage);
        break;
      case 'LEFT_ROOM':
        HandleServerMessages.handleLeftRoom(jsonMessage);
        break;
      default:
        print("Unknown message type: ${jsonMessage['type']}");
    }
  }

   void handleResponse(Map<String, dynamic> message) {
    String operation = message['operation'];
    String result = message['result'];

    // Handle specific operations and results
    switch (operation) {
      case 'IDENTIFY':
        if (result == 'USER_ALREADY_EXISTS') {
          String existingUsername = message['extra'];
          print('nombre de usuario "$existingUsername" ya existe. \n');
          promptForValidUsername();
        } else if (result == 'SUCCESS') {
          String newUsername = message['extra'];
          client.updateUsername(newUsername); 
        }
        break;

      case 'NEW_USER':
        HandleServerMessages.handleNewUser(message);
        break;

      case 'TEXT':
        if (result == 'NO_SUCH_USER') {
          HandleServerMessages.handleNoSuchUser(message);
        }
        break; 

      case 'NEW_ROOM':
        if (result == 'SUCCESS') {
          HandleServerMessages.handleNewRoomSuccess(message);
        } else if (result == 'ROOM_ALREADY_EXISTS') {
          HandleServerMessages.handleRoomAlreadyExists(message);
        }
        break;

      case 'INVITE':
        if (result == 'NO_SUCH_ROOM') {
          HandleServerMessages.handleNoSuchRoom(message);
        } else if (result == 'NO_SUCH_USER') {
          HandleServerMessages.handleNoSuchUser(message);
        }
        break;

      case 'JOIN_ROOM':
        if (result == 'SUCCESS') {
          HandleServerMessages.handleJoinRoomSuccess(message);
        } else if (result == 'NO_SUCH_ROOM') {
          HandleServerMessages.handleNoSuchRoom(message);
        } else if (result == 'NOT_INVITED') {
          HandleServerMessages.handleNotInvited(message);
        }
        break;

        case 'ROOM_USERS':
        if (result == 'NO_SUCH_ROOM') {
          HandleServerMessages.handleNoSuchRoom(message);
        } else if (result == 'NOT_JOINED') {
          HandleServerMessages.handleNotJoinedRoom(message);
        }
        break;

      case 'ROOM_TEXT':
        if (result == 'NO_SUCH_ROOM') {
          HandleServerMessages.handleNoSuchRoom(message); 
        } else if (result == 'NOT_JOINED') {
          HandleServerMessages.handleNotJoinedRoom(message);
        }
        break;

      case 'LEAVE_ROOM':
        if (result == 'NO_SUCH_ROOM') {
          HandleServerMessages.handleNoSuchRoom(message); 
        } else if (result == 'NOT_JOINED') {
          HandleServerMessages.handleNotJoinedRoom(message); 
        }
        break;


      default:
        print("Unknown operation: $operation");
    }
  }

  void promptForValidUsername() {
    String? newUsername;
    // Keep asking until a valid username is provided
    do {
      print('Ingresa un nuevo nombre de usuario (no debe ser vacio ni mayor a 8 caracteres)) :');
      newUsername = stdin.readLineSync();
    } while (newUsername == null || newUsername.trim().isEmpty);

    // Once a valid username is entered, call set_username
    client.set_username(newUsername.trim());
  }
}

