package jobs

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"

	util "github.com/naufalsuryasumirat/ayoradio/util"
)

var regexIp = regexp.MustCompile("(([0-9]+).){3}([0-9]+)")
var regexMac = regexp.MustCompile("([[:alnum:]]{2}:){5}([[:alnum:]]{2})")

var blacklisted = make(map[string]bool)

var mu sync.RWMutex

func LoadBlacklistedDevices() {
	devices := util.GetBlacklistedDevices()

	mu.Lock()
    defer mu.Unlock()
	for _, device := range devices {
		blacklisted[device] = true
	}
}

func ScanLocalDevices() []string {
	device := os.Getenv("AYORADIO_INTERFACE")
	scanCmd := fmt.Sprintf("sudo arp-scan --interface=%s --localnet", device)
	cmd := exec.Command("/bin/sh", "-c", scanCmd)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

    _ = cmd.Run()

	out, err := cmd.Output()
	if err != nil {
		log.Println(err)
	}

	var locals []string

    mu.RLock()
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Fields(line)
		if len(tokens) < 3 {
			continue
		}

		ip := tokens[0]
		isIp := regexIp.Match([]byte(ip))

		mac := tokens[1]
		isMac := regexMac.Match([]byte(mac))

		if isIp && isMac && !blacklisted[mac] {
			locals = append(locals, mac)
		}

		_ = tokens[2:] // provider
	}
    mu.RUnlock()

	return util.ExistDevices(locals)
}

func init() {
	LoadBlacklistedDevices()
	for device := range blacklisted {
		fmt.Printf("blacklisted: %s\n", device)
	}
}
