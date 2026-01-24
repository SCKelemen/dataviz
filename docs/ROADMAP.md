# DataViz Roadmap

## Current Status (v1.0.0)

### Implemented ✅
**Evolution:**
- Line plot ✅
- Area ✅

**Ranking:**
- Barplot (vertical, stacked) ✅

**Correlation:**
- Scatter ✅
- Heatmap (linear + weeks view) ✅

**Part of a Whole:**
- Pie chart (MCP only, needs consolidation) ⚠️

**Other:**
- Stat cards ✅

## Immediate Tasks (v1.1.0)

### Priority 1: Consolidate Existing Code
- [ ] Move pie/donut charts from `mcp/charts/` to `charts/`
- [ ] Refactor MCP server to call `charts/` package (thin wrapper)
- [ ] Eliminate duplicate chart implementations in MCP (~500 lines)
- [ ] Add pie chart to main library API
- [ ] Update tests and documentation

**Rationale:** Single source of truth, reduce maintenance burden

## Future Chart Types

### Distribution (Priority: Medium)

| Chart Type | Use Case | Complexity | Priority |
|------------|----------|------------|----------|
| **Histogram** | Frequency distribution | Low | High |
| **Boxplot** | Statistical summary (quartiles, outliers) | Medium | High |
| **Violin** | Distribution shape with density | Medium | Medium |
| **Density** | Smooth distribution curve | Medium | Medium |
| **Ridgeline** | Multiple overlapping distributions | Medium | Low |

**Implementation notes:**
- Histogram: Bar chart variant with automatic binning
- Boxplot: Requires quartile calculations (25%, 50%, 75%, min, max)
- Violin: Combine boxplot with kernel density estimation
- Density: Kernel density estimation (KDE)
- Ridgeline: Multiple density plots with Y-offset

### Correlation (Priority: High)

| Chart Type | Use Case | Complexity | Priority |
|------------|----------|------------|----------|
| **Scatter** | X/Y correlation | Low | ✅ Done |
| **Heatmap** | Matrix correlation | Low | ✅ Done |
| **Bubble** | 3D data (X, Y, size) | Low | High |
| **Connected scatter** | Path over time | Low | High |
| **Correlogram** | Multiple correlation matrices | High | Medium |
| **Density 2D** | 2D density contours | High | Low |

**Implementation notes:**
- Bubble: Extend scatter plot with size dimension
- Connected scatter: Scatter + line between points in sequence
- Correlogram: Grid of heatmaps/scatter plots
- Density 2D: Requires 2D KDE and contour calculation

### Ranking (Priority: Medium)

| Chart Type | Use Case | Complexity | Priority |
|------------|----------|------------|----------|
| **Barplot** | Compare values | Low | ✅ Done |
| **Lollipop** | Cleaner bars with stems | Low | High |
| **Spider/Radar** | Multi-dimensional comparison | Medium | Medium |
| **Parallel** | Multi-dimensional ranking | Medium | Medium |
| **Circular Barplot** | Radial bars | Medium | Low |
| **Wordcloud** | Text frequency | Medium | Low |

**Implementation notes:**
- Lollipop: Line + circle marker at end
- Spider/Radar: Polar coordinates with polygons
- Parallel: Multiple Y-axes with connected lines
- Circular Barplot: Polar bar chart
- Wordcloud: Font size by frequency, collision detection

### Part of a Whole (Priority: High)

| Chart Type | Use Case | Complexity | Priority |
|------------|----------|------------|----------|
| **Pie chart** | Proportions | Low | ⚠️ In MCP |
| **Doughnut** | Pie with center hole | Low | ⚠️ In MCP |
| **Treemap** | Hierarchical proportions | Medium | High |
| **Dendrogram** | Hierarchical clustering | High | Medium |
| **Circular packing** | Nested circles | High | Low |

**Implementation notes:**
- Treemap: Rectangle packing algorithm (squarified)
- Dendrogram: Tree layout with branch lengths
- Circular packing: Circle packing algorithm

### Evolution (Priority: High)

