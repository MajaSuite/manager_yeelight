package device

import (
	"fmt"
)

// Effects
const (
	Sudden string = "sudden"
	Smooth string = "smooth"
)

// Action
const (
	Increase string = "increase"
	Decrease string = "decrease"
	Circle   string = "circle"
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

type lanRequest struct {
	Id     int         `json:"id"`
	Method string      `json:"method"`
	Params interface{} `json:"params,omitempty"`
}

func (r lanRequest) String() string {
	return fmt.Sprintf("{id:%d, method:%s, params:%s}", r.Id, r.Method, r.Params)
}

type lanResponse struct {
	Id     int        `json:"id"`
	Result []string   `json:"result,omitempty"`
	Error  *respError `json:"error,omitempty"`
}

func (r lanResponse) String() string {
	var e string
	if r.Error != nil {
		e = fmt.Sprintf(",\"error\":%s", r.Error.String())
	}
	return fmt.Sprintf(`{"id":%d,"result":%v,"error":%s}`,
		r.Id, r.Result, e)
}

type respError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (r respError) String() string {
	return fmt.Sprintf("{code:%d, message:%s}", r.Code, r.Message)
}
