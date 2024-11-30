package met

import "time"

type Config struct {
	DbName          string
	MaxPoolSize     uint64
	MinPoolSize     uint64
	MaxConnIdleTime time.Duration
}