| Chart Type | Use Case | Complexity | Priority |
|------------|----------|------------|----------|
| **Line plot** | Time series | Low | ✅ Done |
| **Area** | Filled time series | Low | ✅ Done |
| **Stacked area** | Multi-series cumulative | Medium | High |
| **Streamchart** | Flowing stacked area | Medium | Medium |

**Implementation notes:**
- Stacked area: Calculate cumulative Y values
- Streamchart: Center baseline with symmetrical stacking

### Flow (Priority: Low)

| Chart Type | Use Case | Complexity | Priority |
|------------|----------|------------|----------|
| **Sankey** | Flow between nodes | High | High |
| **Network** | Graph visualization | High | Medium |
| **Chord diagram** | Circular relationships | High | Medium |
| **Arc diagram** | Simplified network | Medium | Low |
| **Edge bundling** | Hierarchical edge routing | High | Low |

**Implementation notes:**
- Sankey: Flow layout algorithm (energy minimization)
- Network: Force-directed layout (d3-force equivalent)
- Chord diagram: Circular layout with bezier curves
- Arc diagram: Semi-circle with arcs
- Edge bundling: Hierarchical clustering of edges

## Implementation Strategy

### Phase 1: Consolidation (v1.1.0)
**Timeline:** 1-2 weeks
**Goal:** Single source of truth for all charts

- Move pie/donut to main library
- Refactor MCP to thin wrapper
- Add comprehensive tests
- Document all chart APIs

### Phase 2: High-Priority Extensions (v1.2.0)
**Timeline:** 4-6 weeks
**Goal:** Cover most common use cases

**Add:**
1. **Histogram** - Frequency distribution (1 week)
2. **Boxplot** - Statistical summaries (1 week)
3. **Bubble chart** - 3D scatter (1 week)
4. **Lollipop chart** - Clean ranking (1 week)
5. **Treemap** - Hierarchical proportions (2 weeks)

**Rationale:**
- Histogram/Boxplot: Statistical analysis (common in data science)
- Bubble: 3D data visualization (extend scatter)
- Lollipop: Modern alternative to bars
- Treemap: Hierarchical data (disk usage, portfolios)

### Phase 3: Advanced Charts (v1.3.0)
**Timeline:** 6-8 weeks
**Goal:** Specialized visualizations

**Add:**
1. **Stacked area** - Multi-series evolution (1 week)
2. **Connected scatter** - Path visualization (1 week)
3. **Violin plot** - Distribution with density (2 weeks)
4. **Spider/Radar** - Multi-dimensional comparison (2 weeks)
5. **Sankey diagram** - Flow visualization (2-3 weeks)

**Rationale:**
- Stacked area: Financial data, resource usage
- Connected scatter: Movement over time
- Violin: Distribution analysis (statistics, ML)
- Spider/Radar: Skills, attributes comparison
- Sankey: Energy flows, user journeys

### Phase 4: Specialized Charts (v2.0.0)
**Timeline:** 8-12 weeks
**Goal:** Comprehensive visualization library

**Add remaining charts:**
- Network graphs
- Correlogram
- Ridgeline
- Streamchart
- Density 2D
- Circular layouts
- Wordcloud
- Edge bundling

**Rationale:** Complete coverage of D3.js chart types

## Architecture: Observable Plot-Style Design

### Motivation
The current library has a split between:
- **Main library**: Time-series focused (time.Time → pixels)
- **MCP server**: Generic coordinates (float64 → pixels)

This creates type mismatches and code duplication. Adopting Observable Plot's architecture unifies both approaches.

### Core Components

#### 1. Scales (Priority: High)
**Purpose:** Map data values to visual coordinates

**Scale Types:**
- **Linear** - Continuous numeric mapping (e.g., 0-100 → 0-400px)
- **Time** - Temporal mapping (time.Time → pixels)
- **Log** - Logarithmic scale for exponential data
- **Pow** - Power scale with configurable exponent
- **Sqrt** - Square root scale
- **Ordinal** - Discrete categories (e.g., ["A", "B", "C"] → [0, 100, 200])
- **Band** - Ordinal with bandwidth (for bar charts)
- **Point** - Ordinal without padding

