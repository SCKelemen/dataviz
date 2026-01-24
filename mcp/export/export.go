package export

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"strings"

	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

// Format represents an export format
type Format string

const (
	FormatSVG  Format = "svg"
	FormatPNG  Format = "png"
	FormatJPEG Format = "jpeg"
	FormatJPG  Format = "jpg"
)

// ExportOptions configures export settings
type ExportOptions struct {
	Format  Format
	Width   int // For raster formats, 0 = use SVG dimensions
	Height  int // For raster formats, 0 = use SVG dimensions
	Quality int // For JPEG, 0-100 (default 90)
	DPI     int // Dots per inch (default 96)
}

// DefaultOptions returns sensible defaults
func DefaultOptions() ExportOptions {
	return ExportOptions{
		Format:  FormatSVG,
		Quality: 90,
		DPI:     96,
	}
}

// Export converts SVG to the specified format
func Export(svgData string, opts ExportOptions) ([]byte, error) {
	// For SVG, just return the data
	if opts.Format == FormatSVG {
		return []byte(svgData), nil
	}

	// For raster formats, we need to rasterize
	return rasterize(svgData, opts)
}

// rasterize converts SVG to a raster image format
func rasterize(svgData string, opts ExportOptions) ([]byte, error) {
	// Parse SVG
	icon, err := oksvg.ReadIconStream(strings.NewReader(svgData))
	if err != nil {
		return nil, fmt.Errorf("failed to parse SVG: %w", err)
	}

	// Determine output dimensions
	width := opts.Width
	height := opts.Height

	if width == 0 || height == 0 {
		// Use SVG dimensions
		w := int(icon.ViewBox.W)
		h := int(icon.ViewBox.H)

		if width == 0 && height == 0 {
			// Use original dimensions
			width = w
			height = h
		} else if width == 0 {
			// Scale width to maintain aspect ratio
			width = int(float64(height) * float64(w) / float64(h))
		} else if height == 0 {
			// Scale height to maintain aspect ratio
			height = int(float64(width) * float64(h) / float64(w))
		}
	}

	// Set dimensions
	icon.SetTarget(0, 0, float64(width), float64(height))

	// Create image
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Create scanner
	scanner := rasterx.NewScannerGV(width, height, img, img.Bounds())

	// Rasterize
	raster := rasterx.NewDasher(width, height, scanner)
	icon.Draw(raster, 1.0)

	// Encode to target format
	var buf bytes.Buffer
	switch opts.Format {
	case FormatPNG:
		encoder := png.Encoder{CompressionLevel: png.DefaultCompression}
		if err := encoder.Encode(&buf, img); err != nil {
			return nil, fmt.Errorf("failed to encode PNG: %w", err)
		}
	case FormatJPEG, FormatJPG:
		quality := opts.Quality
		if quality == 0 {
			quality = 90
		}
		if quality < 1 {
			quality = 1
		}
		if quality > 100 {
			quality = 100
		}
		jpegOpts := &jpeg.Options{Quality: quality}
		if err := jpeg.Encode(&buf, img, jpegOpts); err != nil {
			return nil, fmt.Errorf("failed to encode JPEG: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported format: %s", opts.Format)
	}

	return buf.Bytes(), nil
}

// GetMimeType returns the MIME type for a format
func GetMimeType(format Format) string {
	switch format {
	case FormatSVG:
		return "image/svg+xml"
	case FormatPNG:
		return "image/png"
	case FormatJPEG, FormatJPG:
		return "image/jpeg"
	default:
		return "application/octet-stream"
	}
}

// GetFileExtension returns the file extension for a format
func GetFileExtension(format Format) string {
	switch format {
	case FormatSVG:
		return ".svg"
	case FormatPNG:
		return ".png"
	case FormatJPEG, FormatJPG:
		return ".jpg"
	default:
		return ".bin"
	}
}

// ParseFormat parses a format string
func ParseFormat(s string) (Format, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	switch s {
	case "svg":
		return FormatSVG, nil
	case "png":
		return FormatPNG, nil
	case "jpeg", "jpg":
		return FormatJPEG, nil
	default:
		return "", fmt.Errorf("unknown format: %s", s)
	}
}
