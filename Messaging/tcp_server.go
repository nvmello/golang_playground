package main

// Import required packages for networking, concurrency, etc.
import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

// Define a struct to represent a connected client
type ClientConnection struct {
	id            string
	connection    net.Conn
	subscriptions map[string]bool
}

type Server struct {
	clients sync.Map                       // Maps client addresses to ClientConnection objects [clientID, ClientConnection]
	topics  map[string][]*ClientConnection // Maps topic names to lists of subscribed clients
	mutex   sync.Mutex                     // Protects access to the topics map
}

func main() {
	// Initialize your server
	var server *Server = &Server{
		topics: make(map[string][]*ClientConnection),
	}

	// Start listening for TCP connections on localhost:9000
	var listener net.Listener
	var err error
	listener, err = net.Listen("tcp", "localhost:9000")

	// Handle any errors that might occur when starting the server
	if err != nil {
		fmt.Println("error starting server", err)
		return
	}

	// Remember to defer closing the listener when the program exits
	defer listener.Close()

	// Print a message indicating server started successfully
	fmt.Println("Server started on localhost:9000")

	// Begin an infinite loop to accept new connections
	for {
		// Accept a new connection
		var conn net.Conn
		// Handle any errors that might occur when accepting
		conn, err = listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Launch a goroutine to handle this connection
		fmt.Println("New GoRoutine")
		go handleConnection(server, conn)
	}
}

func handleConnection(server *Server, conn net.Conn) {
	//Create a new client
	clientID := conn.RemoteAddr().String()
	client := &ClientConnection{
		id:            clientID,
		connection:    conn,
		subscriptions: make(map[string]bool),
	}
	//Store the client in our clients map
	server.clients.Store(client.id, client)
	fmt.Println("Handling new connection from", clientID)

	scanner := bufio.NewScanner(conn)

	const maxCap = 1024 * 1024
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, maxCap)

	defer func() {
		// Cleanup when connection closes
		server.mutex.Lock()
		for topic := range client.subscriptions {
			subscribers := server.topics[topic]
			var newList []*ClientConnection
			for _, c := range subscribers {
				if c.id != client.id {
					newList = append(newList, c)
				}
			}
			server.topics[topic] = newList
		}
		server.mutex.Unlock()
		server.clients.Delete(clientID)
		conn.Close()
		fmt.Println("Closed connection from:", clientID)
	}()

	for scanner.Scan() {
		input := scanner.Text()
		if input == "" {
			continue
		}

		// Process this single command
		handleCommand(server, client, input)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from client:", err)
	}
}

func handleCommand(server *Server, client *ClientConnection, input string) {
	input = strings.TrimSpace(input)
	parts := strings.SplitN(input, " ", 3)

	if len(parts) < 2 {
		client.connection.Write([]byte("Invalid command\n"))
		return
	}

	command := strings.ToUpper(parts[0])
	topic := parts[1]

	switch command {

	case "SUB":
		server.mutex.Lock()
		server.topics[topic] = append(server.topics[topic], client)
		server.mutex.Unlock()
		client.subscriptions[topic] = true
		client.connection.Write([]byte("Subbed to " + topic + "\n"))

	case "PUBLISH":
		if len(parts) < 3 {
			client.connection.Write([]byte("Missing message for publish\n"))
			return
		}
		message := parts[2]

		server.mutex.Lock()
		subscribers := server.topics[topic]
		server.mutex.Unlock()

		for _, sub := range subscribers {
			if sub.id != client.id {
				sub.connection.Write([]byte("[" + topic + "] " + message + "\n"))
			}
		}
	default:
		client.connection.Write([]byte("Unknown command\n"))

	}
}
