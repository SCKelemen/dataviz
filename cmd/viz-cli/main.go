package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/SCKelemen/dataviz/charts"
	design "github.com/SCKelemen/design-system"
)

const usage = `viz-cli - Data visualization tool with SVG and terminal output

Usage:
  viz-cli [options]

Options:
  -type string
        Visualization type: heatmap, line-graph, bar-chart, stat-card (default "heatmap")
  -format string
        Output format: svg, terminal (default "terminal")
  -data string
        Path to JSON data file (or use stdin with -)
  -theme string
        Theme: default, midnight, nord, paper, wrapped (default "default")
  -width int
        Width in pixels (SVG) or characters (terminal) (default 80)
  -height int
        Height in pixels (SVG) or characters (terminal) (default 24)
  -color string
        Primary color for visualization (hex format) (default "#3B82F6")

Data Formats:
  Heatmap:      {"days": [{"date": "2024-01-01T00:00:00Z", "count": 10}, ...], "type": "linear"}
  Line Graph:   {"points": [{"date": "2024-01-01T00:00:00Z", "value": 100}, ...], "color": "#3B82F6"}
  Bar Chart:    {"bars": [{"value": 100, "secondary": 50, "label": "Item 1"}, ...], "color": "#3B82F6"}
  Stat Card:    {"title": "Total", "value": "1,234", "subtitle": "past month", "color": "#3B82F6"}

Examples:
  # Terminal heatmap from file
  viz-cli -type heatmap -format terminal -data contributions.json

  # SVG line graph from stdin
  cat metrics.json | viz-cli -type line-graph -format svg > output.svg

  # Terminal bar chart with custom theme
  viz-cli -type bar-chart -data repos.json -theme midnight
`

type Config struct {
	vizType  string
	format   string
	dataFile string
	theme    string
	width    int
	height   int
	color    string
}

func main() {
	cfg := parseFlags()

	if cfg.dataFile == "" || cfg.dataFile == "-" {
		fmt.Fprintln(os.Stderr, "Reading from stdin...")
	}

	// Read data
	data, err := readData(cfg.dataFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading data: %v\n", err)
		os.Exit(1)
	}

	// Get design tokens
	tokens := getTheme(cfg.theme)

	// Create bounds and config
	bounds := charts.Bounds{X: 0, Y: 0, Width: cfg.width, Height: cfg.height}
	renderConfig := charts.RenderConfig{
		DesignTokens: tokens,
		Color:        cfg.color,
		Theme:        cfg.theme,
	}

	// Choose renderer
	var output charts.Output
	switch cfg.format {
	case "svg":
		output = renderSVG(cfg.vizType, data, bounds, renderConfig)
	case "terminal":
		output = renderTerminal(cfg.vizType, data, bounds, renderConfig)
	default:
		fmt.Fprintf(os.Stderr, "Unknown format: %s\n", cfg.format)
		os.Exit(1)
	}

	// Output result
	fmt.Print(output.String())
}

func parseFlags() Config {
	cfg := Config{}

	flag.StringVar(&cfg.vizType, "type", "heatmap", "Visualization type")
	flag.StringVar(&cfg.format, "format", "terminal", "Output format")
	flag.StringVar(&cfg.dataFile, "data", "-", "Data file path")
	flag.StringVar(&cfg.theme, "theme", "default", "Theme name")
	flag.IntVar(&cfg.width, "width", 80, "Width")
	flag.IntVar(&cfg.height, "height", 24, "Height")
	flag.StringVar(&cfg.color, "color", "#3B82F6", "Primary color")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
	}

	flag.Parse()

	return cfg
}

func readData(path string) ([]byte, error) {
	if path == "" || path == "-" {
		return io.ReadAll(os.Stdin)
	}
	return os.ReadFile(path)
}

func getTheme(name string) *design.DesignTokens {
	switch name {
	case "midnight":
		return design.MidnightTheme()
	case "nord":
		return design.NordTheme()
	case "paper":
		return design.PaperTheme()
	case "wrapped":
		return design.WrappedTheme()
	default:
		return design.DefaultTheme()
	}
}

func renderSVG(vizType string, data []byte, bounds charts.Bounds, config charts.RenderConfig) charts.Output {
	renderer := charts.NewSVGRenderer()
	return renderVisualization(renderer, vizType, data, bounds, config)
}

func renderTerminal(vizType string, data []byte, bounds charts.Bounds, config charts.RenderConfig) charts.Output {
	renderer := charts.NewTerminalRenderer()
	return renderVisualization(renderer, vizType, data, bounds, config)
}

func renderVisualization(r charts.Renderer, vizType string, data []byte, bounds charts.Bounds, config charts.RenderConfig) charts.Output {
	switch vizType {
	case "heatmap":
		var heatmapData charts.HeatmapData
		if err := json.Unmarshal(data, &heatmapData); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing heatmap data: %v\n", err)
			os.Exit(1)
		}
		return r.RenderHeatmap(heatmapData, bounds, config)

	case "line-graph":
		var lineData charts.LineGraphData
		if err := json.Unmarshal(data, &lineData); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing line graph data: %v\n", err)
			os.Exit(1)
		}
		return r.RenderLineGraph(lineData, bounds, config)

	case "bar-chart":
		var barData charts.BarChartData
		if err := json.Unmarshal(data, &barData); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing bar chart data: %v\n", err)
			os.Exit(1)
		}
		return r.RenderBarChart(barData, bounds, config)

	case "stat-card":
		var statData charts.StatCardData
		if err := json.Unmarshal(data, &statData); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing stat card data: %v\n", err)
			os.Exit(1)
		}
		return r.RenderStatCard(statData, bounds, config)

	default:
		fmt.Fprintf(os.Stderr, "Unknown visualization type: %s\n", vizType)
		os.Exit(1)
		return nil
	}
}

// Helper to create sample data for testing
func createSampleHeatmap() charts.HeatmapData {
	days := make([]charts.ContributionDay, 30)
	startDate := time.Now().AddDate(0, 0, -30)
	for i := 0; i < 30; i++ {
		days[i] = charts.ContributionDay{
			Date:  startDate.AddDate(0, 0, i),
			Count: (i * 3) % 20,
		}
	}
	return charts.HeatmapData{
		Days:      days,
		StartDate: startDate,
		EndDate:   time.Now(),
		Type:      "linear",
	}
}
