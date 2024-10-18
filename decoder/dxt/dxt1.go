package dxt

import (
	"image/color"

	. "github.com/ahephi/dds/decoder/dxt/internal"
)

type dxt1 struct {
	ColorDecoder
}

func (*dxt1) BlockSize() byte {
	return 8
}

func (d *dxt1) DecodeBlock(buffer []byte) {
	d.BlockColor(buffer[0:8:8])
}

func (d *dxt1) Pixel(index byte) color.Color {
	return d.PixelColor(index)
}
