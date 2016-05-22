package main

import (
	"fmt"
	"log"
	"time"

	recorder "github.com/docwhat/docker-gc/memrecorder"
	dockerclient "github.com/fsouza/go-dockerclient"
)

var tagger *recorder.MemRecorder

func main() {
	fmt.Println("Press Control-C to exit...")
	tagger = recorder.NewMemRecorder()
	schedule()
}

func schedule() {
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
			recordImage(event)
		case "container":
			recordContainer(event)
		case "network":
		default:
			log.Printf("Discarding (%s %s) %v", event.Action, event.Type, event.Actor.Attributes)
		}
	}
}

func record(repoTag string, event *dockerclient.APIEvents) {
	name := normalizeRepoTag(repoTag)
	time := time.Unix(event.Time, event.TimeNano)
	log.Printf("Recorded: %-9s %-8s %s", event.Type, event.Action, name)
	tagger.SawImageTagAt(name, time)
}

func recordImage(event *dockerclient.APIEvents) {
	if tagName, ok := event.Actor.Attributes["name"]; ok {
		record(tagName, event)
	}
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
