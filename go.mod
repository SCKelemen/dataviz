module github.com/SCKelemen/dataviz

go 1.25.4

require (
	// SCKelemen foundation libraries
	github.com/SCKelemen/color v1.0.0
	github.com/SCKelemen/design-system v0.1.0

	// SCKelemen rendering stack
	github.com/SCKelemen/layout v1.1.0
	github.com/SCKelemen/svg v0.2.0
	github.com/SCKelemen/text v1.0.0
	github.com/SCKelemen/units v1.0.2
	github.com/modelcontextprotocol/go-sdk v1.2.0
	github.com/srwiley/oksvg v0.0.0-20220731023508-a61f04f16b76
	github.com/srwiley/rasterx v0.0.0-20220730225603-2ab79fcdd4ef
)

require (
	github.com/SCKelemen/unicode v1.0.1-0.20251225190048-233be2b0d647 // indirect
	github.com/google/jsonschema-go v0.3.0 // indirect
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
	golang.org/x/image v0.34.0 // indirect
	golang.org/x/net v0.0.0-20211118161319-6a13c67c3ce4 // indirect
	golang.org/x/oauth2 v0.30.0 // indirect
	golang.org/x/text v0.32.0 // indirect
)

// Exclude problematic test-only dependency (used only in layout tests)
exclude github.com/SCKelemen/wpt-test-gen v0.0.0-00010101000000-000000000000

exclude github.com/SCKelemen/wpt-test-gen v0.0.0-20251213153317-6265321ae2a1

// Use replace for layout to avoid wpt-test-gen issues
replace github.com/SCKelemen/layout => github.com/SCKelemen/layout v1.1.0

// Local development - replace with local paths when developing
// Commented out for CI/CD - uncomment for local development
// replace github.com/SCKelemen/text => /Users/samuel.kelemen/Code/github.com/SCKelemen/text
