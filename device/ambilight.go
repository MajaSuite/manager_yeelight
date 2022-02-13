package device

import (
	"fmt"
)

type AmbilightDevice struct {
	YeelightDevice
	Bright      int  `json:"bright"`
	ColorMode   int  `json:"mode"`
	ColorTemp   int  `json:"temp"`
	BgPower     bool `json:"bgpower"`
	BgBright    int  `json:"bgbright"`
	BgColorMode int  `json:"bgmode"`
	BgColorTemp int  `json:"bgtemp"`
	BgRgb       int  `json:"bgrgb"`
	BgHue       int  `json:"bghue"`
	BgSat       int  `json:"bgsat"`
}

func NewAmbilightDevice() Device {
	return nil
}

func (rgb *AmbilightDevice) String() string {
	return fmt.Sprintf(`{%s,"bright":%d,"mode":%d,"temp":%d}`,
		rgb.YeelightDevice.String(), rgb.Bright, rgb.ColorMode, rgb.ColorTemp)
}

func (d *AmbilightDevice) Retain() string {
	return fmt.Sprintf(`{%s}`, d.YeelightDevice.Retain())
}

func (orig *AmbilightDevice) CompareAndUpdate(dev *Device) error {
	source := (*dev).(*AmbilightDevice)

	if orig.Name != source.Name {
		//
	}
	if orig.Power != source.Power {
		//
	}
	if orig.Bright != source.Bright {
		//
	}
	if orig.ColorMode != source.ColorMode {
		//
	}
	if orig.ColorTemp != source.ColorTemp {
		//
	}
	// todo ....
	// add other checks (bg)...

	return nil
}

func (orig *AmbilightDevice) RunMethod(dev *Device, method string, effect string, duration int) error {
	source := (*dev).(*AmbilightDevice)

	var err error
	switch method {
	case "set_ct_abx":
	case "set_bright":
	case "set_power":
	case "toggle":
		_, err = orig.Run(method, []interface{}{})
		if err == nil {
			orig.Power = !source.Power
		}
	case "set_default":
	case "start_cf":
	case "stop_cf":
	case "set_scene":
	case "cron_add":
	case "cron_get":
	case "cron_del":
	case "set_adjust":
	case "set_music":
	case "set_name":
		_, err = orig.Run(method, []interface{}{source.Name})
		if err == nil {
			orig.Name = source.Name
		}
	case "dev_toggle":
	case "adjust_bright":
	case "adjust_ct":
	case "adjust_color":
	case "bg_set_rgb":
	case "bg_set_hsv":
	case "bg_set_ct_abx":
	case "bg_start_cf":
	case "bg_stop_cf":
	case "bg_set_scene":
	case "bg_set_default":
	case "bg_set_bright":
	case "bg_set_power":
	case "bg_set_adjust":
	case "bg_toggle":
	case "bg_adjust_bright":
	case "bg_adjust_ct":
	case "bg_adjust_color":
	}

	return nil
}
