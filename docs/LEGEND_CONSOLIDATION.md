# Legend Consolidation Plan

## Current State Analysis

### Legend Implementations in Codebase

| Chart Type | File | Location | Style | Position | Features |
|------------|------|----------|-------|----------|----------|
| **Pie Chart** | `charts/piechart.go` | Lines 103-120 | Vertical list | Bottom | Color swatch + label + percentage |
| **Stat Card** | `charts/statcard.go` | Lines 56-85 | Horizontal list | Header (top-right) | Color swatch + label |
| **MCP Line Chart** | `mcp/charts/charts.go` | Lines 167-182 | Horizontal list | Top-left | Line sample + label |
| **Line Graph** | `charts/linegraph.go` | N/A | ❌ None | N/A | No legend support |
| **Area Chart** | `charts/areachart.go` | N/A | ❌ None | N/A | No legend support |
| **Bar Chart** | `charts/barchart.go` | N/A | ❌ None | N/A | No legend support |
| **Scatter Plot** | `charts/scatterplot.go` | N/A | ❌ None | N/A | No legend support |
| **Heatmap** | `charts/heatmap.go` | N/A | ❌ None | N/A | No legend support |

### Current Legend Type (types.go:119-124)

```go
type LegendItem struct {
	Color string
	Label string
	X     int // Optional X position (if 0, will be auto-positioned)
}
```

**Used by:** Stat cards only

**Limitations:**
- No positioning configuration (hardcoded)
- No symbol type (only implicit color swatch)
- Manual X positioning is brittle
- Not used by other charts

### Implementation Patterns

#### 1. Pie Chart Legend (charts/piechart.go:103-120)
```go
// Draw legend if enabled
if showLegend {
    legendY := height - legendHeight + 10
    for i, slice := range data.Slices {
        color := colors[i%len(colors)]
        y := legendY + i*25

        // Color box
        sb.WriteString(fmt.Sprintf(`<rect x="20" y="%d" width="15" height="15" fill="%s"/>`,
            y, color))

        // Label
        percentage := (slice.Value / total) * 100
        label := fmt.Sprintf("%s (%.1f%%)", slice.Label, percentage)
        sb.WriteString(fmt.Sprintf(`<text x="40" y="%d" font-family="Arial, sans-serif" font-size="12" fill="#333">%s</text>`,
            y+12, label))
    }
}
```

**Characteristics:**
- Vertical layout (one item per row)
- Fixed position: bottom-left
- Fixed spacing: 25px between items
- Symbol: 15x15px color swatch
- Hardcoded fonts and colors

#### 2. Stat Card Legend (charts/statcard.go:56-85)
```go
// Legends (if provided)
if data.Legend1 != "" && data.TrendColor != "" {
    legendSize := 10.0
    legendSpacing := 15.0

    // Calculate positions
    totalLegendWidth := legendSize + textSpacing + maxTextWidth
    legendX := float64(width) - float64(designTokens.Layout.CardPaddingRight) - totalLegendWidth
    legendY := 15

    // First legend
    sb.WriteString(fmt.Sprintf(`<rect x="%.2f" y="%.2f" width="%.2f" height="%.2f" fill="%s"/>`,
        legendX, legendY, legendSize, legendSize, data.TrendColor))
    sb.WriteString(fmt.Sprintf(`<text x="%.2f" y="%.2f" font-family="Arial, sans-serif" font-size="10" fill="#6B7280">%s</text>`,
        legendX+legendSize+textSpacing, legendY+legendSize-2, data.Legend1))

    // Second legend (if exists)
    // ...
}
```

**Characteristics:**
- Horizontal layout (stacked vertically if multiple)
- Position: top-right, aligned to card padding
- Uses design tokens for some measurements
- Symbol: 10x10px color swatch
- Text measurement approximation (7px per char)

#### 3. MCP Line Chart Legend (mcp/charts/charts.go:167-182)
```go
// Draw legend
legendX := margin
legendY := 50.0
for i, series := range config.Series {
    seriesColor := series.Color
    if seriesColor == "" {
        seriesColor = colors[i%len(colors)]
    }

    xOffset := float64(i * 120)
    sb.WriteString(fmt.Sprintf(`  <line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="%s" stroke-width="2"/>`,
        legendX+xOffset, legendY, legendX+xOffset+20, legendY, seriesColor))
    sb.WriteString(fmt.Sprintf(`  <text x="%.2f" y="%.2f" font-size="12" fill="#374151">%s</text>`,
        legendX+xOffset+25, legendY+4, series.Name))
}
```

**Characteristics:**
- Horizontal layout
- Position: top-left
- Fixed spacing: 120px between items
- Symbol: 20px line sample
- Hardcoded fonts and colors

