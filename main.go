package main

// Package is called aw
import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	translate "cloud.google.com/go/translate/apiv3"
	aw "github.com/deanishe/awgo"
	translatepb "google.golang.org/genproto/googleapis/cloud/translate/v3"
)

const (
	envWorkflowID      = "alfred_workflow_bundleid"
	envWorkflowCache   = "alfred_workflow_cache"
	envWorkflowData    = "alfred_workflow_data"
	envWorkflowVersion = "alfred_workflow_version"
	envAlfredVersion   = "alfred_version"
)

var defaultSettings = map[string]string{
	envWorkflowID:      "com.keyneston.translate",
	envWorkflowCache:   filepath.Join(os.TempDir(), "alfred_cache"),
	envWorkflowData:    filepath.Join(os.TempDir(), "alfred_data"),
	envWorkflowVersion: "v0.0.1",
	envAlfredVersion:   "4",
}

func init() {
	for k, v := range defaultSettings {
		if envVar := os.Getenv(k); envVar == "" {
			os.Setenv(k, v)
		}
	}
}

type Workflow struct {
	*aw.Workflow
}

// Run wraps the underlying aw.Workflow.Run
func (w *Workflow) Run() {
	w.Workflow.Run(w.run)
}

// Your workflow starts here
func (w *Workflow) run() {
	if len(os.Args) < 2 {
		return
	}

	client, err := translate.NewTranslationClient(context.Background())
	if err != nil {
		w.Fatalf("error creating translate client: %v", err)
	}

	result, err := client.TranslateText(context.Background(),
		&translatepb.TranslateTextRequest{
			Contents: os.Args[1:],
		})
	if err != nil {
		w.Fatalf("error getting results: %v", err)
	}

	// Do thing here!
	for _, i := range result.Translations {
		w.NewItem(fmt.Sprintf("%s [%s]", i.TranslatedText, i.DetectedLanguageCode))
	}
	w.SendFeedback()
}

func main() {
	wf := &Workflow{aw.New()}
	wf.Run()
}
