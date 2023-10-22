package config

import (
	"b2n3/package/util"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
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

func LoadConfigFile(name string) (*Config, error) {
	buf, err := os.ReadFile(fmt.Sprintf("%s/%s.json", configPath, strings.TrimSuffix(name, ".json")))
	if err != nil {
		return nil, err
	}
	conf := &Config{}
	err = json.Unmarshal(buf, &conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func ModifyConfiguration(config *Config) error {
	Conf = config
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

func DeleteConfigurationFile(name string) error {
	err := os.Remove(fmt.Sprintf("%s/%s.json", configPath, strings.TrimSuffix(name, ".json")))
	if err != nil {
		return errors.New("删除配置文件失败")
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
