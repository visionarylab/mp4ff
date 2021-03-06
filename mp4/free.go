package mp4

import (
	"fmt"
	"io"
	"io/ioutil"
)

// FreeBox - Free Box
type FreeBox struct {
	notDecoded []byte
}

// DecodeFree - box-specific decode
func DecodeFree(hdr *boxHeader, startPos uint64, r io.Reader) (Box, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return &FreeBox{data}, nil
}

// Type - box type
func (b *FreeBox) Type() string {
	return "free"
}

// Size - calculated size of box
func (b *FreeBox) Size() uint64 {
	return uint64(boxHeaderSize + len(b.notDecoded))
}

// Encode - write box to w
func (b *FreeBox) Encode(w io.Writer) error {
	err := EncodeHeader(b, w)
	if err != nil {
		return err
	}
	_, err = w.Write(b.notDecoded)
	return err
}

func (b *FreeBox) Dump(w io.Writer, indent, indentStep string) error {
	_, err := fmt.Fprintf(w, "%s%s size=%d\n", indent, b.Type(), b.Size())
	return err
}
