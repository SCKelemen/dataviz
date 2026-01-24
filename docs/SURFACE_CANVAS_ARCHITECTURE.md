# Surface, Canvas, and Context Architecture

## Overview

This document defines the rendering abstraction layer for dataviz, introducing three key concepts: **Surface**, **Canvas**, and **Context**. These abstractions enable charts to render consistently across different targets (CLI, web, GitHub, etc.) while respecting each target's capabilities and constraints.

## Motivation

Current issues:
- No abstraction for rendering targets (hardcoded SVG strings)
- No layout system for multi-chart compositions
- No way to adapt to different surface capabilities
- Color/font/size decisions are implementation details scattered throughout

We need:
- **Surface**: Where we're rendering (CLI, web, GitHub README)
- **Canvas**: What we're rendering into (viewport with margins/legends/grid)
- **Context**: How we're rendering (design tokens, color spaces, fonts)

## Core Concepts

### Surface

A **Surface** represents the rendering target with its capabilities and constraints.

```go
package surface

import (
    "github.com/SCKelemen/color"
    "github.com/SCKelemen/units"
)

// Surface represents a rendering target
type Surface interface {
    // Capabilities
    Type() SurfaceType
    SupportsColor() bool
    SupportsInteractivity() bool
    SupportsAnimation() bool
    ColorSpace() color.Space

    // Constraints
    MaxWidth() units.Length
    MaxHeight() units.Length
    Resolution() Resolution // DPI or character grid

    // Output
    Render(canvas *Canvas) (string, error)
}

// SurfaceType identifies the kind of rendering target
type SurfaceType int

const (
    SurfaceWeb SurfaceType = iota       // Web browser (SVG)
    SurfaceCLI                          // Terminal (ANSI + Unicode)
    SurfaceGitHubReadme                 // GitHub markdown (SVG with constraints)
    SurfaceGitHubRepo                   // GitHub repo badge area
    SurfacePrint                        // Print/PDF (high DPI)
    SurfaceSlack                        // Slack/Discord (limited colors)
    SurfaceEmail                        // Email clients (very limited)
)

// Resolution describes the surface's pixel density or character grid
type Resolution struct {
    Type       ResolutionType
    PixelDPI   float64 // For pixel-based surfaces
    CharWidth  int     // For character-based surfaces (terminal)
    CharHeight int
}

type ResolutionType int

const (
    ResolutionPixel ResolutionType = iota
    ResolutionCharacter
)
```

#### Surface Examples

##### Web Surface (SVG)
```go
type WebSurface struct {
    width  units.Length
    height units.Length
}

func NewWebSurface(width, height units.Length) *WebSurface {
    return &WebSurface{width: width, height: height}
}

func (s *WebSurface) Type() SurfaceType {
    return SurfaceWeb
}

func (s *WebSurface) SupportsColor() bool {
    return true // Full color support
}

func (s *WebSurface) SupportsInteractivity() bool {
    return true // Can add JavaScript
}

func (s *WebSurface) ColorSpace() color.Space {
    return color.SRGB // Web uses sRGB by default
}

func (s *WebSurface) MaxWidth() units.Length {
    return s.width
}

func (s *WebSurface) Render(canvas *Canvas) (string, error) {
    // Generate SVG
    return canvas.RenderSVG()
}
```

##### CLI Surface (Terminal)
```go
type CLISurface struct {
    width  int // Characters
    height int // Lines
}

func NewCLISurface(width, height int) *CLISurface {
    return &CLISurface{width: width, height: height}
}

func (s *CLISurface) Type() SurfaceType {
    return SurfaceCLI
}

func (s *CLISurface) SupportsColor() bool {
    return true // ANSI 256 colors + truecolor
}

func (s *CLISurface) SupportsInteractivity() bool {
    return false // Static output (bubbletea provides separate interactivity)
}

func (s *CLISurface) ColorSpace() color.Space {
    return color.SRGB // Terminal emulators use sRGB
}

func (s *CLISurface) Resolution() Resolution {
    return Resolution{
        Type:       ResolutionCharacter,
        CharWidth:  s.width,
        CharHeight: s.height,
    }
}

func (s *CLISurface) Render(canvas *Canvas) (string, error) {
    // Generate ANSI/Unicode/Braille
    return canvas.RenderTerminal(s.width, s.height)
}
```

