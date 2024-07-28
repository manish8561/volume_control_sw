# Desktop Uitilty Volume Control App

This is desktop volume control app for your system developed in golang with fyne lib. Short cut to close window.

Link to download app.
https://apps.fyne.io/apps/com.github.manish8561.volume_control_sw.html

## Operating System Dependency

For Linux
```bash
sudo apt install pulseaudio-utils
```
For MacOS, Windows coming soon.

## Dependencies go packages

fyne/v2
volume

### To run in dev enviroment

desktop

```go
go run main.go
```

mobile view

```go
go run -tags mobile main.go
```

### Create build using fyne for ubuntu

fyne package -os linux -icon ./assets/volume.png

### To install after build in linux system extract file and run below command with Makefile

```bash
sudo make install
```

to uninstall

```bash
sudo make uninstall
```
