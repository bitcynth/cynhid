package main

import (
	"log"

	"github.com/bitcynth/cynhid"
)

var usbPacketGetVersion = []byte{0x12, 0x20}

func main() {
	err := cynhid.Init()
	if err != nil {
		log.Fatalf("failed to initialize hidapi: %v", err)
	}

	devInfos, err := cynhid.Enumerate(0x04d9, 0x0348)
	if err != nil {
		log.Fatal(err)
	}

	var devPath string
	for _, devInfo := range devInfos {
		if devInfo.InterfaceNumber == 1 {
			devPath = devInfo.Path
		}
	}

	dev, err := cynhid.OpenPath(devPath)
	if err != nil {
		log.Fatal(err)
	}

	dev.SetNonblocking(true)

	dev.Write(padData(usbPacketGetVersion, 64), 64)

	dev.Read(64)

	dev.Close()

	cynhid.Exit()
}

func padData(data []byte, length int) []byte {
	dataLen := len(data)

	for i := 0; i < length-dataLen; i++ {
		data = append(data, 0x00)
	}

	return data
}
