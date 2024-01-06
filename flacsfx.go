package main

import (
	"bytes"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/mewkiz/flac"
)

// Function to display help text and exit
func help(err int) {
	os.Stdout.WriteString(
		"Usage: flacsfx [options...]\n"+
		" -o <file>     Destination file\n"+
		" -info         Show stream info",
	)
	os.Exit(err)
}

// Function to translate NChannels to channel layout
func layoutLookup(nChannels uint8) (layout string) {
	switch nChannels {
		case 1:
			layout = "mono"
		case 2:
			layout = "stereo"
		case 3:
			layout = "2.1"
		case 4:
			layout = "4.0"
		case 5:
			layout = "5.0"
		case 6:
			layout = "5.1"
		case 7:
			layout = "6.1"
		case 8:
			layout = "7.1"
	}
	return
}

func main() {

	// Check for invalid number of arguments
	if len(os.Args) > 4 {
		help(-1)
	}

	// Initialize uninitialized flags
	var info bool

	// Push arguments to flag pointers
	for i := 1; i < len(os.Args); i++ {
		if strings.HasPrefix(os.Args[i], "-") {
			switch strings.TrimPrefix(os.Args[i], "-") {
				case "o":
					i++
					wavName = os.Args[i]
					continue
				case "info":
					info = true
					continue
				default:
					help(-1)
			}
		} else {help(-1)}
	}

	// Parse FLAC byte stream
	flacStream, err := flac.New(bytes.NewReader(flacRaw))
	if err != nil {
		os.Stdout.WriteString("There was a problem parsing the FLAC stream.")
		os.Exit(1)
	}
	defer flacStream.Close()

	// Display stream info and exit
	if info {
		os.Stdout.WriteString(
			"codec_name=flac\n"+
			"codec_long_name=FLAC (Free Lossless Audio Codec)\n"+
			"codec_type=audio\n"+
			"sample_rate="+strconv.FormatUint(uint64(flacStream.Info.SampleRate), 10)+"\n"+
			"channels="+strconv.Itoa(int(flacStream.Info.NChannels))+"\n"+
			"channel_layout="+layoutLookup(flacStream.Info.NChannels)+"\n"+
			"time_base=1/"+strconv.FormatUint(uint64(flacStream.Info.SampleRate), 10)+"\n"+
			"duration_ts="+strconv.FormatUint(uint64(flacStream.Info.NSamples), 10)+"\n"+
			"duration="+strconv.FormatFloat(float64(flacStream.Info.NSamples)/float64(flacStream.Info.SampleRate), 'f', -1, 64)+"\n"+
			"bits_per_raw_sample="+strconv.Itoa(int(flacStream.Info.BitsPerSample)),
		)
		os.Exit(0)
	}

	// Initialize WAV writer
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

	// Initialize WAV encoder
	wavEncoder := wav.NewEncoder(
		wavWriter,
		int(flacStream.Info.SampleRate),
		int(flacStream.Info.BitsPerSample),
		int(flacStream.Info.NChannels),
		1,
	)
	defer wavEncoder.Close()

	// Decode FLAC samples, encode to WAV, and write to file
	var data []int
	for {
		// Decode FLAC audio samples
		frame, err := flacStream.ParseNext()
		if err != nil {
			if err == io.EOF {
				break
			}
			os.Stdout.WriteString("There was a problem decoding the FLAC stream.")
			os.Exit(3)
		}

		// Encode WAV audio samples
		data = data[:0]
		for i := 0; i < frame.Subframes[0].NSamples; i++ {
			for _, subframe := range frame.Subframes {
				sample := int(subframe.Samples[i])
				if frame.BitsPerSample == 8 {sample += 0x80}
				data = append(data, sample)
			}
		}		

		// Write the encoded WAV stream to file
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
