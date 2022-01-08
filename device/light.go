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
	deviceType Type
	token      []byte
	id         uint32
	ip         string
	model      string
	FwVer      int
	Nam        string
	support    []string
	Power      bool
	Bright     int
	ColorMode  int
	ColorTemp  int
	Rgb        int
	Hue        int
	Sat        int
	conn       net.Conn
	counter    uint32
}

// Create device object. ip should set separatly.
func NewYeeLight(id uint32, token []byte, model string, fw int, support []string, power bool, bright int, colorMode int,
	temp int, rgb int, hue int, sat int, name string) *YeeLight {

	return &YeeLight{
		id:         id,
		deviceType: CheckDevice(model),
		token:      token,
		model:      model,
		FwVer:      fw,
		support:    support,
		Power:      power,
		Bright:     bright,
		ColorMode:  colorMode,
		ColorTemp:  temp,
		Rgb:        rgb,
		Hue:        hue,
		Sat:        sat,
		Nam:        name,
	}
}

func (y *YeeLight) Type() Type {
	return y.deviceType
}

func (y *YeeLight) ID() uint32 {
	return y.id
}

func (y *YeeLight) IP() string {
	return y.ip
}

func (y *YeeLight) SetIP(ip string) error {
	if conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, YeeLightPort), Timeout); err == nil {
		y.ip = ip
		y.conn = conn

		conn.SetReadDeadline(time.Now().Add(Timeout))
		if n, err := y.GetProp([]Prop{"name"}); len(n) > 0 && err == nil {
			y.Nam = string(n[0])
		}
		return nil
	} else {
		return err
	}
}

func (y *YeeLight) Token() []byte {
	return y.token
}

func (y *YeeLight) Model() string {
	return y.model
}

func (y *YeeLight) Name() string {
	if y.Nam == "" {
		if n, err := y.GetProp([]Prop{"name"}); len(n) > 0 && err == nil {
			y.Nam = string(n[0])
		}
	}
	return y.Nam
}

func (y *YeeLight) SetName(name string) error {
	y.counter++

	var props []string
	props = append(props, name)
	req := &request{Id: y.counter, Method: "set_name", Params: props}
	res, err := y.command(req)
	if err != nil || res == nil {
		return utils.ErrIOError
	}

	if res.Error != nil {
		log.Println("receive error from command: ", res.Error.Message)
		return ErrInvalidCommand
	}

	y.Nam = name

	return nil
}

func (y *YeeLight) String() string {
	return fmt.Sprintf("{type:\"%s\", id:%x, ip:\"%s\", name:\"%s\", model:\"%s\", ver:%d, support:\"%v\", power:\"%v\", "+
		"bright:%d, colorMode:%d, temp:%d, rgb:%d, hue:%d, sat:%d}", y.deviceType, y.id, y.ip,
		y.Nam, y.model, y.FwVer, y.support, y.Power, y.Bright, y.ColorMode, y.ColorTemp, y.Rgb, y.Hue, y.Sat)
}

func (y *YeeLight) Close() error {
	return y.conn.Close()
}

