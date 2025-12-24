# Changelog

## [Unreleased]

### Added
- MySQL support with embedded binaries
- Cross-platform library path handling (DYLD_LIBRARY_PATH, LD_LIBRARY_PATH, PATH)
- Windows-specific process management (netstat, taskkill)
- PowerShell installer for Windows
- Integration tests for PostgreSQL, MySQL, and Redis
- Interactive credential customization prompt

### Changed
- Updated README to include MySQL examples and Windows installation
- Improved interactive mode with visual feedback after selections

## [0.1.0] - 2025-12-21

### Added
- PostgreSQL, MySQL, and Redis support with embedded binaries (zero dependencies)
- Multi-engine architecture supporting multiple database types
- Interactive mode with engine selection and optional instance naming
- Commands: start, stop, pause, resume, list, url, status
- Instance management by name or ID
- Auto port allocation and isolated data directories
- Persistent and ephemeral modes (--persist flag)
- Pause/resume functionality to save resources
- Beautiful CLI UI with Bubble Tea (spinners, prompts, tables)
- Cross-platform support (macOS, Linux, Windows)
