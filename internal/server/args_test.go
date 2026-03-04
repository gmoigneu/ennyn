package server

import "testing"

func TestValidateHost(t *testing.T) {
	tests := []struct {
		host string
		want bool
	}{
		{"myapp", true},
		{"my-app", true},
		{"a", true},
		{"app123", true},
		{"123", true},
		{"a-b-c", true},
		{"-app", false},
		{"app-", false},
		{"My-App", false},
		{"my_app", false},
		{"my.app", false},
		{"", false},
		{"-", false},
	}

	for _, tt := range tests {
		t.Run(tt.host, func(t *testing.T) {
			if got := ValidateHost(tt.host); got != tt.want {
				t.Errorf("ValidateHost(%q) = %v, want %v", tt.host, got, tt.want)
			}
		})
	}
}

func TestValidatePort(t *testing.T) {
	tests := []struct {
		input   string
		want    int
		wantErr bool
	}{
		{"80", 80, false},
		{"3000", 3000, false},
		{"1", 1, false},
		{"65535", 65535, false},
		{"0", 0, true},
		{"65536", 0, true},
		{"99999", 0, true},
		{"-1", 0, true},
		{"abc", 0, true},
		{"", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ValidatePort(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePort(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidatePort(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    *Config
		wantErr bool
	}{
		{
			name: "basic invocation",
			args: []string{"myapp", "3000", "node", "server.js"},
			want: &Config{Host: "myapp", AppPort: 3000, AppCmd: "node", AppArgs: []string{"server.js"}, ProxyPort: 7890},
		},
		{
			name: "custom proxy port",
			args: []string{"-p", "8000", "myapp", "3000", "node", "server.js"},
			want: &Config{Host: "myapp", AppPort: 3000, AppCmd: "node", AppArgs: []string{"server.js"}, ProxyPort: 8000},
		},
		{
			name: "long port flag",
			args: []string{"-port", "9000", "api", "8080", "python", "-m", "flask", "run"},
			want: &Config{Host: "api", AppPort: 8080, AppCmd: "python", AppArgs: []string{"-m", "flask", "run"}, ProxyPort: 9000},
		},
		{
			name: "command with no extra args",
			args: []string{"web", "4000", "my-server"},
			want: &Config{Host: "web", AppPort: 4000, AppCmd: "my-server", AppArgs: []string{}, ProxyPort: 7890},
		},
		{
			name:    "missing arguments",
			args:    []string{"myapp"},
			wantErr: true,
		},
		{
			name:    "missing command",
			args:    []string{"myapp", "3000"},
			wantErr: true,
		},
		{
			name:    "invalid host",
			args:    []string{"MyApp", "3000", "node", "server.js"},
			wantErr: true,
		},
		{
			name:    "invalid port",
			args:    []string{"myapp", "notanumber", "node", "server.js"},
			wantErr: true,
		},
		{
			name:    "port out of range",
			args:    []string{"myapp", "99999", "node", "server.js"},
			wantErr: true,
		},
		{
			name:    "no arguments",
			args:    []string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseArgs(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseArgs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want == nil {
				return
			}
			if got.Host != tt.want.Host {
				t.Errorf("Host = %q, want %q", got.Host, tt.want.Host)
			}
			if got.AppPort != tt.want.AppPort {
				t.Errorf("AppPort = %d, want %d", got.AppPort, tt.want.AppPort)
			}
			if got.AppCmd != tt.want.AppCmd {
				t.Errorf("AppCmd = %q, want %q", got.AppCmd, tt.want.AppCmd)
			}
			if len(got.AppArgs) != len(tt.want.AppArgs) {
				t.Errorf("AppArgs len = %d, want %d", len(got.AppArgs), len(tt.want.AppArgs))
			} else {
				for i := range got.AppArgs {
					if got.AppArgs[i] != tt.want.AppArgs[i] {
						t.Errorf("AppArgs[%d] = %q, want %q", i, got.AppArgs[i], tt.want.AppArgs[i])
					}
				}
			}
			if got.ProxyPort != tt.want.ProxyPort {
				t.Errorf("ProxyPort = %d, want %d", got.ProxyPort, tt.want.ProxyPort)
			}
		})
	}
}
