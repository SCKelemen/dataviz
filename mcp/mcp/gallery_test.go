package mcp

import (
	"context"
	"testing"

	"github.com/SCKelemen/dataviz/internal/gallery"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestHandleGallery(t *testing.T) {
	server, err := NewServer()
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	tests := []struct {
		name        string
		galleryType string
		wantErr     bool
	}{
		{
			name:        "bar gallery",
			galleryType: "bar",
			wantErr:     false,
		},
		{
			name:        "line gallery",
			galleryType: "line",
			wantErr:     false,
		},
		{
			name:        "scatter gallery",
			galleryType: "scatter",
			wantErr:     false,
		},
		{
			name:        "radar gallery",
			galleryType: "radar",
			wantErr:     false,
		},
		{
			name:        "invalid gallery type",
			galleryType: "invalid",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := map[string]interface{}{
				"gallery_type": tt.galleryType,
			}

			request := createTestRequest(t, "generate_gallery", args)
			result, err := server.handleGallery(context.Background(), request)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Error("Expected result but got nil")
				return
			}

			if len(result.Content) != 1 {
				t.Errorf("Expected 1 content item, got %d", len(result.Content))
				return
			}

			textContent, ok := result.Content[0].(*mcp.TextContent)
			if !ok {
				t.Error("Expected TextContent")
				return
			}

			if textContent.Text == "" {
				t.Error("Expected non-empty SVG content")
			}

			// Verify it's wrapped in markdown code block
			if len(textContent.Text) < 10 {
				t.Error("SVG content too short")
			}

			if textContent.Text[:7] != "```svg\n" {
				t.Error("SVG not properly wrapped in markdown code block")
			}
		})
	}
}

func TestGalleryRegistry(t *testing.T) {
	// Verify all expected gallery types are registered
	expectedTypes := []string{
		"bar", "area", "stacked-area", "lollipop", "histogram",
		"pie", "boxplot", "violin", "treemap", "icicle", "ridgeline",
		"line", "scatter", "connected-scatter", "statcard",
		"radar", "streamchart", "candlestick", "sunburst", "circle-packing",
		"heatmap",
	}

	for _, galleryType := range expectedTypes {
		t.Run(galleryType, func(t *testing.T) {
			config, ok := gallery.GalleryRegistry[galleryType]
			if !ok {
				t.Errorf("Gallery type %s not found in registry", galleryType)
				return
			}

			if config.Name != galleryType {
				t.Errorf("Expected name %s, got %s", galleryType, config.Name)
			}

			if config.Title == "" {
				t.Error("Gallery title is empty")
			}

			if config.Layout == nil {
				t.Error("Gallery layout is nil")
			}

			if len(config.Variants) == 0 {
				t.Error("Gallery has no variants")
			}
		})
	}
}

func TestGalleryGeneration(t *testing.T) {
	// Test that gallery generation actually produces valid SVG
	testCases := []string{"bar", "line", "radar"}

	for _, galleryType := range testCases {
		t.Run(galleryType, func(t *testing.T) {
			config, ok := gallery.GalleryRegistry[galleryType]
			if !ok {
				t.Fatalf("Gallery type %s not found", galleryType)
			}

			svg, err := gallery.GenerateGallery(config)
			if err != nil {
				t.Fatalf("Failed to generate gallery: %v", err)
			}

			// Basic SVG validation
			if svg == "" {
				t.Error("Generated SVG is empty")
			}

			if len(svg) < 100 {
				t.Error("Generated SVG is suspiciously short")
			}

			// Should start with SVG tag
			if len(svg) < 5 || svg[:5] != "<svg " {
				t.Error("Generated output doesn't start with <svg tag")
			}

			// Should end with closing SVG tag
			if len(svg) < 6 || svg[len(svg)-6:] != "</svg>" {
				t.Error("Generated output doesn't end with </svg>")
			}
		})
	}
}
