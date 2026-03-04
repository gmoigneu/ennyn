package ca

import (
	"bufio"
	"crypto/sha1"
	"crypto/x509"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/smallstep/truststore"
)

func sha1Sum(data []byte) [20]byte {
	return sha1.Sum(data)
}

// Install generates or loads the CA and leaf cert, then installs the CA into the system trust store.
func Install() error {
	fmt.Println("Ennyn Trust Setup")
	fmt.Println("=================")
	fmt.Println()
	fmt.Println("This will:")
	fmt.Println("  1. Generate a local Certificate Authority (if needed)")
	fmt.Println("  2. Generate a wildcard certificate for *.localhost")
	fmt.Println("  3. Install the CA into your system trust store")
	fmt.Println()
	fmt.Println("Step 3 requires administrator access (sudo).")
	fmt.Println()

	// Generate or load CA
	authority, generated, err := LoadOrGenerate()
	if err != nil {
		return fmt.Errorf("CA setup: %w", err)
	}

	if generated {
		fmt.Println("Generated new CA certificate.")
	} else {
		fmt.Println("Using existing CA certificate.")
	}

	// Generate or load leaf cert (force regen if CA was just generated)
	_, err = LoadOrGenerateLeaf(authority, generated)
	if err != nil {
		return fmt.Errorf("leaf cert setup: %w", err)
	}

	if generated {
		fmt.Println("Generated wildcard certificate for *.localhost.")
	} else {
		fmt.Println("Using existing wildcard certificate for *.localhost.")
	}

	// Install CA into system trust store
	fmt.Println()
	fmt.Println("Installing CA into system trust store...")
	if err := truststore.Install(authority.Cert, truststore.WithPrefix("Ennyn")); err != nil {
		return fmt.Errorf("trust store installation: %w", err)
	}
	fmt.Println("CA installed in system trust store.")

	// WSL detection
	if isWSL() {
		fmt.Println()
		fmt.Println("WSL detected. Your browser likely runs on Windows and uses the Windows trust store.")
		fmt.Print("Also install CA in the Windows trust store? [y/N] ")

		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(strings.ToLower(answer))

		if answer == "y" || answer == "yes" {
			if err := installWindowsTrust(authority.Cert); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Windows trust store installation failed: %v\n", err)
				fmt.Println("You can manually import the CA certificate on the Windows side:")
				dir, _ := StorageDir()
				fmt.Printf("  File: %s\n", filepath.Join(dir, "ca.crt"))
			} else {
				fmt.Println("CA installed in Windows trust store.")
			}
		}
	}

	// Summary
	dir, _ := StorageDir()
	fmt.Println()
	fmt.Println("Setup complete!")
	fmt.Println()
	fmt.Printf("  CA certificate:   %s\n", filepath.Join(dir, "ca.crt"))
	fmt.Printf("  CA key:           %s\n", filepath.Join(dir, "ca.key"))
	fmt.Printf("  Leaf certificate: %s\n", filepath.Join(dir, "localhost.crt"))
	fmt.Printf("  Leaf key:         %s\n", filepath.Join(dir, "localhost.key"))
	fmt.Println()
	fmt.Println("HTTPS is now ready for *.localhost domains.")

	return nil
}

// Uninstall removes the CA from the trust store and deletes all cert files.
func Uninstall() error {
	dir, err := StorageDir()
	if err != nil {
		return err
	}

	certPath := filepath.Join(dir, "ca.crt")
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		fmt.Println("No CA found. Nothing to uninstall.")
		return nil
	}

	// Load the CA to uninstall from trust store
	authority, err := load(certPath, filepath.Join(dir, "ca.key"))
	if err != nil {
		return fmt.Errorf("loading CA for removal: %w", err)
	}

	fmt.Println("Removing CA from system trust store...")
	if err := truststore.Uninstall(authority.Cert, truststore.WithPrefix("Ennyn")); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: trust store removal may have failed: %v\n", err)
	} else {
		fmt.Println("CA removed from system trust store.")
	}

	// WSL: also remove from Windows trust store
	if isWSL() {
		_ = uninstallWindowsTrust(authority.Cert)
	}

	// Delete all cert files
	files := []string{"ca.crt", "ca.key", "localhost.crt", "localhost.key"}
	for _, f := range files {
		os.Remove(filepath.Join(dir, f))
	}

	// Try to remove the directory (only succeeds if empty)
	os.Remove(dir)

	fmt.Println("All certificates and keys removed.")
	return nil
}

// isWSL checks if we're running under Windows Subsystem for Linux.
func isWSL() bool {
	data, err := os.ReadFile("/proc/version")
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(data)), "microsoft")
}

// installWindowsTrust uses powershell.exe to import the CA cert into the Windows cert store.
func installWindowsTrust(cert *x509.Certificate) error {
	dir, err := StorageDir()
	if err != nil {
		return err
	}
	certPath := filepath.Join(dir, "ca.crt")

	// Convert WSL path to Windows path
	out, err := exec.Command("wslpath", "-w", certPath).Output()
	if err != nil {
		return fmt.Errorf("converting path: %w", err)
	}
	winPath := strings.TrimSpace(string(out))

	cmd := exec.Command("powershell.exe", "-Command",
		fmt.Sprintf(`Import-Certificate -FilePath '%s' -CertStoreLocation Cert:\LocalMachine\Root`, winPath))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// uninstallWindowsTrust removes the CA cert from the Windows cert store.
func uninstallWindowsTrust(cert *x509.Certificate) error {
	thumbprint := fmt.Sprintf("%X", sha1Sum(cert.Raw))
	cmd := exec.Command("powershell.exe", "-Command",
		fmt.Sprintf(`Get-ChildItem Cert:\LocalMachine\Root | Where-Object { $_.Thumbprint -eq '%s' } | Remove-Item`, thumbprint))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
