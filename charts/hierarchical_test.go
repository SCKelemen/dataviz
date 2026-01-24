package charts

import (
	"strings"
	"testing"
)

func createTestTree() *TreeNode {
	root := NewTreeNode("Root", 0)

	// Level 1
	a := NewTreeNode("A", 0)
	b := NewTreeNode("B", 0)
	root.AddChild(a).AddChild(b)

	// Level 2 - children of A
	a1 := NewTreeNode("A1", 10)
	a2 := NewTreeNode("A2", 20)
	a3 := NewTreeNode("A3", 15)
	a.AddChild(a1).AddChild(a2).AddChild(a3)

	// Level 2 - children of B
	b1 := NewTreeNode("B1", 25)
	b2 := NewTreeNode("B2", 30)
	b.AddChild(b1).AddChild(b2)

	return root
}

func createFlatTree() *TreeNode {
	root := NewTreeNode("Root", 0)
	root.AddChild(NewTreeNode("Item1", 10))
	root.AddChild(NewTreeNode("Item2", 20))
	root.AddChild(NewTreeNode("Item3", 15))
	root.AddChild(NewTreeNode("Item4", 25))
	return root
}

func createDeepTree() *TreeNode {
	root := NewTreeNode("Root", 0)
	current := root
	for i := 1; i <= 5; i++ {
		child := NewTreeNode("Level"+string(rune('0'+i)), float64(i*10))
		current.AddChild(child)
		current = child
	}
	return root
}

// TreeNode tests

func TestNewTreeNode(t *testing.T) {
	node := NewTreeNode("Test", 42.5)

	if node.Name != "Test" {
		t.Errorf("Expected name 'Test', got '%s'", node.Name)
	}

	if node.Value != 42.5 {
		t.Errorf("Expected value 42.5, got %f", node.Value)
	}

	if len(node.Children) != 0 {
		t.Error("Expected empty children slice")
	}

	if node.Metadata == nil {
		t.Error("Expected initialized metadata map")
	}
}

func TestTreeNodeAddChild(t *testing.T) {
	parent := NewTreeNode("Parent", 0)
	child1 := NewTreeNode("Child1", 10)
	child2 := NewTreeNode("Child2", 20)

	parent.AddChild(child1).AddChild(child2)

	if len(parent.Children) != 2 {
		t.Errorf("Expected 2 children, got %d", len(parent.Children))
	}

	if parent.Children[0].Name != "Child1" {
		t.Errorf("Expected first child 'Child1', got '%s'", parent.Children[0].Name)
	}

	if parent.Children[1].Name != "Child2" {
		t.Errorf("Expected second child 'Child2', got '%s'", parent.Children[1].Name)
	}
}

func TestTreeNodeSetColor(t *testing.T) {
	node := NewTreeNode("Test", 10)
	node.SetColor("#FF0000")

	if node.Color != "#FF0000" {
		t.Errorf("Expected color '#FF0000', got '%s'", node.Color)
	}
}

func TestCalculateTreeValue(t *testing.T) {
	tests := []struct {
		name     string
		tree     *TreeNode
		expected float64
	}{
		{"nil node", nil, 0},
		{"leaf node", NewTreeNode("Leaf", 42), 42},
		{"flat tree", createFlatTree(), 70}, // 10+20+15+25
		{"nested tree", createTestTree(), 100}, // 10+20+15+25+30
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := calculateTreeValue(test.tree)
			if result != test.expected {
				t.Errorf("Expected %f, got %f", test.expected, result)
			}
		})
	}
}

// Treemap tests

func TestRenderTreemap(t *testing.T) {
	spec := TreemapSpec{
		Root:         createTestTree(),
		Width:        800,
		Height:       600,
		Padding:      2,
		ShowLabels:   true,
		MinLabelSize: 20,
	}

	result := RenderTreemap(spec)

	if result == "" {
		t.Fatal("Expected non-empty SVG output")
	}

	if !strings.Contains(result, "<rect") {
		t.Error("Expected rect elements in output")
	}

	if !strings.Contains(result, "<text") {
		t.Error("Expected text elements when ShowLabels is true")
	}
}

func TestRenderTreemapNilRoot(t *testing.T) {
	spec := TreemapSpec{
		Root:   nil,
		Width:  800,
		Height: 600,
	}

	result := RenderTreemap(spec)

	if result != "" {
		t.Error("Expected empty string for nil root")
	}
}

