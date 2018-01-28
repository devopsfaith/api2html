package generator

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestTmplScanner(t *testing.T) {
	pwd := os.Getenv("PWD")
	sourceFolder := fmt.Sprintf("%s/test/sources", pwd)

	iso := "en-US"
	paths := []string{
		fmt.Sprintf("%s/global", sourceFolder),
		fmt.Sprintf("%s/%s", sourceFolder, iso),
	}
	scanner := NewScanner(paths)

	tmpls := scanner.Scan()
	log.Println(tmpls)
	if len(tmpls) != 2 {
		t.Errorf("[%s] unexpected scan result. have %d nodes, want %d", iso, len(tmpls), 2)
		return
	}
	if tmpls[0].Path != paths[0] {
		t.Errorf("[%s - 0] unexpected path. have %s, want %s", iso, tmpls[0].Path, paths[0])
		return
	}
	if len(tmpls[0].Content) != 4 {
		t.Errorf("[%s - 0] unexpected content size. have %d, want 4", iso, len(tmpls[0].Content))
		return
	}
	if tmpls[1].Path != paths[1] {
		t.Errorf("[%s - 1] unexpected path. have %s, want %s", iso, tmpls[1].Path, paths[1])
		return
	}
	if len(tmpls[1].Content) != 1 {
		t.Errorf("[%s - 1] unexpected content size. have %d, want 1", iso, len(tmpls[1].Content))
		return
	}

	iso = "en-UK"
	paths = []string{
		fmt.Sprintf("%s/global", sourceFolder),
		fmt.Sprintf("%s/%s", sourceFolder, iso),
	}
	scanner = NewScanner(paths)

	tmpls = scanner.Scan()
	if len(tmpls) != 2 {
		t.Errorf("[%s] unexpected scan result. have %d nodes, want %d", iso, len(tmpls), 2)
		return
	}
	if tmpls[0].Path != paths[0] {
		t.Errorf("[%s - 0] unexpected path. have %s, want %s", iso, tmpls[0].Path, paths[0])
		return
	}
	if len(tmpls[0].Content) != 4 {
		t.Errorf("[%s - 0] unexpected content size. have %d, want 4", iso, len(tmpls[0].Content))
		return
	}
	if tmpls[1].Path != paths[1] {
		t.Errorf("[%s - 1] unexpected path. have %s, want %s", iso, tmpls[1].Path, paths[1])
		return
	}
	if len(tmpls[1].Content) != 1 {
		t.Errorf("[%s - 1] unexpected content size. have %d, want 1", iso, len(tmpls[1].Content))
		return
	}

	iso = "es-ES"
	paths = []string{
		fmt.Sprintf("%s/global", sourceFolder),
		fmt.Sprintf("%s/%s", sourceFolder, iso),
	}
	scanner = NewScanner(paths)

	tmpls = scanner.Scan()
	if len(tmpls) != 2 {
		t.Errorf("[%s] unexpected scan result. have %d nodes, want %d", iso, len(tmpls), 2)
		return
	}
	if tmpls[0].Path != paths[0] {
		t.Errorf("[%s - 0] unexpected path. have %s, want %s", iso, tmpls[0].Path, paths[0])
		return
	}
	if len(tmpls[0].Content) != 4 {
		t.Errorf("[%s - 0] unexpected content size. have %d, want 4", iso, len(tmpls[0].Content))
		return
	}
	if tmpls[1].Path != paths[1] {
		t.Errorf("[%s - 1] unexpected path. have %s, want %s", iso, tmpls[1].Path, paths[1])
		return
	}
	if len(tmpls[1].Content) != 1 {
		t.Errorf("[%s - 1] unexpected content size. have %d, want 1", iso, len(tmpls[1].Content))
		return
	}

	iso = "unknown"
	paths = []string{
		fmt.Sprintf("%s/global", sourceFolder),
		fmt.Sprintf("%s/%s", sourceFolder, iso),
	}
	scanner = NewScanner(paths)

	tmpls = scanner.Scan()
	if len(tmpls) != 1 {
		t.Errorf("[%s] unexpected scan result. have %d nodes, want %d", iso, len(tmpls), 2)
		return
	}
	if tmpls[0].Path != paths[0] {
		t.Errorf("[%s - 0] unexpected path. have %s, want %s", iso, tmpls[0].Path, paths[0])
		return
	}
	if len(tmpls[0].Content) != 4 {
		t.Errorf("[%s - 0] unexpected content size. have %d, want 4", iso, len(tmpls[0].Content))
		return
	}
}
