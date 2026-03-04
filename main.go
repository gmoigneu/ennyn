package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/gmoigneu/ennyn/internal/ca"
	"github.com/gmoigneu/ennyn/internal/config"
	"github.com/gmoigneu/ennyn/internal/process"
	"github.com/gmoigneu/ennyn/internal/proxy"
	"github.com/gmoigneu/ennyn/internal/registry"
	"github.com/gmoigneu/ennyn/internal/serve"
	"github.com/gmoigneu/ennyn/internal/server"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "server":
		os.Exit(runServer(os.Args[2:]))
	case "serve":
		runServe(os.Args[2:])
	case "trust":
		runTrust(os.Args[2:])
	case "list":
		runList()
	case "stop":
		runStop(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr, "  ennyn serve [flags]                                         Start all services from ennyn.yml")
	fmt.Fprintln(os.Stderr, "  ennyn stop [<host>]                                         Stop a service, or all services from ennyn.yml")
	fmt.Fprintln(os.Stderr, "  ennyn server [flags] <host> <app-port> <app-cmd> [app-args...]  Start a single service")
	fmt.Fprintln(os.Stderr, "  ennyn list                                                  List running instances")
	fmt.Fprintln(os.Stderr, "  ennyn trust [--uninstall]                                   Manage local CA and certificates")
}

func runServer(args []string) int {
	cfg, err := server.ParseArgs(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	// Check for duplicate host
	if err := registry.CheckDuplicate(cfg.Host); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	// Start child process
	proc, err := process.Start(cfg.AppCmd, cfg.AppArgs...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	// Create proxy handler
	handler := proxy.NewHandler(cfg.Host, cfg.AppPort)

	addr := fmt.Sprintf("127.0.0.1:%d", cfg.ProxyPort)
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	// Start HTTP listener
	go func() {
		log.Printf("proxy listening on %s → %s.localhost:%d → 127.0.0.1:%d",
			addr, cfg.Host, cfg.ProxyPort, cfg.AppPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	// Register instance
	inst := registry.Instance{
		PID:       os.Getpid(),
		Host:      cfg.Host,
		ProxyPort: cfg.ProxyPort,
		AppPort:   cfg.AppPort,
		AppCmd:    cfg.AppCmd,
		StartedAt: time.Now(),
	}
	if err := registry.Register(inst); err != nil {
		log.Printf("warning: failed to register instance: %v", err)
	}
	defer registry.Deregister(cfg.Host)

	// Wait for shutdown signal or child exit
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	exitCode := 0
	select {
	case sig := <-quit:
		log.Printf("received %s, shutting down...", sig)
	case err := <-proc.Wait():
		code := process.ExitCode(err)
		if code != 0 {
			log.Printf("app exited with code %d, shutting down...", code)
		} else {
			log.Println("app exited, shutting down...")
		}
		exitCode = code
	}

	// Shutdown HTTP server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("HTTP shutdown error: %v", err)
	}

	// Stop child process (if still running)
	proc.Stop()

	log.Println("stopped")
	return exitCode
}

func runServe(args []string) {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	var proxyPort int
	fs.IntVar(&proxyPort, "port", 7890, "proxy listen port")
	fs.IntVar(&proxyPort, "p", 7890, "proxy listen port (shorthand)")
	fs.Parse(args)

	os.Exit(serve.Run("ennyn.yml", proxyPort))
}

func runTrust(args []string) {
	fs := flag.NewFlagSet("trust", flag.ExitOnError)
	uninstall := fs.Bool("uninstall", false, "remove CA and certificates from trust store")
	fs.Parse(args)

	if *uninstall {
		if err := ca.Uninstall(); err != nil {
			log.Fatalf("uninstall failed: %v", err)
		}
		return
	}

	if err := ca.Install(); err != nil {
		log.Fatalf("trust setup failed: %v", err)
	}
}

func runList() {
	instances, err := registry.List()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if len(instances) == 0 {
		fmt.Println("No running instances.")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "HOST\tPID\tPROXY PORT\tAPP PORT\tCOMMAND\tUPTIME")
	for _, inst := range instances {
		uptime := time.Since(inst.StartedAt).Truncate(time.Second)
		fmt.Fprintf(w, "%s\t%d\t%d\t%d\t%s\t%s\n",
			inst.Host, inst.PID, inst.ProxyPort, inst.AppPort, inst.AppCmd, uptime)
	}
	w.Flush()
}

func runStop(args []string) {
	if len(args) == 0 {
		// Try loading config for stop-all mode
		cfg, err := config.Load("ennyn.yml")
		if err != nil {
			fmt.Fprintln(os.Stderr, "Usage: ennyn stop [<host>]")
			fmt.Fprintln(os.Stderr, "  Without arguments, reads ennyn.yml to stop all services.")
			os.Exit(1)
		}
		stopAll(cfg)
		return
	}

	host := args[0]
	stopOne(host)
}

func stopAll(cfg *config.Config) {
	hadError := false
	for _, svc := range cfg.Services {
		inst, alive, err := registry.Get(svc.Host)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: error: %v\n", svc.Host, err)
			hadError = true
			continue
		}
		if !alive {
			fmt.Printf("%s: not running\n", svc.Host)
			continue
		}

		proc, err := os.FindProcess(inst.PID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: error finding process %d: %v\n", svc.Host, inst.PID, err)
			hadError = true
			continue
		}
		if err := proc.Signal(syscall.SIGTERM); err != nil {
			fmt.Fprintf(os.Stderr, "%s: error sending SIGTERM to PID %d: %v\n", svc.Host, inst.PID, err)
			hadError = true
			continue
		}
		fmt.Printf("%s: sent SIGTERM (PID %d)\n", svc.Host, inst.PID)
	}
	if hadError {
		os.Exit(1)
	}
}

func stopOne(host string) {
	inst, alive, err := registry.Get(host)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if !alive {
		fmt.Fprintf(os.Stderr, "no running instance for %q\n", host)
		os.Exit(1)
	}

	proc, err := os.FindProcess(inst.PID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error finding process %d: %v\n", inst.PID, err)
		os.Exit(1)
	}

	if err := proc.Signal(syscall.SIGTERM); err != nil {
		fmt.Fprintf(os.Stderr, "error sending SIGTERM to PID %d: %v\n", inst.PID, err)
		os.Exit(1)
	}

	fmt.Printf("Sent SIGTERM to %s (PID %d)\n", host, inst.PID)
}
