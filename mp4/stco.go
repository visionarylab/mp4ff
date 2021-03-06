package mp4

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
)

// StcoBox - Chunk Offset Box (stco - mandatory)
//
// Contained in : Sample Table box (stbl)
//
// This is the 32bits version of the box, the 64bits version (co64) is not decoded.
//
// The table contains the offsets (starting at the beginning of the file) for each chunk of data for the current track.
// A chunk contains samples, the table defining the allocation of samples to each chunk is stsc.
type StcoBox struct {
	Version     byte
	Flags       uint32
	ChunkOffset []uint32
}

// DecodeStco - box-specific decode
func DecodeStco(hdr *boxHeader, startPos uint64, r io.Reader) (Box, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	versionAndFlags := binary.BigEndian.Uint32(data[0:4])
	b := &StcoBox{
		Version:     byte(versionAndFlags >> 24),
		Flags:       versionAndFlags & flagsMask,
		ChunkOffset: []uint32{},
	}
	ec := binary.BigEndian.Uint32(data[4:8])
	for i := 0; i < int(ec); i++ {
		chunk := binary.BigEndian.Uint32(data[(8 + 4*i):(12 + 4*i)])
		b.ChunkOffset = append(b.ChunkOffset, chunk)
	}
	return b, nil
}

// Type - box-specific type
func (b *StcoBox) Type() string {
	return "stco"
}

// Size - box-specific size
func (b *StcoBox) Size() uint64 {
	return uint64(boxHeaderSize + 8 + len(b.ChunkOffset)*4)
}

// Encode - box-specific encode
func (b *StcoBox) Encode(w io.Writer) error {
	err := EncodeHeader(b, w)
	if err != nil {
		return err
	}
	buf := makebuf(b)
	sw := NewSliceWriter(buf)
	versionAndFlags := (uint32(b.Version) << 24) + b.Flags
	sw.WriteUint32(versionAndFlags)
	sw.WriteUint32(uint32(len(b.ChunkOffset)))
	for i := range b.ChunkOffset {
		sw.WriteUint32(b.ChunkOffset[i])
	}
	_, err = w.Write(buf)
	return err
}

func (s *StcoBox) Dump(w io.Writer, indent, indentStep string) error {
	_, err := fmt.Fprintf(w, "%s%s size=%d\n", indent, s.Type(), s.Size())
	return err
}
