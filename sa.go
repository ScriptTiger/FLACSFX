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

	// Name wav file after executable
	wavName = strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))+".wav"

	// Slurp file into memory
	fileData, _ := os.ReadFile(filePath)

	// Determine length of transcoder / start of flac data
	tcSize := int64(1000000)
	for {
		if fileData[tcSize-1] == 0x00 &&
			fileData[tcSize] == 0x66 &&
			fileData[tcSize+1] == 0x4c &&
			fileData[tcSize+2] == 0x61 &&
			fileData[tcSize+3] == 0x43 &&
			fileData[tcSize+4] == 0x00 {
				break
		}
		tcSize++
	}

	// Calculate length of flac data
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