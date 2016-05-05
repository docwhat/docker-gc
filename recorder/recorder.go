package recorder

import (
	"log"

	dockerclient "github.com/fsouza/go-dockerclient"
)

func Start() {
	client := newClient()

	listener := make(chan *dockerclient.APIEvents)
	if err := client.AddEventListener(listener); err != nil {
		log.Fatalf("Unable to listen for events: %s", err)
	}

	for {
		event := <-listener
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

func newClient() *dockerclient.Client {
	client, err := dockerclient.NewClientFromEnv()
	if err != nil {
		log.Fatalf("Unable to connect to docker daemon: %s", err)
	}
	return client
}

func recordImage(event *dockerclient.APIEvents) {
	if name, ok := event.Actor.Attributes["name"]; ok {
		log.Printf("Image:     %s %s", event.Action, name)
	}
}

func recordContainer(event *dockerclient.APIEvents) {
	if name, ok := event.Actor.Attributes["image"]; ok {
		log.Printf("Container: %s %s", event.Action, name)
	}
}
