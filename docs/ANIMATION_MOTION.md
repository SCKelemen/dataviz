# Animation and Motion Architecture

## Overview

This document defines the animation and motion strategy for dataviz. Animation is **separate from interactivity** - animations enhance visual communication and data storytelling, while interactions respond to user input.

## Core Principles

1. **Separate Animation from Interactivity**
   - Animation: Transitions, reveals, emphasis (automatic, declarative)
   - Interactivity: Hover, click, drag (user-triggered, requires JavaScript)

2. **Use Design System Motion Tokens**
   - Duration, easing, delay all come from SCKelemen/design-system
   - Consistent motion language across all visualizations

3. **SVG-Native Animations**
   - Use SVG `<animate>`, `<animateTransform>`, `<animateMotion>` elements
   - No CSS or JavaScript required (works in static contexts)
   - Graceful degradation when animation disabled

4. **Accessibility First**
   - Respect `prefers-reduced-motion` via CSS media query
   - Provide instant rendering option (no animation)
   - Never rely on animation for critical information

## Animation Types

### 1. Reveal Animations (Data Entry)

Show data progressively to guide attention and aid comprehension.

```go
// RevealAnimation defines how chart elements appear
type RevealAnimation struct {
    Type     RevealType     // Sequential, Simultaneous, Cascade
    Duration time.Duration  // From design tokens
    Easing   EasingFunction // From design tokens
    Delay    time.Duration  // Stagger delay between elements
}

type RevealType int
const (
    RevealSequential   RevealType = iota // One after another
    RevealSimultaneous                    // All at once
    RevealCascade                         // Staggered with overlap
    RevealNone                            // Instant (accessibility)
)
```

**Example: Bar Chart Reveal**
```svg
<!-- Bar grows from bottom to final height -->
<rect x="10" y="50" width="20" height="0" fill="#3b82f6">
    <animate
        attributeName="height"
        from="0"
        to="150"
        dur="0.6s"
        fill="freeze"
        calcMode="spline"
        keySplines="0.16, 1, 0.3, 1"
    />
    <animate
        attributeName="y"
        from="200"
        to="50"
        dur="0.6s"
        fill="freeze"
        calcMode="spline"
        keySplines="0.16, 1, 0.3, 1"
    />
</rect>
```

**Example: Pie Chart Reveal**
```svg
<!-- Slice draws from startAngle to endAngle -->
<path d="M..." fill="#3b82f6">
    <animate
        attributeName="d"
        from="M 200,200 L 200,100 A 0,0 0 0,1 200,100 Z"
        to="M 200,200 L 200,100 A 100,100 0 0,1 300,200 Z"
        dur="0.8s"
        fill="freeze"
        calcMode="spline"
        keySplines="0.16, 1, 0.3, 1"
    />
</path>
```

**Example: Line Chart Reveal**
```svg
<!-- Line draws from left to right using stroke-dasharray -->
<path d="M..." stroke="#3b82f6" fill="none" stroke-width="2"
      stroke-dasharray="1000" stroke-dashoffset="1000">
    <animate
        attributeName="stroke-dashoffset"
        from="1000"
        to="0"
        dur="1.2s"
        fill="freeze"
        calcMode="spline"
        keySplines="0.16, 1, 0.3, 1"
    />
</path>
```

### 2. Transition Animations (Data Updates)

Smoothly animate between data states.

```go
// TransitionAnimation defines how data changes are animated
type TransitionAnimation struct {
    Duration time.Duration  // From design tokens
    Easing   EasingFunction // From design tokens
    Property AnimProperty   // What to animate (position, size, color, opacity)
}

type AnimProperty int
const (
    AnimPosition AnimProperty = iota // x, y coordinates
    AnimSize                         // width, height
    AnimColor                        // fill, stroke
    AnimOpacity                      // fill-opacity, stroke-opacity
    AnimTransform                    // translate, scale, rotate
)
```

