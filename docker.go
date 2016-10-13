// +build go1.7

package main

import (
	"fmt"
	"strings"
	"time"

	clientLib "github.com/fsouza/go-dockerclient"
)

// ImageNameScanner is a function used when going through image names.
type ImageNameScanner func(tag string)

// SawImageHandler is a method for dealing with an image being seen at a specific time.
type SawImageHandler func(tag string, when time.Time)

// Docker is a wrapper around a subset of docker commands need by
// docker-gc.
type Docker struct {
	client *clientLib.Client
	logger Logger
}

// NewDocker initializes and returns a new Docker instance.
func NewDocker(logger Logger) *Docker {
	client, err := clientLib.NewClientFromEnv()

	if err != nil {
		logger.Fatal("Unable to create docker client: %s", err)
	}

	if err := client.Ping(); err != nil {
		logger.Fatal("Cannot connect to the Docker daemon. Is the docker daemon running on this host?\n\t%s", err)
	}

	return &Docker{client: client, logger: logger}
}

// NormalizeImageName returns a normalized image name. The closest thing to a canonical name.
func NormalizeImageName(repoTag string) string {
	repo, tag := clientLib.ParseRepositoryTag(repoTag)
	if tag == "" {
		tag = "latest"
	}
	return fmt.Sprintf("%s:%s", repo, tag)
}

// eventToTime returns a time object from an event.
//
// Events has weird rules for the time due to different versions of
// the docker API.
func eventToTime(event *clientLib.APIEvents) time.Time {
	if event.TimeNano != 0 {
		return time.Unix(0, event.TimeNano)
	} else if event.Time != 0 {
		return time.Unix(event.Time, 0)
	}
	return time.Now()
}

// RemoveImage deletes an image from docker.
//
// Returns true if it succeeds.
func (d *Docker) RemoveImage(tag string) bool {
	if err := d.client.RemoveImage(tag); err != nil {
		return false
	}
	return true
}

// HasImage returns `true` if an image exists in docker.
func (d *Docker) HasImage(tag string) bool {
	var imagesFound int

	if images, err := d.client.ListImages(clientLib.ListImagesOptions{Filter: tag}); err != nil {
		d.logger.Fatal("Failed when asking about %v: %s", tag, err)
	} else {
		imagesFound = len(images)
	}

	return (imagesFound != 0)
}

// HasContainerWithImage returns true if the image is in use by a container.
func (d *Docker) HasContainerWithImage(tag string) bool {
	var containersFound int

	opts := clientLib.ListContainersOptions{Filters: make(map[string][]string)}
	opts.Filters["ancestor"] = []string{tag}
	opts.All = true
	if containers, err := d.client.ListContainers(opts); err != nil {
		d.logger.Fatal("Failed when asking about containers with image %v: %s", tag, err)
	} else {
		containersFound = len(containers)
	}

	return (containersFound != 0)
}

// ScanDanglingImages returns all images that are dangling.
func (d *Docker) ScanDanglingImages(safetyDuration time.Duration, f ImageNameScanner) {
	opts := clientLib.ListImagesOptions{Filters: make(map[string][]string)}
	opts.Filters["dangling"] = []string{"true"}
	if images, err := d.client.ListImages(opts); err != nil {
		d.logger.Fatal("Failed to scan for dangling images: %s", err)
	} else {
		safeTime := time.Now().Add(safetyDuration)
		for _, image := range images {
			if safeTime.Before(time.Unix(image.Created, 0)) {
				continue
			}
			if d.HasContainerWithImage(image.ID) {
				continue
			}
			f(image.ID)
		}
	}
}

// ScanAllImageNames runs a function on all images.
func (d *Docker) ScanAllImageNames(f ImageNameScanner) {
	opts := clientLib.ListImagesOptions{All: true}
	if images, err := d.client.ListImages(opts); err != nil {
		d.logger.Fatal("Unable to list images: %s", err)
	} else {
		for _, image := range images {
			for _, tag := range image.RepoTags {
				if strings.EqualFold(tag, "<none>:<none>") {
					continue
				}
				if strings.EqualFold(tag, image.ID) {
					continue
				}
				f(NormalizeImageName(tag))
			}
		}
	}
}

// ScanAllContainerImageNames runs a function on all container images.
//
// It may scan the same image name more than once.
func (d *Docker) ScanAllContainerImageNames(f ImageNameScanner) {
	opts := clientLib.ListContainersOptions{All: true}
	if containers, err := d.client.ListContainers(opts); err != nil {
		d.logger.Fatal("Unable to list containers: %s", err)
	} else {
		for _, container := range containers {
			tag := NormalizeImageName(container.Image)
			f(tag)
		}
	}
}

// HandleImageNameEvents runs a function on incoming events that include an image name.
func (d *Docker) HandleImageNameEvents(f SawImageHandler) {
	listener := make(chan *clientLib.APIEvents, 40)
	if err := d.client.AddEventListener(listener); err != nil {
		d.logger.Fatal("Unable to listen for events: %s", err)
	}

	for {
		event := <-listener

		if event == nil {
			d.logger.Fatal("Lost connection to docker host")
		}

		switch event.Type {
		case "image":
			f(NormalizeImageName(event.Actor.ID), eventToTime(event))
		case "container":
			if tag, ok := event.Actor.Attributes["image"]; ok {
				f(NormalizeImageName(tag), eventToTime(event))
			}
		}
	}
}
