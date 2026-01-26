package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/SCKelemen/dataviz/internal/gallery"
)

func main() {
	if err := generateGalleries(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func generateGalleries() error {
	outputDir := "examples-gallery"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	// Generate all galleries from the registry
	for name, config := range gallery.GalleryRegistry {
		fmt.Printf("Generating %s gallery...\n", name)

		svg, err := gallery.GenerateGallery(config)
		if err != nil {
			fmt.Printf("  ✗ Failed: %v\n", err)
			continue
		}

		outputPath := filepath.Join(outputDir, name+"-gallery.svg")
		if err := os.WriteFile(outputPath, []byte(svg), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", outputPath, err)
		}
		fmt.Printf("  ✓ %s\n", outputPath)
	}

	fmt.Println("✓ Gallery generation complete!")
	return nil
}
