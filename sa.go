package main

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"strings"
)

var wavName string
var flacRaw []byte

func init() {

	// Locate executable
	filePath, _ := os.Executable()
	filePath, _ = filepath.EvalSymlinks(filePath)

	// Slurp file into memory
	fileData, _ := os.ReadFile(filePath)

	// Determine length of transcoder / start of FLAC data
	// Find TITLE tag and use for file name if present
	tcSize := int64(1000000)
	var titleTagBuilder strings.Builder
	for {
		if fileData[tcSize-1] == 0x00 &&
		fileData[tcSize] == 0x66 &&
		fileData[tcSize+1] == 0x4c &&
		fileData[tcSize+2] == 0x61 &&
		fileData[tcSize+3] == 0x43 &&
		fileData[tcSize+4] == 0x00 {
			titleTag := tcSize+32
			for {
				if fileData[titleTag] == 0x00 &&
				fileData[titleTag+1] == 0x54 &&
				fileData[titleTag+2] == 0x49 &&
				fileData[titleTag+3] == 0x54 &&
				fileData[titleTag+4] == 0x4C &&
				fileData[titleTag+5] == 0x45 &&
				fileData[titleTag+6] == 0x3D {
					titleTag = titleTag+7
					for {
						if fileData[titleTag+1] == 0x00 {break}
						titleTagBuilder.WriteByte(fileData[titleTag])
						titleTag++
					}
					break
				}
				if (int(titleTag) == len(fileData)-6) ||
				(fileData[titleTag] == 0x00 &&
				fileData[titleTag+1] == 0x00 &&
				fileData[titleTag+2] == 0x00 &&
				fileData[titleTag+3] == 0x00 &&
				fileData[titleTag+4] == 0x00 &&
				fileData[titleTag+5] == 0x00 &&
				fileData[titleTag+6] == 0x00 &&
				fileData[titleTag+7] == 0x00 &&
				fileData[titleTag+8] == 0x00 &&
				fileData[titleTag+9] == 0x00 &&
				fileData[titleTag+10] == 0x00 &&
				fileData[titleTag+11] == 0x00 &&
				fileData[titleTag+12] == 0x00 &&
				fileData[titleTag+13] == 0x00 &&
				fileData[titleTag+14] == 0x00 &&
				fileData[titleTag+15] == 0x00) {
					break
				}
				titleTag++
			}
		break
		}
		if int(tcSize) == len(fileData)-6 {
			os.Stdout.WriteString("No FLAC data found.")
			os.Exit(5)
		}
		tcSize++
	}

	// Name wav file after executable, or title tag if present
	if titleTagBuilder.Len() == 0 {
		wavName = strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))+".wav"
	} else {
		wavName = titleTagBuilder.String()+".wav"
	}

	// Calculate length of FLAC data
	flacSize := int64(len(fileData))-tcSize

	// Create file data reader
	reader := bytes.NewReader(fileData)

	// Shift position of reader to the start of the flac data
	_, _ = reader.Seek(tcSize, os.SEEK_SET)

	// Create byte array to hold flac data
	flacRaw = make([]byte, flacSize)

	// Load empty byte arrary with flac data
	binary.Read(reader, binary.LittleEndian, &flacRaw)
}