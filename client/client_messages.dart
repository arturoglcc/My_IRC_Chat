import 'writer.dart'; 

class ClientMessages {
  final Writer writter; 

  ClientMessages(this.writter); 


  void processMessage(String input) {
    if (input.startsWith('@all')) {
      _handlePublicMessage(input);
    } else if (input.contains('-->')) {
      _handleRoomTextMessage(input);
    }else if (input.startsWith('@')) {
      _handlePrivateMessage(input);
    }  else {
    if (input.isEmpty) return;
    List<String> parts = input.split(' ');
    String command = parts[0];

  switch (command) {
      case '/status':
        _handleStatusCommand(parts);
        break;

      case '/users':
        _sendUsersCommand();
        break;
        
      case '/cr':
          _handleCreateRoom(parts);
          break; 

      case '/inv':
          _handleInviteCommand(parts);
          break;

      case '/join':
          _handleJoinRoom(parts);
          break;    

      case '/ru':
          _handleRoomUsers(parts);
          break;    

      case '/lr':
          _handleLeaveRoom(parts);
          break;    
      case '/leave':
        _handleDisconnect();
        break;
      case '/help':
      _handleHelp();
      default:
        print('Unknown command.');
    }
  }
}

  // Handle the /status command and uppercase the status argument
  void _handleStatusCommand(List<String> parts) {
    if (parts.length > 1) {
      String status = parts[1].toUpperCase();
      if (_isValidStatus(status)) {
        _sendStatusMessage(status);
      } else {
        print('Invalid status. Valid options are: AWAY, ACTIVE, BUSY.');
      }
    } else {
      print('Please provide a status.');
    }
  }

  // Validate if the status is valid
  bool _isValidStatus(String status) {
    return ['AWAY', 'ACTIVE', 'BUSY'].contains(status);
  }

    // Convert the status message into a JSON and send it
  void _sendStatusMessage(String status) {
    Map<String, dynamic> statusMessage = {
      'type': 'STATUS',
      'status': status
    };

    writter.sendJsonMessage(statusMessage);
  }

    // Handle the /users command and send the JSON message
  void _sendUsersCommand() {
    Map<String, dynamic> usersMessage = {
      'type': 'USERS'
    };

    writter.sendJsonMessage(usersMessage); 
  }

  void _handlePrivateMessage(String input) {
    List<String> parts = input.split(' ');
    String recipient = parts[0].substring(1); // Extract username after '@'
    String messageText = parts.sublist(1).join(' '); // Join the rest of the message

    if (recipient.isNotEmpty && messageText.isNotEmpty) {
      _sendTextMessage(recipient, messageText);
    } else {
      print('Invalid message format. Use @username followed by your message.');
    }
  }

    // Convert the private message to JSON and send it
  void _sendTextMessage(String username, String text) {
    Map<String, dynamic> textMessage = {
      'type': 'TEXT',
      'username': username,
      'text': text
    };
  writter.sendJsonMessage(textMessage); 
  }

    // Handle the @all message and send public text
  void _handlePublicMessage(String input) {
    String messageText = input.substring(4).trim(); // Extract text after @all

    if (messageText.isNotEmpty) {
      sendPublicTextMessage(messageText);
    } else {
      print('No puedes enviar un mensaje vacio.');
    }
  }

    // Convert the public text message into a JSON and send it
  void sendPublicTextMessage(String text) {
    Map<String, dynamic> publicTextMessage = {
      'type': 'PUBLIC_TEXT',
      'text': text
    };

    writter.sendJsonMessage(publicTextMessage); // Use the writer's sendJsonMessage
  }

// Handle the /cr (create room) command
void _handleCreateRoom(List<String> parts) {
  if (parts.length > 1) {
    // Join the parts starting from index 1 (the room name) to include spaces
    String roomName = parts.sublist(1).join(' ');
    sendCreateRoomMessage(roomName);
  } else {
    print('Please provide a room name.');
  }
}


