# Relative Sizing in Gallery Generator

## Overview

The gallery generator has been refactored to use relative sizing with the `SCKelemen/units` package to avoid error accumulation and rounding errors. Pixel dimensions are calculated at the last possible moment.

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

## Completed Refactorings

- ✅ **Pie Gallery**: 3 columns, single row
- ✅ **Bar Gallery**: 2 columns, single row

## Remaining Galleries (19)

- Line (2x2 grid)
- Scatter (2x3 grid)
- Connected Scatter (1x5 row)
- Area (1x2 row)
- Stacked Area (1x2 row)
- Heatmap (1x2 row)
- StatCard (2x3 grid)
- BoxPlot (1x2 row)
- Histogram (1x2 row)
- Violin (1x2 row)
- Lollipop (1x2 row)
- Candlestick (1x2 row)
- Treemap (1x2 row)
- Sunburst (1x2 row)
- Circle Packing (1x2 row)
- Icicle (1x2 row)
- Radar (1x2 row)
- Streamchart (1x2 row)
- Ridgeline (1x2 row)

## Benefits Summary

1. **Precision**: Percentage calculations avoid floating-point accumulation
2. **Consistency**: All galleries use same calculation method
3. **Maintainability**: Change base dimensions in one place
4. **Flexibility**: Easy to adjust grid sizes and margins
5. **Correctness**: Pixels calculated once at final render time

## Next Steps

To complete the refactoring:
1. Apply `CalculateSingleRowDimensions` to all 2-column galleries
2. Apply `CalculateGridDimensions` to Line and Scatter galleries
3. Verify all galleries render correctly
4. Remove old hardcoded dimension calculations
