package device

import (
	"bufio"
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
	Method      string          `json:"method,omitempty"`
	Effect      string          `json:"effect,omitempty"`
	Duration    int             `json:"duration,omitempty"`
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

func readResponse(debug bool, conn net.Conn, id int) ([]string, error) {
	//	conn.SetReadDeadline(time.Now().Add(time.Second * 10))
	response, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		if debug {
			log.Println("lan read", err)
		}
		return nil, err
	}

	log.Println("run lan response", strings.Trim(string(response), "\r\n"))

	var resp lanResponse
	if err := json.Unmarshal(response, &resp); err != nil {
		if debug {
			log.Println("lan error unmarshal message", err)
		}
		return nil, err
	}

	if resp.Id == id {
		if resp.Error != nil {
			return nil, fmt.Errorf("lan error from device %s", resp.Error.Message)
		}
		if resp.Result != nil && len(resp.Result) > 0 {
			return resp.Result, nil
		}
	} else {
		return readResponse(debug, conn, id)
	}

	return nil, ErrWrongParameter
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

func (d *YeelightDevice) SetIP(ip string) error {
	if ip == "" {
		return ErrIpUnknown
	}

	if d.Ip != "" {
		if d.debug {
			log.Printf("%s allready have ip (%s)", d.Id, d.IP())
		}
		return ErrAlreadyStarted
	}

	d.Ip = ip

	return nil
}

func (d *YeelightDevice) Run(method string, props []interface{}) ([]string, error) {
	if d.IP() == "" {
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

func (d *YeelightDevice) CompareAndUpdate(dev *Device) error {
	// nothing to do here
	return nil
}

func (d *YeelightDevice) RunMethod(dev *Device, method string, effect string, duration int) error {
	// do nothing here
	return nil
}
