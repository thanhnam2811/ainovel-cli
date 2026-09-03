package host

import (
	"strings"
	"testing"

	"github.com/voocel/ainovel-cli/internal/blueprint"
	"github.com/voocel/ainovel-cli/internal/domain"
	storepkg "github.com/voocel/ainovel-cli/internal/store"
)

const storyFactoryBlueprint = `{"schema":"story-factory.narrative-blueprint.v1","story":{"name":"X"},"narrative_design":{},"pacing_contract":{"depth_target":"fast_commercial","chapter_word_range":[900,1300],"scene_count_range":[2,3]},"volumes":[{"index":1,"title":"V1","theme":"T","arcs":[{"index":1,"title":"A1","goal":"G","chapters":[{"chapter":1,"title":"C1","core_event":"30 days before: buy the RV; payoff BASE_BUILDING","hook":"The system wakes","scenes":["Rebirth","Buy RV"]}]}]}]}`

func TestSeedNarrativeBlueprintPersistsApprovedOutline(t *testing.T) {
	doc, _, err := blueprint.Parse(storyFactoryBlueprint)
	if err != nil {
		t.Fatal(err)
	}
	st := storepkg.NewStore(t.TempDir())
	if err := st.Init(); err != nil {
		t.Fatal(err)
	}
	if err := st.Progress.Init(0); err != nil {
		t.Fatal(err)
	}

	if err := seedNarrativeBlueprint(st, doc); err != nil {
		t.Fatal(err)
	}
	outline, err := st.Outline.LoadLayeredOutline()
	if err != nil || !blueprint.SameOutline(outline, doc.Volumes) {
		t.Fatalf("outline was not preserved: %#v err=%v", outline, err)
	}
	progress, err := st.Progress.Load()
	if err != nil {
		t.Fatal(err)
	}
	if progress.Phase != domain.PhaseOutline || !progress.Layered || progress.TotalChapters != 1 {
		t.Fatalf("progress = %#v", progress)
	}
}

func TestPreserveRawRequirementKeepsCanonicalInput(t *testing.T) {
	raw := "chapter 1 must begin 30 days before the apocalypse"
	task := preserveRawRequirement(raw, "summarized task")
	if !strings.HasPrefix(task, raw) || !strings.Contains(task, "summarized task") {
		t.Fatalf("task = %q", task)
	}
}
