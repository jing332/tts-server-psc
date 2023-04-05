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
	response, err := http.Get("http://" + s.Address + "/api/sync/pull")
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
	response, err := http.Post("http://"+s.Address+"/api/sync/push", "text/javascript", strings.NewReader(code))
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return errors.New("HTTP状态码不等于OK：" + response.Status)
	}

	return nil
}

func (s *SyncClient) ActionDebug() error {
	return s.Action("debug")
}

func (s *SyncClient) ActionUI() error {
	return s.Action("ui")
}

func (s *SyncClient) Action(name string) error {
	response, err := http.Get("http://" + s.Address + "/api/sync/action?action=" + name)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return errors.New("HTTP状态码不等于OK：" + response.Status)
	}
	return nil
}
