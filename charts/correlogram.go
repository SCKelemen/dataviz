package charts

import (
	"fmt"
	"math"

	"github.com/SCKelemen/svg"
	"github.com/SCKelemen/units"
)

// CorrelationMatrix represents a correlation matrix
type CorrelationMatrix struct {
	Variables []string    // Variable names
	Matrix    [][]float64 // Correlation values (must be square matrix)
}

// CorrelogramSpec configures correlogram rendering
type CorrelogramSpec struct {
	Data           CorrelationMatrix
	Width          float64
	Height         float64
	ShowValues     bool   // Show correlation values in cells
	ShowDiagonal   bool   // Show diagonal (always 1.0)
	TriangleMode   string // "full", "upper", "lower"
	ColorScheme    string // "redblue", "bluered", "coolwarm"
	Title          string
	CellPadding    float64 // Padding between cells
}

// RenderCorrelogram generates an SVG correlogram (correlation matrix heatmap)
func RenderCorrelogram(spec CorrelogramSpec) string {
	if len(spec.Data.Variables) == 0 || len(spec.Data.Matrix) == 0 {
		return ""
	}

	n := len(spec.Data.Variables)

	// Validate matrix is square
	for _, row := range spec.Data.Matrix {
		if len(row) != n {
			return ""
		}
	}

	// Set defaults
	if spec.TriangleMode == "" {
		spec.TriangleMode = "full"
	}
	if spec.ColorScheme == "" {
		spec.ColorScheme = "redblue"
	}
	if spec.CellPadding == 0 {
		spec.CellPadding = 2
	}

	// Calculate dimensions
	margin := 100.0 // Need space for labels
	chartWidth := spec.Width - (2 * margin)
	chartHeight := spec.Height - (2 * margin)

	cellWidth := chartWidth / float64(n)
	cellHeight := chartHeight / float64(n)
	cellSize := math.Min(cellWidth, cellHeight)

	// Adjust chart dimensions to be square
	chartWidth = cellSize * float64(n)
	chartHeight = cellSize * float64(n)

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

	// Draw correlation cells
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			// Check triangle mode
			skipCell := false
			if spec.TriangleMode == "upper" && i > j {
				skipCell = true
			} else if spec.TriangleMode == "lower" && i < j {
				skipCell = true
			}

			// Check diagonal
			if !spec.ShowDiagonal && i == j {
				skipCell = true
			}

			if skipCell {
				continue
			}

			correlation := spec.Data.Matrix[i][j]

			// Calculate position
			x := margin + float64(j)*cellSize
			y := margin + float64(i)*cellSize

			// Get color based on correlation value
			cellColor := getCorrelationColor(correlation, spec.ColorScheme)

			// Draw cell
			cellStyle := svg.Style{
				Fill:        cellColor,
				Stroke:      "#ffffff",
				StrokeWidth: spec.CellPadding,
			}
			result += svg.Rect(x, y, cellSize, cellSize, cellStyle) + "\n"

			// Draw correlation value if enabled
			if spec.ShowValues {
				textColor := "#ffffff"
				// Use black text for light colors (correlations near 0)
				if math.Abs(correlation) < 0.3 {
					textColor = "#000000"
				}

				valueStyle := svg.Style{
					Fill:             textColor,
					FontSize:         units.Px(math.Min(12, cellSize*0.3)),
					FontFamily:       "sans-serif",
					TextAnchor:       svg.TextAnchorMiddle,
					DominantBaseline: svg.DominantBaselineMiddle,
				}
				valueText := fmt.Sprintf("%.2f", correlation)
				result += svg.Text(valueText, x+cellSize/2, y+cellSize/2, valueStyle) + "\n"
			}
		}
	}

	// Draw row labels (left side)
	labelStyle := svg.Style{
		FontSize:         units.Px(10),
		FontFamily:       "sans-serif",
		TextAnchor:       svg.TextAnchorEnd,
		DominantBaseline: svg.DominantBaselineMiddle,
	}
	for i, varName := range spec.Data.Variables {
		y := margin + float64(i)*cellSize + cellSize/2
		result += svg.Text(varName, margin-10, y, labelStyle) + "\n"
	}

	// Draw column labels (bottom)
	for j, varName := range spec.Data.Variables {
		x := margin + float64(j)*cellSize + cellSize/2
		y := margin + chartHeight + 10

		// Rotate label for better fit
		result += fmt.Sprintf(`<text x="%.2f" y="%.2f" text-anchor="start" font-size="10" font-family="sans-serif" transform="rotate(45 %.2f %.2f)">%s</text>`,
			x, y, x, y, varName) + "\n"
	}

	// Draw color scale legend
	result += renderCorrelationLegend(spec.Width-margin+20, margin, 20, chartHeight, spec.ColorScheme)

	return result
}

