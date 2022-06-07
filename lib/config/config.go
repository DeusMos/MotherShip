package config

import (
	"errors"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

var (
	// Get the desired configuration value.
	Get struct {
		Server struct {
			Host string `yaml:"host"`
			Port uint16 `yaml:"port"`
		} `yaml:"server"`

		Machine struct {
			ID int `yaml:"id"`
		} `yaml:"machine"`
		Reporting struct {
			Customer  string   `yaml:"customer"`
			Plant     string   `yaml:"plant"`
			ID        string   `yaml:"id"`
			Endpoints []string `yaml:"endpoints"`
		} `yaml:"reporting"`
		Logging struct {
			EnableConsole bool   `yaml:"enableConsole"`
			EnableHistory bool   `yaml:"enableHistory"`
			HistoryLimit  uint32 `yaml:"historyLimit"`
		} `yaml:"logging"`
	}
)

func Save(filename string) (err error) {
	my_yaml, err := yaml.Marshal(Get)
	if err != nil {
		return
	}
	ioutil.WriteFile(filename, my_yaml, 0)
	return
}

// Load the configuration file.
func Load(filename string) (err error) {
	// Read contents from the configuration file.
	var contents []byte
	if contents, err = ioutil.ReadFile(filename); nil != err {
		return
	}

	// Parse contents into 'Get' configuration.
	if err = yaml.Unmarshal(contents, &Get); nil != err {
		return
	}

	if Get.Reporting.ID == "" {
		err = errors.New("reporting:id must be set in config.yml")
		return
	}
	return
}
func init() {
	setDefaults()
}

func setDefaults() {
	Get.Server.Host = "woodfam.us"
	Get.Server.Port = 8888
	Get.Reporting.Plant = "Walla"
	Get.Reporting.Customer = "key"
	Get.Reporting.ID = "dummyID"
	Get.Reporting.Endpoints = []string{"https://fini.key.net/g6stats/api/stats"}
	Get.Logging.EnableConsole = true
	Get.Logging.EnableHistory = true
	Get.Logging.HistoryLimit = 10000

}