**API Design:**
```go
type Scale interface {
    Domain() []interface{}      // Input range
    Range() []float64           // Output range
    Apply(value interface{}) float64
    Invert(position float64) interface{}
}

// Example usage
xScale := scales.NewTimeScale(
    scales.Domain(startTime, endTime),
    scales.Range(0, 800),
)

yScale := scales.NewLinearScale(
    scales.Domain(0, 100),
    scales.Range(400, 0), // Inverted for SVG coordinates
)
```

**Implementation Notes:**
- Scales are pure functions (immutable)
- Support clamping to prevent out-of-bounds
- Support nice ticks for axis generation
- Color scales for heatmaps

#### 2. Marks (Priority: High)
**Purpose:** Visual primitives that consume scales

**Mark Types:**
- **Dot** - Points with markers (circle, square, diamond, etc.)
- **Line** - Connected path (straight or smooth)
- **Area** - Filled region between line and baseline
- **Bar** - Rectangles (vertical or horizontal)
- **Rect** - Generic rectangles for heatmaps, treemaps
- **Text** - Labels and annotations
- **Rule** - Lines for axes, grids, reference lines
- **Link** - Connections between points (for networks)
- **Arrow** - Directional indicators

**API Design:**
```go
type Mark interface {
    Render(data []interface{}, xScale, yScale Scale) string
}

// Example: Compose marks with scales
chart := plot.New(
    plot.Width(800),
    plot.Height(400),
    plot.X(scales.Time(data, "date"), scales.Range(0, 800)),
    plot.Y(scales.Linear(data, "value"), scales.Range(400, 0)),
    plot.Marks(
        marks.Line(data, marks.Stroke("#3b82f6"), marks.StrokeWidth(2)),
        marks.Dot(data, marks.Fill("#3b82f6"), marks.R(3)),
    ),
)
```

**Rationale:**
- Marks are data-agnostic (scales handle data types)
- Compose multiple marks for rich visualizations
- Reusable across chart types

#### 3. Projections (Priority: Medium)
**Purpose:** Transform coordinate systems

**Projection Types:**
- **Geographic**: Mercator, Albers USA, Orthographic, etc.
- **Polar**: Polar coordinates for radar/spider charts
- **Cartesian**: Standard X/Y (default, identity transform)
- **3D**: Perspective projection for 3D visualizations (future)

**API Design:**
```go
type Projection interface {
    Project(x, y float64) (px, py float64)
    Invert(px, py float64) (x, y float64)
}

// Geographic projection
proj := projections.NewMercator(
    projections.Center(-98, 38),
    projections.Scale(1000),
)

// Polar projection for radar chart
proj := projections.NewPolar(
    projections.Origin(400, 400),
    projections.Radius(200),
)
```

**Use Cases:**
- Geographic maps (choropleth, bubble maps)
- Radar/spider charts (polar coordinates)
- Circular layouts (chord diagrams, sunburst)
- 3D scatter plots (future)

#### 4. Transforms (Priority: High)
**Purpose:** Data transformations before rendering

**Transform Types:**
- **Bin** - Histogram binning
- **Group** - Group by category
- **Stack** - Stacked layouts (area, bar)
- **Normalize** - Normalize to 100%
- **WindowY** - Rolling window (moving average, etc.)
- **SelectFirst/Last** - Filter to first/last of group
- **Filter** - Conditional filtering
- **Sort** - Order data
- **Smooth** - LOESS, regression
- **Hexbin** - 2D binning for dense scatter plots

**API Design:**
```go
type Transform interface {
    Apply(data []interface{}) []interface{}
}

// Example: Histogram with binning
chart := plot.New(
    plot.Data(rawData),
    plot.Transform(
        transforms.BinX(transforms.Thresholds(20)),
        transforms.GroupY(transforms.Count()),
    ),
    plot.X(scales.Linear(), scales.Range(0, 800)),
    plot.Y(scales.Linear(), scales.Range(400, 0)),
    plot.Marks(marks.Bar()),
)

// Example: Stacked area chart
chart := plot.New(
    plot.Data(timeSeriesData),
    plot.Transform(transforms.StackY()),
    plot.Marks(marks.Area()),
)
```

