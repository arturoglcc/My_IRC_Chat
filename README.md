Usa make build-all para contruir un ejecutable de el servidor y del cliente.

Usa make clean para borrar los ejecutables hechos.

Usa server/server para iniciar el servidor.

Usa client/client para iniciar el cliente.

SERVER
punto de entrada:
El servidor en su punto de entrada lee si el usuario especifico un puerto para conectarse, si lo especifico y es valido se conecta a este puerto, sino se conecta a el puerto 1234 y crea un servidor

Server
Define las estructuras server, client y room, maneja la conexión de nuevos usuarios y se asegura de que su primer mensaje sea identificandose, si todo sale bien escucha mensajes del cliente en una go rutina, 
si algo falla se lo avisa al cliente. Ademas que termina el servidor con gracia si ha sido cerrado abruptamente.

listener
Escucha los mensajes del cliente y manda a llamar a una función dependiendo de que mensaje se ha recibido

users
Define dos comportamientos que puede tener un usuario, crear un cuarto e invitar a alguien a un cuarto  

handle-messages
maneja los mensajes que recibe del cliente mediante una función previamente llamada por listener

writer
Escribe la mayoria de mensajes que se le envian al cliente

auxiliar-writer
define algunas operaciones comunes al momento de escribir mensajes

writer-errors
Escribe mensajes que son errores de uso del cliente




CLIENT
punto de entrada
si el usuario especifica a que maquina quiere conectarse y a que puerto quiere conectarse y son validos, se conecta ahi, sino, se intenta conectar a localhost y al puerto 1234.
Despues crea instancias cliente, user imput listener y server listener

client
tiene las operaciones update-username y disconnect

server-listener
unicamente recibe los mensajes que manda el servidor y los manda a procesar a una instancia de servermessages

server-messages
identifica el tipo de mensaje, su operación y resultado y en base a eso lo manda a procesar a handle-server-messages

handle-server-messages
tiene la mayoria de operaciones necesarias para avisar al usuario del estado de su conexión al servidor

user-input-listener
unicamente recibe lo que haya escrito el usuario y lo procesa con una instancia de client-messages

client-messages
aqui tanto se identifica que tipo de mensaje es como se procesa el mensaje para enviar algo al servidor

writer
unicamente hace jsons y se los envia al servidor.



COMO USAR EL CLIENTE
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


