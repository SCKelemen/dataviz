package layout

import (
	"fmt"
	"sort"
	"strings"

	"github.com/SCKelemen/dataviz/transforms"
	"github.com/SCKelemen/svg"
	"github.com/SCKelemen/units"
)

// FacetLayout defines how facets are arranged
type FacetLayout string

const (
	FacetWrap   FacetLayout = "wrap"   // Wrap into rows (specify columns)
	FacetLayoutGrid FacetLayout = "grid"   // Fixed grid (specify rows and cols)
	FacetCustom FacetLayout = "custom" // Custom positioning
)

// ScaleSharing defines how scales are shared across facets
type ScaleSharing string

const (
	ScaleShareNone ScaleSharing = "none" // Independent scales per facet
	ScaleShareX    ScaleSharing = "x"    // Share X scale only
	ScaleShareY    ScaleSharing = "y"    // Share Y scale only
	ScaleShareXY   ScaleSharing = "xy"   // Share both X and Y scales
)

// Facet represents a specification for creating small multiples
type Facet struct {
	// Field to facet by (DataPoint.Group or other field)
	Field string

	// How to arrange facets
	Layout FacetLayout

	// Number of columns (for wrap layout)
	NCols int

	// Number of rows (for grid layout)
	NRows int

	// Scale sharing strategy
	ScaleSharing ScaleSharing

	// Whether to show titles for each facet
	ShowTitles bool

	// Custom ordering for facet values
	Order []string

	// Gap between facets
	Gap units.Length

	// Margin around each facet
	FacetMargin Margin
}

// NewFacet creates a new facet specification
func NewFacet(field string) *Facet {
	return &Facet{
		Field:        field,
		Layout:       FacetWrap,
		NCols:        2,
		NRows:        0,
		ScaleSharing: ScaleShareNone,
		ShowTitles:   true,
		Gap:          units.Px(10),
		FacetMargin:  Uniform(units.Px(5)),
	}
}

// WithLayout sets the facet layout
func (f *Facet) WithLayout(layout FacetLayout) *Facet {
	f.Layout = layout
	return f
}

// WithCols sets the number of columns (wrap layout)
func (f *Facet) WithCols(cols int) *Facet {
	f.NCols = cols
	return f
}

// WithRows sets the number of rows (grid layout)
func (f *Facet) WithRows(rows int) *Facet {
	f.NRows = rows
	return f
}

// WithScaleSharing sets the scale sharing strategy
func (f *Facet) WithScaleSharing(sharing ScaleSharing) *Facet {
	f.ScaleSharing = sharing
	return f
}

// WithTitles sets whether to show titles
func (f *Facet) WithTitles(show bool) *Facet {
	f.ShowTitles = show
	return f
}

// WithOrder sets custom ordering for facet values
func (f *Facet) WithOrder(order []string) *Facet {
	f.Order = order
	return f
}

// WithGap sets the gap between facets
func (f *Facet) WithGap(gap units.Length) *Facet {
	f.Gap = gap
	return f
}

// WithFacetMargin sets the margin around each facet
func (f *Facet) WithFacetMargin(margin Margin) *Facet {
	f.FacetMargin = margin
	return f
}

// FacetData holds data split by facet values
type FacetData struct {
	Value string              // Facet category value
	Data  []transforms.DataPoint // Data points for this facet
	Index int                 // Position in facet order
}

