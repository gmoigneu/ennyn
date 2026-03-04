package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "ennyn.yml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoadValidConfig(t *testing.T) {
	path := writeConfig(t, `
services:
  - host: myapp
    port: 3000
    command: npm run dev
  - host: api
    port: 8080
    command: go run ./cmd/api
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Services) != 2 {
		t.Fatalf("expected 2 services, got %d", len(cfg.Services))
	}
	if cfg.Services[0].Host != "myapp" || cfg.Services[0].Port != 3000 {
		t.Errorf("service 0: got host=%q port=%d", cfg.Services[0].Host, cfg.Services[0].Port)
	}
	if cfg.Services[1].Host != "api" || cfg.Services[1].Port != 8080 {
		t.Errorf("service 1: got host=%q port=%d", cfg.Services[1].Host, cfg.Services[1].Port)
	}
}

func TestLoadEmptyServices(t *testing.T) {
	path := writeConfig(t, "services: []\n")
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Services) != 0 {
		t.Fatalf("expected 0 services, got %d", len(cfg.Services))
	}
}

func TestLoadInvalidHost(t *testing.T) {
	path := writeConfig(t, `
services:
  - host: MyApp
    port: 3000
    command: npm run dev
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for invalid host")
	}
}

func TestLoadInvalidPort(t *testing.T) {
	path := writeConfig(t, `
services:
  - host: myapp
    port: 99999
    command: npm run dev
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for port out of range")
	}
}

func TestLoadMissingCommand(t *testing.T) {
	path := writeConfig(t, `
services:
  - host: myapp
    port: 3000
    command: ""
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing command")
	}
}

func TestLoadDuplicateHosts(t *testing.T) {
	path := writeConfig(t, `
services:
  - host: myapp
    port: 3000
    command: npm run dev
  - host: myapp
    port: 8080
    command: go run .
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for duplicate hosts")
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := Load("/nonexistent/ennyn.yml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestSplitCommand(t *testing.T) {
	svc := Service{Command: "python -m flask run"}
	exe, args := svc.SplitCommand()
	if exe != "python" {
		t.Errorf("expected exe=python, got %q", exe)
	}
	if len(args) != 3 || args[0] != "-m" || args[1] != "flask" || args[2] != "run" {
		t.Errorf("expected args=[-m flask run], got %v", args)
	}
}

func TestSplitCommandSingle(t *testing.T) {
	svc := Service{Command: "node"}
	exe, args := svc.SplitCommand()
	if exe != "node" {
		t.Errorf("expected exe=node, got %q", exe)
	}
	if len(args) != 0 {
		t.Errorf("expected no args, got %v", args)
	}
}
