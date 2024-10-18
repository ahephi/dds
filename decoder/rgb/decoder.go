package rgb

import (
	"encoding/binary"
	"image"
	"io"
	"math/bits"

	"github.com/ahephi/dds/header"
)

type Decoder struct {
	ddpfHeader header.DDPFHeader
	bounds     image.Point
}

func New(header *header.Header) *Decoder {
	return &Decoder{
		ddpfHeader: header.DDPFHeader,
		bounds:     image.Pt(int(header.Width), int(header.Height)),
	}
}

func (d *Decoder) Decode(r io.Reader) (image.Image, error) {
	rgba := image.NewNRGBA(image.Rectangle{Max: d.bounds})
	if rgba.Rect.Empty() {
		return rgba, nil
	}

	pixelBytes := int(d.ddpfHeader.RgbBitCount / 8)

	rMask := d.ddpfHeader.RBitMask
	rBitCount := bits.OnesCount32(rMask)
	rDistance := uint32(bits.TrailingZeros32(rMask))
	rMax := float64((uint64(1) << rBitCount) - 1)

	gMask := d.ddpfHeader.GBitMask
	gBitCount := bits.OnesCount32(gMask)
	gDistance := uint32(bits.TrailingZeros32(gMask))
	gMax := float64((uint64(1) << gBitCount) - 1)

	bMask := d.ddpfHeader.BBitMask
	bBitCount := bits.OnesCount32(bMask)
	bDistance := uint32(bits.TrailingZeros32(bMask))
	bMax := float64((uint64(1) << bBitCount) - 1)

	aMask := d.ddpfHeader.ABitMask
	aBitCount := bits.OnesCount32(aMask)
	aDistance := uint32(bits.TrailingZeros32(aMask))
	aMax := float64((uint64(1) << aBitCount) - 1)

	switch d.ddpfHeader.PixelFlags.F {
	case header.DDPFAlphaPixels | header.DDPFRGB:
		for y := 0; y < d.bounds.Y; y++ {
			for x := 0; x < d.bounds.X; x++ {
				raw, err := nextPixel(r, pixelBytes)
				if err != nil {
					return nil, err
				}

				offset := rgba.PixOffset(x, y)

				rv := float64((raw&rMask)>>rDistance) / rMax
				gv := float64((raw&gMask)>>gDistance) / gMax
				bv := float64((raw&bMask)>>bDistance) / bMax
				av := float64((raw&aMask)>>aDistance) / aMax
				rgba.Pix[offset+0] = byte(rv * 0xff)
				rgba.Pix[offset+1] = byte(gv * 0xff)
				rgba.Pix[offset+2] = byte(bv * 0xff)
				rgba.Pix[offset+3] = byte(av * 0xff)
			}
		}
	case header.DDPFRGB:
		for y := 0; y < d.bounds.Y; y++ {
			for x := 0; x < d.bounds.X; x++ {
				raw, err := nextPixel(r, pixelBytes)
				if err != nil {
					return nil, err
				}

				offset := rgba.PixOffset(x, y)

				rv := float64((raw&rMask)>>rDistance) / rMax
				gv := float64((raw&gMask)>>gDistance) / gMax
				bv := float64((raw&bMask)>>bDistance) / bMax
				rgba.Pix[offset+0] = byte(rv * 0xff)
				rgba.Pix[offset+1] = byte(gv * 0xff)
				rgba.Pix[offset+2] = byte(bv * 0xff)
				rgba.Pix[offset+3] = 0xff
			}
		}
	}

	return rgba, nil
}

func nextPixel(r io.Reader, bytes int) (uint32, error) {
	v := [4]byte{}

	if _, err := r.Read(v[:bytes]); err != nil {
		return 0, err
	}

	a := binary.LittleEndian.Uint32(v[:])
	return a, nil
}
