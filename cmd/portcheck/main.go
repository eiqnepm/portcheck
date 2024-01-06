package main

import (
	"log"
	"net"
	u "net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Eiqnepm/portcheck/internal/deluge"
	"github.com/Eiqnepm/portcheck/internal/network"
	"github.com/Eiqnepm/portcheck/internal/qbit"
)

func env(key string, defaultValue string) (value string) {
	value = os.Getenv(key)
	if value == "" {
		value = defaultValue
	}

	return
}

func main() {
	log.SetFlags(log.LstdFlags)

	client := env("CLIENT", "qBittorrent")

	clientPort, err := strconv.Atoi(env("CLIENT_PORT", "6881"))
	if err != nil {
		log.Fatal(err)
	}

	clientWebScheme := env("CLIENT_WEBUI_SCHEME", "http")
	clientWebHost := env("CLIENT_WEBUI_HOST", "localhost")
	clientWebPort := env("CLIENT_WEBUI_PORT", "8080")
	if !strings.EqualFold(client, "qBittorrent") {
		clientWebPort = env("CLIENT_WEBUI_PORT", "8112")
	}
	clientWebUrl := u.URL{
		Scheme: clientWebScheme,
		Host:   net.JoinHostPort(clientWebHost, clientWebPort),
	}

	qbitUsername := env("CLIENT_USERNAME", "admin")
	clientPassword := env("CLIENT_PASSWORD", "adminadmin")
	if !strings.EqualFold(client, "qBittorrent") {
		clientPassword = env("CLIENT_PASSWORD", "deluge")
	}
	t, err := strconv.Atoi(env("TIMEOUT", "300"))
	if err != nil {
		log.Fatal(err)
	}

	timeout := time.Duration(t) * time.Second
	t, err = strconv.Atoi(env("DIAL_TIMEOUT", "5"))
	if err != nil {
		log.Fatal(err)
	}

	dialTimeout := time.Duration(t) * time.Second

	firstLoop := true
	for {
		if !firstLoop {
			time.Sleep(timeout)
		}

		firstLoop = false

		outboundIp, err := network.GetOutboundIP()
		if err != nil {
			log.Println(err)
			continue
		}

		err = network.QueryPort(outboundIp, clientPort, dialTimeout)
		if err == nil {
			continue
		}

		log.Println(err)

		if !strings.EqualFold(client, "qBittorrent") {
			func() {
				sesh, err := deluge.Login(clientWebUrl, clientPassword)
				if err != nil {
					log.Println(err)
					return
				}

				defer func(sesh deluge.Session) {
					err := sesh.Logout()
					if err != nil {
						log.Println(err)
					}
				}(sesh)

				err = sesh.SetPreference("listen_ports", []int{0})
				if err != nil {
					log.Println(err)
					return
				}

				err = sesh.SetPreference("listen_ports", []int{clientPort})
				if err != nil {
					log.Println(err)
					return
				}
			}()

			continue
		}

		func() {
			session, err := qbit.Login(clientWebUrl, qbitUsername, clientPassword)
			if err != nil {
				log.Println(err)
				return
			}

			defer func(session qbit.Session) {
				err := session.Logout()
				if err != nil {
					log.Println(err)
				}
			}(session)

			err = session.SetPreference("listen_port", 0)
			if err != nil {
				log.Println(err)
				return
			}

			err = session.SetPreference("listen_port", clientPort)
			if err != nil {
				log.Println(err)
				return
			}
		}()
	}
}
