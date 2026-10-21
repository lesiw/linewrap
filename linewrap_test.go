package linewrap

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

//go:embed testdata
var testdata embed.FS

// testdataDir materializes the embedded testdata into a fresh temp
// directory, restoring the .go extensions the fixtures are stored
// without. Fixtures live as .input in the repo so the formatter never
// touches them, but analysistest needs real .go files to parse. The
// returned dir holds the same src/ layout analysistest expects.
func testdataDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	walk := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		content, err := testdata.ReadFile(path)
		if err != nil {
			return err
		}
		// Drop the testdata/ prefix; the embedded tree already carries src/.
		repoPath, _ := strings.CutPrefix(path, "testdata/")
		name, _ := strings.CutSuffix(repoPath, ".input")
		out := filepath.Join(dir, name+".go")
		if err := os.MkdirAll(filepath.Dir(out), 0777); err != nil {
			return err
		}
		return os.WriteFile(out, content, 0666)
	}
	if err := fs.WalkDir(testdata, ".", walk); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestAnalysisTest(t *testing.T) {
	analysistest.Run(t, testdataDir(t), Analyzer, "a")
}