**Example: Bar Height Change**
```svg
<rect x="10" y="50" width="20" height="150" fill="#3b82f6">
    <!-- When data updates from 150 to 200 -->
    <animate
        attributeName="height"
        from="150"
        to="200"
        begin="dataUpdate.begin"
        dur="0.4s"
        fill="freeze"
    />
    <animate
        attributeName="y"
        from="50"
        to="0"
        begin="dataUpdate.begin"
        dur="0.4s"
        fill="freeze"
    />
</rect>
```

### 3. Emphasis Animations (Attention)

Highlight specific data points or patterns.

```go
// EmphasisAnimation draws attention to elements
type EmphasisAnimation struct {
    Type     EmphasisType   // Pulse, Glow, Bounce, Wiggle
    Duration time.Duration  // From design tokens
    Repeat   RepeatMode     // Once, Loop, Count
}

type EmphasisType int
const (
    EmphasisPulse  EmphasisType = iota // Scale up/down
    EmphasisGlow                        // Opacity pulse
    EmphasisBounce                      // Vertical bounce
    EmphasisWiggle                      // Horizontal shake
)
```

**Example: Pulse Animation**
```svg
<circle cx="100" cy="100" r="5" fill="#ef4444">
    <animateTransform
        attributeName="transform"
        type="scale"
        values="1; 1.3; 1"
        dur="0.6s"
        repeatCount="3"
        additive="sum"
    />
    <animate
        attributeName="opacity"
        values="1; 0.6; 1"
        dur="0.6s"
        repeatCount="3"
    />
</circle>
```

### 4. Path Morphing (Shape Transitions)

Animate between different chart types or configurations.

```go
// MorphAnimation animines path shape changes
type MorphAnimation struct {
    FromPath string        // Starting SVG path
    ToPath   string        // Ending SVG path
    Duration time.Duration // From design tokens
    Easing   EasingFunction
}
```

**Example: Bar to Line Chart Morph**
```svg
<path fill="none" stroke="#3b82f6">
    <animate
        attributeName="d"
        from="M 0,100 L 0,150 L 20,150 L 20,100 Z"
        to="M 0,100 L 100,80 L 200,120"
        dur="0.8s"
        fill="freeze"
    />
</path>
```

## Design System Integration

### Motion Tokens

All timing and easing values come from the design system:

```go
// From SCKelemen/design-system
type MotionTokens struct {
    // Duration tokens
    DurationInstant   time.Duration // 0ms - No animation
    DurationQuick     time.Duration // 100ms - Micro-interactions
    DurationNormal    time.Duration // 200ms - Default transitions
    DurationModerate  time.Duration // 400ms - Data reveals
    DurationSlow      time.Duration // 600ms - Complex animations
    DurationDeliberate time.Duration // 1000ms - Story-driven reveals

    // Easing tokens (cubic-bezier values)
    EaseStandard      EasingFunction // (0.4, 0.0, 0.2, 1) - Standard
    EaseDecelerate    EasingFunction // (0.0, 0.0, 0.2, 1) - Enter
    EaseAccelerate    EasingFunction // (0.4, 0.0, 1, 1)   - Exit
    EaseExpressive    EasingFunction // (0.16, 1, 0.3, 1)  - Bouncy

    // Stagger delays
    StaggerTiny       time.Duration // 25ms
    StaggerSmall      time.Duration // 50ms
    StaggerMedium     time.Duration // 100ms
    StaggerLarge      time.Duration // 200ms
}
```

### Easing Functions

Map design tokens to SVG keySplines:

```go
type EasingFunction struct {
    Name       string
    CubicBezier [4]float64
    KeySplines string // For SVG animate calcMode="spline"
}

var (
    EaseStandard = EasingFunction{
        Name:       "standard",
        CubicBezier: [4]float64{0.4, 0.0, 0.2, 1.0},
        KeySplines: "0.4 0 0.2 1",
    }

    EaseExpressive = EasingFunction{
        Name:       "expressive",
        CubicBezier: [4]float64{0.16, 1.0, 0.3, 1.0},
        KeySplines: "0.16 1 0.3 1",
    }
)
```

