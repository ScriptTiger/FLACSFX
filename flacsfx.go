package main

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/mewkiz/flac"
	"github.com/mewkiz/flac/frame"
	"github.com/ScriptTiger/mixerInG"
)

// Memory structures

// Structure to store flac stream information
type flacInfo struct {
	start int64
	size int64
	file *os.File
	stream *flac.Stream
	outFile *os.File
	wavEnc *wav.Encoder
	title string
	format *audio.Format
	sampleRate int
	numChans int
	bitDepth int
	numSamples int
	numDecodedSamples int
	frame *frame.Frame
	currentSubframe int
	numFrameSamples int
	currentSample int
	rawBuffer []int
	rawBufferSize int
	bufferSize int
	intBuffer *audio.IntBuffer
	floatBuffer *audio.FloatBuffer
}

// Functions for flacInfo

// Function to create a new flacInfo
func newFlacInfo() (*flacInfo) {return &flacInfo{currentSubframe: -1, currentSample: -1}}

// Function to flush flacInfo raw buffer to intBuffer
func (i *flacInfo) flush() (bool) {
	if i.bufferSize != 0 {return false}
	i.bufferSize = i.rawBufferSize
	if i.bufferSize == 0 {return false}
	if i.bufferSize != cap(i.rawBuffer) {
		i.rawBuffer = i.rawBuffer[:i.bufferSize]
		i.intBuffer.Data = i.intBuffer.Data[:i.bufferSize]
		i.floatBuffer.Data = i.floatBuffer.Data[:i.bufferSize]
	}
	i.intBuffer.Data = i.rawBuffer
	i.rawBufferSize = 0
	i.floatBuffer = i.intBuffer.AsFloatBuffer()
	return true
}

// Function to buffer samples
func (i *flacInfo) bufferSample(sample int) (bool) {
	i.rawBuffer[i.rawBufferSize] = sample
	i.rawBufferSize++
	if i.rawBufferSize == cap(i.rawBuffer) {
		i.flush()
		return true
	} else {return false}
}

// Function to write buffer to encoder
func (i *flacInfo) writeBuffer() {
	if i.bufferSize == 0 {return}
	if i.bufferSize != cap(i.intBuffer.Data) {i.intBuffer.Data = i.intBuffer.Data[:i.bufferSize]}
	i.wavEnc.Write(i.intBuffer)
	i.bufferSize = 0
}

// Function to parse next frame
func (i *flacInfo) walkFrames() (bool) {
	var err error
	if i.currentSample != -1 && i.currentSubframe != -1 {return true}
	i.frame, err = i.stream.ParseNext()
	if err != nil {return false}
	i.numFrameSamples = i.frame.Subframes[0].NSamples
	return true
}

// Function to walk samples
func (i *flacInfo) walkSamples() (bool) {
	if i.currentSubframe != -1 {return true}
	i.currentSample++
	i.numDecodedSamples++
	if i.currentSample == i.numFrameSamples {
		i.currentSample = -1
		return false
	}
	return true
}

// Function to walk subframes
func (i *flacInfo) walkSubframes() (bool) {
	i.currentSubframe++
	if i.currentSubframe == i.numChans {
		i.currentSubframe = -1
		return false
	}
	return true
}

// Function to systematically walk frames, samples, and subframes
func (i *flacInfo) walk() (bool) {
	if i.numDecodedSamples >= i.numSamples && i.currentSample == -1 && i.currentSubframe == -1 {return false}
	for i.walkFrames() {for i.walkSamples() {for i.walkSubframes() {return true}}}
	return false
}

// General functions

// Function to display help text and exit
func help(err int) {
	os.Stdout.WriteString(
		"Usage: flacsfx [options...]\n"+
		" -flac               Output FLAC\n"+
		" -mix                Output mix\n"+
		" -o <directory|file> Destination directory, or file for mix\n"+
		" -b <16|24|32>       Bit depth of mix\n"+
		" -info               Show stream info",
	)
	os.Exit(err)
}

// Function to translate NChannels to channel layout
func layoutLookup(nChannels int) (layout string) {
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
		default:
			layout = "unknown"
	}
	return
}

// Function to scan a buffered I/O reader for a specific block of data
func fileScanner(block *[]byte, reader *bufio.Reader) (err error) {
	for i := 0; i < 6; i++ {(*block)[i] = (*block)[i+1]}
	(*block)[6], err = reader.ReadByte()
	return err
}

// Function to close all opened and indexed files before exiting
func exit(index *[]flacInfo, err int) {
	for _, flacStream := range (*index) {
		flacStream.file.Close()
		if flacStream.outFile != nil {flacStream.outFile.Close()}
	}
	os.Exit(err)
}

// Function to write headers to all wav files and close wav encoders
func wavClose(index *[]flacInfo) {for _, flacStream := range (*index) {flacStream.wavEnc.Close()}}

