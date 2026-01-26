package textutil

import "github.com/SCKelemen/text"

// ElideLabel truncates labels that exceed maxWidth pixels
// using middle ellipsis strategy for better readability.
//
// Example:
//
//	label := ElideLabel("Cloud Engineering Leadership", 90)
//	// Returns: "Cloud Engi...dership"
func ElideLabel(label string, maxWidth float64) string {
	txt := text.NewTerminal()
	return txt.Elide(label, maxWidth)
}

// ElideLabelEnd truncates labels at the end if they exceed maxWidth.
//
// Example:
//
//	label := ElideLabelEnd("Very long description text", 50)
//	// Returns: "Very long descri..."
func ElideLabelEnd(label string, maxWidth float64) string {
	txt := text.NewTerminal()
	return txt.ElideEnd(label, maxWidth)
}

// ElideLabelWithContext automatically detects label type and applies
// the most appropriate elision strategy.
func ElideLabelWithContext(label string, maxWidth float64) string {
	txt := text.NewTerminal()
	return txt.ElideAuto(label, maxWidth)
}
