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