##### GitHub README Surface
```go
type GitHubReadmeSurface struct {
    width  int // Pixels (GitHub renders SVGs)
    height int
}

func NewGitHubReadmeSurface() *GitHubReadmeSurface {
    return &GitHubReadmeSurface{
        width:  800,  // GitHub README typical width
        height: 400,
    }
}

func (s *GitHubReadmeSurface) Type() SurfaceType {
    return SurfaceGitHubReadme
}

func (s *GitHubReadmeSurface) SupportsColor() bool {
    return true
}

func (s *GitHubReadmeSurface) SupportsInteractivity() bool {
    return false // GitHub strips JavaScript from SVGs
}

func (s *GitHubReadmeSurface) SupportsAnimation() bool {
    return false // GitHub strips animations
}

func (s *GitHubReadmeSurface) Render(canvas *Canvas) (string, error) {
    // Generate sanitized SVG (no scripts, no animations)
    svg, err := canvas.RenderSVG()
    if err != nil {
        return "", err
    }
    return sanitizeForGitHub(svg), nil
}
```

### Canvas

A **Canvas** represents the rendering viewport, including margins, legends, titles, and layout for single or multiple charts.

```go
package canvas

import (
    "github.com/SCKelemen/layout"
    "github.com/SCKelemen/units"
    "github.com/SCKelemen/dataviz/charts/legends"
)

// Canvas represents the rendering viewport
type Canvas struct {
    Width   units.Length
    Height  units.Length
    Margins Margins

    // Content
    Charts  []ChartElement
    Legend  *legends.Legend
    Title   *Title
    Axes    *Axes

    // Layout
    LayoutMode LayoutMode
    Grid       *GridLayout // For multi-chart layouts

    // Context
    Context *Context
}

// Margins define space around the chart content
type Margins struct {
    Top    units.Length
    Right  units.Length
    Bottom units.Length
    Left   units.Length
}

// ChartElement represents a chart in the canvas
type ChartElement struct {
    Chart    interface{} // Any chart type (line, bar, etc.)
    Position Rect        // Position within canvas
    ZIndex   int         // Stacking order
}

// Rect defines a rectangular region
type Rect struct {
    X      units.Length
    Y      units.Length
    Width  units.Length
    Height units.Length
}

// Title represents the canvas title
type Title struct {
    Text      string
    Subtitle  string
    Position  TitlePosition
    FontSize  units.Length
    FontFamily string
    Color     color.Color
}

type TitlePosition int

const (
    TitleTop TitlePosition = iota
    TitleBottom
    TitleLeft
    TitleRight
)

// LayoutMode defines how charts are arranged
type LayoutMode int

const (
    LayoutSingle LayoutMode = iota // Single chart with margins
    LayoutGrid                     // Grid of charts
    LayoutFlex                     // Flexbox layout
    LayoutAbsolute                 // Absolute positioning
)

// GridLayout defines a grid of charts
type GridLayout struct {
    Rows    int
    Cols    int
    Gap     units.Length
    Align   layout.Align
}
```

#### Canvas Examples

