package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/MajaSuite/mqtt/packet"
	"github.com/MajaSuite/mqtt/transport"
	"log"
	"manager_yeelight/device"
	"manager_yeelight/ssdp"
	"manager_yeelight/utils"
	"strings"
	"time"
)

var (
	srv       = flag.String("mqtt", "127.0.0.1:1883", "mqtt server address")
	clientid  = flag.String("clientid", "yeelight-1", "client id for mqtt server")
	keepalive = flag.Int("keepalive", 30, "keepalive timeout for mqtt server")
	login     = flag.String("login", "", "login string for mqtt server")
	pass      = flag.String("pass", "", "password string for mqtt server")
	debug     = flag.Bool("debug", false, "print debuging hex dumps")
)

func main() {
	flag.Parse()

	// connect to mqtt
	log.Println("try connect to mqtt")
	var id uint16 = 1
	mqtt, err := transport.Connect(*srv, *clientid, uint16(*keepalive), *login, *pass, *debug)
	if err != nil {
		panic("can't connect to mqtt server " + err.Error())
	}
	go mqtt.Start()

	log.Println("subscribe to managed topics")
	sp := packet.NewSubscribe()
	sp.Id = 1
	sp.Topics = []packet.SubscribePayload{{Topic: "yeelight/#", QoS: 1}}
	mqtt.Sendout <- sp

	devices := make(map[uint32]device.Device)

	// receive updates from devices
	update := make(chan device.Device)
	go func() {
		for {
			for dev := range update {
				log.Println("update device", dev)
				p := packet.NewPublish()
				id++
				p.Id = id
				p.Topic = fmt.Sprintf("yeelight/%x", dev.Id())
				p.QoS = 1
				p.Payload = dev.String()
				mqtt.Sendout <- p
			}
		}
	}()

	// fetch command data from mqtt server
	go func() {
		for {
			for pkt := range mqtt.Broker {
				if pkt.Type() == packet.PUBLISH {
					var dev device.DeviceJson
					topics := strings.Split(pkt.(*packet.PublishPacket).Topic, "/")
					if err := json.Unmarshal([]byte(pkt.(*packet.PublishPacket).Payload), &dev); err == nil {
						dev.Id = utils.ConvertHex(dev.StringId)
						if dev.Cmd != "" {
							// process command from hub
							devliceId := utils.ConvertHex(topics[1])
							log.Printf("run command %s (%s %s %s) for device: %x", dev.Cmd, dev.Value1, dev.Value2,
								dev.Value3, devliceId)

							if devices[devliceId] != nil {
								devices[devliceId].Cmd(dev.Cmd, dev.Value1, dev.Value2, dev.Value3)
							}
						} else {
							// restore devices from mqtt
							if devices[dev.Id] == nil {
								log.Printf("restore device %x", dev.Id)
								d := device.NewYeeLight(dev.Id, dev.Model, dev.Name, dev.Support)
								d.Version = dev.Version
								d.Power = dev.Power
								d.Bright = dev.Bright
								d.ColorMode = dev.ColorMode
								d.ColorTemp = dev.ColorTemp
								d.Rgb = dev.Rgb
								d.Hue = dev.Hue
								d.Sat = dev.Sat
								devices[dev.Id] = d
							}
						}
					} else {
						log.Println(err)
					}
				}
			}
		}
	}()

	time.Sleep(3 * time.Second)

	// start device discovery
	discovery := ssdp.NewDiscovery()

	for lines := range discovery.Reporter {
		deviceId := utils.ConvertHex(lines["id"])
		if devices[deviceId] == nil {
			dev := device.CreateDevice(lines)
			if dev != nil {
				dev.Start(lines["ip"], update)
				devices[deviceId] = dev

				log.Printf("new device: %s", dev.String())

				p := packet.NewPublish()
				id++
				p.Id = id
				p.Topic = fmt.Sprintf("yeelight/%x", dev.Id())
				p.QoS = 1
				p.Payload = dev.String()
				p.Retain = true
				mqtt.Sendout <- p
			} else {
				log.Printf("unknown device: %s (%s) at %s", lines["type"], lines["model"], lines["ip"])
			}
		} else {
			dev := devices[deviceId]
			if err := dev.Start(lines["ip"], update); err != device.ErrAlreadyStarted {
				log.Printf("start device: %s", dev.String())
				//if err := device.UpdateDevice(dev, lines); err != nil {
				//	p := packet.NewPublish()
				//	id++
				//	p.Id = id
				//	p.Topic = fmt.Sprintf("yeelight/%x", dev.Id())
				//	p.QoS = 1
				//	p.Payload = dev.String()
				//	mqtt.Sendout <- p
				//}
			}
		}
	}
}
