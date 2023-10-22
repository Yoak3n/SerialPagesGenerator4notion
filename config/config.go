package config

import (
	"b2n3/package/util"
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DatabaseID string `json:"database_id"`
	Token      string `json:"token"`
}

var Conf *Config

const configPath = "data/.config"

func init() {
	util.CreateDirNotExists(configPath)
	Conf = &Config{}
}

func LoadOptions() []string {
	dir, err := os.ReadDir(configPath)
	if err != nil {
		return nil
	}
	var files []string
	for _, fp := range dir {
		if !fp.IsDir() {
			files = append(files, fp.Name())
		}
	}
	return files
}

func LoadConfigFile(name string) error {
	buf, err := os.ReadFile(fmt.Sprintf("%s/%s.json", configPath, name))
	if err != nil {
		return err
	}
	err = json.Unmarshal(buf, &Conf)
	if err != nil {
		return err
	}
	return nil
}

func ModifyConfiguration(config *Config, name string) error {
	Conf = config
	err := writeConfigFile(name)
	if err != nil {
		return err
	}
	return nil
}

func CreateConfiguration(config *Config, name string) error {
	Conf.DatabaseID = config.DatabaseID
	Conf.Token = config.Token
	err := writeConfigFile(name)
	if err != nil {
		return err
	}
	return nil
}

func writeConfigFile(name string) error {
	fp, err := os.Create(fmt.Sprintf("%s/%s.json", configPath, name))
	if err != nil {
		return err
	}
	defer fp.Close()
	_, err = fp.Write([]byte(fmt.Sprintf(`{"database_id":"%s","token":"%s"}`, Conf.DatabaseID, Conf.Token)))
	return err
}
