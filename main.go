package main

import (
	"log"

	"github.com/google/gousb"
)

const vendorId = 0x0a67
const productId = 0x2114

func main() {
	log.Println("Initializing usb context")

	ctx := gousb.NewContext()
	defer ctx.Close()

	device, err := ctx.OpenDeviceWithVIDPID(vendorId, productId)
	if err != nil {
		log.Fatalf("Could not open USB device: %v", err)
	}

	defer device.Close()

	if device == nil {
		log.Fatalf("No device found with VID:PID %4x:%04x", vendorId, productId)
	}

	log.Printf("Opened device %04x:%04x\n", vendorId, productId)

	err = device.SetAutoDetach(true)
	if err != nil {
		log.Fatalf("Error turning on auto detach: %v", err)
	}

	cfg, err := device.Config(1)
	if err != nil {
		log.Fatalf("Error setting config: %v", err)
	}

	defer cfg.Close()

	intf, err := cfg.Interface(1, 0)
	if err != nil {
		log.Fatalf("Error claiming interface: %v", err)
	}

	defer intf.Close()

	inEndpoint, err := intf.InEndpoint(1)
	if err != nil {
		log.Fatalf("Error opening IN endpoint: %v", err)
	}

	buf := make([]byte, 64)

	log.Println("Waiting for MIDI data...")

	for {
		n, err := inEndpoint.Read(buf)
		if err != nil {
			log.Fatalf("Error reading from USB: %v", err)
		}
		log.Printf("Read %d bytes: %x\n", n, buf[:n])
	}
}
