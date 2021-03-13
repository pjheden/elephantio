package pi

import (
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

type Pi struct {
	// Buttons []*Button
}

// func (p *Pi) AddButton(task string, pin int) error {
// 	pb, err := NewButton(pin)
// 	if err != nil {
// 		return err
// 	}

// 	p.Buttons = append(p.Buttons, pb)
// 	return nil
// }

func (p *Pi) Close() {
	// Unmap gpio memory when done
	rpio.Close()
}

func New() (*Pi, error) {
	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		return nil, err
	}

	return &Pi{}, nil
}

type Button struct {
	Pin         rpio.Pin
	heldDown    bool
	lastPressed time.Time // TODO: implement this so the button have a "cooldown", it helps with the light flickering
	cooldown    time.Duration
}

func NewButton(p int) (*Button, error) {
	pin := rpio.Pin(p)
	pin.Input() // Input mode
	return &Button{Pin: pin, cooldown: 100 * time.Millisecond}, nil
}

func (b *Button) IsJustPressed() bool {
	// Read state from pin (High / Low)
	pressed := b.Pin.Read() == rpio.Low
	if pressed {
		if b.heldDown {
			return false
		} else {
			if time.Since(b.lastPressed) < b.cooldown {
				return false
			}
			b.heldDown = true
			b.lastPressed = time.Now()
			return true
		}
	}
	b.heldDown = false
	return false
}

type LED struct {
	Pin  rpio.Pin
	IsOn bool
}

func NewLED(p int) (*LED, error) {
	pin := rpio.Pin(p)
	pin.Output() // output mode
	pin.Low()    // make sure it is off

	return &LED{Pin: pin, IsOn: false}, nil
}

func (l *LED) Toggle() {
	// Read state from pin (High / Low)
	l.Pin.Toggle()
	l.IsOn = !l.IsOn
}

func (l *LED) On() {
	l.Pin.High()
	l.IsOn = true
}
func (l *LED) Off() {
	l.Pin.Low()
	l.IsOn = false
}
