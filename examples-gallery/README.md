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

### Area Chart Gallery (`area-gallery.svg`)
- **Simple Area**: Basic area chart with filled region
- **Different Color**: Area chart with alternative color scheme

### Stacked Area Gallery (`stacked-area-gallery.svg`)
- **Standard Stacked**: Standard stacked area chart with multiple series
- **Smooth Curves**: Stacked area with Bezier curve smoothing

### Heatmap Gallery (`heatmap-gallery.svg`)
- **Linear Heatmap**: Horizontal heatmap showing contributions over time
- **Weeks Heatmap**: GitHub-style calendar heatmap (grid of weeks)

### Stat Card Gallery (`statcard-gallery.svg`)
Demonstrates 6 different trend patterns:
- **Positive Trend**: Rising revenue with positive change indicator
- **Negative Trend**: Declining active users with negative change indicator
- **Steady Growth**: Consistent upward trend
- **Flat Trend**: Stable metrics with no significant change
- **Rising**: Increasing metric (not necessarily positive)
- **Declining**: Decreasing metric over time

### Box Plot Gallery (`boxplot-gallery.svg`)
- **Vertical Box Plot**: Standard vertical box-and-whisker plots with outliers and mean
- **Horizontal Box Plot**: Horizontal box plots with confidence interval notches

### Histogram Gallery (`histogram-gallery.svg`)
- **Count Histogram**: Frequency distribution showing counts
- **Density Histogram**: Normalized histogram showing probability density

### Violin Plot Gallery (`violin-gallery.svg`)
- **Basic Violin Plot**: Kernel density estimation with mirrored distribution
- **Violin + Box Plot**: Violin plot with embedded box plot showing quartiles and median

### Lollipop Chart Gallery (`lollipop-gallery.svg`)
- **Vertical Lollipop**: Standard vertical lollipop chart with value labels
- **Horizontal Lollipop**: Horizontal orientation with value labels

### Candlestick Gallery (`candlestick-gallery.svg`)
- **Candlestick Chart**: OHLC data visualized as candlesticks (rising in green, falling in red)
- **OHLC Chart**: Traditional OHLC bar chart representation

### Treemap Gallery (`treemap-gallery.svg`)
- **Standard Treemap**: Squarified treemap showing hierarchical data as nested rectangles
- **With Padding**: Treemap with padding between rectangles for better visual separation

### Sunburst Gallery (`sunburst-gallery.svg`)
- **Full Sunburst**: Radial partition chart showing hierarchy as concentric rings
- **With Inner Radius**: Sunburst with hollow center (donut style) for better readability

### Circle Packing Gallery (`circle-packing-gallery.svg`)
- **Standard Packing**: Hierarchical circle packing with nested circles
- **With Padding**: Circle packing with padding between circles

### Icicle Gallery (`icicle-gallery.svg`)
- **Vertical Icicle**: Top-down icicle chart showing hierarchy as stacked rectangles
- **Horizontal Icicle**: Left-to-right icicle chart orientation

### Radar Chart Gallery (`radar-gallery.svg`)
- **With Grid**: Spider/radar chart with concentric grid circles
- **Without Grid**: Clean radar chart without background grid

### Streamchart Gallery (`streamchart-gallery.svg`)
- **Center Layout**: Streamchart with centered baseline (flowing river effect)
- **Smooth Curves**: Streamchart with Bezier curve smoothing

### Ridgeline Gallery (`ridgeline-gallery.svg`)
- **Standard Ridgeline**: Joy plot showing overlapping density distributions
- **With Fill**: Ridgeline plot with filled areas under curves

## Generating Galleries

Galleries are automatically generated using the gallery generator tool:

```bash
cd /path/to/dataviz
go run ./cmd/gallery-gen
```

Or build and run:

```bash
go build ./cmd/gallery-gen
./gallery-gen
```

The generator uses a generic configuration-driven system:
- **Gallery configurations** define variants declaratively
- **Layout strategies** handle positioning (single-row, grid, vertical stack)
- **Generic renderer** generates consistent output
- Built on the SCKelemen/svg library for type-safe SVG creation

See `cmd/gallery-gen/GALLERY_SYSTEM.md` for architecture details.

## Usage

These galleries serve as visual documentation of chart capabilities and can be:
- Embedded in documentation
- Used as reference examples
- Shared with users to demonstrate available options
- Used for visual regression testing

## Implementation

Gallery generation code is organized across four files in `cmd/gallery-gen/`:

- **`main.go`**: Core utilities and entry point
- **`gallery.go`**: Generic rendering engine
- **`layouts.go`**: Layout strategy implementations
- **`configs.go`**: Gallery configuration registry

Each gallery configuration:
1. Defines sample data providers for each variant
2. Specifies a layout strategy (single-row, grid, or vertical stack)
3. Provides chart renderers that use the main library functions
4. Declares positioning offsets and labels

The generic system:
1. Calculates dimensions based on layout strategy
2. Renders background, title, and labels consistently
3. Positions variants using layout strategies
4. Wraps everything in a single SVG document

The galleries use the same rendering functions as the main library, ensuring they accurately represent actual chart output.

### Architecture Benefits

- **No code duplication**: Structural code written once
- **Easy to extend**: New galleries in ~30-50 lines
- **Type-safe**: Compiler enforces consistency
- **Maintainable**: Single source of truth for structure

See [cmd/gallery-gen/GALLERY_SYSTEM.md](../cmd/gallery-gen/GALLERY_SYSTEM.md) for complete documentation.
