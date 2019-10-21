package mp4

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
)

// StscBox - Sample To Chunk Box (stsc - mandatory)
//
// A chunk contains samples. This table defines to which chunk a sample is associated.
// Each entry is defined by :
//
//   * first chunk : all chunks starting at this index up to the next first chunk have the same sample count/description
//   * samples per chunk : number of samples in the chunk
//   * description id : description (see the sample description box - stsd)
type StscBox struct {
	Version             byte
	Flags               uint32
	FirstChunk          []uint32
	SamplesPerChunk     []uint32
	SampleDescriptionID []uint32
}

// DecodeStsc - box-specific decode
func DecodeStsc(hdr *boxHeader, startPos uint64, r io.Reader) (Box, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	versionAndFlags := binary.BigEndian.Uint32(data[0:4])
	b := &StscBox{
		Version:             byte(versionAndFlags >> 24),
		Flags:               versionAndFlags & flagsMask,
		FirstChunk:          []uint32{},
		SamplesPerChunk:     []uint32{},
		SampleDescriptionID: []uint32{},
	}
	ec := binary.BigEndian.Uint32(data[4:8])
	for i := 0; i < int(ec); i++ {
		fc := binary.BigEndian.Uint32(data[(8 + 12*i):(12 + 12*i)])
		spc := binary.BigEndian.Uint32(data[(12 + 12*i):(16 + 12*i)])
		sdi := binary.BigEndian.Uint32(data[(16 + 12*i):(20 + 12*i)])
		b.FirstChunk = append(b.FirstChunk, fc)
		b.SamplesPerChunk = append(b.SamplesPerChunk, spc)
		b.SampleDescriptionID = append(b.SampleDescriptionID, sdi)
	}
	return b, nil
}

// Type box-specific type
func (b *StscBox) Type() string {
	return "stsc"
}

// Size - box-specfic size
func (b *StscBox) Size() uint64 {
	return uint64(boxHeaderSize + 8 + len(b.FirstChunk)*12)
}

// Dump - box-specific dump
func (b *StscBox) Dump() {
	fmt.Println("Sample to Chunk:")
	for i := range b.SamplesPerChunk {
		fmt.Printf(" #%d : %d samples per chunk starting @chunk #%d \n", i, b.SamplesPerChunk[i], b.FirstChunk[i])
	}
}

// Encode - box-specific encode
func (b *StscBox) Encode(w io.Writer) error {
	err := EncodeHeader(b, w)
	if err != nil {
		return err
	}
	buf := makebuf(b)
	sw := NewSliceWriter(buf)
	versionAndFlags := (uint32(b.Version) << 24) + b.Flags
	sw.WriteUint32(versionAndFlags)
	sw.WriteUint32(uint32(len(b.FirstChunk)))
	for i := range b.FirstChunk {
		sw.WriteUint32(b.FirstChunk[i])
		sw.WriteUint32(b.SamplesPerChunk[i])
		sw.WriteUint32(b.SampleDescriptionID[i])
	}
	_, err = w.Write(buf)
	return err
}