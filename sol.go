package sol

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/merliot/dean"
	"github.com/merliot/device"
	"github.com/merliot/device/modbus"
)

//go:embed css go.mod html js
var fs embed.FS

//template NEED TO ADD THIS

const (
	regBattInfo  = 0x3045
	regLoadInfo  = 0x304A
	regSolarInfo = 0x304E
	regSysInfo   = 0x9027
)

type System struct {
	Config uint16
}

type Battery struct {
	Remaining float32
	Volts     float32
	Amps      float32
}

type LoadInfo struct {
	Volts float32
	Amps  float32
}

type Solar struct {
	Volts float32
	Amps  float32
}

type msgSystem struct {
	Path   string
	System System
}

type msgStatus struct {
	Path   string
	Status string
}

type msgBattery struct {
	Path    string
	Battery Battery
}

type msgLoadInfo struct {
	Path     string
	LoadInfo LoadInfo
}

type msgSolar struct {
	Path  string
	Solar Solar
}

type Sol struct {
	*device.Device
	*modbus.Modbus `json:"-"`
	Status         string
	Battery        Battery
	LoadInfo       LoadInfo
	Solar          Solar
	System         System
}

var targets = []string{"demo", "nano-rp2040"}

func New(id, model, name string) dean.Thinger {
	println("NEW SOLDEVICE")
	s := &Sol{}
	s.Device = device.New(id, model, name, targets).(*device.Device)
	s.Modbus = modbus.New(s)
	s.Status = "OK"
	return s
}

func (s *Sol) save(msg *dean.Msg) {
	msg.Unmarshal(s).Broadcast()
}

func (s *Sol) getState(msg *dean.Msg) {
	s.Path = "state"
	msg.Marshal(s).Reply()
}

func (s *Sol) Subscribers() dean.Subscribers {
	return dean.Subscribers{
		"state":          s.save,
		"get/state":      s.getState,
		"update/system":  s.save,
		"update/battery": s.save,
		"update/load":    s.save,
		"update/solar":   s.save,
	}
}

func (s *Sol) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.API(w, r, s)
}

func serial(b []byte) string {
	return fmt.Sprintf("%02X%02X-%02X%02X", b[0], b[1], b[2], b[3])
}

func swap(b []byte) uint16 {
	return (uint16(b[0]) << 8) | uint16(b[1])
}

func volts(b []byte) float32 {
	return float32(swap(b)) * 0.01
}

func amps(b []byte) float32 {
	return float32(swap(b)) * 0.01
}

func percent(b []byte) float32 {
	return float32(swap(b)) * 1
}

func (s *Sol) readSystem(sys *System) error {
	// System Info (34 bytes)
	regs, err := s.ReadRegisters(regSysInfo, 1)
	if err != nil {
		return err
	}
	sys.Config = swap(regs[0:2])

	return nil
}

func (s *Sol) readDynamic(b *Battery, l *LoadInfo, pv *Solar) error {

	// Controller Dynamic Info (6 bytes)
	regs, err := s.ReadRegisters(regBattInfo, 3)
	if err != nil {
		return err
	}
	b.Remaining = percent(regs[0:2]) //FIX
	b.Volts = volts(regs[2:4])       //FIX
	b.Amps = amps(regs[4:6])         //FIX

	// Solar information ()
	regs, err = s.ReadRegisters(regSolarInfo, 2)
	if err != nil {
		return err
	}
	pv.Volts = volts(regs[0:2])
	pv.Amps = amps(regs[2:4])

	// Load Information (4 bytes)
	regs, err = s.ReadRegisters(regLoadInfo, 2)
	if err != nil {
		return err
	}
	l.Volts = volts(regs[0:1])
	l.Amps = amps(regs[1:2])

	return nil
}

func (s *Sol) sendStatus(i *dean.Injector, newStatus string) {
	if s.Status == newStatus {
		return
	}

	var status = msgStatus{Path: "update/status"}
	var msg dean.Msg

	status.Status = newStatus
	i.Inject(msg.Marshal(status))
}

func (s *Sol) sendSystem(i *dean.Injector) {
	var system = msgSystem{Path: "update/system"}
	var msg dean.Msg

	// sendSystem blocks until we get a good system info read

	for {
		if err := s.readSystem(&system.System); err != nil {
			s.sendStatus(i, err.Error())
			continue
		}
		i.Inject(msg.Marshal(system))
		break
	}

	s.sendStatus(i, "OK")
}

func (s *Sol) sendDynamic(i *dean.Injector) {
	var battery = msgBattery{Path: "update/battery"}
	var loadInfo = msgLoadInfo{Path: "update/load"}
	var solar = msgSolar{Path: "update/solar"}
	var msg dean.Msg

	err := s.readDynamic(&battery.Battery,
		&loadInfo.LoadInfo, &solar.Solar)
	if err != nil {
		s.sendStatus(i, err.Error())
		return
	}

	// If anything has changed, send update msg(s)

	if battery.Battery != s.Battery {
		i.Inject(msg.Marshal(battery))
	}
	if loadInfo.LoadInfo != s.LoadInfo {
		i.Inject(msg.Marshal(loadInfo))
	}
	if solar.Solar != s.Solar {
		i.Inject(msg.Marshal(solar))
	}

	s.sendStatus(i, "OK")
}
