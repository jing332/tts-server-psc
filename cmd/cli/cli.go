package main

import (
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
	tts_server_plugin_sync_client "tts-server-plugin-sync-client"
)

var (
	path    string
	address string
	debug   bool
	preUi   bool
)

func main() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true})

	mSet := flag.NewFlagSet("", flag.ExitOnError)
	mSet.StringVar(&path, "path", "", "文件路径")
	mSet.StringVar(&address, "addr", "", "地址:端口")
	mSet.BoolVar(&debug, "debug", false, "push完成后自动debug")
	mSet.BoolVar(&preUi, "ui", false, "push完成后自动预览UI")

	err := mSet.Parse(os.Args[2:])
	if err != nil {
		log.Fatal(err)
	}

	action := os.Args[1]
	if action != "pull" && action != "push" && action != "watch" {
		log.Fatal("需要动作，必须为 push, pull, watch 之一")
	}

	if path == "" {
		log.Fatal("请指定 -path 参数")
	}

	if address == "" {
		log.Fatal("请指定 -addr 参数")
	}

	if action == "pull" {
		err := pull()
		if err != nil {
			log.Fatal(err)
		}
	} else if action == "push" {
		err := push()
		if err != nil {
			log.Fatal(err)
		}
	} else if action == "watch" {
		err := watch()
		if err != nil {
			log.Fatal(err)
		}
	}

}

func pull() error {
	client := tts_server_plugin_sync_client.SyncClient{Address: address}
	code, err := client.Pull()
	if err != nil {
		return fmt.Errorf("拉取代码失败：%v", err)
	}

	err = writeTxtFile(path, code)
	if err != nil {
		return fmt.Errorf("写入文件失败：%v", err)
	}
	return nil
}

func push() error {
	client := tts_server_plugin_sync_client.SyncClient{Address: address}
	code, err := readTxtFile(path)
	if err != nil {
		return fmt.Errorf("读取文件失败：%v", err)
	}
	if code == "" {
		return fmt.Errorf("代码为空：%v", path)
	}
	err = client.Push(code)
	if err != nil {
		return err
	}

	if debug {
		err := client.ActionDebug()
		if err != nil {
			return err
		}
	} else if preUi {
		err := client.ActionUI()
		if err != nil {
			return err
		}
	}

	return nil
}

//goland:noinspection GoUnhandledErrorResult
func watch() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("创建文件监听器失败：%v", err)
	}
	defer watcher.Close()

	err = watcher.Add(path)
	if err != nil {
		return fmt.Errorf("添加文件到监听器失败：%v", err)
	}

	var cancel func()
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				continue
			}

			if event.Has(fsnotify.Write) {
				if cancel != nil {
					cancel()
				}
				cancel = delay(func() {
					log.Infoln("push...")
					if err := push(); err != nil {
						log.Warnf("推送失败：%v\n", err)
					}
				})
			}
		}
	}

}

func delay(block func()) (cancel func()) {
	t := time.AfterFunc(time.Millisecond*500, func() {
		block()
	})

	return func() {
		t.Stop()
	}
}

func readTxtFile(path string) (string, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(bytes), err
}

func writeTxtFile(path, txt string) error {
	err := os.WriteFile(path, []byte(txt), os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
