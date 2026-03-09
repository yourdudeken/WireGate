---
description: Build and install the wiregate tool locally
---
// turbo-all
1. Update Go dependencies
   ```bash
   go mod tidy
   ```
2. Build the binary
   ```bash
   go build -o wiregate main.go
   ```
3. Global installation (Option A: Symlink)
   ```bash
   sudo ln -sf $(pwd)/wiregate /usr/local/bin/wiregate
   ```
4. Verify the build
   ```bash
   wiregate hub list
   ```