##### Single Chart Canvas
```go
// Create a canvas for a single line chart
canvas := &Canvas{
    Width:  units.Px(800),
    Height: units.Px(400),
    Margins: Margins{
        Top:    units.Px(40),  // Space for title
        Right:  units.Px(20),
        Bottom: units.Px(50),  // Space for X-axis labels
        Left:   units.Px(60),  // Space for Y-axis labels
    },
    LayoutMode: LayoutSingle,
    Title: &Title{
        Text:     "Revenue Over Time",
        Position: TitleTop,
    },
    Legend: legends.New(
        []legends.LegendItem{
            legends.Item("2023", legends.Line(mustHex("#3b82f6"))),
            legends.Item("2024", legends.Line(mustHex("#10b981"))),
        },
        legends.WithPosition(legends.PositionTopRight),
    ),
    Context: DefaultContext(),
}

// Add the chart
canvas.AddChart(lineChart, Rect{
    X:      units.Px(60),  // After left margin
    Y:      units.Px(40),  // After top margin
    Width:  units.Px(720), // Total - margins
    Height: units.Px(310), // Total - margins
})
```

##### Multi-Chart Canvas (Grid)
```go
// Create a 2x2 grid of charts
canvas := &Canvas{
    Width:  units.Px(1200),
    Height: units.Px(800),
    Margins: Margins{
        Top:    units.Px(60),
        Right:  units.Px(40),
        Bottom: units.Px(40),
        Left:   units.Px(40),
    },
    LayoutMode: LayoutGrid,
    Grid: &GridLayout{
        Rows: 2,
        Cols: 2,
        Gap:  units.Px(20),
    },
    Title: &Title{
        Text:     "Quarterly Dashboard",
        Position: TitleTop,
    },
    Context: DefaultContext(),
}

// Charts are automatically positioned in the grid
canvas.AddChart(revenueChart, nil)  // Grid position [0,0]
canvas.AddChart(profitChart, nil)   // Grid position [0,1]
canvas.AddChart(expensesChart, nil) // Grid position [1,0]
canvas.AddChart(marginChart, nil)   // Grid position [1,1]
```

##### Faceted Charts Canvas
```go
// Small multiples: one chart per category
categories := []string{"North", "South", "East", "West"}

canvas := &Canvas{
    Width:  units.Px(1600),
    Height: units.Px(400),
    LayoutMode: LayoutGrid,
    Grid: &GridLayout{
        Rows: 1,
        Cols: 4,
        Gap:  units.Px(10),
    },
    Title: &Title{
        Text: "Sales by Region",
    },
    Context: DefaultContext(),
}

for _, category := range categories {
    chart := createLineChart(dataByCategory[category])
    canvas.AddChart(chart, nil) // Auto-positioned in grid
}
```

### Context

A **Context** holds rendering state and configuration shared across a rendering session.

```go
package context

import (
    "github.com/SCKelemen/color"
    design "github.com/SCKelemen/design-system"
)

// Context holds shared rendering configuration
type Context struct {
    // Design tokens
    DesignTokens *design.DesignTokens
    MotionTokens *design.MotionTokens

    // Color management
    ColorSpace   color.Space          // OKLCH, sRGB, etc.
    ColorProfile *color.Profile       // ICC profile (optional)
    GamutMapping color.GamutMapping   // How to handle out-of-gamut colors

    // Typography
    FontMetrics  *FontMetrics
    TextRenderer TextRenderer

    // Rendering hints
    AntiAlias    bool
    HighQuality  bool // Trade speed for quality

    // Accessibility
    HighContrast bool
    ColorBlind   color.ColorBlindnessType
}

// FontMetrics provides text measurement
type FontMetrics interface {
    MeasureText(text string, fontSize units.Length, fontFamily string) TextDimensions
}

type TextDimensions struct {
    Width  units.Length
    Height units.Length
    Ascent units.Length
    Descent units.Length
}

// TextRenderer handles text rendering
type TextRenderer interface {
    RenderText(text string, x, y units.Length, style TextStyle) string
}

type TextStyle struct {
    FontSize   units.Length
    FontFamily string
    FontWeight string
    Color      color.Color
    Align      TextAlign
}

type TextAlign int

const (
    TextAlignLeft TextAlign = iota
    TextAlignCenter
    TextAlignRight
)

// DefaultContext returns a context with sensible defaults
func DefaultContext() *Context {
    return &Context{
        DesignTokens: design.DefaultTheme(),
        ColorSpace:   color.OKLCH,  // Use perceptually uniform color space
        GamutMapping: color.ClipToGamut,
        FontMetrics:  &DefaultFontMetrics{},
        TextRenderer: &SVGTextRenderer{},
        AntiAlias:    true,
        HighQuality:  true,
    }
}
```

