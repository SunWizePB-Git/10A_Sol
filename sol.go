package sol

import (
	"github.com/merliot/dean"
	"github.com/merliot/device"
	"github.com/merliot/device/modbus"
)

type Sol struct {
	*device.Device
	*modbus.Modbus `json:"-"`
	
}

var targets = []string{"demo", "nano-rp2040"}

func New(id, model, name string) dean.Thinger {
	println("NEW SOLDEVICE")
	s := &Sol{}
	s.Device = device.New(id, model, name, targets).(*device.Device)
	s.Modbus = modbus.New(s)
	return s
}
