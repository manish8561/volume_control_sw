package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	volume "github.com/itchyny/volume-go"
)

// show volume from slider
func handleSliderChange(value float64, label *widget.Label) {
	label.SetText(fmt.Sprintf("%.2f", value))

	err := volume.SetVolume(int(value))
	if err != nil {
		fmt.Printf("set volume failed: %+v", err)
	}
}

// create desktop layout for the code
func desktopLayout(hello *widget.Label, slider *widget.Slider, label2 *widget.Label, btn1 *widget.Button) *fyne.Container {
	return container.NewGridWithRows(2,
		container.NewCenter(hello),
		container.NewVBox(
			container.NewHBox(label2, btn1),
			slider,
		),
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
		})))

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
	w.SetContent(desktopLayout(hello, slider, valueLabel, btn1))
	w.Resize(fyne.NewSize(600, 200))

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
