package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Patron struct {
	Name               string
	Email              string
	Discord            string
	PatronStatus       string
	FollowsYou         string
	FreeMember         string
	FreeTrial          string
	LifetimeAmount     string
	PledgeAmount       string
	ChargeFrequency    string
	Tier               string
	Addressee          string
	Street             string
	City               string
	State              string
	Zip                string
	Country            string
	Phone              string
	PatronageSinceDate string
	LastChargeDate     string
	LastChargeStatus   string
	AdditionalDetails  string
	UserID             string
	LastUpdated        string
	Currency           string
	MaxPosts           string
	AccessExpiration   string
	NextChargeDate     string
	FullCountryName    string
	SubscriptionSource string
}

func getCSVPath(baseDir string) (string, error) {
	csvPath := filepath.Join(baseDir, "pledges.csv")
	fmt.Print("Looking for a file named 'pledges.csv' in the current directory...\n")
	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		fmt.Print("File 'pledges.csv' not found. Please enter the path to your CSV file (with .csv at the end): ")
		fmt.Scanln(&csvPath)
		if _, err := os.Stat(csvPath); os.IsNotExist(err) {
			return "", fmt.Errorf("file '%s' not found", csvPath)
		}
	}
	return csvPath, nil
}

func confirmAndCleanOutputDir(outputDir string) error {
	if stat, err := os.Stat(outputDir); err == nil && stat.IsDir() {
		var response string
		fmt.Printf("Output directory '%s' already exists. Delete it and continue? (y/N): ", outputDir)
		fmt.Scanln(&response)
		if strings.ToLower(strings.TrimSpace(response)) != "y" {
			return fmt.Errorf("aborted by user")
		}
		err := os.RemoveAll(outputDir)
		if err != nil {
			return fmt.Errorf("error deleting output directory: %v", err)
		}
	}
	return nil
}

func readCSVFile(csvPath string) ([][]string, error) {
	file, err := os.Open(csvPath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV: %v", err)
	}
	return records, nil
}

func parsePatrons(records [][]string) ([]Patron, int) {
	var patrons []Patron
	var freeTierCount int
	for i, record := range records {
		if i == 0 {
			continue // Skip header
		}
		if len(record) < 30 {
			continue // Skip malformed rows
		}
		patron := Patron{
			Name:               record[0],
			Email:              record[1],
			Discord:            record[2],
			PatronStatus:       record[3],
			FollowsYou:         record[4],
			FreeMember:         record[5],
			FreeTrial:          record[6],
			LifetimeAmount:     record[7],
			PledgeAmount:       record[8],
			ChargeFrequency:    record[9],
			Tier:               record[10],
			Addressee:          record[11],
			Street:             record[12],
			City:               record[13],
			State:              record[14],
			Zip:                record[15],
			Country:            record[16],
			Phone:              record[17],
			PatronageSinceDate: record[18],
			LastChargeDate:     record[19],
			LastChargeStatus:   record[20],
			AdditionalDetails:  record[21],
			UserID:             record[22],
			LastUpdated:        record[23],
			Currency:           record[24],
			MaxPosts:           record[25],
			AccessExpiration:   record[26],
			NextChargeDate:     record[27],
			FullCountryName:    record[28],
			SubscriptionSource: record[29],
		}
		if strings.Contains(patron.Tier, "Free") {
			freeTierCount++
		}
		patrons = append(patrons, patron)
	}
	return patrons, freeTierCount
}

func filterPatrons(patrons []Patron, now time.Time) ([]Patron, int, int) {
	var filteredPatrons []Patron
	var expiredAccessCount, unpaidStatusCount int
	for _, patron := range patrons {
		if strings.Contains(patron.Tier, "Free") {
			continue
		}
		if strings.ToLower(strings.TrimSpace(patron.LastChargeStatus)) != "paid" {
			unpaidStatusCount++
			continue
		}
		if patron.AccessExpiration != "" {
			expiration, err := time.Parse("2006-01-02 15:04:05", patron.AccessExpiration)
			if err == nil && expiration.Before(now) {
				expiredAccessCount++
				continue
			}
		}
		filteredPatrons = append(filteredPatrons, patron)
	}
	return filteredPatrons, expiredAccessCount, unpaidStatusCount
}

func groupAndSortByTier(patrons []Patron) map[string][]Patron {
	tierGroups := make(map[string][]Patron)
	for _, patron := range patrons {
		tier := strings.TrimSpace(patron.Tier)
		if tier == "" {
			continue
		}
		tierGroups[tier] = append(tierGroups[tier], patron)
	}
	for tier := range tierGroups {
		sort.Slice(tierGroups[tier], func(i, j int) bool {
			return strings.ToLower(tierGroups[tier][i].Name) < strings.ToLower(tierGroups[tier][j].Name)
		})
	}
	return tierGroups
}

func writeTierFiles(outputDir string, tierGroups map[string][]Patron) (int, error) {
	totalPaying := 0
	for tier, patrons := range tierGroups {
		filename := strings.ReplaceAll(tier, " ", "_")
		filename = strings.ReplaceAll(filename, "/", "_")
		filename = strings.ReplaceAll(filename, "\\", "_")
		filename = filepath.Join(outputDir, filename+".txt")
		file, err := os.Create(filename)
		if err != nil {
			fmt.Printf("Error creating file %s: %v\n", filename, err)
			continue
		}
		for _, patron := range patrons {
			_, err := file.WriteString(patron.Name + "\n")
			if err != nil {
				fmt.Printf("Error writing to file %s: %v\n", filename, err)
			}
		}
		file.Close()
		fmt.Printf("Created %s with %d patrons\n", filename, len(patrons))
		totalPaying += len(patrons)
	}
	return totalPaying, nil
}

func main() {
	baseDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current working directory: %v\n", err)
		fmt.Print("Press Enter to exit...")
		fmt.Scanln()
		return
	}
	outputDir := filepath.Join(baseDir, "output")

	csvPath, err := getCSVPath(baseDir)
	if err != nil {
		fmt.Println(err)
		fmt.Print("Press Enter to exit...")
		fmt.Scanln()
		return
	}

	if err := confirmAndCleanOutputDir(outputDir); err != nil {
		fmt.Println(err)
		return
	}

	records, err := readCSVFile(csvPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(records) < 2 {
		fmt.Println("CSV file is empty or has no data rows")
		return
	}

	patrons, freeTierCount := parsePatrons(records)
	filteredPatrons, expiredAccessCount, unpaidStatusCount := filterPatrons(patrons, time.Now().UTC())
	tierGroups := groupAndSortByTier(filteredPatrons)

	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		return
	}

	totalPaying, _ := writeTierFiles(outputDir, tierGroups)

	fmt.Println("----- Summary -----")
	fmt.Printf("Total paying patrons: %d\n", totalPaying)
	fmt.Printf("Total free tier patrons: %d\n", freeTierCount)
	fmt.Printf("Skipped due to expired access: %d\n", expiredAccessCount)
	fmt.Printf("Skipped due to unpaid status: %d\n", unpaidStatusCount)
	fmt.Println("-------------------")
	fmt.Println("Processing complete! Your files are in the 'output' directory.")
	fmt.Print("Press Enter to exit...")
	fmt.Scanln()
}
