package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	"github.com/ChristianHering/GoEssentials"
	"github.com/libvirt/libvirt-go"
)

type Config struct {
	ModKeys []string
	Domains []Domain
}

type Domain struct {
	Name         string
	StartKey     string
	SuspendKey   string
	ForceQuitKey string
}

type InputEvent struct {
	Time  syscall.Timeval
	Type  uint16
	Code  uint16
	Value int32
}

var config = Config{
	ModKeys: []string{"KEY_LEFTCTRL", "KEY_RIGHTCTRL"},
	Domains: []Domain{
		Domain{
			Name:         "arch",
			StartKey:     "KEY_KP1",
			SuspendKey:   "KEY_KP2",
			ForceQuitKey: "KEY_KP3",
		},
		Domain{
			Name:         "win",
			StartKey:     "KEY_KP4",
			SuspendKey:   "KEY_KP5",
			ForceQuitKey: "KEY_KP6",
		},
		Domain{
			Name:         "gentoo",
			StartKey:     "KEY_KP7",
			SuspendKey:   "KEY_KP8",
			ForceQuitKey: "KEY_KP9",
		},
	},
}

var domains []*libvirt.Domain

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}

	err = os.MkdirAll(filepath.Join(homeDir, "/.config/aucht"), 0644)
	if err != nil {
		log.Fatalln(err)
	}

	err = goessentials.GetConfig(filepath.Join(homeDir, "/.config/aucht/config.json"), &config)
	if err != nil && err != goessentials.ErrorConfigFileUnset {
		log.Fatalln(err)
	} else if err != nil {
		log.Println(err)
	}

	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		log.Fatalln("Failed to establish a connection to libvirt")
	}
	defer conn.Close()

	for i := 0; i < len(config.Domains); i++ {
		d, err := conn.LookupDomainByName(config.Domains[i].Name)
		if err != nil {
			log.Println("Failed to lookup domain", config.Domains[i].Name)

			domains = append(domains, nil)
		}

		domains = append(domains, d)
	}

	resolved := "/dev/input/event%d"

	for i := 0; true; i++ {
		path := fmt.Sprintf("/sys/class/input/event%d/device/name", i)

		buff, err := ioutil.ReadFile(path)
		if err == os.ErrNotExist {
			panic("No keyboards found")
		} else if err != nil {
			continue
		}

		deviceName := strings.ToLower(string(buff))

		if strings.Contains(deviceName, "mouse") {
			continue
		} else if strings.Contains(deviceName, "keyboard") {
			resolved = fmt.Sprintf(resolved, i)

			break
		}
	}

	f, err := os.OpenFile(resolved, os.O_RDWR, os.ModeCharDevice)
	if err != nil {
		panic(err)
	}

	eventChan := make(chan InputEvent)

	go func() {
		for {
			buffer := make([]byte, int(unsafe.Sizeof(InputEvent{})))

			_, err = f.Read(buffer)
			if err != nil {
				panic(err)
			}

			event := &InputEvent{}

			err = binary.Read(bytes.NewBuffer(buffer), binary.LittleEndian, event)
			if err != nil {
				panic(err)
			}

			eventChan <- *event
		}
	}()

	var isModDown bool

	for {
		event := <-eventChan

		keypressEvent := keycodeMappings[event.Code]

		if keypressEvent == "KEY_RESERVED" || keypressEvent == "KEY_3" {
			continue
		}

		for i := 0; i < len(config.ModKeys); i++ {
			if keypressEvent == config.ModKeys[i] {
				isModDown = event.Value != 0

				break
			}
		}

		if isModDown && event.Value == 1 {
			for i := 0; i < len(config.Domains); i++ {
				switch keypressEvent {
				case config.Domains[i].StartKey:
					err = domains[i].Create()
					if err != nil {
						log.Println("Failed to start domain ", config.Domains[i].Name)
					}
				case config.Domains[i].SuspendKey:
					err = domains[i].PMSuspendForDuration(1, 0, 0)
					if err != nil {
						log.Println("Failed to suspend domain ", config.Domains[i].Name)
					}
				case config.Domains[i].ForceQuitKey:
					err = domains[i].Destroy()
					if err != nil {
						log.Println("Failed to force stop domain ", config.Domains[i].Name)
					}
				}
			}
		}
	}
}
