# instant-db

Instant, isolated PostgreSQL instances for development. Zero configuration.

## Features

- 🚀 Start PostgreSQL instances in seconds
- 🔒 Fully isolated - no conflicts with existing databases
- 🧹 Clean shutdown - zero traces left behind
- 💻 Cross-platform - macOS, Linux, Windows
- 🎯 Auto-assigned ports - no configuration needed
- 📦 Persistent or ephemeral - your choice

## Requirements

**None!** PostgreSQL binaries are automatically downloaded on first use.

- Binaries are cached in `~/.embedded-postgres-go/`
- Works on macOS, Linux, and Windows
- No manual installation needed

## Installation

```bash
# Clone the repository
git clone https://github.com/db-toolkit/instant-db.git
cd instant-db

# Build
go build -o instant-db src/instantdb/cmd/instantdb/main.go

# Move to PATH (optional)
sudo mv instant-db /usr/local/bin/
```

## Quick Start

```bash
# Start a PostgreSQL instance
instant-db start

# Output:
# 🚀 Starting PostgreSQL instance...
# 
# ✅ PostgreSQL instance started successfully!
# 
#   Instance ID:  a1b2c3d4e5f6...
#   Name:         postgres-a1b2c3d4
#   Port:         54321
#   Connection:   postgresql://localhost:54321/postgres

# Get connection URL
instant-db url a1b2c3d4e5f6

# List all instances
instant-db list

# Check instance status
instant-db status a1b2c3d4e5f6

# Stop and clean up
instant-db stop a1b2c3d4e5f6
```

## Usage

### Start an instance

```bash
# Basic start (auto-assigned port, ephemeral)
instant-db start

# With custom name
instant-db start --name myapp

# With specific port
instant-db start --port 5432

# Persistent data (survives stop)
instant-db start --persist

# All options
instant-db start --name myapp --port 5432 --persist
```

### List instances

```bash
instant-db list

# Output:
# 📋 Running Instances (2)
# 
#   • myapp
#     ID:     a1b2c3d4e5f6...
#     Port:   5432
#     Status: running
```

### Get connection URL

```bash
instant-db url a1b2c3d4e5f6

# Output:
# postgresql://localhost:5432/postgres
```

### Check status

```bash
instant-db status a1b2c3d4e5f6

# Output:
# 📊 Instance Status: a1b2c3d4e5f6
# 
#   Running:  ✅ Yes
#   Healthy:  ✅ Yes
#   Message:  ok
```

### Stop instance

```bash
instant-db stop a1b2c3d4e5f6

# Output:
# 🛑 Stopping instance a1b2c3d4e5f6...
# ✅ Instance stopped successfully!
```

## Use Cases

### Quick Testing
```bash
instant-db start
# Run your tests
instant-db stop <id>
```

### Feature Branch Development
```bash
instant-db start --name feature-auth --persist
# Develop your feature
instant-db stop <id>
```

### Integration with Migrator
```bash
# Start database
instant-db start --name myapp

# Get URL and set environment
export DATABASE_URL=$(instant-db url <id>)

# Run migrations
migrator init
migrator makemigrations "initial"
migrator migrate

# Clean up
instant-db stop <id>
```

## How It Works

1. **Isolated Data Directories** - Each instance gets its own data directory in `~/.instant-db/data/`
2. **Auto Port Allocation** - Finds available ports automatically
3. **Metadata Tracking** - Instance info stored in `~/.instant-db/*.json`
4. **Clean Shutdown** - Graceful SIGTERM with fallback to SIGKILL
5. **Zero Traces** - Non-persistent instances are completely removed on stop

## Architecture

```
instant-db/
├── cmd/instantdb/          # CLI commands
│   ├── main.go            # Entry point
│   └── commands/          # Modular commands
├── internal/
│   ├── engines/           # Database engines
│   │   ├── engine.go      # Engine interface
│   │   └── postgres.go    # PostgreSQL implementation
│   ├── types/             # Shared types
│   │   ├── config.go
│   │   ├── instance.go
│   │   └── status.go
│   └── utils/             # Utilities
│       ├── id.go          # ID generation
│       ├── network.go     # Port allocation
│       ├── process.go     # Process management
│       └── storage.go     # Metadata storage
```

## Roadmap

- [x] PostgreSQL support
- [ ] MySQL support
- [ ] SQLite support
- [ ] MongoDB support
- [ ] Snapshot/restore functionality
- [ ] Clone from existing database
- [ ] Pre-load data from SQL file
- [ ] Docker-free embedded binaries

## Contributing

Contributions welcome! Please open an issue or PR.

## License

MIT

## Part of DB Toolkit

instant-db is part of the [DB Toolkit](https://github.com/db-toolkit) organization - a collection of modern database tools for developers.

**Other tools:**
- [Migrator](https://github.com/db-toolkit/migrator) - Database migrations made simple
- [DB Toolkit Desktop](https://github.com/db-toolkit/db-toolkit-electron) - Cross-platform database management GUI
