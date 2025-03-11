package config

import (
	"encoding/json"
	"log"
	"os"
)

type configType struct {
	Token  string `json:"Token"`
	Prefix string `json:"Prefix"`
}

var Config configType

func ReadFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("Unable to read the file ", err)
		return err
	}

	err = json.Unmarshal(data, &Config)
	if err != nil {
		log.Fatal("Unable to unmarshal the file ", err)
		return err
	}

	return nil
}
