# BaconC2

<img src="https://github.com/RabbITCybErSeC/BaconC2/blob/main/images/baconc2_withoutbackground.png" width="230">

Bacon is a small C2 framework, which aims to support stealthy communication techniques.

> **Development status:** This project is currently in **active development**. There is **no beta release yet** — features, interfaces, and behavior may change as the project evolves.

**Note:** This repository is intended as a learning and research platform. Use responsibly and only on systems where you have explicit permission to test.

## Why?

BaconC2 was created as a learning platform to experiment with advanced C2 techniques, focusing on stealth, security, and flexibility. It aims to simulate real-world C2 scenarios while providing a foundation for research into malware development and network pivoting.


## Features

- **Modular Client Architecture**:
  - Separates concerns into packages (`agent`, `executor`, `commands`, `transport`, `sysinfo`, `queue`, `models`) for maintainability and extensibility.
  - Supports asynchronous command execution and result queuing for efficient operation.
- **Stealthy Communication**:
  - Implements multiple transport protocols:
    - HTTP for standard command/result exchange.
    - To be done: WebSocket for real-time interactive shell sessions.
    - To be done: UDP for lightweight, connectionless communication with shell session support.
  - to be done: Configurable beacon intervals with jitter potential for evading detection.
- **Command Execution**:
  - Built-in commands like `sys_info` for detailed system telemetry (e.g., network interfaces, CPU, memory, uptime).
  - Interactive shell support for `bash`, `sh`, `powershell`, and `cmd`, enabling real-time command execution.
  - Extensible `CommandRegistry` for adding custom commands.
- **Agent Capabilities**:
  - Gathers minimal system info (hostname, IP, OS, protocol) during beaconing.
  - Provides extended system info (architecture, disk, processes, etc.) on demand.
  - Supports session management for interactive shells with session IDs.
- **Security**:
  - JSON-based payloads for flexible data exchange.
  - Thread-safe transport and queue implementations.
  - Planned encryption for shell sessions and sensitive data (e.g., AES).


## Goals

- Develop an intuitive web interface for managing connected agents and executing commands.
- Explore unconventional C2 channels (e.g., DNS, ICMP) with secure encryption.
- Enhance agent capabilities with features like memory dumps, token stealing, and process injection.
- Enable network pivoting by chaining agents for lateral movement.
- Improve reliability with persistent queues, retry mechanisms, and adaptive beaconing.
- Add comprehensive testing and observability (logging, metrics, tracing).

## Getting Started

*Instructions for setting up BaconC2 will be added as the project matures. Currently, the client is under active development, focusing on core functionality and transport protocols.*


## BaconC2 API Documentation

### Authentication

#### Login
- **Endpoint**: `POST /api/v1/auth/login`
- **Description**: Authenticate and receive a JWT token
- **Headers**: 
  - `Content-Type: application/json`
- **Body**: 
  ```json
  {
    "username": "string",
    "password": "string"
  }
  ```
- **Response**: JWT token for authenticated requests

### Frontend API

### List Agents
- **Endpoint**: `GET /api/v1/frontend/agents`
- **Description**: Get list of all agents
- **Headers**: 
  - `Authorization: Bearer <jwt_token>`
- **Response**: List of agent objects

### General API

#### Add Command to Queue
- **Endpoint**: `GET /api/v1/general/queue/command/:id`
- **Description**: Add a command to the execution queue
- **Headers**: 
  - `Authorization: Bearer <jwt_token>`
- **Query Parameters**: 
  - `command`: The command to execute
- **Response**: Command queued status

#### Example:

```bash
curl -X POST http://localhost:8081/api/agents/command?id={id} -H 'Content-Type: application/json' -d '{ "command": "whoami"}' -v
```

## Contributing

Contributions are welcome! Please submit issues or pull requests to the [GitHub repository](https://github.com/RabbITCybErSeC/BaconC2). Focus areas include:

- Implementing new transport protocols.
- Adding advanced agent features (e.g., persistence, evasion).
- Enhancing security with encryption and authentication.
- Writing tests and documentation.

## License to be done

*To be determined. For now, BaconC2 is a research project and not intended for production use.*

## Disclaimer

BaconC2 is for educational and research purposes only. Unauthorized use in production environments or against systems without permission is illegal and unethical.
