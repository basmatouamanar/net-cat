package projectlib

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

// Members nodes
type Member struct {
	ID         int
	Name       string
	Connection net.Conn
}

// GenerateGreeting returns the welcome message and ASCII art
// that is displayed to a new client when they connect.
func GenerateGreeting() string {
	reset := "\033[0m"
	boldGreen := "\033[1;32m"
	cyan := "\033[36m"

	greeting := boldGreen + "Welcome to TCP-Chat!\n" + reset
	greeting += cyan + "         _nnnn_\n"
	greeting += "        dGGGGMMb\n"
	greeting += "       @p~qp~~qMb\n"
	greeting += "       M|@||@) M|\n"
	greeting += "       @,----.JM|\n"
	greeting += "      JS^\\__/  qKL\n"
	greeting += "     dZP        qKRb\n"
	greeting += "    dZP          qKKb\n"
	greeting += "   fZP            SMMb\n"
	greeting += "   HZM            MMMM\n"
	greeting += "   FqM            MMMM\n"
	greeting += "__| \".        |\\dS\"qML\n"
	greeting += " |    `.       | `' \\Zq\n"
	greeting += "_)      \\.___.,|     .'\n"
	greeting += "\\____   )MMMMMP|   .'\n"
	greeting += "     `-'       `--'\n" + reset

	return greeting
}

// Help returns a formatted help message for the TCP-Chat application.
// It includes a brief description of the application and a list of available commands
func Help() string {
	reset := "\033[0m"
	boldCyan := "\033[1;36m"
	boldYellow := "\033[1;33m"

	helpMessage := boldCyan + "\n\nWelcome to the TCP-Chat application!\n" + reset
	helpMessage += "\tThis is a simple project used to establish a connection \n\tbetween different users in the LAN or on the same machine.\n\n"
	helpMessage += boldYellow + "Options:\n" + reset
	helpMessage += "\t@help        ==> Show this help window\n"
	helpMessage += "\t@exit        ==> Disconnect from the chat\n"
	helpMessage += "\t@rename name ==> Change your display name\n\n"

	return helpMessage
}

// BroadcastMessage sends the given message to all connected clients
// except the sender identified by senderID.
func BroadcastMessage(senderID int, message string, members []Member) {
	green := "\033[32m"
	cyan := "\033[36m"
	reset := "\033[0m"

	for _, member := range members {
		if member.ID != senderID {

			_, err := member.Connection.Write([]byte(message))
			if err != nil {
				fmt.Println("Error writing to", member.Name, ":", err)
				return
			}

			now := time.Now()
			formatted := now.Format("2006-01-02 15:04:05")
			_, err = member.Connection.Write([]byte(cyan + "[" + formatted + "]" + reset + green + "[" + member.Name + "]:" + reset))
			if err != nil {
				fmt.Println("Error writing to", member.Name, ":", err)
				return
			}
		}
	}
}

// LogMessages appends a new message to the "chat.log" file.
func LogMessages(mu *sync.Mutex, message string) error {
	mu.Lock()
	defer mu.Unlock()

	f, err := os.OpenFile("chat.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(message)
	return err
}

// GetMessageHistory reads the entire chat history from "messages.txt".
func GetMessageHistory(mu *sync.Mutex) (string, error) {
	mu.Lock()
	defer mu.Unlock()
	content, err := os.ReadFile("chat.log")
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// ClearMessageHistory removes all content from "messages.txt".
func ClearMessageHistory(mu *sync.Mutex) error {
	mu.Lock()
	defer mu.Unlock()
	err := os.WriteFile("chat.log", []byte(""), 0o644)
	if err != nil {
		return err
	}
	return nil
}

// CheckName prompts a connected client to enter a name and validates it.
// It sends a message to the client asking for their name, reads the input,
// trims whitespace and newline characters, and checks if the name is already used
// by other members or if it is empty. Access to the members slice is protected
// with a mutex to ensure safe concurrent access.
func CheckName(connection net.Conn, members *[]Member, mu *sync.Mutex) (string, error) {
	reset := "\033[0m"
	boldGreen := "\033[1;32m"

	_, err := connection.Write([]byte(boldGreen + "[ENTER YOUR NAME]:" + reset))
	if err != nil {
		fmt.Println("Error writing to client :", err)
		return "", errors.New("writingError")
	}

	buf := make([]byte, 1024)

	n, err := connection.Read(buf)
	if err != nil {
		fmt.Println("Error reading :", err)
	}

	name := strings.TrimSpace(strings.ReplaceAll(string(buf[:n]), "\n", ""))

	alreadyExists := false
	mu.Lock()
	defer mu.Unlock()
	for i := 0; i < len(*members); i++ {
		if (*members)[i].Name == name {
			alreadyExists = true
		}
	}
	if alreadyExists {
		return "", errors.New("alreadyExistsError")
	}
	if name == "" {
		return "", errors.New("emptyNameError")
	}
	return name, nil
}

// AddMember adds a new member to the chat room if there is space available.
// If the chat room is full (maximum 10 members), it repeatedly notifies the client
// that they must wait until a spot becomes free, pausing for 3 seconds between messages.
// Access to the members slice is protected with a mutex to ensure safe concurrent updates.
func AddMember(name string, id int, connection net.Conn, members *[]Member, mu *sync.Mutex) error {
	red := "\033[31m"
	reset := "\033[0m"

	for len(*members) == 10 {
		_, err := connection.Write([]byte(red +
			"There is no empty place in this chat room. Please wait until a spot is free.\n" +
			reset))
		if err != nil {
			return errors.New("write error: " + err.Error())
		}
		time.Sleep(3 * time.Second)
	}
	mu.Lock()
	*members = append(*members, Member{ID: id, Name: name, Connection: connection})
	defer mu.Unlock()
	return nil
}

// RemoveMember removes a member from the chat room based on their ID.
func RemoveMember(id int, members *[]Member, mu *sync.Mutex) {
	mu.Lock()
	defer mu.Unlock()
	for i := 0; i < len(*members); i++ {
		if id == (*members)[i].ID {
			*members = append((*members)[0:i], (*members)[i+1:]...)
			return
		}
	}
}

// Rename allows a member to change their display name in the chat room.
// It finds the member by their ID, prompts them to enter a new name using CheckName,
// and ensures the new name is not empty or already in use. The function sends
// appropriate error messages back to the client if the chosen name is invalid.
// Once a valid name is entered, it updates the member's name in the members slice.
func Rename(id int, members *[]Member, mu *sync.Mutex, connection net.Conn) string {
	memberName := ""
	var err error
	red := "\033[31m"
	reset := "\033[0m"

	defer mu.Unlock()
	for i := 0; i < len(*members); i++ {
		if id == (*members)[i].ID {
			for memberName == "" {
				memberName, err = CheckName(connection, members, mu)
				if err != nil {
					if err.Error() == "alreadyExistsError" {
						connection.Write([]byte(red + "You chose a name that is already used. Please try again.\n" + reset))
					} else if err.Error() == "emptyNameError" {
						connection.Write([]byte(red + "You chose an empty name. Please try again.\n" + reset))
					} else {
						connection.Write([]byte(red + "Please try again.\n" + reset))
					}
				}
			}
			mu.Lock()
			(*members)[i].Name = memberName
			return memberName
		}
	}
	return memberName
}
