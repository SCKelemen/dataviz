# Dataviz

Reusable data visualization components with dual rendering support (SVG + Terminal).

## Features

- **Multiple Visualization Types**: Heatmaps, line graphs, bar charts, stat cards, area charts, scatter plots
- **Dual Rendering**: Both SVG (web) and Terminal (CLI) output from the same data
- **Smooth Curves**: Bezier curve interpolation with tension control
- **Custom Markers**: 8+ marker types (circle, square, diamond, triangle, cross, x, dot)
- **Design System Integration**: Works with SCKelemen/design-system tokens
- **Color Space Support**: OKLCH and other perceptually uniform gradients
- **Flexible Data Structures**: Type-safe data models for all visualization types
- **Comprehensive Tests**: 38.6% code coverage with unit and benchmark tests

## Installation

```bash
go get github.com/SCKelemen/dataviz
```

## Visualization Types

### Heatmaps
- **Linear**: 30-day horizontal heatmap with contribution intensity
- **Weeks**: GitHub-style year heatmap (53 weeks × 7 days grid)
- HSL-based color adjustment with 4-level contrast curve
- Automatic scaling and sizing

### Line Graphs
- Time series visualization with data points
- **Smooth curves**: Bezier interpolation with configurable tension (0-1)
- **Custom markers**: Circle, square, diamond, triangle shapes at data points
- Optional gradient fills (vertical fade to transparent)
- Grid lines with value labels
- Smooth rounded line caps and joins

### Area Charts
- Filled regions under curves with baseline support
- **Smooth curves**: Bezier interpolation for organic shapes
- Optional border lines
- Gradient fills with transparency
- Perfect for showing cumulative data or ranges

### Scatter Plots
- Individual data points with customizable markers
- **8 marker types**: circle, square, diamond, triangle, cross, x, dot
- **Per-point sizing**: Custom size for each data point
- **Point labels**: Optional text labels for specific points
- Ideal for showing correlation or distribution

### Bar Charts
- Vertical bars with automatic scaling
- Stacked bar support (primary + secondary values)
- Automatic color lightening for stacked visualization

### Stat Cards
- Statistics display with title, value, subtitle
- Change indicators with arrows
- Optional legends (dual color support)
- Mini trend graphs (single or dual stacked bars)

## Usage

### Basic Setup

```go
import (
    "github.com/SCKelemen/dataviz"
    design "github.com/SCKelemen/design-system"
)

// Create renderer (SVG or Terminal)
renderer := dataviz.NewSVGRenderer()
// Or for terminal output:
// renderer := dataviz.NewTerminalRenderer()

// Configure
bounds := dataviz.Bounds{X: 0, Y: 0, Width: 400, Height: 200}
config := dataviz.RenderConfig{
    DesignTokens: design.DefaultTheme(),
    Color:        "#3B82F6",
    Theme:        "default",
}
```

### Heatmap

```go
heatmapData := dataviz.HeatmapData{
    Days: []dataviz.ContributionDay{
        {Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), Count: 10},
        {Date: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC), Count: 15},
        // ... more days
    },
    StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
    EndDate:   time.Date(2024, 1, 30, 0, 0, 0, 0, time.UTC),
    Type:      "linear", // or "weeks" for GitHub-style grid
}
output := renderer.RenderHeatmap(heatmapData, bounds, config)
```

### Line Graph with Smooth Curves

```go
lineData := dataviz.LineGraphData{
    Points: []dataviz.TimeSeriesData{
        {Date: time.Now(), Value: 100},
        {Date: time.Now().AddDate(0, 0, 1), Value: 125},
        {Date: time.Now().AddDate(0, 0, 2), Value: 115},
        {Date: time.Now().AddDate(0, 0, 3), Value: 140},
    },
    Color:       "#3B82F6",
    FillColor:   "rgba(59, 130, 246, 0.1)",
    UseGradient: true,
    Smooth:      true,    // Enable smooth curves
    Tension:     0.3,     // Curve tension (0-1, default 0.3)
    MarkerType:  "diamond", // Add markers at data points
    MarkerSize:  4,       // Marker size in pixels
}
output := renderer.RenderLineGraph(lineData, bounds, config)
```

### Area Chart

```go
areaData := dataviz.AreaChartData{
    Points: []dataviz.TimeSeriesData{
        {Date: time.Now(), Value: 100},
        {Date: time.Now().AddDate(0, 0, 1), Value: 125},
        {Date: time.Now().AddDate(0, 0, 2), Value: 115},
    },
    Color:     "#10B981",
    FillColor: "#10B981",
    Smooth:    true,
    Tension:   0.3,
    BaselineY: 0, // Y value for baseline (default 0)
}
output := renderer.RenderAreaChart(areaData, bounds, config)
```

### Scatter Plot with Custom Markers

