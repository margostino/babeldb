package common

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

func (s *Server) GetAddress() string {
	return s.Host + ":" + s.Port
}

func GetServerConfig(path string) []*Server {
	return GetConfig(path).Servers
}

func GetClientConfig(path string) []*Client {
	return GetConfig(path).Clients
}

func GetStrikerConfig(path string) *Striker {
	return GetConfig(path).Striker
}

func GetConfig(file string) *Configuration {
	var configuration Configuration
	ymlFile, err := ioutil.ReadFile(file)

	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	ymlFile = []byte(os.ExpandEnv(string(ymlFile)))
	err = yaml.Unmarshal(ymlFile, &configuration)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return &configuration
}
