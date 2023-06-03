package qbit

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	u "net/url"
)

type Session struct {
	url    u.URL
	client http.Client
}

func Login(url u.URL, username string, password string) (session Session, err error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return
	}

	session = Session{
		url: url,
		client: http.Client{
			Jar: jar,
		},
	}

	data := u.Values{
		"username": {username},
		"password": {password},
	}

	session.url.Path = "api/v2/auth/login"
	resp, err := session.client.PostForm(session.url.String(), data)
	if err != nil {
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	if resp.StatusCode == 200 {
		return
	}

	return Session{}, errors.New(resp.Status)
}

func (session Session) Logout() (err error) {
	session.url.Path = "api/v2/auth/logout"
	resp, err := session.client.Post(session.url.String(), "", nil)
	if err != nil {
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	if resp.StatusCode == 200 {
		return
	}

	return errors.New(resp.Status)
}

func (session Session) SetPreference(preference string, value any) (err error) {
	j, err := json.Marshal(map[string]any{preference: value})
	if err != nil {
		return
	}

	data := u.Values{
		"json": {string(j)},
	}

	session.url.Path = "api/v2/app/setPreferences"
	resp, err := session.client.PostForm(session.url.String(), data)
	if err != nil {
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	if resp.StatusCode == 200 {
		log.Println(string(j))
		return
	}

	return errors.New(resp.Status)
}