## Integration with Observable Plot Architecture

### How Surface/Canvas/Context Relates to Scales/Marks

```
┌─────────────────────────────────────────────────────┐
│                    Application                      │
│  (User creates charts with data)                    │
└──────────────────┬──────────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────────┐
│              Observable Plot Layer                  │
│  • Scales (map data → visual coordinates)          │
│  • Marks (visual primitives: Line, Dot, Bar)       │
│  • Transforms (bin, stack, group)                  │
│  • Legends (auto-generated from marks)             │
└──────────────────┬──────────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────────┐
│              Canvas Layer                           │
│  • Layout (single chart vs grid)                   │
│  • Margins (space for axes, titles)                │
│  • Composition (multiple charts)                   │
└──────────────────┬──────────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────────┐
│              Context Layer                          │
│  • Design tokens (colors, typography)              │
│  • Color space (OKLCH for scales)                  │
│  • Font metrics (text measurement)                 │
└──────────────────┬──────────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────────┐
│              Surface Layer                          │
│  • Capabilities (color, interaction, animation)    │
│  • Constraints (size, resolution)                  │
│  • Output format (SVG, ANSI, etc.)                 │
└─────────────────────────────────────────────────────┘
```

### Example: Perceptual Color Scales with OKLCH

```go
// Context specifies OKLCH for perceptually uniform gradients
ctx := &Context{
    ColorSpace: color.OKLCH,
}

// Scale uses context's color space for interpolation
colorScale := scales.NewColorScale(
    scales.Domain(0, 100),
    scales.Range(
        mustHex("#3b82f6"), // Blue
        mustHex("#ef4444"), // Red
    ),
    scales.Interpolation(ctx.ColorSpace), // OKLCH interpolation
)

// Result: Smooth perceptual gradient instead of muddy RGB interpolation
color1 := colorScale.Apply(0)   // Pure blue
color2 := colorScale.Apply(50)  // Perceptually middle color (not muddy)
color3 := colorScale.Apply(100) // Pure red
```

### Example: Surface-Aware Rendering

```go
// Chart works the same regardless of surface
chart := plot.New(
    plot.X(scales.Linear(data, "x")),
    plot.Y(scales.Linear(data, "y")),
    plot.Marks(
        marks.Line(data),
        marks.Dot(data),
    ),
)

// Canvas handles layout
canvas := &Canvas{
    Width:  units.Px(800),
    Height: units.Px(400),
    Charts: []ChartElement{{Chart: chart}},
}

// Surface adapts output
webSurface := NewWebSurface(units.Px(800), units.Px(400))
svg, _ := webSurface.Render(canvas) // SVG output

cliSurface := NewCLISurface(80, 24)
ansi, _ := cliSurface.Render(canvas) // Terminal output

// Same chart, different surfaces, appropriate output!
```

## Implementation Using SCKelemen Foundation

### Canvas Uses layout Package

