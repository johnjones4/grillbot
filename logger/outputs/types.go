package outputs

import "main/device"

type Output interface {
	Listener() device.Listener
	Close()
}
