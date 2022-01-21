package device

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"manager_yeelight/utils"
	"net"
	"strings"
	"time"
)

const (
	connectTimeout = 5 * time.Second
	yeeLightPort   = "55443"
)

const (
	NO_TYPE Type = iota
	BULB
	LIGHTSTRIP
	NIGHTLIGHT
	CEILINGLIGHT
	DOWNLIGHT
	SPOTLIGHT
	DESKLIGHT
	CURTAIN
	PLUG
	FAN
)

type Type int

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
	case CURTAIN:
		return "Curtain"
	case PLUG:
		return "Power plug"
	case FAN:
		return "Fan"
	}
	return "unknown"
}

type Device struct {
	Id        string          `json:"id"`
	Ip        string          `json:"ip"`
	Type      Type            `json:"type"`
	Model     string          `json:"model"`
	Name      string          `json:"name"`
	Version   int             `json:"ver"`
	Support   string          `json:"support"`
	Power     bool            `json:"power"`
	Bright    int             `json:"bright"`
	ColorMode int             `json:"mode"`
	ColorTemp int             `json:"temp"`
	Rgb       int             `json:"rgb"`
	Hue       int             `json:"hue"`
	Sat       int             `json:"sat"`
	Cmd       string          `json:"cmd,omitempty"`
	Value1    string          `json:"value1,omitempty"`
	Value2    string          `json:"value2,omitempty"`
	Value3    string          `json:"value3,omitempty"`
	support   map[string]bool `json:"-"`
	conn      net.Conn        `json:"-"`
	status    net.Conn        `json:"-"`
	counter   int             `json:"-"`
}

func (d *Device) IsStarted() bool {
	if d.counter > 0 {
		return true
	}
	return false
}

func (d *Device) CheckDevice(model string) Type {
	switch model {
	case "ctmt1":
		return CURTAIN
	case "procm1":
		return CURTAIN
	case "ywc01":
		return CURTAIN
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

	//case "yeelink.airp.5862":

	//case "yeelink.bhf_light.v1":
	//case "yeelink.bhf_light.v2":
	//case "yeelink.bhf_light.v3":
	//case "yeelink.bhf_light.v4":
	//case "yeelink.bhf_light.v5":
	//case "yeelink.bhf_light.v6":
	//case "yeelink.bhf_light.v7":
	//case "yeelink.bhf_light.v8":
	//case "yeelink.bhf_light.v9":

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

	case "plug":
		return PLUG
	case "prosw":
		return PLUG
	case "scene":
		return PLUG
	case "vf1":
		return FAN
	case "vf3":
		return FAN
	case "vf4":
		return FAN
	case "vf5":
		return FAN
	}

	return NO_TYPE
}

func (d *Device) ConvertSupport() error {
	if d.Support == "" {
		return fmt.Errorf("empty support array")
	}

	d.support = make(map[string]bool)
	sp := strings.Split(d.Support, " ")
	for _, v := range sp {
		d.support[v] = true
	}
	return nil
}

func (d *Device) Run(method string, v1 string, v2 string, v3 string) ([]string, error) {
	if d.counter == 0 {
		return nil, ErrNotStarted
	}

	if !d.support[method] {
		return nil, ErrInvalidCommand
	}

	var props = []string{}
	if v1 != "" {
		props = append(props, v1)
	}
	if v2 != "" {
		props = append(props, v2)
	}
	if v3 != "" {
		props = append(props, v3)
	}

	d.counter++
	req := &lanRequest{Id: d.counter, Method: method, Params: props}
	request, err := json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	request = append(request, '\r')
	request = append(request, '\n')

	log.Println("lan request", strings.Trim(string(request), "\r\n"))

	for x := 0; x < 5; x++ {
		if _, err := d.conn.Write(request); err != nil {
			log.Println("lan write", err)
			if x > 5 {
				break // too much tries, break it
			}
			time.Sleep(time.Second * 10)
			d.conn, _ = net.DialTimeout("tcp", net.JoinHostPort(d.Ip, yeeLightPort), connectTimeout)
			continue
		}
		break
	}

	response, err := bufio.NewReader(bufio.NewReader(d.conn)).ReadBytes('\n')
	if err != nil {
		log.Println("lan read", err)
		return nil, err
	}

	log.Println("lan response", strings.Trim(string(response), "\r\n"))

	var resp lanResponse
	if err := json.Unmarshal(response, &resp); err != nil {
		return nil, err
	}

	if resp.Id == d.counter {
		if resp.Error != nil {
			return nil, fmt.Errorf("lan error from device %s", resp.Error.Message)
		}
		if resp.Result != nil && len(resp.Result) > 0 {
			return resp.Result, nil
		}
	} else {
		response, err := bufio.NewReader(bufio.NewReader(d.conn)).ReadBytes('\n')
		if err != nil {
			return nil, err
		}

		log.Println("lan additional", strings.Trim(string(response), "\r\n"))
		if err := json.Unmarshal(response, &resp); err != nil {
			return nil, err
		}
		if resp.Result != nil && len(resp.Result) > 0 {
			return resp.Result, nil
		}
	}

	return nil, nil
}

