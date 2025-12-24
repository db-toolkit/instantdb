# Instant DB

A CLI tool that spins up isolated database instances instantly for development, with zero configuration.

## Features

- 🚀 Start PostgreSQL, MySQL, or Redis instances in seconds
- 🔒 Fully isolated - no conflicts with existing databases
- 🧹 Clean shutdown - zero traces left behind
- 💻 Cross-platform - macOS, Linux, Windows
- 🎯 Auto-assigned ports - no configuration needed
- 📦 Persistent or ephemeral - your choice
- ⏸️ Pause/resume instances to save resources
- 🎨 Beautiful interactive CLI with prompts and tables
- 📛 Manage instances by name or ID

## Supported Databases

- **PostgreSQL** - Embedded binaries, zero dependencies
- **MySQL** - Embedded binaries, zero dependencies
- **Redis** - Embedded binaries, zero dependencies

## Installation

### macOS / Linux

```bash
# Quick install
curl -sSL https://raw.githubusercontent.com/db-toolkit/instantdb/main/install.sh | bash

# Or manual install
git clone https://github.com/db-toolkit/instant-db.git
cd instant-db
go build -o instant-db src/instantdb/cmd/instantdb/main.go
sudo mv instant-db /usr/local/bin/
```

### Windows

```powershell
# Quick install (run in PowerShell as Administrator)
irm https://raw.githubusercontent.com/db-toolkit/instantdb/main/install.ps1 | iex

# Or download manually from releases
# https://github.com/db-toolkit/instantdb/releases
```

## Quick Start

```bash
# Start an instance (interactive - choose PostgreSQL, MySQL, or Redis)
instant-db start

# Start PostgreSQL with a name
instant-db start -e postgres --name my-app

# Start MySQL with a name
instant-db start -e mysql --name my-db

# Start Redis with a name
instant-db start -e redis --name my-cache

# List all instances
instant-db list

# Get connection URL (by name or ID)
instant-db url my-app

# Pause instance (stops process, keeps data)
instant-db pause my-cache

# Resume paused instance
instant-db resume my-cache

# Stop and clean up
instant-db stop my-app
```

## Commands

```bash
# Start a new instance (interactive)
instant-db start

# Start with flags (non-interactive)
instant-db start -e postgres --name myapp -u myuser --password mypass --port 5432 --persist
instant-db start -e mysql --name mydb -u root --password mypass --persist
instant-db start -e redis --name mycache --password mypass --persist

# Stop instance (removes data unless --persist was used)
instant-db stop <name-or-id>

# Pause instance (stops process, always keeps data)
instant-db pause <name-or-id>

# Resume paused instance
instant-db resume <name-or-id>

# List all instances
instant-db list

# Get connection URL
instant-db url <name-or-id>

# Check instance status
instant-db status <name-or-id>

# Show version
instant-db --version

# Show help
instant-db --help
```

## Examples

### PostgreSQL
```bash
# Start PostgreSQL
instant-db start -e postgres --name my-postgres

# Get connection URL
instant-db url my-postgres
# Output: postgresql://postgres:postgres@localhost:54321/postgres

# Connect with psql
psql $(instant-db url my-postgres)
```

### Redis
```bash
# Start Redis
instant-db start -e redis --name my-redis --password secret

# Get connection URL
instant-db url my-redis
# Output: redis://:secret@localhost:63791

# Connect with redis-cli
redis-cli -h 127.0.0.1 -p 63791 -a secret
```

### MySQL
```bash
# Start MySQL
instant-db start -e mysql --name my-mysql

# Get connection URL
instant-db url my-mysql
# Output: mysql://root@127.0.0.1:50762/mysql

# Connect with mysql client
mysql -h 127.0.0.1 -P 50762 -u root
```

## Contributing

Contributions welcome! Please open an issue or PR.

## License

MIT

## Part of DB Toolkit

instant-db is part of the [DB Toolkit](https://github.com/db-toolkit) organization - a collection of modern database tools for developers.

**Other tools:**
- [Migrator](https://github.com/db-toolkit/migrator) - Database migrations made simple
- [DB Toolkit Desktop](https://github.com/db-toolkit/db-toolkit-electron) - Cross-platform database management GUI
