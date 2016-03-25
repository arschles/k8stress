package main

// Config is the envconfig compatible configuration struct.
// See github.com/kelseyhightower/envconfig for more information
type Config struct {
	NumGoroutines int    `envconfig:"NUM_GOROUTINES" default:"10000"`
	Namespace     string `envconfig:"NAMESPACE" default:"k8stress"`
}
