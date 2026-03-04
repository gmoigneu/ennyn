package serve

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/gmoigneu/ennyn/internal/config"
	"github.com/gmoigneu/ennyn/internal/process"
	"github.com/gmoigneu/ennyn/internal/proxy"
	"github.com/gmoigneu/ennyn/internal/registry"
)

// Result holds the outcome of starting a single service.
type Result struct {
	Host   string
	Port   int
	Status string // "started", "skipped (already running)", or error message
	Err    error
}

// Run loads the config, starts child processes, and runs a single proxy server
// routing all services. It blocks until interrupted.
func Run(configPath string, proxyPort int) int {
	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	if len(cfg.Services) == 0 {
		fmt.Println("No services defined in config.")
		return 0
	}

	// Start child processes and build routes
	var routes []proxy.Route
	var managers []*process.Manager
	var results []Result

	for _, svc := range cfg.Services {
		// Check if already running
		_, alive, err := registry.Get(svc.Host)
		if err != nil {
			results = append(results, Result{Host: svc.Host, Port: svc.Port, Status: err.Error(), Err: err})
			continue
		}
		if alive {
			results = append(results, Result{Host: svc.Host, Port: svc.Port, Status: "skipped (already running)"})
			continue
		}

		exe, args := svc.SplitCommand()
		proc, err := process.Start(exe, args...)
		if err != nil {
			results = append(results, Result{Host: svc.Host, Port: svc.Port, Status: err.Error(), Err: err})
			continue
		}

		managers = append(managers, proc)
		routes = append(routes, proxy.Route{Host: svc.Host, AppPort: svc.Port})

		// Register instance
		inst := registry.Instance{
			PID:       os.Getpid(),
			Host:      svc.Host,
			ProxyPort: proxyPort,
			AppPort:   svc.Port,
			AppCmd:    exe,
			StartedAt: time.Now(),
		}
		if regErr := registry.Register(inst); regErr != nil {
			log.Printf("warning: failed to register %s: %v", svc.Host, regErr)
		}

		results = append(results, Result{Host: svc.Host, Port: svc.Port, Status: "started"})
	}

	// Print summary
	PrintResults(results)

	if len(routes) == 0 {
		fmt.Println("No services to start.")
		return 1
	}

	// Create multi-host proxy handler
	handler := proxy.NewMultiHandler(routes)
	addr := fmt.Sprintf("127.0.0.1:%d", proxyPort)
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		log.Printf("proxy listening on %s", addr)
		for _, r := range routes {
			log.Printf("  %s.localhost:%d → 127.0.0.1:%d", r.Host, proxyPort, r.AppPort)
		}
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	// Register cleanup for all hosts
	defer func() {
		for _, r := range routes {
			registry.Deregister(r.Host)
		}
	}()

	// Log child exits but keep running — only shutdown on signal
	for _, m := range managers {
		go func(m *process.Manager) {
			if err := <-m.Wait(); err != nil {
				log.Printf("child process exited: %v", err)
			}
		}(m)
	}

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	exitCode := 0
	sig := <-quit
	log.Printf("received %s, shutting down...", sig)

	// Shutdown HTTP server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if shutErr := srv.Shutdown(ctx); shutErr != nil {
		log.Printf("HTTP shutdown error: %v", shutErr)
	}

	// Stop all child processes
	for _, m := range managers {
		m.Stop()
	}

	log.Println("stopped")
	return exitCode
}

// PrintResults prints a summary table of service startup results.
func PrintResults(results []Result) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "HOST\tPORT\tSTATUS")
	for _, r := range results {
		fmt.Fprintf(w, "%s\t%d\t%s\n", r.Host, r.Port, r.Status)
	}
	w.Flush()
}

// HasErrors returns true if any result has an error.
func HasErrors(results []Result) bool {
	for _, r := range results {
		if r.Err != nil {
			return true
		}
	}
	return false
}
