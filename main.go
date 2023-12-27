package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
)

// GetLatestRelease TODO: Make this function use os_list_v3.json
func GetLatestRelease() string {
	resp, err := http.Get("https://downloads.raspberrypi.com/raspios_lite_arm64/release_notes.txt")
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	buf := bufio.NewReader(resp.Body)
	verRaw, err := buf.ReadBytes(':')
	if err != nil {
		log.Fatal(err)
	}
	return string(verRaw[:len(verRaw)-1])
}

func main() {
	releaseRaw, err := os.ReadFile("/srv/release")
	if err != nil {
		log.Println("Downloading Raspberry Pi OS release notes")
		releaseRaw = []byte(GetLatestRelease())
		err := os.WriteFile("/srv/release", releaseRaw, 0o644)
		if err != nil {
			log.Fatal(err)
		}
	}
	release := string(releaseRaw)
	fmt.Printf("Release: %s\n", release)
}
