package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidate_AcceptsGoryukFile(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, "test.env")
	if err := os.WriteFile(envFile, []byte("GORYUK=TRUE\nFOO=bar\n"), 0644); err != nil {
		t.Fatal(err)
	}

	result, err := Validate([]string{envFile})
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 valid env, got %d", len(result))
	}
	if result[0].file != envFile {
		t.Fatalf("expected %s, got %s", envFile, result[0].file)
	}
}

func TestValidate_RejectsPlainEnvFile(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, "test.env")
	if err := os.WriteFile(envFile, []byte("FOO=bar\nBAZ=qux\n"), 0644); err != nil {
		t.Fatal(err)
	}

	result, err := Validate([]string{envFile})
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 0 {
		t.Fatalf("expected 0 valid envs, got %d", len(result))
	}
}
