package annotations

import (
	"github.com/SCKelemen/dataviz/scales"
	"github.com/SCKelemen/svg"
	"github.com/SCKelemen/units"
)

// Annotation represents a visual annotation on a chart
type Annotation interface {
	// Render renders the annotation to SVG
	Render(xScale, yScale scales.Scale) string
}

// Position represents how to position an annotation
type Position string

const (
	// PositionData positions relative to data coordinates (uses scales)
	PositionData Position = "data"

	// PositionPixel positions using absolute pixel coordinates
	PositionPixel Position = "pixel"

	// PositionRelative positions as percentage (0-1) of plot area
	PositionRelative Position = "relative"
)

// Anchor defines where text is anchored
type Anchor string

const (
	AnchorStart  Anchor = "start"
	AnchorMiddle Anchor = "middle"
	AnchorEnd    Anchor = "end"
)

// Baseline defines vertical text alignment
type Baseline string

const (
	BaselineTop        Baseline = "hanging"
	BaselineMiddle     Baseline = "middle"
	BaselineBottom     Baseline = "baseline"
	BaselineAlphabetic Baseline = "alphabetic"
)

// AnnotationStyle contains common styling for annotations
type AnnotationStyle struct {
	Stroke       string
	StrokeWidth  float64
	Fill         string
	Opacity      float64
	StrokeDash   string
	FontSize     units.Length
	FontFamily   string
	FontWeight   string
	TextAnchor   Anchor
	TextBaseline Baseline
	Color        string
}

// DefaultAnnotationStyle returns the default annotation style
func DefaultAnnotationStyle() AnnotationStyle {
	return AnnotationStyle{
		Stroke:       "#666",
		StrokeWidth:  1,
		Fill:         "#666",
		Opacity:      0.8,
		StrokeDash:   "",
		FontSize:     units.Px(12),
		FontFamily:   "sans-serif",
		FontWeight:   "normal",
		TextAnchor:   AnchorMiddle,
		TextBaseline: BaselineMiddle,
		Color:        "#333",
	}
}

// toSVGStyle converts AnnotationStyle to svg.Style
func (as AnnotationStyle) toSVGStyle() svg.Style {
	return svg.Style{
		Stroke:       as.Stroke,
		StrokeWidth:  as.StrokeWidth,
		Fill:         as.Fill,
		Opacity:      as.Opacity,
		FontSize:     as.FontSize,
		FontFamily:   as.FontFamily,
		FontWeight:   svg.FontWeight(as.FontWeight),
		TextAnchor:   svg.TextAnchor(as.TextAnchor),
	}
}

// AnnotationLayer groups multiple annotations
type AnnotationLayer struct {
	Annotations []Annotation
}

// NewAnnotationLayer creates a new annotation layer
func NewAnnotationLayer() *AnnotationLayer {
	return &AnnotationLayer{
		Annotations: make([]Annotation, 0),
	}
}

// Add adds an annotation to the layer
func (al *AnnotationLayer) Add(annotation Annotation) *AnnotationLayer {
	al.Annotations = append(al.Annotations, annotation)
	return al
}

// Render renders all annotations in the layer
func (al *AnnotationLayer) Render(xScale, yScale scales.Scale) string {
	var result string
	for _, annotation := range al.Annotations {
		result += annotation.Render(xScale, yScale)
	}
	return result
}
