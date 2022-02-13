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

func (rgb *RgbDevice) String() string {
	return fmt.Sprintf(`{%s,"bright":%d,"mode":%d,"temp":%d,"rgb":%d,"hue":%d,"sat":%d}`,
		rgb.YeelightDevice.String(), rgb.Bright, rgb.ColorMode, rgb.ColorTemp, rgb.Rgb, rgb.Hue, rgb.Sat)
}

func (d *RgbDevice) Retain() string {
	return fmt.Sprintf(`{%s}`, d.YeelightDevice.Retain())
}

func (orig *RgbDevice) CompareAndUpdate(dev *Device) error {
	source := (*dev).(*RgbDevice)

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
	if orig.Rgb != source.Rgb {
		//
	}
	if orig.Hue != source.Hue {
		//
	}
	if orig.Sat != source.Sat {
		//
	}

	return nil
}

func (orig *RgbDevice) RunMethod(dev *Device, method string, effect string, duration int) error {
	source := (*dev).(*RgbDevice)

	var err error
	switch method {
	case "set_ct_abx":
	case "set_rgb":
	case "set_hsv":
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
	}

	return err
}
