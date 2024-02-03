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

func env(key string, defaultValue string) (value string) {
	value = os.Getenv(key)
	if value == "" {
		value = defaultValue
	}

	return
}

func getLocalAddr() (string, error) {
	conn, err := net.Dial("udp", "255.255.255.255:0")
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

func queryPort(network string, ip string, port int, timeout time.Duration) error {
	conn, err := net.DialTimeout(network, net.JoinHostPort(ip, strconv.Itoa(port)), timeout)
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

	network := env("NETWORK", "tcp")

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

	for range time.Tick(timeout) {
		localAddr, err := getLocalAddr()
		if err != nil {
			log.Println(err)
			continue
		}

		// log.Println(localAddr)

		err = queryPort(network, localAddr, port, dialTimeout)
		if err == nil {
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