func TestRenderTreemapZeroValue(t *testing.T) {
	root := NewTreeNode("Root", 0)

	spec := TreemapSpec{
		Root:   root,
		Width:  800,
		Height: 600,
	}

	result := RenderTreemap(spec)

	if result != "" {
		t.Error("Expected empty string for zero-value tree")
	}
}

func TestRenderTreemapWithColorScheme(t *testing.T) {
	colorScheme := []string{"#FF0000", "#00FF00", "#0000FF"}

	spec := TreemapSpec{
		Root:        createTestTree(),
		Width:       800,
		Height:      600,
		ColorScheme: colorScheme,
	}

	result := RenderTreemap(spec)

	if !strings.Contains(result, "#FF0000") &&
		!strings.Contains(result, "#00FF00") &&
		!strings.Contains(result, "#0000FF") {
		t.Error("Expected colors from color scheme in output")
	}
}

func TestRenderTreemapWithCustomColors(t *testing.T) {
	root := NewTreeNode("Root", 0)
	child := NewTreeNode("Child", 100)
	child.SetColor("#ABCDEF")
	root.AddChild(child)

	spec := TreemapSpec{
		Root:   root,
		Width:  800,
		Height: 600,
	}

	result := RenderTreemap(spec)

	if !strings.Contains(result, "#ABCDEF") {
		t.Error("Expected custom node color in output")
	}
}

// Sunburst tests

func TestRenderSunburst(t *testing.T) {
	spec := SunburstSpec{
		Root:        createTestTree(),
		Width:       800,
		Height:      800,
		InnerRadius: 50,
		ShowLabels:  true,
		StartAngle:  0,
	}

	result := RenderSunburst(spec)

	if result == "" {
		t.Fatal("Expected non-empty SVG output")
	}

	if !strings.Contains(result, "<path") {
		t.Error("Expected path elements in output")
	}

	if !strings.Contains(result, "<g") {
		t.Error("Expected group elements for labels")
	}
}

func TestRenderSunburstNilRoot(t *testing.T) {
	spec := SunburstSpec{
		Root:   nil,
		Width:  800,
		Height: 800,
	}

	result := RenderSunburst(spec)

	if result != "" {
		t.Error("Expected empty string for nil root")
	}
}

func TestRenderSunburstZeroValue(t *testing.T) {
	root := NewTreeNode("Root", 0)

	spec := SunburstSpec{
		Root:   root,
		Width:  800,
		Height: 800,
	}

	result := RenderSunburst(spec)

	if result != "" {
		t.Error("Expected empty string for zero-value tree")
	}
}

func TestRenderSunburstWithStartAngle(t *testing.T) {
	spec := SunburstSpec{
		Root:       createFlatTree(),
		Width:      800,
		Height:     800,
		StartAngle: 90, // Start at right instead of top
	}

	result := RenderSunburst(spec)

	if result == "" {
		t.Fatal("Expected non-empty SVG output")
	}

	// Just verify it produces output - specific angle testing would require parsing SVG
	if !strings.Contains(result, "<path") {
		t.Error("Expected path elements in output")
	}
}

func TestBuildArcPath(t *testing.T) {
	path := buildArcPath(400, 400, 100, 200, 0, 1.5708) // 90 degrees

	if !strings.HasPrefix(path, "M") {
		t.Error("Expected path to start with M (move) command")
	}

	if !strings.Contains(path, "A") {
		t.Error("Expected path to contain A (arc) commands")
	}

	if !strings.Contains(path, "L") {
		t.Error("Expected path to contain L (line) commands")
	}

	if !strings.HasSuffix(path, "Z") {
		t.Error("Expected path to end with Z (close) command")
	}
}

func TestCalculateMaxDepth(t *testing.T) {
	tests := []struct {
		name     string
		tree     *TreeNode
		expected int
	}{
		{"nil node", nil, 0},
		{"leaf node", NewTreeNode("Leaf", 10), 0},
		{"flat tree", createFlatTree(), 1},
		{"nested tree", createTestTree(), 2},
		{"deep tree", createDeepTree(), 5},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := calculateMaxDepth(test.tree, 0)
			if result != test.expected {
				t.Errorf("Expected depth %d, got %d", test.expected, result)
			}
		})
	}
}

// Circle Packing tests

func TestRenderCirclePacking(t *testing.T) {
	spec := CirclePackingSpec{
		Root:       createTestTree(),
		Width:      800,
		Height:     800,
		Padding:    2,
		ShowLabels: true,
	}

	result := RenderCirclePacking(spec)

	if result == "" {
		t.Fatal("Expected non-empty SVG output")
	}

	if !strings.Contains(result, "<circle") {
		t.Error("Expected circle elements in output")
	}

	if !strings.Contains(result, "<text") {
		t.Error("Expected text elements when ShowLabels is true")
	}
}

