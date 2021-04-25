package module

import (
	"log"
	"time"

	"github.com/pjheden/elephantio/config"
	"github.com/pjheden/elephantio/pi"
	"github.com/pjheden/elephantio/task"
)

type Module struct {
	Task        *task.Task
	LED         *pi.LED
	Button      *pi.Button
	CompletedOn time.Time
	Config      *config.Config
}

func (m *Module) Update() bool {
	isJustCompleted := false
	// Check if button is pressed
	if m.Button.IsJustPressed() {
		log.Println("Button is pressed for module ", m.Task.Name)
		m.CompletedOn = time.Now()
		m.LED.Off()
		isJustCompleted = true
	} else if time.Since(m.CompletedOn) > (time.Duration(m.Task.Interval) * time.Hour) {
		// Check if we should light based on configs
		m.LED.On()
		// TODO: save light on LOG to db (This will require rework of logs. Something that keeps track of all events and we have to parse this list to find each modules latest "completed on")
	}
	return isJustCompleted
}
