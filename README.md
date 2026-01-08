# Dataviz

Reusable data visualization components with dual rendering support (SVG + Terminal).

## Features

- **Multiple Visualization Types**: Heatmaps, line graphs, bar charts, stat cards
- **SVG Rendering**: Complete SVG rendering implementation for web
- **Design System Integration**: Works with SCKelemen/design-system tokens
- **Color Space Support**: OKLCH and other perceptually uniform gradients via render-svg
- **Flexible Data Structures**: Type-safe data models for all visualization types

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
- Optional gradient fills (vertical fade to transparent)
- Grid lines with value labels
- Smooth rounded line caps and joins

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

```go
import (
    "github.com/SCKelemen/dataviz"
    design "github.com/SCKelemen/design-system"
)

// Create renderer
renderer := dataviz.NewSVGRenderer()

// Configure
bounds := dataviz.Bounds{X: 0, Y: 0, Width: 400, Height: 70}
config := dataviz.RenderConfig{
    DesignTokens: design.DefaultTheme(),
    Color:        "#E5E7EB",
}

// Render heatmap
heatmapData := dataviz.HeatmapData{
    Days: []dataviz.ContributionDay{
        {Date: time.Now(), Count: 10},
        // ... more days
    },
    Type: "linear", // or "weeks"
}
output := renderer.RenderHeatmap(heatmapData, bounds, config)
svgString := output.String()

// Render line graph
lineData := dataviz.LineGraphData{
    Points: []dataviz.TimeSeriesData{
        {Date: time.Now(), Value: 100},
        // ... more points
    },
    Color:       "#3B82F6",
    FillColor:   "rgba(59, 130, 246, 0.1)",
    UseGradient: true,
}
output = renderer.RenderLineGraph(lineData, bounds, config)
```

## Dependencies

- [github.com/SCKelemen/color](https://github.com/SCKelemen/color) - Color manipulation
- [github.com/SCKelemen/layout](https://github.com/SCKelemen/layout) - Layout engine
- [github.com/SCKelemen/render-svg](https://github.com/SCKelemen/render-svg) - SVG primitives
- [github.com/SCKelemen/design-system](https://github.com/SCKelemen/design-system) - Design tokens

## Roadmap

- [ ] Terminal renderer implementation (for CLI output)
- [ ] Card container extraction
- [ ] Examples and documentation
- [ ] Tests
- [ ] Additional visualization types
