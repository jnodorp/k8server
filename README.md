# k8server

k8server provides some boilerplate code when building [Kubernetes](https://kubernetes.io/)-friendly servers. This
includes:

- [Prometheus](https://prometheus.io/) Metrics
- [Liveness](https://kubernetes.io/docs/concepts/configuration/liveness-readiness-startup-probes/#liveness-probe) and
  [readiness](https://kubernetes.io/docs/concepts/configuration/liveness-readiness-startup-probes/#readiness-probe)
  probes
- Graceful shutdown

## Usage

```go
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintln(w, "hello, world!")
	})

    k8server.Run(context.Background(), mux)
}
```

## Defaults

By default, k8server is serving the provided mux on port 8080. The `/livez`, `/readyz`, and `/metrics` endpoints are
served on port 8081.

## Logging

k8server uses [slog](https://pkg.go.dev/log/slog).

To enable JSON logging, add

```go
slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})))
```

to your code.
