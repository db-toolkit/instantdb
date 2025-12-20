# instant-db

Instant, isolated PostgreSQL instances for development. Zero configuration.

## Features

- 🚀 Start PostgreSQL instances in seconds
- 🔒 Fully isolated - no conflicts with existing databases
- 🧹 Clean shutdown - zero traces left behind
- 💻 Cross-platform - macOS, Linux, Windows
- 🎯 Auto-assigned ports - no configuration needed
- 📦 Persistent or ephemeral - your choice
- ⏸️ Pause/resume instances to save resources

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
# Start a PostgreSQL instance (interactive prompts for credentials)
instant-db start

# List all instances
instant-db list

# Get connection URL
instant-db url <instance-id>

# Pause instance (stops process, keeps data)
instant-db pause <instance-id>

# Resume paused instance
instant-db resume <instance-id>

# Stop and clean up
instant-db stop <instance-id>
```

## Commands

```bash
# Start a new instance (interactive)
instant-db start

# Start with flags (non-interactive)
instant-db start -u myuser -password mypass --name myapp --port 5432 --persist

# Stop instance (removes data unless --persist was used)
instant-db stop <instance-id>

# Pause instance (stops process, always keeps data)
instant-db pause <instance-id>

# Resume paused instance
instant-db resume <instance-id>

# List all instances
instant-db list

# Get connection URL
instant-db url <instance-id>

# Check instance status
instant-db status <instance-id>

# Show version
instant-db --version

# Show help
instant-db --help
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
instant-db pause <id>  # Free resources when not working
instant-db resume <id> # Continue later
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

1. **Embedded PostgreSQL** - Binaries downloaded automatically on first use
2. **Isolated Data Directories** - Each instance gets its own data directory in `~/.instant-db/data/`
3. **Auto Port Allocation** - Finds available ports automatically
4. **Metadata Tracking** - Instance info stored in `~/.instant-db/*.json`
5. **Clean Shutdown** - Graceful process termination
6. **Zero Traces** - Non-persistent instances are completely removed on stop

## Roadmap

- [x] PostgreSQL support
- [x] Pause/resume functionality
- [ ] MySQL support
- [ ] MongoDB support
- [ ] Snapshot/restore functionality
- [ ] Clone from existing database
- [ ] Pre-load data from SQL file

## Contributing

Contributions welcome! Please open an issue or PR.

## License

MIT

## Part of DB Toolkit

instant-db is part of the [DB Toolkit](https://github.com/db-toolkit) organization - a collection of modern database tools for developers.

**Other tools:**
- [Migrator](https://github.com/db-toolkit/migrator) - Database migrations made simple
- [DB Toolkit Desktop](https://github.com/db-toolkit/db-toolkit-electron) - Cross-platform database management GUI
