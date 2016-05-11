package memrecorder

import "time"

// MemRecorder stores all seen image tags in memory.
type MemRecorder struct {
	imageTags map[string]time.Time
}

// An ImageTagSweeper is a method to run on each image tag.
type ImageTagSweeper func(tag string, when time.Time)

// NewMemRecorder initializes a new MemRecorder for use.
func NewMemRecorder() *MemRecorder {
	var r MemRecorder
	r.imageTags = make(map[string]time.Time)
	return &r
}

// SawImageTagAt records when a tag was last seen.
//
// Duplicate or older times will be ignored.
func (r *MemRecorder) SawImageTagAt(tag string, when time.Time) {
	if old, ok := r.imageTags[tag]; ok {
		if when.Before(old) {
			return // We don't need to adjust the value.
		}
	}

	r.imageTags[tag] = when
}

// Sweep runs a function on all tag and timestamp pairs.
func (r *MemRecorder) Sweep(sweeper ImageTagSweeper) {
	for tag, when := range r.imageTags {
		sweeper(tag, when)
	}
}

// func Start() {
//   client := newClient()

//   listener := make(chan *dockerclient.APIEvents)
//   if err := client.AddEventListener(listener); err != nil {
//     log.Fatalf("Unable to listen for events: %s", err)
//   }

//   for {
//     event := <-listener
//     switch event.Type {
//     case "image":
//       recordImage(event)
//     case "container":
//       recordContainer(event)
//     case "network":
//     default:
//       log.Printf("Discarding (%s %s) %v", event.Action, event.Type, event.Actor.Attributes)
//     }
//   }
// }

// func newClient() *dockerclient.Client {
//   client, err := dockerclient.NewClientFromEnv()
//   if err != nil {
//     log.Fatalf("Unable to connect to docker daemon: %s", err)
//   }
//   return client
// }

// func recordImage(event *dockerclient.APIEvents) {
//   if name, ok := event.Actor.Attributes["name"]; ok {
//     log.Printf("Image:     %s %s", event.Action, name)
//   }
// }

// func recordContainer(event *dockerclient.APIEvents) {
//   if name, ok := event.Actor.Attributes["image"]; ok {
//     log.Printf("Container: %s %s", event.Action, name)
//   }
// }
