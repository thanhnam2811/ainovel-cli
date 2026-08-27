package assets

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestVietnameseGenrePresetsAreLocalized(t *testing.T) {
	cases := []struct {
		style       string
		styleMarker string
		refMarker   string
		arcMarker   string
	}{
		{"fantasy", "## Phong cách fantasy / phiêu lưu", "# Tham chiếu bổ sung — fantasy / xianxia", "## Tham chiếu arc — fantasy / xianxia"},
		{"romance", "## Phong cách romance / tình cảm", "# Tham chiếu bổ sung — romance / tình cảm", "## Tham chiếu arc — romance / tình cảm"},
		{"suspense", "## Phong cách suspense / trinh thám", "# Tham chiếu bổ sung — suspense / trinh thám", "## Tham chiếu arc — suspense / trinh thám"},
	}

	for _, tc := range cases {
		t.Run(tc.style, func(t *testing.T) {
			b := Load(tc.style, LoadOptions{})
			if err := ApplyLocaleForStyle(&b, "vi", tc.style); err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(b.Styles[tc.style], tc.styleMarker) {
				t.Fatalf("style preset %s was not localized", tc.style)
			}
			if !strings.Contains(b.References.StyleReference, tc.refMarker) {
				t.Fatalf("genre style reference %s was not localized", tc.style)
			}
			if !strings.Contains(b.References.ArcTemplates, tc.arcMarker) {
				t.Fatalf("genre arc templates %s were not localized", tc.style)
			}
		})
	}
}

func TestVietnameseGenreReferencePreservesUserOverridePrecedence(t *testing.T) {
	home := t.TempDir()
	book := t.TempDir()
	rel := filepath.Join("genres", "fantasy", "style-references.md")

	write := func(root, text string) {
		t.Helper()
		path := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	write(home, "GLOBAL GENRE OVERRIDE")
	write(book, "BOOK GENRE OVERRIDE")
	opts := LoadOptions{HomeStyleDir: home, BookStyleDir: book}
	b := Load("fantasy", opts)
	if err := ApplyLocaleForStyle(&b, "vi", "fantasy", opts); err != nil {
		t.Fatal(err)
	}
	if b.References.StyleReference != "BOOK GENRE OVERRIDE" {
		t.Fatalf("book genre override must win, got %q", b.References.StyleReference)
	}
}

func TestVietnameseDefaultStyleDoesNotInventGenreReferences(t *testing.T) {
	b := Load("default", LoadOptions{})
	if err := ApplyLocaleForStyle(&b, "vi", "default"); err != nil {
		t.Fatal(err)
	}
	if b.References.StyleReference != "" || b.References.ArcTemplates != "" {
		t.Fatal("default style must not invent genre-specific references")
	}
}
