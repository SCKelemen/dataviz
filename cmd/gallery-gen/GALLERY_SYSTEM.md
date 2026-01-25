# Generic Gallery System

## Overview

The gallery generator has been refactored from 21 separate hardcoded functions (~2,200 lines) to a generic, configuration-driven system (~2,200 lines across 4 organized files). This eliminates code duplication while maintaining 100% backward compatibility.

## Architecture

### Core Components

The system consists of four key files:

1. **`main.go`** (176 lines): Core utilities and entry point
2. **`gallery.go`** (121 lines): Generic rendering engine
3. **`layouts.go`** (81 lines): Layout strategy implementations
4. **`configs.go`** (1,825 lines): Gallery configurations

### Design Patterns

#### 1. Layout Strategy Pattern

Interfaces allow flexible positioning logic:

```go
type LayoutStrategy interface {
    CalculateDimensions() GalleryDimensions
    GetCellPosition(variantIndex int) (x, y float64)
}
```

**Implementations:**
- `SingleRowLayout`: Side-by-side variants (most galleries)
- `GridLayout`: Multi-row grids (Line, Scatter, StatCard)
- `VerticalStackLayout`: Vertical stacking (Heatmap)

#### 2. Configuration-Driven Rendering

Each gallery is defined declaratively:

```go
type GalleryConfig struct {
    Name          string
    Title         string
    Layout        LayoutStrategy
    Variants      []VariantConfig
    LabelOffsetY  float64
    ChartOffsetX  float64
    ChartOffsetY  float64
}

type VariantConfig struct {
    Label         string
    DataProvider  func() interface{}       // Generates chart data
    ChartRenderer func(data interface{}, ctx RenderContext) string
}
```

#### 3. Central Registry

All galleries registered in a single map:

```go
var GalleryRegistry = map[string]GalleryConfig{
    "bar":               BarGallery,
    "line":              LineGallery,
    "scatter":           ScatterGallery,
    // ... 18 more galleries
}
```

### Generic Generation Flow

```go
func GenerateGallery(config GalleryConfig) (string, error) {
    // 1. Calculate dimensions from layout
    dims := config.Layout.CalculateDimensions()

    // 2. Render background and title
    // 3. For each variant:
    //    - Get data from DataProvider
    //    - Create RenderContext with dimensions
    //    - Call ChartRenderer
    //    - Position with Layout.GetCellPosition()

    // 4. Wrap in SVG and return
}
```

## Code Reduction

### Before: Hardcoded Functions
- **21 separate functions**: ~100 lines each
- **Total**: ~2,200 lines of repetitive boilerplate
- **Duplication**: Title rendering, background, positioning, labels
- **Maintenance**: Change one gallery → manually update 21 functions

### After: Generic System
- **main.go**: 176 lines (utilities)
- **gallery.go**: 121 lines (generic engine)
- **layouts.go**: 81 lines (3 strategies)
- **configs.go**: 1,825 lines (declarative configs)
- **Total**: 2,203 lines (well organized)
- **Reduction**: ~2,035 lines of duplicated code eliminated (92% from main.go)

## Adding a New Gallery

### Before (Old System)
~100-150 lines of boilerplate:

```go
func generateMyGallery() (string, error) {
    data := /* create data */

    // Calculate dimensions
    w, h := 600, 400
    totalWidth := w * 2
    totalHeight := h + 80

    var content string

    // Background
    content += svg.Rect(0, 0, float64(totalWidth), float64(totalHeight),
        svg.Style{Fill: "#ffffff"})

    // Title
    titleStyle := svg.Style{
        FontSize: units.Px(20),
        FontWeight: "bold",
        FontFamily: "sans-serif",
        Fill: "#000000",
        TextAnchor: "middle",
    }
    content += svg.Text("My Gallery", float64(totalWidth)/2, 30, titleStyle)

    // Label style
    labelStyle := svg.Style{ /* ... */ }

    // Variant 1
    cellX := 25.0
    content += svg.Group(
        svg.Text("Variant 1", /* ... */) +
        svg.Group(/* render chart */, /* ... */),
        fmt.Sprintf("translate(%.2f, %.2f)", cellX, 60),
        svg.Style{},
    )

    // Variant 2 (repeat boilerplate)
    cellX += float64(w)
    content += svg.Group( /* ... */ )

    return wrapSVG(content, totalWidth, totalHeight), nil
}
```

