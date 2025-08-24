# 🖧 TCP Chat Project

This project is a simple TCP-based chat server built in Go.  
Before diving into the code, it’s important to understand how **TCP (Transmission Control Protocol)** works.  

---

## 🔹 What is TCP?  

TCP is **not a network** by itself. It is a **protocol** used in computer networks to guarantee that data:  
- Arrives in the correct order.  
- Is delivered reliably without loss.  
- Is reassembled properly at the destination.  

TCP works by:  
- Splitting data into small chunks (called **segments**).  
- Numbering each segment.  
- Checking that each segment arrives safely.  
- Requesting retransmission if something is missing.  
- Reassembling everything back in the correct order.  

---

## 🔹 The TCP Handshake (3-Way Connection Setup)  

When two machines want to communicate using TCP, they first establish a connection through a **3-way handshake**:  

1. **SYN** → The sender says: *"Hello, I want to connect."*  
2. **SYN-ACK** → The receiver replies: *"Okay, I understand, let’s start."*  
3. **ACK** → The sender confirms: *"Perfect, we can now exchange data."*  

Once this handshake is done:  
- Data is exchanged in numbered segments.  
- Each segment is acknowledged (ACK).  
- Lost segments are retransmitted.  
- At the end, the connection is closed gracefully.  

---

## 🔹 How TCP Works (Step by Step Example)  

Let’s imagine **Machine A** wants to connect to **Machine B**:  

- **Step 1: Machine A → Machine B**  
  Machine A sends a **SYN** with an initial sequence number (ISN), e.g., `2000`.  
  👉 *"I want to open a connection, starting at 2000."*  

- **Step 2: Machine B → Machine A**  
  Machine B replies with:  
  - **SYN = 1** (and chooses its own ISN, e.g., `1000`).  
  - **ACK = 2001** (`ISN_A + 1`).  
  👉 *"I accept your request (ACK=2001), and my starting number is 1000."*  

- **Step 3: Machine A → Machine B**  
  Machine A responds with:  
  - **ACK = 1001** (`ISN_B + 1`).  
  👉 *"I received your 1000, so here’s 1001. The connection is now ready!"*  

📌 In short:  
- Machine A: `2000 → 2001`  
- Machine B: `1000 → 1001`  

➡️ Both machines are synchronized and ready to communicate reliably.  

---

## 🚀 Next Steps in the Project  

This repository implements the server side of a **TCP chat application**.  
The current features:  
- Listening for clients on port `8989`.  
- Accepting a connection from a client.  
- Sending a logo/message to the client when connected.  

The upcoming features to implement:  
- Handle multiple clients at the same time.  
- Exchange messages between clients in real-time.  
- Properly close connections when clients disconnect.  

---
