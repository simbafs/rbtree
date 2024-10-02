package main

import "github.com/charmbracelet/lipgloss"

func (t RBTree) View() string {
	return t.Root.View(t.Nil)
}

func (node *RBNode) View(Nil *RBNode) string {
	style := lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
	if node == nil {
		return ""
	}

	if node.Color == Red {
		style = style.Background(lipgloss.Color("#FF0000"))
	} else {
		style = style.Background(lipgloss.Color("#000000"))
	}

	if node == Nil {
		return style.Render(node.Str())
	}

	l := node.Left.View(Nil)
	r := node.Right.View(Nil)
	self := style.Width(lipgloss.Width(l) + lipgloss.Width(r) - 2). // -2 for the border
									Render(node.Str())

	return lipgloss.JoinVertical(lipgloss.Center,
		self,
		lipgloss.JoinHorizontal(lipgloss.Top, l, r),
	)
}