### After (Generic System)
~30-50 lines of configuration:

```go
var MyGallery = GalleryConfig{
    Name:  "my-gallery",
    Title: "My Gallery",
    Layout: &SingleRowLayout{
        Cols:       2,
        BaseWidth:  600,
        BaseHeight: 400,
        StartX:     25.0,
    },
    Variants: []VariantConfig{
        {
            Label: "Variant 1",
            DataProvider: func() interface{} {
                return MyData{/* ... */}
            },
            ChartRenderer: func(data interface{}, ctx RenderContext) string {
                myData := data.(MyData)
                return charts.RenderMyChart(myData, 0, 0,
                    int(ctx.ChartWidth), int(ctx.ChartHeight))
            },
        },
        {
            Label: "Variant 2",
            DataProvider: func() interface{} {
                return MyData{/* different config */}
            },
            ChartRenderer: func(data interface{}, ctx RenderContext) string {
                myData := data.(MyData)
                return charts.RenderMyChart(myData, 0, 0,
                    int(ctx.ChartWidth), int(ctx.ChartHeight))
            },
        },
    },
    ChartOffsetX: 0.0,
    ChartOffsetY: 30.0,
}
```

Then register it:
```go
var GalleryRegistry = map[string]GalleryConfig{
    // ... existing galleries
    "my-gallery": MyGallery,
}
```

## Layout Strategies

### SingleRowLayout

For side-by-side variants:

```go
Layout: &SingleRowLayout{
    Cols:       2,          // Number of variants
    BaseWidth:  600,        // Width per variant
    BaseHeight: 400,        // Height of row
    StartX:     25.0,       // Initial X offset
}
```

**Used by**: Bar, Area, Pie, BoxPlot, Violin, Treemap, Sunburst, Radar, etc.

### GridLayout

For multi-row grids:

```go
Layout: &GridLayout{
    Cols:       3,          // Columns
    Rows:       2,          // Rows
    BaseWidth:  450,        // Base width per cell
    BaseHeight: 350,        // Base height per cell
}
```

**Used by**: Line (2x2), Scatter (2x3), ConnectedScatter (2x3), StatCard (2x3)

### VerticalStackLayout

For vertical stacking:

```go
Layout: &VerticalStackLayout{
    Rows:       2,          // Number of rows
    BaseWidth:  800,        // Total width
    RowHeight:  250,        // Height per row
    RowSpacing: 20,         // Space between rows
}
```

**Used by**: Heatmap

## RenderContext

Provides chart dimensions to renderers:

```go
type RenderContext struct {
    ChartWidth  float64    // Available width for chart
    ChartHeight float64    // Available height for chart
    OffsetX     float64    // X offset within cell
    OffsetY     float64    // Y offset within cell
    Tokens      interface{} // Design tokens
}
```

Renderers use this to calculate final pixel dimensions:

```go
ChartRenderer: func(data interface{}, ctx RenderContext) string {
    chartW := int(ctx.ChartWidth - 50)   // Leave margin
    chartH := int(ctx.ChartHeight - 80)
    return charts.RenderMyChart(data, 0, 0, chartW, chartH)
}
```

## Helper Functions

Located in `configs.go`:

```go
// Time parsing for time-series data
func mustParseTime(s string) time.Time {
    t, _ := time.Parse("2006-01-02", s)
    return t
}

// Trend data generation for stat cards
func makeTrendData(values []int) []charts.TimeSeriesData {
    // Returns weekly time series data
}

// Violin plot data generation
func generateViolinValues(mean, stddev float64) []float64 {
    // Returns ~100 values approximating normal distribution
}

// Tree structure for hierarchical charts
func createSampleTree() *charts.TreeNode {
    // Returns sample tree with branches and leaves
}

// Heatmap contribution data
func generateHeatmapData() charts.HeatmapData {
    // Returns 365 days of contribution data
}
```

Located in `main.go`:

```go
// Core dimension calculation (also in layouts.go)
func CalculateSingleRowDimensions(cols int, baseWidth, baseHeight float64) GalleryDimensions

func CalculateGridDimensions(cols, rows int, baseWidth, baseHeight float64) GalleryDimensions
```

## Gallery Examples

### Simple Gallery (Bar)

2 variants, single row:

