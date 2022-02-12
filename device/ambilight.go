package device

import "fmt"

type AmbilightDevice struct {
	YeelightDevice
	Bright    int `json:"bright"`
	ColorMode int `json:"mode"`
	ColorTemp int `json:"temp"`
	//...
}

func (d *AmbilightDevice) Retain() string {
	return fmt.Sprintf(`{%s}`, d.YeelightDevice.Retain())
}
