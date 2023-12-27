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
	release := ""
	releaseRaw, err := os.ReadFile("/srv/release")
	if err != nil {
		log.Println("Downloading Raspberry Pi OS release notes")
		resp, err := http.Get("https://downloads.raspberrypi.com/raspios_lite_arm64/release_notes.txt")
		defer resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
		buf := bufio.NewReader(resp.Body)
		releaseRaw, err := buf.ReadBytes(':')
		if err != nil {
			log.Fatal(err)
		}
		releaseRaw = releaseRaw[:len(releaseRaw)-1]
		release = string(releaseRaw)
		err = os.WriteFile("/srv/release", releaseRaw, 0o644)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		release = string(releaseRaw)
	}
	log.Printf("Verison: %s", release)
	return string(releaseRaw)
}

func main() {
	fmt.Println(GetLatestRelease())
}
