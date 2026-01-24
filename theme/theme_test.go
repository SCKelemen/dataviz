package theme

import (
	"strings"
	"testing"

	design "github.com/SCKelemen/design-system"
)

func TestNew(t *testing.T) {
	tokens := design.DefaultTheme()
	theme := New(tokens)

	if theme == nil {
		t.Fatal("Expected non-nil theme")
	}

	if theme.Tokens != tokens {
		t.Error("Expected tokens to be set")
	}

	if len(theme.ColorScheme.Categorical) == 0 {
		t.Error("Expected categorical colors to be set")
	}

	if len(theme.ColorScheme.Sequential) == 0 {
		t.Error("Expected sequential colors to be set")
	}

	if theme.Typography.TitleFont == "" {
		t.Error("Expected title font to be set")
	}
}

func TestNewWithNilTokens(t *testing.T) {
	theme := New(nil)

	if theme == nil {
		t.Fatal("Expected non-nil theme even with nil tokens")
	}

	if theme.Tokens == nil {
		t.Error("Expected default tokens to be set")
	}
}

func TestGetColor(t *testing.T) {
	theme := Default()

	// Test that colors cycle through palette
	color1 := theme.GetColor(0)
	color2 := theme.GetColor(1)
	color3 := theme.GetColor(len(theme.ColorScheme.Categorical))

	if color1 == "" {
		t.Error("Expected non-empty color")
	}

	if color1 == color2 {
		t.Error("Expected different colors for different indices")
	}

	// Should cycle back to first color
	if color1 != color3 {
		t.Error("Expected color to cycle back to first")
	}
}

func TestGetSequentialColor(t *testing.T) {
	theme := Default()

	tests := []struct {
		t        float64
		name     string
	}{
		{0.0, "start"},
		{0.5, "middle"},
		{1.0, "end"},
		{-0.5, "below_zero"},
		{1.5, "above_one"},
	}

	for _, test := range tests {
		color := theme.GetSequentialColor(test.t)
		if color == "" {
			t.Errorf("Expected non-empty color for t=%f (%s)", test.t, test.name)
		}

		// Check that it's a valid hex color
		if !strings.HasPrefix(color, "#") {
			t.Errorf("Expected hex color for t=%f (%s), got %s", test.t, test.name, color)
		}
	}
}

func TestGetDivergingColor(t *testing.T) {
	theme := Default()

	tests := []struct {
		t        float64
		name     string
	}{
		{-1.0, "left_extreme"},
		{-0.5, "left"},
		{0.0, "center"},
		{0.5, "right"},
		{1.0, "right_extreme"},
		{-2.0, "below_range"},
		{2.0, "above_range"},
	}

	for _, test := range tests {
		color := theme.GetDivergingColor(test.t)
		if color == "" {
			t.Errorf("Expected non-empty color for t=%f (%s)", test.t, test.name)
		}

		if !strings.HasPrefix(color, "#") {
			t.Errorf("Expected hex color for t=%f (%s), got %s", test.t, test.name, color)
		}
	}
}

func TestTitleStyle(t *testing.T) {
	theme := Default()
	style := theme.TitleStyle()

	if style.FontFamily == "" {
		t.Error("Expected font family to be set")
	}

	if style.FontSize.Value == 0 {
		t.Error("Expected font size to be set")
	}

	if style.Fill == "" {
		t.Error("Expected text color to be set")
	}
}

func TestBodyStyle(t *testing.T) {
	theme := Default()
	style := theme.BodyStyle()

	if style.FontFamily == "" {
		t.Error("Expected font family to be set")
	}

	if style.FontSize.Value == 0 {
		t.Error("Expected font size to be set")
	}
}

func TestLabelStyle(t *testing.T) {
	theme := Default()
	style := theme.LabelStyle()

	if style.FontFamily == "" {
		t.Error("Expected font family to be set")
	}

	if style.FontSize.Value == 0 {
		t.Error("Expected font size to be set")
	}
}

func TestGridStyle(t *testing.T) {
	theme := Default()
	style := theme.GridStyle()

	if style.Stroke == "" {
		t.Error("Expected stroke color to be set")
	}

	if style.StrokeWidth == 0 {
		t.Error("Expected stroke width to be set")
	}

	if style.Opacity == 0 {
		t.Error("Expected opacity to be set")
	}
}

func TestAxisStyle(t *testing.T) {
	theme := Default()
	style := theme.AxisStyle()

	if style.Stroke == "" {
		t.Error("Expected stroke color to be set")
	}

	if style.StrokeWidth == 0 {
		t.Error("Expected stroke width to be set")
	}
}

func TestDarkModeColors(t *testing.T) {
	tokens := design.DefaultTheme()
	tokens.Mode = "dark"
	theme := New(tokens)

	// Check that dark mode generates appropriate colors
	if theme.ColorScheme.GridColor == "" {
		t.Error("Expected grid color to be set")
	}

	// Dark mode should have darker grid colors
	if !strings.HasPrefix(theme.ColorScheme.GridColor, "#") {
		t.Error("Expected valid hex color")
	}
}

