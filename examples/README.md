# Dataviz Examples

Complete working examples demonstrating all features of the dataviz library.

## Running the Examples

```bash
# Navigate to examples directory
cd dataviz/examples

# Run any example
go run smooth_curves.go
go run custom_markers.go
go run area_scatter.go
go run dual_rendering.go
```

## Examples Overview

### 1. Smooth Curves (`smooth_curves.go`)

Demonstrates smooth curve interpolation with tension control.

**Features Shown:**
- Bezier curve interpolation
- Tension parameter (0.0 to 1.0)
- Comparison of different tension values
- How tension affects curve shape

**Key Learning:**
- Use `Smooth: true` to enable curves
- `Tension: 0.3` is the recommended default
- Lower tension = tighter curves
- Higher tension = looser curves

```bash
go run smooth_curves.go
```

### 2. Custom Markers (`custom_markers.go`)

Shows all 7 available marker types for data visualization.

**Features Shown:**
- Circle, square, diamond, triangle markers
- Cross, x, and dot markers
- SVG and terminal rendering comparison
- When to use each marker type

**Key Learning:**
- Set `MarkerType` to customize point shapes
- `MarkerSize` controls size in pixels
- Terminal uses Unicode symbols (●■◆▲+×)
- SVG uses actual shapes

```bash
go run custom_markers.go
```

### 3. Area Charts & Scatter Plots (`area_scatter.go`)

Demonstrates the two newest visualization types.

**Features Shown:**
- Area charts with smooth curves
- Scatter plots with custom markers
- Per-point sizing in scatter plots
- Point labels
- Dual rendering comparison

**Key Learning:**
- Area charts perfect for cumulative data
- Scatter plots ideal for correlation/distribution
- Use `ScatterPoint.Size` for custom sizing
- Add `Label` to highlight specific points

```bash
go run area_scatter.go
```

### 4. Dual Rendering (`dual_rendering.go`)

Complete demonstration of SVG and Terminal rendering from the same data.

**Features Shown:**
- SVG output for web applications
- Terminal output for CLI tools
- Rendering method comparison
- All visualization types in terminal
- Saving SVG to file

**Key Learning:**
- Same data works for both renderers
- Switch renderer in one line of code
- SVG: High resolution, scalable, web-ready
- Terminal: No graphics, works everywhere
- Terminal uses block characters and Unicode

```bash
go run dual_rendering.go
```

## Sample Output

### Smooth Curves
```
=== Smooth Curves Example ===

Tension: 0.0
SVG Length: 1234 characters

Tension: 0.3
SVG Length: 1256 characters

Tension: 0.6
SVG Length: 1289 characters
...
```

### Terminal Rendering

```
Line Graph:
        │•
      │•│
      ││•
    │•│
    ││•
  │•│
  │││
│•││•
```

### Scatter Plot
```
         ▲
       ▲
        ▲
     ▲    Peak
      ▲
   ▲
```

### Area Chart
```
         █
       █ █
       ███
     █ ███
     █████
   █ █████
```

## Integration with Other Examples

These examples complement the [viz CLI tool](../../cli/examples/viz/):

```bash
# Compare example output with CLI tool
go run dual_rendering.go > /tmp/example.svg
cd ../../cli/examples/viz
go run main.go -type line-graph -format svg -data data/linegraph_smooth.json > /tmp/cli.svg
```

## Extending the Examples

Each example is self-contained and easy to modify:

1. **Change data**: Edit the `points` arrays
2. **Try different themes**: Use `design.MidnightTheme()` or `design.NordTheme()`
3. **Adjust sizing**: Modify `bounds` values
4. **Mix renderers**: Create both SVG and Terminal output
5. **Add more points**: Test with larger datasets

## Dependencies

All examples require:
- `github.com/SCKelemen/dataviz`
- `github.com/SCKelemen/design-system`

No additional setup needed - just `go run`.

## Next Steps

After exploring these examples:

1. Check out the [main viz CLI tool](../../cli/examples/viz/) for a complete application
2. Read the [dataviz README](../README.md) for API documentation
3. Explore [design-system themes](../../design-system/) for styling options
4. See the [svg package](../../svg/) for low-level SVG generation
