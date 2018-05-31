package generator

import (
	"fmt"
	"os"
	"testing"
)

func TestTmplScanner(t *testing.T) {

	tt := []struct {
		iso   string
		nodes int
	}{
		{"en-US", 2},
		{"en-UK", 2},
		{"es-ES", 2},
		{"unknown", 1},
	}

	pwd := os.Getenv("PWD")
	sourceFolder := fmt.Sprintf("%s/test/sources", pwd)

	for _, tc := range tt {
		t.Run(tc.iso, func(t *testing.T) {
			paths := []string{
				fmt.Sprintf("%s/global", sourceFolder),
				fmt.Sprintf("%s/%s", sourceFolder, tc.iso),
			}
			scanner := NewScanner(paths)

			tmpls := scanner.Scan()
			if len(tmpls) != tc.nodes {
				t.Errorf("[%s] unexpected scan result. have %d nodes, want %d", tc.iso, len(tmpls), tc.nodes)
				return
			}

			if tmpls[0].Path != paths[0] {
				t.Errorf("[%s - 0] unexpected path. have %s, want %s", tc.iso, tmpls[0].Path, paths[0])
				return
			}
			if len(tmpls[0].Content) != 4 {
				t.Errorf("[%s - 0] unexpected content size. have %d, want 4", tc.iso, len(tmpls[0].Content))
				return
			}
			if tc.nodes > 1 {
				if tmpls[1].Path != paths[1] {
					t.Errorf("[%s - 1] unexpected path. have %s, want %s", tc.iso, tmpls[1].Path, paths[1])
					return
				}
				if len(tmpls[1].Content) != 1 {
					t.Errorf("[%s - 1] unexpected content size. have %d, want 1", tc.iso, len(tmpls[1].Content))
					return
				}
			}
		})
	}
}