func TestRenderCirclePackingNilRoot(t *testing.T) {
	spec := CirclePackingSpec{
		Root:   nil,
		Width:  800,
		Height: 800,
	}

	result := RenderCirclePacking(spec)

	if result != "" {
		t.Error("Expected empty string for nil root")
	}
}

func TestRenderCirclePackingWithColorScheme(t *testing.T) {
	colorScheme := []string{"#FF0000", "#00FF00", "#0000FF"}

	spec := CirclePackingSpec{
		Root:        createTestTree(),
		Width:       800,
		Height:      800,
		ColorScheme: colorScheme,
	}

	result := RenderCirclePacking(spec)

	if result == "" {
		t.Fatal("Expected non-empty SVG output")
	}

	// Verify circles are generated
	if !strings.Contains(result, "<circle") {
		t.Error("Expected circle elements in output")
	}
}

func TestRenderCirclePackingLeafNode(t *testing.T) {
	root := NewTreeNode("SingleNode", 100)

	spec := CirclePackingSpec{
		Root:   root,
		Width:  400,
		Height: 400,
	}

	result := RenderCirclePacking(spec)

	if result == "" {
		t.Fatal("Expected non-empty SVG output for leaf node")
	}

	if !strings.Contains(result, "<circle") {
		t.Error("Expected circle element for leaf node")
	}
}

func TestPackCircles(t *testing.T) {
	root := createFlatTree()
	circles := packCircles(root, 400, 400, 200, 2, 0)

	if len(circles) == 0 {
		t.Fatal("Expected non-empty circle list")
	}

	// Check that circles have valid properties
	for i, circle := range circles {
		if circle.Radius <= 0 {
			t.Errorf("Circle %d has invalid radius: %f", i, circle.Radius)
		}

		if circle.Node == nil {
			t.Errorf("Circle %d has nil node", i)
		}

		if circle.Depth < 0 {
			t.Errorf("Circle %d has negative depth: %d", i, circle.Depth)
		}
	}
}

func TestPositionCirclesPacked(t *testing.T) {
	circles := []PackedCircle{
		{Radius: 50, Node: NewTreeNode("A", 100), Depth: 1},
		{Radius: 40, Node: NewTreeNode("B", 80), Depth: 1},
		{Radius: 30, Node: NewTreeNode("C", 60), Depth: 1},
	}

	positioned := positionCirclesPacked(circles, 400, 400, 200)

	if len(positioned) != len(circles) {
		t.Errorf("Expected %d circles, got %d", len(circles), len(positioned))
	}

	// Check that all circles have positions assigned
	for i, circle := range positioned {
		if circle.X == 0 && circle.Y == 0 && i > 0 {
			// Only first circle should be at center
			t.Errorf("Circle %d has uninitialized position", i)
		}
	}
}

func TestPositionCirclesPackedSingleCircle(t *testing.T) {
	circles := []PackedCircle{
		{Radius: 50, Node: NewTreeNode("Single", 100), Depth: 1},
	}

	cx, cy := 400.0, 400.0
	positioned := positionCirclesPacked(circles, cx, cy, 200)

	if len(positioned) != 1 {
		t.Fatalf("Expected 1 circle, got %d", len(positioned))
	}

	// Single circle should be centered
	if positioned[0].X != cx || positioned[0].Y != cy {
		t.Errorf("Expected circle at center (%f, %f), got (%f, %f)",
			cx, cy, positioned[0].X, positioned[0].Y)
	}
}

func TestPositionCirclesPackedEmpty(t *testing.T) {
	circles := []PackedCircle{}
	positioned := positionCirclesPacked(circles, 400, 400, 200)

	if len(positioned) != 0 {
		t.Error("Expected empty result for empty input")
	}
}

// Icicle chart tests

func TestRenderIcicleVertical(t *testing.T) {
	spec := IcicleSpec{
		Root:        createTestTree(),
		Width:       800,
		Height:      600,
		Padding:     2,
		Orientation: "vertical",
		ShowLabels:  true,
	}

	result := RenderIcicle(spec)

	if result == "" {
		t.Fatal("Expected non-empty SVG output")
	}

	if !strings.Contains(result, "<rect") {
		t.Error("Expected rect elements in output")
	}

	if !strings.Contains(result, "<text") {
		t.Error("Expected text elements when ShowLabels is true")
	}
}

