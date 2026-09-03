package blueprint

import (
	"strings"
	"testing"
)

const valid = `{"schema":"story-factory.narrative-blueprint.v1","story":{"name":"X"},"narrative_design":{},"pacing_contract":{"depth_target":"fast_commercial","chapter_word_range":[900,1300],"scene_count_range":[2,3]},"volumes":[{"index":1,"title":"V1","theme":"T","arcs":[{"index":1,"title":"A1","goal":"G","chapters":[{"chapter":1,"title":"C1","core_event":"30 days before: buy the RV; payoff BASE_BUILDING","hook":"The system wakes","scenes":["Rebirth","Buy RV"]}]}]}]}`

func TestParseRecognizesStoryFactoryBlueprint(t *testing.T) {
	doc, recognized, err := Parse(valid)
	if err != nil || !recognized {
		t.Fatalf("Parse() recognized=%v err=%v", recognized, err)
	}
	if got := doc.Volumes[0].Arcs[0].Chapters[0].CoreEvent; !strings.Contains(got, "buy the RV") {
		t.Fatalf("core event was not preserved: %q", got)
	}
	if task := CanonicalPlannerTask(valid); !strings.Contains(task, valid) || !strings.Contains(task, "Do not redesign") {
		t.Fatalf("planner task did not preserve source: %s", task)
	}
}

func TestParseLeavesNaturalLanguagePromptsAlone(t *testing.T) {
	if doc, recognized, err := Parse("Write a survival novel"); err != nil || recognized || doc != nil {
		t.Fatalf("Parse() = %#v, %v, %v", doc, recognized, err)
	}
}

func TestParseRejectsBrokenChapterSequence(t *testing.T) {
	broken := strings.Replace(valid, `"chapter":1`, `"chapter":2`, 1)
	if _, recognized, err := Parse(broken); !recognized || err == nil {
		t.Fatalf("Parse() recognized=%v err=%v", recognized, err)
	}
}
