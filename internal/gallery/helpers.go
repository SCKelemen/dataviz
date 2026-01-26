package gallery

import (
	"time"

	"github.com/SCKelemen/dataviz/charts"
)

// mustParseTime is a helper function for parsing dates in gallery configs
func mustParseTime(s string) time.Time {
	t, _ := time.Parse("2006-01-02", s)
	return t
}

// createSampleTree returns a sample tree structure for hierarchical charts
func createSampleTree() *charts.TreeNode {
	return &charts.TreeNode{
		Name:  "Root",
		Value: 100,
		Children: []*charts.TreeNode{
			{
				Name:  "Branch A",
				Value: 40,
				Children: []*charts.TreeNode{
					{Name: "Leaf A1", Value: 15},
					{Name: "Leaf A2", Value: 12},
					{Name: "Leaf A3", Value: 13},
				},
			},
			{
				Name:  "Branch B",
				Value: 35,
				Children: []*charts.TreeNode{
					{Name: "Leaf B1", Value: 20},
					{Name: "Leaf B2", Value: 15},
				},
			},
			{
				Name:  "Branch C",
				Value: 25,
				Children: []*charts.TreeNode{
					{Name: "Leaf C1", Value: 10},
					{Name: "Leaf C2", Value: 8},
					{Name: "Leaf C3", Value: 7},
				},
			},
		},
	}
}
