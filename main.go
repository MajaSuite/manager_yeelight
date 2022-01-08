package main

import (
	"flag"
	"fmt"
	"github.com/MajaSuite/mqtt/packet"
	"github.com/MajaSuite/mqtt/transport"
	"log"
	"manager_yeelight/device"
	"manager_yeelight/ssdp"
	"manager_yeelight/utils"
)

var (
	cert      = flag.String("cert", "server.crt", "path to server certificate")
	key       = flag.String("key", "server.key", "path to server private key")
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
	sp.Topics = []packet.SubscribePayload{{Topic: "yeelight/#", QoS: 1}}
	mqtt.Sendout <- sp

	// fetch command data from mqtt server
	go func() {
		for {
			for pkt := range mqtt.Broker {
				if pkt.Type() == packet.PUBLISH {
					log.Println("NEW PUBLISH: ", pkt)
				}
			}
		}
	}()

	// start device discovery
	discovery := ssdp.NewDiscovery()

	devices := make(map[string]*device.YeeLight)

	for lines := range discovery.Reporter {

		if devices[lines["ip"]] == nil {
			dev := device.NewYeeLight(utils.ConvertHex(lines["id"]),
				nil,
				lines["model"],
				utils.ConvertInt(lines["fw_ver"]),
				utils.ConverArray(lines["support"]),
				utils.ConvertBool(lines["power"]),
				utils.ConvertInt(lines["bright"]),
				utils.ConvertInt(lines["color_mode"]),
				utils.ConvertInt(lines["ct"]),
				utils.ConvertInt(lines["rgb"]),
				utils.ConvertInt(lines["hue"]),
				utils.ConvertInt(lines["sat"]),
				lines["name"])

			dev.SetIP(lines["ip"])
			devices[lines["ip"]] = dev

			log.Printf("new device: %s", dev.String())

			p := packet.NewPublish()
			id++
			p.Id = id
			p.Topic = fmt.Sprintf("yeelight/%x", dev.ID())
			p.QoS = 1
			p.Payload = dev.String()
			p.Retain = true
			mqtt.Sendout <- p

		} else {
			var change = false
			dev := devices[lines["ip"]]

			if dev.FwVer != utils.ConvertInt(lines["fw_ver"]) {
				change = true
				dev.FwVer = utils.ConvertInt(lines["fw_ver"])
			}
			if dev.Power != utils.ConvertBool(lines["power"]) {
				change = true
				dev.Power = utils.ConvertBool(lines["power"])
			}
			if dev.Bright != utils.ConvertInt(lines["bright"]) {
				change = true
				dev.Bright = utils.ConvertInt(lines["bright"])
			}
			if dev.ColorMode != utils.ConvertInt(lines["color_mode"]) {
				change = true
				dev.ColorMode = utils.ConvertInt(lines["color_mode"])
			}
			if dev.ColorTemp != utils.ConvertInt(lines["ct"]) {
				change = true
				dev.ColorTemp = utils.ConvertInt(lines["ct"])
			}
			if dev.Rgb != utils.ConvertInt(lines["rgb"]) {
				change = true
				dev.Rgb = utils.ConvertInt(lines["rgb"])
			}
			if dev.Hue != utils.ConvertInt(lines["hue"]) {
				change = true
				dev.Hue = utils.ConvertInt(lines["hue"])
			}
			if dev.Sat != utils.ConvertInt(lines["sat"]) {
				change = true
				dev.Sat = utils.ConvertInt(lines["sat"])
			}
			if dev.Nam != lines["name"] {
				change = true
				dev.Nam = lines["name"]
			}

			if change {
				p := packet.NewPublish()
				id++
				p.Id = id
				p.Topic = fmt.Sprintf("yeelight/%x", dev.ID())
				p.QoS = 1
				p.Payload = dev.String()
				mqtt.Sendout <- p
			}
		}
	}
}
