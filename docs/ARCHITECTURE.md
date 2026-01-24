# Architecture

## Design Philosophy

This is a **general-purpose layout and rendering engine** with a high-level charting API, not just a charting library.

### Core Principles

1. **Layered Architecture**
   - Low-level: Layout engine (flexbox, grid, text)
   - Mid-level: Renderers (SVG, terminal)
   - High-level: Charts API (uses layout + renderers)

2. **Dual Output Modes**
   - Every layout can render to SVG **and** Terminal
   - Consistent results across output modes
   - PNG/JPEG export via SVG rasterization

3. **General-Purpose Rendering**
   - Not limited to charts
   - Can render any flexbox or grid layout
   - Text rendering with proper Unicode handling
   - Custom visualizations

4. **Optional Design Tokens**
   - Design system is opt-in
   - Rendering works fine without tokens
   - Tokens provide systematic styling

5. **Data-Source Agnostic**
   - Generic data types (interface{}, float64)
   - Works with data from any source
   - MCP server for AI agent integration

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                   Application Layer                             │
│  ┌──────────────┐  ┌──────────────┐  ┌────────────────────┐   │
│  │  cmd/viz-cli │  │cmd/dataviz-  │  │  User Applications │   │
│  │              │  │     mcp      │  │                    │   │
│  └──────┬───────┘  └──────┬───────┘  └─────────┬──────────┘   │
│         │                 │                     │              │
└─────────┼─────────────────┼─────────────────────┼──────────────┘
          │                 │                     │
          ▼                 ▼                     ▼
┌─────────────────────────────────────────────────────────────────┐
│                  High-Level APIs (Optional)                     │
│  ┌──────────────┐  ┌──────────────┐  ┌────────────────────┐   │
│  │   charts/    │  │     tui/     │  │       mcp/         │   │
│  │  Line, Bar,  │  │  Interactive │  │   MCP Protocol     │   │
│  │  Scatter,    │  │  Dashboard   │  │   Server for AI    │   │
│  │  Heatmap,Pie │  │  Components  │  │   Agents           │   │
│  └──────┬───────┘  └──────┬───────┘  └─────────┬──────────┘   │
│         │                 │                     │              │
│         └─────────────────┼─────────────────────┘              │
│                           │                                    │
└───────────────────────────┼────────────────────────────────────┘
                            │ uses
┌───────────────────────────▼────────────────────────────────────┐
│                   Core Rendering Engine                        │
│                                                                 │
│  ┌────────────────────────────────────────────────────────┐   │
│  │  layout/ - CSS Grid, Flexbox, Text Layout             │   │
│  │  • Renderer-agnostic layout calculation                │   │
│  │  • Type-safe CSS units (github.com/SCKelemen/units)    │   │
│  │  • Unicode-aware (github.com/SCKelemen/unicode)        │   │
│  └────────────────────────────────────────────────────────┘   │
│                           │                                    │
│                           ▼                                    │
│  ┌──────────────────┐          ┌──────────────────────────┐   │
│  │   render/svg/    │          │   render/terminal/       │   │
│  │  • SVG output    │          │   • Terminal output      │   │
│  │  • General       │          │   • Box drawing chars    │   │
│  │    purpose       │          │   • 24-bit color         │   │
│  │  • Used by       │          │   • Same layout API      │   │
│  │    charts/       │          │                          │   │
│  └──────────────────┘          └──────────────────────────┘   │
│                           │                                    │
│                           ▼                                    │
│  ┌────────────────────────────────────────────────────────┐   │
│  │  export/ - PNG/JPEG conversion from SVG                │   │
│  │  • Rasterization via oksvg + rasterx                   │   │
│  │  • Auto-dimensions from viewBox                        │   │
│  │  • Quality settings for JPEG                           │   │
│  └────────────────────────────────────────────────────────┘   │
│                                                                 │
│  ┌────────────────────────────────────────────────────────┐   │
│  │  design/ - Design Tokens (Optional)                    │   │
│  │  • Themes: midnight, nord, paper, wrapped              │   │
│  │  • Color, typography, spacing tokens                   │   │
│  │  • Radix UI integration                                │   │
│  └────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Foundation Libraries                         │
│  ┌──────────────┐  ┌──────────────┐  ┌────────────────────┐   │
│  │   unicode    │  │    color     │  │       units        │   │
│  │  10 UAX/UTS  │  │    OKLCH     │  │   Type-safe CSS    │   │
│  │              │  │              │  │                    │   │
│  └──────────────┘  └──────────────┘  └────────────────────┘   │
│  ┌──────────────┐  ┌──────────────┐                           │
│  │     svg      │  │     text     │                           │
│  │  Primitives  │  │  Operations  │                           │
│  └──────────────┘  └──────────────┘                           │
└─────────────────────────────────────────────────────────────────┘
```

## Package Dependencies

```
Foundation (external, stable):
  github.com/SCKelemen/unicode  (v1.0.1)
  github.com/SCKelemen/color    (v1.0.0)
  github.com/SCKelemen/units    (v1.0.2)
  github.com/SCKelemen/svg      (v1.0.0)
  github.com/SCKelemen/text     (v1.0.0)

