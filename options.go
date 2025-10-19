package k8server

import "time"

type options struct {
	// Port used to listen on the provided server mux.
	Port int

	// ManagementPort used to server metrics and probes.
	ManagementPort int

	// Timeout for graceful shutdown.
	Timeout time.Duration
}

func defaults() *options {
	return &options{
		Port:           8080,
		ManagementPort: 8081,
		// The default timeout is half the default termination grace period of 30 seconds used by Kubernetes. This gives
		// us plenty of headroom.
		Timeout: 15 * time.Second,
	}
}

type Option interface {
	apply(*options)
}

type portOption int

func (o portOption) apply(opts *options) {
	opts.Port = int(o)
}

type managementPortOption int

func (o managementPortOption) apply(opts *options) {
	opts.ManagementPort = int(o)
}

func WithManagementPort(managementPort int) Option {
	return managementPortOption(managementPort)
}

type timeoutOption time.Duration

func (o timeoutOption) apply(opts *options) {
	opts.Timeout = time.Duration(o)
}

// WithTimeout when waiting for requests to complete after receiving an interrupt or termination signal.
func WithTimeout(timeout time.Duration) Option {
	return timeoutOption(timeout)
}
