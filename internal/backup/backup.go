package backup

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/yourdudeken/wiregate/internal/config"
	"github.com/yourdudeken/wiregate/internal/ssh"
)

func Run(cfg *config.Config) (string, error) {
	// 1. Create temporary directory for gathering files
	tmpDir, err := os.MkdirTemp("", "wiregate-backup-*")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmpDir)

	// 2. Copy local config
	if err := copyFile("config.yaml", filepath.Join(tmpDir, "config.yaml")); err != nil {
		return "", fmt.Errorf("failed to copy config.yaml: %v", err)
	}

	// 3. Fetch certificates from VPS
	fmt.Println("Fetching certificates from VPS...")
	client := ssh.NewClient(cfg.VPS.SSHUser, cfg.VPS.IP, cfg.VPS.SSHKey)
	// Traefik stores them in ./letsencrypt/acme.json relative to the compose file
	// We assume they are in /root/letsencrypt or similar if deployed via tool
	vpsCertPath := "letsencrypt" // Default relative path on VPS
	localCertPath := filepath.Join(tmpDir, "letsencrypt")
	
	if err := client.Fetch(vpsCertPath, localCertPath); err != nil {
		fmt.Printf("Warning: Could not fetch certificates from VPS: %v\n", err)
	}

	// 4. Create ZIP archive
	timestamp := time.Now().Format("20060102-150405")
	archiveName := fmt.Sprintf("wiregate-backup-%s.zip", timestamp)
	archivePath := filepath.Join(os.TempDir(), archiveName)

	if err := zipSource(tmpDir, archivePath); err != nil {
		return "", fmt.Errorf("failed to create zip: %v", err)
	}

	// 5. Store locally if configured
	if cfg.Backup.LocalPath != "" {
		finalPath := filepath.Join(cfg.Backup.LocalPath, archiveName)
		if err := os.MkdirAll(cfg.Backup.LocalPath, 0755); err != nil {
			return "", err
		}
		if err := copyFile(archivePath, finalPath); err != nil {
			return "", err
		}
		fmt.Printf("Backup saved locally: %s\n", finalPath)
	}

	return archivePath, nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func zipSource(source, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name, _ = filepath.Rel(source, path)
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})

	return nil
}