```go
var BarGallery = GalleryConfig{
    Name:  "bar",
    Title: "Bar Chart Gallery",
    Layout: &SingleRowLayout{
        Cols: 2, BaseWidth: 850, BaseHeight: 450, StartX: 50.0,
    },
    Variants: []VariantConfig{
        {Label: "Simple Bars", DataProvider: /* ... */, ChartRenderer: /* ... */},
        {Label: "Stacked Bars", DataProvider: /* ... */, ChartRenderer: /* ... */},
    },
    ChartOffsetY: 30.0,
}
```

### Grid Gallery (Scatter)

6 variants in 2x3 grid:

```go
var ScatterGallery = GalleryConfig{
    Name:  "scatter",
    Title: "Scatter Plot Gallery",
    Layout: &GridLayout{
        Cols: 3, Rows: 2, BaseWidth: 450, BaseHeight: 350,
    },
    Variants: []VariantConfig{
        {Label: "Marker: circle", /* ... */},
        {Label: "Marker: square", /* ... */},
        {Label: "Marker: diamond", /* ... */},
        {Label: "Marker: triangle", /* ... */},
        {Label: "Marker: cross", /* ... */},
        {Label: "Marker: x", /* ... */},
    },
    ChartOffsetY: 25.0,
}
```

### Special Case (Heatmap)

Vertical stack layout:

```go
var HeatmapGallery = GalleryConfig{
    Name:  "heatmap",
    Title: "Heatmap Gallery",
    Layout: &VerticalStackLayout{
        Rows: 2, BaseWidth: 800, RowHeight: 250, RowSpacing: 20,
    },
    Variants: []VariantConfig{
        {Label: "Linear Heatmap", /* ... */},
        {Label: "Weeks Heatmap (GitHub Style)", /* ... */},
    },
    ChartOffsetY: 25.0,
}
```

## Benefits

### Maintainability
- **Single source of truth**: Gallery structure defined once
- **Global updates**: Change title style everywhere in one place
- **Clear separation**: Layout logic separate from chart rendering
- **Type safety**: Compiler enforces consistency

### Extensibility
- **New gallery**: ~30-50 lines vs ~100-150 lines
- **New layout**: Implement `LayoutStrategy` interface
- **No boilerplate**: Generic engine handles structure
- **Easy testing**: Configuration is data, not code

### Code Quality
- **DRY principle**: Zero duplication of structural code
- **Strategy pattern**: Flexible layout strategies
- **Declarative**: Configs describe what, not how
- **Readable**: Configuration clearly shows gallery structure

## Migration History

### Batch 1: Foundation + 5 Galleries
- Created `gallery.go`, `layouts.go` infrastructure
- Migrated: Bar, Area, StackedArea, Lollipop, Histogram

### Batch 2: 6 Single-Row Galleries
- Migrated: Pie, BoxPlot, Violin, Treemap, Icicle, Ridgeline

### Batch 3: 4 Grid Galleries
- Migrated: Line (2x2), Scatter (2x3), ConnectedScatter (2x3), StatCard (2x3)

### Batch 4: 6 Special Cases
- Migrated: Radar, StreamChart, Candlestick, Sunburst, CirclePacking, Heatmap

### Cleanup
- Deleted all 21 old functions (~2,035 lines)
- Simplified `generateGalleries()` to use registry
- Final: 2,203 lines across 4 organized files

## Testing

All galleries verified with:

```bash
go build ./cmd/gallery-gen
./gallery-gen
```

Output: 21 SVG files in `examples-gallery/` directory

Verification methods:
- Byte-level comparison with old output
- Visual inspection of generated SVGs
- Dimension validation tests

## Future Enhancements

Possible improvements:

1. **CLI filtering**: `./gallery-gen bar line scatter`
2. **YAML configs**: External configuration files
3. **Parallel generation**: Generate galleries concurrently
4. **Theme support**: Dark mode, custom color schemes
5. **Snapshot tests**: Automated regression testing
6. **Custom layouts**: User-defined layout strategies

## Summary

The generic gallery system provides:
- ✅ **92% code reduction** in main.go
- ✅ **Zero duplication** of structural code
- ✅ **100% backward compatible** output
- ✅ **Type-safe** configuration
- ✅ **Easy to extend** with new galleries
- ✅ **Maintainable** single source of truth
- ✅ **Well-documented** architecture

New galleries can be added in ~30-50 lines compared to ~100-150 lines before, with no risk of structural inconsistencies.
