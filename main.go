package main

import (
	"bufio"
	"fmt"
	"github.com/ulikunitz/xz"
	"io"
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
	return release
}

func DownloadRelease(release string) {
	url := fmt.Sprintf("https://downloads.raspberrypi.com/raspios_lite_arm64/images/raspios_lite_arm64-%s/%s-raspios-bookworm-arm64-lite.info", release, release)
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	r, err := xz.NewReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fi, err := os.OpenFile("/srv/raspios.img", 0, 0o644)
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(fi, r)
	if err != nil {
		log.Fatal(err)
	}
}

func EnsureDownload(release string) {
	if _, err := os.Stat("/srv/raspios.img"); err != nil {
		DownloadRelease(release)
	}
}

func main() {
	release := GetLatestRelease()
	EnsureDownload(release)
}
