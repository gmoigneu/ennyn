## Context

Ennyn currently has a minimal HTTP server (`main.go`) with no subcommands. To serve HTTPS with trusted certificates, we need a local CA. mkcert pioneered this approach; the `smallstep/truststore` library (extracted from mkcert, used by Caddy) handles the hard part of platform-specific trust store installation.

## Goals / Non-Goals

**Goals:**
- Generate a local root CA (ECDSA P-256) and persist it to `~/.config/ennyn/`
- Install the CA certificate into the OS trust store via `smallstep/truststore`
- Detect WSL and offer to install in the Windows trust store alongside Linux
- Provide `ennyn trust` (install) and `ennyn trust --uninstall` (remove)
- Print clear post-install instructions

**Non-Goals:**
- Firefox NSS trust store support (deferred)
- Java trust store support (deferred)
- Automatic re-trust on CA expiry (manual re-run of `ennyn trust` is fine)
- HTTPS listener changes (separate change)
- Custom hostname resolution — `.localhost` subdomains resolve to `127.0.0.1` automatically per RFC 6761, no `/etc/hosts` changes needed

## Decisions

**ECDSA P-256 over RSA for CA key**
Faster key generation, smaller certificates, sufficient security for a local dev CA. mkcert defaults to RSA 3072 for compatibility, but since Ennyn only targets modern systems, P-256 is fine. Caddy also uses ECDSA.

**`smallstep/truststore` for trust store installation**
Handles macOS Keychain, Linux ca-certificates, and Windows cert store. Battle-tested in Caddy. Avoids reimplementing platform-specific trust logic. BSD-compatible license. Alternative: vendoring mkcert's trust store code directly — more control but more maintenance.

**CA validity: 10 years**
Long enough that developers rarely need to re-run setup. mkcert uses the same duration. The CA never leaves the local machine so long validity is not a security concern.

**Storage at `~/.config/ennyn/ca/`**
Follows XDG Base Directory spec. Respects `$XDG_CONFIG_HOME` if set. Files: `ca.crt` (PEM-encoded certificate), `ca.key` (PEM-encoded private key). Directory created with 0700 permissions. Key file written with 0600 permissions.

**Subcommand architecture with Go `flag` subcommands**
Refactor `main.go` to use manual subcommand dispatch (`os.Args[1]` matching) rather than pulling in cobra/urfave. Keep it simple — we only have two commands: the default server start and `trust`. If the command set grows, we can adopt a framework later.

**WSL detection via `/proc/version`**
Check if `/proc/version` contains "microsoft" (case-insensitive). This is the standard detection method. When WSL is detected, use `powershell.exe` to invoke certificate trust in the Windows host. Offer this as a prompt — don't force it.

## Risks / Trade-offs

**`smallstep/truststore` is pre-v1.0** → It's been stable for years and is used by Caddy in production. Pin the dependency version.

**Sudo prompts may confuse users** → Print a clear explanation before invoking sudo: "Installing Ennyn CA into your system trust store. This requires administrator access."

**WSL cross-environment trust is fragile** → The `powershell.exe` approach may not work in all WSL configurations. Make it optional with a clear error message if it fails. Users can always manually import the `.crt` file on the Windows side.

**CA key stored on disk** → Acceptable for a local dev tool. File permissions (0600) prevent other users from reading it. Document that `~/.config/ennyn/ca/ca.key` should not be shared.

**All hostnames use `.localhost` TLD (RFC 6761)**
All Ennyn-managed services use `<app>.localhost` hostnames. This has major benefits:
- Resolves to `127.0.0.1` automatically on all platforms — no `/etc/hosts` editing
- Browsers treat `*.localhost` as a Secure Context (Service Workers, etc. work even over HTTP)
- A single wildcard X.509 certificate (`*.localhost`) covers all apps
- X.509 wildcards are single-level, which is exactly our use case (`myapp.localhost` matches, nested subdomains like `api.myapp.localhost` would not — but we don't need them)

**Wildcard leaf certificate generated during `ennyn trust`**
Rather than generating leaf certs on-demand per hostname, we generate a single `*.localhost` cert at setup time. SANs include: `*.localhost`, `localhost`, `127.0.0.1`, `::1`. This cert is stored alongside the CA at `~/.config/ennyn/ca/localhost.crt` and `localhost.key`. Valid for 825 days (within Apple's limit for trusted certificates). The proxy will use this single cert for all TLS connections.
