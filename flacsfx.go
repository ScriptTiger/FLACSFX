package main

import (
	"bytes"
	"io"
	"os"
	"strings"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/mewkiz/flac"
)

// Function to display help text and exit
func help(err int) {
	os.Stdout.WriteString(
		"Usage: flacsfx [options...]\n"+
		" -o <file>     Destination file\n")
	os.Exit(err)
}

func main() {

	// Check for invalid number of arguments
	if len(os.Args) == 2 || len(os.Args) > 3 {
		help(-1)
	}

	// Push arguments to flag pointers
	for i := 1; i < len(os.Args); i++ {
		if strings.HasPrefix(os.Args[i], "-") {
			switch strings.TrimPrefix(os.Args[i], "-") {
				case "o":
					i++
					wavName = os.Args[i]
					continue
				default:
					help(-1)
			}
		} else {help(-1)}
	}

	//Parse FLAC byte stream
	flacStream, err := flac.New(bytes.NewReader(flacRaw))
	if err != nil {
		os.Stdout.WriteString("There was a problem parsing the FLAC stream.")
		os.Exit(1)
	}
	defer flacStream.Close()

	//Initialize WAV writer
	var wavWriter *os.File
	if wavName == "-" {
		wavWriter = os.Stdout
	} else {
		wavWriter, err = os.Create(wavName)
		if err != nil {
			os.Stdout.WriteString("There was a problem creating the new WAV file.")
			os.Exit(2)
		}
		os.Stdout.WriteString("Extracting \""+wavName+"\"...")
	}
	defer wavWriter.Close()

	//Initialize WAV encoder
	wavEncoder := wav.NewEncoder(
		wavWriter,
		int(flacStream.Info.SampleRate),
		int(flacStream.Info.BitsPerSample),
		int(flacStream.Info.NChannels),
		1,
	)
	defer wavEncoder.Close()

	//Decode FLAC samples, encode to WAV, and write to file
	var data []int
	for {
		//Decode FLAC audio samples
		frame, err := flacStream.ParseNext()
		if err != nil {
			if err == io.EOF {
				break
			}
			os.Stdout.WriteString("There was a problem decoding the FLAC stream.")
			os.Exit(3)
		}

		//Encode WAV audio samples
		data = data[:0]
		for i := 0; i < frame.Subframes[0].NSamples; i++ {
			for _, subframe := range frame.Subframes {
				sample := int(subframe.Samples[i])
				if frame.BitsPerSample == 8 {sample += 0x80}
				data = append(data, sample)
			}
		}		

		//Write the encoded WAV stream to file
		wavEncoder.Write(
			&audio.IntBuffer{
				Format: &audio.Format{
					NumChannels: int(flacStream.Info.NChannels),
					SampleRate: int(flacStream.Info.SampleRate),
				},
				Data: data,
				SourceBitDepth: int(flacStream.Info.BitsPerSample),
			},
		)
		if err != nil {
			os.Stdout.WriteString("There was a problem writing the WAV stream to the file.")
			os.Exit(4)
		}
	}
}
