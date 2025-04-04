package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"time"
	"os/exec"
)

var stopChan chan bool

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Countdown Timer")

	label := widget.NewLabel("00:00")
	progressBar := widget.NewProgressBar()
	progressBar.Max = 1.0

	entry := widget.NewEntry()
	entry.SetPlaceHolder("Enter time in seconds")

	startButton := widget.NewButton("Start", func() {
		duration, err := time.ParseDuration(entry.Text + "s")
		if err != nil {
			label.SetText("Invalid time")
			return
		}
		stopChan = make(chan bool)
		go startTimer(label, progressBar, duration)
	})

	stopButton := widget.NewButton("Stop", func() {
		if stopChan != nil {
			stopChan <- true
		}
	})

	myWindow.SetContent(container.NewVBox(
		entry,
		label,
		progressBar,
		startButton,
		stopButton,
	))

	myWindow.Resize(fyne.NewSize(200, 150))
	myWindow.ShowAndRun()
}

func startTimer(label *widget.Label, progressBar *widget.ProgressBar, duration time.Duration) {
	startTime := time.Now()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-stopChan:
			label.SetText("Stopped")
			progressBar.SetValue(0)
			return
		case <-ticker.C:
			elapsed := time.Since(startTime)
			remaining := duration - elapsed
			if remaining <= 0 {
				label.SetText("00:00")
				progressBar.SetValue(0)
				playSound()
				return
			}
			minutes := int(remaining.Minutes())
			seconds := int(remaining.Seconds()) % 60
			label.SetText(fmt.Sprintf("%02d:%02d", minutes, seconds))
			progressBar.SetValue(float64(elapsed) / float64(duration))
		}
	}
}

func playSound() {
	cmd := exec.Command("say", "jikan death")
	cmd.Run()
}
