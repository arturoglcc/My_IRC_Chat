# Directorios de los archivos fuente
SERVER_SRC = ./server/*.go
CLIENT_SRC = ./client/main.dart

# Nombre del ejecutable del servidor y cliente
SERVER_EXECUTABLE = ./server/server   # Ejecutable del servidor dentro de la carpeta 'server'
CLIENT_EXECUTABLE = ./client/client   # Ejecutable del cliente dentro de la carpeta 'client'

# Comando para compilar el servidor (usando Go)
build-server:
	@echo "Compilando servidor..."
	go build -o $(SERVER_EXECUTABLE) $(SERVER_SRC)
	@echo "Servidor compilado: $(SERVER_EXECUTABLE)"

# Comando para compilar el cliente (usando Dart)
build-client:
	@echo "Compilando cliente..."
	dart compile exe $(CLIENT_SRC) -o $(CLIENT_EXECUTABLE)
	@echo "Cliente compilado: $(CLIENT_EXECUTABLE)"

# Comando para compilar ambos: servidor y cliente
build-all: build-server build-client
	@echo "Compilados ambos: servidor y cliente"

# Limpieza de los ejecutables generados
clean:
	@echo "Eliminando ejecutables..."
	rm -f $(SERVER_EXECUTABLE) $(CLIENT_EXECUTABLE)
	@echo "Limpieza completa"

