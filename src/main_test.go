package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestGetCSVPath_FileExists(t *testing.T) {
	tmpDir := t.TempDir()
	csvFile := filepath.Join(tmpDir, "pledges.csv")
	err := os.WriteFile(csvFile, []byte("header1,header2\n"), 0644)
	if err != nil {
		t.Fatalf("failed to create csv: %v", err)
	}
	path, err := getCSVPath(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != csvFile {
		t.Errorf("expected %s, got %s", csvFile, path)
	}
}

func TestConfirmAndCleanOutputDir_DirExists(t *testing.T) {
	tmpDir := t.TempDir()
	// Simulate user input "y"
	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	w.WriteString("y\n")
	w.Close()
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }()

	err := os.MkdirAll(tmpDir, 0755)
	if err != nil {
		t.Fatalf("failed to create dir: %v", err)
	}
	err = confirmAndCleanOutputDir(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(tmpDir); !os.IsNotExist(err) {
		t.Errorf("directory should be deleted")
	}
}

func TestConfirmAndCleanOutputDir_Abort(t *testing.T) {
	tmpDir := t.TempDir()
	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	w.WriteString("n\n")
	w.Close()
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }()

	err := confirmAndCleanOutputDir(tmpDir)
	if err == nil || !strings.Contains(err.Error(), "aborted") {
		t.Errorf("expected abort error, got %v", err)
	}
}

func TestReadCSVFile(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "test.csv")
	content := "a,b,c\n1,2,3\n"
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to write csv: %v", err)
	}
	records, err := readCSVFile(tmpFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(records) != 2 || records[1][0] != "1" {
		t.Errorf("unexpected records: %v", records)
	}
}

func TestParsePatrons(t *testing.T) {
	records := [][]string{
		{"Name", "Email", "Discord", "PatronStatus", "FollowsYou", "FreeMember", "FreeTrial", "LifetimeAmount", "PledgeAmount", "ChargeFrequency", "Tier", "Addressee", "Street", "City", "State", "Zip", "Country", "Phone", "PatronageSinceDate", "LastChargeDate", "LastChargeStatus", "AdditionalDetails", "UserID", "LastUpdated", "Currency", "MaxPosts", "AccessExpiration", "NextChargeDate", "FullCountryName", "SubscriptionSource"},
		{"Alice", "a@b.com", "", "", "", "", "", "", "", "", "Gold", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"Bob", "b@b.com", "", "", "", "", "", "", "", "", "Free", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
	}
	patrons, freeCount := parsePatrons(records)
	if len(patrons) != 2 {
		t.Errorf("expected 2 patrons, got %d", len(patrons))
	}
	if freeCount != 1 {
		t.Errorf("expected 1 free tier, got %d", freeCount)
	}
}

func TestFilterPatrons(t *testing.T) {
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	patrons := []Patron{
		{Name: "A", Tier: "Gold", LastChargeStatus: "paid", AccessExpiration: ""},
		{Name: "B", Tier: "Free", LastChargeStatus: "paid", AccessExpiration: ""},
		{Name: "C", Tier: "Silver", LastChargeStatus: "unpaid", AccessExpiration: ""},
		{Name: "D", Tier: "Gold", LastChargeStatus: "paid", AccessExpiration: "2023-01-01 00:00:00"},
	}
	filtered, expired, unpaid := filterPatrons(patrons, now)
	if len(filtered) != 1 || filtered[0].Name != "A" {
		t.Errorf("unexpected filtered: %v", filtered)
	}
	if expired != 1 {
		t.Errorf("expected 1 expired, got %d", expired)
	}
	if unpaid != 1 {
		t.Errorf("expected 1 unpaid, got %d", unpaid)
	}
}

func TestGroupAndSortByTier(t *testing.T) {
	patrons := []Patron{
		{Name: "Charlie", Tier: "Gold"},
		{Name: "Alice", Tier: "Gold"},
		{Name: "Bob", Tier: "Silver"},
	}
	groups := groupAndSortByTier(patrons)
	if len(groups) != 2 {
		t.Errorf("expected 2 groups, got %d", len(groups))
	}
	if groups["Gold"][0].Name != "Alice" {
		t.Errorf("expected Alice first in Gold, got %s", groups["Gold"][0].Name)
	}
}

func TestWriteTierFiles(t *testing.T) {
	tmpDir := t.TempDir()
	groups := map[string][]Patron{
		"Gold":   {{Name: "Alice"}, {Name: "Bob"}},
		"Silver": {{Name: "Charlie"}},
	}
	total, err := writeTierFiles(tmpDir, groups)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 3 {
		t.Errorf("expected 3 patrons, got %d", total)
	}
	// Check files exist and contents
	goldFile := filepath.Join(tmpDir, "Gold.txt")
	data, err := ioutil.ReadFile(goldFile)
	if err != nil {
		t.Fatalf("failed to read Gold.txt: %v", err)
	}
	if !strings.Contains(string(data), "Alice") || !strings.Contains(string(data), "Bob") {
		t.Errorf("Gold.txt missing patron names: %s", string(data))
	}
}
