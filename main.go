package main

import (
	"fmt"
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	volume "github.com/itchyny/volume-go"

	"github.com/gordonklaus/portaudio"
)

/**
 * get input and output devices portaudio
 */
func getInputOutputDevices() (in, out []string) {
	input, output := []string{}, []string{}

	// Initialize PortAudio
	if err := portaudio.Initialize(); err != nil {
		fmt.Println("Initialize 1: ", err)
	}
	defer portaudio.Terminate()

	// Get a list of all devices
	devices, err := portaudio.Devices()
	if err != nil {
		fmt.Println("devices error1: ", err)
	}
	// List all devices
	for i, device := range devices {
		fmt.Printf("%d: %s\n", i, device.Name)
		fmt.Println("--------------------------------------------")
		if device.MaxInputChannels > 0 {
			fmt.Println("  Input channels:", device.Name)
			input = append(input, device.Name)
		}
		if device.MaxOutputChannels > 0 {
			fmt.Println("  Output channels:", device.Name)
			output = append(output, device.Name)

		}
		fmt.Printf("  Default sample rate: %f\n\n", device.DefaultSampleRate)

	}
	return input, output
}

/**
 * get input and output devices portaudio
 */
func setInputOutputDevices(deviceName string, deviceType int) {

	// Initialize PortAudio
	if err := portaudio.Initialize(); err != nil {
		fmt.Println("Initialize 2: ", err)
	}
	defer portaudio.Terminate()

	// Get a list of all devices
	devices, err := portaudio.Devices()
	if err != nil {
		fmt.Println("devices error2: ", err)
	}

	// List all devices
	for i, device := range devices {
		if strings.ToLower(device.Name) == strings.ToLower(deviceName) && deviceType == 1 && device.MaxInputChannels > 0 {
			fmt.Printf("Input channels: %d\n", device.MaxInputChannels, i)
			fmt.Printf("%+v", device)
			fmt.Println()
			fmt.Printf("%+v", device.HostApi)
			// Set up a stream with the selected devices
			stream, err := portaudio.OpenStream(portaudio.StreamParameters{
				Input: portaudio.StreamDeviceParameters{
					Device:   device,
					Channels: 1,
					Latency:  device.DefaultLowInputLatency,
				},
				SampleRate:      device.DefaultSampleRate,
				FramesPerBuffer: 64,
			}, processAudio)

			if err != nil {
				fmt.Println("Input device setting: ", err)
			}
			if err != nil {
				log.Fatalf("Error opening stream: %v", err)
			}
			defer stream.Close()
			return

		}
		if strings.ToLower(device.Name) == strings.ToLower(deviceName) && deviceType == 2 && device.MaxOutputChannels > 0 {

			fmt.Printf("Output channels: %d\n", i)
			fmt.Printf("%+v", device)
			fmt.Println()
			fmt.Printf("%+v", device.HostApi)

			// Set up a stream with the selected devices
			stream, err := portaudio.OpenStream(portaudio.StreamParameters{

				Output: portaudio.StreamDeviceParameters{
					Device:   device,
					Channels: 1,
					Latency:  device.DefaultLowOutputLatency,
				},
				SampleRate:      device.DefaultSampleRate,
				FramesPerBuffer: 64,
			}, processAudio)

			if err != nil {
				fmt.Println("Output device setting: ", err)
			}
			if err != nil {
				log.Fatalf("Error opening stream: %v", err)
			}
			defer stream.Close()
			return
		}

	}

}

func processAudio(in, out []float32) {
	// This is where you'd process the audio data
	for i := range out {
		out[i] = in[i]
	}
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
	// on load volume
	vol, err := volume.GetVolume()
	if err != nil {
		fmt.Printf("get volume failed: %+v", err)
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

	// input output devices
	input, output := getInputOutputDevices()

	//adding drop down for input devices
	comboIn := widget.NewSelect(input, func(value string) {
		fmt.Println("Select set to", value)
		setInputOutputDevices(value, 1)
	})
	comboIn.SetSelected("default")

	//adding drop down for output devices
	comboOut := widget.NewSelect(output, func(value string) {
		fmt.Println("Select set to", value)
		setInputOutputDevices(value, 2)
	})
	comboOut.SetSelected("default")

	w.SetContent(desktopLayout(hello, slider, valueLabel, btn1, comboOut, comboIn))
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
