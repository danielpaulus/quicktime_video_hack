package coremedia

import (
	"bytes"
	"encoding/binary"
	"os"
)

/*
Thank you http://soundfile.sapp.org/doc/WaveFormat/ for the amazing explanation of the WAV file format:

The canonical WAVE format starts with the RIFF header:

0         4   ChunkID          Contains the letters "RIFF" in ASCII form
                               (0x52494646 big-endian form).
4         4   ChunkSize        36 + SubChunk2Size, or more precisely:
                               4 + (8 + SubChunk1Size) + (8 + SubChunk2Size)
                               This is the size of the rest of the chunk
                               following this number.  This is the size of the
                               entire file in bytes minus 8 bytes for the
                               two fields not included in this count:
                               ChunkID and ChunkSize.
8         4   Format           Contains the letters "WAVE"
                               (0x57415645 big-endian form).

*/
type riffHeader struct {
	ChunkID   uint32
	ChunkSize uint32
	Format    uint32
}

//newRiffHeader get a RIFF header set up for creating a WAVE file
func newRiffHeader(size int) riffHeader {
	return riffHeader{ChunkID: 0x46464952, Format: 0x45564157, ChunkSize: uint32(36 + size)}
}

//serialize this RiffHeader into the given target bytes.Buffer
func (rh riffHeader) serialize(target *bytes.Buffer) error {
	return binary.Write(target, binary.LittleEndian, rh)
}

/*
The "WAVE" format consists of two subchunks: "fmt " and "data":
The "fmt " subchunk describes the sound data's format:

12        4   Subchunk1ID      Contains the letters "fmt "
                               (0x666d7420 big-endian form).
16        4   Subchunk1Size    16 for PCM.  This is the size of the
                               rest of the Subchunk which follows this number.
20        2   AudioFormat      PCM = 1 (i.e. Linear quantization)
                               Values other than 1 indicate some
                               form of compression.
22        2   NumChannels      Mono = 1, Stereo = 2, etc.
24        4   SampleRate       8000, 44100, etc.
28        4   ByteRate         == SampleRate * NumChannels * BitsPerSample/8
32        2   BlockAlign       == NumChannels * BitsPerSample/8
                               The number of bytes for one sample including
                               all channels. I wonder what happens when
                               this number isn't an integer?
34        2   BitsPerSample    8 bits = 8, 16 bits = 16, etc.
          2   ExtraParamSize   if PCM, then doesn't exist
          X   ExtraParams      space for extra parameters
*/
type fmtSubChunk struct {
	SubChunkID    uint32
	SubChunkSize  uint32
	AudioFormat   uint16
	NumChannels   uint16
	SampleRate    uint32
	ByteRate      uint32
	BlockAlign    uint16
	BitsPerSample uint16
}

//NewFmtSubChunk generates the Fmt Subchunk for creating a WAV file
func newFmtSubChunk() fmtSubChunk {
	result := fmtSubChunk{SubChunkID: 0x20746d66, SubChunkSize: 16, AudioFormat: 1, NumChannels: 2, SampleRate: 48000, BitsPerSample: 16}
	result.ByteRate = result.SampleRate * uint32(result.NumChannels) * uint32(result.BitsPerSample) / 8
	result.BlockAlign = result.NumChannels * (result.BitsPerSample / 8)
	return result
}

//Serialize this RiffHeader into the given target bytes.Buffer
func (fmsc fmtSubChunk) serialize(target *bytes.Buffer) error {
	return binary.Write(target, binary.LittleEndian, fmsc)
}

/*
The "data" subchunk contains the size of the data and the actual sound:

36        4   Subchunk2ID      Contains the letters "data"
                               (0x64617461 big-endian form).
40        4   Subchunk2Size    == NumSamples * NumChannels * BitsPerSample/8
                               This is the number of bytes in the data.
                               You can also think of this as the size
                               of the read of the subchunk following this
                               number.
44        *   Data             The actual sound data.
*/
func writeWavDataSubChunkHeader(target *bytes.Buffer, dataLength int) error {
	err := binary.Write(target, binary.BigEndian, uint32(0x64617461))
	if err != nil {
		return err
	}
	err = binary.Write(target, binary.LittleEndian, uint32(dataLength))
	if err != nil {
		return err
	}
	return nil
}

//WriteWavHeader creates a wave file header using the given length and writes it at the BEGINNING of the wavFile.
//Please make sure that the file has enough zero bytes before the audio data.
func WriteWavHeader(length int, wavFile *os.File) error {
	buffer := bytes.NewBuffer(make([]byte, 100))
	buffer.Reset()

	riffHeader := newRiffHeader(length)
	err := riffHeader.serialize(buffer)
	if err != nil {
		return err
	}

	fmtSubChunk := newFmtSubChunk()
	err = fmtSubChunk.serialize(buffer)
	if err != nil {
		return err
	}

	err = writeWavDataSubChunkHeader(buffer, length)
	if err != nil {
		return err
	}
	_, err = wavFile.WriteAt(buffer.Bytes(), 0)
	return err
}