func (d *Device) String() string {
	var support string
	for v := range d.support {
		if support == "" {
			support = v
		} else {
			support = support + " " + v
		}
	}
	return fmt.Sprintf(`{"id":"%s","ip":"%s","model":"%s","name":"%s","ver":%d,"support":"%v","power":%v,"bright":%d,"mode":%d,"temp":%d,"rgb":%d,"hue":%d,"sat":%d}`,
		d.Id, d.Ip, d.Model, d.Name, d.Version, support, d.Power, d.Bright, d.ColorMode, d.ColorTemp, d.Rgb, d.Hue, d.Sat)
}

func (d *Device) Start(update chan *Device) error {
	if d.conn != nil || d.counter > 0 {
		return ErrAlreadyStarted
	}

	if d.Ip == "" {
		return ErrIpUnknown
	}

	log.Printf("start device %s (%s)", d.Id, d.Ip)
	d.counter++

	// connect to device
	go func() {
		for d.counter > 0 {
			conn, err := net.DialTimeout("tcp", net.JoinHostPort(d.Ip, yeeLightPort), connectTimeout)
			if err == nil {
				d.conn = conn

				res, _ := d.Run("get_prop", "name", "", "")
				if len(res) > 0 {
					d.Name = res[0]
				}

				update <- d
				return
			}

			time.Sleep(time.Second * 10)
		}
	}()

	// read device status. separate connection
	go func() {
		log.Println("status thread connect", d.Ip)
		d.status, _ = net.DialTimeout("tcp", net.JoinHostPort(d.Ip, yeeLightPort), connectTimeout)

		for {
			if d.status == nil {
				log.Println("status thread reconnect", d.Ip)
				status, err := net.DialTimeout("tcp", net.JoinHostPort(d.Ip, yeeLightPort), connectTimeout)
				if err != nil {
					time.Sleep(time.Second * 10)
					continue
				}
				d.status = status
			}

			response, err := bufio.NewReader(bufio.NewReader(d.status)).ReadBytes('\n')
			if err != nil {
				log.Println("status connection", d.Ip, err)
				d.status = nil
				time.Sleep(time.Second * 10)
				continue
			}

			var resp lanRequest
			if err := json.Unmarshal(response, &resp); err != nil {
				log.Printf("error unmarshall notification record (%v)", err)
				continue
			}
			log.Println("new notification", string(response))

			if resp.Method == "props" {
				for k, v := range resp.Params.(map[string]interface{}) {
					switch k {
					case "power":
						d.Power = utils.ConvertBool(v.(string))
						update <- d
					case "name":
						d.Name = v.(string)
						update <- d
					case "hue":
						d.Hue = v.(int)
						update <- d
					case "sat":
						d.Sat = v.(int)
						update <- d
					case "rgb":
						d.Rgb = v.(int)
						update <- d
					case "color_mode":
						d.ColorMode = v.(int)
						update <- d
					case "ct":
						d.ColorTemp = v.(int)
						update <- d
					case "bright":
						d.Bright = v.(int)
						update <- d
						// todo some extra also may exists
						// i.e. {"flow_params":"0,1,3000,1,16711680,100,3000,1,65280,100,3000,1,255,100,3000,1,9055202,100","flowing":1}
					}
				}
			}
		}
	}()

	return nil
}

func (d *Device) Close() error {
	return d.conn.Close()
}

func NewDevice(buf *bytes.Buffer, addr net.Addr) *Device {
	device := &Device{
		Ip: strings.Split(addr.String(), ":")[0],
	}

	for {
		chunk, err := buf.ReadBytes('\n')
		if err == io.EOF {
			break
		}

		line := strings.Split(string(chunk), ": ")
		if len(line) >= 2 {
			v := strings.Trim(string(chunk)[len(line[0])+1:], " \r\n")
			switch line[0] {
			case "id":
				device.Id = strings.TrimLeft(strings.Replace(v, "0x", "", -1), "0")
			case "model":
				device.Model = v
				device.Type = device.CheckDevice(v)
			case "name":
				device.Name = v
			case "fw_ver":
				device.Version = utils.ConvertInt(v)
			case "power":
				device.Power = utils.ConvertBool(v)
			case "bright":
				device.Bright = utils.ConvertInt(v)
			case "color_mode":
				device.ColorMode = utils.ConvertInt(v)
			case "ct":
				device.ColorTemp = utils.ConvertInt(v)
			case "rgb":
				device.Rgb = utils.ConvertInt(v)
			case "hue":
				device.Hue = utils.ConvertInt(v)
			case "sat":
				device.Sat = utils.ConvertInt(v)
			case "support":
				device.Support = v
				device.ConvertSupport()
			}
		}
	}

	return device
}
