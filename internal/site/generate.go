package site

import (
	_ "embed"
	"html/template"
	"os"
	"path/filepath"

	"github.com/davidcollom/zwift-ical/internal/events"
)

//go:embed templates/index.html.tmpl
var indexTemplateContent string

func RenderIndex(events []events.Event, outputPath string) error {
	tmpl, err := template.New("index").Parse(indexTemplateContent)
	if err != nil {
		return err
	}
	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()
	return tmpl.Execute(f, events)
}

func RenderIndexLinks(paths []string, outputPath string) error {
	tmpl, err := template.New("index_links").Parse(indexTemplateContent)
	if err != nil {
		return err
	}
	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()
	return tmpl.Execute(f, paths)
}

func WriteICal(icalData, outputPath string) error {
	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(icalData)
	return err
}

func GenerateSite(events []events.Event, icalData string, publicDir string) error {
	indexOutput := filepath.Join(publicDir, "index.html")
	icalOutput := filepath.Join(publicDir, "zwift.ics")
	if err := RenderIndex(events, indexOutput); err != nil {
		return err
	}
	if err := WriteICal(icalData, icalOutput); err != nil {
		return err
	}
	return nil
}