```go
import "github.com/SCKelemen/layout"

func (c *Canvas) calculateLayout() layout.Node {
    switch c.LayoutMode {
    case LayoutSingle:
        return c.layoutSingle()
    case LayoutGrid:
        return c.layoutGrid()
    case LayoutFlex:
        return c.layoutFlex()
    }
}

func (c *Canvas) layoutGrid() layout.Node {
    return layout.Grid(
        layout.Rows(c.Grid.Rows),
        layout.Cols(c.Grid.Cols),
        layout.Gap(c.Grid.Gap),
        layout.Children(
            mapCharts(c.Charts, func(chart ChartElement) layout.Node {
                return chart.ToLayoutNode()
            }),
        ),
    )
}

func (c *Canvas) layoutSingle() layout.Node {
    children := []layout.Node{}

    // Title
    if c.Title != nil {
        children = append(children, c.Title.ToLayoutNode())
    }

    // Main chart with margins
    children = append(children, layout.Box(
        layout.Padding(
            c.Margins.Top,
            c.Margins.Right,
            c.Margins.Bottom,
            c.Margins.Left,
        ),
        layout.Child(c.Charts[0].ToLayoutNode()),
    ))

    // Legend (absolutely positioned)
    if c.Legend != nil {
        children = append(children, c.Legend.ToLayoutNode())
    }

    return layout.Stack(layout.Children(children...))
}
```

### Context Uses color Package for Gradients

```go
import "github.com/SCKelemen/color"

// Perceptually uniform color interpolation
func (ctx *Context) InterpolateColor(start, end color.Color, t float64) color.Color {
    switch ctx.ColorSpace {
    case color.OKLCH:
        // Use OKLCH for perceptually uniform gradients
        return color.Mix(start, end, t, color.OKLCH)
    case color.SRGB:
        // Use sRGB for web-standard gradients
        return color.Mix(start, end, t, color.SRGB)
    default:
        return color.Mix(start, end, t, color.Linear)
    }
}

// Example: Heatmap color scale
func (ctx *Context) CreateHeatmapScale(min, max float64) func(float64) color.Color {
    lowColor := mustHex("#3b82f6")  // Blue
    highColor := mustHex("#ef4444") // Red

    return func(value float64) color.Color {
        t := (value - min) / (max - min)
        return ctx.InterpolateColor(lowColor, highColor, t)
    }
}
```

### Surface Uses cli Package for Terminal Rendering

```go
import "github.com/SCKelemen/cli"

func (s *CLISurface) Render(canvas *Canvas) (string, error) {
    // Use cli package for terminal-optimized rendering
    renderer := cli.NewTerminalRenderer(s.width, s.height)

    for _, chart := range canvas.Charts {
        // Charts implement TerminalRenderable interface
        if tr, ok := chart.Chart.(cli.TerminalRenderable); ok {
            output := tr.RenderTerminal(renderer)
            renderer.Write(output)
        }
    }

    return renderer.String(), nil
}
```

## Usage Examples

### Example 1: Single Chart for Web

```go
// Create chart
chart := createLineChart(data)

// Create canvas
canvas := &Canvas{
    Width:  units.Px(800),
    Height: units.Px(400),
    Margins: Margins{
        Top:    units.Px(40),
        Right:  units.Px(20),
        Bottom: units.Px(50),
        Left:   units.Px(60),
    },
    Title: &Title{Text: "Revenue Over Time"},
    Legend: createLegend(),
    Context: DefaultContext(),
}
canvas.AddChart(chart, nil)

// Render to web surface
surface := NewWebSurface(units.Px(800), units.Px(400))
svg, _ := surface.Render(canvas)

// Save or serve
os.WriteFile("chart.svg", []byte(svg), 0644)
```

### Example 2: Dashboard with Grid Layout

```go
// Create multiple charts
charts := []interface{}{
    createLineChart(revenueData),
    createBarChart(categoryData),
    createScatterPlot(correlationData),
    createHeatmap(activityData),
}

// Create grid canvas
canvas := &Canvas{
    Width:  units.Px(1600),
    Height: units.Px(1200),
    LayoutMode: LayoutGrid,
    Grid: &GridLayout{
        Rows: 2,
        Cols: 2,
        Gap:  units.Px(20),
    },
    Title: &Title{Text: "Q4 Dashboard"},
    Context: DefaultContext(),
}

for _, chart := range charts {
    canvas.AddChart(chart, nil) // Auto-positioned
}

// Render
surface := NewWebSurface(units.Px(1600), units.Px(1200))
svg, _ := surface.Render(canvas)
```

