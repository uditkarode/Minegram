package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"os"
)

func readConfig(fileName string) map[string]string {
	cfg, err := ini.Load(fileName)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	res := make(map[string]string)
	keys := cfg.Section("").KeyStrings()
	for _, j := range keys{
		res[j] = cfg.Section("").Key(j).String()
	}
	return res
}
