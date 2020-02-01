package util

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"strings"
)

var logErr = log.Errorf

func ReadRandomLineFromFile(envName string, baseEndpoint string, prefix string) (string, error) {
	path := os.Getenv(envName)
	if path == "" {
		return "", fmt.Errorf("%s was not set", envName)
	}
	// read in file to list of strings
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return "", err
	}
	scanner := bufio.NewScanner(file)
	words := []string{}
	for scanner.Scan() {
		words = append(words, strings.ToLower(scanner.Text()))
	}
	err = scanner.Err()
	// get random index of list
	return baseEndpoint + prefix + words[rand.Intn(len(words))], err
}
