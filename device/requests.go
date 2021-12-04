package device

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidCommand = errors.New("invalid command")
	ErrWrongParameter = errors.New("invalid parameter")
	ErrNotStarted     = errors.New("device not started")
)

type Effect string

const (
	Sudden Effect = "sudden"
	Smooth Effect = "smooth"
)

type Action string

const (
	Increase Action = "increase"
	Decrease Action = "decrease"
	Circle   Action = "circle"
)

type StateAfterStop int

const (
	StayPrevious StateAfterStop = 0
	StayLast     StateAfterStop = 1
	StayOff      StateAfterStop = 2
)

type SceneClass string

const (
	// means turn on and change the smart LED to specified color
	SceneColor SceneClass = "color"
	// means turn on and change the smart LED to specified color and brightness.
	SceneHsv SceneClass = "hsv"
	// means turn on and change the smart LED to specified temp and brightness.
	SceneCT SceneClass = "ct"
	//means turn on and start a color flow in specified fashion.
	SceneCF SceneClass = "cf"
	// means turn on led to specified brightness and start a sleep timer to turn off the
	// light after the specified minutes.
	SceneAutoDelayOff SceneClass = "auto_delay_off"
)

type AdjustProperty string

const (
	Brightness AdjustProperty = "bright"
	Temp       AdjustProperty = "ct"
	Color      AdjustProperty = "color"
)

type Prop string

const (
	not_exist      Prop = "not_exist"
	power          Prop = "power"          //on: smart LED is turned on / off: smart LED is turned off
	bright         Prop = "bright"         //Brightness percentage. Range 1 ~ 100
	ct             Prop = "ct"             //Color temperature. Range 1700 ~ 6500(k)
	rgb            Prop = "rgb"            //Color. Range 1 ~ 16777215
	hue            Prop = "hue"            //Hue. Range 0 ~ 359
	sat            Prop = "sat"            //Saturation. Range 0 ~ 100
	color_mode     Prop = "color_mode"     //1: rgb mode / 2: color temperature mode / 3: hsv mode
	flowing        Prop = "flowing"        //0: no flow is running / 1:color flow is running
	delayoff       Prop = "delayoff"       //The remaining time of a sleep timer. Range 1 ~ 60 (minutes)
	flow_params    Prop = "flow_params"    //Current flow parameters (only meaningful when 'flowing' is 1)
	music_on       Prop = "music_on"       //1: Music mode is on / 0: Music mode is off
	name           Prop = "name"           //The name of the device set by “set_name” command
	bg_power       Prop = "bg_power"       //Background light power status
	bg_flowing     Prop = "bg_flowing"     //Background light is flowing
	bg_flow_params Prop = "bg_flow_params" // Current flow parameters of background light
	bg_ct          Prop = "bg_ct"          //Color temperature of background light
	bg_lmode       Prop = "bg_lmode"       //1: rgb mode / 2: color temperature mode / 3: hsv mode
	bg_bright      Prop = "bg_bright"      //Brightness percentage of background light
	bg_rgb         Prop = "bg_rgb"         //Color of background light
	bg_hue         Prop = "bg_hue"         //Hue of background light
	bg_sat         Prop = "bg_sat"         //Saturation of background light
	nl_br          Prop = "nl_br"          //Brightness of night mode light
	active_mode    Prop = "active_mode"    //0: daylight mode / 1: moonlight mode (ceiling light only)
)

type request struct {
	Id     uint32      `json:"id"`
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

func (r request) String() string {
	return fmt.Sprintf("{id:%d, method:%s, params:%s}", r.Id, r.Method, r.Params)
}

type response struct {
	Id     uint32      `json:"id"`
	Result []string    `json:"result,omitempty"`
	Params *respParams `json:"params,omitempty"`
	Error  *respError  `json:"error,omitempty"`
}

func (r response) String() string {
	var err string
	if r.Error != nil {
		err = r.Error.String()
	}

	var params string
	if r.Params != nil {
		err = r.Params.String()
	}
	return fmt.Sprintf("{id:%d, result:%s, params:%s, errror:%s}", r.Id, r.Result, params, err)
}

type respError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (r respError) String() string {
	return fmt.Sprintf("{code:%d, message:%s}", r.Code, r.Message)
}

type respParams struct {
	Power string `json:"power,omitempty"`
}

func (r respParams) String() string {
	var sb strings.Builder
	sb.WriteString("{")
	if r.Power != "" {
		sb.WriteString("power:" + r.Power + " ")
	}
	sb.WriteString("}")

	return sb.String()
}
