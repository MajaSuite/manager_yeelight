package main

import (
	"flag"
	"github.com/MajaSuite/mqtt/packet"
	"github.com/MajaSuite/mqtt/transport"
	"log"
	"manager_yeelight/device"
	"strconv"
	"strings"

	//"manager_yeelight/xiaomi"
	"manager_yeelight/ssdp"
)

var (
	cert      = flag.String("cert", "server.crt", "path to server certificate")
	key       = flag.String("key", "server.key", "path to server private key")
	srv       = flag.String("mqtt", "127.0.0.1:1883", "mqtt server address")
	clientid  = flag.String("clientid", "yeelight-1", "client id for mqtt server")
	keepalive = flag.Int("keepalive", 60, "keepalive timeout for mqtt server")
	login     = flag.String("login", "", "login string for mqtt server")
	pass      = flag.String("pass", "", "password string for mqtt server")
	debug     = flag.Bool("debug", false, "print debuging hex dumps")
)

func convertHex(v string) uint32 {
	cleaned := strings.Replace(v, "0x", "", -1)
	res, _ := strconv.ParseUint(cleaned, 16, 64)
	return uint32(res)
}

func convertBool(v string) bool {
	if v == "off" {
		return false
	}
	return true
}

func converArray(v string) []string {
	return nil
}

func convertInt(v string) int {
	if i, err := strconv.Atoi(v); err != nil {
		return 0
	} else {
		return i
	}
}

func main() {
	flag.Parse()

	// connect to mqtt
	log.Println("try connect to mqtt")
	mqtt, err := transport.Connect(*srv, *clientid, uint16(*keepalive), *login, *pass, *debug)
	if err != nil {
		panic("can't connect to mqtt server " + err.Error())
	}

	log.Println("subscribe to managed topics")
	sp := packet.NewSubscribe()
	sp.Topics = []packet.SubscribePayload{
		{Topic: "yeelight", QoS: 1},
		{Topic: "yeelight/#", QoS: 1},
	}
	mqtt.Subscribe(sp)

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
	devices := make(map[string]device.Device)

	for lines := range discovery.Reporter {
		if devices[lines["ip"]] == nil {
			dev := device.NewYeeLight(convertHex(lines["id"]), nil, lines["model"], convertInt(lines["fw_ver"]),
				converArray(lines["support"]), convertBool(lines["power"]), convertInt(lines["bright"]),
				convertInt(lines["color_mode"]), convertInt(lines["ct"]), convertInt(lines["rgb"]), convertInt(lines["hue"]),
				convertInt(lines["sat"]), lines["name"])
			dev.SetIP(lines["ip"])
			devices[lines["ip"]] = dev

			log.Printf("new device:%s\n", dev.String())
		}
	}

	// finish := make(chan os.Signal, 1)
	// signal.Notify(finish, syscall.SIGINT, syscall.SIGTERM)
	// <-finish
}
