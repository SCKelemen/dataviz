package axes

import (
	"fmt"
	"strings"

	"github.com/SCKelemen/svg"
	"github.com/SCKelemen/units"
)

// AxisStyle contains styling options for axis rendering
type AxisStyle struct {
	StrokeColor      string
	StrokeWidth      float64
	TextColor        string
	FontSize         float64
	FontFamily       string
	GridStrokeColor  string
	GridStrokeWidth  float64
	GridDashArray    string
	TitleFontSize    float64
	TitleFontWeight  string
}

// DefaultAxisStyle returns the default axis styling
func DefaultAxisStyle() AxisStyle {
	return AxisStyle{
		StrokeColor:      "#000000",
		StrokeWidth:      1,
		TextColor:        "#000000",
		FontSize:         11,
		FontFamily:       "sans-serif",
		GridStrokeColor:  "#e0e0e0",
		GridStrokeWidth:  1,
		GridDashArray:    "",
		TitleFontSize:    12,
		TitleFontWeight:  "bold",
	}
}

// RenderOptions contains options for rendering an axis
type RenderOptions struct {
	Style    AxisStyle
	Position units.Length // Position perpendicular to axis (e.g., y-position for horizontal axis)
}

// DefaultRenderOptions returns default render options
func DefaultRenderOptions() RenderOptions {
	return RenderOptions{
		Style:    DefaultAxisStyle(),
		Position: units.Px(0),
	}
}

// Render generates SVG markup for this axis
func (a *Axis) Render(opts RenderOptions) string {
	var sb strings.Builder

	ticks := a.Ticks()
	if len(ticks) == 0 {
		return ""
	}

	// Get axis range
	scaleRange := a.scale.Range()
	rangeStart := scaleRange[0].Value
	rangeEnd := scaleRange[1].Value

	// Open group
	sb.WriteString(fmt.Sprintf(`<g class="axis axis-%s">`, a.orientation.String()))
	sb.WriteString("\n")

	// Render based on orientation
	switch a.orientation {
	case AxisOrientationBottom:
		a.renderHorizontalBottom(&sb, ticks, rangeStart, rangeEnd, opts)
	case AxisOrientationTop:
		a.renderHorizontalTop(&sb, ticks, rangeStart, rangeEnd, opts)
	case AxisOrientationLeft:
		a.renderVerticalLeft(&sb, ticks, rangeStart, rangeEnd, opts)
	case AxisOrientationRight:
		a.renderVerticalRight(&sb, ticks, rangeStart, rangeEnd, opts)
	}

	// Close group
	sb.WriteString("</g>\n")

	return sb.String()
}

// renderHorizontalBottom renders a bottom-oriented horizontal axis
func (a *Axis) renderHorizontalBottom(sb *strings.Builder, ticks []Tick, rangeStart, rangeEnd float64, opts RenderOptions) {
	y := opts.Position.Value

	// Main axis line
	style := svg.Style{
		Stroke:      opts.Style.StrokeColor,
		StrokeWidth: opts.Style.StrokeWidth,
	}
	sb.WriteString("  ")
	sb.WriteString(svg.Line(rangeStart, y, rangeEnd, y, style))
	sb.WriteString("\n")

	// Ticks and labels
	for _, tick := range ticks {
		x := tick.Position.Value

		// Tick mark
		sb.WriteString("  ")
		sb.WriteString(svg.Line(x, y, x, y+a.tickSize.Value, style))
		sb.WriteString("\n")

		// Grid line
		if a.showGrid {
			gridStyle := svg.Style{
				Stroke:      opts.Style.GridStrokeColor,
				StrokeWidth: opts.Style.GridStrokeWidth,
			}
			sb.WriteString("  ")
			sb.WriteString(svg.Line(x, y, x, y-a.gridLength.Value, gridStyle))
			sb.WriteString("\n")
		}

		// Label
		labelY := y + a.tickSize.Value + a.tickPadding.Value + opts.Style.FontSize
		textStyle := svg.Style{
			Fill:       opts.Style.TextColor,
			FontSize:   units.Px(opts.Style.FontSize),
			FontFamily: opts.Style.FontFamily,
			TextAnchor: svg.TextAnchorMiddle,
		}
		sb.WriteString("  ")
		sb.WriteString(svg.Text(tick.Label, x, labelY, textStyle))
		sb.WriteString("\n")
	}

	// Title
	if a.title != "" {
		titleY := y + a.tickSize.Value + a.tickPadding.Value + opts.Style.FontSize + opts.Style.TitleFontSize + 5
		titleX := (rangeStart + rangeEnd) / 2

		titleStyle := svg.Style{
			Fill:       opts.Style.TextColor,
			FontSize:   units.Px(opts.Style.TitleFontSize),
			FontFamily: opts.Style.FontFamily,
			FontWeight: svg.FontWeight(opts.Style.TitleFontWeight),
			TextAnchor: svg.TextAnchorMiddle,
		}
		sb.WriteString("  ")
		sb.WriteString(svg.Text(a.title, titleX, titleY, titleStyle))
		sb.WriteString("\n")
	}
}

