package main

import (
	"bytes"
	"io"
	"os"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/mewkiz/flac"
)

func main() {
	//Parse FLAC byte stream
	stream, err := flac.New(bytes.NewReader(flacRaw))
	if err != nil {
		os.Stdout.WriteString("There was a problem parsing the FLAC stream.")
		os.Exit(1)
	}
	defer stream.Close()

	//Initialize wav file writer
	wavWriter, err := os.Create(wavName)
	if err != nil {
		os.Stdout.WriteString("There was a problem creating the new WAV file.")
		os.Exit(2)
	}
	defer wavWriter.Close()

	//Initialize new wav encoder
	wavEncoder := wav.NewEncoder(wavWriter, int(stream.Info.SampleRate), int(stream.Info.BitsPerSample), int(stream.Info.NChannels), 1)
	defer wavEncoder.Close()

	os.Stdout.WriteString("Extracting \""+wavName+"\"...")

	//Decode FLAC samples, encode to WAV, and write to file
	var data []int
	for {
		//Decode FLAC audio samples
		frame, err := stream.ParseNext()
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
		buf := &audio.IntBuffer{
			Format: &audio.Format{
				NumChannels: int(stream.Info.NChannels),
				SampleRate: int(stream.Info.SampleRate),
			},
			Data: data,
			SourceBitDepth: int(stream.Info.BitsPerSample),
		}
		wavEncoder.Write(buf)
		if err != nil {
			os.Stdout.WriteString("There was a problem writing the buffered WAV stream to the file.")
			os.Exit(4)
		}
	}
}
