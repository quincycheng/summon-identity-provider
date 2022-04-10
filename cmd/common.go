package cmd

import (
	"fmt"

	lipgloss "github.com/charmbracelet/lipgloss"
	colorgrad "github.com/mazznoer/colorgrad"
)

/*********************************
	Supportive Functions
**********************************/
var highlightStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#3490CE"))

func printHeader(header string) {
	grad, _ := colorgrad.NewGradient().HtmlColors("#2BB3EE", "grey").Build()

	for i, c := range header {
		fmt.Print(lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(grad.At(float64(i) / float64(len(header)-1)).Hex())).
			Render(string(c)))
	}
}
