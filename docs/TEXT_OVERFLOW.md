# Text Overflow Handling in Dataviz

## Overview

The dataviz library now supports intelligent text overflow handling for axis labels using the SCKelemen/text library. This prevents long labels from overlapping and maintains chart readability.

## Features

### Text Overflow Options

The `AxisStyle` struct now includes a `TextOverflow` field with three options:

- **`TextOverflowEllipsis`** (default): Truncates long labels with "..." in the middle
- **`TextOverflowWrap`**: Wraps text onto multiple lines (planned for future implementation)
- **`TextOverflowClip`**: Clips text without modification (lets SVG handle overflow)

### Automatic Width Calculation

Labels are automatically truncated based on:
- **Font size**: Larger fonts get proportionally more space
- **Axis orientation**:
  - Horizontal axes (bottom/top): ~10 characters typical
  - Vertical axes (left/right): ~15 characters typical

## Usage

### Default Behavior

By default, all axis labels use ellipsis for long text:

```go
import "github.com/SCKelemen/dataviz/axes"

axis := axes.NewAxis(scale, axes.AxisOrientationBottom)
svg := axis.Render(axes.DefaultRenderOptions())
// Labels like "Cloud Engineering Leadership" become "Cloud Eng...dership"
```

### Custom Text Overflow

You can customize the text overflow behavior:

```go
style := axes.DefaultAxisStyle()
style.TextOverflow = axes.TextOverflowClip  // No ellipsis

opts := axes.RenderOptions{
    Style: style,
    Position: units.Px(0),
}

svg := axis.Render(opts)
```

### Disable Ellipsis

To disable text truncation entirely:

```go
style.TextOverflow = axes.TextOverflowClip
```

## Implementation Details

### Text Processing

The `processLabel()` function handles text overflow:

```go
func processLabel(label string, maxWidth float64, overflow TextOverflow) string {
    switch overflow {
    case TextOverflowEllipsis:
        return textutil.ElideLabel(label, maxWidth)
    case TextOverflowWrap:
        // TODO: Implement wrapping
        return textutil.ElideLabel(label, maxWidth)
    case TextOverflowClip:
        return label  // No processing
    }
}
```

### Ellipsis Strategy

The SCKelemen/text library uses **middle ellipsis** by default:
- `"Very Long Category Name"` → `"Very Lon...y Name"`
- Preserves both the start and end of labels for better recognition
- More readable than simple truncation: `"Very Long Categ..."`

## Examples

### Bar Chart with Long Labels

```go
data := charts.BarChartData{
    Bars: []charts.BarData{
        {Label: "Cloud Engineering Leadership", Value: 6},
        {Label: "Incident Management", Value: 8},
        {Label: "Observability Platform", Value: 5},
    },
}

// Labels automatically truncated to fit axis space
svg := charts.RenderBarChartWithAxes(data, ...)
```

Output labels:
- "Cloud Eng...dership" (22 chars → ~12 chars)
- "Incident Management" (fits without truncation)
- "Observabil...latform" (23 chars → ~12 chars)

## Future Enhancements

### Text Wrapping (Planned)

```go
style.TextOverflow = axes.TextOverflowWrap
// Will wrap long labels onto multiple lines:
// "Cloud Engineering
//  Leadership"
```

### Custom Ellipsis

The textutil package supports custom ellipsis characters:

```go
// In future versions:
style.EllipsisCharacter = "…"  // Unicode ellipsis
style.EllipsisCharacter = " [...]"  // Custom ellipsis
```

## Dependencies

- **github.com/SCKelemen/text**: Provides intelligent text elision
- **github.com/SCKelemen/unicode**: Unicode-aware width measurement

## See Also

- [SCKelemen/text Documentation](https://github.com/SCKelemen/text)
- [Axis Rendering](./axes/render.go)
- [Text Utilities](./internal/textutil/elide.go)
