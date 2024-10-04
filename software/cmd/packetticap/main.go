package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gopacket/gopacket"
	"github.com/gopacket/gopacket/layers"
	"github.com/gopacket/gopacket/pcapgo"
	"github.com/kor44/extcap"

	"github.com/buglloc/packetti/software/pkg/packetti"
)

func main() {
	app := extcap.App{
		Usage: "piratcap",
		Version: extcap.VersionInfo{
			Info: "1.0.0",
			Help: "https://github.com/buglloc/PiRAT",
		},
		HelpPage:      "PiRATCap - extcap application to integrate PiRAT with Wireshark or something",
		GetInterfaces: getAllInterfaces,
		GetDLT:        getDLT,
		StartCapture:  startCapture,
	}

	app.Run(os.Args)
}

func getAllInterfaces() ([]extcap.CaptureInterface, error) {
	ifaces, err := packetti.Ifaces()
	if err != nil {
		return nil, fmt.Errorf("unable to get information about interfaces: %w", err)
	}

	extIfaces := make([]extcap.CaptureInterface, len(ifaces))
	for i, iface := range ifaces {
		// we use Name as Value to deal with PiRAt replugs
		extIfaces[i] = extcap.CaptureInterface{
			Value:   iface.Name,
			Display: iface.Path,
		}
	}

	return extIfaces, nil
}

func getDLT(_ string) (extcap.DLT, error) {
	return extcap.DLT{
		Number: packetti.LinkTypeUSBFullSpeed.Int(),
		Name:   packetti.LinkTypeUSBFullSpeed.String(),
	}, nil
}

func startCapture(iface string, pipe io.WriteCloser, _ string, _ map[string]interface{}) error {
	defer func() { _ = pipe.Close() }()

	r, err := packetti.NewDeviceByName(iface)
	if err != nil {
		return fmt.Errorf("open PiRAT reader: %w", err)
	}
	defer func() { _ = r.Close() }()

	w, err := pcapgo.NewNgWriterInterface(
		pipe,
		pcapgo.NgInterface{
			Name:                filepath.Base(iface),
			OS:                  runtime.GOOS,
			LinkType:            layers.LinkType(packetti.LinkTypeUSBFullSpeed.Int()),
			SnapLength:          0, //unlimited
			TimestampResolution: 9,
		},
		pcapgo.NgWriterOptions{
			SectionInfo: pcapgo.NgSectionInfo{
				Hardware:    runtime.GOARCH,
				OS:          runtime.GOOS,
				Application: "PiRAT", //spread the word
			},
		},
	)
	if err != nil {
		return fmt.Errorf("open pcapng writer: %w", err)
	}
	defer func() { _ = w.Flush() }()

	for {
		packet, err := r.Packet()
		if err != nil {
			return fmt.Errorf("read packet: %w", err)
		}

		ci := gopacket.CaptureInfo{
			Length:         len(packet),
			CaptureLength:  len(packet),
			InterfaceIndex: 0,
		}
		err = w.WritePacket(ci, packet)
		if err != nil {
			return fmt.Errorf("write packet: %w", err)
		}
	}
}
