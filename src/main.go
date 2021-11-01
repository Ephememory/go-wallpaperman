package main

import (
	"fmt"
	"os"
	"strings"
	"time"
	"unsafe"

	"github.com/getlantern/systray"

	"golang.org/x/sys/windows"

	"os/exec"
)

// type TimeFrame struct {
// 	Start, End int
// }

//TODO: make lat and long user config or automatic?
const (
// ss_api_url string = "https://api.sunrise-sunset.org/json?lat=31.429340&lng=-87.344050"
// sunrise TimeFrame = new(TimeFrame)
)

var (
	user32DLL           *windows.LazyDLL
	procSystemParamInfo *windows.LazyProc
)

func main() {
	systray.Run(onReady, onExit)

}

func onReady() {
	systray.SetIcon(IconData)

	systray.SetTitle("wallpaperman")
	systray.SetTooltip("Wallpaper Manager")

	openFolder := systray.AddMenuItem("Open containing folder", "Opens the folder the app is running in.")

	go func() {
		<-openFolder.ClickedCh
		cmd := exec.Command(`explorer`, `/select,`, `/`)
		cmd.Run()
	}()

	systray.AddSeparator()

	close := systray.AddMenuItem("Close", "Closes the app")
	go func() {
		<-close.ClickedCh
		systray.Quit()
	}()

	postReady()
}

func onExit() {

}

func mainLoop(day, night *uint16) {
	for {
		currentTime := time.Now().Local()
		timeOfDayImage := day

		if currentTime.Hour() > 6 && currentTime.Hour() < 18 {
			timeOfDayImage = day
		} else {
			timeOfDayImage = night
		}

		procSystemParamInfo.Call(20, 0, uintptr(unsafe.Pointer(timeOfDayImage)), 0x001A)
		time.Sleep(time.Hour * 1)
	}

}

func postReady() {
	user32DLL = windows.NewLazyDLL("user32.dll")
	procSystemParamInfo = user32DLL.NewProc("SystemParametersInfoW")

	directory, directoryError := os.Executable()

	if directoryError != nil {
		os.Exit(1)
		return
	}

	filteredPath := strings.Split(directory, "\\wallpaperman.exe")[0]
	dayImagePath := fmt.Sprintf("%v\\day.png", filteredPath)
	nightImagePath := fmt.Sprintf("%v\\night.png", filteredPath)

	dayTimeImage, dayError := windows.UTF16PtrFromString(dayImagePath)
	nightTimeImage, nightError := windows.UTF16PtrFromString(nightImagePath)

	if dayError != nil || nightError != nil {
		fmt.Println("Could not find bg images!")
		os.Exit(1)
		return
	}

	go mainLoop(dayTimeImage, nightTimeImage)
}
