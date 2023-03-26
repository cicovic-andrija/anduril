package yfm_test

import (
	"strings"
	"testing"

	"github.com/cicovic-andrija/anduril/anduril"
	"github.com/cicovic-andrija/anduril/yfm"
)

func TestParseExpectedInput(t *testing.T) {
	input := `---
tags: [diving, ssi]
title: SSI Peak Performance Buoyancy
created: '2023-03-01T20:01:32.854Z'
modified: '2023-03-03T17:34:50.396Z'
---`
	expectedTitle := "SSI Peak Performance Buoyancy"
	expectedTag0 := "diving"
	expectedTag1 := "ssi"
	expectedCreated := "2023-03-01T20:01:32.854Z"
	expectedModified := "2023-03-03T17:34:50.396Z"
	metadata := &anduril.ArticleMetadata{}
	err := yfm.Parse(strings.NewReader(input), metadata)
	if err != nil {
		t.Fatalf("parsing failed: %v", err)
	}
	if metadata.Title != expectedTitle {
		t.Fatalf("title: expected: %q found: %q", expectedTitle, metadata.Title)
	}
	if len(metadata.Tags) != 2 {
		t.Fatalf("tags: expected len: %d found: %d", 2, len(metadata.Tags))
	}
	if metadata.Tags[0] != expectedTag0 {
		t.Fatalf("tags[0]: expected: %q found: %q", expectedTag0, metadata.Tags[0])
	}
	if metadata.Tags[1] != expectedTag1 {
		t.Fatalf("tags[1]: expected: %q found: %q", expectedTag1, metadata.Tags[1])
	}
	if metadata.Created != expectedCreated {
		t.Fatalf("created: expected: %q found: %q", expectedCreated, metadata.Created)
	}
	if metadata.Modified != expectedModified {
		t.Fatalf("modified: expected: %q found: %q", expectedModified, metadata.Modified)
	}
}

func TestParseWithoutTags(t *testing.T) {
	input := `---
title: SSI Dry Suit
created: '2023-03-01T20:01:32.854Z'
modified: '2023-03-03T17:34:50.396Z'
---`
	metadata := &anduril.ArticleMetadata{}
	err := yfm.Parse(strings.NewReader(input), metadata)
	if err != nil {
		t.Fatalf("parsing failed: %v", err)
	}
	if metadata.Tags != nil || len(metadata.Tags) > 0 {
		t.Fatalf("tags: expected: 0 found: %d", len(metadata.Tags))
	}
	t.Logf("number of tags found: %d", len(metadata.Tags))
}
