package ssh

import (
	"fmt"
	"os"
	"os/exec"
)

type Client struct {
	User    string
	Host    string
	KeyPath string
}

func NewClient(user, host, keyPath string) *Client {
	return &Client{User: user, Host: host, KeyPath: keyPath}
}

func (c *Client) Run(command string) error {
	dest := fmt.Sprintf("%s@%s", c.User, c.Host)
	args := []string{"-o", "BatchMode=yes", "-o", "ConnectTimeout=10", "-o", "StrictHostKeyChecking=no"}
	if c.KeyPath != "" {
		args = append(args, "-i", c.KeyPath)
	}
	args = append(args, dest, command)

	cmd := exec.Command("ssh", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (c *Client) Output(command string) (string, error) {
	dest := fmt.Sprintf("%s@%s", c.User, c.Host)
	args := []string{"-o", "BatchMode=yes", "-o", "ConnectTimeout=5", "-o", "StrictHostKeyChecking=no"}
	if c.KeyPath != "" {
		args = append(args, "-i", c.KeyPath)
	}
	args = append(args, dest, command)

	cmd := exec.Command("ssh", args...)
	out, err := cmd.Output()
	return string(out), err
}

func (c *Client) Copy(src, dst string) error {
	dest := fmt.Sprintf("%s@%s:%s", c.User, c.Host, dst)
	args := []string{"-o", "BatchMode=yes", "-r"}
	if c.KeyPath != "" {
		args = append(args, "-i", c.KeyPath)
	}
	args = append(args, src, dest)

	cmd := exec.Command("scp", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (c *Client) Fetch(src, dst string) error {
	remote := fmt.Sprintf("%s@%s:%s", c.User, c.Host, src)
	args := []string{"-o", "BatchMode=yes", "-o", "ConnectTimeout=10", "-r"}
	if c.KeyPath != "" {
		args = append(args, "-i", c.KeyPath)
	}
	args = append(args, remote, dst)

	cmd := exec.Command("scp", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
