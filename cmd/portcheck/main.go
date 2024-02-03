package main

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func env(key string, defaultValue string) (value string) {
	value = os.Getenv(key)
	if value == "" {
		value = defaultValue
	}

	return
}

func getLocalAddr() (string, error) {
	conn, err := net.Dial("tcp", "255.255.255.255:0")
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

		err = queryPort(network, localAddr, port, dialTimeout)
		if err == nil {
			continue
		}

		log.Println(err)

		args := []string{
			"ps",
			"--format",
			"json",
			"--filter",
			"label=io.github.eiqnepm.portcheck.enable=true",
		}

		output, err := exec.Command("docker", args...).Output()
		if err != nil {
			log.Println(err)
			continue
		}

		ids := []string{"restart"}
		for _, line := range strings.Split(string(output), "\n") {
			var container struct {
				ID string `json:"ID"`
			}

			if err := json.Unmarshal([]byte(line), &container); err != nil {
				log.Println(line)
				log.Println(err)
				continue
			}

			log.Println(container.ID)
			ids = append(ids, container.ID)
		}

		if _, err := exec.Command("docker", ids...).Output(); err != nil {
			log.Println(err)
			continue
		}
	}
}
