package main

import (
	"fmt"
)

// Config is the envconfig compatible configuration struct.
// See github.com/kelseyhightower/envconfig for more information
type Config struct {
	NumGoroutines int    `envconfig:"NUM_GOROUTINES" default:"10000"`
	TimeSec       int    `envconfig:"TIME_SEC" default:"86400"` // 1 day
	Namespace     string `envconfig:"NAMESPACE" default:"k8stress"`
}

func (c Config) String() string {
	return fmt.Sprintf("Num Goroutines = %d, Namespace = %s, runtime = %ds", c.NumGoroutines, c.Namespace, c.TimeSec)
}
