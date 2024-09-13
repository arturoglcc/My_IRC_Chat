import 'client.dart'; // To access the Client class

class HandleServerMessages {
  final Client client;

  HandleServerMessages(this.client);


  static void handleNewUser(Map<String, dynamic> message) {
    String username = message['username'];
    print('Usuario "$username" se ha unido al servidor.');
  }

  // Handle "NEW_STATUS" type messages
  static void handleNewStatus(Map<String, dynamic> message) {
    String username = message['username'];
    String status = message['status'];
    print('Usuario "$username" cambió su estado a: "$status".');
  }

  // Handle "USER_LIST" type messages
  static void handleUserList(Map<String, dynamic> message) {
    Map<String, dynamic> users = message['users'];
    StringBuffer userList = StringBuffer();

    users.forEach((username, status) {
      userList.writeln('$username: $status');
    });

    print('User list:\n$userList');
  }

  // Handle "TEXT_FROM" type messages
  static void handleTextFrom(Map<String, dynamic> message) {
    String sender = message['username'];
    String text = message['text'];
    print('$sender: $text');
  }

  // Handle "NO_SUCH_USER" result
  static void handleNoSuchUser(Map<String, dynamic> message) {
    String username = message['extra'];
    print('Error: El usuario "$username" no existe.');
  }

    // Handle "PUBLIC_TEXT_FROM" type messages
  static void handlePublicTextFrom(Map<String, dynamic> message) {
    String username = message['username'];
    String text = message['text'];

    // Notify the user with the formatted message
    print('$username [general]: $text');
  }

    // Handle success for "NEW_ROOM" operation
  static void handleNewRoomSuccess(Map<String, dynamic> message) {
    String roomName = message['extra'];
    print('Cuarto "$roomName" Se ha creado exitosamente.');
  }

  // Handle room already exists case
  static void handleRoomAlreadyExists(Map<String, dynamic> message) {
    String roomName = message['extra'];
    print('El cuarto "$roomName" ya existe.');
  }

    // Handle "INVITATION" type messages
  static void handleInvitation(Map<String, dynamic> message) {
    String username = message['username'];
    String roomName = message['roomname'];
    print('Has sido invitado a $roomName por $username.');
  }

    // Handle "NO_SUCH_ROOM" result for INVITE operation
  static void handleNoSuchRoom(Map<String, dynamic> message) {
    String roomName = message['extra'];
    print('No existe un cuarto llamado $roomName.');
  }

  // Handle "JOIN_ROOM" success result
  static void handleJoinRoomSuccess(Map<String, dynamic> message) {
    String roomName = message['extra'];
    print('Te has unido a $roomName.');
  }

    // Handle "JOINED_ROOM" type messages
  static void handleJoinedRoom(Map<String, dynamic> message) {
    String username = message['username'];
    String roomName = message['roomname'];
    print('$username se ha unido a la $roomName.');
  }

  // Handle "NOT_JOINED" result for ROOM_TEXT operation
  static void handleNotJoinedRoom(Map<String, dynamic> message) {
    String roomName = message['extra'];
    print('No eres parte de $roomName.');
  }

    // Handle "NOT_INVITED" result for JOIN_ROOM operation
  static void handleNotInvited(Map<String, dynamic> message) {
    String roomName = message['extra'];
    print('No has sido invitado a $roomName.');
  }

  // Handle "ROOM_USER_LIST" type messages
  static void handleRoomUserList(Map<String, dynamic> message) {
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
  static void handleRoomTextFrom(Map<String, dynamic> message) {
    String username = message['username'];
    String roomName = message['roomname'];
    String text = message['text'];

    print('$username [$roomName]: $text');
  }

  // Handle "LEFT_ROOM" type messages
  static void handleLeftRoom(Map<String, dynamic> message) {
    String username = message['username'];
    String roomName = message['roomname'];

    print('$username se ha ido de $roomName.');
  }

  // Handle "DISCONNECTED" type messages
  static void handleDisconnected(Map<String, dynamic> message) {
    String username = message['username'];
    print('Se ha desconectado $username.');
  }
}