### Problems with Current Approach

1. **Code Duplication**: Same rendering logic repeated 3 times
2. **Inconsistent APIs**: Different ways to configure legends
3. **Hardcoded Positioning**: No flexible positioning options
4. **Hardcoded Styling**: Fonts, colors, sizes are hardcoded strings
5. **Limited Symbols**: Only color swatches and line samples, no markers
6. **No Automatic Generation**: Must manually create legend data
7. **Poor Text Measurement**: Approximations instead of proper measurement
8. **No Layout System**: Manual pixel math instead of flexbox

## Proposed Solution

### Unified Legend API

```go
package legends

import (
    "github.com/SCKelemen/layout"
    "github.com/SCKelemen/color"
    "github.com/SCKelemen/units"
    design "github.com/SCKelemen/design-system"
)

// Legend represents a chart legend
type Legend struct {
    Items    []LegendItem
    Position LegendPosition
    Layout   LegendLayout
    Style    *LegendStyle
}

// LegendItem represents a single legend entry
type LegendItem struct {
    Label  string
    Symbol Symbol
    Value  string // Optional value display (e.g., "45%" for pie slices)
}

// Symbol represents the visual indicator for a legend item
type Symbol interface {
    Render() layout.Node
}

// Common symbol types
type ColorSwatch struct {
    Color color.Color
    Size  units.Length
}

type LineSample struct {
    Color  color.Color
    Width  units.Length
    Length units.Length
    Dash   []float64 // Optional dash pattern
}

type MarkerSymbol struct {
    Type   string       // "circle", "square", "diamond", etc.
    Color  color.Color
    Size   units.Length
}

// LegendPosition defines where the legend appears
type LegendPosition int

const (
    PositionTopLeft LegendPosition = iota
    PositionTopRight
    PositionTopCenter
    PositionBottomLeft
    PositionBottomRight
    PositionBottomCenter
    PositionLeft
    PositionRight
    PositionNone // Hide legend
)

// LegendLayout defines how items are arranged
type LegendLayout int

const (
    LayoutVertical LegendLayout = iota   // Stack vertically
    LayoutHorizontal                     // Flow horizontally
    LayoutGrid                           // Grid layout
    LayoutAuto                           // Auto-detect based on position
)

// LegendStyle contains styling configuration
type LegendStyle struct {
    Background   color.Color
    Border       color.Color
    Padding      units.Length
    ItemSpacing  units.Length
    Font         *design.TypographyToken
    TextColor    color.Color
}

// Render generates the legend as a layout node
func (l *Legend) Render(tokens *design.DesignTokens) layout.Node {
    // Use SCKelemen/layout for flexbox positioning
    // Returns a layout.Node that can be composed into the chart
}

// Auto-generate legend from chart data
func FromData(data interface{}, opts ...Option) *Legend {
    // Inspect data structure and generate legend items
    // Returns a Legend with sensible defaults
}
```

### Usage Examples

#### Example 1: Pie Chart with Unified Legend
```go
// Current API (maintain for backwards compatibility)
svg := RenderPieChart(data, 0, 0, 400, 400, "Sales by Region", false, true, true)

// New API with explicit legend
legend := legends.New(
    legends.Position(legends.PositionBottomLeft),
    legends.Layout(legends.LayoutVertical),
    legends.Items(
        legends.Item("North", legends.ColorSwatch("#FF6B6B")),
        legends.Item("South", legends.ColorSwatch("#4ECDC4")),
        legends.Item("East", legends.ColorSwatch("#45B7D1")),
    ),
)

svg := RenderPieChartWithLegend(data, 0, 0, 400, 400, "Sales by Region", legend)

// Or auto-generate from data
legend := legends.FromData(data,
    legends.Position(legends.PositionBottomLeft),
    legends.SymbolType(legends.SymbolColorSwatch),
)
```

#### Example 2: Multi-Series Line Chart
```go
legend := legends.New(
    legends.Position(legends.PositionTopRight),
    legends.Layout(legends.LayoutVertical),
    legends.Items(
        legends.Item("Revenue", legends.LineSample("#3b82f6")),
        legends.Item("Profit", legends.LineSample("#10b981")),
        legends.Item("Expenses", legends.LineSample("#ef4444", legends.Dashed())),
    ),
)
```

#### Example 3: Scatter Plot with Markers
```go
legend := legends.New(
    legends.Position(legends.PositionTopRight),
    legends.Items(
        legends.Item("Group A", legends.Marker("circle", "#3b82f6")),
        legends.Item("Group B", legends.Marker("square", "#10b981")),
        legends.Item("Group C", legends.Marker("diamond", "#f59e0b")),
    ),
)
```

