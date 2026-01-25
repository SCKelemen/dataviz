# Relative Sizing in Gallery Generator

> **Note**: The gallery system has been further refactored into a generic configuration-driven system. See [GALLERY_SYSTEM.md](./GALLERY_SYSTEM.md) for the complete architecture overview.

## Overview

The gallery generator uses relative sizing with the `SCKelemen/units` package to avoid error accumulation and rounding errors. Pixel dimensions are calculated at the last possible moment. This document describes the dimension calculation system that underlies the generic gallery framework.

## Key Improvements

### Before: Hardcoded Pixels
```go
w, h := 800, 350
totalWidth := w * 3               // Hardcoded multiplication
x := col * w                      // Manual positioning
charts.RenderPieChart(data, 0, 0, w, h-70, ...)  // Direct pixel values
```

**Problems:**
- Error accumulation through multiple calculations
- Rounding errors in positioning
- Inflexible layout
- Difficult to maintain consistent spacing

### After: Relative Units
```go
dims := CalculateSingleRowDimensions(3, 800, 350)
chartW := int(dims.ChartWidth)    // Pixels calculated once at end
cellX += dims.ColWidth            // Use calculated offsets
charts.RenderPieChart(data, 0, 0, chartW, chartH, ...)
```

**Benefits:**
- Percentage-based calculations: `units.Percent(100.0 / cols)`
- Single pixel conversion at render time
- Precise positioning without accumulation
- Consistent spacing across all galleries

## Helper Functions

### `CalculateSingleRowDimensions(cols int, baseWidth, baseHeight float64) GalleryDimensions`

For single-row galleries (pie, bar, scatter, etc.)

**Parameters:**
- `cols`: Number of columns
- `baseWidth`: Base width per chart
- `baseHeight`: Base height for charts

**Returns:** `GalleryDimensions` with all calculated values

**Example:**
```go
// 3 columns, 800px wide each, 350px tall
dims := CalculateSingleRowDimensions(3, 800, 350)

// Results in:
// - TotalWidth: 2400px (800 * 3)
// - TotalHeight: 430px (350 + title + margin)
// - ChartWidth: 750px (with padding)
// - ColWidth: 800px
```

### `CalculateGridDimensions(cols, rows int, baseWidth, baseHeight float64) GalleryDimensions`

For multi-row grid galleries (line graphs, scatter with many items, etc.)

**Uses percentage-based calculations:**
```go
colPct := units.Percent(100.0 / float64(cols))
rowPct := units.Percent(100.0 / float64(rows))
```

**Benefits:**
- Grid cells are exactly equal size
- No rounding errors accumulate across rows/columns
- Margins calculated as percentages

### `GalleryDimensions` Struct

```go
type GalleryDimensions struct {
    TotalWidth   float64  // Total viewBox width
    TotalHeight  float64  // Total viewBox height
    ChartWidth   float64  // Width for each chart
    ChartHeight  float64  // Height for each chart
    ColWidth     float64  // Width of each column
    RowHeight    float64  // Height of each row
    TitleY       float64  // Y position for gallery title
    ChartStartY  float64  // Y position where charts begin
    BottomMargin float64  // Bottom margin size
}
```

## Migration Pattern

### 1. Replace hardcoded dimensions
```go
// Before:
w, h := 800, 350
totalWidth := w * 3
totalHeight := h + 30

// After:
dims := CalculateSingleRowDimensions(3, 800, 350)
```

### 2. Use calculated dimensions
```go
// Before:
svg.Rect(0, 0, float64(totalWidth), float64(totalHeight), ...)

// After:
svg.Rect(0, 0, dims.TotalWidth, dims.TotalHeight, ...)
```

### 3. Calculate chart dimensions once
```go
// Convert to int for chart renderers (pixels resolved here)
chartW := int(dims.ChartWidth)
chartH := int(dims.ChartHeight - 70)
```

### 4. Use relative positioning
```go
// Before:
translate := fmt.Sprintf("translate(%d, 50)", w*2)

// After:
cellX := dims.ColWidth * 2  // Calculated offset
translate := fmt.Sprintf("translate(%.2f, %.2f)", cellX, dims.ChartStartY)
```

### 5. Convert to int at final output
```go
// Before:
return wrapSVG(content, totalWidth, totalHeight), nil

// After:
return wrapSVG(content, int(dims.TotalWidth), int(dims.TotalHeight)), nil
```

## Gallery Patterns

### Pattern 1: Single Row (3 columns)
**Examples:** Pie, Scatter (basic)

```go
dims := CalculateSingleRowDimensions(3, 800, 350)
cellX := 0.0
for _, chart := range charts {
    // Render at cellX offset
    cellX += dims.ColWidth
}
```

### Pattern 2: Single Row (2 columns)
**Examples:** Bar, Area, Stacked Area

```go
dims := CalculateSingleRowDimensions(2, 850, 450)
```

### Pattern 3: Grid (2x2, 2x3, etc.)
**Examples:** Line (2x2), Scatter with many variations (2x3)