Core Rendering:
  layout/          → unicode, units, text
  render/svg/      → layout, svg
  render/terminal/ → layout, text, unicode
  export/          → oksvg, rasterx
  design/          → color

High-Level APIs:
  charts/          → layout, render/svg, render/terminal, design
  tui/             → charts, design, bubbletea
  mcp/             → charts, export, mcp-sdk

Binaries:
  cmd/viz-cli/     → tui, charts
  cmd/dataviz-mcp/ → mcp
```

## Data Flow

### 1. General-Purpose Rendering

```
User creates layout
       │
       ▼
layout/ calculates positions
       │
       ├─────────────┬─────────────┐
       ▼             ▼             ▼
render/svg/   render/terminal/  (custom)
       │             │
       ▼             ▼
     SVG         Terminal
       │           output
       ▼
   export/
       │
       ▼
   PNG/JPEG
```

### 2. High-Level Charting

```
User provides data
       │
       ▼
charts/ creates layout
       │
       ▼
layout/ calculates positions
       │
       ├─────────────┬─────────────┐
       ▼             ▼             ▼
render/svg/   render/terminal/  export/
       │             │             │
       ▼             ▼             ▼
     SVG         Terminal      PNG/JPEG
```

### 3. MCP Server (AI Agent Workflow)

```
Claude Code (MCP Client)
       │
       ▼
mcp/ receives data
       │
       ▼
charts/ generates chart
       │
       ▼
render/svg/ produces SVG
       │
       ├─────────────┐
       │             ▼
       │         export/
       │             │
       ▼             ▼
     SVG         PNG/JPEG
       │             │
       └─────┬───────┘
             ▼
    Return to Claude Code
