package main

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strings"
	"svg-to-ruida/converter"
	"testing"
)

func TestConverter(t *testing.T) {
	demonstrations, err := filepath.Glob("Demonstrations/example*")
	if err != nil {
		t.Fatalf("Failed to find demonstration directories: %v", err)
	}

	for _, dir := range demonstrations {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			t.Fatalf("Failed to read directory %s: %v", dir, err)
		}

		var svgPath, rdPath string
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".svg") {
				svgPath = filepath.Join(dir, file.Name())
			}
			if strings.HasSuffix(file.Name(), ".rd") {
				rdPath = filepath.Join(dir, file.Name())
			}
		}

		if svgPath == "" || rdPath == "" {
			t.Fatalf("Could not find svg/rd pair in %s", dir)
		}

		t.Run(filepath.Base(dir), func(t *testing.T) {
			svgData, err := ioutil.ReadFile(svgPath)
			if err != nil {
				t.Fatalf("Failed to read SVG file %s: %v", svgPath, err)
			}

			expectedRdData, err := ioutil.ReadFile(rdPath)
			if err != nil {
				t.Fatalf("Failed to read RD file %s: %v", rdPath, err)
			}

			actualRdData := converter.Convert(svgData)

			if !bytes.Equal(expectedRdData, actualRdData) {
				t.Errorf("Converted RD data does not match expected for %s", filepath.Base(dir))
				// Optional: Write the actual output for debugging
				// ioutil.WriteFile(filepath.Join(dir, "actual.rd"), actualRdData, 0644)
			}
		})
	}
}
