# GoWire

A Redis-inspired in-memory key-value server built from scratch using Go.

The goal of this project is to understand how Redis-like systems work internally, exploring TCP networking, protocol design, command execution, concurrency, and data storage.

## Architecture

The following diagram illustrates the current architecture of GoWire, from client connections to command execution and storage.

![GoWire Architecture](./Arch%20Diagram%20v1.png)

## Features

### Networking
- TCP server listening on `localhost:6140`
- Multiple concurrent client connections
- Independent goroutine for each client connection

### RESP Protocol
- RESP command parsing
- RESP response serialization
- Reader and Writer abstractions

### Supported Commands

#### Strings

| Command | Description |
|---|---|
| `PING` | Check server connectivity |
| `SET key value` | Store a key-value pair |
| `GET key` | Retrieve a value |
| `DEL key` | Delete a key |
| `EXISTS key` | Check whether a key exists |

#### Hashes

| Command | Description |
|---|---|
| `HSET key field value` | Set a hash field |
| `HGET key field` | Retrieve a hash field |
| `HDEL key field` | Delete a hash field |
| `HEXISTS key field` | Check whether a hash field exists |

### Concurrency
- Goroutines for handling multiple clients
- `sync.RWMutex` for thread-safe shared storage
- Read locks for concurrent read operations
- Write locks for storage modifications

## Project Structure

```text
gowire/
├── cleint/
    ├── go.mod
    ├── main.go
    ├── resp.go
├── server/
    ├── keyValues/
        ├── keyValues.go
    ├── go.mod
    ├── main.go
    ├── resp_parse.go
    ├── resp_write.go
    ├── dispatch.go
    ├── validate.go
├── README.md
└── Arch Diagram v1.png
```

*The structure may evolve as the project grows.*

## Getting Started

### Prerequisites

- Go installed on your machine
- A terminal

### Run the Server

Clone the repository:

```bash
git clone <your-repository-url>
cd gowire
```

Start the server:

```bash
go run .
```

The server listens on:

```text
localhost:6140
```

### Connect Using a Redis Client

If your server implements the compatible RESP commands, you can connect using `redis-cli`:

```bash
redis-cli -p 6140
```

Try the following commands:

```redis
PING

SET name John
GET name
EXISTS name
DEL name

HSET user:1 name John
HGET user:1 name
HEXISTS user:1 name
HDEL user:1 name
```

## Learning Resources

This project was inspired by [Build Redis from Scratch](https://www.build-redis-from-scratch.dev/en/introduction) by Ahmed Ashraf.

## What's Next?

This project is a work in progress. Some areas I'd like to explore further:

- Additional Redis commands
- More advanced data structures (Lists and Sets)
- Data persistence
- Improved error handling
- Expiration and TTL support
- More comprehensive testing

## Acknowledgements

Thanks to Ahmed Ashraf for the excellent *Build Redis from Scratch* learning resource.

---

**GoWire — Learning how databases work, one command at a time.** 🚀