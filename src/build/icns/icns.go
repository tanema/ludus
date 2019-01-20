package icns

// copied from github.com/tonyhb/goicns to make simpler modifications

import (
	"bytes"
	"encoding/binary"
	"image/png"
	"os"

	"github.com/nfnt/resize"
)

const (
	icp4 = 16  // 16x16
	icp5 = 32  // 32x32
	icp6 = 64  // 64x64
	ic07 = 128 // 128x128
	ic08 = 256 // 256x256, 10.5
	ic09 = 512 // 512x512, 10.5

	ic11 = 32   // 16x16@2x retina
	ic12 = 64   // 32x32@2x retina
	ic13 = 256  // 128x128@2x retina
	ic14 = 512  // 256x256@2x retina
	ic10 = 1024 // 1024x1024/512x512@2x
)

var (
	icnsHeader = []byte{0x69, 0x63, 0x6e, 0x73}
	sizes      = []int{16, 32, 64, 128, 256, 512, 1024}
	sizeToType = map[int][]string{
		16:   {"icp4"},
		32:   {"icp5", "ic11"},
		64:   {"icp6", "ic12"},
		128:  {"ic07"},
		256:  {"ic08", "ic13"},
		512:  {"ic09", "ic14"},
		1024: {"ic10"},
	}
)

// ConvertPNG converts a png file to a icns file
func ConvertPNG(f *os.File) ([]byte, error) {
	img, err := png.Decode(f)
	if err != nil {
		return nil, err
	}

	icns := new(bytes.Buffer)
	for _, s := range sizes {
		imgBuf := new(bytes.Buffer)
		resized := resize.Resize(uint(s), uint(s), img, resize.MitchellNetravali)
		if err := png.Encode(imgBuf, resized); err != nil {
			return nil, err
		}

		lenByt := make([]byte, 4, 4)
		binary.BigEndian.PutUint32(lenByt, uint32(imgBuf.Len()+8))
		for _, ostype := range sizeToType[s] {
			if _, err := icns.Write([]byte(ostype)); err != nil {
				return nil, err
			}
			if _, err := icns.Write(lenByt); err != nil {
				return nil, err
			}
			if _, err := icns.Write(imgBuf.Bytes()); err != nil {
				return nil, err
			}
		}
	}

	lenByt := make([]byte, 4, 4)
	binary.BigEndian.PutUint32(lenByt, uint32(icns.Len()+8))

	data := icnsHeader
	data = append(data, lenByt...)
	data = append(data, icns.Bytes()...)

	return data, nil
}