### Usage in Context

```go
// From SURFACE_CANVAS_ARCHITECTURE.md
type Context struct {
    DesignTokens *design.DesignTokens
    ColorSpace   color.Space
    FontMetrics  *FontMetrics

    // Motion configuration
    MotionTokens     *MotionTokens
    ReducedMotion    bool           // From prefers-reduced-motion
    AnimationEnabled bool           // Global animation toggle
}

// Helper to get appropriate duration
func (ctx *Context) GetRevealDuration() time.Duration {
    if ctx.ReducedMotion || !ctx.AnimationEnabled {
        return 0 // Instant
    }
    return ctx.MotionTokens.DurationModerate // 400ms
}
```

## Chart-Specific Animation Patterns

### Pie/Donut Charts

```go
type PieChartAnimation struct {
    RevealMode RevealType // Sequential (one slice at a time) or Cascade
    Duration   time.Duration
    StartAngle float64 // Where to start reveal (usually -90° for top)
}

// Each slice animates its path from 0 width to final width
func (p *PieChart) RenderAnimated(ctx *Context) string {
    duration := ctx.GetRevealDuration()

    for i, slice := range p.Slices {
        delay := time.Duration(i) * ctx.MotionTokens.StaggerMedium

        // Animate path from zero arc to final arc
        svg += fmt.Sprintf(`
            <path d="%s" fill="%s">
                <animate
                    attributeName="d"
                    from="%s"
                    to="%s"
                    begin="%dms"
                    dur="%dms"
                    fill="freeze"
                    calcMode="spline"
                    keySplines="%s"
                />
            </path>
        `, slice.Path, slice.Color,
           slice.ZeroPath, slice.Path,
           delay.Milliseconds(), duration.Milliseconds(),
           ctx.MotionTokens.EaseExpressive.KeySplines)
    }
}
```

### Bar Charts

```go
type BarChartAnimation struct {
    RevealMode RevealType // Sequential, Cascade, or Simultaneous
    Direction  AnimDirection // BottomUp, TopDown, LeftRight
}

// Bars grow from baseline to final height
func (b *BarChart) RenderAnimated(ctx *Context) string {
    duration := ctx.GetRevealDuration()

    for i, bar := range b.Bars {
        delay := time.Duration(i) * ctx.MotionTokens.StaggerSmall

        svg += fmt.Sprintf(`
            <rect x="%f" y="%f" width="%f" height="0" fill="%s">
                <animate
                    attributeName="height"
                    from="0"
                    to="%f"
                    begin="%dms"
                    dur="%dms"
                    fill="freeze"
                    calcMode="spline"
                    keySplines="%s"
                />
                <animate
                    attributeName="y"
                    from="%f"
                    to="%f"
                    begin="%dms"
                    dur="%dms"
                    fill="freeze"
                    calcMode="spline"
                    keySplines="%s"
                />
            </rect>
        `, bar.X, bar.Y+bar.Height, bar.Width, bar.Color,
           bar.Height,
           delay.Milliseconds(), duration.Milliseconds(),
           ctx.MotionTokens.EaseExpressive.KeySplines,
           bar.Y+bar.Height, bar.Y,
           delay.Milliseconds(), duration.Milliseconds(),
           ctx.MotionTokens.EaseExpressive.KeySplines)
    }
}
```

### Line Charts

