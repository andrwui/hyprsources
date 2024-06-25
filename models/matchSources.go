package models

import (
	"fmt"
	"strings"
)

func MatchSources(jsonSources []SourceFile, hyprSources []string) SourceFilePointerArray {
	newSources := make([]*SourceFile, 0)

	if len(jsonSources) < 1 {
		fmt.Printf("No length in jsonSources")
		for _, hyprSource := range hyprSources {
			newSource := createSourceFromPath(hyprSource)
			newSources = append(newSources, newSource)
		}
	} else {
		for _, jsonSource := range jsonSources {
			newSources = append(newSources, &jsonSource)
		}
		for _, hyprSource := range hyprSources {
			if !findIfCreated(jsonSources, hyprSource) {
				newSource := createSourceFromPath(hyprSource)
				newSources = append(newSources, newSource)
			}
		}
	}
	return newSources
}

func createSourceFromPath(source string) *SourceFile {

	active := !strings.HasPrefix(strings.TrimSpace(source), "#")
	cleanPath := strings.TrimSpace(strings.Split(source, "=")[1])

	return &SourceFile{
		Name:     cleanPath,
		Path:     cleanPath,
		Selected: active,
		Priority: "low"}
}

func findIfCreated(jsonSources []SourceFile, source string) bool {
	for _, jsonSource := range jsonSources {
		if strings.TrimSpace(strings.Split(source, "=")[1]) == jsonSource.Path {
			return true
		}
	}
	return false
}
