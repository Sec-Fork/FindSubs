package Common

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// readDomainFromFile 从文件中读取域名列表
func ReadDomainFromFile(filepath string) ([]string, error) {
	var domains []string

	file, err := os.Open(filepath)
	if err != nil {
		return domains, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		domain := strings.TrimSpace(scanner.Text())
		if domain != "" {
			domains = append(domains, domain)
		}
	}

	if err := scanner.Err(); err != nil {
		return domains, err
	}

	return domains, nil
}

func GetApiKey() (ApiKey string, ApiID string, ApiSecret string) {
	file, err := os.Open("api-key.ini")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()
	config := make(map[string]string)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "#") && strings.Contains(line, "=") {
			pair := strings.SplitN(line, "=", 2)
			config[pair[0]] = strings.TrimSpace(pair[1])
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error:", err)
		return
	}
	return config["ApiKey"], config["ApiID"], config["ApiSecret"]
}