// Split splits data into facets based on the field
func (f *Facet) Split(data []transforms.DataPoint) []FacetData {
	if len(data) == 0 {
		return nil
	}

	// Group data by facet field value
	groups := make(map[string][]transforms.DataPoint)

	for _, d := range data {
		var key string
		switch f.Field {
		case "Group", "group":
			key = d.Group
		case "Label", "label":
			key = d.Label
		default:
			// For custom fields, try to extract from Data interface
			if d.Data != nil {
				if m, ok := d.Data.(map[string]interface{}); ok {
					if v, exists := m[f.Field]; exists {
						key = fmt.Sprintf("%v", v)
					}
				}
			}
			// Fallback to Group if no key found
			if key == "" {
				key = d.Group
			}
		}

		if key == "" {
			key = "default"
		}

		groups[key] = append(groups[key], d)
	}

	// Convert to FacetData slice
	facets := make([]FacetData, 0, len(groups))
	for value, data := range groups {
		facets = append(facets, FacetData{
			Value: value,
			Data:  data,
		})
	}

	// Apply ordering
	if len(f.Order) > 0 {
		// Custom order specified
		ordered := make([]FacetData, 0, len(facets))
		orderMap := make(map[string]int)
		for i, v := range f.Order {
			orderMap[v] = i
		}

		// Add ordered facets
		for _, value := range f.Order {
			for _, fd := range facets {
				if fd.Value == value {
					fd.Index = len(ordered)
					ordered = append(ordered, fd)
					break
				}
			}
		}

		// Add any remaining facets not in order
		for _, fd := range facets {
			found := false
			for _, o := range ordered {
				if o.Value == fd.Value {
					found = true
					break
				}
			}
			if !found {
				fd.Index = len(ordered)
				ordered = append(ordered, fd)
			}
		}

		facets = ordered
	} else {
		// Sort alphabetically
		sort.Slice(facets, func(i, j int) bool {
			return facets[i].Value < facets[j].Value
		})
		for i := range facets {
			facets[i].Index = i
		}
	}

	return facets
}

// CalculateDimensions calculates grid dimensions based on layout and number of facets
func (f *Facet) CalculateDimensions(numFacets int) (rows, cols int) {
	switch f.Layout {
	case FacetLayoutGrid:
		if f.NRows > 0 && f.NCols > 0 {
			return f.NRows, f.NCols
		}
		if f.NRows > 0 {
			cols = (numFacets + f.NRows - 1) / f.NRows
			return f.NRows, cols
		}
		if f.NCols > 0 {
			rows = (numFacets + f.NCols - 1) / f.NCols
			return rows, f.NCols
		}
		// Default: try to make square
		cols = int(float64(numFacets) + 0.5)
		rows = (numFacets + cols - 1) / cols
		return rows, cols

	case FacetWrap:
		fallthrough
	default:
		cols = f.NCols
		if cols <= 0 {
			cols = 2
		}
		rows = (numFacets + cols - 1) / cols
		return rows, cols
	}
}

// FacetRenderer renders multiple charts in a faceted layout
type FacetRenderer struct {
	// Facet specification
	Facet *Facet

	// Grid layout
	Grid *FacetGrid

	// Chart renderer function
	// Takes data and bounds, returns SVG content
	ChartRenderer func(data []transforms.DataPoint, bounds Rect) string

	// Title renderer function (optional)
	TitleRenderer func(value string, bounds Rect) string
}

// NewFacetRenderer creates a new facet renderer
func NewFacetRenderer(facet *Facet, width, height units.Length) *FacetRenderer {
	// Calculate dimensions
	rows, cols := facet.CalculateDimensions(1) // Will recalculate with actual data

	// Create grid
	grid := NewFacetGrid(width, height, rows, cols)
	grid.SetGap(facet.Gap)
	grid.SetFacetMargin(facet.FacetMargin)
	grid.SetShowTitles(facet.ShowTitles)
	grid.SetScaleSharing(facet.ScaleSharing)

	return &FacetRenderer{
		Facet: facet,
		Grid:  grid,
	}
}

