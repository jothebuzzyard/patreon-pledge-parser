package main

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"
)

func ExportNamesSVG(names []string, outputPath string, settings Settings) error {
	var (
		width        = settings.Width
		margin       = settings.Margin
		colGap       = settings.ColGap
		fontSize     = settings.FontSize
		lineHeight   = settings.LineHeight
		columns      = settings.Columns
		fontFamily   = settings.FontFamily
		columnColors = settings.ColumnColors
	)
	r := rand.New(rand.NewSource(int64(len(names))))
	// Sort names alphabetically
	sort.Slice(names, func(i, j int) bool {
		return strings.ToLower(names[i]) < strings.ToLower(names[j])
	})

	// Split names into columns
	nPerCol := (len(names) + columns - 1) / columns
	colNames := make([][]string, columns)
	for i, name := range names {
		col := i / nPerCol
		if col >= columns {
			col = columns - 1
		}
		colNames[col] = append(colNames[col], name)
	}

	// Calculate SVG height
	maxColLen := 0
	for _, col := range colNames {
		if len(col) > maxColLen {
			maxColLen = len(col)
		}
	}
	height := margin*2 + maxColLen*lineHeight

	// Prepare SVG file
	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// SVG header
	fmt.Fprintf(f, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d">`, width, height)
	fmt.Fprintf(f, `<rect width="100%%" height="100%%" fill="none"/>`)

	// Column positions
	colWidth := (width - margin*2 - colGap*(columns-1)) / columns

	for colIdx, col := range colNames {
		x := margin + colIdx*(colWidth+colGap) + colWidth/2
		for rowIdx, name := range col {
			y := margin + fontSize + rowIdx*lineHeight
			color := getColorForColumn(colIdx, columnColors, settings.RandomizeSVGColors, r)
			fmt.Fprintf(f,
				`<text x="%d" y="%d" font-family="%s" font-size="%d" fill="none" stroke="#000" stroke-width="2" paint-order="stroke" text-anchor="middle">%s</text>`,
				x, y, fontFamily, fontSize, escapeXML(name))
			fmt.Fprintf(f,
				`<text x="%d" y="%d" font-family="%s" font-size="%d" fill="%s" stroke="none" text-anchor="middle">%s</text>`,
				x, y, fontFamily, fontSize, color, escapeXML(name))
		}
	}

	fmt.Fprint(f, `</svg>`)
	return nil
}

func getColorForColumn(colIdx int, columnColors []string, randomize bool, r *rand.Rand) string {
	if randomize {
		return columnColors[r.Intn(len(columnColors))]
	}
	return columnColors[colIdx%len(columnColors)]
}

// Helper to escape XML special characters in names
func escapeXML(s string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		"\"", "&quot;",
		"'", "&apos;",
	)
	return replacer.Replace(s)
}
