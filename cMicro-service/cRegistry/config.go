package cRegistry

import "time"

type config struct {
	Addrs        []string
	Timeout      time.Duration
	RegistryPath string
	HeartBeat    int64
}

type Option func(config *config)

func WithAddrs(addrs []string) Option {
	return func(config *config) {
		config.Addrs = addrs
	}
}
func WithTimeout(timeout time.Duration) Option {
	return func(config *config) {
		config.Timeout = timeout
	}
}
func WithRegistryPath(registryPath string) Option {
	return func(config *config) {
		config.RegistryPath = registryPath
	}
}
func WithHeartBeat(heartBeat int64) Option {
	return func(config *config) {
		config.HeartBeat = heartBeat
	}
}

const (
	MaxServiceNum          = 10
	MaxSyncServiceInterval = time.Second * 10
)