// getCorrelationColor returns a color based on correlation value and color scheme
func getCorrelationColor(correlation float64, scheme string) string {
	// Clamp correlation to [-1, 1]
	if correlation < -1 {
		correlation = -1
	}
	if correlation > 1 {
		correlation = 1
	}

	// Map correlation to color
	switch scheme {
	case "redblue":
		// -1 = red, 0 = white, +1 = blue
		return interpolateRedBlue(correlation)
	case "bluered":
		// -1 = blue, 0 = white, +1 = red
		return interpolateRedBlue(-correlation)
	case "coolwarm":
		// -1 = cool (blue), 0 = neutral, +1 = warm (red)
		return interpolateCoolWarm(correlation)
	default:
		return interpolateRedBlue(correlation)
	}
}

// interpolateRedBlue maps correlation to red-white-blue color scale
func interpolateRedBlue(value float64) string {
	if value < 0 {
		// Negative: white to red
		intensity := uint8(255 * (1 + value))
		return fmt.Sprintf("rgb(255,%d,%d)", intensity, intensity)
	} else {
		// Positive: white to blue
		intensity := uint8(255 * (1 - value))
		return fmt.Sprintf("rgb(%d,%d,255)", intensity, intensity)
	}
}

// interpolateCoolWarm maps correlation to cool-warm color scale
func interpolateCoolWarm(value float64) string {
	if value < 0 {
		// Negative: neutral to cool (blue)
		t := -value
		r := uint8(220 * (1 - t) + 59*t)
		g := uint8(220 * (1 - t) + 76*t)
		b := uint8(220 * (1 - t) + 192*t)
		return fmt.Sprintf("rgb(%d,%d,%d)", r, g, b)
	} else {
		// Positive: neutral to warm (red)
		t := value
		r := uint8(220 * (1 - t) + 220*t)
		g := uint8(220 * (1 - t) + 50*t)
		b := uint8(220 * (1 - t) + 47*t)
		return fmt.Sprintf("rgb(%d,%d,%d)", r, g, b)
	}
}

// renderCorrelationLegend draws a vertical color scale legend
func renderCorrelationLegend(x, y, width, height float64, scheme string) string {
	var result string

	steps := 50
	stepHeight := height / float64(steps)

	// Draw color gradient bars
	for i := 0; i < steps; i++ {
		// Map step to correlation value (-1 to +1)
		correlation := 1.0 - (2.0 * float64(i) / float64(steps-1))
		color := getCorrelationColor(correlation, scheme)

		yPos := y + float64(i)*stepHeight

		rectStyle := svg.Style{
			Fill:   color,
			Stroke: "none",
		}
		result += svg.Rect(x, yPos, width, stepHeight+1, rectStyle) + "\n"
	}

	// Draw border around legend
	borderStyle := svg.Style{
		Fill:        "none",
		Stroke:      "#374151",
		StrokeWidth: 1,
	}
	result += svg.Rect(x, y, width, height, borderStyle) + "\n"

	// Draw scale labels
	labelStyle := svg.Style{
		FontSize:         units.Px(10),
		FontFamily:       "sans-serif",
		DominantBaseline: svg.DominantBaselineMiddle,
	}

	// +1 label (top)
	result += svg.Text("+1", x+width+5, y, labelStyle) + "\n"

	// 0 label (middle)
	result += svg.Text("0", x+width+5, y+height/2, labelStyle) + "\n"

	// -1 label (bottom)
	result += svg.Text("-1", x+width+5, y+height, labelStyle) + "\n"

	return result
}

// CalculateCorrelationMatrix computes Pearson correlation matrix from data
// data is a slice of variables, where each variable is a slice of values
func CalculateCorrelationMatrix(variables []string, data [][]float64) CorrelationMatrix {
	n := len(variables)
	matrix := make([][]float64, n)
	for i := range matrix {
		matrix[i] = make([]float64, n)
	}

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i == j {
				matrix[i][j] = 1.0
			} else {
				matrix[i][j] = pearsonCorrelation(data[i], data[j])
			}
		}
	}

	return CorrelationMatrix{
		Variables: variables,
		Matrix:    matrix,
	}
}

// pearsonCorrelation computes Pearson correlation coefficient between two variables
func pearsonCorrelation(x, y []float64) float64 {
	if len(x) != len(y) || len(x) == 0 {
		return 0
	}

	n := float64(len(x))

	// Calculate means
	var sumX, sumY float64
	for i := range x {
		sumX += x[i]
		sumY += y[i]
	}
	meanX := sumX / n
	meanY := sumY / n

	// Calculate correlation
	var numerator, denomX, denomY float64
	for i := range x {
		dx := x[i] - meanX
		dy := y[i] - meanY
		numerator += dx * dy
		denomX += dx * dx
		denomY += dy * dy
	}

	if denomX == 0 || denomY == 0 {
		return 0
	}

	return numerator / math.Sqrt(denomX*denomY)
}
