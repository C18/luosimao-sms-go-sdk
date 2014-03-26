package luosimao

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	AuthUsername = "api"
)

type Luosimao struct {
	ApiKey string
	Suffix string
}

type SendResult struct {
	Err int
	Msg string
}

type StatusResult struct {
	Error   int
	Deposit int `json:",string"`
}

func New(key string, suffix string) *Luosimao {
	return &Luosimao{key, suffix}
}

func (this *Luosimao) Send(mobile, content string) error {
	params := url.Values{}
	params.Add("mobile", mobile)
	params.Add("message", content+this.Suffix)

	buf := bytes.NewBuffer([]byte(params.Encode()))

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://sms-api.luosimao.com/v1/send.json", buf)
	if err != nil {
		return err
	}

	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(AuthUsername, this.ApiKey)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	result := &SendResult{}
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	if result.Err != 0 {
		return errors.New(fmt.Sprintf("%d", result.Err))
	}
	return nil
}

func (this *Luosimao) Status() (int, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://sms-api.luosimao.com/v1/status.json", nil)
	req.SetBasicAuth(AuthUsername, this.ApiKey)
	if err != nil {
		return 0, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return 0, err
	}

	result := &StatusResult{}
	if err := json.Unmarshal(body, result); err != nil {
		return 0, err
	}

	if result.Error != 0 {
		return 0, errors.New(fmt.Sprintf("%d", result.Error))
	}
	return result.Deposit, nil
}
