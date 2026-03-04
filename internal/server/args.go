package server

import (
	"flag"
	"fmt"
	"regexp"
	"strconv"
)

var hostPattern = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]*[a-z0-9])?$`)

// Config holds the parsed arguments for the server subcommand.
type Config struct {
	Host      string
	AppPort   int
	AppCmd    string
	AppArgs   []string
	ProxyPort int
}

// ParseArgs parses the arguments following "ennyn server".
// Expected: [flags] <host> <app-port> <app-cmd> [app-args...]
func ParseArgs(args []string) (*Config, error) {
	fs := flag.NewFlagSet("server", flag.ContinueOnError)

	var proxyPort int
	fs.IntVar(&proxyPort, "port", 7890, "proxy listen port")
	fs.IntVar(&proxyPort, "p", 7890, "proxy listen port (shorthand)")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	positional := fs.Args()
	if len(positional) < 3 {
		return nil, fmt.Errorf("usage: ennyn server [flags] <host> <app-port> <app-cmd> [app-args...]\n\nRequires at least 3 arguments: host, app-port, and command")
	}

	host := positional[0]
	if !ValidateHost(host) {
		return nil, fmt.Errorf("invalid host %q: must be lowercase alphanumeric and hyphens, cannot start or end with a hyphen", host)
	}

	appPort, err := ValidatePort(positional[1])
	if err != nil {
		return nil, err
	}

	return &Config{
		Host:      host,
		AppPort:   appPort,
		AppCmd:    positional[2],
		AppArgs:   positional[3:],
		ProxyPort: proxyPort,
	}, nil
}

// ValidateHost checks that the host is lowercase alphanumeric with hyphens,
// not starting or ending with a hyphen.
func ValidateHost(host string) bool {
	return hostPattern.MatchString(host)
}

// ValidatePort parses and validates a port number string (1–65535).
func ValidatePort(s string) (int, error) {
	port, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid port %q: must be a number", s)
	}
	if port < 1 || port > 65535 {
		return 0, fmt.Errorf("port %d out of range: must be 1–65535", port)
	}
	return port, nil
}
