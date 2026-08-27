package assets

import (
	"crypto/sha1"
	"fmt"
	"testing"
)

// TestVietnameseLocalizedAssetsMatchReviewedUpstream catches the dangerous case
// where an upstream merge changes a translated prompt/reference/style while the
// Vietnamese overlay silently stays stale. The expected values are Git blob
// SHAs from the upstream assets at the time the Vietnamese version was reviewed.
//
// A failure is not fixed by blindly updating the SHA: review the upstream diff,
// port behavior/meaning into assets/locales/vi, update contract tests if needed,
// then pin the new reviewed blob SHA.
func TestVietnameseLocalizedAssetsMatchReviewedUpstream(t *testing.T) {
	checks := []struct {
		name string
		read func() ([]byte, error)
		sha  string
	}{
		// Runtime prompts.
		{"prompts/architect-short.md", func() ([]byte, error) { return promptsFS.ReadFile("prompts/architect-short.md") }, "1656a6d07183fc8e9efe932f01095e982ea5eaf8"},
		{"prompts/architect-long.md", func() ([]byte, error) { return promptsFS.ReadFile("prompts/architect-long.md") }, "31486b98691fca066b99b1c0dcc24e34bbf6ae7c"},
		{"prompts/writer.md", func() ([]byte, error) { return promptsFS.ReadFile("prompts/writer.md") }, "e43842383503783fd697e225ed6f3f77365ff8b9"},
		{"prompts/editor.md", func() ([]byte, error) { return promptsFS.ReadFile("prompts/editor.md") }, "6ffd6047ceaef96e444cd9817bc2ec4bf2b4dc63"},
		{"prompts/import-segment.md", func() ([]byte, error) { return promptsFS.ReadFile("prompts/import-segment.md") }, "74331a321a7fee86b43f8927b9fb13eb0da83431"},
		{"prompts/import-analyze.md", func() ([]byte, error) { return promptsFS.ReadFile("prompts/import-analyze.md") }, "5126b49da36a651b80502d8b8100855ef1fc09c9"},
		{"prompts/import-synthesize.md", func() ([]byte, error) { return promptsFS.ReadFile("prompts/import-synthesize.md") }, "ac5111303b180d63835eac468655f712486abc4b"},
		{"prompts/import-range.md", func() ([]byte, error) { return promptsFS.ReadFile("prompts/import-range.md") }, "a852a889d2e50fbc0ddc64ea622ed4af33d195e9"},
		{"prompts/simulation-source.md", func() ([]byte, error) { return promptsFS.ReadFile("prompts/simulation-source.md") }, "9b0f322f4e88824edbb2100087f70ce8a32ad959"},
		{"prompts/simulation-merge.md", func() ([]byte, error) { return promptsFS.ReadFile("prompts/simulation-merge.md") }, "369ba1688d3456dd5b1aec5a91f2a26061d48955"},
		{"prompts/revision-analyze.md", func() ([]byte, error) { return promptsFS.ReadFile("prompts/revision-analyze.md") }, "bcf8585619cb9a2dfa0effdda2f612f35353b94d"},
		{"prompts/arbiter-plan-start.md", func() ([]byte, error) { return promptsFS.ReadFile("prompts/arbiter-plan-start.md") }, "0200b759f5cc4ac0f4f8faf49f80021b6b30b910"},
		{"prompts/arbiter-intervention.md", func() ([]byte, error) { return promptsFS.ReadFile("prompts/arbiter-intervention.md") }, "6e2d65b4cf2e496217b1d89bca204f99369881d2"},
		{"prompts/arbiter-failure.md", func() ([]byte, error) { return promptsFS.ReadFile("prompts/arbiter-failure.md") }, "3a05c5700f70ed49df6ee706783e24ebb246bc18"},

		// Generic runtime reference pack.
		{"references/anti-ai-tone.md", func() ([]byte, error) { return referencesFS.ReadFile("references/anti-ai-tone.md") }, "5629d95ccaae2bfc5c42429b6dce370fad3b8816"},
		{"references/chapter-guide.md", func() ([]byte, error) { return referencesFS.ReadFile("references/chapter-guide.md") }, "8b1e1c34500f1b5b2002617ecaaef26abc15d457"},
		{"references/chapter-template.md", func() ([]byte, error) { return referencesFS.ReadFile("references/chapter-template.md") }, "c0b452c4b494357b97130d78598f155a2ad1792e"},
		{"references/character-template.md", func() ([]byte, error) { return referencesFS.ReadFile("references/character-template.md") }, "42ade51fdbee51caeeabb2f9c23294e28bf45079"},
		{"references/consistency.md", func() ([]byte, error) { return referencesFS.ReadFile("references/consistency.md") }, "24072fc0103e5bc4b066de824bc19fbd220fc77c"},
		{"references/content-expansion.md", func() ([]byte, error) { return referencesFS.ReadFile("references/content-expansion.md") }, "d2ec2a5c2279ec865834dbc5a0d4b801d9364cc5"},
		{"references/dialogue-writing.md", func() ([]byte, error) { return referencesFS.ReadFile("references/dialogue-writing.md") }, "1f1fd9ec143add5ee36a58dd2e5df6e5dd87269c"},
		{"references/differentiation.md", func() ([]byte, error) { return referencesFS.ReadFile("references/differentiation.md") }, "aecbd111e5769f6c90b8127fba79fbee9b7967a0"},
		{"references/hook-techniques.md", func() ([]byte, error) { return referencesFS.ReadFile("references/hook-techniques.md") }, "12958e804cc1663118dbb7d25ad793be0d653abf"},
		{"references/longform-planning.md", func() ([]byte, error) { return referencesFS.ReadFile("references/longform-planning.md") }, "a26e7aafbcfbd7b89c920733790a63eb61b76ffb"},
		{"references/outline-template.md", func() ([]byte, error) { return referencesFS.ReadFile("references/outline-template.md") }, "7efc73fc938432d0a5a12db651f0cdf3fc8a9649"},
		{"references/quality-checklist.md", func() ([]byte, error) { return referencesFS.ReadFile("references/quality-checklist.md") }, "057ddb2191497f1c57e85bd1cb6006185b0c671f"},

		// Built-in style/voice assets.
		{"voice.md", func() ([]byte, error) { return voiceFS.ReadFile("voice.md") }, "fbc10db37570118d2b2453d9f625968badac024e"},
		{"styles/default.md", func() ([]byte, error) { return stylesFS.ReadFile("styles/default.md") }, "50d569c1620bfe71e076acc7f8765fcc54462280"},
		{"styles/fantasy.md", func() ([]byte, error) { return stylesFS.ReadFile("styles/fantasy.md") }, "05a333ed3b09f196428ea2817b6e26a298d54520"},
		{"styles/romance.md", func() ([]byte, error) { return stylesFS.ReadFile("styles/romance.md") }, "438ec7202328bb9450835e836ee47dda62a82ddd"},
		{"styles/suspense.md", func() ([]byte, error) { return stylesFS.ReadFile("styles/suspense.md") }, "d108470bdf90065c789791e24e45f633588dbd1a"},

		// Genre-specific runtime references.
		{"references/genres/fantasy/arc-templates.md", func() ([]byte, error) { return referencesFS.ReadFile("references/genres/fantasy/arc-templates.md") }, "e76c7f6b39090ab5f157fcd5d13d7718e69a4ff9"},
		{"references/genres/fantasy/style-references.md", func() ([]byte, error) { return referencesFS.ReadFile("references/genres/fantasy/style-references.md") }, "c81ff462dd556827be6216f1bb6fe0a7c9268373"},
		{"references/genres/romance/arc-templates.md", func() ([]byte, error) { return referencesFS.ReadFile("references/genres/romance/arc-templates.md") }, "2e438aa3382d0635fc9b746d7c5a64be95ab3ae2"},
		{"references/genres/romance/style-references.md", func() ([]byte, error) { return referencesFS.ReadFile("references/genres/romance/style-references.md") }, "49bdb997cdb14e776855463fd93a9d021b9fb2e4"},
		{"references/genres/suspense/arc-templates.md", func() ([]byte, error) { return referencesFS.ReadFile("references/genres/suspense/arc-templates.md") }, "45078e3d396ae4f0ad0e38bc0cbf44b05fbacdf5"},
		{"references/genres/suspense/style-references.md", func() ([]byte, error) { return referencesFS.ReadFile("references/genres/suspense/style-references.md") }, "5eb776bdc0a820471c0c7196de422dd4f1a2321c"},
	}

	for _, tc := range checks {
		t.Run(tc.name, func(t *testing.T) {
			data, err := tc.read()
			if err != nil {
				t.Fatal(err)
			}
			if got := gitBlobSHA(data); got != tc.sha {
				t.Fatalf("upstream asset drifted: got blob %s, reviewed %s; review upstream diff and port it into the Vietnamese overlay before updating this pin", got, tc.sha)
			}
		})
	}
}

func gitBlobSHA(data []byte) string {
	header := []byte(fmt.Sprintf("blob %d\x00", len(data)))
	h := sha1.New() // Git blob object ID, not a security primitive.
	_, _ = h.Write(header)
	_, _ = h.Write(data)
	return fmt.Sprintf("%x", h.Sum(nil))
}
