package main

import (
	"fmt"
	"os"
	"time"
)

type app struct {
	config   appConfig
	docker   *Docker
	recorder Recorder
	logger   Logger
}

func main() {
	config, err := newAppConfig(os.Args[1:])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(2)
	}
	logger := NewLogger(config)
	main := app{config: config, docker: NewDocker(logger), logger: logger}

	main.recorder = NewMemRecorder()

	logger.Info("Press Control-C to exit...")

	main.scanAllImageNames()

	// go recordingSchedule()
	go main.deleteDanglingLoop()
	go main.deleteOldImagesLoop()
	go main.listenForEvents()

	// Sleep forever
	select {}
}

func (m app) deleteDangling() {
	m.docker.ScanDanglingImages(m.config.DangleSafeDuration, func(tag string) {
		m.logger.Info("Removing dangling image %s", tag)
		m.docker.RemoveImage(tag)
	})
}

func (m app) deleteDanglingLoop() {
	for {
		time.Sleep(m.config.SweeperTime)
		m.deleteDangling()
	}
}

func (m app) scanAllImageNames() {
	m.docker.ScanAllImageNames(func(tag string) {
		m.logger.Info("Saw image %s", tag)
		m.recorder.SawImageTag(tag)
	})
}

func (m app) scanContainerImageNames() {
	m.docker.ScanAllContainerImageNames(func(tag string) {
		m.logger.Info("Saw image %s in container", tag)
		m.recorder.SawImageTag(tag)
	})
}

func (m app) deleteOldImagesHandler(tag string, lastSeen time.Time) bool {
	age := time.Since(lastSeen)

	if age < (4 * time.Second) {
		return false
	}

	if age < m.config.MaxAgeOfImages {
		m.logger.Info("Not deleting %s because it is only %s old", tag, age)
		return false // don't delete from recorder
	}

	if m.docker.HasImage(tag) {
		m.logger.Info("Removing image %s because it is %s old", tag, age)
		return m.docker.RemoveImage(tag)
	}

	m.logger.Info("** Someone already removed image: %v (%v old)", tag, time.Since(lastSeen))
	return true // delete from recorder
}

// Scan for images to delete
func (m app) deleteOldImagesLoop() {
	for {
		time.Sleep(m.config.SweeperTime)
		m.scanContainerImageNames()
		m.recorder.Sweep(m.deleteOldImagesHandler)
	}
}

func (m app) listenForEvents() {
	m.docker.HandleImageNameEvents(func(tag string, when time.Time) {
		m.logger.Info("Event: %s at %s", tag, when)
		m.recorder.SawImageTagAt(tag, when)
	})
}
