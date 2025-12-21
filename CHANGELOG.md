# Changelog

## [0.1.0] - 2025-12-21

### Added
- PostgreSQL support with embedded binaries (zero dependencies)
- Redis support with embedded binaries (zero dependencies)
- Multi-engine architecture supporting multiple database types
- Interactive mode with engine selection, optional instance naming
- Commands: start, stop, pause, resume, list, url, status
- Instance management by name or ID
- Auto port allocation and isolated data directories
- Persistent and ephemeral modes (--persist flag)
- Pause/resume functionality to save resources
- Beautiful CLI UI with Bubble Tea (spinners, prompts, tables)
- Cross-platform support (macOS, Linux, Windows)
- Automated GitHub Actions workflows for building binaries
- RDB snapshot persistence for Redis instances

### Technical
- Go-based CLI with Cobra command framework
- Embedded PostgreSQL using fergusstrange/embedded-postgres
- Embedded Redis using pre-built binaries from GitHub releases
- Instance metadata stored in ~/.instant-db/
- Graceful signal handling (Ctrl+C)
- Health checks with connection pooling
