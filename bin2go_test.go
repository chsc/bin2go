package main

import (
	"testing"
)

var camelCaseTT = []struct {
	input string
	want  string
}{
	// Regular
	{"/regular/root/path/file.ext", "regularRootPathFileExt"},
	{"regular/relative/path/file.ext", "regularRelativePathFileExt"},

	// Already in a form of casing
	{"camelCased", "camelCased"},
	{"TitledCamelCased", "titledCamelCased"},
	{"snake_cased", "snakeCased"},

	// Special cases
	{".startsWithDot", "startsWithDot"},
	{"._starts_with_dot_underscore", "startsWithDotUnderscore"},
	{"finishes_with_dot_underscore._", "finishesWithDotUnderscore"},
	{"._wrapped_with_dot_underscore._", "wrappedWithDotUnderscore"},
	{"has space in it", "hasSpaceInIt"},
	{"_non_char+in/_path\\/_file+/_filename..ext/_file_", "nonCharInPathFileFilenameExtFile"},

	// Unicode
	{"日本語_nihõŋɡo", "日本語Nihõŋɡo"},
	{"Ægis.époustouflant", "ægisÉpoustouflant"},
	{"ιι_κκ_λλ_μμ/ιι_κκ_λλ_μμ/ιι_κκ_λλ_μμ.ππππ", "ιιΚκΛλΜμΙιΚκΛλΜμΙιΚκΛλΜμΠπππ"},

	// Unicode with non-digits, non-letters
	{"ιι;κκ;λλ;μμ·ιι;κκ;λλ;μμ·ιι;κκ;λλ;μμ;ππππ", "ιιΚκΛλΜμΙιΚκΛλΜμΙιΚκΛλΜμΠπππ"},
}

func TestCamelCase(t *testing.T) {
	for _, testCase := range camelCaseTT {
		got := camelCase(testCase.input)
		if testCase.want != got {
			t.Errorf("For input '%s', want '%s' got '%s'", testCase.input, testCase.want, got)
		}
	}
}
