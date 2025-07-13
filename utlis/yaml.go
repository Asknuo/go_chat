package utlis

import (
	"gochat/global"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadYaml() ([]byte, error) {
	return os.ReadFile("config.yaml")

}
func SaveYmal() error {
	byteData, err := yaml.Marshal(global.Config)
	if err != nil {
		return err
	}
	return os.WriteFile("config.yaml", byteData, 0644)
}
