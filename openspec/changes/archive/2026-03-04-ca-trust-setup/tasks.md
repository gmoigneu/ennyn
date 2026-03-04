## 1. Subcommand Architecture

- [x] 1.1 Refactor `main.go` to dispatch subcommands: `trust` runs CA workflow, no subcommand (or `-p`/`-port`) starts the HTTP server, unknown subcommand prints usage and exits non-zero
- [x] 1.2 Move HTTP server logic into a `runServer()` function called from the default path

## 2. CA Generation

- [x] 2.1 Create `internal/ca/` package with a function to resolve the CA storage directory (`$XDG_CONFIG_HOME/ennyn/ca/` or `~/.config/ennyn/ca/`)
- [x] 2.2 Implement CA key pair generation (ECDSA P-256) and self-signed root certificate creation (10-year validity, CN "Ennyn Local CA", IsCA true)
- [x] 2.3 Implement CA persistence: write `ca.crt` and `ca.key` as PEM files, directory 0700, key file 0600
- [x] 2.4 Implement CA loading: if `ca.crt` and `ca.key` already exist, load and return them without regenerating

## 3. Wildcard Leaf Certificate

- [x] 3.1 Implement leaf certificate generation: ECDSA P-256 key, signed by CA, 825-day validity, SANs: `*.localhost`, `localhost`, `127.0.0.1`, `::1`
- [x] 3.2 Store leaf cert as `localhost.crt` and `localhost.key` in the CA storage directory with 0600 key permissions
- [x] 3.3 On CA regeneration, also regenerate the leaf cert; on reuse of existing CA, reuse existing leaf cert if present

## 4. Trust Store Installation

- [x] 4.1 Add `smallstep/truststore` dependency via `go get`
- [x] 4.2 Implement `Install()` function that calls `truststore.Install()` with the CA certificate
- [x] 4.3 Implement `Uninstall()` function that calls `truststore.Uninstall()` and deletes the CA and leaf cert files from disk
- [x] 4.4 Wire up `ennyn trust` to generate/load CA, generate/load leaf cert, then call Install; `ennyn trust --uninstall` calls Uninstall

## 5. WSL Detection

- [x] 5.1 Implement WSL detection by reading `/proc/version` for "microsoft" (case-insensitive)
- [x] 5.2 When WSL is detected, prompt the user (stdin y/n) to also install the CA in the Windows trust store
- [x] 5.3 If user accepts, invoke `powershell.exe` to import the certificate into the Windows cert store; handle and report errors gracefully

## 6. User Output

- [x] 6.1 Print pre-install message explaining what will happen and that sudo may be required
- [x] 6.2 Print post-install summary: CA storage path, which trust stores were updated, leaf cert location, and that HTTPS is ready for `*.localhost`
- [x] 6.3 Print appropriate messages for uninstall and for the "no CA found" uninstall case

## 7. Verification

- [x] 7.1 `go build ./...` compiles without errors
- [x] 7.2 `go vet ./...` passes
- [x] 7.3 `ennyn` with no args still starts the HTTP server on port 7890
- [x] 7.4 `ennyn trust` generates CA files, leaf cert, and installs CA into trust store
- [x] 7.5 `ennyn trust` again reuses existing CA and leaf cert without regenerating
- [x] 7.6 `ennyn trust --uninstall` removes trust and deletes all cert files
- [x] 7.7 `ennyn foobar` prints usage and exits non-zero
- [x] 7.8 Leaf cert has correct SANs: `*.localhost`, `localhost`, `127.0.0.1`, `::1`