```

## Key Design Decisions

### 1. Why Layered Architecture?

**Problem**: Mixing layout logic with rendering creates tight coupling.

**Solution**: Separate concerns:
- `layout/` - What to position and where (renderer-agnostic)
- `render/*/` - How to output the positioned elements
- `charts/` - High-level API for common use cases

**Benefits**:
- Can add new renderers without touching layout
- Can use layout without charts
- Testable in isolation

### 2. Why Dual Output (SVG + Terminal)?

**Problem**: Different use cases need different outputs.

**Solution**: Single layout engine, multiple renderers.

**Benefits**:
- CLI tools get terminal output
- Web apps get SVG
- Documentation gets high-quality images (PNG/JPEG)
- Consistent results across outputs

### 3. Why General-Purpose (Not Just Charts)?

**Problem**: Charts are just one type of visualization.

**Solution**: Build general rendering engine, then add charts on top.

**Benefits**:
- Can render any layout (dashboards, infographics, diagrams)
- Charts are not special-cased
- Users can build custom visualizations
- More flexibility

### 4. Why Optional Design Tokens?

**Problem**: Not everyone needs systematic styling.

**Solution**: Make design tokens opt-in via `design/` package.

**Benefits**:
- Simple use cases don't pay for complexity
- Advanced use cases get consistency
- Themes provided but not required

### 5. Why Monorepo?

**Problem**: Three separate repos had:
- Duplicate chart implementations (~500 lines)
- Inconsistent features (dataviz had gradients, mcp didn't)
- Coordination overhead (separate releases, versions)

**Solution**: Single monorepo with multiple packages.

**Benefits**:
- Single source of truth for rendering
- Atomic commits across packages
- Consistent versioning
- Shared CI/CD
- Easier maintenance

## Rendering Details

### SVG Renderer (`render/svg/`)

The SVG renderer produces standard SVG 1.1 output:

```xml
<svg xmlns="http://www.w3.org/2000/svg" width="800" height="600" viewBox="0 0 800 600">
  <!-- Layout elements -->
  <rect x="0" y="0" width="800" height="600" fill="#ffffff"/>

  <!-- Smooth curves using Catmull-Rom splines -->
  <path d="M 100,500 C 120,480 140,460 150,450 ..." stroke="#3b82f6" fill="none"/>

  <!-- Gradients -->
  <defs>
    <linearGradient id="gradient1" x1="0%" y1="0%" x2="0%" y2="100%">
      <stop offset="0%" stop-color="#3b82f6" stop-opacity="0.8"/>
      <stop offset="100%" stop-color="#3b82f6" stop-opacity="0.2"/>
    </linearGradient>
  </defs>

  <!-- Text with proper typography -->
  <text x="400" y="50" font-family="Inter, sans-serif" font-size="24" text-anchor="middle">
    Chart Title
  </text>
</svg>
```

Features:
- Standard SVG (works in browsers, design tools)
- Smooth curves via `svg.SmoothLinePath()`
- Gradients (linear, radial)
- Proper text rendering with Unicode support
- Responsive via viewBox

### Terminal Renderer (`render/terminal/`)

The terminal renderer uses ANSI escape codes and Unicode:

```
╭─────────────────────────────────────╮
│  Chart Title                         │
├─────────────────────────────────────┤
│                                     │
│  ▲                                 │
│  │        ╭───╮                    │
│  │   ╭───╯     ╰───╮              │
│  │ ╭╯                 ╰──╮         │
│  └─────────────────────────────>   │
│    Jan   Feb   Mar   Apr           │
╰─────────────────────────────────────╯
```

Features:
- Box-drawing characters (Unicode)
- 24-bit RGB colors (true color terminals)
- Unicode-aware text (handles emoji, CJK, etc.)
- Responsive to terminal width

### Export (`export/`)

Converts SVG to raster formats:

1. **Parse SVG** using `oksvg`
2. **Rasterize** using `rasterx` (scan line converter)
3. **Encode** to PNG or JPEG

Quality:
- PNG: Lossless with configurable compression
- JPEG: Lossy with quality 1-100 (default 90)

Dimensions:
- Auto-calculate from SVG viewBox
- Custom width/height
- Maintain aspect ratio

## Testing Strategy

### Unit Tests
- Each package has its own `*_test.go` files
- Test layout calculation in isolation
- Test SVG generation without layout
- Test export with known inputs

### Integration Tests
- Test full pipeline: data → layout → render → export
- Test charts end-to-end
- Test MCP server with real requests
- Test CLI tool with real files

### Visual Tests
- Generate reference SVGs
- Compare against known-good outputs
- Human review of samples

## Performance Considerations

### Layout Engine
- O(n) layout calculation for flexbox
- O(n×m) for grid (n rows, m columns)
- Memoization for repeated layouts

### SVG Rendering
- String building (not DOM manipulation)
- Minimal allocations
- Reusable buffers

### Terminal Rendering
- Single-pass rendering
- Buffered output (no flicker)
- Efficient Unicode handling

### Export
- Streaming SVG parsing
- Parallel rasterization (future)
- Configurable quality vs size

## Future Enhancements

### Planned Features
1. **More chart types**: Sankey, treemap, network graphs
2. **Animation**: Smooth transitions for web output
3. **3D rendering**: Perspective transforms for SVG
4. **PDF export**: Direct PDF generation (no SVG intermediate)
5. **Canvas renderer**: HTML5 Canvas output mode

### Architectural Improvements
1. **Plugin system**: Custom chart types via extensions
2. **Streaming rendering**: Large datasets without loading all in memory
3. **GPU acceleration**: For export rasterization
4. **WASM build**: Run in browser or edge functions

## Related Documentation

- [README.md](../README.md) - Getting started, examples, FAQ
- [MIGRATION_GUIDE.md](MIGRATION_GUIDE.md) - Migrating from old repos
- [MCP_INTEGRATION.md](MCP_INTEGRATION.md) - Using the MCP server
- [DESIGN_TOKENS.md](DESIGN_TOKENS.md) - Design system guide

## Questions?

See FAQ in [README.md](../README.md#faq) or open an issue.
