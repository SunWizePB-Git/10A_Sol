package sol

import (
	"github.com/merliot/device"
	"github.com/merliot/device/modbus"
)

type Sol struct {
	*device.Device
	*modbus.Modbus `json:"-"`
	
}
