package tts_server_plugin_sync_client

import (
	"errors"
	"io"
	"net/http"
	"strings"
)

type SyncClient struct {
	Address string
}

//goland:noinspection GoUnhandledErrorResult
func (s *SyncClient) Pull() (string, error) {
	response, err := http.Get("http://" + s.Address + "/api/plugin/pull")
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", errors.New("HTTP状态码不等于OK：" + response.Status)
	}

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (s *SyncClient) Push(code string) error {
	response, err := http.Post("http://"+s.Address+"/api/plugin/push", "text/javascript", strings.NewReader(code))
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return errors.New("HTTP状态码不等于OK：" + response.Status)
	}

	return nil
}

func (s *SyncClient) ActionDebug() error {
	response, err := http.Get("http://" + s.Address + "/api/plugin/action-debug")
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return errors.New("HTTP状态码不等于OK：" + response.Status)
	}
	return nil
}

func (s *SyncClient) ActionUI() error {
	response, err := http.Get("http://" + s.Address + "/api/plugin/action-ui")
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return errors.New("HTTP状态码不等于OK：" + response.Status)
	}
	return nil
}