```go
type LineChartAnimation struct {
    DrawMode DrawMode // LeftToRight, RightToLeft, FromCenter
}

// Line draws using stroke-dasharray technique
func (l *LineChart) RenderAnimated(ctx *Context) string {
    duration := ctx.GetRevealDuration()

    // Calculate path length
    pathLength := calculatePathLength(l.Path)

    svg += fmt.Sprintf(`
        <path d="%s" stroke="%s" fill="none"
              stroke-dasharray="%f"
              stroke-dashoffset="%f">
            <animate
                attributeName="stroke-dashoffset"
                from="%f"
                to="0"
                dur="%dms"
                fill="freeze"
                calcMode="spline"
                keySplines="%s"
            />
        </path>
    `, l.Path, l.Color,
       pathLength, pathLength,
       pathLength,
       duration.Milliseconds(),
       ctx.MotionTokens.EaseStandard.KeySplines)

    // Animate fill area with opacity
    if l.FillColor != "" {
        svg += fmt.Sprintf(`
            <path d="%s" fill="%s" opacity="0">
                <animate
                    attributeName="opacity"
                    from="0"
                    to="0.3"
                    begin="%dms"
                    dur="%dms"
                    fill="freeze"
                />
            </path>
        `, l.AreaPath, l.FillColor,
           duration.Milliseconds(), // Start after line finishes
           ctx.MotionTokens.DurationQuick.Milliseconds())
    }
}
```

### Area Charts

```go
// Similar to line but fill area animates with vertical clip-path
func (a *AreaChart) RenderAnimated(ctx *Context) string {
    duration := ctx.GetRevealDuration()

    // Create clip path that sweeps left to right
    clipID := fmt.Sprintf("areaClip-%d", a.ID)

    svg += fmt.Sprintf(`
        <defs>
            <clipPath id="%s">
                <rect x="0" y="0" width="0" height="%f">
                    <animate
                        attributeName="width"
                        from="0"
                        to="%f"
                        dur="%dms"
                        fill="freeze"
                        calcMode="spline"
                        keySplines="%s"
                    />
                </rect>
            </clipPath>
        </defs>
        <path d="%s" fill="%s" clip-path="url(#%s)" />
    `, clipID, a.Height,
       a.Width,
       duration.Milliseconds(),
       ctx.MotionTokens.EaseStandard.KeySplines,
       a.Path, a.FillColor, clipID)
}
```

### Scatter Plots

```go
type ScatterPlotAnimation struct {
    RevealMode RevealType // Cascade (staggered appearance)
}

// Points fade in and scale from 0
func (s *ScatterPlot) RenderAnimated(ctx *Context) string {
    duration := ctx.GetRevealDuration()

    for i, point := range s.Points {
        delay := time.Duration(i) * ctx.MotionTokens.StaggerTiny

        svg += fmt.Sprintf(`
            <circle cx="%f" cy="%f" r="%f" fill="%s" opacity="0">
                <animate
                    attributeName="opacity"
                    from="0"
                    to="1"
                    begin="%dms"
                    dur="%dms"
                    fill="freeze"
                />
                <animateTransform
                    attributeName="transform"
                    type="scale"
                    from="0 0"
                    to="1 1"
                    begin="%dms"
                    dur="%dms"
                    fill="freeze"
                    calcMode="spline"
                    keySplines="%s"
                    additive="sum"
                />
            </circle>
        `, point.X, point.Y, point.Radius, point.Color,
           delay.Milliseconds(), duration.Milliseconds(),
           delay.Milliseconds(), duration.Milliseconds(),
           ctx.MotionTokens.EaseExpressive.KeySplines)
    }
}
```

## Accessibility Considerations

### Prefers-Reduced-Motion

Respect user preferences for reduced motion:

```svg
<style>
@media (prefers-reduced-motion: reduce) {
    * {
        animation-duration: 0.01ms !important;
        animation-iteration-count: 1 !important;
        transition-duration: 0.01ms !important;
    }
}
</style>
```

In Go:

```go
func (ctx *Context) ShouldAnimate() bool {
    return ctx.AnimationEnabled && !ctx.ReducedMotion
}

func (chart *Chart) Render(ctx *Context) string {
    if ctx.ShouldAnimate() {
        return chart.RenderAnimated(ctx)
    }
    return chart.RenderStatic(ctx)
}
```

