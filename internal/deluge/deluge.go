package deluge

import (
	"bytes"
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
	client *http.Client
}

func (sesh Session) post(unmarshmellowed any) error {
	marshmellowed, err := json.Marshal(unmarshmellowed)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", sesh.url.String(), bytes.NewBuffer(marshmellowed))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := sesh.client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Println(string(body))
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}
	return nil
}

func Login(url u.URL, password string) (Session, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return Session{}, err
	}
	sesh := Session{
		url: url,
		client: &http.Client{
			Jar: jar,
		},
	}
	sesh.url.Path = "json"
	err = sesh.post(map[string]interface{}{
		"method": "auth.login",
		"params": []string{
			password,
		},
		"id": 0,
	})
	return sesh, err
}

func (sesh Session) SetPreference(preference string, value any) error {
	return sesh.post(map[string]any{
		"method": "core.set_config",
		"params": []map[string]any{
			{preference: value},
		},
		"id": 0,
	})
}

func (sesh Session) Logout() error {
	return sesh.post(map[string]interface{}{
		"method": "auth.delete_session",
		"params": []string{},
		"id":     0,
	})
}
