package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoadSettings_DefaultsWhenFileMissing(t *testing.T) {
	settings := LoadSettings("nonexistent_settings.conf")

	if !settings.ExportSVG {
		t.Errorf("Expected ExportSVG to be true by default")
	}
	if !settings.ExportTXT {
		t.Errorf("Expected ExportTXT to be true by default")
	}
	if settings.OutputDir != "output" {
		t.Errorf("Expected OutputDir to be 'output', got '%s'", settings.OutputDir)
	}
	if settings.DefaultCSVFile != "pledges.csv" {
		t.Errorf("Expected DefaultCSVFile to be 'pledges.csv', got '%s'", settings.DefaultCSVFile)
	}
	if settings.Width != 1161 {
		t.Errorf("Expected Width to be 1161, got %d", settings.Width)
	}
	if settings.Margin != 26 {
		t.Errorf("Expected Margin to be 26, got %d", settings.Margin)
	}
	if settings.ColGap != 54 {
		t.Errorf("Expected ColGap to be 54, got %d", settings.ColGap)
	}
	if settings.FontSize != 16 {
		t.Errorf("Expected FontSize to be 16, got %d", settings.FontSize)
	}
	if settings.LineHeight != 20 {
		t.Errorf("Expected LineHeight to be 20, got %d", settings.LineHeight)
	}
	if settings.Columns != 3 {
		t.Errorf("Expected Columns to be 3, got %d", settings.Columns)
	}
	if settings.FontFamily != "Trebuchet MS, Arial, sans-serif" {
		t.Errorf("Expected FontFamily to be 'Trebuchet MS, Arial, sans-serif', got '%s'", settings.FontFamily)
	}
	expectedColors := []string{"#3aff22", "#c622ff", "#a8ff21"}
	if !reflect.DeepEqual(settings.ColumnColors, expectedColors) {
		t.Errorf("Expected ColumnColors to be %v, got %v", expectedColors, settings.ColumnColors)
	}
	if !settings.RandomizeSVGColors {
		t.Errorf("Expected RandomizeSVGColors to be true by default")
	}
	if settings.UserColorMap != nil {
		t.Errorf("Expected UserColorMap to be nil by default")
	}
}

func TestLoadSettings_ParsesFileCorrectly(t *testing.T) {
	tmpDir := t.TempDir()
	confPath := filepath.Join(tmpDir, "settings.conf")
	content := `
EXPORT_SVG=false
EXPORT_TXT=false
OUTPUT_DIR=mydir
DEFAULT_CSV_FILE=myfile.csv
SVG_WIDTH=800
SVG_MARGIN_TO_EDGE=10
SVG_COLUMN_GAP=5
SVG_FONTSIZE=12
SVG_LINEHEIGHT=15
SVG_COLUMNS=2
SVG_FONTFAMILY=Arial
SVG_COLUMN_COLORS=#111,#222,#333
SVG_RANDOMIZE_COLORS=false
USER_COLOR_MAP=alice:#fff,bob:#000
`
	if err := os.WriteFile(confPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write temp settings.conf: %v", err)
	}

	settings := LoadSettings(confPath)

	if settings.ExportSVG {
		t.Errorf("Expected ExportSVG to be false")
	}
	if settings.ExportTXT {
		t.Errorf("Expected ExportTXT to be false")
	}
	if settings.OutputDir != "mydir" {
		t.Errorf("Expected OutputDir to be 'mydir', got '%s'", settings.OutputDir)
	}
	if settings.DefaultCSVFile != "myfile.csv" {
		t.Errorf("Expected DefaultCSVFile to be 'myfile.csv', got '%s'", settings.DefaultCSVFile)
	}
	if settings.Width != 800 {
		t.Errorf("Expected Width to be 800, got %d", settings.Width)
	}
	if settings.Margin != 10 {
		t.Errorf("Expected Margin to be 10, got %d", settings.Margin)
	}
	if settings.ColGap != 5 {
		t.Errorf("Expected ColGap to be 5, got %d", settings.ColGap)
	}
	if settings.FontSize != 12 {
		t.Errorf("Expected FontSize to be 12, got %d", settings.FontSize)
	}
	if settings.LineHeight != 15 {
		t.Errorf("Expected LineHeight to be 15, got %d", settings.LineHeight)
	}
	if settings.Columns != 2 {
		t.Errorf("Expected Columns to be 2, got %d", settings.Columns)
	}
	if settings.FontFamily != "Arial" {
		t.Errorf("Expected FontFamily to be 'Arial', got '%s'", settings.FontFamily)
	}
	expectedColors := []string{"#111", "#222", "#333"}
	if !reflect.DeepEqual(settings.ColumnColors, expectedColors) {
		t.Errorf("Expected ColumnColors to be %v, got %v", expectedColors, settings.ColumnColors)
	}
	if settings.RandomizeSVGColors {
		t.Errorf("Expected RandomizeSVGColors to be false")
	}
	expectedMap := map[string]string{"alice": "#fff", "bob": "#000"}
	if !reflect.DeepEqual(settings.UserColorMap, expectedMap) {
		t.Errorf("Expected UserColorMap to be %v, got %v", expectedMap, settings.UserColorMap)
	}
}

func TestLoadSettings_IgnoresCommentsAndBlankLines(t *testing.T) {
	tmpDir := t.TempDir()
	confPath := filepath.Join(tmpDir, "settings.conf")
	content := `
# This is a comment
EXPORT_SVG=false

# Another comment
EXPORT_TXT=true

`
	if err := os.WriteFile(confPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write temp settings.conf: %v", err)
	}

	settings := LoadSettings(confPath)
	if settings.ExportSVG {
		t.Errorf("Expected ExportSVG to be false")
	}
	if !settings.ExportTXT {
		t.Errorf("Expected ExportTXT to be true")
	}
}

func TestLoadSettings_InvalidLinesAreIgnored(t *testing.T) {
	tmpDir := t.TempDir()
	confPath := filepath.Join(tmpDir, "settings.conf")
	content := `
INVALID_LINE
SVG_WIDTH=900
ANOTHER_INVALID_LINE
`
	if err := os.WriteFile(confPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write temp settings.conf: %v", err)
	}

	settings := LoadSettings(confPath)
	if settings.Width != 900 {
		t.Errorf("Expected Width to be 900, got %d", settings.Width)
	}
}
