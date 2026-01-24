package charts

import (
	"fmt"
	"math"
	"sort"

	"github.com/SCKelemen/svg"
	"github.com/SCKelemen/units"
)

// WordCloudWord represents a word in the word cloud
type WordCloudWord struct {
	Text       string
	Frequency  float64 // Frequency/weight of the word
	Color      string  // Optional custom color
	Angle      float64 // Optional rotation angle in degrees (0 = horizontal)
}

// WordCloudSpec configures wordcloud rendering
type WordCloudSpec struct {
	Words        []WordCloudWord
	Width        float64
	Height       float64
	MinFontSize  float64   // Minimum font size (default: 12)
	MaxFontSize  float64   // Maximum font size (default: 72)
	FontFamily   string    // Font family (default: sans-serif)
	DefaultColor string    // Default word color
	Layout       string    // "spiral", "horizontal" (default: spiral)
	Title        string
}

// RenderWordCloud generates an SVG word cloud
func RenderWordCloud(spec WordCloudSpec) string {
	if len(spec.Words) == 0 {
		return ""
	}

	// Set defaults
	if spec.MinFontSize == 0 {
		spec.MinFontSize = 12
	}
	if spec.MaxFontSize == 0 {
		spec.MaxFontSize = 72
	}
	if spec.FontFamily == "" {
		spec.FontFamily = "sans-serif"
	}
	if spec.DefaultColor == "" {
		spec.DefaultColor = "#3b82f6"
	}
	if spec.Layout == "" {
		spec.Layout = "spiral"
	}

	// Find min/max frequency
	minFreq := spec.Words[0].Frequency
	maxFreq := spec.Words[0].Frequency
	for _, word := range spec.Words {
		if word.Frequency < minFreq {
			minFreq = word.Frequency
		}
		if word.Frequency > maxFreq {
			maxFreq = word.Frequency
		}
	}

	freqRange := maxFreq - minFreq
	if freqRange == 0 {
		freqRange = 1
	}

	// Sort words by frequency (descending) for better placement
	sortedWords := make([]WordCloudWord, len(spec.Words))
	copy(sortedWords, spec.Words)
	sort.Slice(sortedWords, func(i, j int) bool {
		return sortedWords[i].Frequency > sortedWords[j].Frequency
	})

	// Calculate font sizes
	words := make([]wordWithPosition, len(sortedWords))
	for i, word := range sortedWords {
		// Scale font size based on frequency
		normalizedFreq := (word.Frequency - minFreq) / freqRange
		fontSize := spec.MinFontSize + normalizedFreq*(spec.MaxFontSize-spec.MinFontSize)

		words[i] = wordWithPosition{
			word:     word,
			fontSize: fontSize,
		}
	}

	// Layout words
	switch spec.Layout {
	case "horizontal":
		layoutHorizontal(words, spec.Width, spec.Height)
	default: // "spiral"
		layoutSpiral(words, spec.Width, spec.Height)
	}

	var result string

	// Draw title
	if spec.Title != "" {
		titleStyle := svg.Style{
			FontSize:         units.Px(16),
			FontFamily:       "sans-serif",
			FontWeight:       "bold",
			TextAnchor:       svg.TextAnchorMiddle,
			DominantBaseline: svg.DominantBaselineHanging,
		}
		result += svg.Text(spec.Title, spec.Width/2, 10, titleStyle) + "\n"
	}

	// Default colors for variety
	defaultColors := []string{"#3b82f6", "#10b981", "#f59e0b", "#ef4444", "#8b5cf6", "#ec4899", "#06b6d4", "#84cc16"}

	// Draw words
	for i, word := range words {
		// Get word color
		wordColor := word.word.Color
		if wordColor == "" {
			// Use varied colors for visual interest
			if spec.DefaultColor != "" {
				wordColor = spec.DefaultColor
			} else {
				wordColor = defaultColors[i%len(defaultColors)]
			}
		}

		// Create text style
		textStyle := svg.Style{
			Fill:             wordColor,
			FontSize:         units.Px(word.fontSize),
			FontFamily:       spec.FontFamily,
			FontWeight:       "bold",
			TextAnchor:       svg.TextAnchorMiddle,
			DominantBaseline: svg.DominantBaselineMiddle,
		}

		// Apply rotation if specified
		if word.word.Angle != 0 {
			result += fmt.Sprintf(`<text x="%.2f" y="%.2f" font-size="%.0f" font-family="%s" font-weight="bold" text-anchor="middle" dominant-baseline="middle" fill="%s" transform="rotate(%.1f %.2f %.2f)">%s</text>`,
				word.x, word.y, word.fontSize, spec.FontFamily, wordColor, word.word.Angle, word.x, word.y, word.word.Text) + "\n"
		} else {
			result += svg.Text(word.word.Text, word.x, word.y, textStyle) + "\n"
		}
	}

	return result
}

