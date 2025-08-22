# Netcat Chat Server

A simple, multi-client chat server built in Go using TCP for communication. This project enables multiple users to connect to a central server, send messages, and interact with each other in real-time. It also supports commands for managing usernames, viewing message history, and exiting the chat.

## Features

- **TCP-based communication**: Built on a TCP server for real-time messaging between clients.
- **Real-time messaging**: Any message sent by one user is broadcasted to all other connected users.
- **User management**:
  - Users can set their username (with uniqueness validation).
  - Users can rename their username during their session.
  - Users can exit the chat gracefully.
- **Message history**: All messages are saved, allowing new clients to view the entire message history when they join.
- **Help system**: The chat includes a built-in help system to guide users with available commands.
- **Concurrency**: The server can handle multiple client connections at once using Go’s goroutines and mutexes for synchronization.

## Prerequisites

- **Go 1.16 or higher** – The project is developed using Go, so you’ll need to have it installed.
- **Basic knowledge of Go and TCP networking** – This project uses Go’s `net` package for socket communication.
- **Basic understanding of concurrency** – The server handles multiple clients concurrently using goroutines and mutex locks.

## Installation

To set up the project locally:

1. Clone th repo ,
2. usage : go run . (port)
3. from other terminal or other device :
- ** nc [ip_adress] [port] 

- **congratulations now you are in net cat 
