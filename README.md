# INSTANT-DB

A CLI tool that spins up isolated database instances instantly for development, with zero configuration.

## Features

- 🚀 Start PostgreSQL instances in seconds
- 🔒 Fully isolated - no conflicts with existing databases
- 🧹 Clean shutdown - zero traces left behind
- 💻 Cross-platform - macOS, Linux, Windows
- 🎯 Auto-assigned ports - no configuration needed
- 📦 Persistent or ephemeral - your choice
- ⏸️ Pause/resume instances to save resources

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

## Contributing

Contributions welcome! Please open an issue or PR.

## Part of DB Toolkit

instant-db is part of the [DB Toolkit](https://github.com/db-toolkit) organization - a collection of modern database tools for developers.

**Other tools:**
- [Migrator](https://github.com/db-toolkit/migrator) - Database migrations made simple
- [DB Toolkit Desktop](https://github.com/db-toolkit/db-toolkit-electron) - Cross-platform database management GUI
