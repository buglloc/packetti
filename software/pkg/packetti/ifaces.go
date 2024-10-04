package packetti

import (
	"fmt"
	"strings"

	"go.bug.st/serial/enumerator"
)

const (
	piRatVID = "2E8A"
	piRatPID = "5052"
)

type Iface struct {
	Path string
	Name string
}

func Ifaces() ([]Iface, error) {
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		return nil, fmt.Errorf("list detailed ports: %w", err)
	}

	var out []Iface
	for _, port := range ports {
		if !port.IsUSB {
			continue
		}

		if !strings.EqualFold(port.VID, piRatVID) {
			continue
		}

		if !strings.EqualFold(port.PID, piRatPID) {
			continue
		}

		out = append(out, Iface{
			Path: port.Name,
			Name: fmt.Sprintf("PiRAT:%s", port.SerialNumber),
		})
	}

	return out, nil
}
