package device

import "fmt"

type LightDevice struct {
	YeelightDevice
	Bright    int `json:"bright"`
	ColorMode int `json:"mode"`
	ColorTemp int `json:"temp"`
}

func (d *LightDevice) String() string {
	return fmt.Sprintf(`{%s,"bright":%d,"mode":%d,"temp":%d}`,
		d.YeelightDevice.String(), d.Bright, d.ColorMode, d.ColorTemp)
}

func (d *LightDevice) Retain() string {
	return fmt.Sprintf(`{%s}`, d.YeelightDevice.Retain())
}

func NewBWDevice(debug bool, model string, id string, ip string, name string, support string, power bool, ver int,
	bright int, mode int, temp int) Device {
	d := &LightDevice{
		YeelightDevice: *NewYeeLightDevice(debug, model, id, ip, name, support, power, ver),
		Bright:         bright,
		ColorMode:      mode,
		ColorTemp:      temp,
	}
	d.ConvertSupport(support)
	return d
}
