package main

// Package is called aw
import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"

	translate "cloud.google.com/go/translate/apiv3"
	aw "github.com/deanishe/awgo"
	"golang.org/x/text/language"
	translatepb "google.golang.org/genproto/googleapis/cloud/translate/v3"
)

const (
	envWorkflowID      = "alfred_workflow_bundleid"
	envWorkflowCache   = "alfred_workflow_cache"
	envWorkflowData    = "alfred_workflow_data"
	envWorkflowVersion = "alfred_workflow_version"
	envAlfredVersion   = "alfred_version"

	googleCredentials = "GOOGLE_APPLICATION_CREDENTIALS" // hack until I get settings working
)

// config keys
const (
	KeyDefaultTargetLanguage = "default_target_language"
	KeyGoogleProjectID       = "google_project_id"
)

var defaultSettings = map[string]string{
	googleCredentials:  "/Users/tabitha/.google/translate.json",
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

func (w Workflow) GetProjectID() string {
	return path.Join("projects", w.Config.GetString(KeyGoogleProjectID))
}

// Run wraps the underlying aw.Workflow.Run
func (w *Workflow) Run() {
	w.Workflow.Run(w.run)
}

// Your workflow starts here
func (w *Workflow) run() {
	ctx := context.Background()
	if len(os.Args) < 2 {
		return
	}

	targetLanguage := w.Config.GetString(KeyDefaultTargetLanguage, "en")

	_, err := language.Parse(targetLanguage) // Make sure language code is valid
	if err != nil {
		w.Fatalf("error validating language %q: %v", targetLanguage, err)
	}

	client, err := translate.NewTranslationClient(ctx)
	if err != nil {
		w.Fatalf("error creating translate client: %v", err)
	}
	defer client.Close()

	req := &translatepb.TranslateTextRequest{
		// SourceLanguageCode: sourceLang.String(),
		TargetLanguageCode: targetLanguage,
		MimeType:           "text/plain",
		Contents:           os.Args[1:],
		Parent:             w.GetProjectID(),
	}

	client.TranslateText(ctx, req)
	result, err := client.TranslateText(context.Background(), req)
	if err != nil {
		w.Fatalf("error getting results: %v", err)
	}

	// Do thing here!
	for _, i := range result.Translations {
		result := i.TranslatedText
		item := w.NewItem(result)
		if i.DetectedLanguageCode != "" {
			item.Subtitle(fmt.Sprintf("[%s=>%s]", i.DetectedLanguageCode, targetLanguage))
		}
		item.Copytext(i.TranslatedText)
		item.Var("translation", i.TranslatedText)
		item.Valid(true)
		item.Arg(i.TranslatedText)
	}
	w.SendFeedback()
}

func main() {
	wf := &Workflow{aw.New()}
	wf.Run()
}
