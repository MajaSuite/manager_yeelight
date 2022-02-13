package device

import (
	"fmt"
)

type LightDevice struct {
	YeelightDevice
	Bright    int `json:"bright"`
	ColorMode int `json:"mode"`
	ColorTemp int `json:"temp"`
}

func NewLightDevice(debug bool, model string, id string, ip string, name string, support string, power bool, ver int,
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

func (d *LightDevice) String() string {
	return fmt.Sprintf(`{%s,"bright":%d,"mode":%d,"temp":%d}`,
		d.YeelightDevice.String(), d.Bright, d.ColorMode, d.ColorTemp)
}

func (d *LightDevice) Retain() string {
	return fmt.Sprintf(`{%s}`, d.YeelightDevice.Retain())
}

func (orig *LightDevice) CompareAndUpdate(dev *Device) error {
	source := (*dev).(*LightDevice)

	if orig.Name != source.Name {
		_, err := orig.Run("set_name", []interface{}{source.Name})
		if err == nil {
			orig.Name = source.Name
		}
	}
	if orig.Power != source.Power {
		power := "off"
		if source.Power {
			power = "on"
		}
		_, err := orig.Run("set_power", []interface{}{power, defaultEffect, defaultEffectDuration, 0})
		if err == nil {
			orig.Power = source.Power
		}
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

	return nil
}

func (orig *LightDevice) RunMethod(dev *Device, method string, effect string, duration int) error {
	source := (*dev).(*LightDevice)

	var err error
	switch method {
	case "set_ct_abx":
	case "set_bright":
	case "set_power":
		power := "off"
		if source.Power {
			power = "on"
		}
		_, err := orig.Run("set_power", []interface{}{power, defaultEffect, defaultEffectDuration, 0})
		if err == nil {
			orig.Power = source.Power
		}
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
