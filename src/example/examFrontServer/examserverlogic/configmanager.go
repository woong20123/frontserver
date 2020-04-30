package examserverlogic

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// ServerConfig is
type ServerConfig struct {
	ServerIP   string `json:"serverip"`
	ServerPort int    `json:"serverport"`

	// ServerMode is How the server behaves
	// "main" : 자체적으로 서버가 동작합니다.
	// "front" : front 서버로 동작합니다.
	ServerMode string `json:"servermode"`
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
	svrConfigWrite.ServerIP = "0.0.0.0"
	svrConfigWrite.ServerPort = 20224
	svrConfigWrite.ServerMode = "main"

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
func (cm *ConfigManager) ReadConfig(path string) {

	// Config이 있는지 확인합니다.
	if _, err := os.Stat(path); os.IsNotExist(err) {
		Logger().Println("Not Find ConfigPath = ", path)
		Logger().Println("config path Replace  to \"./Example.json\" ")
		path = "./Example.json"
	}

	b, err := ioutil.ReadFile(path)
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
