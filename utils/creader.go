package utils

import (
	"fmt"
	"os"

	"gopkg.in/ini.v1"
)

func ReadConfig(fileName string) map[string]string {
	cfg, err := ini.Load(fileName)
	if err != nil {
		fmt.Printf("Failed to read config file: %v", err)
		os.Exit(1)
	}
	res := make(map[string]string)
	keys := cfg.Section("").KeyStrings()
	for _, j := range keys {
		res[j] = cfg.Section("").Key(j).String()
	}
	return res
}
