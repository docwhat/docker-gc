package config

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

// Config stores the configuration from cli flags and environment variables.
type Config struct {
	MaxAgeOfImages time.Duration
	SweeperTime    time.Duration
	Verbosity      uint8
	Version        string
}

// NewConfig initializes a Config object from the cli flags and environment variables.
func NewConfig() Config {
	config := Config{Version: version}
	kingpin.CommandLine.Writer(os.Stdout)
	kingpin.HelpFlag.Short('h')
	kingpin.CommandLine.Help = "The missing docker garbage collector."
	kingpin.CommandLine.Author("Christian Höltje")
	kingpin.Version(config.Version)

	kingpin.Flag("max-image-age", "How old to allow images to be before deletion. (Env: DOCKER_GC_MAX_IMAGE_AGE)").Short('m').Default("168h").OverrideDefaultFromEnvar("DOCKER_GC_MAX_IMAGE_AGE").DurationVar(&config.MaxAgeOfImages)
	kingpin.Flag("sweeper-time", "How much time between running checks to delete images. (Env: DOCKER_GC_SWEEPER_TIME)").Short('s').Default("15m").OverrideDefaultFromEnvar("DOCKER_GC_SWEEPER_TIME").DurationVar(&config.SweeperTime)
	kingpin.Flag("verbosity", "How much logging to stderr. 0 = none. 9 = maximal (Env: DOCKER_GC_VERBOSITY)").Short('v').Default("1").OverrideDefaultFromEnvar("DOCKER_GC_VERBOSITY").Uint8Var(&config.Verbosity)

	kingpin.Parse()

	if config.Verbosity == 0 {
		log.SetOutput(ioutil.Discard)
	}
	return config
}
