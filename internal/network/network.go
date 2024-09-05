package network

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func CheckPort(port int) (success bool, err error) {
	resp, err := http.Get("https://portcheck.transmissionbt.com/" + url.PathEscape(strconv.Itoa(port)))
	if err != nil {
		return success, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println(err)
		}
	}()

	if resp.StatusCode != 200 {
		return success, errors.New("port check status: " + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return success, err
	}

	switch string(body) {
	case "0":
		return false, nil
	case "1":
		return true, nil
	default:
		return success, err
	}
}
