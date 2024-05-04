# Desktop Uitilty Volume Control App

this is desktop volume control app for your system developed in golang with fyne lib. Short cut to close window

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
