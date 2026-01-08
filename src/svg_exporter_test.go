package main

import (
	"os"
	"strings"
	"testing"
)

func TestExportNamesSVG_Basic(t *testing.T) {
	names := []string{"Alice", "Bob", "Charlie"}
	settings := Settings{
		Width:        400,
		Margin:       10,
		ColGap:       10,
		FontSize:     16,
		LineHeight:   24,
		Columns:      2,
		FontFamily:   "Arial",
		ColumnColors: []string{"#ff0000", "#00ff00"},
	}
	tmpfile, err := os.CreateTemp("", "testsvg-*.svg")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	err = ExportNamesSVG(names, tmpfile.Name(), settings)
	if err != nil {
		t.Fatalf("ExportNamesSVG failed: %v", err)
	}

	data, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to read SVG file: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "<svg") {
		t.Error("SVG header missing")
	}
	for _, name := range names {
		if !strings.Contains(content, name) {
			t.Errorf("SVG missing name: %s", name)
		}
	}
}

func TestExportNamesSVG_ColumnColors(t *testing.T) {
	names := []string{"A", "B", "C", "D"}
	settings := Settings{
		Width:        400,
		Margin:       10,
		ColGap:       10,
		FontSize:     16,
		LineHeight:   24,
		Columns:      2,
		FontFamily:   "Arial",
		ColumnColors: []string{"#111111", "#222222"},
	}
	tmpfile, err := os.CreateTemp("", "testsvg-*.svg")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	err = ExportNamesSVG(names, tmpfile.Name(), settings)
	if err != nil {
		t.Fatalf("ExportNamesSVG failed: %v", err)
	}

	data, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to read SVG file: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "#111111") || !strings.Contains(content, "#222222") {
		t.Error("SVG missing expected column colors")
	}
}

func TestExportNamesSVG_UserColorMap(t *testing.T) {
	names := []string{"X", "Y"}
	settings := Settings{
		Width:        200,
		Margin:       5,
		ColGap:       5,
		FontSize:     12,
		LineHeight:   18,
		Columns:      1,
		FontFamily:   "Arial",
		ColumnColors: []string{"#000"},
		UserColorMap: map[string]string{"X": "#abc123"},
	}
	tmpfile, err := os.CreateTemp("", "testsvg-*.svg")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	err = ExportNamesSVG(names, tmpfile.Name(), settings)
	if err != nil {
		t.Fatalf("ExportNamesSVG failed: %v", err)
	}

	data, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to read SVG file: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "#abc123") {
		t.Error("SVG missing user color map color")
	}
}

func TestEscapeXML(t *testing.T) {
	input := `Tom & Jerry <Cartoon> "Best" 'Show'`
	expected := `Tom &amp; Jerry &lt;Cartoon&gt; &quot;Best&quot; &apos;Show&apos;`
	got := escapeXML(input)
	if got != expected {
		t.Errorf("escapeXML failed: got %q, want %q", got, expected)
	}
}
