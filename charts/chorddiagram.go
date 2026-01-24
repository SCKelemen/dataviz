package charts

import (
	"fmt"
	"math"

	"github.com/SCKelemen/svg"
	"github.com/SCKelemen/units"
)

// ChordEntity represents an entity in the chord diagram
type ChordEntity struct {
	ID    string
	Label string
	Color string // Optional custom color
}

// ChordRelation represents a relationship between two entities
type ChordRelation struct {
	Source string  // Source entity ID
	Target string  // Target entity ID
	Value  float64 // Relationship strength (determines chord width)
}

// ChordDiagramSpec configures chord diagram rendering
type ChordDiagramSpec struct {
	Entities      []ChordEntity
	Relations     []ChordRelation
	Width         float64
	Height        float64
	InnerRadius   float64 // Inner radius for entity arcs (auto if 0)
	ArcWidth      float64 // Width of entity arcs (default: 20)
	ArcPadding    float64 // Padding between entity arcs in degrees (default: 2)
	DefaultColor  string  // Default entity color
	ShowLabels    bool    // Show entity labels
	Title         string
}

// RenderChordDiagram generates an SVG chord diagram
func RenderChordDiagram(spec ChordDiagramSpec) string {
	if len(spec.Entities) == 0 || len(spec.Relations) == 0 {
		return ""
	}

	// Set defaults
	if spec.ArcWidth == 0 {
		spec.ArcWidth = 20
	}
	if spec.ArcPadding == 0 {
		spec.ArcPadding = 2
	}
	if spec.DefaultColor == "" {
		spec.DefaultColor = "#3b82f6"
	}

	// Calculate center and radius
	centerX := spec.Width / 2
	centerY := spec.Height / 2
	margin := 80.0
	maxRadius := math.Min(spec.Width, spec.Height)/2 - margin

	// Set inner radius
	innerRadius := spec.InnerRadius
	if innerRadius == 0 {
		innerRadius = maxRadius - spec.ArcWidth
	}

	// Build entity map
	entityMap := make(map[string]*ChordEntity)
	for i := range spec.Entities {
		entityMap[spec.Entities[i].ID] = &spec.Entities[i]
	}

	// Calculate total value for each entity (sum of all relationships)
	entityTotals := make(map[string]float64)
	for _, rel := range spec.Relations {
		entityTotals[rel.Source] += rel.Value
		if rel.Source != rel.Target { // Don't double-count self-loops
			entityTotals[rel.Target] += rel.Value
		}
	}

	// Calculate total of all relationships
	totalValue := 0.0
	for _, total := range entityTotals {
		totalValue += total
	}
	if totalValue == 0 {
		totalValue = 1
	}

	// Calculate total available angle (360 - padding between arcs)
	totalPadding := spec.ArcPadding * float64(len(spec.Entities))
	availableAngle := 360.0 - totalPadding

	// Assign arc angles to each entity based on their total value
	entityArcs := make(map[string]arcSegment)
	currentAngle := 0.0

	for _, entity := range spec.Entities {
		total := entityTotals[entity.ID]
		arcAngle := (total / totalValue) * availableAngle

		entityArcs[entity.ID] = arcSegment{
			startAngle: currentAngle,
			endAngle:   currentAngle + arcAngle,
			value:      total,
		}

		currentAngle += arcAngle + spec.ArcPadding
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

	// Draw chords (relationships) first, so they appear below arcs
	for _, rel := range spec.Relations {
		sourceArc, sourceExists := entityArcs[rel.Source]
		targetArc, targetExists := entityArcs[rel.Target]

		if !sourceExists || !targetExists {
			continue
		}

		// Calculate chord positions
		// For simplicity, center the chord on each arc
		sourceAngle := (sourceArc.startAngle + sourceArc.endAngle) / 2
		targetAngle := (targetArc.startAngle + targetArc.endAngle) / 2

		// Create chord ribbon
		chordPath := createChordRibbon(centerX, centerY, innerRadius, sourceAngle, targetAngle, rel.Value/sourceArc.value*20)

		// Get chord color (use source entity color)
		chordColor := spec.DefaultColor
		if entity, ok := entityMap[rel.Source]; ok && entity.Color != "" {
			chordColor = entity.Color
		}

		chordStyle := svg.Style{
			Fill:        chordColor,
			FillOpacity: 0.5,
			Stroke:      "none",
		}
		result += svg.Path(chordPath, chordStyle) + "\n"
	}

	// Draw entity arcs
	defaultColors := []string{"#3b82f6", "#10b981", "#f59e0b", "#ef4444", "#8b5cf6", "#ec4899"}

	for idx, entity := range spec.Entities {
		arc, exists := entityArcs[entity.ID]
		if !exists {
			continue
		}

		// Get entity color
		entityColor := entity.Color
		if entityColor == "" {
			entityColor = defaultColors[idx%len(defaultColors)]
		}

		// Draw arc
		arcPath := createAnnularSector(centerX, centerY, innerRadius, maxRadius, arc.startAngle, arc.endAngle)

		arcStyle := svg.Style{
			Fill:        entityColor,
			Stroke:      "#ffffff",
			StrokeWidth: 1,
		}
		result += svg.Path(arcPath, arcStyle) + "\n"

		// Draw label
		if spec.ShowLabels && entity.Label != "" {
			// Calculate label position (middle of arc, outside radius)
			midAngle := (arc.startAngle + arc.endAngle) / 2
			angleRad := (midAngle - 90) * math.Pi / 180

			labelDistance := maxRadius + 15
			labelX := centerX + labelDistance*math.Cos(angleRad)
			labelY := centerY + labelDistance*math.Sin(angleRad)

			// Adjust text anchor based on position
			var textAnchor string
			if math.Abs(math.Cos(angleRad)) < 0.1 {
				textAnchor = "middle"
			} else if math.Cos(angleRad) > 0 {
				textAnchor = "start"
			} else {
				textAnchor = "end"
			}

			labelStyle := svg.Style{
				FontSize:         units.Px(11),
				FontFamily:       "sans-serif",
				TextAnchor:       svg.TextAnchor(textAnchor),
				DominantBaseline: svg.DominantBaselineMiddle,
			}
			result += svg.Text(entity.Label, labelX, labelY, labelStyle) + "\n"
		}
	}

	return result
}

// arcSegment stores the angular segment for an entity
type arcSegment struct {
	startAngle float64
	endAngle   float64
	value      float64
}

// createChordRibbon creates a bezier ribbon connecting two points on a circle
func createChordRibbon(cx, cy, radius, angle1Deg, angle2Deg, width float64) string {
	// Convert angles to radians
	angle1 := (angle1Deg - 90) * math.Pi / 180
	angle2 := (angle2Deg - 90) * math.Pi / 180

	// Calculate points on the circle
	x1 := cx + radius*math.Cos(angle1)
	y1 := cy + radius*math.Sin(angle1)
	x2 := cx + radius*math.Cos(angle2)
	y2 := cy + radius*math.Sin(angle2)

	// Calculate perpendicular offsets for width
	offset1X := -width/2 * math.Sin(angle1)
	offset1Y := width/2 * math.Cos(angle1)
	offset2X := -width/2 * math.Sin(angle2)
	offset2Y := width/2 * math.Cos(angle2)

	// Start and end points with width
	x1a := x1 + offset1X
	y1a := y1 + offset1Y
	x1b := x1 - offset1X
	y1b := y1 - offset1Y
	x2a := x2 + offset2X
	y2a := y2 + offset2Y
	x2b := x2 - offset2X
	y2b := y2 - offset2Y

	// Create quadratic bezier curve through center
	// Control point is at the center for a nice arc
	path := fmt.Sprintf("M %.2f %.2f ", x1a, y1a)
	path += fmt.Sprintf("Q %.2f %.2f %.2f %.2f ", cx, cy, x2a, y2a)
	path += fmt.Sprintf("L %.2f %.2f ", x2b, y2b)
	path += fmt.Sprintf("Q %.2f %.2f %.2f %.2f ", cx, cy, x1b, y1b)
	path += "Z"

	return path
}

// ChordDiagramFromMatrix creates a chord diagram from an adjacency matrix
// entities are the entity labels
// matrix is a square matrix where matrix[i][j] is the relationship from entity i to entity j
func ChordDiagramFromMatrix(entities []string, matrix [][]float64, width, height float64) string {
	n := len(entities)
	if n == 0 || len(matrix) != n {
		return ""
	}

	// Validate matrix is square
	for _, row := range matrix {
		if len(row) != n {
			return ""
		}
	}

	// Create entities
	chordEntities := make([]ChordEntity, n)
	for i, label := range entities {
		chordEntities[i] = ChordEntity{
			ID:    fmt.Sprintf("entity_%d", i),
			Label: label,
		}
	}

	// Create relations from matrix
	var relations []ChordRelation
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if matrix[i][j] > 0 {
				relations = append(relations, ChordRelation{
					Source: fmt.Sprintf("entity_%d", i),
					Target: fmt.Sprintf("entity_%d", j),
					Value:  matrix[i][j],
				})
			}
		}
	}

	spec := ChordDiagramSpec{
		Entities:   chordEntities,
		Relations:  relations,
		Width:      width,
		Height:     height,
		ShowLabels: true,
	}

	return RenderChordDiagram(spec)
}
