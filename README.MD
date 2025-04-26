# go_remox_server

**go_remox_server** is a lightweight remote input server written in Go.  
It allows remote control of the mouse, keyboard and multimedia functions via network.


## ✨ Features

- 🖱️ Mouse control
- ⌨️ Keyboard control (In progress)
- 📦 Supports two message formats: **Protobuf** and **JSON**
- 🔒 Simple and efficient architecture written in pure Go

## 🛠️ Platform Support

| Platform | Mouse Support | Keyboard Support | Media Control Support |
|----------|---------------|------------------|-----------------------|
| Windows  | ✅ Supported   | 🚧 In Progress   | ❓ Planned             |
| macOS    | ✅ Supported   | 🚧 In Progress   | ❓ Planned             |
| Linux    | ❓ Planned     | ❓ Planned        | ❓ Planned             |

## 🚀 Getting Started

```bash
go run main.go
```

Make sure to create a .env file in the root directory before running the server.

## ⚙️ Configuration

The server reads settings from a .env file:

```env
HOST=127.0.0.1              # Server IP address
PORT=8080                   # Listening port
MESSAGE_FORMAT=protobuf     # Message format: protobuf or json
TEST_HTML_FILE=tester.html  # Path to optional test HTML file
```

## 📦 Usage

This server is intended to be used with a compatible client (such as a mobile app or desktop UI) that sends control
commands.