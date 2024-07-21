# Desktop Uitilty Volume Control App

This is desktop volume control app for your system developed in golang with fyne lib. Short cut to close window.

Link to download app.
https://apps.fyne.io/apps/com.github.manish8561.volume_control_sw.html

## Pre-requirements for OS
for macOS
```bash
brew install portaudio
```
for linux/ubuntu
```bash
sudo apt install portaudio19-dev
```

## Dependencies go packages

fyne/v2
volume
github.com/gordonklaus/portaudio

### To run in dev enviroment

desktop

```go
go run main.go
```

mobile view

```go
go run -tags mobile main.go
```

### Create package using fyne for ubuntu

fyne package -os linux -icon ./assets/volume.png

### To install after build in linux system extract file and run below command with Makefile

```bash
sudo make install
```

to uninstall

```bash
sudo make uninstall
```
