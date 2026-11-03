package app

import (
	"strings"
	"testing"
)

func TestFullScriptNeverUsesStaleSegments(t *testing.T) {
	s := Script{InputMode: "text", Text: "First sentence. Second sentence! The complete ending?", Segments: []Segment{{Text: "Old preview only.", Voice: "alba"}}}
	if err := validateScript(&s); err != nil {
		t.Fatal(err)
	}
	parts := splitScript(s)
	if len(parts) != 3 || speechContent(segmentText(parts)) != speechContent(s.Text) {
		t.Fatal("full script was truncated", parts)
	}
	if len(s.Segments) != 0 {
		t.Fatal("stale chunks survived text-mode validation")
	}
}

func TestLegacyAmbiguousScriptsDoNotSilentlyTruncate(t *testing.T) {
	s := Script{Text: "First sentence. This ending must also be spoken.", Segments: []Segment{{Text: "First sentence."}}}
	if err := validateScript(&s); err == nil {
		t.Fatal("ambiguous legacy payload must be rejected")
	}
	s.InputMode = "segments"
	s.Segments = append(s.Segments, Segment{Text: "The edited ending."})
	if err := validateScript(&s); err != nil {
		t.Fatal(err)
	}
	if s.Text != "First sentence.\nThe edited ending." {
		t.Fatal("segment edits did not update the complete script", s.Text)
	}
	legacy := Script{Segments: s.Segments}
	if err := validateScript(&legacy); err != nil || legacy.InputMode != "segments" {
		t.Fatal("legacy segment-only clients must remain supported", err)
	}
}

func TestSentenceChunkingPreservesAllNonWhitespace(t *testing.T) {
	text := "... First?!\n" + strings.Repeat("Café 中文 words ", 120) + "THE FINAL ENDING!!!"
	s := Script{Text: text, Voice: "alba"}
	parts := splitScript(s)
	if speechContent(segmentText(parts)) != speechContent(text) {
		t.Fatal("chunking dropped text or punctuation")
	}
	for _, s := range parts {
		if len(s.Text) > 500 {
			t.Fatal("chunk exceeds CPU worker limit")
		}
	}
}
