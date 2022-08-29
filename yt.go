package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

var url, format, output string

func die(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func FileExists(fn string) bool {
	stat, err := os.Stat(fn)
	if os.IsNotExist(err) {
		return false
	}

	return !stat.IsDir()
}

// grab the mp3 and mp4 url from youtube-dl
func UrlDownload(u string) (string, string) {
	if FileExists("url.txt") {
		err := os.Remove("url.txt")
		die(err)
	}
	file, err := os.Create("url.txt")
	die(err)

	cmd := exec.Command("youtube-dl", "--get-url", url)
	cmd.Stdout = file
	err = cmd.Run()
	die(err)

	file.Close()

	file, err = os.Open("url.txt")
	die(err)

	s := bufio.NewScanner(file)
	s.Split(bufio.ScanLines)

	var lines []string

	for s.Scan() {
		lines = append(lines, s.Text())
	}

	mp4_url := lines[0]
	mp3_url := lines[1]

	err = os.Remove("url.txt")
	die(err)

	return mp4_url, mp3_url
}

// Grab the title from youtube-dl
func VidTitle(u string) string {
	out, err := exec.Command("youtube-dl", "e", u).Output()
	die(err)

	title := string(out)
	title = strings.Replace(title, "\n", "", -1)

	return title
}

func VidDownload(wg *sync.WaitGroup, u, fn string) {
	defer wg.Done()
	fmt.Println("[+]Downloading Video")
	if FileExists(fn) {
		err := os.Remove(fn)
		die(err)
	}

	f, err := os.Create(fn)
	die(err)

	c := http.Client{Timeout: 60 * time.Second}

	resp, err := c.Get(u)
	if err != nil {
		f.Close()
		log.Fatal(err)
	}

	defer resp.Body.Close()

	if _, err = io.Copy(f, resp.Body); err != nil {
		f.Close()
		log.Fatal(err)
	}

	f.Close()
	fmt.Println("[+]Download Complete")
}

func AudioDownload(wg *sync.WaitGroup, u, fn string) {
	defer wg.Done()
	fmt.Println("[+]Downloading Audio")

	if FileExists(fn) {
		err := os.Remove(fn)
		die(err)
	}

	f, err := os.Create(fn)
	die(err)

	resp, err := http.Get(u)
	if err != nil {
		f.Close()
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if _, err = io.Copy(f, resp.Body); err != nil {
		f.Close()
		log.Fatal(err)
	}

	f.Close()
	fmt.Println("[+]Download Complete")
}

// Merge the mp3 audio with the mp4 video
func Merge(out_fn, mp4_path, mp3_path string) {
	if out_fn == "" {
		if FileExists(out_fn) {
			err := os.Remove(out_fn)
			die(err)
		}
	}

	cmd := exec.Command("ffmpeg", "-i", mp4_path, "-i", mp3_path,
		"-map",
		"0:v", "1:a",
		"-c:v",
		"copy",
		"-c:a",
		"copy",
		"-y", out_fn)
	err := cmd.Run()
	die(err)

	err = os.Remove(mp4_path)
	die(err)

	err = os.Remove(mp3_path)
	die(err)
}

func WebmConverter(in, out string) {
	cmd := exec.Command("ffmpeg", "-i",
		in,
		"-vn",
		"-ab",
		"128k",
		"-ar",
		"44100",
		"-y",
		out)
	err := cmd.Run()
	die(err)

	err = os.Remove(in)
	die(err)
}

// flags
func init() {
	flag.StringVar(&url, "u", "", "input url")
	flag.StringVar(&format, "f", "mp3", "mp4 or mp3")
	flag.StringVar(&output, "o", "", "output filename")
}

func main() {
	flag.Parse()

	if url == "" {
		fmt.Println("[-]You forgot the url")
		os.Exit(1)
	}

	switch format {
	case "mp4":
		// get it straight from youtube-dl
		mp4_url, mp3_url := UrlDownload(url)
		tmp_webm_audio := "tmp_audio.webm"
		tmp_video := "tmp_video.mp4"

		var wg sync.WaitGroup
		wg.Add(2)

		// grab both the audio and video
		go VidDownload(&wg, mp4_url, tmp_video)
		go AudioDownload(&wg, mp3_url, tmp_webm_audio)
		wg.Wait()

		// convert webm to mp3
		fmt.Println("[+]Converting file...")
		tmp_audio_mp3 := "tmp_audio.mp3"
		WebmConverter(tmp_webm_audio, tmp_audio_mp3)

		if output == "" {
			output = VidTitle(url)
			output = output + ".mp4"
		}

		Merge(output, tmp_video, tmp_audio_mp3)

		fmt.Println("[+]Done.")

	}
	fmt.Println("vim-go")
}