func main() {

	// Check for invalid number of arguments
	if len(os.Args) > 6 {help(1)}

	// Argument declarations
	var (
		outName *string
		flacenc bool
		mix bool
		info bool
		bitDepth int
		err error
	)

	// Argument handling
	for i := 1; i < len(os.Args); i++ {
		if strings.HasPrefix(os.Args[i], "-") {
			switch strings.TrimPrefix(os.Args[i], "-") {
				case "flac":
					if flacenc {help(2)}
					flacenc = true
					continue
				case "mix":
					if mix {help(3)}
					mix = true
					continue
				case "o":
					if outName != nil {help(4)}
					i++
					outName = &os.Args[i]
					continue
				case "b":
					if bitDepth > 0 {help(5)}
					i++
					bitDepth, err = strconv.Atoi(os.Args[i])
					if err != nil ||
					(bitDepth != 16 &&
					bitDepth != 24 &&
					bitDepth != 32) {help(6)}
					continue
				case "info":
					if info {help(7)}
					info = true
					break
				case "":
					if outName != nil {help(8)}
					outName = &os.Args[i]
					continue
				default:
					help(9)
			}
		} else {
			if outName != nil {help(10)}
			outName = &os.Args[i]
			continue
		}
	}

	// Locate executable
	filePath, _ := os.Executable()
	filePath, _ = filepath.EvalSymlinks(filePath)

	// Open file
	sfxFile, _ := os.Open(filePath)

	// Set default output directory if not given as argument
	if outName == nil {
		outName = new(string)
		*outName = strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
		if mix {*outName = *outName+".wav"}
	}

	// Create output directory, or do nothing if it already exists
	if !info && *outName != "-" {
		if mix {os.MkdirAll(filepath.Dir(*outName), 0755)
		} else {os.MkdirAll(*outName, 0755)}
	}

	// Determine total file size
	fileInfo, _ := sfxFile.Stat()
	sfxTotalSize := fileInfo.Size()

	// Set common buffer capacity
	bufferCap := 8000

	// Build the index of embedded FLAC streams
	readPoint := int64(1500000)
	sfxFile.Seek(readPoint, io.SeekStart)
	sfxReader := bufio.NewReader(sfxFile)
	block := make([]byte, 7)
	isMixable := true
	var titleBuilder strings.Builder
	var index []flacInfo
	for {
		i := len(index)
		if fileScanner(&block, sfxReader) != nil {
			if i == 0 {
				os.Stdout.WriteString("No emedded FLAC stream found.\n")
				sfxFile.Close()
				exit(&index, 11)
			}
			index[i-1].size = sfxTotalSize-index[i-1].start
			if i < 2 {isMixable = false}
			break
		}
		if string(block) == "fLaC\x00\x00\x00" {
			currentPoint := readPoint-6
			index = append(index, *newFlacInfo())
			index[i].start = currentPoint
			index[i].file, _ = os.Open(filePath)
			defer index[i].file.Close()
			index[i].file.Seek(currentPoint, io.SeekStart)
			for i, _ := range block {block[i] = '\x00'}
			flacFileReader := bufio.NewReader(index[i].file)
			for {
				if fileScanner(&block, flacFileReader) != nil {
					index[i].title = strconv.Itoa(i)+"_"+strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
					break
				}
				if string(block) == "\x00TITLE=" {
					for {
						titleByte, _ := flacFileReader.ReadByte()
						if titleByte == '\x00' {
							titleSize := titleBuilder.Len()
							if titleSize == 0 {
								index[i].title = strconv.Itoa(i)+"_"+strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
							} else {
								index[i].title = string([]byte(titleBuilder.String())[0:titleSize-1])
								titleBuilder.Reset()
							}
							break
						}
						titleBuilder.WriteByte(titleByte)
					}
					break
				}
				if string(block) == strings.Repeat("\x00", 16) {
					index[i].title = strconv.Itoa(i)+"_"+strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
					break
				}
			}
			index[i].file.Seek(currentPoint, io.SeekStart)
			index[i].stream, err = flac.New(index[i].file)
			if err != nil {
				index = index[:i]
				continue
			}
			index[i].sampleRate = int(index[i].stream.Info.SampleRate)
			index[i].numChans = int(index[i].stream.Info.NChannels)
			index[i].bitDepth = int(index[i].stream.Info.BitsPerSample)
			index[i].numSamples = int(index[i].stream.Info.NSamples)
			index[i].format = &audio.Format{
				NumChannels: index[i].numChans,
				SampleRate: index[i].sampleRate,
			}
			if !info {
				if flacenc {index[i].outFile, err = os.Create(filepath.Join(*outName, index[i].title)+".flac")
				} else {
					if !mix {
						index[i].outFile, err = os.Create(filepath.Join(*outName, index[i].title)+".wav")
						index[i].wavEnc = wav.NewEncoder(
							index[i].outFile,
							index[i].sampleRate,
							index[i].bitDepth,
							index[i].numChans,
							1,
						)
					}
					index[i].rawBuffer = make([]int, bufferCap)
					index[i].intBuffer = &audio.IntBuffer{Format: index[i].format, Data: make([]int, bufferCap)}
					index[i].floatBuffer = &audio.FloatBuffer{Format: index[i].format, Data: make([]float64, bufferCap)}
				}
				if err != nil {
					os.Stdout.WriteString ("A problem was encountered while attempting to write a file.\n")
					exit(&index, 12)
				}
			}
			if i > 0 {
				if *index[i].format != *index[0].format {isMixable = false}
				index[i-1].size = currentPoint-index[i-1].start
			}
		}
		readPoint++
	}

	// Close sfxFile as it's no longer needed
	sfxFile.Close()

	// Display stream info if requested and exit
	if info {
		os.Stdout.WriteString("mixable="+strconv.FormatBool(isMixable)+"\n")
		for i, flacStream := range index {
			os.Stdout.WriteString(
				"[STREAM]\n"+
				"index="+strconv.Itoa(i)+"\n"+
				"title="+flacStream.title+"\n"+
				"codec_name=flac\n"+
				"codec_long_name=FLAC (Free Lossless Audio Codec)\n"+
				"codec_type=audio\n"+
				"sample_rate="+strconv.Itoa(flacStream.sampleRate)+"\n"+
				"channels="+strconv.Itoa(flacStream.numChans)+"\n"+
				"channel_layout="+layoutLookup(flacStream.numChans)+"\n"+
				"time_base=1/"+strconv.Itoa(flacStream.sampleRate)+"\n"+
				"duration_ts="+strconv.Itoa(flacStream.numSamples)+"\n"+
				"duration="+strconv.FormatFloat(float64(flacStream.numSamples)/float64(flacStream.sampleRate), 'f', -1, 64)+"\n"+
				"bits_per_raw_sample="+strconv.Itoa(flacStream.bitDepth)+"\n"+
				"[/STREAM]\n",
			)
		}
		exit(&index, 0)
	}

	// Store number of tracks for quick reference
	numTracks := len(index)

	// Dump flac files if requested and exit
	if flacenc {
		for _, flacStream := range index {
			os.Stdout.WriteString("Extracting \""+flacStream.title+".flac\"...\n")
			flacStream.file.Seek(flacStream.start, io.SeekStart)
			io.CopyN(flacStream.outFile, flacStream.file, flacStream.size)
		}
		exit(&index, 0)
	}

	// Reject mix requests if tracks are not mixable
	if mix && !isMixable {
		if numTracks < 2 {
			os.Stdout.WriteString("You need at least 2 tracks to mix.\n")
		} else {os.Stdout.WriteString("The tracks cannot be mixed due to having incompatible formats.\n")}
		exit(&index, 13)
	}

	// Set up mix properties if needed
	var (
		outFile *os.File
		wavEnc *wav.Encoder
		mixFloatBuffer *audio.FloatBuffer
		sourceTracks []*mixerInG.TrackInfo
	)
	if mix {
		format := index[0].format
		sampleRate := index[0].sampleRate
		numChans := index[0].numChans
		if bitDepth == 0 {bitDepth = 24}

		// Create mix out
		if *outName == "-" {outFile = os.Stdout
		} else {
			outFile, err = os.Create(*outName)
			if err != nil {
				os.Stdout.WriteString("A problem was encountered while attempting to create the mix file.\n")
				exit(&index, 14)
			}
		}

		// Create mix wav encoder
		wavEnc = wav.NewEncoder(
			outFile,
			sampleRate,
			bitDepth,
			numChans,
			1,
		)

		// Create mix buffer
		mixFloatBuffer = &audio.FloatBuffer{Format: format, Data: make([]float64, bufferCap)}

		// Initialize TrackInfos slice for source tracks
		sourceTracks = make([]*mixerInG.TrackInfo, numTracks)
	}

	// Decode FLAC audio samples to buffers and feed to individual track wav encoders, or to mix and wav encoder if mix requested
	for {
		for i, _ := range index {
			if !mix {os.Stdout.WriteString("Extracting \""+index[i].title+".wav\"...\n")}
			for index[i].walk() {
				sample := int(index[i].frame.Subframes[index[i].currentSubframe].Samples[index[i].currentSample])
				if index[i].frame.BitsPerSample == 8 {sample += 0x80}
				if index[i].bufferSample(sample) {
					if mix {break
					} else {index[i].writeBuffer()}
				}
			}
			if !mix {
				index[i].flush()
				index[i].writeBuffer()
			}
		}
		if !mix {break}

		// Mix tracks to a mix track
		for i, track := range index {
			sourceTracks[i] = &mixerInG.TrackInfo{
				BitDepth: track.bitDepth,
				BufferSize: track.bufferSize,
				FloatBuffer: track.floatBuffer,
			}
		}
		mixTrackSize := mixerInG.Mix(mixFloatBuffer, sourceTracks, bitDepth, false)
		wavEnc.Write(mixFloatBuffer.AsIntBuffer())
		if mixTrackSize < bufferCap {break}
		for i, _ := range sourceTracks {index[i].bufferSize = 0}

	}
	if mix {
		wavEnc.Close()
		outFile.Close()
	} else {wavClose(&index)}
	exit(&index, 0)
}