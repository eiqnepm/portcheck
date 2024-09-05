package main

import (
	"log"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/eiqnepm/portcheck/internal/deluge"
	"github.com/eiqnepm/portcheck/internal/network"
	"github.com/eiqnepm/portcheck/internal/qbit"
)

func env(key string, defaultValue string) (value string) {
	value = os.Getenv(key)
	if value == "" {
		value = defaultValue
	}

	return
}

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)

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

	clientWebUrl := url.URL{
		Scheme: clientWebScheme,
		Host:   net.JoinHostPort(clientWebHost, clientWebPort),
	}

	qbitUsername := env("CLIENT_USERNAME", "admin")
	clientPassword := env("CLIENT_PASSWORD", "adminadmin")
	if !strings.EqualFold(client, "qBittorrent") {
		clientPassword = env("CLIENT_PASSWORD", "deluge")
	}

	startupTimeout := func() time.Duration {
		startupTimeout, err := strconv.Atoi(env("STARTUP_TIMEOUT", "10"))
		if err != nil {
			log.Fatal(err)
		}

		return time.Duration(startupTimeout) * time.Second
	}()

	checkTimeout := func() time.Duration {
		timeout, err := strconv.Atoi(env("TIMEOUT", "300"))
		if err != nil {
			log.Fatal(err)
		}

		return time.Duration(timeout) * time.Second
	}()

	log.Printf("Port checks starting in %v second(s)\n", startupTimeout.Seconds())
	time.Sleep(startupTimeout)
	for c := time.Tick(checkTimeout); ; <-c {
		success, err := network.CheckPort(clientPort)
		if err != nil {
			log.Println(err)
			continue
		}

		if success {
			continue
		}

		log.Println("Port is closed")
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
