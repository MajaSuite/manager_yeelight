package device

import (
	"fmt"
)

type RgbDevice struct {
	YeelightDevice
	Bright    int `json:"bright"`
	ColorMode int `json:"mode"`
	ColorTemp int `json:"temp"`
	Rgb       int `json:"rgb"`
	Hue       int `json:"hue"`
	Sat       int `json:"sat"`
}

func (rgb *RgbDevice) String() string {
	return fmt.Sprintf(`{%s,"bright":%d,"mode":%d,"temp":%d,"rgb":%d,"hue":%d,"sat":%d}`,
		rgb.YeelightDevice.String(), rgb.Bright, rgb.ColorMode, rgb.ColorTemp, rgb.Rgb, rgb.Hue, rgb.Sat)
}

func (d *RgbDevice) Retain() string {
	return fmt.Sprintf(`{%s}`, d.YeelightDevice.Retain())
}

func NewRgbDevice(debug bool, model string, id string, ip string, name string, support string, power bool, ver int,
	bright int, mode int, temp int, rgb int, hue int, sat int) Device {
	d := &RgbDevice{
		YeelightDevice: *NewYeeLightDevice(debug, model, id, ip, name, support, power, ver),
		Bright:         bright,
		ColorMode:      mode,
		ColorTemp:      temp,
		Rgb:            rgb,
		Hue:            hue,
		Sat:            sat,
	}
	d.ConvertSupport(support)
	return d
}
