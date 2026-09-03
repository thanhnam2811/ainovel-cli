package main

import "testing"

func TestParseCLIOptionsAcceptsHeadlessChapterLimit(t *testing.T) {
	opts, args, err := parseCLIOptions([]string{"--headless", "--max-chapters", "3", "--prompt", "outline"})
	if err != nil {
		t.Fatal(err)
	}
	if len(args) != 0 || !opts.Headless || opts.MaxChapters != 3 || opts.Prompt != "outline" {
		t.Fatalf("unexpected parse result: opts=%+v args=%v", opts, args)
	}
}

func TestParseCLIOptionsRejectsInvalidChapterLimit(t *testing.T) {
	for _, argv := range [][]string{
		{"--headless", "--max-chapters", "0"},
		{"--headless", "--max-chapters", "many"},
		{"--max-chapters", "1"},
	} {
		if _, _, err := parseCLIOptions(argv); err == nil {
			t.Fatalf("expected error for %v", argv)
		}
	}
}
