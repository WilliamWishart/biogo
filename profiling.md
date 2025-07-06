# Profiling Guide: CPU and Memory Profiling in Go

This guide explains how to use the built-in Go profiler to analyze CPU and memory usage in your project.

---

## 1. Enabling Profiling

Profiling is now configurable! You can enable CPU and memory profiling in two ways:

- **Environment variable:**
  ```sh
  BIOGO_PROFILE=1 go run main.go
  ```
- **Command-line flag:**
  ```sh
  go run main.go --profile
  ```

Profiling is only enabled if either the environment variable or the flag is set. If not set, the program runs without profiling overhead.

---

## 2. Generating Profiles

When profiling is enabled, your program generates `cpu.prof` and `mem.prof` files in the project root.

- **CPU Profile:** `cpu.prof`
- **Memory Profile:** `mem.prof`

---

## 3. Analyzing CPU Profile

### Basic Usage

1. **Run your program** to generate `cpu.prof` (see above).
2. **Analyze with pprof:**
   ```sh
   go tool pprof cpu.prof
   ```
3. **At the pprof prompt**, use commands like:
   - `top` — Show top functions by CPU usage
   - `list <funcname>` — Show source-level breakdown for a function
   - `web` — Open a graphical visualization (requires Graphviz)

### With Binary (for source context)

If you want to see annotated source, specify the binary:
```sh
go build -o biogo
# Then:
go tool pprof biogo cpu.prof
```

### Web UI

For an interactive browser UI:
```sh
go tool pprof -http=:8080 cpu.prof
```
Then open [http://localhost:8080](http://localhost:8080) in your browser.

---

## 4. Analyzing Memory Profile

1. **Run your program** to generate `mem.prof` (see above).
2. **Analyze with pprof:**
   ```sh
   go tool pprof mem.prof
   ```
3. **Use the same commands** as above (`top`, `list`, `web`, etc.).

---

## 5. Useful pprof Commands

- `top` — Show top resource consumers
- `list <funcname>` — Show annotated source for a function
- `web` — Generate a graph (requires Graphviz)
- `help` — List all commands

---

## 6. Tips

- Profiling adds some overhead and will slow your program down (CPU profiling typically adds 5–20% overhead).
- Always run profiling on a representative workload.
- Use `runtime.GC()` before writing heap profiles for up-to-date stats.
- For more info, see the [Go pprof documentation](https://golang.org/pkg/runtime/pprof/).

---

**Happy profiling!**
