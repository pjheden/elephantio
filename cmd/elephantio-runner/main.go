package main

import (
	"log"

	"github.com/pjheden/elephantio/database"
	"github.com/pjheden/elephantio/pi"
)

func main() {

	p, err := pi.New()
	if err != nil {
		log.Fatal("can't open pi: ", err)
	}
	defer p.Close()

	dbPath := "./db.sqlite"
	db, err := database.New(dbPath)
	if err != nil {
		log.Fatal("can't open db: ", err)
	}

	ms, err := db.Modules()
	if err != nil {
		log.Fatal(err)
	}

	for _, m := range ms {
		log.Printf("got ms %+v\n", m)

		// setup pi connections configs
		// Button
		log.Println("setup button on pin ", m.Config.ButtonPin)
		btn, err := pi.NewButton(m.Config.ButtonPin)
		if err != nil {
			log.Fatal(err)
		}
		m.Button = btn

		// LED
		log.Println("setup LED on pin ", m.Config.LEDPin)
		l, err := pi.NewLED(m.Config.LEDPin)
		if err != nil {
			log.Fatal(err)
		}
		m.LED = l
	}

	// Board loop
	for {
		for _, m := range ms {
			m.Update()
		}
	}

}
