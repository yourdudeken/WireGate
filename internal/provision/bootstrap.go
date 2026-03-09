package provision

import (
	"fmt"
	"github.com/yourdudeken/wiregate/internal/ssh"
)

func Bootstrap(client *ssh.Client) error {
	fmt.Println("Starting VPS Provisioning...")

	commands := []struct {
		desc string
		cmd  string
	}{
		{"Updating system packages", "sudo apt-get update -y && sudo apt-get upgrade -y"},
		{"Installing dependencies", "sudo apt-get install -y curl gnupg lsb-release iptables wireguard-tools fail2ban"},
		{"Checking for Docker", "docker --version >/dev/null 2>&1 || (curl -fsSL https://get.docker.com -o get-docker.sh && sudo sh get-docker.sh)"},
		{"Enabling IP forwarding", "echo 'net.ipv4.ip_forward=1' | sudo tee /etc/sysctl.d/99-wiregate.conf && sudo sysctl -p /etc/sysctl.d/99-wiregate.conf"},
		{"Configuring Firewall (UFW)", "sudo ufw allow 22/tcp && sudo ufw allow 80/tcp && sudo ufw allow 443/tcp && sudo ufw allow 51820/udp && sudo ufw --force enable"},
		{"Configuring Fail2Ban", "sudo cp /etc/fail2ban/jail.conf /etc/fail2ban/jail.local && sudo systemctl enable fail2ban && sudo systemctl start fail2ban"},
	}

	for _, c := range commands {
		fmt.Printf("[%s]...\n", c.desc)
		if err := client.Run(c.cmd); err != nil {
			return fmt.Errorf("failed during '%s': %w", c.desc, err)
		}
	}

	fmt.Println("VPS Provisioning Complete!")
	return nil
}
