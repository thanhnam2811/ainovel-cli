// Package blueprint recognizes and validates typed narrative handoffs from Story Factory.
package blueprint

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/voocel/ainovel-cli/internal/domain"
)

const Schema = "story-factory.narrative-blueprint.v1"

type PacingContract struct {
	DepthTarget                           string `json:"depth_target"`
	ChapterWordRange                      []int  `json:"chapter_word_range"`
	SceneCountRange                       []int  `json:"scene_count_range"`
	PrimaryPayoffsPerChapter              int    `json:"primary_payoffs_per_chapter"`
	PayoffDeadlinePercent                 int    `json:"payoff_deadline_percent"`
	MaxExplanatoryParagraphsWithoutAction int    `json:"max_explanatory_paragraphs_without_action"`
	TechnicalDetailRule                   string `json:"technical_detail_rule"`
}

type Document struct {
	Schema          string                 `json:"schema"`
	Story           map[string]any         `json:"story"`
	NarrativeDesign map[string]any         `json:"narrative_design"`
	PacingContract  PacingContract         `json:"pacing_contract"`
	Volumes         []domain.VolumeOutline `json:"volumes"`
}

// Parse returns recognized=false for ordinary natural-language prompts.
func Parse(raw string) (*Document, bool, error) {
	var header struct {
		Schema string `json:"schema"`
	}
	if err := json.Unmarshal([]byte(raw), &header); err != nil || header.Schema == "" {
		return nil, false, nil
	}
	if header.Schema != Schema {
		return nil, true, fmt.Errorf("unsupported narrative blueprint schema %q", header.Schema)
	}
	var doc Document
	if err := json.Unmarshal([]byte(raw), &doc); err != nil {
		return nil, true, fmt.Errorf("decode narrative blueprint: %w", err)
	}
	if err := doc.Validate(); err != nil {
		return nil, true, err
	}
	return &doc, true, nil
}

func (d *Document) Validate() error {
	if len(d.Volumes) == 0 {
		return fmt.Errorf("narrative blueprint requires at least one volume")
	}
	if d.PacingContract.DepthTarget == "" {
		return fmt.Errorf("narrative blueprint requires pacing_contract.depth_target")
	}
	expectedChapter := 1
	for vi, volume := range d.Volumes {
		if volume.Index != vi+1 || strings.TrimSpace(volume.Title) == "" || strings.TrimSpace(volume.Theme) == "" {
			return fmt.Errorf("narrative blueprint volume %d is invalid", vi+1)
		}
		if len(volume.Arcs) == 0 {
			return fmt.Errorf("narrative blueprint volume %d has no arcs", volume.Index)
		}
		for ai, arc := range volume.Arcs {
			if arc.Index != ai+1 || strings.TrimSpace(arc.Title) == "" || strings.TrimSpace(arc.Goal) == "" {
				return fmt.Errorf("narrative blueprint volume %d arc %d is invalid", volume.Index, ai+1)
			}
			if len(arc.Chapters) == 0 {
				return fmt.Errorf("narrative blueprint volume %d arc %d has no chapters", volume.Index, arc.Index)
			}
			for _, chapter := range arc.Chapters {
				if chapter.Chapter != expectedChapter {
					return fmt.Errorf("narrative blueprint chapter sequence: got %d, want %d", chapter.Chapter, expectedChapter)
				}
				if strings.TrimSpace(chapter.Title) == "" || strings.TrimSpace(chapter.CoreEvent) == "" || strings.TrimSpace(chapter.Hook) == "" {
					return fmt.Errorf("narrative blueprint chapter %d requires title, core_event, and hook", chapter.Chapter)
				}
				if len(chapter.Scenes) < 2 || len(chapter.Scenes) > 3 {
					return fmt.Errorf("narrative blueprint chapter %d requires 2-3 scenes", chapter.Chapter)
				}
				expectedChapter++
			}
		}
	}
	return nil
}

func CanonicalPlannerTask(raw string) string {
	return strings.TrimSpace(raw) + "\n\n" + `MANDATORY HANDOFF RULES:
- The JSON blueprint above is canonical source data, not optional inspiration.
- Its volumes/arcs/chapters have already been approved. Do not redesign, summarize, reorder, skip, or replace them.
- Generate the remaining book metadata and foundation around this outline.
- Do not call save_foundation with outline or layered_outline; the host has persisted it deterministically.
- Preserve pacing_contract in every chapter plan and draft.`
}

func SameOutline(got []domain.VolumeOutline, want []domain.VolumeOutline) bool {
	gotJSON, gotErr := json.Marshal(got)
	wantJSON, wantErr := json.Marshal(want)
	return gotErr == nil && wantErr == nil && string(gotJSON) == string(wantJSON)
}