### Animation Control API

Provide explicit animation control:

```go
// AnimationConfig allows fine-grained control
type AnimationConfig struct {
    Enabled       bool
    RevealType    RevealType
    Duration      time.Duration // Override design token
    Easing        EasingFunction
    StaggerDelay  time.Duration
    RepeatCount   int // 0 = no repeat, -1 = infinite
}

// Chart-level animation config
func (chart *PieChart) WithAnimation(config AnimationConfig) *PieChart {
    chart.Animation = config
    return chart
}

// Usage
pie := NewPieChart(data).
    WithAnimation(AnimationConfig{
        Enabled:      true,
        RevealType:   RevealCascade,
        Duration:     600 * time.Millisecond,
        Easing:       EaseExpressive,
        StaggerDelay: 50 * time.Millisecond,
    })
```

## Implementation Guidelines

### 1. Always Provide Static Fallback

Every chart must render without animation:

```go
func (chart *Chart) Render(ctx *Context) string {
    if ctx.ShouldAnimate() {
        return chart.RenderAnimated(ctx)
    }
    return chart.RenderStatic(ctx)
}
```

### 2. Use Design Token Durations

Never hardcode timing values:

```go
// ❌ BAD
duration := 400 * time.Millisecond

// ✅ GOOD
duration := ctx.MotionTokens.DurationModerate
```

### 3. Coordinate Multiple Animations

Use `begin` attribute to sequence:

```svg
<!-- Bar grows first -->
<rect id="bar" ...>
    <animate id="barGrow" ... dur="0.4s" />
</rect>

<!-- Label appears after bar -->
<text ...>
    <animate begin="barGrow.end" ... dur="0.2s" />
</text>
```

### 4. Optimize Path Morphing

Only morph paths with same number of commands:

```go
// Both paths must have same structure
fromPath := "M 0,0 L 100,0 L 100,100 L 0,100 Z"
toPath   := "M 0,0 L 100,0 L 100,50 L 0,80 Z"

// ✅ GOOD - same number of commands (M + 3L + Z)
```

### 5. Test with Reduced Motion

Always test with motion preferences disabled:

```bash
# Chrome DevTools: Rendering panel
# Enable "Emulate CSS prefers-reduced-motion"

# Firefox about:config
# Set ui.prefersReducedMotion to 1
```

## Animation in Different Surfaces

### Web Surface (Full Animation Support)

```go
type WebSurface struct {
    SupportsAnimation    bool // true
    SupportsCSS          bool // true
    SupportsJavaScript   bool // true (but avoid for declarative anims)
}

func (s *WebSurface) Render(chart *Chart, ctx *Context) string {
    // Full SVG animation support
    return chart.RenderAnimated(ctx)
}
```

### GitHub README Surface (Limited Animation)

```go
type GitHubReadmeSurface struct {
    SupportsAnimation    bool // Partial (SVG animate only, no CSS)
    SupportsCSS          bool // false (stripped by sanitizer)
    SupportsJavaScript   bool // false
}

func (s *GitHubReadmeSurface) Render(chart *Chart, ctx *Context) string {
    // GitHub strips CSS and JS but keeps SVG animate
    // Use simple animations, avoid complex sequences
    ctx.MotionTokens.DurationModerate = 300 * time.Millisecond // Shorter
    return chart.RenderAnimated(ctx)
}
```

### CLI/Terminal Surface (No Animation)

```go
type CLISurface struct {
    SupportsAnimation bool // false (static output)
}

func (s *CLISurface) Render(chart *Chart, ctx *Context) string {
    // Terminal cannot animate SVG
    ctx.AnimationEnabled = false
    return chart.RenderStatic(ctx)
}
```

### Print Surface (No Animation)

