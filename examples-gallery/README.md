# Chart Gallery

This directory contains gallery SVGs showcasing all the different variations and styling options for each chart type. Each gallery displays multiple versions of a chart side-by-side to demonstrate the available configuration options.

## Gallery Files

### Pie Chart Gallery (`pie-gallery.svg`)
- **Regular Pie**: Standard pie chart with percentage labels
- **Donut Chart**: Pie chart with hollow center (donut mode)
- **Custom Colors**: Pie chart with custom color scheme

### Bar Chart Gallery (`bar-gallery.svg`)
- **Simple Bars**: Basic bar chart with single data series
- **Stacked Bars**: Bar chart with stacked segments (Open/Closed)

### Line Graph Gallery (`line-gallery.svg`)
- **Simple Line**: Basic line chart
- **Smoothed**: Line chart with Bezier curve smoothing (tension = 0.3)
- **With Markers**: Line chart with circle markers at data points
- **Filled Area**: Line chart with semi-transparent area fill

### Scatter Plot Gallery (`scatter-gallery.svg`)
Demonstrates all available marker types:
- **Circle**: Standard circular markers
- **Square**: Square markers
- **Diamond**: Diamond-shaped markers
- **Triangle**: Triangular markers
- **Cross**: Cross (+) markers
- **X**: X-shaped markers

### Connected Scatter Gallery (`connected-scatter-gallery.svg`)
Demonstrates all available line styles:
- **Solid**: Standard solid line
- **Dashed**: Dashed line pattern (10,5)
- **Dotted**: Dotted line pattern (2,3)
- **Dash-Dot**: Alternating dash-dot pattern (10,5,2,5)
- **Long Dash**: Long dashed pattern (20,5)

## Generating Galleries

Galleries are automatically generated using the gallery generator tool:

```bash
go run ./cmd/gallery-gen
```

The generator uses the SCKelemen/svg library for type-safe SVG element creation and the chart rendering functions from the `charts/` package.

## Usage

These galleries serve as visual documentation of chart capabilities and can be:
- Embedded in documentation
- Used as reference examples
- Shared with users to demonstrate available options
- Used for visual regression testing

## Implementation

Gallery generation code is located in `cmd/gallery-gen/main.go`. Each gallery:
1. Creates sample data appropriate for the chart type
2. Renders multiple chart variations using different configuration options
3. Positions charts side-by-side using SVG transforms
4. Adds labels to identify each variation
5. Wraps everything in a single SVG document

The galleries use the same rendering functions as the main library, ensuring they accurately represent actual chart output.