// renderHorizontalTop renders a top-oriented horizontal axis
func (a *Axis) renderHorizontalTop(sb *strings.Builder, ticks []Tick, rangeStart, rangeEnd float64, opts RenderOptions) {
	y := opts.Position.Value

	// Main axis line
	style := svg.Style{
		Stroke:      opts.Style.StrokeColor,
		StrokeWidth: opts.Style.StrokeWidth,
	}
	sb.WriteString("  ")
	sb.WriteString(svg.Line(rangeStart, y, rangeEnd, y, style))
	sb.WriteString("\n")

	// Ticks and labels
	for _, tick := range ticks {
		x := tick.Position.Value

		// Tick mark (upward)
		sb.WriteString("  ")
		sb.WriteString(svg.Line(x, y, x, y-a.tickSize.Value, style))
		sb.WriteString("\n")

		// Grid line
		if a.showGrid {
			gridStyle := svg.Style{
				Stroke:      opts.Style.GridStrokeColor,
				StrokeWidth: opts.Style.GridStrokeWidth,
			}
			sb.WriteString("  ")
			sb.WriteString(svg.Line(x, y, x, y+a.gridLength.Value, gridStyle))
			sb.WriteString("\n")
		}

		// Label (above tick)
		labelY := y - a.tickSize.Value - a.tickPadding.Value
		textStyle := svg.Style{
			Fill:       opts.Style.TextColor,
			FontSize:   units.Px(opts.Style.FontSize),
			FontFamily: opts.Style.FontFamily,
			TextAnchor: svg.TextAnchorMiddle,
		}
		sb.WriteString("  ")
		sb.WriteString(svg.Text(tick.Label, x, labelY, textStyle))
		sb.WriteString("\n")
	}

	// Title
	if a.title != "" {
		titleY := y - a.tickSize.Value - a.tickPadding.Value - opts.Style.TitleFontSize - 5
		titleX := (rangeStart + rangeEnd) / 2

		titleStyle := svg.Style{
			Fill:       opts.Style.TextColor,
			FontSize:   units.Px(opts.Style.TitleFontSize),
			FontFamily: opts.Style.FontFamily,
			FontWeight: svg.FontWeight(opts.Style.TitleFontWeight),
			TextAnchor: svg.TextAnchorMiddle,
		}
		sb.WriteString("  ")
		sb.WriteString(svg.Text(a.title, titleX, titleY, titleStyle))
		sb.WriteString("\n")
	}
}

// renderVerticalLeft renders a left-oriented vertical axis
func (a *Axis) renderVerticalLeft(sb *strings.Builder, ticks []Tick, rangeStart, rangeEnd float64, opts RenderOptions) {
	x := opts.Position.Value

	// Main axis line
	style := svg.Style{
		Stroke:      opts.Style.StrokeColor,
		StrokeWidth: opts.Style.StrokeWidth,
	}
	sb.WriteString("  ")
	sb.WriteString(svg.Line(x, rangeStart, x, rangeEnd, style))
	sb.WriteString("\n")

	// Ticks and labels
	for _, tick := range ticks {
		y := tick.Position.Value

		// Tick mark (leftward)
		sb.WriteString("  ")
		sb.WriteString(svg.Line(x, y, x-a.tickSize.Value, y, style))
		sb.WriteString("\n")

		// Grid line
		if a.showGrid {
			gridStyle := svg.Style{
				Stroke:      opts.Style.GridStrokeColor,
				StrokeWidth: opts.Style.GridStrokeWidth,
			}
			sb.WriteString("  ")
			sb.WriteString(svg.Line(x, y, x+a.gridLength.Value, y, gridStyle))
			sb.WriteString("\n")
		}

		// Label (left of tick)
		labelX := x - a.tickSize.Value - a.tickPadding.Value
		textStyle := svg.Style{
			Fill:              opts.Style.TextColor,
			FontSize:          units.Px(opts.Style.FontSize),
			FontFamily:        opts.Style.FontFamily,
			TextAnchor:        svg.TextAnchorEnd,
			DominantBaseline:  svg.DominantBaselineMiddle,
		}
		sb.WriteString("  ")
		sb.WriteString(svg.Text(tick.Label, labelX, y, textStyle))
		sb.WriteString("\n")
	}

	// Title (rotated)
	if a.title != "" {
		titleX := x - a.tickSize.Value - a.tickPadding.Value - 30 // Approximate label width
		titleY := (rangeStart + rangeEnd) / 2

		titleStyle := svg.Style{
			Fill:       opts.Style.TextColor,
			FontSize:   units.Px(opts.Style.TitleFontSize),
			FontFamily: opts.Style.FontFamily,
			FontWeight: svg.FontWeight(opts.Style.TitleFontWeight),
			TextAnchor: svg.TextAnchorMiddle,
		}
		sb.WriteString("  ")
		sb.WriteString(fmt.Sprintf(`<text x="%g" y="%g" transform="rotate(-90 %g %g)"%s>%s</text>`,
			titleX, titleY, titleX, titleY,
			formatStyleAttrs(titleStyle),
			escapeXML(a.title)))
		sb.WriteString("\n")
	}
}