**Rationale:**
- Keep marks simple, move complexity to transforms
- Composable data pipeline
- Statistical operations before visualization

#### 5. Legends (Priority: High)
**Purpose:** Consistent legend rendering across charts

**Current State:**
- Scattered implementations (pie chart has inline legend, line graphs have none)
- No unified positioning or styling
- Inconsistent APIs

**Proposed API:**
```go
type Legend struct {
    Position  LegendPosition // TopLeft, TopRight, BottomLeft, BottomRight, Right, Left
    Items     []LegendItem
    Style     LegendStyle
    Padding   int
}

type LegendItem struct {
    Label  string
    Symbol Symbol // Color swatch, marker, line sample
    Value  string // Optional value display
}

type LegendPosition int
const (
    LegendTopLeft LegendPosition = iota
    LegendTopRight
    LegendBottomLeft
    LegendBottomRight
    LegendRight  // Vertical, aligned right
    LegendLeft   // Vertical, aligned left
    LegendBottom // Horizontal, centered bottom
)

// Usage
chart := plot.New(
    plot.Marks(
        marks.Line(data1, marks.Stroke("#3b82f6")),
        marks.Line(data2, marks.Stroke("#10b981")),
    ),
    plot.Legend(
        legend.Position(LegendTopRight),
        legend.Items(
            legend.Item("Series 1", legend.ColorSwatch("#3b82f6")),
            legend.Item("Series 2", legend.ColorSwatch("#10b981")),
        ),
    ),
)
```

**Features:**
- Automatic legend generation from marks
- Manual legend customization
- Interactive legends (toggle visibility - future)
- Export-friendly (SVG groups with IDs)

**Current Legend Implementations:**
- **Pie chart**: Vertical list with color swatches, positioned at bottom
- **Line graphs**: No legend (TODO)
- **Stat cards**: Header legend items with manual X positioning
- **MCP line chart**: Horizontal legend at top

**Migration Path:**
1. Extract common legend rendering to `charts/legend.go`
2. Refactor existing charts to use unified legend API
3. Add legend support to charts that lack it
4. Support automatic legend generation from data

#### 6. Curves (Priority: Medium)
**Purpose:** Line interpolation methods

**Current State:**
- Basic smooth curves using Bezier (tension: 0-1)
- Linear interpolation (default)

