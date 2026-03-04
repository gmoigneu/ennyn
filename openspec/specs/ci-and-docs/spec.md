### Requirement: CI builds binaries for all target platforms
The GitHub Actions workflow SHALL build binaries for linux/amd64, darwin/arm64, and windows/amd64 on every push to main and on version tags.

#### Scenario: Push to main triggers build
- **WHEN** code is pushed to the main branch
- **THEN** the workflow SHALL build, vet, and test the project, and produce binaries for all 3 platforms

#### Scenario: Tag push uploads artifacts
- **WHEN** a tag matching `v*` is pushed
- **THEN** the workflow SHALL build binaries and upload them as workflow artifacts
