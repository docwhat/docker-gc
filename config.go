package main

import (
	"os"
	"time"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

// Config stores the configuration from cli flags and environment variables.
type appConfig struct {
	MaxAgeOfImages     time.Duration
	SweeperTime        time.Duration
	DangleSafeDuration time.Duration
	Quiet              bool
	Version            string
}

// NewConfig initializes a Config object from the cli flags and environment variables.
func newAppConfig(args []string) appConfig {
	config := appConfig{Version: version}

	app := kingpin.New("docker-gc", "The missing docker garbage collector.")
	app.Writer(os.Stdout)
	app.HelpFlag.Short('h')
	app.Author("Christian Höltje")
	app.Version(config.Version)
	app.VersionFlag.Short('v')

	app.
		Flag("max-image-age", "How old to allow images to be before deletion. (Env: DOCKER_GC_MAX_IMAGE_AGE)").
		Short('m').
		Default("168h").
		OverrideDefaultFromEnvar("DOCKER_GC_MAX_IMAGE_AGE").
		DurationVar(&config.MaxAgeOfImages)

	app.
		Flag("sweeper-time", "How much time between running checks to delete images. (Env: DOCKER_GC_SWEEPER_TIME)").
		Short('s').
		Default("15m").
		OverrideDefaultFromEnvar("DOCKER_GC_SWEEPER_TIME").
		DurationVar(&config.SweeperTime)

	app.
		Flag("dangle-safe-duration", "How old should a dangle image be before deletion. (Env: DOCKER_GC_DANGLE_SAFE_DURATION)").
		Short('d').
		Default("30m").
		OverrideDefaultFromEnvar("DOCKER_GC_DANGLE_SAFE_DURATION").
		DurationVar(&config.DangleSafeDuration)

	app.
		Flag("quiet", "Don't show any output. (Env: DOCKER_GC_QUIET)").
		Short('q').
		Default("false").
		OverrideDefaultFromEnvar("DOCKER_GC_QUIET").
		BoolVar(&config.Quiet)

	if command, err := app.Parse(args); err != nil {
		app.Usage(nil)
	} else {
		if command != "" {
			app.Usage(nil)
		}
	}

	if config.SweeperTime < (4 * time.Second) {
		app.Fatalf("You must set the sweeper-time to greater or equal than 4s")
	}

	return config
}
