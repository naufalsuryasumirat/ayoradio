package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
    "regexp"
	"strings"

	"github.com/joho/godotenv"
)

// run gocron to run the function for arp-scan
func main() {
	envs, err := godotenv.Read(".local.env")
	if err != nil {
		panic(err)
	}

	device := envs["AYORADIO_INTERFACE"]
	scanCmd := fmt.Sprintf("sudo arp-scan --interface=%s --localnet", device)
	cmd := exec.Command("/bin/sh", "-c", scanCmd)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}

    regexIp := regexp.MustCompile("(([0-9]+).){3}([0-9]+)")
    regexMac := regexp.MustCompile("([[:alnum:]]{2}:){5}([[:alnum:]]{2})")
    scanner := bufio.NewScanner(strings.NewReader(string(out)))
    for scanner.Scan() {
        line := scanner.Text()
        tokens := strings.Fields(line)
        if len(tokens) < 3 {
            continue
        }

        ip := tokens[0]
        isIp := regexIp.Match([]byte(ip))
        fmt.Printf("%s :: %t\n", ip, isIp)

        mac := tokens[1]
        isMac := regexMac.Match([]byte(mac))
        fmt.Printf("%s :: %t\n", mac, isMac)

        if isIp && isMac {
            fmt.Printf("format: %v\n", tokens)
        }

        provider := tokens[2:]
        fmt.Println(ip, mac, provider)
    }
}

