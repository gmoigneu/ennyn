package registry

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

// Instance represents a running ennyn server instance.
type Instance struct {
	PID       int       `json:"pid"`
	Host      string    `json:"host"`
	ProxyPort int       `json:"proxyPort"`
	AppPort   int       `json:"appPort"`
	AppCmd    string    `json:"appCmd"`
	StartedAt time.Time `json:"startedAt"`
}

// RunDir returns the directory where instance state files are stored.
func RunDir() (string, error) {
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("cannot determine home directory: %w", err)
		}
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(base, "ennyn", "run"), nil
}

// Register writes a state file for the given instance.
func Register(inst Instance) error {
	dir, err := RunDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating run directory: %w", err)
	}

	data, err := json.MarshalIndent(inst, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling instance: %w", err)
	}

	path := filepath.Join(dir, inst.Host+".json")
	return os.WriteFile(path, data, 0644)
}

// Deregister removes the state file for the given hostname.
func Deregister(host string) error {
	dir, err := RunDir()
	if err != nil {
		return err
	}
	path := filepath.Join(dir, host+".json")
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// IsAlive checks whether a process with the given PID is running.
func IsAlive(pid int) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	err = proc.Signal(syscall.Signal(0))
	return err == nil
}

// List returns all live instances, removing stale state files.
func List() ([]Instance, error) {
	dir, err := RunDir()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var live []Instance
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}

		inst, err := loadInstance(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}

		if IsAlive(inst.PID) {
			live = append(live, *inst)
		} else {
			// Stale — clean up
			os.Remove(filepath.Join(dir, e.Name()))
		}
	}
	return live, nil
}

// Get loads a single instance by hostname. Returns the instance and true if alive,
// or nil and false if not found or stale (stale files are cleaned up).
func Get(host string) (*Instance, bool, error) {
	dir, err := RunDir()
	if err != nil {
		return nil, false, err
	}

	path := filepath.Join(dir, host+".json")
	inst, err := loadInstance(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, false, nil
		}
		return nil, false, err
	}

	if IsAlive(inst.PID) {
		return inst, true, nil
	}

	// Stale — clean up
	os.Remove(path)
	return nil, false, nil
}

// CheckDuplicate returns an error if the host is already registered with a live process.
// If the state file is stale, it is cleaned up and nil is returned.
func CheckDuplicate(host string) error {
	inst, alive, err := Get(host)
	if err != nil {
		return err
	}
	if alive {
		return fmt.Errorf("host %q is already running (PID %d)", inst.Host, inst.PID)
	}
	return nil
}

func loadInstance(path string) (*Instance, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var inst Instance
	if err := json.Unmarshal(data, &inst); err != nil {
		return nil, err
	}
	return &inst, nil
}
