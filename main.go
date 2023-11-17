package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sleepycrew/wg_exporter/internal"
	"golang.zx2c4.com/wireguard/wgctrl"
)

var (
	peersGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "wg_peers",
		Help: "The total number of processed events",
	}, []string{"device_name", "device_public_key", "peer_public_key", "peer_allowed_ips"})
)

func recordWgPeers() {
	go func() {
		for {
			client, err := wgctrl.New()
			devices, err := client.Devices()
			if err != nil {
				log.Fatal("could not create wgctrl client")
			}
			for _, wgDevice := range devices {
				_, err := internal.GetDeviceInfo(client, wgDevice.Name)
				if err != nil {
					log.Println("failed to get device info")
				}
				peers := internal.GetPeers(client, wgDevice.Name)
				if peers == nil {
					log.Println("failed to get peers")
					continue
				}

				for _, peer := range peers {
					var allowedIps []string

					for _, ip := range peer.AllowedIps {
						allowedIps = append(allowedIps, ip.String())
					}
					peerCount := float64(len(allowedIps))

					peersGauge.WithLabelValues(wgDevice.Name, wgDevice.PublicKey.String(), peer.PublicKey, strings.Join(allowedIps, ",")).Set(peerCount)
				}
			}

			time.Sleep(15 * time.Second)
		}
	}()
}

func main() {
	go recordWgPeers()
	http.Handle("/metrics", promhttp.Handler())
	log.Println("listening on :2112")
	http.ListenAndServe(":2112", nil)
}
