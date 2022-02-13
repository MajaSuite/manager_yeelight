package device

import (
	"errors"
	"time"
)

const (
	connectTimeout        = 2 * time.Second
	yeeLightPort          = "55443"
	defaultEffect         = Smooth
	defaultEffectDuration = 3
)

var (
	ErrEmptyString    = errors.New("empty string")
	ErrCantConnect    = errors.New("can't connect to device")
	ErrInvalidCommand = errors.New("invalid command")
	ErrWrongParameter = errors.New("invalid parameter")
	ErrNotStarted     = errors.New("device not started")
	ErrIpUnknown      = errors.New("device ip address unknown")
	ErrAlreadyStarted = errors.New("device already started")
)

const (
	NO_TYPE Type = iota
	LIGHT_DEVICE
	RGB_DEVICE
	AMBILIGHT_DEVICE
)

type Type int

type Device interface {
	Type() Type
	Model() string
	ID() string
	IP() string
	SetIP(ip string) error
	Run(method string, props []interface{}) ([]string, error)
	CompareAndUpdate(dev *Device) error
	RunMethod(dev *Device, method string, effect string, duration int) error
	Close() error
	String() string
	Retain() string
}

func (t Type) String() string {
	switch t {
	case LIGHT_DEVICE:
		return "Light"
	case RGB_DEVICE:
		return "RGB light"
	case AMBILIGHT_DEVICE:
		return "Light with ambilight"
	}
	return "unknown"
}

/* Define deviceType by model string
 */
