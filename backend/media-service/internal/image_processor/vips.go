package image_processor

import (
	"context"

	"github.com/davidbyttow/govips/v2/vips"
)

type Processor struct{}

func New() *Processor {
	vips.Startup(nil)
	return &Processor{}
}

func (p *Processor) ResizeToWidth(
	ctx context.Context,
	input []byte,
	width int,
) ([]byte, int, int, error) {

	select {
	case <-ctx.Done():
		return nil, 0, 0, ctx.Err()
	default:
	}

	img, err := vips.NewImageFromBuffer(input)
	if err != nil {
		return nil, 0, 0, err
	}
	defer img.Close()

	origWidth := img.Width()
	origHeight := img.Height()

	// если уже меньше — не увеличиваем
	if origWidth <= width {
		width = origWidth
	}

	scale := float64(width) / float64(origWidth)
	height := int(float64(origHeight) * scale)

	err = img.Resize(scale, vips.KernelLanczos3)
	if err != nil {
		return nil, 0, 0, err
	}

	buf, _, err := img.ExportWebp(&vips.WebpExportParams{
		Quality: 85,
	})
	if err != nil {
		return nil, 0, 0, err
	}

	return buf, width, height, nil
}

func (p *Processor) Close() {
	vips.Shutdown()
}
