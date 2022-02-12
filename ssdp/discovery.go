package ssdp

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"manager_yeelight/device"
	"manager_yeelight/utils"
	"net"
	"strings"
	"time"

	"golang.org/x/net/ipv4"
)

const (
	ssdpDiscoveryAddr    = "239.255.255.250"
	ssdpDiscoveryPort    = 1982
	ssdpDiscoveryHeader  = "ssdp:discover"
	ssdpDiscoveryService = "wifi_bulb"
)

func NewDiscovery(debug bool, discovery chan device.Device) error {
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
	if err != nil {
		log.Println("error listen udp", err)
		panic(err)
	}
	packetConn := ipv4.NewPacketConn(conn)

	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", ssdpDiscoveryAddr, ssdpDiscoveryPort))
	if err != nil {
		log.Println("error resolv udp address", err)
		panic(err)
	}

	multicast, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		log.Println("error listen multicast udp", err)
		panic(err)
	}

	packetMulticast := ipv4.NewPacketConn(multicast)

	ifaces, err := net.Interfaces()
	for _, iface := range ifaces {
		if err := packetConn.JoinGroup(&iface, &net.UDPAddr{IP: net.ParseIP(ssdpDiscoveryAddr)}); err != nil {
			if debug {
				log.Printf("join error on %s : %s", iface.Name, err)
			}
		} else {
			if debug {
				log.Printf("listen on interface %s", iface.Name)
			}
		}
	}

	go func() {
		for {
			request := fmt.Sprintf("M-SEARCH * HTTP/1.1\r\nHOST: %s:%d\r\nMAN: \"%s\"\r\nST: %s\r\n",
				ssdpDiscoveryAddr, ssdpDiscoveryPort, ssdpDiscoveryHeader, ssdpDiscoveryService)

			if _, err := packetMulticast.WriteTo([]byte(request), nil, addr); err != nil {
				log.Println("error write discovery request", err)
				panic(err)
			}

			time.Sleep(time.Second * 60)
		}
	}()

	for {
		buffer := make([]byte, 2048)
		_, _, src, err := packetMulticast.ReadFrom(buffer)
		if err != nil {
			log.Println("error read multicast packet", err)
			panic(err)
		}

		var ip = strings.Split(src.String(), ":")[0]
		var id, model, name, support string
		var ver, bright, mode, temp, rgb, hue, sat int
		var power bool
		buf := bytes.NewBuffer(buffer)
		for {
			chunk, err := buf.ReadBytes('\n')
			if err == io.EOF {
				break
			}

			line := strings.Split(string(chunk), ": ")
			if len(line) >= 2 {
				v := strings.Trim(string(chunk)[len(line[0])+1:], " \r\n")
				switch line[0] {
				case "id":
					id = strings.TrimLeft(strings.Replace(v, "0x", "", -1), "0")
				case "model":
					model = v
				case "name":
					name = v
				case "fw_ver":
					ver = utils.ConvertInt(v)
				case "power":
					power = utils.ConvertBool(v)
				case "bright":
					bright = utils.ConvertInt(v)
				case "color_mode":
					mode = utils.ConvertInt(v)
				case "ct":
					temp = utils.ConvertInt(v)
				case "rgb":
					rgb = utils.ConvertInt(v)
				case "hue":
					hue = utils.ConvertInt(v)
				case "sat":
					sat = utils.ConvertInt(v)
				case "support":
					support = v
				}
			}
		}

		discovery <- device.CreateDevice(debug, model, id, ip, name, support, power, ver, bright, mode, temp, rgb, hue, sat)
	}
}
