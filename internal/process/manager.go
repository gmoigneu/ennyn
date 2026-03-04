package process

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
)

// Manager manages a child process lifecycle.
type Manager struct {
	cmd  *exec.Cmd
	done chan error
}

// Start launches the command as a child process in its own process group,
// inheriting stdout and stderr.
func Start(name string, args ...string) (*Manager, error) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("starting %s: %w", name, err)
	}

	m := &Manager{
		cmd:  cmd,
		done: make(chan error, 1),
	}

	go func() {
		m.done <- cmd.Wait()
	}()

	return m, nil
}

// Wait returns a channel that receives the process exit error (nil if clean exit).
func (m *Manager) Wait() <-chan error {
	return m.done
}

// ExitCode returns the exit code from the child process error.
// Returns 0 for nil error, 1 if the code can't be determined.
func ExitCode(err error) int {
	if err == nil {
		return 0
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode()
	}
	return 1
}

// Stop sends SIGTERM to the child process group, waits up to 5 seconds,
// then sends SIGKILL if still running.
func (m *Manager) Stop() error {
	// Send SIGTERM to the process group
	pgid, err := syscall.Getpgid(m.cmd.Process.Pid)
	if err != nil {
		// Process may have already exited
		return nil
	}
	syscall.Kill(-pgid, syscall.SIGTERM)

	select {
	case <-m.done:
		return nil
	case <-time.After(5 * time.Second):
		// Force kill the process group
		syscall.Kill(-pgid, syscall.SIGKILL)
		<-m.done
		return nil
	}
}
