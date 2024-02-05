package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func env(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

func getLocalAddr(version int) (string, error) {
	var (
		network string
		address string
	)

	switch version {
	case 4:
		network = "udp4"
		address = "10.0.0.0:0"
	case 6:
		network = "udp6"
		address = "[fd00::]:0"
	}

	conn, err := net.Dial(network, address)
	if err != nil {
		return "", nil
	}

	defer func(conn *net.Conn) {
		err := (*conn).Close()
		if err != nil {
			log.Println(err)
		}
	}(&conn)

	host, _, err := net.SplitHostPort(conn.LocalAddr().String())
	if err != nil {
		return "", err
	}

	return host, nil
}

func checkPort(ip string, port int, timeout time.Duration) error {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, strconv.Itoa(port)), timeout)
	if err != nil {
		return err
	}

	if err := conn.Close(); err != nil {
		log.Println(err)
	}

	return nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	port, err := strconv.Atoi(env("PORT", "6881"))
	if err != nil {
		log.Fatal(err)
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

	version, err := strconv.Atoi(env("IP_VERSION", "4"))
	if err != nil {
		log.Fatal(err)
	}

	for range time.Tick(timeout) {
		localAddr, err := getLocalAddr(version)
		if err != nil {
			log.Println(err)
			continue
		}

		if err = checkPort(localAddr, port, dialTimeout); err == nil {
			continue
		}

		log.Println(err)

		func() {
			ctx := context.Background()
			cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
			if err != nil {
				log.Println(err)
				return
			}

			defer func(cli *client.Client) {
				if err := cli.Close(); err != nil {
					log.Println(err)
				}
			}(cli)

			filter := filters.NewArgs()
			filter.Add("label", "io.github.eiqnepm.portcheck.enable=true")
			containers, err := cli.ContainerList(ctx, container.ListOptions{Filters: filter})
			if err != nil {
				log.Println(err)
				return
			}

			for _, con := range containers {
				fmt.Println(con.ID)
				if err := cli.ContainerRestart(context.Background(), con.ID, container.StopOptions{}); err != nil {
					log.Println(err)
				}
			}
		}()
	}
}
