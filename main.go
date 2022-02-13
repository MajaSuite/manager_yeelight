package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/MajaSuite/mqtt/client"
	"github.com/MajaSuite/mqtt/packet"
	"log"
	"manager_yeelight/device"
	"manager_yeelight/ssdp"
	"manager_yeelight/utils"
	"strings"
)

var (
	srv       = flag.String("mqtt", "127.0.0.1:1883", "mqtt server address")
	clientid  = flag.String("clientid", "yeelight-1", "client id for mqtt server")
	keepalive = flag.Int("keepalive", 30, "keepalive timeout for mqtt server")
	login     = flag.String("login", "", "login string for mqtt server")
	pass      = flag.String("pass", "", "password string for mqtt server")
	debug     = flag.Bool("debug", false, "print debuging hex dumps")
	qos       = flag.Int("qos", 0, "qos to send/receive from mqtt")
)

func main() {
	flag.Parse()

	log.Println("starting manager_yeelight")

	// connect to mqtt
	log.Println("try connect to mqtt")
	var mqttId uint16 = 1
	mqtt, err := client.Connect(*srv, *clientid, uint16(*keepalive), false, *login, *pass /* *debug */, false)
	if err != nil {
		panic("can't connect to mqtt server ")
	}

	log.Println("subscribe to managed topics")
	sp := packet.NewSubscribe()
	sp.Id = mqttId
	sp.Topics = []packet.SubscribePayload{{Topic: "yeelight/#", QoS: packet.QoS(*qos)}}
	mqtt.Send <- sp
	mqttId++

	log.Println("start yeelight discovery")
	discovery := make(chan device.Device)
	go ssdp.NewDiscovery(*debug, discovery)

	devices := make(map[string]device.Device)

	for {
		select {
		case pkt := <-mqtt.Receive:
			if pkt.Type() == packet.PUBLISH {
				topics := strings.Split(pkt.(*packet.PublishPacket).Topic, "/")
				var msg map[string]interface{}
				if err := json.Unmarshal([]byte(pkt.(*packet.PublishPacket).Payload), &msg); err == nil {
					dev := device.CreateDevice(*debug, utils.ConvertToString(msg, "model"), topics[1], "",
						utils.ConvertToString(msg, "name"), utils.ConvertToString(msg, "support"),
						utils.ConvertToBool(msg, "power"), utils.ConvertToInt(msg, "ver"),
						utils.ConvertToInt(msg, "bright"), utils.ConvertToInt(msg, "mode"),
						utils.ConvertToInt(msg, "temp"), utils.ConvertToInt(msg, "rgb"),
						utils.ConvertToInt(msg, "hue"), utils.ConvertToInt(msg, "sat"))

					if dev != nil {
						if devices[topics[1]] == nil {
							devices[dev.ID()] = dev

							if *debug {
								log.Println("new from mqtt", dev)
							}
						} else {
							p := packet.NewPublish()
							p.Id = mqttId
							p.Topic = fmt.Sprintf("yeelight/%s", dev.ID())
							p.QoS = packet.QoS(*qos)
							mqttId++

							method := utils.ConvertToString(msg, "method")
							if method != "" {
								log.Println("run direct command for", topics[1])
								switch dev.Type() {
								case device.LIGHT_DEVICE:
									err = devices[topics[1]].(*device.LightDevice).RunMethod(&dev, method,
										utils.ConvertToString(msg, "effect"), utils.ConvertToInt(msg, "duration"))
									p.Payload = devices[topics[1]].(*device.LightDevice).String()
								case device.RGB_DEVICE:
									err = devices[topics[1]].(*device.RgbDevice).RunMethod(&dev, method,
										utils.ConvertToString(msg, "effect"), utils.ConvertToInt(msg, "duration"))
									p.Payload = devices[topics[1]].(*device.RgbDevice).String()
								case device.AMBILIGHT_DEVICE:
									err = devices[topics[1]].(*device.AmbilightDevice).RunMethod(&dev, method,
										utils.ConvertToString(msg, "effect"), utils.ConvertToInt(msg, "duration"))
									p.Payload = devices[topics[1]].(*device.AmbilightDevice).String()
								default:
									log.Println("unknown device, can't run method")
								}

								if err == nil {
									mqtt.Send <- p
								} else {
									log.Println("error apply changes", err)
								}
							} else {
								switch dev.Type() {
								case device.LIGHT_DEVICE:
									err = devices[topics[1]].(*device.LightDevice).CompareAndUpdate(&dev)
									p.Payload = devices[topics[1]].(*device.LightDevice).String()
								case device.RGB_DEVICE:
									err = devices[topics[1]].(*device.RgbDevice).CompareAndUpdate(&dev)
									p.Payload = devices[topics[1]].(*device.RgbDevice).String()
								case device.AMBILIGHT_DEVICE:
									err = devices[topics[1]].(*device.AmbilightDevice).CompareAndUpdate(&dev)
									p.Payload = devices[topics[1]].(*device.AmbilightDevice).String()
								default:
									log.Println("unknown device, can't change state")
								}
							}
						}
					}
				} else {
					log.Println("mqtt message unmarshall error ", err)
				}
			}

		case dev := <-discovery:
			if dev == nil || dev.Type() == device.NO_TYPE { // undefined (...and unsupported) device
				continue
			}

			p := packet.NewPublish()
			p.Id = mqttId
			p.Topic = fmt.Sprintf("yeelight/%s", dev.ID())
			p.QoS = packet.QoS(*qos)
			mqttId++

			if devices[dev.ID()] == nil {
				if *debug {
					log.Println("new device from discovery", dev.Retain())
				}
				devices[dev.ID()] = dev
				p.Retain = true
				p.Payload = dev.Retain()
			} else {
				// update existing
				devices[dev.ID()] = dev
				p.Payload = dev.String()
			}

			mqtt.Send <- p
		}
	}
}
