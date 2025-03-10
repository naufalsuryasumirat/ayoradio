package util

import (
	"database/sql"
	"fmt"
	"log"
    "regexp"
	"strings"
)

var regexMac = regexp.MustCompile("([[:alnum:]]{2}:){5}([[:alnum:]]{2})")

func TryAddDevice(addr string, blacklisted bool) error {
    if !regexMac.Match([]byte(addr)) {
        return fmt.Errorf("not a mac address")
    }

	res, err := db.Exec(
        fmt.Sprintf(
            "INSERT OR IGNORE INTO devices (mac_address, blacklisted) VALUES (?, %d);",
            func() int {
                if blacklisted {
                    return 1
                } else {
                    return 0
                }
            }(),
        ),
        addr,
    )
    if err != nil {
        return err
    }

    aff, err := res.RowsAffected()
    if err != nil {
        return err
    }

    if aff == 0 {
        return fmt.Errorf("device already exist")
    }

    return nil
}

func AddDevice(addr string) {
	_, err := db.Exec("INSERT OR IGNORE INTO devices (mac_address, blacklisted) VALUES (?, 0);", addr)
	if err != nil {
		log.Println(err)
	}
}

func AddBlacklistedDevice(addr string) {
    _, err := db.Exec("INSERT OR IGNORE INTO devices (mac_address, blacklisted) VALUES (?, 1);", addr)
    if err != nil {
        log.Println(err)
    }
}

func ExistDevice(addr string) bool {
    var exist int
    err := db.QueryRow("SELECT COUNT(1) FROM devices WHERE mac_address = ?", addr).Scan(&exist)
	if err != nil {
        if err != sql.ErrNoRows {
            log.Println(err)
        }
        return false
	}

    return exist == 1
}

func ExistDevices(addrs []string) []string {
    var res []string
    if addrs == nil || len(addrs) == 0 {
        return res
    }

    query := fmt.Sprintf(
        "SELECT mac_address FROM devices WHERE mac_address IN (%s);",
        strings.Repeat("?, ", len(addrs)-1) + "?",
    )
    args := make([]interface{}, len(addrs))
    for i, v := range addrs {
        args[i] = v
    }

    row, _ := db.Query(query, args...)
    defer row.Close()

    for row.Next() {
        var mac string
        row.Scan(&mac)
        res = append(res, mac)
    }

    return res
}

func GetBlacklistedDevices() []string {
    row, _ := db.Query("SELECT mac_address FROM devices WHERE blacklisted = 1;")
    defer row.Close()

    var res []string
    for row.Next() {
        var mac string
        row.Scan(&mac)
        res = append(res, mac)
    }

    return res
}