func CheckDevice(model string) Type {
	switch model {
	case "ceiling1":
		return LIGHT_DEVICE
	case "ceiling2":
		return LIGHT_DEVICE
	case "ceiling3":
		return LIGHT_DEVICE
	case "ceiling4":
		return LIGHT_DEVICE
	case "ceiling5":
		return LIGHT_DEVICE
	case "ceiling6":
		return LIGHT_DEVICE
	case "ceiling7":
		return LIGHT_DEVICE
	case "ceiling8":
		return LIGHT_DEVICE
	case "ceiling9":
		return LIGHT_DEVICE
	case "ceiling10":
		return LIGHT_DEVICE
	case "ceiling11":
		return LIGHT_DEVICE
	case "ceiling12":
		return LIGHT_DEVICE
	case "ceiling13":
		return LIGHT_DEVICE
	case "ceiling14":
		return LIGHT_DEVICE
	case "ceiling15":
		return LIGHT_DEVICE
	case "ceiling16":
		return LIGHT_DEVICE
	case "ceiling17":
		return LIGHT_DEVICE
	case "ceiling18":
		return LIGHT_DEVICE
	case "ceiling19":
		return LIGHT_DEVICE
	case "ceiling20":
		return RGB_DEVICE // ambi
	case "ceiling21":
		return LIGHT_DEVICE
	case "ceiling22":
		return LIGHT_DEVICE
	case "ceiling23":
		return LIGHT_DEVICE
	case "ceiling24":
		return LIGHT_DEVICE
	case "ceila":
		return LIGHT_DEVICE
	case "ceilb":
		return LIGHT_DEVICE
	case "ceilc":
		return LIGHT_DEVICE
	case "ceild":
		return LIGHT_DEVICE
	case "ceil26":
		return LIGHT_DEVICE
	case "ceil27":
		return LIGHT_DEVICE
	case "ceil28":
		return LIGHT_DEVICE
	case "ceil29":
		return LIGHT_DEVICE
	case "ceil30":
		return LIGHT_DEVICE
	case "ceil31":
		return LIGHT_DEVICE
	case "ceil32":
		return LIGHT_DEVICE
	case "ceil33":
		return LIGHT_DEVICE
	case "ceil34":
		return LIGHT_DEVICE
	case "ceil35":
		return LIGHT_DEVICE
	case "ceil36":
		return LIGHT_DEVICE
	case "mono1":
		return LIGHT_DEVICE
	case "mono4":
		return LIGHT_DEVICE
	case "mono5":
		return LIGHT_DEVICE
	case "mono6":
		return LIGHT_DEVICE
	case "monoa":
		return LIGHT_DEVICE
	case "monob":
		return LIGHT_DEVICE
	case "color1":
		return RGB_DEVICE
	case "color2":
		return RGB_DEVICE
	case "color3":
		return RGB_DEVICE
	case "color4":
		return RGB_DEVICE
	case "color5":
		return RGB_DEVICE
	case "color6":
		return RGB_DEVICE
	case "color7":
		return RGB_DEVICE
	case "color8":
		return RGB_DEVICE
	case "colora":
		return RGB_DEVICE
	case "colorb":
		return RGB_DEVICE
	case "colorc":
		return RGB_DEVICE
	case "strip1":
		return RGB_DEVICE
	case "strip2":
		return RGB_DEVICE
	case "strip4":
		return RGB_DEVICE
	case "strip5":
		return RGB_DEVICE
	case "strip6":
		return RGB_DEVICE
	case "strip7":
		return RGB_DEVICE
	case "strip8":
		return RGB_DEVICE
	case "strip9":
		return RGB_DEVICE
	case "stripa":
		return RGB_DEVICE

		//case "yeelink.bhf_light.v1":
		//case "yeelink.bhf_light.v2":
		//case "yeelink.bhf_light.v3":
		//case "yeelink.bhf_light.v4":
		//case "yeelink.bhf_light.v5":
		//case "yeelink.bhf_light.v6":
		//case "yeelink.bhf_light.v7":
		//case "yeelink.bhf_light.v8":
		//case "yeelink.bhf_light.v9":

		//case "29":
		//case "bslamp1"
		//case "bslamp2":
		//case "bslamp3":

		//case "nl1":
		//case "nl2":
		//case "panel1":
		//case "panel3":
		//case "plate1":
		//case "plate2":
		//case "proct1":
		//case "proct2":
		//case "proct3":
		//case "sp1grp":
		//case "spec1":
		//case "spot1":
		//case "spot1":
		//case "spot2":

		//case "ct2"
		//case "cta"
		//case "cta":
		//case "dd005":
		//case "dn2grp"
		//case "dn2grp":
		//case "dn3grp":
		//case "dnlight2"
		//case "dnlight2":
		//case "fancl1":
		//case "fancl2":
		//case "fancl3":
		//case "fancl5":
		//case "fancl6":

		//case "lamp1":
		//case "lamp2":
		//case "lamp3":
		//case "lamp4"
		//case "lamp5"
		//case "lamp7":
		//case "lamp9":
		//case "lamp10":
		//case "lamp11":
		//case "lamp12":
		//case "lamp13":
		//case "lamp14":
		//case "lamp15":
		//case "lamp16":
		//case "lamp17":
		//case "lamp18":
		//case "lamp19":
		//case "lamp20":
		//case "lamp21"
		//case "lamp22":
		//case "lampb":
		//case "lamps":
		//case "lampv":

		//case "light3":
		//case "light4":
		//case "light5":

		//case "mb1grp"
		//case "mb1grp":
		//case "mb2grp"
		//case "mb2grp":
		//case "mb3grp"
		//case "mbulb3":
		//case "mbulb4":
		//case "mbulb5":

		//case "meshbulb1"
		//case "meshbulb1":
		//case "meshbulb2"
		//case "meshdev":

		//case "ml1":
		//case "ml2":
		//case "ml3":
		//case "mla":

		//case "tlight":
		//case "tmbulb":
		//case "yct01":
		//case "ydim01":
		//case "yrgb01":
	}

	return NO_TYPE
}

/* Create device using sets of values
 */
func CreateDevice(debug bool, model string, id string, ip string, name string, support string, power bool, ver int,
	bright int, mode int, temp int, rgb int, hue int, sat int) Device {
	var dev Device

	switch CheckDevice(model) {
	case RGB_DEVICE:
		dev = NewRgbDevice(debug, model, id, ip, name, support, power, ver, bright, mode, temp, rgb, hue, sat)
	case LIGHT_DEVICE:
		dev = NewLightDevice(debug, model, id, ip, name, support, power, ver, bright, mode, temp)
	case AMBILIGHT_DEVICE:
		// todo
	}

	return dev
}
