package mp4

import (
	"fmt"
	"io"
	"io/ioutil"
)

// UnknownBox - box that we don't know how to parse
type UnknownBox struct {
	name       string
	size       uint64
	notDecoded []byte
}

// DecodeUnknown - decode an unknown box
func DecodeUnknown(hdr *boxHeader, startPos uint64, r io.Reader) (Box, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return &UnknownBox{hdr.name, hdr.size, data}, nil
}

// Type - return box type
func (b *UnknownBox) Type() string {
	return b.name
}

// Size - return calculated size
func (b *UnknownBox) Size() uint64 {
	return b.size
}

// Encode - write box to w
func (b *UnknownBox) Encode(w io.Writer) error {
	err := EncodeHeader(b, w)
	if err != nil {
		return err
	}
	_, err = w.Write(b.notDecoded)
	return err
}

func (b *UnknownBox) Dump(w io.Writer, indent, indentStep string) error {
	_, err := fmt.Fprintf(w, "%s%s size=%d\n", indent, b.Type(), b.Size())
	return err
}
