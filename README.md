[![Say Thanks!](https://img.shields.io/badge/Say%20Thanks-!-1EAEDB.svg)](https://docs.google.com/forms/d/e/1FAIpQLSfBEe5B_zo69OBk19l3hzvBmz3cOV6ol1ufjh0ER1q3-xd2Rg/viewform)

# FLACSFX: FLAC SelF-eXtracting archive
FLACSFX is a minimal FLAC-to-WAV transcoder to transcode an embedded FLAC file to a WAV file. The FLAC file can either be embedded at build time using the `embed.go` or embedded later by appending a FLAC file to a stand-alone FLACSFX executable built with the `sa.go`. This allows you to quickly and easily send losslessly compressed FLAC audio to someone who needs it as a WAV in order to reduce the file size in transit and increase transfer speeds, while also not requiring any technical know-how or additional software on the part of the recipient.

Usage: `flacsfx [options...]`
Argument                  | Description
--------------------------|-----------------------------------------------------------------------------------------------------
 `-o <file>`              | Destination file

`-` can be used in place of `<file>` to designate standard output as the destination, allowing you to quickly pipe the WAV data to a compatible application, such as VLC, without having to extract it to an actual file.

# Appending a FLAC file to a stand-alone FLACSFX executable, recommended for most users
Download the latest pre-built release for the intended target system:  
https://github.com/ScriptTiger/FLACSFX/releases/latest

For appending a FLAC file to a FLACSFX executable, issue one of the following commands.

For Windows:
```
copy /b "FLACSFX.exe"+"file.flac" "MyFLACSFX.exe"
```

For Linux and Mac:
```
cat "FLACSFX" "file.flac" > "MyFLACSFX"
```

# Building a FLACSFX application in Go, for Windows users
If you would like to embed a FLAC file at build time, you can use the `Build-embed.cmd`. To use, simply place a FLAC file into the repository's root directory and execute the `Build-embed.cmd`. When the resultant application is executed, it will transcode the embedded FLAC file to a WAV file.

If you would like to build a stand-alone FLACSFX executable and append a FLAC file to it later, execute the `Build-SA.cmd`.

# Building a FLACSFX application in Go, for non-Windows users
**Step 1:**  
Navigate to the root directory of this repostory and open a terminal session with the root directory as the current working directory. Ensure the `GOARCH` and `GOOS` environmental variables are set to their appropriate values for the desired target system. Possible values for `GOARCH` include `amd64` and `386`. Possible values for `GOOS` include `windows`, `linux`, and `darwin` (Mac).

**Step 2:**  
If this is the first time you are using this project, you will need to initialize the go module by issuing the following commands consecutively into the terminal.
```
go mod init main
go mod tidy
```

**Step 3 (Skip if not embedding at build time):**  
Ensure the desired FLAC file that will be embedded is placed within the root directory. Create a file named `embed.go` within the root directory following the below template, making sure to replace `file.flac` and `file.wav` with their respective values. `file.flac` should be replaced by the name of the desired FLAC file you wish to embed, `file.wav` should be replaced by the desired name of the WAV file which the application will create.
```
package main

import _ "embed"

//go:embed "file.flac"
var flacRaw []byte
var wavName string = "file.wav"
```

**Step 4:**  
Build the application by issuing one of the following commands into the terminal, making sure to replace `MyFLACSFX` with the desired name of the application file.

For embedding a FLAC file at build time:
```
go build -ldflags="-s -w" -o "MyFLACSFX" flacsfx.go embed.go
```

For creating a stand-alone FLACSFX executable which you can append a FLAC file to later:
```
go build -ldflags="-s -w" -o "MyFLACSFX" flacsfx.go sa.go
```

# More About ScriptTiger

For more ScriptTiger scripts and goodies, check out ScriptTiger's GitHub Pages website:  
https://scripttiger.github.io/

[![Donate](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](https://www.paypal.com/cgi-bin/webscr?cmd=_s-xclick&hosted_button_id=MZ4FH4G5XHGZ4)