**Curve Types:**
- **Linear** - Straight lines between points
- **Step** - Step function (step-before, step-after, step)
- **Basis** - B-spline (smooth, doesn't pass through points)
- **Cardinal** - Cardinal spline (passes through points, configurable tension)
- **Catmull-Rom** - Catmull-Rom spline (current "smooth" implementation)
- **Monotone** - Monotone cubic interpolation (prevents overshoot)
- **Natural** - Natural cubic spline
- **Bump** - Bump curve for sankey/flow diagrams

**API Design:**
```go
// Current API (maintain for compatibility)
data := LineGraphData{
    Smooth:  true,
    Tension: 0.3,
}

// New unified curves API
data := plot.Line(
    plot.Data(points),
    plot.Curve(curves.CatmullRom(0.3)), // or curves.Monotone(), curves.Step(), etc.
)
```

**Rationale:**
- Different use cases need different interpolation
- Monotone curves prevent overshoot for financial data
- Step curves for discrete data
- Basis curves for aesthetic smoothing

#### 7. Formats (Priority: Medium)
**Purpose:** Format numbers and dates for axes, labels, tooltips

**Format Types:**
- **Number**: SI prefix (1K, 1M), fixed decimals, percentage, currency
- **Date**: Locale-aware formatting (RFC3339, custom patterns)
- **Duration**: Time spans (1h 30m)
- **Bytes**: File sizes (1.5 MB, 2.3 GB)

**API Design:**
```go
type Formatter interface {
    Format(value interface{}) string
}

// Usage in axes
xAxis := axis.Bottom(
    axis.Scale(xScale),
    axis.Format(formats.Date("Jan 2")),
)

yAxis := axis.Left(
    axis.Scale(yScale),
    axis.Format(formats.Number(formats.Precision(1), formats.SI())),
)

// Examples
formats.Number(formats.Precision(2))        // "1.23"
formats.Number(formats.SI())                // "1.2K", "3.4M"
formats.Percentage(formats.Precision(0))    // "45%"
formats.Currency("$", formats.Precision(2)) // "$1.23"
formats.Date("Jan 2, 2006")                 // "Jan 15, 2024"
```

**Current State:**
- Hardcoded number formatting in chart implementations
- No date formatting options
- Inconsistent precision across charts

#### 8. Intervals (Priority: Medium)
**Purpose:** Confidence intervals, error bars, ranges

**Interval Types:**
- **Error bars** - Vertical/horizontal bars for uncertainty
- **Confidence bands** - Shaded regions around lines
- **Range bars** - Min/max ranges
- **Quantile bands** - Quartile ranges (for boxplots)

**API Design:**
```go
// Error bars
plot.Marks(
    marks.ErrorY(data,
        marks.Y("mean"),
        marks.Y1("lower"),
        marks.Y2("upper"),
        marks.Stroke("#999"),
    ),
)

// Confidence band
plot.Marks(
    marks.Area(data,
        marks.Y1("lower95"),
        marks.Y2("upper95"),
        marks.Fill("#3b82f6"),
        marks.Opacity(0.2),
    ),
    marks.Line(data, marks.Y("mean")),
)
```

**Use Cases:**
- Statistical charts (regression with confidence bands)
- Scientific plots (measurement error)
- Financial charts (high/low/close)
- Weather data (temperature ranges)

#### 9. Markers (Priority: Low - Already Implemented)
**Purpose:** Point symbols for scatter plots, line graphs

**Current State:** ✅ **Already implemented**
- Circle, square, diamond, triangle, cross, x, dot
- Configurable size
- Used in line graphs and scatter plots

**Enhancement Opportunities:**
- Custom SVG markers (stars, arrows, etc.)
- Marker orientation (for vector fields)
- Marker scaling by data value (bubble charts)
- Marker images (for pictorial charts)

#### 10. Interactions (Priority: Low - Future)
**Purpose:** Mouse events, zoom, pan, tooltips

**Interaction Types:**
- **Tooltips** - Hover to show data values
- **Crosshair** - Follow cursor across chart
- **Zoom** - Mouse wheel or drag to zoom
- **Pan** - Drag to pan
- **Brush** - Select region
- **Click** - Select points/bars
- **Hover** - Highlight on hover

**Note:** Interactions are web-only (not applicable to SVG output without JavaScript)

**Potential Approaches:**
1. **Static SVG + JS**: Generate SVG with data attributes, add interactivity client-side
2. **WASM**: Compile Go to WASM for browser-native interactivity
3. **Terminal**: Limited interactivity via bubbletea (already possible with tui package)

**Priority:** Low for now - focus on static visualization first

#### 11. Facets (Priority: Medium)
**Purpose:** Small multiples, trellis plots

**Facet Types:**
- **Grid** - Rows and columns by category
- **Wrap** - Flow layout with wrapping
- **Vertical** - Stack vertically
- **Horizontal** - Stack horizontally

**API Design:**
```go
chart := plot.New(
    plot.Data(data),
    plot.Facet(
        facet.Grid(
            facet.Row("region"),    // Facet rows by region
            facet.Col("category"),  // Facet columns by category
        ),
    ),
    plot.Marks(marks.Bar()),
)

// Or facet wrap
chart := plot.New(
    plot.Facet(facet.Wrap("region", facet.Cols(3))),
)
```

**Use Cases:**
- Compare trends across categories
- Multi-dimensional analysis
- Geographic breakdown
- A/B test comparisons

**Implementation Notes:**
- Requires layout engine to position sub-charts
- Shared or independent axes
- Shared or independent scales
- Automatic title generation

#### 12. Charts (Priority: High - Already Implemented)
**Purpose:** High-level convenience APIs for common chart types

**Current State:** ✅ **Already implemented**
- Line, Area, Bar, Scatter, Heatmap, Pie/Donut, Stat Cards

**Future:** Built on top of scales/marks/transforms
- Charts become compositions of marks
- Example: Bar chart = Bar mark + Linear scales + Ordinal x-scale
- Example: Histogram = Bin transform + Bar mark + Linear scales

**Rationale:**
- Keep high-level API for common use cases
- Power users can compose custom charts from marks
- Backwards compatibility with current API

### Implementation Phases

**Phase 1: Core Scales and Marks (v1.5.0)** - 6-8 weeks
- Implement basic scales (Linear, Time, Ordinal, Band)
- Implement core marks (Dot, Line, Area, Bar, Rect, Text)
- Refactor existing charts to use scales/marks internally
- Maintain backward compatibility with current API
- **Components:** Scales, Marks, Markers (already done)

**Phase 2: Transforms, Legends, Formats (v1.6.0)** - 4-6 weeks
- Implement key transforms (Bin, Stack, Group, Filter)
- Unified legend API across all charts
- Number and date formatters for axes
- Refactor all charts to use unified legends
- Add automatic legend generation
- **Components:** Transforms, Legends, Formats

**Phase 3: Curves and Intervals (v1.7.0)** - 4-6 weeks
- Implement curve interpolation methods (Monotone, Step, Basis, Cardinal)
- Add interval support (error bars, confidence bands)
- Enhance line and area charts with curve options
- Add statistical chart types using intervals
- **Components:** Curves, Intervals

**Phase 4: Projections and Facets (v1.8.0)** - 6-8 weeks
- Polar projection for radar/spider charts
- Geographic projections (Mercator, Albers)
- Add geographic chart types (choropleth, bubble maps)
- Implement faceting (grid, wrap, stack)
- Small multiples support
- **Components:** Projections, Facets

**Phase 5: Advanced Features (v2.0.0)** - 8-10 weeks
- Advanced transforms (Smooth/LOESS, Hexbin, Normalize, WindowY)
- Advanced scales (Log, Pow, Sqrt, Color)
- Advanced chart types (density plots, contour plots)
- Performance optimizations (stream processing, sampling)
- Comprehensive examples and documentation
- **Components:** Advanced Transforms, Advanced Scales

**Phase 6: Interactions (v2.1.0+)** - Future
- Tooltip generation (SVG with data attributes)
- WASM-based interactivity for web
- Terminal interactions via bubbletea
- **Components:** Interactions (web-focused)

**Component Priority Summary:**
| Component | Priority | Status | Phase |
|-----------|----------|--------|-------|
| Markers | Low | ✅ Implemented | v1.0.0 |
| Charts | High | ✅ Implemented | v1.0.0 |
| Scales | High | Planned | v1.5.0 |
| Marks | High | Planned | v1.5.0 |
| Transforms | High | Planned | v1.6.0 |
| Legends | High | Planned | v1.6.0 |
| Formats | Medium | Planned | v1.6.0 |
| Curves | Medium | Planned | v1.7.0 |
| Intervals | Medium | Planned | v1.7.0 |
| Projections | Medium | Planned | v1.8.0 |
| Facets | Medium | Planned | v1.8.0 |
| Interactions | Low | Future | v2.1.0+ |

**Success Criteria:**
- ✅ Unified API for time-series and generic data
- ✅ Zero code duplication between main library and MCP
- ✅ All existing charts refactored to use scales/marks
- ✅ New charts can be composed from marks
- ✅ Consistent legends across all chart types
- ✅ Flexible data transformations via transforms
- ✅ Geographic visualization support
- ✅ Small multiples via faceting
- ✅ Statistical accuracy (intervals, curves)

**Inspiration:**
- Observable Plot: https://observablehq.com/plot
- D3.js scales: https://github.com/d3/d3-scale
- Vega-Lite: https://vega.github.io/vega-lite/
- Grammar of Graphics: Wilkinson (1999)

## Design Principles

### 1. Consistent API
All charts follow the same pattern:
```go
type ChartData struct {
    // Chart-specific data
}

type ChartConfig struct {
    Width  int
    Height int
    Title  string
    // Chart-specific config
}

func RenderChart(data ChartData, config ChartConfig) string {
    // SVG output
}

func RenderChartTerminal(data ChartData, config ChartConfig) string {
    // Terminal output
}
```

### 2. Dual Output (Where Feasible)
Most charts support both outputs:
- **SVG** - Web, documentation, high-quality (all charts)
- **Terminal** - CLI, SSH, logs (simple charts only)

**Terminal-friendly charts:**
- Line, Area, Bar, Histogram, Boxplot, Scatter (basic), Heatmap, Lollipop

**SVG-only charts (too complex for terminal):**
- Pie/Donut, Treemap, Spider/Radar, Sankey, Network, Chord, Dendrogram, Circular layouts, Wordcloud, Flow diagrams

**Rationale:** Terminal has limited resolution and character-based rendering. Complex spatial layouts, circular arrangements, and graph networks don't translate well to text.

### 3. Statistical Accuracy
Charts with statistical components (boxplot, violin, density) should:
- Use established algorithms (Tukey fences, Silverman's rule)
- Match R/Python output
- Include tests against known datasets

### 4. Performance
- Efficient algorithms (O(n log n) or better)
- Stream processing for large datasets
- Configurable detail level

### 5. Accessibility
- Color-blind safe palettes
- High contrast terminal mode
- Screen reader friendly SVG (ARIA labels)

## Chart Priority Matrix

```
High Impact + Easy Implementation:
├── Histogram ⭐⭐⭐
├── Bubble chart ⭐⭐⭐
├── Lollipop ⭐⭐⭐
├── Stacked area ⭐⭐⭐
└── Connected scatter ⭐⭐⭐

High Impact + Medium Implementation:
├── Boxplot ⭐⭐
├── Treemap ⭐⭐
├── Violin ⭐⭐
├── Spider/Radar ⭐⭐
└── Sankey ⭐⭐

Medium Impact:
├── Streamchart ⭐
├── Correlogram ⭐
├── Ridgeline ⭐
├── Network graph ⭐
└── Wordcloud ⭐

Low Impact (specialized):
├── Circular packing
├── Dendrogram
├── Chord diagram
├── Arc diagram
├── Edge bundling
└── Density 2D
```

## Dependencies

### Required for Advanced Charts
- **Statistics**: Histogram binning, KDE, quartiles
  - Could use: gonum.org/v1/gonum/stat
  - Or implement: Custom lightweight stats package

- **Layout algorithms**: Treemap, force-directed, Sankey
  - Implement: Custom algorithms (refer to d3 implementations)

- **Geometry**: Convex hulls, voronoi, delaunay
  - Could use: github.com/paulmach/go.geo
  - Or implement: Custom for specific needs

**Decision:** Prefer custom lightweight implementations to avoid heavy dependencies.

## Testing Strategy

### Unit Tests
- Each chart type has test file
- Test with known datasets
- Verify SVG structure
- Test edge cases (empty, single point, etc.)

### Visual Regression Tests
- Snapshot testing for SVG output
- Compare against reference images
- Use in CI/CD

### Statistical Validation
- Compare boxplot output with R's boxplot()
- Compare histogram bins with numpy.histogram()
- Validate density curves against scipy

### Performance Benchmarks
- Benchmark each chart with 100, 1K, 10K, 100K points
- Track memory usage
- Optimize hot paths

## Documentation

### For Each Chart
- **Overview**: What is it, when to use it
- **Data format**: Expected input structure
- **Configuration**: All available options
- **Examples**: SVG and Terminal output
- **Algorithm**: Statistical methods used
- **Limitations**: What it can't do

### Tutorials
- Getting started with common charts
- Advanced features (smooth curves, markers, gradients)
- Using design tokens
- Terminal rendering
- MCP integration

## Success Metrics

**v1.1.0:** MCP consolidated, single source of truth
**v1.2.0:** 5 new high-priority charts, 90%+ test coverage
**v1.3.0:** 10+ total chart types, comprehensive docs
**v2.0.0:** 30+ chart types, D3.js feature parity

**Community:**
- 100+ GitHub stars
- 10+ external contributors
- Package used in production

**Quality:**
- 90%+ test coverage
- Zero P0 bugs
- All charts statistically validated

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for:
- How to add a new chart type
- Code style guidelines
- Testing requirements
- PR process

## Questions?

Open an issue with the `question` or `feature-request` label.
