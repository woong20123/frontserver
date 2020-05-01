package examserverlogic

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// ServerConfig is
type ServerConfig struct {
	// ServerIP is examplefrontserver's IP info
	ServerIP string `json:"serverip"`
	// ServerPort is examplefrontserver's Port info
	ServerPort int `json:"serverport"`

	// ServerMode is How the server behaves
	// "main" : 단일 서버로 동작합니다. 로직의 처리가 이곳에서 완료 됩니다.
	// "front" : front 서버로 동작합니다. 패킷을 목적지 서버로 전달하는 역활을 합니다.
	ServerMode string `json:"servermode"`

	// BackEndChatServerAddr is back-end Chat Server's address
	BackEndChatServerIP   string `json:"backendChatServerIP"`
	BackEndChatServerPort int    `json:"backendChatServerPort"`
}

// MakeExampleConfig is
func MakeExampleConfig() {
	bReadExampleFile := true
	exampleFilePath := "./Example.json"

	// ExampleConfig이 있는지 확인합니다.
	if _, err := os.Stat(exampleFilePath); os.IsNotExist(err) {
		bReadExampleFile = false
	}

	svrConfigWrite := ServerConfig{}
	svrConfigWrite.ServerIP = "127.0.0.1"
	svrConfigWrite.ServerPort = 20224
	svrConfigWrite.ServerMode = "main"
	svrConfigWrite.BackEndChatServerIP = "127.0.0.1"
	svrConfigWrite.BackEndChatServerPort = 20414

	// ExampleConfig 파일이 있다면 비교해서 다른점이 없다면 리턴
	if bReadExampleFile == true {
		svrConfigRead := ServerConfig{}
		b, err := ioutil.ReadFile(exampleFilePath)
		if err == nil {
			json.Unmarshal(b, &svrConfigRead)

			if svrConfigRead == svrConfigWrite {
				return
			}
		}
	}

	doc, _ := json.Marshal(svrConfigWrite)

	err := ioutil.WriteFile(exampleFilePath, doc, os.FileMode(0644))
	if err != nil {
		Logger().Println(err)
	}
}

// ConfigManager is
type ConfigManager struct {
	svrConfig *ServerConfig
}

func newConfigMgr() *ConfigManager {
	configMgr := new(ConfigManager)
	configMgr.svrConfig = new(ServerConfig)
	return configMgr
}

// ReadConfig is
func (cm *ConfigManager) ReadConfig(filename string) {

	// Config이 있는지 확인합니다
	current_path, _ := os.Getwd()
	target_path := current_path + filename

	if _, err := os.Stat(target_path); os.IsNotExist(err) {
		Logger().Println("Not Find ConfigPath = ", target_path)
		Logger().Println("config path Replace  to \"./Example.json\" ")
		target_path = "./Example.json"
	}

	b, err := ioutil.ReadFile(target_path)
	if err != nil {
		Logger().Println(err)
	}

	err = json.Unmarshal(b, cm.svrConfig)
	if err != nil {
		Logger().Println(err)
	}
}

// ServerConfig is
func (cm *ConfigManager) ServerConfig() *ServerConfig {
	return cm.svrConfig
}
