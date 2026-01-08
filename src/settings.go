package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Settings holds configuration values with type safety
type Settings struct {
	ExportSVG      bool
	ExportTXT      bool
	OutputDir      string
	DefaultCSVFile string

	Width              int
	Margin             int
	ColGap             int
	FontSize           int
	LineHeight         int
	Columns            int
	FontFamily         string
	ColumnColors       []string
	RandomizeSVGColors bool
}

// LoadSettings loads settings from settings.conf and returns a Settings object
func LoadSettings(path string) Settings {
	fmt.Print("Loading settings.conf\n")
	settings := Settings{
		ExportSVG:      true,
		ExportTXT:      true,
		OutputDir:      "output",
		DefaultCSVFile: "pledges.csv",

		Width:              1161,
		Margin:             26,
		ColGap:             54,
		FontSize:           16,
		LineHeight:         20,
		Columns:            3,
		FontFamily:         "Trebuchet MS, Arial, sans-serif",
		ColumnColors:       []string{"#3aff22", "#c622ff", "#a8ff21"},
		RandomizeSVGColors: true,
	}

	file, err := os.Open(path)
	if err != nil {
		fmt.Print("No settings.conf found. Loading defaults. \n If you would like to change the behaviour of this exporter, create a settings.conf file in the same directory.\n See README for help! \n")
		return settings // Use defaults if file missing
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		switch key {
		case "EXPORT_SVG":
			settings.ExportSVG = strings.ToLower(val) == "true"
		case "EXPORT_TXT":
			settings.ExportTXT = strings.ToLower(val) == "true"
		case "OUTPUT_DIR":
			settings.OutputDir = val
		case "DEFAULT_CSV_FILE":
			settings.DefaultCSVFile = val
		case "SVG_WIDTH":
			fmt.Sscanf(val, "%d", &settings.Width)
		case "SVG_MARGIN_TO_EDGE":
			fmt.Sscanf(val, "%d", &settings.Margin)
		case "SVG_COLUMN_GAP":
			fmt.Sscanf(val, "%d", &settings.ColGap)
		case "SVG_FONTSIZE":
			fmt.Sscanf(val, "%d", &settings.FontSize)
		case "SVG_LINEHEIGHT":
			fmt.Sscanf(val, "%d", &settings.LineHeight)
		case "SVG_COLUMNS":
			fmt.Sscanf(val, "%d", &settings.Columns)
		case "SVG_FONTFAMILY":
			settings.FontFamily = val
		case "SVG_COLUMN_COLORS":
			settings.ColumnColors = strings.Split(val, ",")
		case "SVG_RANDOMIZE_COLORS":
			settings.RandomizeSVGColors = strings.ToLower(val) == "true"
		}
	}
	return settings
}
