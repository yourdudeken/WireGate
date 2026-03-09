# Changelog

All notable changes to the `wiregate` project will be documented in this file.

## [2.2.2] - 2026-02-04
### Added
- **Peer SSH Tunneling**: New `ssh [peer]` command allows users to jump directly from the VPS hub into any home server node via the encrypted WireGuard network.
- **SSH User Orchestration**: Added `--ssh-user` support to peer management for automated login.

## [2.2.1] - 2026-01-31
### Fixed
- **CI/CD Permissions**: Fixed missing `contents: write` permissions in Release workflow preventing automatic tagging.

## [2.2.0] - 2026-01-31
### Added
- **Proactive Monitoring**: New `monitor` command with background auditing and real-time alerts.
    - Discord Webhook integration
    - Telegram Bot integration
    - VPS, Peer, and Service port health polling
- **Automated Backups**: New `backup` command for disaster recovery.
    - Zips local `config.yaml` and remote Let's Encrypt certificates
    - Supports local storage and S3-compatible buckets (AWS, R2)
- **Web UI Authentication (Basic Auth)**: Dashboard security using the `--password` flag or `WG_ADMIN_PASS` environment variable.
- **Multi-Hub Support (Contexts)**: Use the `-c` or `--config` flag to manage multiple independent gateways from the same CLI.
- **Service Templates**: One-click configuration for popular apps (Plex, Home Assistant, Jellyfin, Pi-hole, etc.).
- **Zero-Setup DNS (sslip.io)**: Automated hostname mapping to VPS IP for users without custom domains.
- **Local Setup Engine**: New `setup` command to automate local environment preparation (UFW config, Docker verification).
- **Improved Health Auditing**: Deep-packet verification of service connectivity in the `check` command.
- **Enhanced Log Filtering**: Peer-specific log streaming and improved remote log orchestration.

## [2.1.0] - 2026-01-27
### Added
- **Web UI Dashboard**: Modern, dark-themed web interface for visual gateway management.
    - Real-time status monitoring with live statistics
    - Visual peer and service management
    - Configuration editor with form validation
    - Responsive design for desktop and mobile
    - RESTful API backend with embedded static files
- **Enhanced CLI Commands**:
    - `web`: Launch the web-based dashboard interface
    - `up [peer]`: Start local peer tunnel via Docker Compose
    - `down [peer]`: Stop local peer tunnel
    - `service update`: Modify existing service port mappings
- **Improved Configuration Management**:
    - Dot-notation support in `config` command (e.g., `vps.ip`, `proxy.email`)
    - Enhanced validation and error messages

### Changed
- Removed all emojis from codebase for enterprise compatibility
- Updated README with comprehensive command reference and web UI documentation
- Improved CLI help text and command descriptions

## [2.0.0] - 2026-01-27
### Added
- **Multi-Node Support**: Refactored architecture to support Hub-and-Spoke topology (one VPS, many home servers).
- **Command-Driven UI**: Complete CLI lifecycle management. Users no longer need to edit files or SSH manually.
- **Peer Management**: 
    - New `peer` command group for managing multiple home servers
    - Automatic IP assignment for new peers
    - Key generation and rotation for all peers
- **Service Management**: 
    - `service add` with `--peer` flag for multi-node routing
    - `service remove` for cleanup
    - `service list` with peer association display
- **Security Hardening**: 
    - Automated **Fail2Ban** installation and configuration on VPS
    - Custom SSH identity file support via `--key` and `config vps.key`
    - UFW firewall automation with explicit port rules
- **Observability Suite**:
    - `check` command for live end-to-end connectivity auditing
    - `logs` command for streaming remote (VPS) and local (Home) container logs
    - Enhanced `status` command with production-readiness audit
- **CLI Aesthetics**: 
    - Professional ASCII banner
    - Clean, emoji-free output for enterprise logs

### Changed
- Refactored `config.yaml` to support a `Peers` slice instead of a single `Home` object
- Updated all WireGuard and Traefik templates for multi-peer support
- Improved `deploy` command with better error reporting and automated SCP/SSH orchestration
- Enhanced `generate` command to create peer-specific deployment folders

### Breaking Changes
- Configuration file structure changed from `Home` to `Peers` array
- Service definitions now require `peer_name` field
- Deployment files now generated in `deploy/peers/[name]` instead of `deploy/home`

## [1.0.0] - 2026-01-25
### Added
- Initial release of the `wiregate` CLI
- Basic WireGuard point-to-point tunnel orchestration
- Support for Traefik reverse proxy with automated Let's Encrypt
- Core commands: `init`, `generate`, `deploy`, `status`
- Docker-based deployments for both VPS and Home server
- Automated WireGuard key generation
- Basic service routing configuration

---
*Note: This project follows Semantic Versioning.*
