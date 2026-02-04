# W-G Gateway

```
 __      __  ________        ________         __                                           
/  \    /  |/  _____/       /  _____/  ____ _/  |_  ____ __  _  __ _____  ___.__.          
\   \/\/   /   \  ___      /   \  ____/ __ \\   __\/ __ \\ \/ \/ // __ \<   |  |          
 \        /    \    \     \    \_\  \  ___/ |  | \  ___/ \     /\  ___/ \___  |          
  \__/\  / \______  /      \______  /\___  >|__|  \___  > \/\_/  \___  >/ ____|          
       \/         \/               \/     \/           \/             \/ \/               
```

**WireGuard VPS-to-Home Gateway Tool**

`wg-gateway` is a production-grade infrastructure tool designed to securely expose home servers (behind CGNAT, 4G, or dynamic IPs) to the public internet. It automates the orchestration of WireGuard tunnels and Traefik reverse proxies using a strictly **command-driven workflow**.

---

## Core Philosophy: Zero Manual Configuration

*   **No SSH Required**: The tool manages the VPS via automated orchestration.
*   **No File Editing**: All configurations are handled via the CLI.
*   **Hub-and-Spoke**: One VPS can securely tunnel to multiple distributed home server nodes.
*   **Zero-Setup DNS**: Built-in `sslip.io` support for those without a domain.
*   **Security First**: Automated hardening with Fail2Ban, UFW, and context-aware authentication.

---

## Feature Highlights

*   **Automated VPS Provisioning**: One command to transform a fresh Ubuntu VPS into a secured gateway.
*   **Multi-Node Support**: Connect multiple home servers (peers) to a single VPS hub.
*   **Observability Suite**: Real-time health checks, Handshake auditing, and centralized log streaming.
*   **Proactive Monitoring**: Background monitoring with Discord and Telegram webhook alerts.
*   **Automated Backups**: Snapshot your configuration and Let's Encrypt certificates to local storage or S3.
*   **Service Templates**: One-click configuration for Plex, Home Assistant, Jellyfin, and more.
*   **Multi-Hub Support**: Manage multiple gateway environments (contexts) from a single CLI.

---

## Installation

### The One-Liner (Recommended)
This script detects your OS and architecture, downloads the latest release, and installs it to your path.
```bash
curl -sSfL https://raw.githubusercontent.com/yourdudeken/wg-gateway/main/scripts/install.sh | sh
```

### Manual Installation (Go)
```bash
go install github.com/yourdudeken/wg-gateway@latest
```

---

## Quick Start (5-Step Workflow)

### 1. Initialize & Setup
Initialize your project and run the setup command to configure your local firewall (UFW) to allow tunnel traffic.
```bash
wg-gateway init --ip 1.2.3.4 --user root --key ~/.ssh/id_ed25519 --email admin@example.com
wg-gateway setup
```

### 2. Add a Peer (Node)
Add your home server node. providing a unique name.
```bash
wg-gateway peer add warehouse-lab
```

### 3. Add Services
Route domains to your peers. Use built-in templates for common apps.
```bash
wg-gateway service add-template plex myplex
# Result: myplex.1.2.3.4.sslip.io -> Local Port 32400
```

### 4. Deploy to VPS
Bootstrap the remote server and sync configurations.
```bash
wg-gateway deploy --bootstrap
```

### 5. Start the Tunnel
Launch the secure tunnel on your local machine.
```bash
wg-gateway up home
```

---

## Advanced Management

### Proactive Monitoring & Alerts
The tool can run in "Watcher" mode to monitor the health of your VPS, tunnels, and services. Configure alerts via Discord or Telegram.

```bash
# Configure alerts
wg-gateway config monitor.discord.url "https://discord.com/api/webhooks/..."
wg-gateway config monitor.discord.enabled true
wg-gateway config monitor.interval 10

# Start the monitor
wg-gateway monitor
```

### Automated Backups
Protect your configuration and SSL certificates from data loss.

```bash
# Set backup location
wg-gateway config backup.local_path ./backups

# Run backup
wg-gateway backup
```
*Note: This zips your local `config.yaml` and fetches the `letsencrypt/acme.json` file from your VPS.*

### Multi-Hub Contexts
Manage multiple independent VPS gateways (e.g., US and Europe) using the `-c` flag.

```bash
# List all hub contexts
wg-gateway hub list

# Switch context for a specific command
wg-gateway -c europe.yaml status
```

### Web UI Dashboard
A modern, password-protected graphical interface for managing your gateway.

```bash
wg-gateway web --password MySecurePass
```
Username: `admin` | Default Port: `8080`

---

## Command Reference

### System & Setup
| Command | Description |
|---------|-------------|
| `setup` | Configures local UFW and verifies system readiness |
| `init` | Creates a new gateway configuration file |
| `hub list` | Lists all available gateway contexts |

### Infrastructure
| Command | Description |
|---------|-------------|
| `deploy` | Provisions VPS and syncs configurations (`--bootstrap` for fresh VPS) |
| `status` | Checks production readiness and project overview |
| `check` | Live connectivity audit (Ping peers, check ports) |
| `logs` | Stream logs from VPS or any Home peer |

### Node & Service Management
| Command | Description |
|---------|-------------|
| `peer add` | Register a new home server node (supports `--ssh-user`) |
| `peer list` | View all configured nodes and their tunnel IPs |
| `ssh [peer]` | Securely jump into a peer node from the VPS |
| `service add` | Map a domain to a local port (Supports sslip.io) |
| `service add-template`| Add a service with pre-configured app defaults |
| `rotate-keys` | Regenerate all WireGuard keypairs for the project |

---

## Security

`wg-gateway` implements several hardening measures automatically:
*   **Fail2Ban**: Installed and configured during bootstrap to block brute-force attacks.
*   **UFW Firewall**: Defaults to "deny all" with explicit allows for SSH (22), WireGuard (51820), and Web (80/443).
*   **Basic Auth**: The Web Dashboard is protected by Basic Authentication (Username: `admin`).
*   **Encapsulated Networking**: Internal services are never exposed to the VPS host network; they communicate strictly over the encrypted `wg0` interface.

---

## Contributing

Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