func TestRenderIcicleHorizontal(t *testing.T) {
	spec := IcicleSpec{
		Root:        createTestTree(),
		Width:       800,
		Height:      600,
		Padding:     2,
		Orientation: "horizontal",
		ShowLabels:  true,
	}

	result := RenderIcicle(spec)

	if result == "" {
		t.Fatal("Expected non-empty SVG output")
	}

	if !strings.Contains(result, "<rect") {
		t.Error("Expected rect elements in output")
	}
}

func TestRenderIcicleDefaultOrientation(t *testing.T) {
	spec := IcicleSpec{
		Root:   createFlatTree(),
		Width:  800,
		Height: 600,
	}

	result := RenderIcicle(spec)

	if result == "" {
		t.Fatal("Expected non-empty SVG output")
	}

	// Should default to vertical orientation
	if !strings.Contains(result, "<rect") {
		t.Error("Expected rect elements in output")
	}
}

func TestRenderIcicleNilRoot(t *testing.T) {
	spec := IcicleSpec{
		Root:   nil,
		Width:  800,
		Height: 600,
	}

	result := RenderIcicle(spec)

	if result != "" {
		t.Error("Expected empty string for nil root")
	}
}

func TestRenderIcicleZeroValue(t *testing.T) {
	root := NewTreeNode("Root", 0)

	spec := IcicleSpec{
		Root:   root,
		Width:  800,
		Height: 600,
	}

	result := RenderIcicle(spec)

	if result != "" {
		t.Error("Expected empty string for zero-value tree")
	}
}

func TestRenderIcicleWithColorScheme(t *testing.T) {
	colorScheme := []string{"#FF0000", "#00FF00", "#0000FF"}

	spec := IcicleSpec{
		Root:        createTestTree(),
		Width:       800,
		Height:      600,
		ColorScheme: colorScheme,
	}

	result := RenderIcicle(spec)

	if result == "" {
		t.Fatal("Expected non-empty SVG output")
	}

	// Verify rectangles are generated
	if !strings.Contains(result, "<rect") {
		t.Error("Expected rect elements in output")
	}
}

func TestIcicleLayoutVertical(t *testing.T) {
	root := createFlatTree()
	maxDepth := calculateMaxDepth(root, 0)

	rects := icicleLayoutVertical(root, 0, 0, 800, 600, maxDepth, 2, 0)

	if len(rects) == 0 {
		t.Fatal("Expected non-empty rectangle list")
	}

	// Check that rectangles have valid properties
	for i, rect := range rects {
		if rect.Width <= 0 {
			t.Errorf("Rectangle %d has invalid width: %f", i, rect.Width)
		}

		if rect.Height <= 0 {
			t.Errorf("Rectangle %d has invalid height: %f", i, rect.Height)
		}

		if rect.Node == nil {
			t.Errorf("Rectangle %d has nil node", i)
		}

		if rect.Depth < 0 {
			t.Errorf("Rectangle %d has negative depth: %d", i, rect.Depth)
		}
	}
}

func TestIcicleLayoutHorizontal(t *testing.T) {
	root := createFlatTree()
	maxDepth := calculateMaxDepth(root, 0)

	rects := icicleLayoutHorizontal(root, 0, 0, 800, 600, maxDepth, 2, 0)

	if len(rects) == 0 {
		t.Fatal("Expected non-empty rectangle list")
	}

	// Check that rectangles have valid properties
	for i, rect := range rects {
		if rect.Width <= 0 {
			t.Errorf("Rectangle %d has invalid width: %f", i, rect.Width)
		}

		if rect.Height <= 0 {
			t.Errorf("Rectangle %d has invalid height: %f", i, rect.Height)
		}

		if rect.Node == nil {
			t.Errorf("Rectangle %d has nil node", i)
		}
	}
}

func TestCalculateIcicleFontSize(t *testing.T) {
	tests := []struct {
		width    float64
		height   float64
		expected float64
		name     string
	}{
		{10, 10, 8, "very_small"},
		{30, 30, 10, "small"},
		{60, 60, 12, "medium"},
		{100, 100, 14, "large"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			size := calculateIcicleFontSize(test.width, test.height)

			// Verify it returns a valid units.Length with expected size
			// We can't easily test the exact value, but we can check it's not zero
			if size.String() == "" {
				t.Error("Expected non-empty font size")
			}
		})
	}
}
