package main

import (
	"embed"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// Embed yt-dlp and ffmpeg binaries
//
//go:embed assets/yt-dlp* assets/ffmpeg*
var embeddedFiles embed.FS

func getBinaryName(base string) string {
	if runtime.GOOS == "windows" {
		return base + ".exe"
	}
	return base
}

func extractBinary(name string) (string, error) {
	binaryPath := filepath.Join(os.TempDir(), name)
	data, err := embeddedFiles.ReadFile("assets/" + name)
	if err != nil {
		return "", fmt.Errorf("failed to read embedded %s binary: %w", name, err)
	}
	if err := ioutil.WriteFile(binaryPath, data, 0755); err != nil {
		return "", fmt.Errorf("failed to write %s binary: %w", name, err)
	}
	return binaryPath, nil
}

func fetchVideoInfo(ytDlpPath, url string) (string, string, string, error) {
	cmd := exec.Command(ytDlpPath, "--print", "%(title)s\n%(uploader)s\n%(duration)s", url)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", "", "", fmt.Errorf("error fetching video info: %w", err)
	}

	info := strings.Split(string(output), "\n")
	if len(info) < 3 {
		return "", "", "", fmt.Errorf("unexpected info format")
	}

	title := info[0]
	uploader := info[1]
	duration := info[2]

	return title, uploader, duration, nil
}

func downloadContent(ytDlpPath, ffmpegPath, url string, isAudio bool, wg *sync.WaitGroup, status *widget.Label) {
	defer wg.Done()

	format := "bestvideo[ext=mp4]+bestaudio[ext=m4a]/best[ext=mp4]"
	if isAudio {
		format = "bestaudio[ext=m4a]/bestaudio"
	}

	cmd := exec.Command(ytDlpPath, "-f", format, "--ffmpeg-location", ffmpegPath, "-o", "downloaded_content.%(ext)s", url)
	status.SetText("Downloading...")

	if err := cmd.Run(); err != nil {
		status.SetText("Download failed: " + err.Error())
		return
	}
	status.SetText("Download complete!")
}

func main() {
	a := app.New()
	w := a.NewWindow("YouTube Downloader by LewdLillyVT")
	w.Resize(fyne.NewSize(400, 400))

	ytDlpPath, err := extractBinary(getBinaryName("yt-dlp"))
	if err != nil {
		dialog.ShowError(err, w)
		return
	}
	ffmpegPath, err := extractBinary(getBinaryName("ffmpeg"))
	if err != nil {
		dialog.ShowError(err, w)
		return
	}
	defer os.Remove(ytDlpPath)
	defer os.Remove(ffmpegPath)

	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("Enter YouTube URL")

	titleLabel := widget.NewLabel("Title: ")
	uploaderLabel := widget.NewLabel("Uploader: ")
	durationLabel := widget.NewLabel("Duration: ")
	statusLabel := widget.NewLabel("")

	fetchInfoButton := widget.NewButton("Fetch Info", func() {
		url := urlEntry.Text
		if url == "" {
			dialog.ShowInformation("Error", "Please enter a valid YouTube URL", w)
			return
		}

		go func() {
			title, uploader, duration, err := fetchVideoInfo(ytDlpPath, url)
			if err != nil {
				statusLabel.SetText("Failed to fetch video info")
				fmt.Println("Error fetching video info:", err)
				return
			}

			titleLabel.SetText("Title: " + title)
			uploaderLabel.SetText("Uploader: " + uploader)
			durationLabel.SetText("Duration: " + duration)
		}()
	})

	downloadVideoButton := widget.NewButton("Download Video", func() {
		url := urlEntry.Text
		if url == "" {
			dialog.ShowInformation("Error", "Please enter a valid YouTube URL", w)
			return
		}

		var wg sync.WaitGroup
		wg.Add(1)
		go downloadContent(ytDlpPath, ffmpegPath, url, false, &wg, statusLabel)
		wg.Wait()
	})

	downloadAudioButton := widget.NewButton("Download Audio", func() {
		url := urlEntry.Text
		if url == "" {
			dialog.ShowInformation("Error", "Please enter a valid YouTube URL", w)
			return
		}

		var wg sync.WaitGroup
		wg.Add(1)
		go downloadContent(ytDlpPath, ffmpegPath, url, true, &wg, statusLabel)
		wg.Wait()
	})

	w.SetContent(container.NewVBox(
		urlEntry,
		fetchInfoButton,
		titleLabel,
		uploaderLabel,
		durationLabel,
		statusLabel,
		downloadVideoButton,
		downloadAudioButton,
	))

	w.ShowAndRun()
}
