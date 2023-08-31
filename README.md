[![Say Thanks!](https://img.shields.io/badge/Say%20Thanks-!-1EAEDB.svg)](https://docs.google.com/forms/d/e/1FAIpQLSfBEe5B_zo69OBk19l3hzvBmz3cOV6ol1ufjh0ER1q3-xd2Rg/viewform)

# FLACSFX: FLAC SelF-eXtracting archive
FLACSFX builds a minimal FLAC-to-WAV transcoder with an embedded FLAC byte stream. To use, simply place a FLAC file into the repository's root directory and execute the `Build.cmd` (for non-Windows users, please see below for further instructions on how to build manually). When the resultant application is executed, it will transcode the embedded FLAC byte stream to a WAV file. This allows you to quickly and easily send losslessly compressed FLAC audio to someone who needs it as a WAV in order to reduce the file size in transit and increase transfer speeds, while also not requiring any technical know-how or additional software on the part of the recipient.

# Building manually (for non-Windows users)
**Step 1:**  
Navigate to the root directory of this repostory and open a terminal session with the root directory as the current working directory.

**Step 2:**  
Ensure the desired FLAC file that will be embedded is placed within the root directory.

**Step 3:**  
If this is the first time you are using this project, you will need to initialize the go module by issuing the following commands consecutively into the terminal.
```
go mod init main
go mod tidy
```

**Step 4:**  
Create a file named `embed.go` within the root directory following the below template, making sure to replace `file.flac` and `file.wav` with their respective values. `file.flac` should be replaced by the name of the desired FLAC file you wish to embed, `file.wav` should be replaced by the desired name of the WAV file which the application will create.
```
package main

import _ "embed"

//go:embed "file.flac"
var flacRaw []byte
var wavName string = "file.wav"
```

**Step 5:**  
Ensure the `GOARCH` and `GOOS` environmental variables are set to their appropriate values for the desired target system. Possible values for `GOARCH` include `amd64` and `386`. Possible values for `GOOS` include `windows`, `linux`, and `darwin` (Mac). 

**Step 6:**  
Build the application by issuing the following command into the terminal, making sure to replace `MyFLACSFX` with the desired name of the application file.
```
go build -ldflags="-s -w" -o "MyFLACSFX" flacsfx.go embed.go
```

# More About ScriptTiger

For more ScriptTiger scripts and goodies, check out ScriptTiger's GitHub Pages website:  
https://scripttiger.github.io/

[![Donate](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](https://www.paypal.com/cgi-bin/webscr?cmd=_s-xclick&hosted_button_id=MZ4FH4G5XHGZ4)
