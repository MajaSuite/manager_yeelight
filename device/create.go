package device

import (
	"errors"
	"manager_yeelight/utils"
)

var (
	ErrDeviceUpdated = errors.New("device require update")
)

func UpdateDevice(dev interface{}, lines map[string]string) error {
	var updateRequired = false

	if y, ok := dev.(*YeeLight); ok {
		version := utils.ConvertInt(lines["fw_ver"])
		if y.Version != version {
			updateRequired = true
			y.Version = version
		}

		power := utils.ConvertBool(lines["power"])
		if y.Power != power {
			updateRequired = true
			y.Power = power
		}

		bright := utils.ConvertInt(lines["bright"])
		if y.Bright != bright {
			updateRequired = true
			y.Bright = bright
		}

		cm := utils.ConvertInt(lines["color_mode"])
		if y.ColorMode != cm {
			updateRequired = true
			y.ColorMode = cm
		}

		ct := utils.ConvertInt(lines["ct"])
		if y.ColorTemp != ct {
			updateRequired = true
			y.ColorTemp = ct
		}

		rgb := utils.ConvertInt(lines["rgb"])
		if y.Rgb != rgb {
			updateRequired = true
			y.Rgb = rgb
		}

		hue := utils.ConvertInt(lines["hue"])
		if y.Hue != hue {
			updateRequired = true
			y.Hue = hue
		}

		sat := utils.ConvertInt(lines["sat"])
		if y.Sat != sat {
			updateRequired = true
			y.Sat = sat
		}
	}

	if updateRequired {
		return ErrDeviceUpdated
	}
	return nil
}

func CreateDevice(l map[string]string) Device {
	switch CheckDevice(l["model"]) {
	case NO_TYPE:
		return nil
	default:
		d := NewYeeLight(utils.ConvertHex(l["id"]), l["model"], l["name"], l["support"])
		UpdateDevice(d, l)
		return d
	}
}
