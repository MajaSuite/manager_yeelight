package ssdp

import (
	"bytes"
	"fmt"
	"golang.org/x/net/ipv4"
	"log"
	"manager_yeelight/device"
	"net"
	"strings"
	"time"
)

const (
	ssdpDiscoveryAddr    = "239.255.255.250"
	ssdpDiscoveryPort    = 1982
	ssdpDiscoveryHeader  = "ssdp:discover"
	ssdpDiscoveryService = "wifi_bulb"
)

type Discovery struct {
	Reporter chan *device.Device
}

func NewDiscovery() *Discovery {
	log.Println("Start discovery")

	d := &Discovery{
		Reporter: make(chan *device.Device),
	}

	conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
	if err != nil {
		panic(err)
	}

	pkt := ipv4.NewPacketConn(conn)

	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", ssdpDiscoveryAddr, ssdpDiscoveryPort))
	if err != nil {
		panic(err)
	}

	mconn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		return nil
	}
	mpkt := ipv4.NewPacketConn(mconn)

	ifaces, err := net.Interfaces()
	if err != nil {
		log.Println("error getting interfaces", err)
		return nil
	}
	for _, iface := range ifaces {
		ifAddr, err := iface.Addrs()
		if err != nil {
			return nil
		}

		for _, ifaddr := range ifAddr {
			// skip ipv6
			if strings.Contains(ifaddr.String(), "::") {
				continue
			}

			// skip localhost
			if strings.Contains(ifaddr.String(), "127.0.0.1") {
				continue
			}

			log.Printf("listen on interface %s", iface.Name)
			if err := pkt.JoinGroup(&iface, &net.UDPAddr{IP: net.ParseIP(ssdpDiscoveryAddr)}); err != nil {
				log.Println("join error", err)
			}

			go d.notifier(mpkt, addr)
			go d.Listener(mpkt)
		}
	}

	return d
}

func (d *Discovery) Listener(mconn *ipv4.PacketConn) error {
	for {
		buffer := make([]byte, 0x2048)
		_, _, addr, err := mconn.ReadFrom(buffer)
		if err != nil {
			panic(err)
		}

		buf := bytes.NewBuffer(buffer)
		if _, err := buf.ReadBytes('\n'); err != nil {
			continue
		}

		device := device.NewDevice(buf, addr)

		d.Reporter <- device
	}
}

func (d *Discovery) notifier(conn *ipv4.PacketConn, addr *net.UDPAddr) {
	for {
		log.Printf("send discovery request")

		request := fmt.Sprintf("M-SEARCH * HTTP/1.1\r\nHOST: %s:%d\r\nMAN: \"%s\"\r\nST: %s\r\n",
			ssdpDiscoveryAddr, ssdpDiscoveryPort, ssdpDiscoveryHeader, ssdpDiscoveryService)

		if _, err := conn.WriteTo([]byte(request), nil, addr); err != nil {
			panic(err)
		}
		time.Sleep(60 * time.Second)
	}
}
