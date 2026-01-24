package export

import (
	"strings"
	"testing"
)

const testSVG = `<svg xmlns="http://www.w3.org/2000/svg" width="200" height="100" viewBox="0 0 200 100">
  <rect width="200" height="100" fill="#3b82f6"/>
  <circle cx="100" cy="50" r="30" fill="#ffffff"/>
</svg>`

func TestExportSVG(t *testing.T) {
	opts := DefaultOptions()
	opts.Format = FormatSVG

	result, err := Export(testSVG, opts)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	if string(result) != testSVG {
		t.Error("SVG export should return original data")
	}
}

func TestExportPNG(t *testing.T) {
	opts := DefaultOptions()
	opts.Format = FormatPNG
	opts.Width = 400
	opts.Height = 200

	result, err := Export(testSVG, opts)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	// Check PNG header
	if len(result) < 8 {
		t.Fatal("PNG output too small")
	}
	if string(result[:8]) != "\x89PNG\r\n\x1a\n" {
		t.Error("Invalid PNG header")
	}
}

func TestExportJPEG(t *testing.T) {
	opts := DefaultOptions()
	opts.Format = FormatJPEG
	opts.Width = 400
	opts.Height = 200
	opts.Quality = 85

	result, err := Export(testSVG, opts)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	// Check JPEG header
	if len(result) < 2 {
		t.Fatal("JPEG output too small")
	}
	if result[0] != 0xFF || result[1] != 0xD8 {
		t.Error("Invalid JPEG header")
	}
}

func TestExportAutoDimensions(t *testing.T) {
	opts := DefaultOptions()
	opts.Format = FormatPNG
	// Width and Height = 0 should use SVG dimensions

	result, err := Export(testSVG, opts)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	if len(result) == 0 {
		t.Error("Export produced empty result")
	}
}

func TestExportAspectRatio(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
	}{
		{"width only", 400, 0},
		{"height only", 0, 300},
		{"both specified", 800, 600},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := DefaultOptions()
			opts.Format = FormatPNG
			opts.Width = tt.width
			opts.Height = tt.height

			result, err := Export(testSVG, opts)
			if err != nil {
				t.Fatalf("Export failed: %v", err)
			}

			if len(result) == 0 {
				t.Error("Export produced empty result")
			}
		})
	}
}

func TestGetMimeType(t *testing.T) {
	tests := []struct {
		format   Format
		expected string
	}{
		{FormatSVG, "image/svg+xml"},
		{FormatPNG, "image/png"},
		{FormatJPEG, "image/jpeg"},
		{FormatJPG, "image/jpeg"},
	}

	for _, tt := range tests {
		t.Run(string(tt.format), func(t *testing.T) {
			mime := GetMimeType(tt.format)
			if mime != tt.expected {
				t.Errorf("GetMimeType(%s) = %s, want %s", tt.format, mime, tt.expected)
			}
		})
	}
}

func TestGetFileExtension(t *testing.T) {
	tests := []struct {
		format   Format
		expected string
	}{
		{FormatSVG, ".svg"},
		{FormatPNG, ".png"},
		{FormatJPEG, ".jpg"},
		{FormatJPG, ".jpg"},
	}

	for _, tt := range tests {
		t.Run(string(tt.format), func(t *testing.T) {
			ext := GetFileExtension(tt.format)
			if ext != tt.expected {
				t.Errorf("GetFileExtension(%s) = %s, want %s", tt.format, ext, tt.expected)
			}
		})
	}
}

func TestParseFormat(t *testing.T) {
	tests := []struct {
		input    string
		expected Format
		wantErr  bool
	}{
		{"svg", FormatSVG, false},
		{"SVG", FormatSVG, false},
		{" svg ", FormatSVG, false},
		{"png", FormatPNG, false},
		{"PNG", FormatPNG, false},
		{"jpeg", FormatJPEG, false},
		{"jpg", FormatJPEG, false},
		{"JPEG", FormatJPEG, false},
		{"unknown", "", true},
		{"", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			format, err := ParseFormat(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Error("ParseFormat should have returned error")
				}
			} else {
				if err != nil {
					t.Errorf("ParseFormat failed: %v", err)
				}
				if format != tt.expected {
					t.Errorf("ParseFormat(%s) = %s, want %s", tt.input, format, tt.expected)
				}
			}
		})
	}
}

func TestInvalidSVG(t *testing.T) {
	opts := DefaultOptions()
	opts.Format = FormatPNG

	_, err := Export("not valid svg", opts)
	if err == nil {
		t.Error("Export should fail with invalid SVG")
	}
	// Error can be either parse failure or encoding failure
	if !strings.Contains(err.Error(), "failed to parse SVG") &&
	   !strings.Contains(err.Error(), "failed to encode PNG") {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestUnsupportedFormat(t *testing.T) {
	opts := DefaultOptions()
	opts.Format = "webp" // unsupported

	_, err := Export(testSVG, opts)
	if err == nil {
		t.Error("Export should fail with unsupported format")
	}
	if !strings.Contains(err.Error(), "unsupported format") {
		t.Errorf("Unexpected error: %v", err)
	}
}
