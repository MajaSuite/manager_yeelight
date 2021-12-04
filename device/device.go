package device

import "time"

const Timeout = 5 * time.Second

type Type byte

const (
	NO_TYPE Type = iota
	BULB
	LIGHTSTRIP
	NIGHTLIGHT
	CEILINGLIGHT
	DOWNLIGHT
	SPOTLIGHT
	DESKLIGHT
)

func (t Type) String() string {
	switch t {
	case BULB:
		return "Bulb"
	case LIGHTSTRIP:
		return "Light strip"
	case NIGHTLIGHT:
		return "Night lamp"
	case CEILINGLIGHT:
		return "Celing light"
	case DOWNLIGHT:
		return "Embeded downlight"
	case SPOTLIGHT:
		return "Embeded spotlight"
	case DESKLIGHT:
		return "Desk light"
	}
	return "unknown"
}

type Device interface {
	Type() Type
	ID() uint32
	IP() string
	Token() []byte
	Model() string
	Name() string
	String() string
	Close() error
	SetIP(ip string) error
}

func CheckDevice(model string) Type {
	switch model {
	case "yeelink.airp.5862":
		return NO_TYPE // unknown
	//case "yeelink.bhf_light.v1":
	//case "yeelink.bhf_light.v2":
	//case "yeelink.bhf_light.v3":
	//case "yeelink.bhf_light.v4":
	//case "yeelink.bhf_light.v5":
	//case "yeelink.bhf_light.v6":
	//case "yeelink.bhf_light.v7":
	//case "yeelink.bhf_light.v8":
	//case "yeelink.bhf_light.v9":

	//case "yeelink.controller.v2":

	//case "yeelink.curtain.ctmt1":
	//case "yeelink.curtain.ctmt1":
	//case "yeelink.curtain.procm1":
	//case "yeelink.curtain.ywc01":

	//case "yeelink.gateway.pro01":
	//case "yeelink.gateway.v1"
	//case "yeelink.gateway.va":

	//case "yeelink.light.29":
	//case "yeelink.light.bslamp1"
	//case "yeelink.light.bslamp2":
	//case "yeelink.light.bslamp3":

	//case "yeelink.light.ceila":
	//case "yeelink.light.ceilb":
	//case "yeelink.light.ceilc":
	//case "yeelink.light.ceild":
	//case "yeelink.light.ceil26":
	//case "yeelink.light.ceil27":
	//case "yeelink.light.ceil28":
	//case "yeelink.light.ceil29":
	//case "yeelink.light.ceil30":
	//case "yeelink.light.ceil31":
	//case "yeelink.light.ceil32":
	//case "yeelink.light.ceil33":
	//case "yeelink.light.ceil34":
	//case "yeelink.light.ceil35":
	//case "yeelink.light.ceil36":

	case "ceiling1":
		return CEILINGLIGHT
	case "ceiling2":
		return CEILINGLIGHT
	case "ceiling3":
		return CEILINGLIGHT
	case "ceiling4":
		return CEILINGLIGHT
	case "ceiling5":
		return CEILINGLIGHT
	case "ceiling6":
		return CEILINGLIGHT
	case "ceiling7":
		return CEILINGLIGHT
	case "ceiling8":
		return CEILINGLIGHT
	case "ceiling9":
		return CEILINGLIGHT
	case "ceiling10":
		return CEILINGLIGHT
	case "ceiling11":
		return CEILINGLIGHT
	case "ceiling12":
		return CEILINGLIGHT
	case "ceiling13":
		return CEILINGLIGHT
	case "ceiling14":
		return CEILINGLIGHT
	case "ceiling15":
		return CEILINGLIGHT
	case "ceiling16":
		return CEILINGLIGHT
	case "ceiling17":
		return CEILINGLIGHT
	case "ceiling18":
		return CEILINGLIGHT
	case "ceiling19":
		return CEILINGLIGHT
	case "ceiling20":
		return CEILINGLIGHT
	case "ceiling21":
		return CEILINGLIGHT
	case "ceiling22":
		return CEILINGLIGHT
	case "ceiling23":
		return CEILINGLIGHT
	case "ceiling24":
		return CEILINGLIGHT

	case "mono1":
		return BULB
	case "mono4":
		return BULB
	case "mono5":
		return BULB
	case "mono6":
		return BULB
	case "monoa":
		return BULB
	case "monob":
		return BULB

	//case "yeelink.light.nl1":
	//case "yeelink.light.nl2":
	//case "yeelink.light.panel1":
	//case "yeelink.light.panel3":
	//case "yeelink.light.plate1":
	//case "yeelink.light.plate2":
	//case "yeelink.light.proct1":
	//case "yeelink.light.proct2":
	//case "yeelink.light.proct3":
	//case "yeelink.light.sp1grp":
	//case "yeelink.light.spec1":
	//case "yeelink.light.spot1":
	//case "yeelink.light.spot1":
	//case "yeelink.light.spot2":

	case "strip1":
		return LIGHTSTRIP
	case "strip2":
		return LIGHTSTRIP
	case "strip4":
		return LIGHTSTRIP
	case "strip5":
		return LIGHTSTRIP
	case "strip6":
		return LIGHTSTRIP
	case "strip7":
		return LIGHTSTRIP
	case "strip8":
		return LIGHTSTRIP
	case "strip9":
		return LIGHTSTRIP
	case "stripa":
		return LIGHTSTRIP

		//case "yeelink.light.color1"
		//case "yeelink.light.color2":
		//case "yeelink.light.color3":
		//case "yeelink.light.color4":
		//case "yeelink.light.color5":
		//case "yeelink.light.color6":
		//case "yeelink.light.color7":
		//case "yeelink.light.color8":
		//case "yeelink.light.colora":
		//case "yeelink.light.colorb"
		//case "yeelink.light.colorb":
		//case "yeelink.light.colorc":

		//case "yeelink.light.ct2"
		//case "yeelink.light.cta"
		//case "yeelink.light.cta":
		//case "yeelink.light.dd005":
		//case "yeelink.light.dn2grp"
		//case "yeelink.light.dn2grp":
		//case "yeelink.light.dn3grp":
		//case "yeelink.light.dnlight2"
		//case "yeelink.light.dnlight2":
		//case "yeelink.light.fancl1":
		//case "yeelink.light.fancl2":
		//case "yeelink.light.fancl3":
		//case "yeelink.light.fancl5":
		//case "yeelink.light.fancl6":

		//case "yeelink.light.lamp1":
		//case "yeelink.light.lamp2":
		//case "yeelink.light.lamp3":
		//case "yeelink.light.lamp4"
		//case "yeelink.light.lamp5"
		//case "yeelink.light.lamp7":
		//case "yeelink.light.lamp9":
		//case "yeelink.light.lamp10":
		//case "yeelink.light.lamp11":
		//case "yeelink.light.lamp12":
		//case "yeelink.light.lamp13":
		//case "yeelink.light.lamp14":
		//case "yeelink.light.lamp15":
		//case "yeelink.light.lamp16":
		//case "yeelink.light.lamp17":
		//case "yeelink.light.lamp18":
		//case "yeelink.light.lamp19":
		//case "yeelink.light.lamp20":
		//case "yeelink.light.lamp21"
		//case "yeelink.light.lamp22":
		//case "yeelink.light.lampb":
		//case "yeelink.light.lamps":
		//case "yeelink.light.lampv":

		//case "yeelink.light.light3":
		//case "yeelink.light.light4":
		//case "yeelink.light.light5":

		//case "yeelink.light.mb1grp"
		//case "yeelink.light.mb1grp":
		//case "yeelink.light.mb2grp"
		//case "yeelink.light.mb2grp":
		//case "yeelink.light.mb3grp"
		//case "yeelink.light.mbulb3":
		//case "yeelink.light.mbulb4":
		//case "yeelink.light.mbulb5":

		//case "yeelink.light.meshbulb1"
		//case "yeelink.light.meshbulb1":
		//case "yeelink.light.meshbulb2"
		//case "yeelink.light.meshdev":

		//case "yeelink.light.ml1":
		//case "yeelink.light.ml2":
		//case "yeelink.light.ml3":
		//case "yeelink.light.mla":

		//case "yeelink.light.test":
		//case "yeelink.light.tlight":
		//case "yeelink.light.tmbulb":
		//case "yeelink.light.yct01":
		//case "yeelink.light.ydim01":
		//case "yeelink.light.yrgb01":
		//case "yeelink.magnet.ycs01":
		//case "yeelink.magnet.ycs02":
		//case "yeelink.mirror.bm1":
		//case "yeelink.motion.ymt01":
		//case "yeelink.motion.ymt02":
		//case "yeelink.plug.plug":
		//case "yeelink.plug.prosw":
		//case "yeelink.plug.scene":
		//case "yeelink.remote.remote":
		//case "yeelink.remote.yrm01":
		//case "yeelink.remote.yrm02":
		//case "yeelink.switch.prosw1":
		//case "yeelink.switch.sw1":
		//case "yeelink.switch.sw2":
		//case "yeelink.switch.test":
		//case "yeelink.switch.xxx":
		//case "yeelink.switch.ysw01":
		//case "yeelink.uwb.tag1":
		//case "yeelink.ven_fan.vf1":
		//case "yeelink.ven_fan.vf3":
		//case "yeelink.ven_fan.vf4":
		//case "yeelink.ven_fan.vf5":
		//case "yeelink.wifispeaker.v1":
	}

	return NO_TYPE
}
