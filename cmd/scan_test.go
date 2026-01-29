package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScan_FindsEnvFile(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, "test.env")
	if err := os.WriteFile(envFile, []byte("FOO=bar\n"), 0644); err != nil {
		t.Fatal(err)
	}

	files, err := Scan(dir, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}
	if files[0] != envFile {
		t.Fatalf("expected %s, got %s", envFile, files[0])
	}
}

func TestScan_WalksUpward(t *testing.T) {
	parent := t.TempDir()
	child := filepath.Join(parent, "subdir")
	if err := os.Mkdir(child, 0755); err != nil {
		t.Fatal(err)
	}
	envFile := filepath.Join(parent, "test.env")
	if err := os.WriteFile(envFile, []byte("FOO=bar\n"), 0644); err != nil {
		t.Fatal(err)
	}

	files, err := Scan(child, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}
	if files[0] != envFile {
		t.Fatalf("expected %s, got %s", envFile, files[0])
	}
}

func TestScan_FiltersByEnvironment(t *testing.T) {
	dir := t.TempDir()
	testEnv := filepath.Join(dir, "test.env")
	prodEnv := filepath.Join(dir, "prod.env")
	if err := os.WriteFile(testEnv, []byte("FOO=bar\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(prodEnv, []byte("FOO=baz\n"), 0644); err != nil {
		t.Fatal(err)
	}

	files, err := Scan(dir, "test")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}
	if files[0] != testEnv {
		t.Fatalf("expected %s, got %s", testEnv, files[0])
	}
}