func (y *YeeLight) command(req *request) (*response, error) {
	if y.conn == nil {
		return nil, ErrNotStarted
	}

	request, err := json.Marshal(&req)
	request = append(request, '\r')
	request = append(request, '\n')
	if err != nil {
		return nil, err
	}

	log.Println("lan request", strings.Trim(string(request), "\r\n"))

	if _, err = y.conn.Write(request); err != nil {
		return nil, utils.ErrIOError
	}

	resp, err := bufio.NewReader(bufio.NewReader(y.conn)).ReadBytes('\n')
	if err != nil {
		return nil, utils.ErrIOError
	}

	log.Println("lan response", strings.Trim(string(resp), "\r\n"))

	var result response
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (y *YeeLight) GetProp(props []Prop) ([]Prop, error) {
	y.counter++
	req := &request{Id: y.counter, Method: "get_prop", Params: props}
	res, err := y.command(req)
	if err != nil || res == nil {
		return nil, utils.ErrIOError
	}

	if res.Error != nil {
		log.Println("receive error from command: ", res.Error.Message)
		return nil, ErrInvalidCommand
	}

	var ps []Prop
	for _, p := range res.Result {
		ps = append(ps, Prop(p))
	}

	return ps, nil
}

func (y *YeeLight) SetCtAbx(temp int, effect Effect, duration int) error {
	y.counter++

	var props []interface{}
	if temp >= 1700 && temp <= 6500 {
		props = append(props, temp)
	} else {
		return ErrWrongParameter
	}
	props = append(props, effect)
	props = append(props, duration*1000)
	req := &request{Id: y.counter, Method: "set_ct_abx", Params: props}
	res, err := y.command(req)
	if err != nil || res == nil {
		return utils.ErrIOError
	}

	if res.Error != nil {
		log.Println("receive error from command: ", res.Error.Message)
		return ErrInvalidCommand
	}

	// result:["ok"]
	return nil
}

func (y *YeeLight) BgSetCtAbx(temp int, effect Effect, duration int) error {
	y.counter++

	var props []interface{}
	if temp >= 1700 && temp <= 6500 {
		props = append(props, temp)
	} else {
		return ErrWrongParameter
	}
	props = append(props, effect)
	props = append(props, duration*1000)
	req := &request{Id: y.counter, Method: "bg_set_ct_abx", Params: props}
	res, err := y.command(req)
	if err != nil || res == nil {
		return utils.ErrIOError
	}

	if res.Error != nil {
		log.Println("receive error from command: ", res.Error.Message)
		return ErrInvalidCommand
	}

	//result:["ok"]
	return nil
}

// rgb : r << 8 + g << 8 + b
func (y *YeeLight) SetRgb(rgb int, effect Effect, duration int) error {

	y.counter++

	var props []interface{}
	props = append(props, rgb)
	props = append(props, effect)
	props = append(props, duration*1000)
	req := &request{Id: y.counter, Method: "set_rgb", Params: props}
	res, err := y.command(req)
	if err != nil || res == nil {
		return utils.ErrIOError
	}

	if res.Error != nil {
		log.Println("receive error from command: ", res.Error.Message)
		return ErrInvalidCommand
	}

	//result:["ok"]
	return nil
}

func (y *YeeLight) BgSetRgb(rgb int, effect Effect, duration int) error {
	y.counter++

	var props []interface{}
	props = append(props, rgb)
	props = append(props, effect)
	props = append(props, duration*1000)
	req := &request{Id: y.counter, Method: "bg_set_rgb", Params: props}
	res, err := y.command(req)
	if err != nil || res == nil {
		return utils.ErrIOError
	}

	if res.Error != nil {
		log.Println("receive error from command: ", res.Error.Message)
		return ErrInvalidCommand
	}

	//result:["ok"]
	return nil
}

func (y *YeeLight) SetHsv(hue int, sat int, effect Effect, duration int) error {
	y.counter++

	var props []interface{}
	if hue >= 0 && hue <= 359 {
		props = append(props, hue)
	} else {
		return ErrWrongParameter
	}

	if sat >= 0 && sat <= 100 {
		props = append(props, sat)
	} else {
		return ErrWrongParameter
	}

	props = append(props, effect)
	props = append(props, duration*1000)
	req := &request{Id: y.counter, Method: "set_hsv", Params: props}
	res, err := y.command(req)
	if err != nil || res == nil {
		return utils.ErrIOError
	}

	if res.Error != nil {
		log.Println("receive error from command: ", res.Error.Message)
		return ErrInvalidCommand
	}

	// result:["ok"]
	return nil
}

func (y *YeeLight) BgSetHsv(hue int, sat int, effect Effect, duration int) error {
	y.counter++

	var props []interface{}
	if hue >= 0 && hue <= 359 {
		props = append(props, hue)
	} else {
		return ErrWrongParameter
	}

	if sat >= 0 && sat <= 100 {
		props = append(props, sat)
	} else {
		return ErrWrongParameter
	}

	props = append(props, effect)
	props = append(props, duration*1000)
	req := &request{Id: y.counter, Method: "bg_set_hsv", Params: props}
	res, err := y.command(req)
	if err != nil || res == nil {
		return utils.ErrIOError
	}

	if res.Error != nil {
		log.Println("receive error from command: ", res.Error.Message)
		return ErrInvalidCommand
	}

	//result:["ok"]
	return nil
}

func (y *YeeLight) SetBrightness(brightness int, effect Effect, duration int) error {
	y.counter++

	var props []interface{}
	if brightness > 0 && brightness <= 100 {
		props = append(props, brightness)
	} else {
		return ErrWrongParameter
	}
	props = append(props, effect)
	props = append(props, duration*1000)
	req := &request{Id: y.counter, Method: "set_bright", Params: props}
	res, err := y.command(req)
	if err != nil || res == nil {
		return utils.ErrIOError
	}

	if res.Error != nil {
		log.Println("receive error from command: ", res.Error.Message)
		return ErrInvalidCommand
	}

	// result:["ok"]
	return nil
}

func (y *YeeLight) BgSetBrightness(brightness int, effect Effect, duration int) error {
	y.counter++

	var props []interface{}
	if brightness > 0 && brightness <= 100 {
		props = append(props, brightness)
	} else {
		return ErrWrongParameter
	}
	props = append(props, effect)
	props = append(props, duration*1000)
	req := &request{Id: y.counter, Method: "bg_set_bright", Params: props}
	res, err := y.command(req)
	if err != nil || res == nil {
		return utils.ErrIOError
	}

	if res.Error != nil {
		log.Println("receive error from command: ", res.Error.Message)
		return ErrInvalidCommand
	}

	// result:["ok"]
	return nil
}

func (y *YeeLight) SetPower(power bool, effect Effect, duration int, mode int) error {
	y.counter++

	var props []interface{}
	if power {
		props = append(props, "on")
	} else {
		props = append(props, "off")
	}
	props = append(props, effect)
	props = append(props, duration*1000)
	props = append(props, mode)
	req := &request{Id: y.counter, Method: "set_power", Params: props}
	res, err := y.command(req)
	if err != nil || res == nil {
		return utils.ErrIOError
	}

	if res.Error != nil {
		log.Println("receive error from command: ", res.Error.Message)
		return ErrInvalidCommand
	}

	// result:["ok"]
	return nil
}

func (y *YeeLight) BgSetPower(power bool, effect Effect, duration int, mode int) error {
	y.counter++

	var props []interface{}
	if power {
		props = append(props, "on")
	} else {
		props = append(props, "off")
	}
	props = append(props, effect)
	props = append(props, duration*1000)
	props = append(props, mode)
	req := &request{Id: y.counter, Method: "bg_set_power", Params: props}
	res, err := y.command(req)
	if err != nil || res == nil {
		return utils.ErrIOError
	}

	if res.Error != nil {
		log.Println("receive error from command: ", res.Error.Message)
		return ErrInvalidCommand
	}

	// result:["ok"]
	return nil
}

func (y *YeeLight) Toggle() (*bool, error) {
	y.counter++

	req := &request{Id: y.counter, Method: "toggle", Params: []int{}}
	res, err := y.command(req)
	if err != nil || res == nil {
		return nil, utils.ErrIOError
	}

	if res.Error != nil {
		log.Println("receive error from command: ", res.Error.Message)
		return nil, ErrInvalidCommand
	}

	if res.Params != nil {
		var stat bool
		switch res.Params.Power {
		case "ok":
			stat = true
		}
		return &stat, nil
	}

	return nil, utils.ErrIOError
}

func (y *YeeLight) BgToggle() (*bool, error) {
	y.counter++

	req := &request{Id: y.counter, Method: "bg_toggle", Params: []int{}}
	res, err := y.command(req)
	if err != nil || res == nil {
		return nil, utils.ErrIOError
	}

	if res.Params != nil {
		var stat bool
		switch res.Params.Power {
		case "ok":
			stat = true
		}
		return &stat, nil
	}

	return nil, utils.ErrIOError
}

func (y *YeeLight) SetDefault() error {
	y.counter++

	req := &request{Id: y.counter, Method: "set_default", Params: []int{}}
	res, err := y.command(req)
	if err != nil || res == nil {
		return utils.ErrIOError
	}

	if res.Error != nil {
		log.Println("receive error from command: ", res.Error.Message)
		return ErrInvalidCommand
	}

	// result:["ok"]
	return nil
}

func (y *YeeLight) BgSetDefault() error {
	y.counter++

	req := &request{Id: y.counter, Method: "bg_set_default", Params: []int{}}
	res, err := y.command(req)
	if err != nil || res == nil {
		return utils.ErrIOError
	}

	if res.Error != nil {
		log.Println("receive error from command: ", res.Error.Message)
		return ErrInvalidCommand
	}

	// result:["ok"]
	return nil
}

func (y *YeeLight) StartColorFlow(count int, action StateAfterStop, expr string) error {
	y.counter++

	var props []interface{}
	props = append(props, count)
	props = append(props, action)
	props = append(props, expr)
	req := &request{Id: y.counter, Method: "start_cf", Params: props}
	res, err := y.command(req)
	if err != nil || res == nil {
		return utils.ErrIOError
	}

	if res.Error != nil {
		log.Println("receive error from command: ", res.Error.Message)
		return ErrInvalidCommand
	}

	return nil
}

func (y *YeeLight) BgStartColorFlow(count int, action StateAfterStop, expr string) error {
	y.counter++

	var props []interface{}
	props = append(props, count)
	props = append(props, action)
	props = append(props, expr)
	req := &request{Id: y.counter, Method: "bg_start_cf", Params: props}
	res, err := y.command(req)
	if err != nil || res == nil {
		return utils.ErrIOError
	}

	if res.Error != nil {
		log.Println("receive error from command: ", res.Error.Message)
		return ErrInvalidCommand
	}

	return nil
}

func (y *YeeLight) StopColorFlow() error {
	y.counter++

	req := &request{Id: y.counter, Method: "stop_cf", Params: []int{}}
	res, err := y.command(req)
	if err != nil || res == nil {
		return utils.ErrIOError
	}

	if res.Error != nil {
		log.Println("receive error from command: ", res.Error.Message)
		return ErrInvalidCommand
	}

	// result:["ok"]
	return nil
}

func (y *YeeLight) BgStopColorFlow() error {
	y.counter++

	req := &request{Id: y.counter, Method: "bg_stop_cf", Params: []int{}}
	res, err := y.command(req)
	if err != nil || res == nil {
		return utils.ErrIOError
	}

	if res.Error != nil {
		log.Println("receive error from command: ", res.Error.Message)
		return ErrInvalidCommand
	}

	// result:["ok"]
	return nil
}

func (y *YeeLight) SetScene(class SceneClass, value1 int, value2 int, value3 int) error {
	y.counter++

	var props []interface{}
	props = append(props, class)
	props = append(props, value1)
	props = append(props, value2)
	props = append(props, value3)
	req := &request{Id: y.counter, Method: "set_scene", Params: props}
	res, err := y.command(req)
	if err != nil || res == nil {
		return utils.ErrIOError
	}

	if res.Error != nil {
		log.Println("receive error from command: ", res.Error.Message)
		return ErrInvalidCommand
	}

	//result:["ok"]
	return nil
}

func (y *YeeLight) BgSetScene(class SceneClass, value1 int, value2 int, value3 int) error {
	y.counter++

	var props []interface{}
	props = append(props, class)
	props = append(props, value1)
	props = append(props, value2)
	props = append(props, value3)
	req := &request{Id: y.counter, Method: "bg_set_scene", Params: props}
	res, err := y.command(req)
	if err != nil || res == nil {
		return utils.ErrIOError
	}

	if res.Error != nil {
		log.Println("receive error from command: ", res.Error.Message)
		return ErrInvalidCommand
	}

	//result:["ok"]
	return nil
}

func (y *YeeLight) SetAdjust(action Action, prop AdjustProperty) error {
	y.counter++

	if prop == Color {
		action = Circle
	}

	var props []interface{}
	props = append(props, action)
	props = append(props, prop)
	req := &request{Id: y.counter, Method: "set_adjust", Params: props}
	res, err := y.command(req)
	if err != nil || res == nil {
		return utils.ErrIOError
	}

	if res.Error != nil {
		log.Println("receive error from command: ", res.Error.Message)
		return ErrInvalidCommand
	}

	//result["ok"]
	return nil
}

func (y *YeeLight) BgSetAdjust(action Action, prop AdjustProperty) error {
	y.counter++

	if prop == Color {
		action = Circle
	}

	var props []interface{}
	props = append(props, action)
	props = append(props, prop)
	req := &request{Id: y.counter, Method: "bg_set_adjust", Params: props}
	res, err := y.command(req)
	if err != nil || res == nil {
		return utils.ErrIOError
	}

	if res.Error != nil {
		log.Println("receive error from command: ", res.Error.Message)
		return ErrInvalidCommand
	}

	//result:["ok"]
	return nil
}

func (y *YeeLight) SetMusic(action bool, host string, port int) error {
	y.counter++

	var props []interface{}

	if action {
		props = append(props, 1)
	} else {
		props = append(props, 0)
	}
	props = append(props, host)
	props = append(props, port)
	req := &request{Id: y.counter, Method: "set_music", Params: props}
	res, err := y.command(req)
	if err != nil || res == nil {
		return utils.ErrIOError
	}

	if res.Error != nil {
		log.Println("receive error from command: ", res.Error.Message)
		return ErrInvalidCommand
	}

	// result:["ok"]
	return nil
}

func (y *YeeLight) DevToggle() error {
	y.counter++

	req := &request{Id: y.counter, Method: "dev_toggle", Params: []int{}}
	res, err := y.command(req)
	if err != nil || res == nil {
		return utils.ErrIOError
	}

	if res.Error != nil {
		log.Println("receive error from command: ", res.Error.Message)
		return ErrInvalidCommand
	}

	// result: ["ok"]
	return nil
}

func (y *YeeLight) AdjustBright(percentage int, duration int) error {
	y.counter++

	var props []int

	if percentage >= -100 && percentage <= 100 {
		props = append(props, percentage)
	} else {
		return ErrWrongParameter
	}
	props = append(props, duration*1000)
	req := &request{Id: y.counter, Method: "adjust_bright", Params: props}
	res, err := y.command(req)
	if err != nil || res == nil {
		return utils.ErrIOError
	}

	if res.Error != nil {
		log.Println("receive error from command: ", res.Error.Message)
		return ErrInvalidCommand
	}

	// result:["ok"]
	return nil
}

func (y *YeeLight) BgAdjustBright(percentage int, duration int) error {
	y.counter++

	var props []int

	if percentage >= -100 && percentage <= 100 {
		props = append(props, percentage)
	} else {
		return ErrWrongParameter
	}
	props = append(props, duration*1000)
	req := &request{Id: y.counter, Method: "bg_adjust_bright", Params: props}
	res, err := y.command(req)
	if err != nil || res == nil {
		return utils.ErrIOError
	}

	if res.Error != nil {
		log.Println("receive error from command: ", res.Error.Message)
		return ErrInvalidCommand
	}

	// result:["ok"]
	return nil
}

func (y *YeeLight) AdjustColorTemp(percentage int, duration int) error {
	y.counter++

	var props []int

	if percentage >= -100 && percentage <= 100 {
		props = append(props, percentage)
	} else {
		return ErrWrongParameter
	}
	props = append(props, duration*1000)
	req := &request{Id: y.counter, Method: "adjust_ct", Params: props}
	res, err := y.command(req)
	if err != nil || res == nil {
		return utils.ErrIOError
	}

	if res.Error != nil {
		log.Println("receive error from command: ", res.Error.Message)
		return ErrInvalidCommand
	}

	// result:["ok"]
	return nil
}

func (y *YeeLight) BgAdjustColorTemp(percentage int, duration int) error {
	y.counter++

	var props []int

	if percentage >= -100 && percentage <= 100 {
		props = append(props, percentage)
	} else {
		return ErrWrongParameter
	}
	props = append(props, duration*1000)
	req := &request{Id: y.counter, Method: "bg_adjust_ct", Params: props}
	res, err := y.command(req)
	if err != nil || res == nil {
		return utils.ErrIOError
	}

	if res.Error != nil {
		log.Println("receive error from command: ", res.Error.Message)
		return ErrInvalidCommand
	}

	//result:["ok"]
	return nil
}

func (y *YeeLight) AdjustColor(percentage int, duration int) error {
	y.counter++

	var props []int

	if percentage >= -100 && percentage <= 100 {
		props = append(props, percentage)
	} else {
		return ErrWrongParameter
	}
	props = append(props, duration*1000)
	req := &request{Id: y.counter, Method: "adjust_color", Params: props}
	res, err := y.command(req)
	if err != nil || res == nil {
		return utils.ErrIOError
	}

	if res.Error != nil {
		log.Println("receive error from command: ", res.Error.Message)
		return ErrInvalidCommand
	}

	//result:["ok"]
	return nil
}

func (y *YeeLight) BgAdjustColor(percentage int, duration int) error {
	y.counter++

	var props []int

	if percentage >= -100 && percentage <= 100 {
		props = append(props, percentage)
	} else {
		return ErrWrongParameter
	}
	props = append(props, duration*1000)
	req := &request{Id: y.counter, Method: "bg_adjust_color", Params: props}
	res, err := y.command(req)
	if err != nil || res == nil {
		return utils.ErrIOError
	}

	if res.Error != nil {
		log.Println("receive error from command: ", res.Error.Message)
		return ErrInvalidCommand
	}

	//result:["ok"]
	return nil
}
