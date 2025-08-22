package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	lib "net-cat/project-lib"
)

var (
	members  []lib.Member
	latestID int
	mu       sync.Mutex
)

var (
	red   = "\033[31m"
	green = "\033[32m"
	cyan  = "\033[36m"
	reset = "\033[0m"
)

func main() {
	// Get command-line arguments (skip the program name itself)
	args := os.Args[1:]

	// Default port to listen on if none is provided
	port := "8989"

	// If exactly one argument is given → use it as the port
	if len(args) == 1 {
		port = args[0]
	} else if len(args) != 0 {
		// If more than one argument is given → print usage and exit
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}

	// Start listening on the given TCP port
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error starting server :", err)
		return
	}

	// Make sure the listener is closed when main exits
	defer listener.Close()

	fmt.Println("Listening on the port :", port)

	// Clear previous chat history (using a function from your lib package)
	// 'mu' is probably a mutex used to protect concurrent access
	err = lib.ClearMessageHistory(&mu)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Main server loop → wait for clients to connect
	for {
		// Accept a new incoming connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection :", err)
			continue
		}

		// Handle the client connection in a separate goroutine
		// This way multiple clients can be served at the same time
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("client connected from ", conn.RemoteAddr())

	_, err := conn.Write([]byte(lib.GenerateGreeting()))
	if err != nil {
		fmt.Println("Error writing to client :", err)
		return
	}

	memberID := latestID
	latestID++
	memberName := ""
	for memberName == "" {
		memberName, err = lib.CheckName(conn, &members, &mu)
		if err != nil {
			if err.Error() == "alreadyExistsError" {
				conn.Write([]byte(red + "You chose a name that is already used. Please try again.\n" + reset))
			} else if err.Error() == "emptyNameError" {
				conn.Write([]byte(red + "You chose an empty name. Please try again.\n" + reset))
			} else {
				conn.Write([]byte(red + "Please try again.\n" + reset))
			}
		}
	}

	if err := lib.AddMember(memberName, memberID, conn, &members, &mu); err != nil {
		return
	}

	history, err := lib.GetMessageHistory(&mu)
	if err != nil {
		fmt.Println("Error get messages history :", err)
		lib.RemoveMember(memberID, &members, &mu)
		return
	}
	_, err = conn.Write([]byte(history))
	if err != nil {
		fmt.Println("Error writing :", err)
		lib.RemoveMember(memberID, &members, &mu)
		return
	}
	lib.BroadcastMessage(memberID, "\n"+green+memberName+" has joined our chat...\n"+reset, members)

	for {
		now := time.Now()
		formatted := now.Format("2006-01-02 15:04:05")

		_, err = conn.Write([]byte(cyan + "[" + formatted + "]" + reset + green + "[" + memberName + "]:" + reset))
		if err != nil {
			fmt.Println("Error writing :", err)
			break
		}

		buf := make([]byte, 1024)

		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading :", err)
			break
		}

		memberMessage := strings.TrimSpace(string(buf[:n]))
		if memberMessage != "" {
			switch memberMessage {
			case "@help":
				helpMessage := lib.Help()
				_, err = conn.Write([]byte(helpMessage))
				if err != nil {
					fmt.Println("Error writing :", err)
					break
				}
			case "@exit":
				lib.BroadcastMessage(memberID, "\n"+red+memberName+" has left our chat...\n"+reset, members)
				lib.RemoveMember(memberID, &members, &mu)
				return
			case "@rename":
				memberName = lib.Rename(memberID, &members, &mu, conn)
			default:
				message := cyan + "[" + formatted + "]" + reset + green + "[" + memberName + "]:" + reset + string(buf[:n])
				lib.BroadcastMessage(memberID, "\n"+message, members)
				err = lib.LogMessages(&mu, message)
				if err != nil {
					fmt.Println("Error log the new message ", err)
					break
				}
			}
		}
	}

	lib.BroadcastMessage(memberID, "\n"+red+memberName+" has left our chat...\n"+reset, members)
	lib.RemoveMember(memberID, &members, &mu)
}
