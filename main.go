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

	log.Println("starting manager_yeelight ...")

	// connect to mqtt
	log.Println("try connect to mqtt")
	var mqttid uint16 = 1
	mqtt, err := transport.Connect(*srv, *clientid, uint16(*keepalive), *login, *pass, *debug)
	if err != nil {
		panic("can't connect to mqtt server " + err.Error())
	}
	go mqtt.Start()

	log.Println("subscribe to managed topics")
	sp := packet.NewSubscribe()
	sp.Id = mqttid
	sp.Topics = []packet.SubscribePayload{{Topic: "yeelight/#", QoS: 1}}
	mqtt.Sendout <- sp
	mqttid++

	devices := make(map[string]*device.Device)
	updates := make(chan *device.Device)

	// fetch command data from mqtt server
	go func() {
		for {
			for pkt := range mqtt.Broker {
				if pkt.Type() == packet.PUBLISH {
					var d device.Device

					topics := strings.Split(pkt.(*packet.PublishPacket).Topic, "/")
					if err := json.Unmarshal([]byte(pkt.(*packet.PublishPacket).Payload), &d); err == nil {
						if d.Cmd != "" {
							// process command from hub
							if devices[topics[1]] != nil {
								if res, err := devices[topics[1]].Run(d.Cmd, d.Value1, d.Value2, d.Value3); err != nil {
									log.Println("error run command", err)
								} else {
									// send results back to mqtt server
									log.Println("run command result", res)
								}
							}
						} else {
							// restore devices from mqtt
							if devices[topics[1]] == nil {
								log.Printf("restore device %s", topics[1])
								d.Type = d.CheckDevice(d.Model)
								d.ConvertSupport()
								devices[topics[1]] = &d

								// we don't know real ip address of this device, so we can't start it.
							}
						}
					} else {
						// unmarshall error (probably should not been here
						log.Println("unmarshall error ", err)
					}
				}
			}
		}
	}()

	time.Sleep(3 * time.Second)

	// receive updates from devices
	go func() {
		for {
			for d := range updates {
				log.Printf("device %s state changed: %s", d.Id, d.String())
				p := packet.NewPublish()
				p.Id = mqttid
				p.Topic = fmt.Sprintf("yeelight/%s", d.Id)
				p.QoS = 1
				p.Payload = d.String()
				mqtt.Sendout <- p
				mqttid++
			}
		}
	}()

	// start device discovery
	discovery := ssdp.NewDiscovery()
	for d := range discovery.Reporter {
		if devices[d.Id] == nil {
			log.Printf("new device: %s", d.String())
			d.Start(updates)

			p := packet.NewPublish()
			p.Id = mqttid
			p.Topic = fmt.Sprintf("yeelight/%s", d.Id)
			p.QoS = 1
			p.Payload = fmt.Sprintf(`{"id":"%s","model":"%s","name":"%s","ver":%d}`, d.Id, d.Model, d.Name, d.Version)
			p.Retain = true
			mqtt.Sendout <- p
			mqttid++
		} else {
			if !devices[d.Id].IsStarted() {
				devices[d.Id] = d
				d.Start(updates)
			}
		}
	}
}