// Render renders the faceted plot
func (fr *FacetRenderer) Render(data []transforms.DataPoint) string {
	if fr.ChartRenderer == nil {
		return ""
	}

	// Split data into facets
	facets := fr.Facet.Split(data)
	if len(facets) == 0 {
		return ""
	}

	// Recalculate grid dimensions
	rows, cols := fr.Facet.CalculateDimensions(len(facets))
	fr.Grid.rows = rows
	fr.Grid.cols = cols

	// Create SVG
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg viewBox="0 0 %f %f" xmlns="http://www.w3.org/2000/svg">`,
		fr.Grid.bounds.Width.Value, fr.Grid.bounds.Height.Value))
	sb.WriteString("\n")

	// Render each facet
	for _, facet := range facets {
		if facet.Index >= rows*cols {
			break // Grid is full
		}

		row := facet.Index / cols
		col := facet.Index % cols

		// Get cell bounds
		cellBounds := fr.Grid.FacetCell(row, col)

		// Render title if enabled
		if fr.Facet.ShowTitles {
			titleBounds := fr.Grid.FacetTitleArea(row, col)
			if fr.TitleRenderer != nil {
				sb.WriteString(fr.TitleRenderer(facet.Value, titleBounds))
			} else {
				// Default title rendering
				sb.WriteString("  ")
				sb.WriteString(svg.Text(
					facet.Value,
					titleBounds.X.Value+titleBounds.Width.Value/2,
					titleBounds.Y.Value+titleBounds.Height.Value/2,
					svg.Style{
						TextAnchor: "middle",
						FontSize:   units.Px(12),
						FontWeight: "bold",
					},
				))
				sb.WriteString("\n")
			}
		}

		// Render chart
		chartSVG := fr.ChartRenderer(facet.Data, cellBounds)
		sb.WriteString(chartSVG)
	}

	sb.WriteString("</svg>\n")
	return sb.String()
}

// FacetedChart is a convenience wrapper for creating faceted plots
type FacetedChart struct {
	Data          []transforms.DataPoint
	Facet         *Facet
	Width         units.Length
	Height        units.Length
	ChartRenderer func(data []transforms.DataPoint, bounds Rect) string
}

// Render renders the faceted chart
func (fc *FacetedChart) Render() string {
	renderer := NewFacetRenderer(fc.Facet, fc.Width, fc.Height)
	renderer.ChartRenderer = fc.ChartRenderer
	return renderer.Render(fc.Data)
}

// ComputeSharedDomain calculates shared domain across all facets
func ComputeSharedDomain(facets []FacetData, field string) (min, max float64) {
	first := true

	for _, facet := range facets {
		for _, d := range facet.Data {
			var value float64
			switch field {
			case "Y":
				value = d.Y
			case "X":
				if f, ok := d.X.(float64); ok {
					value = f
				} else {
					continue
				}
			default:
				value = d.Y
			}

			if first {
				min = value
				max = value
				first = false
			} else {
				if value < min {
					min = value
				}
				if value > max {
					max = value
				}
			}
		}
	}

	return
}

// GetScaleDomain returns the domain for a scale based on sharing strategy
func GetScaleDomain(facets []FacetData, facetIndex int, field string, sharing ScaleSharing) (min, max float64) {
	switch sharing {
	case ScaleShareNone:
		// Independent scales - use only this facet's data
		if facetIndex >= 0 && facetIndex < len(facets) {
			first := true
			for _, d := range facets[facetIndex].Data {
				var value float64
				switch field {
				case "Y":
					value = d.Y
				case "X":
					if f, ok := d.X.(float64); ok {
						value = f
					} else {
						continue
					}
				default:
					value = d.Y
				}

				if first {
					min = value
					max = value
					first = false
				} else {
					if value < min {
						min = value
					}
					if value > max {
						max = value
					}
				}
			}
		}
		return

	case ScaleShareX:
		if field == "X" {
			return ComputeSharedDomain(facets, field)
		}
		return GetScaleDomain(facets, facetIndex, field, ScaleShareNone)

	case ScaleShareY:
		if field == "Y" {
			return ComputeSharedDomain(facets, field)
		}
		return GetScaleDomain(facets, facetIndex, field, ScaleShareNone)

	case ScaleShareXY:
		return ComputeSharedDomain(facets, field)

	default:
		return ComputeSharedDomain(facets, field)
	}
}
