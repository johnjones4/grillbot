package outputs

import (
	"fmt"
	"main/device"
	"os"
)

type Disk struct {
	file *os.File
}

func NewDisk(path string) (*Disk, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return &Disk{f}, nil
}

func (d *Disk) receiveUpdates(_ *device.Device, m device.Message) {
	d.file.WriteString(fmt.Sprintf("%f,%f\n", m.Temp1, m.Temp2))
}

func (d *Disk) Listener() device.Listener {
	return d.receiveUpdates
}

func (d *Disk) Close() {
	d.file.Close()
}
