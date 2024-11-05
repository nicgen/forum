package lib

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// LoadEnv loads environment variables from a .env file
func LoadEnv(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Ignore empty lines and comments
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		// Split the line into key and value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			log.Printf("Invalid line: %s\n", line)
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Set the environment variable
		err = os.Setenv(key, value)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// AskForPassword prompts the user for a password and verifies it against DB_PASSWORD in the environment
func AskForPassword() bool {
	var password string
	fmt.Print("Enter Password: ")
	fmt.Scanln(&password)

	storedPassword := os.Getenv("DB_PASSWORD")
	return password == storedPassword
}