```go
dims := CalculateGridDimensions(2, 2, 650, 350)
for row := 0; row < 2; row++ {
    for col := 0; col < 2; col++ {
        x := float64(col) * dims.ColWidth
        y := dims.ChartStartY + float64(row) * dims.RowHeight
    }
}
```

## Completed Refactorings (21/21 = 100%)

### Single Row - 3 Columns
- ✅ **Pie Gallery**: 3 variations (regular, donut, custom colors)

### Single Row - 2 Columns
- ✅ **Bar Gallery**: Simple and stacked bars
- ✅ **Area Gallery**: Simple area and different color
- ✅ **Stacked Area Gallery**: Standard and smooth curves
- ✅ **Lollipop Gallery**: Vertical and horizontal
- ✅ **Box Plot Gallery**: Basic and with confidence notches
- ✅ **Histogram Gallery**: Count and density
- ✅ **Violin Gallery**: Basic and with box plot overlay
- ✅ **Candlestick Gallery**: With scale adjustments
- ✅ **Treemap Gallery**: Standard and with padding
- ✅ **Icicle Gallery**: Vertical and horizontal
- ✅ **Radar Gallery**: With and without grid (500x500 square)
- ✅ **Streamchart Gallery**: Center layout and smooth curves
- ✅ **Ridgeline Gallery**: Standard and with fill
- ✅ **Sunburst Gallery**: Full and with inner radius (chartSize-based)
- ✅ **Circle Packing Gallery**: Standard and with padding (chartSize-based)

### Vertical Stack - 1 Column, 2 Rows
- ✅ **Heatmap Gallery**: Linear and weeks heatmaps

### Multi-Row Grid Layouts
- ✅ **Line Gallery** (2x2): Simple, smoothed, markers, filled area
- ✅ **Scatter Gallery** (2x3): 6 marker types (circle, square, diamond, triangle, cross, x)
- ✅ **Connected Scatter Gallery** (2x3): 5 line styles (solid, dashed, dotted, dash-dot, long dash)
- ✅ **StatCard Gallery** (2x3): 6 stat cards with trends

## Gallery Types Summary

All 21 galleries have been successfully refactored to use relative sizing:

- **Single Row Layouts**: 17 galleries using `CalculateSingleRowDimensions()`
- **Vertical Stack Layout**: 1 gallery (Heatmap) using manual relative calculations
- **Grid Layouts**: 4 galleries using `CalculateGridDimensions()`

## Benefits Summary

1. **Precision**: Percentage calculations avoid floating-point accumulation
2. **Consistency**: All galleries use same calculation method
3. **Maintainability**: Change base dimensions in one place
4. **Flexibility**: Easy to adjust grid sizes and margins
5. **Correctness**: Pixels calculated once at final render time

## Progress Notes

### Final Status
- **21 of 21 galleries refactored (100%)**
- All galleries use either `CalculateSingleRowDimensions` or `CalculateGridDimensions`
- Pattern is well-established and documented
- All galleries build and render correctly

### Refactoring Process
Each gallery refactoring involved:
1. Replace `w, h := 600, 400` with `dims := CalculateSingleRowDimensions(2, 600, 400)` or `dims := CalculateGridDimensions(cols, rows, baseWidth, baseHeight)`
2. Replace `totalWidth`/`totalHeight` with `dims.TotalWidth`/`dims.TotalHeight`
3. Replace `float64(w)/2` with `dims.ColWidth/2` for labels
4. Calculate `chartW := int(dims.ChartWidth - padding)` once
5. Use `cellX += dims.ColWidth` for cell positioning (or loop-based for grids)
6. Replace `fmt.Sprintf("translate(%d, 60)", w)` with `fmt.Sprintf("translate(%.2f, %.2f)", cellX, dims.ChartStartY)`
7. Convert final dimensions to int: `wrapSVG(content, int(dims.TotalWidth), int(dims.TotalHeight))`

### Special Considerations Handled
- **Candlestick**: Updated scales to use calculated chartW and chartH
- **Sunburst/CirclePacking**: Preserved chartSize-based calculations within relative sizing framework
- **Heatmap**: Vertical stack layout with manual relative calculations
- **Line/Scatter/ConnectedScatter/StatCard**: Used `CalculateGridDimensions` for multi-row grid layouts

## Results

All 21 gallery SVG files generate successfully with correct dimensions:
- Pixel calculations deferred to render time
- No error accumulation from repeated multiplications
- Consistent spacing and positioning across all galleries
- Maintainable and flexible layout system

## Current Architecture

The relative sizing system is now integrated into a generic gallery framework:

- **`CalculateSingleRowDimensions()`** and **`CalculateGridDimensions()`** are used by `LayoutStrategy` implementations
- **Layout strategies** (`SingleRowLayout`, `GridLayout`, `VerticalStackLayout`) encapsulate dimension calculations
- **Gallery configurations** use layouts declaratively
- **Generic generator** applies dimensions consistently

See [GALLERY_SYSTEM.md](./GALLERY_SYSTEM.md) for the complete architecture documentation.