```go
scatterData := dataviz.ScatterPlotData{
    Points: []dataviz.ScatterPoint{
        {Date: time.Now(), Value: 100, Size: 5, Label: ""},
        {Date: time.Now().AddDate(0, 0, 1), Value: 150, Size: 8, Label: "Peak"},
        {Date: time.Now().AddDate(0, 0, 2), Value: 125, Size: 6, Label: ""},
    },
    Color:      "#F59E0B",
    MarkerType: "triangle", // circle, square, diamond, triangle, cross, x, dot
    MarkerSize: 5,          // Default size (overridden by per-point Size)
}
output := renderer.RenderScatterPlot(scatterData, bounds, config)
```

### Bar Chart

```go
barData := dataviz.BarChartData{
    Bars: []dataviz.BarData{
        {Value: 100, Secondary: 20, Label: "Week 1"},
        {Value: 150, Secondary: 30, Label: "Week 2"},
        {Value: 125, Secondary: 25, Label: "Week 3"},
    },
    Color:   "#3B82F6",
    Label:   "Weekly Stats",
    Stacked: true,
}
output := renderer.RenderBarChart(barData, bounds, config)
```

### Stat Card

```go
statData := dataviz.StatCardData{
    Title:     "Total Users",
    Value:     "12,345",
    Subtitle:  "past 30 days",
    Change:    1234,
    ChangePct: 11.2,
    Color:     "#3B82F6",
    TrendData: []dataviz.TimeSeriesData{
        {Date: time.Now().AddDate(0, 0, -6), Value: 50},
        {Date: time.Now().AddDate(0, 0, -5), Value: 60},
        // ... more trend points
    },
    Legend1: "Primary",
    Legend2: "Secondary",
}
output := renderer.RenderStatCard(statData, bounds, config)
```

## Dual Rendering

### SVG Output (Web)

```go
renderer := dataviz.NewSVGRenderer()
output := renderer.RenderLineGraph(data, bounds, config)
svgString := output.String()
// Use in web applications, save to file, etc.
```

### Terminal Output (CLI)

```go
renderer := dataviz.NewTerminalRenderer()
output := renderer.RenderLineGraph(data, bounds, config)
terminalString := output.String()
// Display in terminal with Unicode/ASCII characters
```

The terminal renderer uses:
- Block characters (█ ░ ▒ ▓) for heatmaps and area charts
- Unicode symbols (● ■ ◆ ▲ + ×) for markers in scatter plots
- ASCII art for line graphs with connecting lines
- ANSI colors when available

## Marker Types

All marker types available for line graphs and scatter plots:

| Marker | SVG | Terminal | Best For |
|--------|-----|----------|----------|
| `circle` | ○ | ● | Data points |
| `square` | ▢ | ■ | Discrete values |
| `diamond` | ◇ | ◆ | Highlights |
| `triangle` | △ | ▲ | Trends |
| `cross` | + | + | Intersections |
| `x` | × | × | Exclusions |
| `dot` | · | ● | Dense data |

## Testing

Comprehensive test suite with 38.6% code coverage:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run benchmarks
go test -bench=. -benchmem
```

Tests cover:
- All visualization types (heatmap, line, area, scatter, bar, stat card)
- Edge cases (empty data, single point, extreme values)
- Smooth curve interpolation
- Marker rendering
- Terminal and SVG rendering

## Performance

Typical rendering performance on modern hardware:

| Visualization | SVG | Terminal |
|---------------|-----|----------|
| Heatmap (365 days) | ~1-2 ms | ~500 μs |
| Line Graph (100 points) | ~300-500 μs | ~200-300 μs |
| Line Graph smooth (100 points) | ~1-2 ms | ~300-400 μs |
| Area Chart (100 points) | ~1-2 ms | ~300-500 μs |
| Scatter Plot (100 points) | ~500 μs - 1 ms | ~200-300 μs |
| Bar Chart (20 bars) | ~200-300 μs | ~100-200 μs |

## Dependencies

- [github.com/SCKelemen/svg](https://github.com/SCKelemen/svg) - SVG generation
- [github.com/SCKelemen/color](https://github.com/SCKelemen/color) - Color manipulation
- [github.com/SCKelemen/layout](https://github.com/SCKelemen/layout) - Layout engine
- [github.com/SCKelemen/design-system](https://github.com/SCKelemen/design-system) - Design tokens

## Examples

See the [viz CLI example](https://github.com/SCKelemen/cli/tree/main/examples/viz) for complete working examples of all visualization types:

```bash
cd cli/examples/viz

# SVG output
go run main.go -type line-graph -format svg -data data/linegraph_smooth.json > output.svg

# Terminal output
go run main.go -type scatter-plot -format terminal -data data/scatterplot.json

# Area chart with theme
go run main.go -type area-chart -data data/areachart.json -theme midnight
```
