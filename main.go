package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/google/gousb"
)

const vendorId = 0x0a67
const productId = 0x2114

func main() {
	errChan := make(chan error)
	go read(errChan)

	for err := range errChan {
		log.Fatalf("Error in reader: %v", err)
	}
}

func read(errChan chan error) {
	err := readLoop()
	if err != nil {
		errChan <- err
	}
}

func readLoop() error {
	log.Println("Initializing usb context")

	ctx := gousb.NewContext()
	defer ctx.Close()

	device, err := ctx.OpenDeviceWithVIDPID(vendorId, productId)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not open USB device: %v", err))
	}

	defer device.Close()

	if device == nil {
		return errors.New(fmt.Sprintf("No device found with VID:PID %4x:%04x", vendorId, productId))
	}

	log.Printf("Opened device %04x:%04x\n", vendorId, productId)

	err = device.SetAutoDetach(true)
	if err != nil {
		return errors.New(fmt.Sprintf("Error turning on auto detach: %v", err))
	}

	cfg, err := device.Config(1)
	if err != nil {
		return errors.New(fmt.Sprintf("Error setting config: %v", err))
	}

	defer cfg.Close()

	intf, err := cfg.Interface(1, 0)
	if err != nil {
		return errors.New(fmt.Sprintf("Error claiming interface: %v", err))
	}

	defer intf.Close()

	inEndpoint, err := intf.InEndpoint(1)
	if err != nil {
		return errors.New(fmt.Sprintf("Error opening IN endpoint: %v", err))
	}

	buf := make([]byte, 64)

	log.Println("Waiting for MIDI data...")

	for {
		n, err := inEndpoint.Read(buf)
		if err != nil {
			return errors.New(fmt.Sprintf("Error reading from USB: %v", err))
		}

		if n%4 != 0 {
			return errors.New(fmt.Sprintf("Data length is not a multiple of 4: %v %v", n, buf))
		}

		if n < 4 {
			continue
		}

		if buf[0] == 0b00001111 {
			continue
		}

		log.Printf("Read %d bytes: %08b\n", n, buf[:n])

		if buf[0] != 0b00001001 {
			return errors.New(fmt.Sprintf("Unkown first byte %0b", buf[0]))
		}

		var status int
		if buf[1] == 0x90 {
			status = NoteOn
		} else {
			return errors.New(fmt.Sprintf("Unkown status byte %0b", buf[1]))
		}

		if status == NoteOn {
			noteNumber := buf[2]
			remainder := noteNumber % 12

			var note string
			if remainder == 0 {
				note = C
			} else if remainder == 1 {
				note = CSharp
			} else if remainder == 2 {
				note = D
			} else if remainder == 3 {
				note = DSharp
			} else if remainder == 4 {
				note = E
			} else if remainder == 5 {
				note = F
			} else if remainder == 6 {
				note = FSharp
			} else if remainder == 7 {
				note = G
			} else if remainder == 8 {
				note = GSharp
			} else if remainder == 9 {
				note = A
			} else if remainder == 10 {
				note = ASharp
			} else if remainder == 11 {
				note = B
			}

			velocity := buf[3]

			if velocity == 0 {
				status = NoteOff
				log.Printf("NoteOff %s\n", note)
			} else {
				log.Printf("NoteOn %s with %d/127 velocity\n", note, int(velocity))
			}

		}
	}
}

const (
	NoteOn  = iota
	NoteOff = iota

	C      = "C"
	CSharp = "C#"
	D      = "D"
	DSharp = "D#"
	E      = "E"
	F      = "F"
	FSharp = "F#"
	G      = "G"
	GSharp = "G#"
	A      = "A"
	ASharp = "A#"
	B      = "B"
)