// renderVerticalRight renders a right-oriented vertical axis
func (a *Axis) renderVerticalRight(sb *strings.Builder, ticks []Tick, rangeStart, rangeEnd float64, opts RenderOptions) {
	x := opts.Position.Value

	// Main axis line
	style := svg.Style{
		Stroke:      opts.Style.StrokeColor,
		StrokeWidth: opts.Style.StrokeWidth,
	}
	sb.WriteString("  ")
	sb.WriteString(svg.Line(x, rangeStart, x, rangeEnd, style))
	sb.WriteString("\n")

	// Ticks and labels
	for _, tick := range ticks {
		y := tick.Position.Value

		// Tick mark (rightward)
		sb.WriteString("  ")
		sb.WriteString(svg.Line(x, y, x+a.tickSize.Value, y, style))
		sb.WriteString("\n")

		// Grid line
		if a.showGrid {
			gridStyle := svg.Style{
				Stroke:      opts.Style.GridStrokeColor,
				StrokeWidth: opts.Style.GridStrokeWidth,
			}
			sb.WriteString("  ")
			sb.WriteString(svg.Line(x, y, x-a.gridLength.Value, y, gridStyle))
			sb.WriteString("\n")
		}

		// Label (right of tick)
		labelX := x + a.tickSize.Value + a.tickPadding.Value
		textStyle := svg.Style{
			Fill:              opts.Style.TextColor,
			FontSize:          units.Px(opts.Style.FontSize),
			FontFamily:        opts.Style.FontFamily,
			TextAnchor:        svg.TextAnchorStart,
			DominantBaseline:  svg.DominantBaselineMiddle,
		}
		sb.WriteString("  ")
		sb.WriteString(svg.Text(tick.Label, labelX, y, textStyle))
		sb.WriteString("\n")
	}

	// Title (rotated)
	if a.title != "" {
		titleX := x + a.tickSize.Value + a.tickPadding.Value + 30 // Approximate label width
		titleY := (rangeStart + rangeEnd) / 2

		titleStyle := svg.Style{
			Fill:       opts.Style.TextColor,
			FontSize:   units.Px(opts.Style.TitleFontSize),
			FontFamily: opts.Style.FontFamily,
			FontWeight: svg.FontWeight(opts.Style.TitleFontWeight),
			TextAnchor: svg.TextAnchorMiddle,
		}
		sb.WriteString("  ")
		sb.WriteString(fmt.Sprintf(`<text x="%g" y="%g" transform="rotate(90 %g %g)"%s>%s</text>`,
			titleX, titleY, titleX, titleY,
			formatStyleAttrs(titleStyle),
			escapeXML(a.title)))
		sb.WriteString("\n")
	}
}

// String generates the complete SVG string for this axis
func (a *Axis) String(opts RenderOptions) string {
	return a.Render(opts)
}

// Helper functions

func formatStyleAttrs(style svg.Style) string {
	var attrs []string

	if style.Fill != "" {
		attrs = append(attrs, fmt.Sprintf(`fill="%s"`, style.Fill))
	}
	if style.FontSize.Value > 0 {
		attrs = append(attrs, fmt.Sprintf(`font-size="%s"`, style.FontSize.String()))
	}
	if style.FontFamily != "" {
		attrs = append(attrs, fmt.Sprintf(`font-family="%s"`, style.FontFamily))
	}
	if string(style.FontWeight) != "" {
		attrs = append(attrs, fmt.Sprintf(`font-weight="%s"`, string(style.FontWeight)))
	}
	if string(style.TextAnchor) != "" {
		attrs = append(attrs, fmt.Sprintf(`text-anchor="%s"`, string(style.TextAnchor)))
	}

	if len(attrs) == 0 {
		return ""
	}
	return " " + strings.Join(attrs, " ")
}

func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}