### Implementation Using SCKelemen Foundation

```go
package legends

import (
    "github.com/SCKelemen/layout"
    "github.com/SCKelemen/svg"
    "github.com/SCKelemen/color"
    "github.com/SCKelemen/units"
    "github.com/SCKelemen/text"
    design "github.com/SCKelemen/design-system"
)

func (l *Legend) Render(tokens *design.DesignTokens) layout.Node {
    // Determine layout direction
    direction := layout.Column
    if l.Layout == LayoutHorizontal {
        direction = layout.Row
    }

    // Build legend items using flexbox
    items := make([]layout.Node, len(l.Items))
    for i, item := range l.Items {
        items[i] = layout.Flex(
            layout.Direction(layout.Row),
            layout.AlignItems(layout.Center),
            layout.Gap(tokens.Layout.SpacingSmall),
            layout.Children(
                item.Symbol.Render(),
                layout.Text(
                    item.Label,
                    layout.Font(tokens.Typography.Body),
                    layout.Color(l.Style.TextColor),
                ),
            ),
        )
    }

    // Create legend container with flexbox
    return layout.Flex(
        layout.Direction(direction),
        layout.Gap(tokens.Layout.SpacingMedium),
        layout.Padding(l.Style.Padding),
        layout.Background(l.Style.Background),
        layout.Border(l.Style.Border, units.Px(1)),
        layout.Children(items...),
    )
}

func (c ColorSwatch) Render() layout.Node {
    return svg.Rect(
        svg.Width(c.Size),
        svg.Height(c.Size),
        svg.Fill(c.Color),
        svg.Stroke(color.Black, units.Px(1)),
    )
}

func (l LineSample) Render() layout.Node {
    return svg.Line(
        svg.X1(units.Px(0)),
        svg.Y1(units.Px(0)),
        svg.X2(l.Length),
        svg.Y2(units.Px(0)),
        svg.Stroke(l.Color),
        svg.StrokeWidth(l.Width),
        svg.StrokeDashArray(l.Dash),
    )
}
```

## Migration Plan

### Phase 1: Extract Legend Package (v1.6.0) - Week 1
- Create `charts/legends/` package
- Implement core types (Legend, LegendItem, Symbol)
- Implement rendering using `layout` and `svg` packages
- Add tests

### Phase 2: Refactor Pie Chart (v1.6.0) - Week 1
- Update `RenderPieChart` to use legends package internally
- Maintain current API for backwards compatibility
- Add optional `*Legend` parameter variant

### Phase 3: Refactor Stat Card (v1.6.0) - Week 2
- Replace inline legend with legends package
- Update `LegendItem` type in `types.go` to alias new type
- Maintain backwards compatibility

### Phase 4: Add Legends to Other Charts (v1.6.0) - Week 2
- Add legend support to: Line, Area, Bar, Scatter, Heatmap
- Use auto-generation for simple cases
- Document legend positioning options

### Phase 5: Refactor MCP Charts (v1.6.0) - Week 3
- Update MCP line chart to use legends package
- Consolidate legend rendering code

### Phase 6: Documentation (v1.6.0) - Week 3
- Create legend guide with examples
- Document all symbol types
- Document positioning options
- Add cookbook recipes

## Benefits

1. **Single Source of Truth**: One legend implementation, used everywhere
2. **Consistency**: All legends look and behave the same
3. **Flexibility**: Position, layout, and styling are configurable
4. **Type Safety**: Use `color.Color`, `units.Length` instead of strings/floats
5. **Layout Engine**: Proper flexbox instead of manual positioning
6. **Text Measurement**: Accurate text sizing via `text` package
7. **Design Tokens**: Consistent styling via `design-system`
8. **Composability**: Legends as `layout.Node` can be composed into charts
9. **Testability**: Pure functions, easy to test in isolation
10. **Extensibility**: Easy to add new symbol types

## Success Criteria

- ✅ All charts use unified legend API
- ✅ Zero hardcoded positioning (all configurable)
- ✅ Legends use `layout` package for positioning
- ✅ Symbols use `svg` package for rendering
- ✅ Text measurement via `text` package
- ✅ Colors use `color.Color` type
- ✅ Design tokens for styling
- ✅ Backwards compatible API
- ✅ Comprehensive tests (>90% coverage)
- ✅ Documentation with examples

## Future Enhancements (v2.0+)

- **Interactive Legends**: Click to toggle series visibility
- **Legend Filtering**: Show/hide legend items
- **Custom Symbols**: User-provided SVG symbols
- **Legend Titles**: Optional title for legend groups
- **Scrollable Legends**: For charts with many series
- **Export**: Export legend separately from chart
