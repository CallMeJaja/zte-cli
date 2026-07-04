package cli

import (
	"bufio"
	"fmt"
	"strings"
)

// MenuItem represents a single item in the menu grid.
type MenuItem struct {
	ID          int
	Name        string
	Description string
}

// Menu represents a screen with a 2-column grid of items.
type Menu struct {
	Title   string
	Items   []MenuItem
	Width   int // Total terminal width (default 78)
}

// NewMenu creates a new menu with default width.
func NewMenu(title string, items []MenuItem) *Menu {
	return &Menu{
		Title: title,
		Items: items,
		Width: 78,
	}
}

// Render prints the menu in 2-column aaPanel-style layout.
func (m *Menu) Render() {
	sep := strings.Repeat("=", m.Width)
	dash := strings.Repeat("-", m.Width)

	fmt.Println(sep)
	if m.Title != "" {
		// Center the title
		padding := (m.Width - len(m.Title)) / 2
		if padding > 0 {
			fmt.Printf("%s%s\n", strings.Repeat(" ", padding), m.Title)
		} else {
			fmt.Println(m.Title)
		}
	}
	fmt.Println(dash)

	// Render items in 2 columns
	half := (len(m.Items) + 1) / 2 // ceiling division
	for i := 0; i < half; i++ {
		left := m.Items[i]
		leftStr := fmt.Sprintf("(%d) %s", left.ID, left.Name)

		rightStr := ""
		rightIdx := i + half
		if rightIdx < len(m.Items) {
			right := m.Items[rightIdx]
			rightStr = fmt.Sprintf("(%d) %s", right.ID, right.Name)
		}

		// Pad left to column width (~40 chars) then print right
		if rightStr != "" {
			fmt.Printf("%-40s%s\n", leftStr, rightStr)
		} else {
			fmt.Println(leftStr)
		}
	}

	fmt.Println(sep)
}

// Prompt shows the input prompt and returns the user's choice.
func (m *Menu) Prompt(scanner *bufio.Scanner) string {
	fmt.Print("Pls enter command number: ")
	if !scanner.Scan() {
		return ""
	}
	return strings.TrimSpace(scanner.Text())
}

// RenderSubHeader prints a separator for data display.
func RenderSubHeader(width int) {
	fmt.Println(strings.Repeat("=", width))
}

// RenderDataHeader prints the title bar for a data view.
func RenderDataHeader(title string, width int) {
	sep := strings.Repeat("=", width)
	dash := strings.Repeat("-", width)
	fmt.Println(sep)
	if title != "" {
		padding := (width - len(title)) / 2
		if padding > 0 {
			fmt.Printf("%s%s\n", strings.Repeat(" ", padding), title)
		} else {
			fmt.Println(title)
		}
	}
	fmt.Println(dash)
}

// RenderDataRow prints a key-value pair in the data view.
func RenderDataRow(key, value string, keyWidth int) {
	if value == "" {
		value = "N/A"
	}
	padding := keyWidth - len(key)
	if padding < 1 {
		padding = 1
	}
	fmt.Printf("  %s%s: %s\n", key, strings.Repeat(" ", padding), value)
}

// RenderDataFooter prints the action bar for data view (Refresh + Back).
func RenderDataFooter(width int, actions ...string) {
	sep := strings.Repeat("=", width)
	fmt.Println(sep)
	if len(actions) > 0 {
		fmt.Println(strings.Join(actions, "  |  "))
	} else {
		fmt.Println("(0) Kembali")
	}
	fmt.Println(sep)
}
