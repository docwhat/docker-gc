package main

import (
	"fmt"
	"log"
	"time"

	recorder "github.com/docwhat/docker-gc/memrecorder"
	dockerclient "github.com/fsouza/go-dockerclient"
)

type imageTagRecorder interface {
	SawImageTagAt(tag string, when time.Time)
}

func main() {
	fmt.Println("Press Control-C to exit...")
	tagger := recorder.NewMemRecorder()
	schedule(tagger)
}

func schedule(tagger imageTagRecorder) {
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
			recordImage(tagger, event)
		case "container":
			recordContainer(tagger, event)
		case "network":
		default:
			log.Printf("Discarding (%s %s) %v", event.Action, event.Type, event.Actor.Attributes)
		}
	}
}

func recordImage(tagger imageTagRecorder, event *dockerclient.APIEvents) {
	if tagName, ok := event.Actor.Attributes["name"]; ok {
		name := normalizeRepoTag(tagName)
		time := time.Unix(event.Time, event.TimeNano)
		log.Printf("Image:     %8s %s", event.Action, name)
		tagger.SawImageTagAt(name, time)
	}
}

func recordContainer(tagger imageTagRecorder, event *dockerclient.APIEvents) {
	if tagName, ok := event.Actor.Attributes["image"]; ok {
		name := normalizeRepoTag(tagName)
		time := time.Unix(event.Time, event.TimeNano)
		log.Printf("Container: %8s %s", event.Action, name)
		tagger.SawImageTagAt(name, time)
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
