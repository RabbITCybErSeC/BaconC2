# Filesystem Commands with State Management

BaconC2 now supports stateful filesystem commands that maintain the agent's current working directory across command executions.

## Architecture

### Agent State
The agent maintains runtime state including:
- **Current Working Directory**: Tracked across `cd` commands
- **Environment Variables**: Can be set/retrieved (future expansion)

### Command Types

#### Stateless Handlers
Traditional handlers that don't need state:
```go
func Handler(cmd models.Command) models.CommandResult
```

#### Stateful Handlers
Handlers that receive execution context with state access:
```go
func Handler(ctx *command_handler.CommandContext) models.CommandResult
```

## Available Commands

### `pwd` - Print Working Directory
Returns the current working directory maintained by the agent.

**Usage:**
```json
{
  "type": "intern",
  "command": "pwd",
  "args": []
}
```

**Example Response:**
```
/home/user/documents
```

---

### `cd` - Change Directory
Changes the agent's current working directory. Supports both absolute and relative paths.

**Usage:**
```json
{
  "type": "intern",
  "command": "cd",
  "args": ["/path/to/directory"]
}
```

**Examples:**

Absolute path:
```json
{"type": "intern", "command": "cd", "args": ["/etc"]}
```

Relative path:
```json
{"type": "intern", "command": "cd", "args": ["../parent"]}
```

Go to home directory:
```json
{"type": "intern", "command": "cd", "args": []}
```

**Response:**
```
Changed directory to: /etc
```

---

### `ls` - List Directory
Lists contents of a directory. Uses current working directory if no path specified.

**Usage:**
```json
{
  "type": "intern",
  "command": "ls",
  "args": ["/optional/path"]
}
```

**Examples:**

List current directory:
```json
{"type": "intern", "command": "ls", "args": []}
```

List specific directory:
```json
{"type": "intern", "command": "ls", "args": ["/var/log"]}
```

List relative path:
```json
{"type": "intern", "command": "ls", "args": ["../sibling"]}
```