// wordWithPosition stores word data with layout information
type wordWithPosition struct {
	word     WordCloudWord
	fontSize float64
	x        float64
	y        float64
}

// layoutSpiral places words in a spiral pattern starting from center
func layoutSpiral(words []wordWithPosition, width, height float64) {
	centerX := width / 2
	centerY := height / 2

	// Simple spiral placement - place words along a spiral
	angle := 0.0
	radius := 0.0
	radiusStep := 5.0
	angleStep := 0.3

	for i := range words {
		// Calculate position on spiral
		x := centerX + radius*math.Cos(angle)
		y := centerY + radius*math.Sin(angle)

		// Clamp to bounds
		margin := words[i].fontSize
		if x < margin {
			x = margin
		}
		if x > width-margin {
			x = width - margin
		}
		if y < margin {
			y = margin
		}
		if y > height-margin {
			y = height - margin
		}

		words[i].x = x
		words[i].y = y

		// Update spiral parameters
		angle += angleStep
		radius += radiusStep

		// Larger words increase the radius more
		if words[i].fontSize > 40 {
			radius += radiusStep * 2
		}
	}
}

// layoutHorizontal places words in rows from top to bottom
func layoutHorizontal(words []wordWithPosition, width, height float64) {
	margin := 20.0
	x := margin
	y := margin + 30.0 // Start below title area

	maxRowHeight := 0.0

	for i := range words {
		// Estimate word width (rough approximation)
		wordWidth := float64(len(words[i].word.Text)) * words[i].fontSize * 0.6

		// Check if word fits on current line
		if x+wordWidth > width-margin && x > margin {
			// Move to next line
			x = margin
			y += maxRowHeight + 10
			maxRowHeight = 0
		}

		// Check if we've run out of vertical space
		if y > height-margin {
			y = height - margin
		}

		words[i].x = x + wordWidth/2
		words[i].y = y + words[i].fontSize/2

		x += wordWidth + 15

		if words[i].fontSize > maxRowHeight {
			maxRowHeight = words[i].fontSize
		}
	}
}

// WordCloudFromFrequencies creates a word cloud from word-frequency pairs
func WordCloudFromFrequencies(words []string, frequencies []float64, width, height float64) string {
	if len(words) != len(frequencies) {
		return ""
	}

	cloudWords := make([]WordCloudWord, len(words))
	for i := range words {
		cloudWords[i] = WordCloudWord{
			Text:      words[i],
			Frequency: frequencies[i],
		}
	}

	spec := WordCloudSpec{
		Words:  cloudWords,
		Width:  width,
		Height: height,
		Layout: "spiral",
	}

	return RenderWordCloud(spec)
}

// RotatedWordCloud creates a word cloud with some words rotated
func RotatedWordCloud(words []string, frequencies []float64, width, height float64) string {
	if len(words) != len(frequencies) {
		return ""
	}

	cloudWords := make([]WordCloudWord, len(words))
	for i := range words {
		// Randomly rotate some words (every 3rd word at 90 degrees)
		angle := 0.0
		if i%3 == 0 {
			angle = 90.0
		}

		cloudWords[i] = WordCloudWord{
			Text:      words[i],
			Frequency: frequencies[i],
			Angle:     angle,
		}
	}

	spec := WordCloudSpec{
		Words:  cloudWords,
		Width:  width,
		Height: height,
		Layout: "spiral",
	}

	return RenderWordCloud(spec)
}

// ColorfulWordCloud creates a word cloud with varied colors
func ColorfulWordCloud(words []string, frequencies []float64, colors []string, width, height float64) string {
	if len(words) != len(frequencies) {
		return ""
	}

	cloudWords := make([]WordCloudWord, len(words))
	for i := range words {
		color := ""
		if i < len(colors) {
			color = colors[i]
		}

		cloudWords[i] = WordCloudWord{
			Text:      words[i],
			Frequency: frequencies[i],
			Color:     color,
		}
	}

	spec := WordCloudSpec{
		Words:  cloudWords,
		Width:  width,
		Height: height,
		Layout: "spiral",
	}

	return RenderWordCloud(spec)
}
