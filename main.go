package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/docwhat/docker-gc/memrecorder"
	"github.com/docwhat/docker-gc/types"

	cli_config "github.com/docwhat/docker-gc/config"
	dockerclient "github.com/fsouza/go-dockerclient"
)

var (
	recorder types.Recorder
	config   cli_config.Config
)

func main() {
	config = cli_config.NewConfig()
	recorder = memrecorder.NewMemRecorder()

	go recordingSchedule()

	go scanImages()
	go scanContainers()

	go deletionSweep()
	go danglingDeletionSchedule()
	fmt.Println("Press Control-C to exit...")

	// Sleep forever
	select {}
}

func hasContainerWithImage(tag string) bool {
	client := newClient()
	var containersFound int

	opts := dockerclient.ListContainersOptions{Filters: make(map[string][]string)}
	opts.Filters["ancestor"] = []string{tag}
	opts.All = true
	if containers, err := client.ListContainers(opts); err != nil {
		log.Fatalf("Failed when asking about containers with image %v: %s", tag, err)
	} else {
		containersFound = len(containers)
	}

	return (containersFound != 0)
}

func danglingDeletionSchedule() {
	for {
		time.Sleep(config.SweeperTime)
		client := newClient()
		opts := dockerclient.ListImagesOptions{Filters: make(map[string][]string)}
		opts.Filters["dangling"] = []string{"true"}
		if images, err := client.ListImages(opts); err != nil {
			log.Printf("Failed to scan for dangling images: %s", err)
		} else {
			tooOld := time.Now().Add(-2 * config.SweeperTime)
			for _, image := range images {
				created := time.Unix(image.Created, 0)
				if created.Before(tooOld) {
					if hasContainerWithImage(image.ID) {
						log.Printf("** Skipping dangling image %v running in container", image.ID)
						continue
					}
					log.Printf("** Removing dangling image %v", image.ID)
					if err := client.RemoveImage(image.ID); err != nil {
						log.Printf("Failed to remove image %v: %s", image.ID, err)
					}
				}
			}
		}
	}
}

func hasImage(tag string) bool {
	client := newClient()
	var imagesFound int

	if images, err := client.ListImages(dockerclient.ListImagesOptions{Filter: tag}); err != nil {
		log.Fatalf("Failed when asking about %v: %s", tag, err)
	} else {
		imagesFound = len(images)
	}

	return (imagesFound != 0)
}

func deletionSweepHandler(tag string, lastSeen time.Time) bool {
	client := newClient()
	tooOld := time.Now().Add(-1 * config.MaxAgeOfImages)

	if !lastSeen.Before(tooOld) {
		return false // don't delete from recorder
	}

	if hasImage(tag) {
		if err := client.RemoveImage(tag); err != nil {
			log.Printf("Failed to remove image %v: %s", tag, err)
			return false // don't delete from recorder
		}

		log.Printf("** Removing old image %v (%v old)", tag, time.Since(lastSeen))
		return true // delete from recorder
	}

	log.Printf("** Someone already removed image: %v (%v old)", tag, time.Since(lastSeen))
	return true // delete from recorder
}

// Scan for images to delete
func deletionSweep() {
	for {
		time.Sleep(config.SweeperTime)
		scanContainers()
		recorder.Sweep(deletionSweepHandler)
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
					recorder.SawImageTag(tag)
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
			recordingLog("container", "image", tag)
			recorder.SawImageTag(tag)
		}
	}
}

func recordingLog(noun string, verb string, tag string) {
	const recordingFmt = "Recorded: %-9s %-9s %s"

	if config.Verbosity >= 2 {
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
			if config.Verbosity >= 5 {
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
	recorder.SawImageTagAt(name, when)
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
