# Dataviz

Reusable data visualization components with dual rendering support (SVG + Terminal).

## Features

- **Multiple Visualization Types**: Heatmaps, line graphs, bar charts, stat cards
- **Dual Rendering**: Generate both SVG (for web) and terminal output (for CLI)
- **Design System Integration**: Works with SCKelemen/design-system tokens
- **Color Space Support**: OKLCH and other perceptually uniform gradients
- **Card Containers**: Reusable card wrappers with headers, legends, footers

## Installation

```bash
go get github.com/SCKelemen/dataviz
```

## Current Status

**Phase 2 In Progress** - Core data structures and types are complete. Visualization rendering functions are being extracted from repobeats.

## Visualization Types

- **Heatmaps**: Linear (30-day) and weeks (GitHub-style grid)
- **Line Graphs**: Time series with optional gradient fills
- **Bar Charts**: Vertical bars with stacking support
- **Stat Cards**: Statistics with trend mini-graphs

## Dependencies

- [github.com/SCKelemen/color](https://github.com/SCKelemen/color) - Color manipulation
- [github.com/SCKelemen/layout](https://github.com/SCKelemen/layout) - Layout engine
- [github.com/SCKelemen/render-svg](https://github.com/SCKelemen/render-svg) - SVG primitives
- [github.com/SCKelemen/design-system](https://github.com/SCKelemen/design-system) - Design tokens
- [github.com/SCKelemen/cli](https://github.com/SCKelemen/cli) - Terminal rendering

## Coming Soon

- Complete SVG renderer implementation
- Terminal renderer implementation
- Examples and documentation
- Tests