```go
type PrintSurface struct {
    SupportsAnimation bool // false (static medium)
}

func (s *PrintSurface) Render(chart *Chart, ctx *Context) string {
    // Print is inherently static
    ctx.AnimationEnabled = false
    return chart.RenderStatic(ctx)
}
```

## Animation Antipatterns

### ❌ Don't: Rely on Animation for Information

```go
// BAD - critical information only shown during animation
<text opacity="1">
    <animate attributeName="opacity" to="0" dur="2s" />
</text>
```

Animation should enhance, not replace, static visual encoding.

### ❌ Don't: Animate Too Much

```go
// BAD - everything animates, overwhelming
chart.Animate(Background, Grid, Axes, Data, Labels, Legend, Title)
```

Focus animation on primary data elements only.

### ❌ Don't: Use Inconsistent Timing

```go
// BAD - arbitrary durations
bar1.Duration = 300 * time.Millisecond
bar2.Duration = 450 * time.Millisecond
bar3.Duration = 200 * time.Millisecond
```

Use design token durations consistently.

### ❌ Don't: Ignore Reduced Motion

```go
// BAD - always animates
func Render() string {
    return RenderAnimated()
}
```

Always check `prefers-reduced-motion` and provide static fallback.

## Animation Best Practices

### ✅ Do: Stagger Similar Elements

```go
for i, bar := range bars {
    delay := time.Duration(i) * ctx.MotionTokens.StaggerSmall
    bar.AnimateWithDelay(delay)
}
```

Creates visual rhythm, guides eye movement.

### ✅ Do: Use Expressive Easing

```go
// Data reveals feel more natural with bouncy easing
ctx.MotionTokens.EaseExpressive // 0.16, 1, 0.3, 1
```

### ✅ Do: Coordinate Reveal Order

```go
// Reveal axes → grid → data → labels
axes.Animate(delay: 0ms)
grid.Animate(delay: 100ms)
data.Animate(delay: 200ms)
labels.Animate(delay: 600ms)
```

Build visual hierarchy through timing.

### ✅ Do: Keep Durations Short

```go
// Most animations should be 200-600ms
// Longer feels sluggish, shorter feels janky
ctx.MotionTokens.DurationModerate // 400ms - good default
```

## Future Enhancements

### Phase 1: Core Animation (v1.7.0)
- [ ] Implement `AnimationConfig` type
- [ ] Add `RenderAnimated()` methods to all charts
- [ ] Support reveal animations (sequential, cascade)
- [ ] Integrate motion tokens from design system
- [ ] Implement `prefers-reduced-motion` detection

### Phase 2: Advanced Animations (v1.8.0)
- [ ] Transition animations for data updates
- [ ] Path morphing for chart type changes
- [ ] Emphasis animations (pulse, glow)
- [ ] Coordinated multi-element sequences

### Phase 3: Interactive Animations (v2.0.0+)
- [ ] Hover-triggered animations
- [ ] Click-triggered state changes
- [ ] Drag interactions with smooth feedback
- [ ] Gesture-based animations (pinch, swipe)

Note: Interactive animations require JavaScript and are **separate** from declarative SVG animations.

## Related Documents

- [SURFACE_CANVAS_ARCHITECTURE.md](SURFACE_CANVAS_ARCHITECTURE.md) - Context provides MotionTokens
- [COLOR_STRATEGY.md](COLOR_STRATEGY.md) - Color transitions using OKLCH
- [ROADMAP.md](ROADMAP.md) - Implementation timeline
- SCKelemen/design-system: Motion tokens specification

## Success Criteria

- ✅ All charts support animation toggle (enabled/disabled)
- ✅ Animations use design system motion tokens
- ✅ Animations respect `prefers-reduced-motion`
- ✅ Static fallback always available
- ✅ Smooth 60fps animations in browser
- ✅ Animations work in GitHub README (SVG-only)
- ✅ Clear separation: Animation ≠ Interactivity
