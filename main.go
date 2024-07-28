package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	volume "github.com/itchyny/volume-go"
)

// get default sink and source
func getDefaultDevice(deviceType string) (string, error) {
	var cmd *exec.Cmd
	if deviceType == "sinks" {
		cmd = exec.Command("pactl", "get-default-sink")
	} else if deviceType == "sources" {
		cmd = exec.Command("pactl", "get-default-source")
	} else {
		return "", fmt.Errorf("invalid device type: %s", deviceType)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to switch %s: %v, output: %s", deviceType, err, output)
	}
	return string(output), nil
}

// switchPulseAudioDevice switches the active PulseAudio sink or source
func switchPulseAudioDevice(deviceType, deviceName string) error {
	var cmd *exec.Cmd
	if deviceType == "sinks" {
		cmd = exec.Command("pactl", "set-default-sink", deviceName)
	} else if deviceType == "sources" {
		cmd = exec.Command("pactl", "set-default-source", deviceName)
	} else {
		return fmt.Errorf("invalid device type: %s", deviceType)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to switch %s: %v, output: %s", deviceType, err, output)
	}
	return nil
}

// getPulseAudioDevices lists all PulseAudio sinks or sources
func getPulseAudioDevices(deviceType string) ([]string, error) {
	cmd := exec.Command("pactl", "list", "short", deviceType)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to list %s: %v", deviceType, err)
	}

	var devices []string
	scanner := bufio.NewScanner(&out)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) > 1 {
			devices = append(devices, fields[1])
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %v", deviceType, err)
	}

	return devices, nil
}

// show volume from slider
func handleSliderChange(value float64, label *widget.Label) {
	label.SetText(fmt.Sprintf("%.2f", value))

	err := volume.SetVolume(int(value))
	if err != nil {
		fmt.Printf("set volume failed: %+v", err)
	}
}

// create desktop layout for the code
func desktopLayout(hello *widget.Label, slider *widget.Slider, label2 *widget.Label, btn1 *widget.Button, comboOut *widget.Select, comboIn *widget.Select) *fyne.Container {
	tabs := container.NewAppTabs(
		container.NewTabItem("Output Devices", container.NewVBox(comboOut)),
		container.NewTabItem("Input Devices", container.NewVBox(comboIn)),
	)
	tabs.SetTabLocation(container.TabLocationTop)

	return container.NewGridWithRows(2,
		container.NewVBox(
			container.NewCenter(hello),
			container.NewHBox(label2, btn1),
			slider,
		),
		container.New(layout.NewVBoxLayout(), tabs),
	)
}

func checkMuteUnmute(btn1 *widget.Button, check bool) {
	//check if volume is already muted or not
	muted, err := volume.GetMuted()
	if err != nil {
		fmt.Printf("get mute value failed: %+v", err)
	}
	if muted {
		btn1.SetText("Unmute")
		if check {
			err = volume.Unmute()
			if err != nil {
				fmt.Printf("unmute failed: %+v", err)
			}
			btn1.SetText("Mute")
		}

	} else {
		btn1.SetText("Mute")
		if check {
			err = volume.Mute()
			if err != nil {
				fmt.Printf("mute failed: %+v", err)
			}
			btn1.SetText("Unmute")
		}
	}
}

func main() {
	// on load volume current value
	vol, err := volume.GetVolume()
	if err != nil {
		fmt.Printf("Get volume failed: %+v", err)
	}

	a := app.New()
	w := a.NewWindow("System Volume Control")
	w.SetMaster()
	//setting volume icon
	w.SetIcon(theme.VolumeUpIcon())

	hello := widget.NewLabel("System Volume")
	// rendering menu
	mainMenu := fyne.NewMainMenu(fyne.NewMenu("File",
		fyne.NewMenuItem("Show", func() {
			hello.SetText("Waiting for this feature!")
		}),
	))

	slider := widget.NewSlider(0, 100)
	slider.SetValue(float64(vol))

	// Create a label to display the current slider value
	valueLabel := widget.NewLabel(fmt.Sprintf("%.2f", slider.Value))

	slider.OnChanged = func(value float64) {
		// fmt.Printf("Slider value changed: %.2f\n", value)
		// Call your custom function here to handle the value change
		handleSliderChange(value, valueLabel)
	}
	var btn1 *widget.Button

	// function to mute and unmute the volume
	btn1 = widget.NewButton("Mute", func() {
		checkMuteUnmute(btn1, true)
	})
	//check if volume is already muted or not on window load
	checkMuteUnmute(btn1, false)

	w.SetMainMenu(mainMenu)

	// Get available sinks (output devices)
	output, err := getPulseAudioDevices("sinks")
	if err != nil {
		fmt.Printf("Error getting sinks: %v", err)
	}
	// Get available sources (input devices)
	input, err := getPulseAudioDevices("sources")
	if err != nil {
		fmt.Printf("Error getting sources: %v", err)
	}

	//adding drop down for input devices
	comboIn := widget.NewSelect(input, func(value string) {
		// fmt.Println("Select set to:-> ", value)
		switchPulseAudioDevice("sources", value)
	})
	//default input device
	defaultDevice, err := getDefaultDevice("sources")
	if err != nil {
		fmt.Printf("Error getting sources: %v", err)
	}

	comboIn.SetSelected(defaultDevice)

	//adding drop down for output devices
	comboOut := widget.NewSelect(output, func(value string) {
		// fmt.Println("Select set to:-> ", value)
		switchPulseAudioDevice("sinks", value)

	})
	defaultDevice, err = getDefaultDevice("sinks")
	if err != nil {
		fmt.Printf("Error getting sinks: %v", err)
	}
	comboOut.SetSelected(defaultDevice)

	content := desktopLayout(hello, slider, valueLabel, btn1, comboOut, comboIn)

	// Create a background rectangle with a specified color

	// Create a padding container using layout.NewPadded
	padding := container.New(layout.NewPaddedLayout(), content)

	w.SetContent(padding)
	w.Resize(fyne.NewSize(600, 300))

	// Define custom keyboard shortcut for the button (e.g., Ctrl+KeyUp)
	ctrlKeyUp := &desktop.CustomShortcut{KeyName: fyne.KeyUp, Modifier: fyne.KeyModifierControl}
	w.Canvas().AddShortcut(ctrlKeyUp, func(shortcut fyne.Shortcut) {
		vol, err := volume.GetVolume()
		if err != nil {
			fmt.Printf("get volume failed: %+v", err)
		}
		if vol >= 100 {
			vol = 100
		}
		vol += 1
		slider.SetValue(float64(vol))
	})
	// Define custom keyboard shortcut for the button (e.g., Ctrl+Down)
	ctrlKeyDown := &desktop.CustomShortcut{KeyName: fyne.KeyDown, Modifier: fyne.KeyModifierControl}
	w.Canvas().AddShortcut(ctrlKeyDown, func(shortcut fyne.Shortcut) {
		vol, err := volume.GetVolume()
		if err != nil {
			fmt.Printf("get volume failed: %+v", err)
		}
		if vol <= 0 {
			vol = 0
		}
		vol -= 1
		slider.SetValue(float64(vol))
	})
	// Define custom keyboard shortcut for the button (e.g., Ctrl+w) to close window
	ctrlKeyW := &desktop.CustomShortcut{KeyName: fyne.KeyW, Modifier: fyne.KeyModifierControl}
	w.Canvas().AddShortcut(ctrlKeyW, func(shortcut fyne.Shortcut) {
		a.Quit()
	})
	// show and run window
	w.ShowAndRun()
	// w.Show()
	// a.Run()
}
