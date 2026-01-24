package layout

import (
	"fmt"
	"sort"
	"strings"

	"github.com/SCKelemen/dataviz/transforms"
	"github.com/SCKelemen/layout"
	"github.com/SCKelemen/svg"
	"github.com/SCKelemen/units"
)

// Facet represents a specification for creating small multiples
type FacetSpec struct {
	// Field to facet by (DataPoint.Group or other field)
	Field string

	// Number of columns
	NCols int

	// Number of rows (0 = auto)
	NRows int

	// Scale sharing strategy
	ScaleSharing ScaleSharing

	// Whether to show titles for each facet
	ShowTitles bool

	// Custom ordering for facet values
	Order []string

	// Gap between facets
	Gap float64

	// Margin around each facet plot area
	FacetMargin float64
}

// ScaleSharing defines how scales are shared across facets
type ScaleSharing string

const (
	ScaleShareNone ScaleSharing = "none" // Independent scales per facet
	ScaleShareX    ScaleSharing = "x"    // Share X scale only
	ScaleShareY    ScaleSharing = "y"    // Share Y scale only
	ScaleShareXY   ScaleSharing = "xy"   // Share both X and Y scales
)

// NewFacetSpec creates a new facet specification
func NewFacetSpec(field string) *FacetSpec {
	return &FacetSpec{
		Field:        field,
		NCols:        2,
		NRows:        0,
		ScaleSharing: ScaleShareNone,
		ShowTitles:   true,
		Gap:          10,
		FacetMargin:  5,
	}
}

// WithCols sets the number of columns
func (f *FacetSpec) WithCols(cols int) *FacetSpec {
	f.NCols = cols
	return f
}

// WithRows sets the number of rows
func (f *FacetSpec) WithRows(rows int) *FacetSpec {
	f.NRows = rows
	return f
}

// WithScaleSharing sets the scale sharing strategy
func (f *FacetSpec) WithScaleSharing(sharing ScaleSharing) *FacetSpec {
	f.ScaleSharing = sharing
	return f
}

// WithTitles sets whether to show titles
func (f *FacetSpec) WithTitles(show bool) *FacetSpec {
	f.ShowTitles = show
	return f
}

// WithOrder sets custom ordering for facet values
func (f *FacetSpec) WithOrder(order []string) *FacetSpec {
	f.Order = order
	return f
}

// WithGap sets the gap between facets
func (f *FacetSpec) WithGap(gap float64) *FacetSpec {
	f.Gap = gap
	return f
}

// WithFacetMargin sets the margin around each facet
func (f *FacetSpec) WithFacetMargin(margin float64) *FacetSpec {
	f.FacetMargin = margin
	return f
}

// FacetData holds data split by facet values
type FacetData struct {
	Value string                     // Facet category value
	Data  []transforms.DataPoint     // Data points for this facet
	Index int                        // Position in facet order
	Node  *layout.Node               // Layout node for this facet
}

// Split splits data into facets based on the field
func (f *FacetSpec) Split(data []transforms.DataPoint) []FacetData {
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

// CalculateDimensions calculates grid dimensions based on number of facets
func (f *FacetSpec) CalculateDimensions(numFacets int) (rows, cols int) {
	cols = f.NCols
	if cols <= 0 {
		cols = 2
	}

	rows = f.NRows
	if rows <= 0 {
		rows = (numFacets + cols - 1) / cols
	}

	return rows, cols
}

// BuildLayout creates a CSS Grid layout for facets
func (f *FacetSpec) BuildLayout(numFacets int, width, height float64) *layout.Node {
	rows, cols := f.CalculateDimensions(numFacets)

	// Create CSS Grid
	root := ChartGridWithGap(rows, cols, f.Gap)
	root.Style.Width = layout.Px(width)
	root.Style.Height = layout.Px(height)

	return root
}

// FacetPlot creates a faceted visualization
type FacetPlot struct {
	Spec   *FacetSpec
	Width  float64
	Height float64
	Data   []transforms.DataPoint

	// Renderer for each facet cell
	CellRenderer func(data []transforms.DataPoint, bounds layout.Rect) string

	// Optional title renderer
	TitleRenderer func(value string, bounds layout.Rect) string
}

// NewFacetPlot creates a new faceted plot
func NewFacetPlot(spec *FacetSpec, width, height float64) *FacetPlot {
	return &FacetPlot{
		Spec:   spec,
		Width:  width,
		Height: height,
	}
}

// WithData sets the data
func (fp *FacetPlot) WithData(data []transforms.DataPoint) *FacetPlot {
	fp.Data = data
	return fp
}

// WithCellRenderer sets the cell renderer
func (fp *FacetPlot) WithCellRenderer(renderer func(data []transforms.DataPoint, bounds layout.Rect) string) *FacetPlot {
	fp.CellRenderer = renderer
	return fp
}

// WithTitleRenderer sets the title renderer
func (fp *FacetPlot) WithTitleRenderer(renderer func(value string, bounds layout.Rect) string) *FacetPlot {
	fp.TitleRenderer = renderer
	return fp
}

// Render renders the faceted plot
func (fp *FacetPlot) Render() string {
	if fp.CellRenderer == nil {
		return ""
	}

	// Split data into facets
	facets := fp.Spec.Split(fp.Data)
	if len(facets) == 0 {
		return ""
	}

	// Build grid layout
	root := fp.Spec.BuildLayout(len(facets), fp.Width, fp.Height)

	// Create nodes for each facet
	for i := range facets {
		facetNode := &layout.Node{
			Style: layout.Style{
				Display: layout.DisplayBlock,
			},
		}

		// Add padding if specified
		if fp.Spec.FacetMargin > 0 {
			facetNode = WithPadding(facetNode, fp.Spec.FacetMargin)
		}

		facets[i].Node = facetNode
		root.AddChild(facetNode)
	}

	// Compute layout
	constraints := layout.Loose(fp.Width, fp.Height)
	ctx := layout.NewLayoutContext(fp.Width, fp.Height, 16)
	layout.Layout(root, constraints, ctx)

	// Render to SVG
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg viewBox="0 0 %f %f" xmlns="http://www.w3.org/2000/svg">`,
		fp.Width, fp.Height))
	sb.WriteString("\n")

	// Render each facet
	for _, facet := range facets {
		// Create group with transform
		sb.WriteString(fmt.Sprintf(`<g transform="translate(%f,%f)">`,
			facet.Node.Rect.X, facet.Node.Rect.Y))
		sb.WriteString("\n")

		// Render title if enabled
		if fp.Spec.ShowTitles {
			titleHeight := 20.0
			titleBounds := layout.Rect{
				X:      0,
				Y:      0,
				Width:  facet.Node.Rect.Width,
				Height: titleHeight,
			}

			if fp.TitleRenderer != nil {
				sb.WriteString(fp.TitleRenderer(facet.Value, titleBounds))
			} else {
				// Default title rendering
				sb.WriteString("  ")
				sb.WriteString(svg.Text(
					facet.Value,
					facet.Node.Rect.Width/2,
					titleHeight/2,
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
		chartBounds := layout.Rect{
			X:      0,
			Y:      0,
			Width:  facet.Node.Rect.Width,
			Height: facet.Node.Rect.Height,
		}

		if fp.Spec.ShowTitles {
			chartBounds.Y = 20
			chartBounds.Height -= 20
		}

		chartSVG := fp.CellRenderer(facet.Data, chartBounds)
		sb.WriteString(chartSVG)

		sb.WriteString("</g>\n")
	}

	sb.WriteString("</svg>\n")
	return sb.String()
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
