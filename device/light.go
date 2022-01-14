package device

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"manager_yeelight/utils"
	"net"
	"strings"
	"time"
)

const (
	YeeLightPort = "55443"
)

type YeeLight struct {
	DeviceId    uint32
	DeviceIp    string
	DeviceType  Type
	DeviceModel string
	Name        string
	Version     int
	Power       bool
	Bright      int
	ColorMode   int
	ColorTemp   int
	Rgb         int
	Hue         int
	Sat         int
	support     map[string]bool
	conn        net.Conn
	counter     int
}

func NewYeeLight(id uint32, model string, name string, support string) *YeeLight {
	y := &YeeLight{
		DeviceId:    id,
		DeviceType:  CheckDevice(model),
		DeviceModel: model,
		Name:        name,
	}
	y.Support(support)
	return y
}

func (g *YeeLight) Type() Type {
	return g.DeviceType
}

func (g *YeeLight) Id() uint32 {
	return g.DeviceId
}

func (y *YeeLight) Start(ip string, update chan Device) error {
	if y.conn != nil {
		return ErrAlreadyStarted
	}

	log.Println("Start device", ip)
	y.DeviceIp = ip

	if conn, err := net.DialTimeout("tcp", net.JoinHostPort(y.DeviceIp, YeeLightPort), Timeout); err != nil {
		return err
	} else {
		y.conn = conn

		// read device status. separate connection
		go func() {
			var status net.Conn

			for {
				if status == nil {
					log.Println("status thread reconnect", ip)
					status, err = net.DialTimeout("tcp", net.JoinHostPort(y.DeviceIp, YeeLightPort), Timeout)
					if err != nil {
						time.Sleep(time.Minute)
						continue
					}
				}

				response, err := bufio.NewReader(bufio.NewReader(status)).ReadBytes('\n')
				if err != nil {
					log.Println("status error", ip, err)
					status = nil
					continue
				}

				var resp notification
				if err := json.Unmarshal(response, &resp); err != nil {
					log.Println("notification error", err)
					continue
				}

				/*
									2022/01/14 16:25:43 lan lanResponse {"method":"props","params":{"power":"on"}}
								{"method":"props","params":{"indicator_on":0}}
							{"method":"props","params":{"hue":36,"rgb":16750848,"flowing":0,"color_mode":1,"bright":1}}
						{"method":"props","params":{"color_mode":2,"ct":3200,"bright":80}}
					{"method":"props","params":{"flow_params":"0,1,3000,1,16711680,100,3000,1,65280,100,3000,1,255,100,3000,1,9055202,100","flowing":1}}

				*/
				log.Println("=============== notification ", string(response))
				if resp.Params.Power != "" {
					y.Power = utils.ConvertBool(resp.Params.Power)
				}
				update <- y
			}
		}()

		res, err := y.Cmd("get_prop", "name", "", "")
		if err != nil {
			return err
		}
		if len(res) > 0 {
			y.Name = res[0]
			update <- y
		}
		return nil
	}
}

func (y *YeeLight) Close() error {
	return y.conn.Close()
}

func (y *YeeLight) Cmd(method string, v1 string, v2 string, v3 string) ([]string, error) {
	if y.conn == nil {
		return nil, ErrNotStarted
	}

	if !y.support[method] {
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

	y.counter++

	req := &lanRequest{Id: y.counter, Method: method, Params: props}
	request, err := json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	request = append(request, '\r')
	request = append(request, '\n')

	log.Println("lan request", strings.Trim(string(request), "\r\n"))

	if _, err := y.conn.Write(request); err != nil {
		log.Println("write", err)
		y.conn, _ = net.DialTimeout("tcp", net.JoinHostPort(y.DeviceIp, YeeLightPort), Timeout)
		return y.Cmd(method, v1, v2, v3)
	}

	response, err := bufio.NewReader(bufio.NewReader(y.conn)).ReadBytes('\n')
	if err != nil {
		log.Println("read", err)
		y.conn, _ = net.DialTimeout("tcp", net.JoinHostPort(y.DeviceIp, YeeLightPort), Timeout)
		return nil, err
	}

	log.Println("lan response", strings.Trim(string(response), "\r\n"))

	var resp lanResponse
	if err := json.Unmarshal(response, &resp); err != nil {
		return nil, err
	}

	if resp.Id == y.counter {
		if resp.Error != nil {
			log.Printf("lan error %s", resp.Error.Message)
			return nil, fmt.Errorf("yeelight error %s", resp.Error.Message)
		}
		if resp.Result != nil && len(resp.Result) > 0 {
			return resp.Result, nil
		}
		return nil, nil
	}

	return nil, nil
}

func (y *YeeLight) Support(v string) {
	y.support = make(map[string]bool)
	sp := strings.Split(v, " ")
	for _, v := range sp {
		y.support[v] = true
	}
}

func (y *YeeLight) String() string {
	var support string
	for v := range y.support {
		if support == "" {
			support = v
		} else {
			support = support + " " + v
		}
	}
	return fmt.Sprintf(`{"id":"%x","ip":"%s","type":"%s","model":"%s","name":"%s","ver":%d,"support":"%v","power":%v,"bright":%d,"mode":%d,"temp":%d,"rgb":%d,"hue":%d,"sat":%d}`,
		y.DeviceId, y.DeviceIp, y.DeviceType, y.DeviceModel, y.Name, y.Version, support, y.Power, y.Bright, y.ColorMode, y.ColorTemp, y.Rgb, y.Hue, y.Sat)
}
