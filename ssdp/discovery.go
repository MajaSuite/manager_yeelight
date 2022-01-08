package ssdp

import (
	"bytes"
	"fmt"
	"golang.org/x/net/ipv4"
	"io"
	"log"
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
	Reporter chan map[string]string
}

func NewDiscovery() *Discovery {
	log.Println("Start discovery")

	d := &Discovery{
		Reporter: make(chan map[string]string),
	}

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

		for _, addr := range ifAddr {
			// skip ipv6
			if strings.Contains(addr.String(), "::") {
				continue
			}

			// skip localhost
			if strings.Contains(addr.String(), "127.0.0.1") {
				continue
			}

			ifAddr := strings.Split(addr.String(), "/")

			go d.Listener(iface, ifAddr[0])
		}
	}

	return d
}

func (d *Discovery) Listener(iface net.Interface, ifAddr string) error {
	log.Printf("listen on %s interface with ip %s", iface.Name, ifAddr)
	addr, err := net.ResolveUDPAddr("udp4",
		fmt.Sprintf("%s:%d", ssdpDiscoveryAddr, ssdpDiscoveryPort))
	if err != nil {
		return err
	}

	conn, err := net.ListenPacket("udp4", fmt.Sprintf("%s:%d", "0.0.0.0", ssdpDiscoveryPort))
	if err != nil {
		return err
		panic(err)
	}

	mcConn := ipv4.NewPacketConn(conn)
	if err := mcConn.JoinGroup(&iface, addr); err != nil {
		panic(err)
	}
	if err := mcConn.SetMulticastLoopback(false); err != nil {
		panic(err)
	}
	if err := mcConn.SetControlMessage(ipv4.FlagDst, true); err != nil {
		panic(err)
	}

	go d.notifier(mcConn, addr)

	for {
		buffer := make([]byte, 0x2048)
		_, _, addr, err := mcConn.ReadFrom(buffer)
		if err != nil {
			panic(err)
		}

		buf := bytes.NewBuffer(buffer)
		if _, err := buf.ReadBytes('\n'); err != nil {
			continue
		}

		lines := make(map[string]string)
		for {
			chunk, err := buf.ReadBytes('\n')
			if err == io.EOF {
				break
			}

			line := strings.Split(string(chunk), ": ")
			if len(line) >= 2 {
				lines[line[0]] = strings.Trim(string(chunk)[len(line[0])+1:], " \r\n")
			}
		}
		lines["ip"] = strings.Split(addr.String(), ":")[0]

		d.Reporter <- lines

	}
}

func (d *Discovery) notifier(conn *ipv4.PacketConn, addr *net.UDPAddr) {
	for {
		log.Printf("send notification")

		request := fmt.Sprintf("M-SEARCH * HTTP/1.1\r\nHOST: %s:%d\r\nMAN: \"%s\"\r\nST: %s\r\n",
			ssdpDiscoveryAddr, ssdpDiscoveryPort, ssdpDiscoveryHeader, ssdpDiscoveryService)

		if _, err := conn.WriteTo([]byte(request), nil, addr); err != nil {
			panic(err)
		}
		time.Sleep(60*time.Second)
	}
}