  // Convert the room creation command into a JSON message and send it
  void sendCreateRoomMessage(String roomName) {
    Map<String, dynamic> newRoomMessage = {
      'type': 'NEW_ROOM',
      'roomname': roomName
    };

    writter.sendJsonMessage(newRoomMessage); // Use the writer's sendJsonMessage
  }

// Handle the /inv command with room names that might contain spaces
void _handleInviteCommand(List<String> parts) {
  // Find where the usernames start (first part that begins with @)
  int usernameStartIndex = parts.indexWhere((part) => part.startsWith('@'));

  if (usernameStartIndex > 1) {
    // Join all parts before the first username as the room name
    String roomName = parts.sublist(1, usernameStartIndex).join(' ');

    // Extract usernames from the parts starting from usernameStartIndex
    List<String> usernames = extractUsernames(parts.sublist(usernameStartIndex));

    if (usernames.isNotEmpty) {
      _sendInviteMessage(roomName, usernames);
    } else {
      print('Please provide at least one valid username to invite.');
    }
  } else {
    print('Please provide a room name and at least one username.');
  }
}

// Extract usernames (remove @ and trim spaces)
List<String> extractUsernames(List<String> parts) {
  return parts
      .map((user) => user.replaceAll('@', '').trim()) // Remove "@" and trim spaces
      .where((user) => user.isNotEmpty) // Only non-empty usernames
      .toList();
}
  // Send the invite message in JSON format
  void _sendInviteMessage(String roomName, List<String> usernames) {
    Map<String, dynamic> inviteMessage = {
      'type': 'INVITE',
      'roomname': roomName,
      'usernames': usernames
    };
    writter.sendJsonMessage(inviteMessage);
  }

    // Handle the /join command
  void _handleJoinRoom(List<String> parts) {
    if (parts.length > 1) {
      String roomName = parts.sublist(1).join(' ');
      _sendJoinRoomMessage(roomName);
    } else {
      print('Please provide a room name.');
    }
  }

  // Send the join room message in JSON format
  void _sendJoinRoomMessage(String roomName) {
    Map<String, dynamic> joinRoomMessage = {
      'type': 'JOIN_ROOM',
      'roomname': roomName
    };

    writter.sendJsonMessage(joinRoomMessage); 
  }

  // Handle the /ru (room users) command
  void _handleRoomUsers(List<String> parts) {
    if (parts.length > 1) {
      String roomName = parts.sublist(1).join(' ');
      _sendRoomUsersMessage(roomName);
    } else {
      print('Please provide a room name.');
    }
  }

  // Send the room users message in JSON format
  void _sendRoomUsersMessage(String roomName) {
    Map<String, dynamic> roomUsersMessage = {
      'type': 'ROOM_USERS',
      'roomname': roomName
    };
    writter.sendJsonMessage(roomUsersMessage); 
  }


  void _handleRoomTextMessage(String input) {
    // Split the input at the '-->'
    List<String> parts = input.split('-->');
    if (parts.length > 1) {
      String roomNamePart = parts[0].trim();
      String textPart = parts.sublist(1).join('-->').trim();

      // Extract room name from @RoomName format
      if (roomNamePart.startsWith('@')) {
        String roomName = roomNamePart.substring(1).trim();
        _sendRoomTextMessage(roomName, textPart);
      } else {
        print('Formato de cuarto invalido.');
      }
    } else {
      print('Texto invalido .');
    }
  }

  // Send the room text message in JSON format
  void _sendRoomTextMessage(String roomName, String text) {
    Map<String, dynamic> roomTextMessage = {
      'type': 'ROOM_TEXT',
      'roomname': roomName,
      'text': text
    };

    writter.sendJsonMessage(roomTextMessage); 
  }

  // Handle the /lr (leave room) command
  void _handleLeaveRoom(List<String> parts) {
    if (parts.length > 1) {
      String roomName = parts.sublist(1).join(' ');
      _sendLeaveRoomMessage(roomName);
    } else {
      print('Por favor escribe un cuarto.');
    }
  }

  // Send the leave room message in JSON format
  void _sendLeaveRoomMessage(String roomName) {
    Map<String, dynamic> leaveRoomMessage = {
      'type': 'LEAVE_ROOM',
      'roomname': roomName
    };

    writter.sendJsonMessage(leaveRoomMessage); 
  }

    // Handle the /leave command
  void _handleDisconnect() {
    Map<String, dynamic> disconnectMessage = {
      'type': 'DISCONNECT'
    };

    writter.sendJsonMessage(disconnectMessage); 
  }

  void _handleHelp() {
      String commands = '''
Client commands:
1. /help - Displays this help message with a list of available commands.
2. /status [away|busy|active] - Set your current status (e.g., away, busy, or active).
3. /users - Displays a list of all users currently connected to the server.
4. @[username] [message] - Send a private message to a specific user.
5. /cr [roomname] - Create a new chat room with the given name.
6. /inv [roomname] @user1 @user2 ... @userN - invite users to a room you are in. 
7. /join [roomname] - Join an existing chat room.
8. /lr [roomname] - Leave the chat room.
9. /ru [roomname] - List all users in a specific chat room.
10. @all [message] - Send a message to all users in the chat.
11. @roomname --> [message] - Send a message to a specific chat room.
12. /leave - Disconnect from the server.
''';

  print(commands);
  }
}