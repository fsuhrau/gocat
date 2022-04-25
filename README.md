# gocat
adb logcat coloring for shell

## Requirements
- Go 1.18+

## Installation
### form source
``` bash
$ go get -u github.com/fsuhrau/gocat
$ go install github.com/fsuhrau/gocat
```

### via brew ( TODO )
``` bash
$ brew install fsuhrau/homebrew-tap/gocat
```

## Usage
``` bash
$ adb logcat | gocat
$ adb logcat | gocat -tag com.yourcompany.appname
$ adb logcat | gocat -tag com.yourcompany.appname -filter includes
```
