package device

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

type YeelightDevice struct {
	deviceType  Type            `json:"type"`
	deviceModel string          `json:"model"`
	Id          string          `json:"-"`
	Ip          string          `json:"-"`
	Name        string          `json:"name"`
	Support     string          `json:"support"`
	Version     int             `json:"ver"`
	Power       bool            `json:"power"`
	support     map[string]bool `json:"-"`
	debug       bool            `json:"-"`
	//status      net.Conn        `json:"-"`
	Method   string `json:"method,omitempty"`
	Effect   string `json:"effect,omitempty"`
	Duration string `json:"duration,omitempty"`
}

func NewYeeLightDevice(debug bool, model string, id string, ip string, name string, support string, power bool,
	ver int) *YeelightDevice {
	d := &YeelightDevice{
		deviceType:  CheckDevice(model),
		deviceModel: model,
		Id:          id,
		Ip:          ip,
		Name:        name,
		Support:     support,
		Power:       power,
		Version:     ver,
		debug:       debug,
	}
	d.ConvertSupport(support)

	return d
}

func (d *YeelightDevice) Type() Type {
	return d.deviceType
}

func (d *YeelightDevice) Model() string {
	return d.deviceModel
}

func (d *YeelightDevice) ID() string {
	return d.Id
}

func (d *YeelightDevice) IP() string {
	return d.Ip
}

func (d *YeelightDevice) ConvertSupport(support string) error {
	if support == "" {
		return ErrEmptyString
	}

	d.support = make(map[string]bool)
	sp := strings.Split(support, " ")
	for _, v := range sp {
		d.support[v] = true
	}
	return nil
}

func (d *YeelightDevice) Close() error {
	//return d.status.Close()
	d.Ip = ""
	return nil
}

func (d *YeelightDevice) String() string {
	return fmt.Sprintf(`"type":%d,"model":"%s","id":"%s","ip":"%s","name":"%s","support":"%s","power":%v,"ver":%d`,
		d.deviceType, d.deviceModel, d.Id, d.Ip, d.Name, d.Support, d.Power, d.Version)
}

func (d *YeelightDevice) Retain() string {
	return fmt.Sprintf(`"type":%d,"model":"%s","id":"%s","name":"%s","support":"%s","ver":%d`,
		d.deviceType, d.deviceModel, d.Id, d.Name, d.Support, d.Version)
}

func (d *YeelightDevice) Connect(ip string, update chan Device) error {
	if d.Ip != "" {
		if d.debug {
			log.Printf("%s allready connected (%s)", d.Id, d.IP())
		}
		return ErrAlreadyStarted
	}

	d.Ip = ip
	if d.debug {
		log.Printf("start device %s (%s)", d.Id, d.Ip)
	}

	//go func() {
	//	for {
	//		if d.debug {
	//			log.Println("connect for notifications")
	//		}
	//
	//		var err error
	//		if d.status, err = net.DialTimeout("tcp", net.JoinHostPort(d.Ip, yeeLightPort), connectTimeout); err != nil {
	//			time.Sleep(time.Second * 5)
	//			continue
	//		}
	//
	//		for {
	//			d.status.SetReadDeadline(time.Now().Add(time.Second * 10))
	//			if line, err := bufio.NewReader(d.status).ReadBytes('\n'); err == nil && len(line) > 0 {
	//				if d.debug {
	//					log.Println("new message from device", strings.Trim(string(line), "\r\n"))
	//				}
	//
	//				var resp lanRequest
	//				if err := json.Unmarshal(line, &resp); err != nil {
	//					log.Println("error unmarshal message", err)
	//					continue
	//				}
	//
	//				//log.Println(resp.Params)
	//				// update <- message
	//
	//				continue
	//			}
	//
	//			d.status.SetWriteDeadline(time.Now().Add(time.Second * 2))
	//			if _, err := d.status.Write([]byte("\n")); err != nil {
	//				d.status.Close()
	//				break
	//			}
	//		}
	//	}
	//}()

	res, err := d.Run("get_prop", []string{"name"})
	if err == nil && len(res) > 0 {
		d.Name = res[0]
		// todo update it?
	}

	if err != nil {
		return err
	}

	return nil
}

func (d *YeelightDevice) Run(method string, props []string) ([]string, error) {
	if /*d.status == nil ||*/ d.IP() == "" {
		return nil, ErrNotStarted
	}

	if !d.support[method] {
		return nil, ErrInvalidCommand
	}

	id := 1

	var conn net.Conn
	var err error
	for i := 0; i < 3; i++ {
		if conn, err = net.DialTimeout("tcp", net.JoinHostPort(d.IP(), yeeLightPort), connectTimeout); err != nil {
			if d.debug {
				log.Println("lan device reconnect to send command")
			}
			time.Sleep(connectTimeout)
			continue
		}
		break
	}
	defer conn.Close()

	if err != nil {
		if d.debug {
			log.Println("lan can't connect to device")
			return nil, ErrCantConnect
		}
	}

	req := lanRequest{Id: id, Method: method, Params: props}
	request, err := json.Marshal(req)
	if err != nil {
		if d.debug {
			log.Println("lan error marshal request", err)
		}
		return nil, err
	}

	request = append(request, "\r\n"...)

	if d.debug {
		log.Println("run lan request", strings.Trim(string(request), "\r\n"))
	}

	if _, err := conn.Write(request); err != nil {
		if d.debug {
			log.Println("lan write", err)
		}
		return nil, err
	}

	return readResponse(d.debug, conn, id)
}
