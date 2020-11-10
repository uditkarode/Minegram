package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

func readConfig(fileName string) map[string]string {
	filename, _ := filepath.Abs(fileName)
	confFile, err := ioutil.ReadFile(filename)
	startSpace := regexp.MustCompile(`^\s`)
	res := make(map[string]string)

	if err != nil {
		fmt.Println("Could not open " + fileName)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(confFile)))

	for scanner.Scan() {
		split := strings.Split(scanner.Text(), "=")
		if len(split) > 1 {
			res[strings.ReplaceAll(split[0], " ", "")] = startSpace.ReplaceAllString(strings.Join(split[1:], "="), "")
		}
	}
	return res
}
