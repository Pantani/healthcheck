package fixtures

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetFixtures(t *testing.T) {
	got, err := GetFixtures("testdata/fixtures.json")
	if err != nil {
		t.Fatalf("GetFixtures() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len(fixtures) = %d, want 1", len(got))
	}
	if got[0].Namespace != "ethereum" {
		t.Errorf("namespace = %q", got[0].Namespace)
	}
	if got.CheckCount() != 1 {
		t.Errorf("CheckCount() = %d, want 1", got.CheckCount())
	}
	if got[0].Tests[0].Method != "GET" {
		t.Errorf("method = %q, want GET", got[0].Tests[0].Method)
	}
}

func TestGetFixturesValidationError(t *testing.T) {
	path := writeFixture(t, `[
  {
    "namespace": "api",
    "host": "localhost:8080",
    "tests": [
      {
        "name": "status",
        "method": "GET",
        "url_path": "health",
        "json_path": "ok",
        "expression": "newValue == true",
        "update_time": "5s"
      }
    ]
  }
]`)

	_, err := GetFixtures(path)
	if err == nil {
		t.Fatal("GetFixtures() error = nil, want validation error")
	}
	if !strings.Contains(err.Error(), "must use http or https") {
		t.Fatalf("GetFixtures() error = %v", err)
	}
}

func TestTestValidateRejectsBadDuration(t *testing.T) {
	test := Test{
		Name:       "status",
		Method:     "GET",
		JSONPath:   "ok",
		Expression: "newValue == true",
		UpdateTime: "soon",
	}

	err := test.Validate()
	if err == nil {
		t.Fatal("Validate() error = nil, want invalid duration")
	}
	if !strings.Contains(err.Error(), "invalid update_time") {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestLoadFixturesMissingFile(t *testing.T) {
	var got Fixtures
	err := loadFixtures("missing.json", &got)
	if err == nil {
		t.Fatal("loadFixtures() error = nil, want error")
	}
}

func writeFixture(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "fixtures.json")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	return path
}
