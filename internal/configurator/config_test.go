package configurator

import (
	"testing"
)

var goodEmptyConfigExample = "{}"

var goodShortConfigExample = `
{
  "preset": "psr12"
}
`

var goodFullConfigExample = `
{
  "preset": "psr12",
  "exclude": [
    "tests/e2e"
  ],
  "notName": [
    "*-test.php"
  ],
  "notPath": [
    "app/conf/user/maps/calc.php"
  ],
  "rules": {
    "simplified_null_return": true,
    "rule_not_found": true,
    "new_with_braces": {
      "anonymous_class": false,
      "named_class": false
    }
  }
}
`

var badFullConfigExample = `
{
  "preset": "psr12"
  "exclude": [
    "tests/e2e"
  ],
  "notName": [
    "*-test.php"
  ],
  "notPath": [
    "app/conf/user/maps/calc.php"
  ],
  "rules": {
    "simplified_null_return": true,
    "rule_not_found": true,
    "new_with_braces": {
      "anonymous_class": false,
      "named_class": false
    }
  }
}
`

func TestParsingEmptyConfig(t *testing.T) {
	config, err := Parse(goodEmptyConfigExample)

	if err != nil {
		t.Fatalf("Parsing %v failed: %v", goodEmptyConfigExample, err)
	}

	if config.Preset != defaultPreset {
		t.Fatalf("Parsing %v failed: got %v, want %v", goodEmptyConfigExample, config.Preset, defaultPreset)
	}

}

func TestParsingShortConfig(t *testing.T) {
	config, err := Parse(goodShortConfigExample)

	if err != nil {
		t.Fatalf("Parsing %v failed: %v", goodShortConfigExample, err)
	}

	if config.Preset != "psr12" {
		t.Fatalf("Parsing %v failed: got %v, want %v", goodShortConfigExample, config.Preset, defaultPreset)
	}

	if config.GetRules() != nil {
		t.Fatalf("Parsing %v failed: got %v, want %v", goodShortConfigExample, config.GetRules(), nil)
	}
}

func TestParsingFullConfig(t *testing.T) {
	config, err := Parse(goodFullConfigExample)
	if err != nil {
		t.Fatalf("Parsing %v failed: %v", goodFullConfigExample, err)
	}
	if len(config.GetRules()) != 3 {
		t.Fatalf("Parsing %v failed: got %v, want %v", goodFullConfigExample, len(config.GetRules()), 3)
	}
}

func TestParsingError(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Fatal(err)
		}
	}()

	_, err := Parse(badFullConfigExample)

	if err == nil {
		t.Fatalf("Parsing %v failed: got nil, want error", badFullConfigExample)
	}

}
