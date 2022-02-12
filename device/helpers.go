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

func readResponse(debug bool, conn net.Conn, id int) ([]string, error) {
	conn.SetReadDeadline(time.Now().Add(time.Second * 10))
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
		// todo we receive notification message, update client
		return readResponse(debug, conn, id)
	}

	return nil, ErrWrongParameter
}

func CreateDevice(debug bool, model string, id string, ip string, name string, support string, power bool, ver int,
	bright int, mode int, temp int, rgb int, hue int, sat int) Device {
	var dev Device

	switch CheckDevice(model) {
	case RGB_DEVICE:
		dev = NewRgbDevice(debug, model, id, ip, name, support, power, ver, bright, mode, temp, rgb, hue, sat)
	case LIGHT_DEVICE:
		dev = NewBWDevice(debug, model, id, ip, name, support, power, ver, bright, mode, temp)
	case AMBILIGHT_DEVICE:
		//
	}

	return dev
}
