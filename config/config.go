package config

type (
	Config struct {
		// How many workers will run in parallel.
		Parallel uint64
		// URLs to fetch and generate MD5 hash for
		URLs []string
	}
)
