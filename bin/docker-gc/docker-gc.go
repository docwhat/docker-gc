package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	recorder "github.com/docwhat/docker-gc/memrecorder"
	dockerclient "github.com/fsouza/go-dockerclient"
)

var tagger *recorder.MemRecorder

// FIXME: These need to switches to the command
// var maxAgeOfImages = 7 * 24 * time.Hour
var maxAgeOfImages = 10 * time.Second
var sweeperTime = 5 * time.Second
var verbosity = 9

func main() {
	tagger = recorder.NewMemRecorder()

	go recordingSchedule()

	go scanImages()
	go scanContainers()

	go deletionSweep()
	go danglingDeletionSchedule()
	fmt.Println("Press Control-C to exit...")

	// Sleep forever
	select {}
}

func danglingDeletionSchedule() {
	for {
		time.Sleep(sweeperTime)
		client := newClient()
		opts := dockerclient.ListImagesOptions{Filters: make(map[string][]string)}
		opts.Filters["dangling"] = []string{"true"}
		if images, err := client.ListImages(opts); err != nil {
			log.Printf("Failed to scan for dangling images: %s", err)
		} else {
			for _, image := range images {
				// FIXME: Check that creation time isn't in the last sweeperTime + 1 minute
				log.Printf("** Removing dangling image %v", image.ID)
				if err := client.RemoveImage(image.ID); err != nil {
					log.Printf("Failed to remove image %v: %s", image.ID, err)
				}
			}
		}
	}
}

// Scan for images to delete
func deletionSweep() {
	var sweeper recorder.ImageTagSweeper = func(tag string, lastSeen time.Time) bool {
		client := newClient()
		tooOld := time.Now().Add(-1 * maxAgeOfImages)
		if lastSeen.Before(tooOld) {
			log.Printf("** Removing old image %v (%v old)", tag, time.Since(lastSeen))
			if images, err := client.ListImages(dockerclient.ListImagesOptions{Filter: tag}); err != nil {
				log.Printf("Failed when asking about %v: %s", tag, err)
			} else {
				if len(images) == 0 {
					log.Printf("   ...%v is already removed", tag)
					return true
				}
			}
			if err := client.RemoveImage(tag); err != nil {
				log.Printf("Failed to remove image %v: %s", tag, err)
				return false
			}
			return true
		}
		return false
	}

	for {
		time.Sleep(sweeperTime)
		scanContainers()
		tagger.Sweep(sweeper)
	}
}

func scanImages() {
	client := newClient()
	opts := dockerclient.ListImagesOptions{All: true}
	if images, err := client.ListImages(opts); err != nil {
		log.Fatalf("Unable to list images: %s", err)
	} else {
		for _, image := range images {
			for _, tag := range image.RepoTags {
				if !strings.EqualFold(tag, "<none>:<none>") {
					tag = normalizeRepoTag(tag)
					recordingLog("existing", "image", tag)
					tagger.SawImageTag(tag)
				}
			}
		}
	}
}

func scanContainers() {
	client := newClient()
	opts := dockerclient.ListContainersOptions{All: true}
	if containers, err := client.ListContainers(opts); err != nil {
		log.Fatalf("Unable to list containers: %s", err)
	} else {
		for _, container := range containers {
			tag := normalizeRepoTag(container.Image)
			recordingLog("running", "container", tag)
			tagger.SawImageTag(tag)
		}
	}
}

func recordingLog(noun string, verb string, tag string) {
	const recordingFmt = "Recorded: %-9s %-9s %s"

	if verbosity >= 2 {
		log.Printf(recordingFmt, noun, verb, tag)
	}
}

func recordingSchedule() {
	client := newClient()

	listener := make(chan *dockerclient.APIEvents)
	if err := client.AddEventListener(listener); err != nil {
		log.Fatalf("Unable to listen for events: %s", err)
	}

	for {
		event := <-listener

		if event == nil {
			log.Fatalf("Lost connection to docker host")
		}

		switch event.Type {
		case "image":
			record(event.Actor.ID, event)
		case "container":
			recordContainer(event)
		case "network":
		default:
			if verbosity >= 5 {
				log.Printf("Discarding event (%s %s) %v", event.Action, event.Type, event.Actor.Attributes)
			}
		}
	}
}

func eventTime(event *dockerclient.APIEvents) time.Time {
	if event.TimeNano != 0 {
		return time.Unix(0, event.TimeNano)
	} else if event.Time != 0 {
		return time.Unix(event.TimeNano, 0)
	}
	return time.Now()
}

func record(repoTag string, event *dockerclient.APIEvents) {
	if strings.HasPrefix(repoTag, "sha256:") {
		return
	}
	name := normalizeRepoTag(repoTag)
	when := eventTime(event)
	recordingLog(event.Type, event.Action, name)
	tagger.SawImageTagAt(name, when)
}

func recordContainer(event *dockerclient.APIEvents) {
	if tagName, ok := event.Actor.Attributes["image"]; ok {
		record(tagName, event)
	}
}

func normalizeRepoTag(repoTag string) string {
	repo, tag := dockerclient.ParseRepositoryTag(repoTag)
	if tag == "" {
		tag = "latest"
	}
	return fmt.Sprintf("%s:%s", repo, tag)
}

func newClient() *dockerclient.Client {
	client, err := dockerclient.NewClientFromEnv()
	if err != nil {
		log.Fatalf("Unable to create docker client: %s", err)
	}

	if err := client.Ping(); err != nil {
		log.Printf("Cannot connect to the Docker daemon. Is the docker daemon running on this host?")
		log.Fatalf("%s", err)
	}
	return client
}
