package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// LoadEnv takes the path of a .env file and uses it to set environment variables
func LoadEnv(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Unable to open env file at", path, ":", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		data := strings.Split(line, "=")
		key := strings.TrimSpace(data[0])
		value := strings.TrimSpace(data[1])

		os.Setenv(key, value)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Error reading env file at", path, ":", err)
	}
}