func TestLightModeColors(t *testing.T) {
	tokens := design.PaperTheme()
	theme := New(tokens)

	// Check that light mode generates appropriate colors
	if theme.ColorScheme.GridColor == "" {
		t.Error("Expected grid color to be set")
	}

	// Light mode should have lighter grid colors
	if !strings.HasPrefix(theme.ColorScheme.GridColor, "#") {
		t.Error("Expected valid hex color")
	}
}

func TestDensityScaling(t *testing.T) {
	compactTokens := design.DefaultTheme()
	compactTokens.Density = "compact"
	compactTheme := New(compactTokens)

	comfortableTokens := design.DefaultTheme()
	comfortableTokens.Density = "comfortable"
	comfortableTheme := New(comfortableTokens)

	// Compact should have smaller font sizes
	if compactTheme.Typography.BodySize.Value >= comfortableTheme.Typography.BodySize.Value {
		t.Error("Expected compact theme to have smaller font sizes")
	}

	// Compact should have smaller point sizes
	if compactTheme.Chart.PointSize >= comfortableTheme.Chart.PointSize {
		t.Error("Expected compact theme to have smaller point sizes")
	}
}

func TestPresetThemes(t *testing.T) {
	presets := []struct {
		name  string
		theme *Theme
	}{
		{"Default", Default()},
		{"Midnight", Midnight()},
		{"Nord", Nord()},
		{"Paper", Paper()},
		{"Wrapped", Wrapped()},
		{"Monochrome Dark", Monochrome(true)},
		{"Monochrome Light", Monochrome(false)},
		{"Ocean Dark", Ocean(true)},
		{"Ocean Light", Ocean(false)},
		{"Forest Dark", Forest(true)},
		{"Forest Light", Forest(false)},
		{"Sunset Dark", Sunset(true)},
		{"Sunset Light", Sunset(false)},
		{"High Contrast Dark", HighContrast(true)},
		{"High Contrast Light", HighContrast(false)},
		{"Scientific", Scientific()},
		{"Minimal", Minimal()},
	}

	for _, preset := range presets {
		if preset.theme == nil {
			t.Errorf("Preset %s returned nil theme", preset.name)
			continue
		}

		if preset.theme.Tokens == nil {
			t.Errorf("Preset %s has nil tokens", preset.name)
		}

		if len(preset.theme.ColorScheme.Categorical) == 0 {
			t.Errorf("Preset %s has empty categorical colors", preset.name)
		}

		if preset.theme.Typography.TitleFont == "" {
			t.Errorf("Preset %s has empty title font", preset.name)
		}
	}
}

func TestFromTokens(t *testing.T) {
	tokens := &design.DesignTokens{
		Color:      "#FF0000",
		Background: "#000000",
		Accent:     "#00FF00",
		FontFamily: "Arial",
		Mode:       "dark",
		Layout:     design.DefaultLayoutTokens(),
	}

	theme := FromTokens(tokens)

	if theme == nil {
		t.Fatal("Expected non-nil theme")
	}

	if theme.Tokens != tokens {
		t.Error("Expected tokens to be set")
	}

	if theme.ColorScheme.TextColor != "#FF0000" {
		t.Errorf("Expected text color #FF0000, got %s", theme.ColorScheme.TextColor)
	}

	if theme.ColorScheme.BackgroundColor != "#000000" {
		t.Errorf("Expected background color #000000, got %s", theme.ColorScheme.BackgroundColor)
	}
}

func TestScientificThemeColors(t *testing.T) {
	theme := Scientific()

	// Scientific theme should use Okabe-Ito colorblind-friendly palette
	if len(theme.ColorScheme.Categorical) < 8 {
		t.Error("Expected at least 8 categorical colors for Okabe-Ito palette")
	}

	// Check for some known Okabe-Ito colors
	colors := theme.ColorScheme.Categorical
	foundOrange := false
	for _, c := range colors {
		if strings.ToUpper(c) == "#E69F00" {
			foundOrange = true
			break
		}
	}

	if !foundOrange {
		t.Error("Expected Okabe-Ito orange color in scientific theme")
	}
}

func TestMinimalTheme(t *testing.T) {
	theme := Minimal()

	// Minimal theme should have subtle styling
	if theme.Chart.GridOpacity > 0.3 {
		t.Errorf("Expected minimal theme to have low grid opacity, got %f", theme.Chart.GridOpacity)
	}

	if theme.Chart.GridStrokeWidth > 1.0 {
		t.Errorf("Expected minimal theme to have thin grid lines, got %f", theme.Chart.GridStrokeWidth)
	}

	if theme.Tokens.Density != "compact" {
		t.Error("Expected minimal theme to use compact density")
	}
}
