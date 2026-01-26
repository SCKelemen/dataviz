# DataViz

**Chart and data visualization library for Go with dual output modes (SVG + Terminal).**

[![License: BearWare 1.0](https://img.shields.io/badge/license-BearWare%201.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.25+-blue.svg)](https://go.dev/dl/)
[![CI](https://github.com/SCKelemen/dataviz/actions/workflows/ci.yml/badge.svg)](https://github.com/SCKelemen/dataviz/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/SCKelemen/dataviz)](https://goreportcard.com/report/github.com/SCKelemen/dataviz)
[![codecov](https://codecov.io/gh/SCKelemen/dataviz/branch/main/graph/badge.svg)](https://codecov.io/gh/SCKelemen/dataviz)

## Overview

DataViz provides high-level charting APIs built on top of the general-purpose [SCKelemen rendering stack](https://github.com/SCKelemen/layout). It focuses on **chart-specific code**: implementations of common chart types and tools for data visualization.

## Architecture

This library provides **chart implementations and visualization tools**:

```
┌─────────────────────────────────────────────────────────┐
│          DataViz Monorepo (Chart-Specific Code)         │
│  ┌─────────────────────────────────────────────────┐   │
│  │  charts/     - Line, Area, Bar, Scatter, Heat-  │   │
│  │                map, Pie/Donut, Stat cards       │   │
│  │  mcp/        - Model Context Protocol server    │   │
│  │  cmd/        - viz-cli, dataviz-mcp binaries    │   │
│  └─────────────────────────────────────────────────┘   │
└────────────────────────┬────────────────────────────────┘
                         │ depends on
┌────────────────────────▼────────────────────────────────┐
│       SCKelemen Rendering Stack (External Deps)         │
│  ┌──────────────────┐    ┌──────────────────────────┐  │
│  │  layout          │    │  design-system           │  │
│  │  • Flexbox       │    │  • Design tokens         │  │
│  │  • CSS Grid      │    │  • Themes (optional)     │  │
│  │  • Text layout   │    │  • Radix UI integration  │  │
│  └──────────────────┘    └──────────────────────────┘  │
│                                                          │
│  ┌──────────────────┐    ┌──────────────────────────┐  │
│  │  cli             │    │  tui                     │  │
│  │  • SVG output    │    │  • Dashboard framework   │  │
│  │  • Terminal out  │    │  • Interactive UI        │  │
│  └──────────────────┘    └──────────────────────────┘  │
└────────────────────────┬────────────────────────────────┘
                         │ depends on
┌────────────────────────▼────────────────────────────────┐
│            Foundation Libraries (External)              │
│  unicode, color, units, svg, text                       │
└─────────────────────────────────────────────────────────┘
```

### Use Cases

**1. Data Visualization & Charting**
- Line graphs, area charts, bar charts, scatter plots, heatmaps, pie/donut charts, stat cards
- Time-series visualization with `time.Time` types
- Smooth curves with configurable tension (Bezier interpolation)
- Custom markers: circle, square, diamond, triangle, cross, x, dot
- Gradients, fills, stacked bars
- Design token integration for consistent styling
- Dual output: SVG for web, Terminal for CLI (where applicable)

**2. AI Agent Integration**
- MCP server for Claude Code and other MCP clients
- Generic data types (interface{}, float64) for multi-source data
- Multi-series line charts, generic XY scatter plots, matrix heatmaps
- Composable with other MCP servers (Omnitron, file systems, APIs)
- MCP acts as thin wrapper around main library where possible

**3. Command-Line Tools**
- viz-cli: Interactive terminal chart viewer
- dataviz-mcp: MCP server for AI agents

## Packages

### Chart Implementations

#### `charts/`
High-level charting API with multiple chart types:
- **Line graphs** with smooth curves, tension control, area fill, gradients, markers
- **Area charts** with smooth curves and gradient fills
- **Bar charts** with stacked support
- **Pie/Donut charts** with percentage labels and legend
- **Scatter plots** with custom markers (7 types)
- **Heatmaps** (linear and GitHub-style weeks view)
- **Stat cards** with change indicators and mini trend graphs
- **Time-series** support with `time.Time` types
- **Dual output**: SVG and Terminal rendering (where applicable)

**Features:**
- Smooth curves: Bezier interpolation with configurable tension (0-1)
- Markers: circle, square, diamond, triangle, cross, x, dot
- Gradients: Vertical fade with opacity control
- Terminal: ANSI colors + Braille dots for high-resolution
- Donut mode: Configurable inner radius for donut charts

Built on top of [SCKelemen/layout](https://github.com/SCKelemen/layout) for positioning and layout.

**Note:** Currently focused on time-series data. See [Roadmap](docs/ROADMAP.md) for future generic coordinate support via Observable Plot-style scales and marks architecture.

### MCP Server

#### `mcp/`
Model Context Protocol server for AI agents:
- **29 chart generation tools** for Claude Code and MCP clients
- **Gallery tool** for generating comparison galleries of chart variations
- Generic data types (interface{}, float64)
- Composable with other MCP servers (Omnitron, file systems, APIs)
- Data-source agnostic design

**Tools Include:**
- Statistical: bar, pie, line, scatter, histogram, boxplot, violin, density, ridgeline
- Hierarchical: treemap, sunburst, icicle, circle_packing, dendrogram
- Financial: candlestick, ohlc
- Specialized: heatmap, radar, parallel, streamchart, sankey, chord, wordcloud
- Gallery: generate_gallery (comparison galleries of chart variants)

**Architecture:**
- **Consolidated charts** (pie, bar): MCP acts as thin wrapper, calls main library
- **Generic charts** (line with multi-series, XY scatter, matrix heatmap): MCP-specific implementations that complement the time-series-focused main library
- **Gallery system**: Reusable internal/gallery package for generating comparison views
- **Future:** Unified approach via Observable Plot-style scales and marks (see [Roadmap](docs/ROADMAP.md))

### Command-Line Tools

#### `cmd/viz-cli/`
Interactive terminal chart viewer:
```bash
viz-cli data.json              # Visualize JSON data
viz-cli --watch data.json      # Watch for changes
viz-cli --output chart.svg     # Export to SVG
```

#### `cmd/dataviz-mcp/`
MCP server binary for Claude Code integration:
```bash
dataviz-mcp                    # Start MCP server
```

Configure in Claude Code:
```json
{
  "mcpServers": {
    "dataviz": {
      "command": "dataviz-mcp"
    }
  }
}
```

## External Dependencies

This library depends on the **SCKelemen rendering stack**:

- **[layout](https://github.com/SCKelemen/layout)** - CSS Grid, Flexbox, text layout
- **[cli](https://github.com/SCKelemen/cli)** - SVG and terminal renderers
- **[tui](https://github.com/SCKelemen/tui)** - Interactive dashboard framework
- **[design-system](https://github.com/SCKelemen/design-system)** - Design tokens and themes (optional)

And foundation libraries:
- **[unicode](https://github.com/SCKelemen/unicode)** - 10 UAX/UTS implementations
- **[color](https://github.com/SCKelemen/color)** - OKLCH perceptually uniform color
- **[units](https://github.com/SCKelemen/units)** - Type-safe CSS units
- **[svg](https://github.com/SCKelemen/svg)** - SVG generation primitives
- **[text](https://github.com/SCKelemen/text)** - Unicode-aware text operations

## Installation

```bash
# As a library
go get github.com/SCKelemen/dataviz

# Build CLI tools
git clone https://github.com/SCKelemen/dataviz
cd dataviz
go build -o viz-cli ./cmd/viz-cli
go build -o dataviz-mcp ./cmd/dataviz-mcp
```

## Usage Examples

### 1. Time-Series Line Chart

```go
import "github.com/SCKelemen/dataviz/charts"

// Time-series line chart
data := []charts.TimeSeriesData{
    {Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), Value: 100},
    {Date: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC), Value: 150},
    {Date: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC), Value: 120},
}

config := charts.LineChartConfig{
    Width: 800,
    Height: 400,
    Title: "Sales Over Time",
    UseGradient: true,
    ShowMarkers: true,
}

// Render to SVG
svgChart := charts.RenderLineChart(data, config)

// Or render to terminal
termChart := charts.RenderLineChartTerminal(data, config)
```

### 3. With Design Tokens (Optional)

```go
import (
    design "github.com/SCKelemen/design-system"
    "github.com/SCKelemen/dataviz/charts"
)

// Use design tokens for consistent styling
theme := design.MidnightTheme()

config := charts.LineChartConfig{
    Width: 800,
    Height: 400,
    Colors: theme.Colors.Chart,
    Typography: theme.Typography,
    Spacing: theme.Spacing,
}

svgChart := charts.RenderLineChart(data, config)
```

### 4. Interactive Dashboard

```go
import (
    "github.com/SCKelemen/tui"
    "github.com/SCKelemen/dataviz/charts"
    tea "github.com/charmbracelet/bubbletea"
)

func main() {
    // Create dashboard model with charts
    model := tui.NewDashboard()
    model.AddChart("Sales", salesChartData)

    // Run with bubbletea
    p := tea.NewProgram(model)
    if err := p.Start(); err != nil {
        log.Fatal(err)
    }
}
```

## Key Design Principles

### 1. Layered Architecture
- **Low-level**: Layout engine (flexbox, grid, text)
- **Mid-level**: Renderers (SVG, terminal)
- **High-level**: Charts API (uses layout + renderers)

This separation allows you to:
- Use the layout engine without charts
- Use the renderers for custom visualizations
- Use the charts API for quick results

### 2. Dual Output Modes
Charts can be rendered to:
- **SVG** - For web, documentation, high-quality printing
- **Terminal** - For CLI tools, SSH sessions, logs (where applicable: line, area, bar, scatter, heatmap)

### 3. Optional Design Tokens
Design tokens are **opt-in**:
- Use them for consistent styling across your app
- Or don't - the rendering engine works fine without them
- Themes: midnight, nord, paper, wrapped

### 4. Data-Source Agnostic
The MCP server and charts API accept generic data:
- `interface{}` for X values (can be time.Time, int, float64, string)
- `float64` for Y values
- Works with data from any source (Omnitron, databases, APIs, files)

### 5. Type-Safe
- Layout uses type-safe CSS units
- Charts use proper types (time.Time for time-series)
- Compile-time safety where possible

## Package Import Paths

```go
// DataViz packages (this monorepo)
import "github.com/SCKelemen/dataviz/charts"  // Chart implementations
import "github.com/SCKelemen/dataviz/mcp"     // MCP server (usually not imported, used as binary)

// External rendering stack (separate repos)
import "github.com/SCKelemen/layout"          // CSS Grid, Flexbox, text layout
import "github.com/SCKelemen/cli"             // SVG and terminal renderers
import "github.com/SCKelemen/tui"             // Dashboard framework
import design "github.com/SCKelemen/design-system"  // Design tokens (optional)

// Foundation libraries (separate repos)
import "github.com/SCKelemen/color"           // OKLCH color operations
import "github.com/SCKelemen/svg"             // SVG primitives
import "github.com/SCKelemen/text"            // Unicode text operations
```

## Project Structure

```
github.com/SCKelemen/dataviz/
├── charts/          # Chart implementations (line, area, bar, scatter, heatmap, pie/donut, stat cards)
├── mcp/             # MCP server implementation
│   ├── charts/      # MCP chart handlers (thin wrappers + generic implementations)
│   ├── types/       # MCP type definitions
│   └── mcp/         # MCP protocol implementation
├── cmd/
│   ├── viz-cli/     # CLI binary for terminal charts
│   └── dataviz-mcp/ # MCP server binary
├── examples/        # Example code and data files
└── docs/            # Documentation (ROADMAP.md, etc.)
```

## Related Projects

This monorepo consolidates three previous repositories:
- **dataviz** (archived) - Original core library
- **viz-cli** (archived) - Original CLI tool
- **dataviz-mcp** (archived) - Original MCP server

### SCKelemen Foundation Libraries
- [SCKelemen/unicode](https://github.com/SCKelemen/unicode) - 10 UAX/UTS implementations (monorepo)
- [SCKelemen/color](https://github.com/SCKelemen/color) - OKLCH perceptually uniform color
- [SCKelemen/units](https://github.com/SCKelemen/units) - Type-safe CSS units
- [SCKelemen/svg](https://github.com/SCKelemen/svg) - SVG generation primitives
- [SCKelemen/clix](https://github.com/SCKelemen/clix) - CLI framework with extensions

## Contributing

Contributions welcome! This is a monorepo with multiple packages, so please:
1. Keep changes to relevant packages
2. Update tests in the same commit
3. Follow existing code style
4. Update documentation

## License

BearWare 1.0 - MIT-compatible license. See [LICENSE](LICENSE) for details.

Help the bear. 🐻🐼🐻‍❄️

## FAQ

### Is this just a charting library?

No! The core is a **general-purpose layout and rendering engine**. Charts are a high-level API built on top. You can use the layout engine to render:
- Flexbox layouts
- CSS Grid layouts
- Text with proper Unicode handling
- Any custom visualization

### Can I use the layout engine without charts?

Yes! Import `layout/` and `render/svg/` (or `render/terminal/`) directly.

### Do I need to use design tokens?

No, design tokens in `design/` are optional. The rendering engine works fine without them.

### What's the difference between render/svg/ and charts/?

- `render/svg/` is a **general-purpose SVG renderer** (can render any layout)
- `charts/` is a **high-level API** for common chart types (uses render/svg/ internally)

### Can I render to PNG/JPEG?

SVG output can be converted to PNG/JPEG using external tools or libraries. The core library focuses on SVG and terminal rendering.

### What's MCP?

Model Context Protocol - a standard for integrating tools with LLMs like Claude. The `mcp/` package implements a server that exposes chart generation as MCP tools.

### Why a monorepo?

- Single source of truth for rendering logic
- Shared layout engine across all outputs
- Atomic commits across packages
- Easier testing and CI/CD
- Consistent versioning

### What happened to the old repos?

The previous repos (dataviz, viz-cli, dataviz-mcp) have been archived and redirect here. Import paths have changed to `github.com/SCKelemen/dataviz/*`.
