package internal

import (
	"errors"
	"golang.zx2c4.com/wireguard/wgctrl"
	"log"
	"net"
	"time"
)

type Peer struct {
	PublicKey         string
	AllowedIps        []net.IPNet
	LastHandshakeTime time.Time
}

type Device struct {
	Ip        net.IP
	Name      string
	PublicKey string
}

func GetPeers(client *wgctrl.Client, deviceName string) []Peer {
	device, err := client.Device(deviceName)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	var peers []Peer
	for _, peer := range device.Peers {
		peers = append(peers, Peer{
			PublicKey:         peer.PublicKey.String(),
			AllowedIps:        peer.AllowedIPs,
			LastHandshakeTime: peer.LastHandshakeTime,
		})

	}

	return peers
}

func GetDeviceInfo(client *wgctrl.Client, deviceName string) (*Device, error) {
	device, err := client.Device(deviceName)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	ip, err := getIpOfInterface(device.Name)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &Device{
		Name:      device.Name,
		PublicKey: device.PublicKey.String(),
		Ip:        ip,
	}, nil

}

func getIpOfInterface(name string) (net.IP, error) {
	iface, err := net.InterfaceByName(name)

	if err != nil {
		return nil, err
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}

	addr, ok := addrs[0].(*net.IPNet)
	if !ok {
		return nil, errors.New("cant get ip from interface")
	}

	return addr.IP, nil
}