### Example 3: Same Chart, Multiple Surfaces

```go
// Create chart once
chart := createLineChart(data)
canvas := createCanvas(chart)

// Render to web
webSurface := NewWebSurface(units.Px(800), units.Px(400))
svg, _ := webSurface.Render(canvas)
os.WriteFile("chart.svg", []byte(svg), 0644)

// Render to terminal
cliSurface := NewCLISurface(80, 24)
ansi, _ := cliSurface.Render(canvas)
fmt.Println(ansi)

// Render to GitHub README
githubSurface := NewGitHubReadmeSurface()
githubSvg, _ := githubSurface.Render(canvas)
os.WriteFile("chart-github.svg", []byte(githubSvg), 0644)

// Same chart, three outputs!
```

### Example 4: Accessibility-Aware Context

```go
// High contrast mode for accessibility
ctx := &Context{
    DesignTokens: design.DefaultTheme(),
    ColorSpace:   color.OKLCH,
    HighContrast: true,
    ColorBlind:   color.Deuteranopia, // Simulate for testing
}

canvas := &Canvas{
    Context: ctx,
    // ... rest of canvas setup
}

// Charts automatically use high-contrast colors
// and avoid problematic color combinations
```

## Implementation Phases

### Phase 1: Core Abstractions (v1.6.0) - 2 weeks
- Define Surface, Canvas, Context interfaces
- Implement WebSurface and CLISurface
- Basic Canvas with single chart layout
- Context with design tokens integration

### Phase 2: Advanced Canvas (v1.7.0) - 2 weeks
- Grid layout implementation using layout package
- Flex layout implementation
- Multi-chart composition
- Legend/title positioning

### Phase 3: Surface Adapters (v1.8.0) - 2 weeks
- GitHubReadmeSurface with sanitization
- GitHubRepoSurface for badges
- PrintSurface for high-DPI output
- EmailSurface with constraints

### Phase 4: Context Enhancements (v1.9.0) - 2 weeks
- Color space integration with scales
- Font metrics from text package
- Accessibility features (high contrast, colorblind simulation)
- Performance hints

### Phase 5: Integration (v2.0.0) - 2 weeks
- Refactor all charts to use Surface/Canvas/Context
- Update examples and documentation
- Performance optimization
- Production-ready release

## Benefits

1. **Separation of Concerns**: Chart logic vs rendering vs layout
2. **Reusability**: Same chart, multiple surfaces
3. **Adaptability**: Surface-aware rendering (respects constraints)
4. **Consistency**: Shared context ensures consistent styling
5. **Composability**: Multiple charts in grids/layouts
6. **Accessibility**: Context-aware high contrast, colorblind modes
7. **Foundation Integration**: Uses layout, color, text, design-system packages
8. **Type Safety**: Units, colors, layouts are type-safe
9. **Testability**: Mock surfaces for testing
10. **Extensibility**: Easy to add new surfaces (Slack, Discord, etc.)

## Success Criteria

- ✅ All charts work with any Surface
- ✅ Canvas handles single and multi-chart layouts
- ✅ Context provides OKLCH color interpolation
- ✅ Terminal and web surfaces produce equivalent output
- ✅ GitHub surfaces respect constraints (no JS, no animation)
- ✅ Layout package used for all positioning
- ✅ Design tokens integrated via Context
- ✅ Accessibility features (high contrast, colorblind)
- ✅ Performance: < 100ms for 10-chart dashboard
- ✅ Documentation with examples for all surfaces

## Related Documents

- [ROADMAP.md](ROADMAP.md) - Observable Plot architecture (Scales, Marks, Transforms)
- [LEGEND_CONSOLIDATION.md](LEGEND_CONSOLIDATION.md) - Unified legend API
- Observable Plot: https://observablehq.com/plot
- D3.js: https://d3js.org
