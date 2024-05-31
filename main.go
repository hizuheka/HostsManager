package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: HostsManager [処理区分] [ホスト名のリスト]")
		os.Exit(1)
	}

	operation := os.Args[1]
	hostnames := strings.Split(os.Args[2], ",")

	hostsFile := "C:\\Windows\\System32\\drivers\\etc\\hosts"

	if err := backupHostsFile(hostsFile); err != nil {
		fmt.Println("Error backing up the hosts file:", err)
		os.Exit(1)
	}

	lines, err := readLines(hostsFile)
	if err != nil {
		fmt.Println("Error reading the hosts file:", err)
		os.Exit(1)
	}

	modifiedLines := modifyHostsLines(lines, operation, hostnames)

	if err := writeLines(hostsFile, modifiedLines); err != nil {
		fmt.Println("Error writing to the hosts file:", err)
		os.Exit(1)
	}

	fmt.Println("Hosts file modified successfully")
}

func backupHostsFile(src string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	dst := filepath.Join(filepath.Dir(src), fmt.Sprintf("hosts_%s", time.Now().Format("20060102_150405")))
	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

func readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

func modifyHostsLines(lines []string, operation string, hostnames []string) []string {
	var modifiedLines []string

	for _, line := range lines {
		hostname := getContainsHostName(line, hostnames)
		if hostname == "" {
			modifiedLines = append(modifiedLines, line)
		} else {
			if operation == "1" {
				// Comment out the line and add a new line with 127.0.0.1
				if !strings.HasPrefix(line, "#") {
					newLine := "#" + line
					fmt.Printf("置換 : %s => %s", line, newLine)
					modifiedLines = append(modifiedLines, newLine)
				}
			} else if operation == "2" {
				// Uncomment the line
				if strings.HasPrefix(line, "#") {
					newLine := strings.TrimPrefix(line, "#")
					fmt.Printf("置換 : %s => %s", line, newLine)
					modifiedLines = append(modifiedLines, newLine)
				}
				// Remove 127.0.0.1 line
				if strings.HasPrefix(line, "127.0.0.1") {
					fmt.Println("削除 : " + line)
				}
			}
		}
	}

	// アクセス禁止用の定義を追加する
	for _, h := range hostnames {
		if operation == "1" {
			// add a new line with 127.0.0.1
			newLine := fmt.Sprintf("127.0.0.1 %s", h)
			fmt.Println("追加 : " + newLine)
			modifiedLines = append(modifiedLines, newLine)
		}
	}
	return modifiedLines
}

func getContainsHostName(line string, hostnames []string) string {
	for _, hostname := range hostnames {
		if strings.Contains(line, hostname) {
			return hostname
		}
	}
	return ""
}
func writeLines(filename string, lines []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, line := range lines {
		if _, err := file.WriteString(line + "\r\n"); err != nil {
			return err
		}
	}

	return nil
}
