package doctor

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestScanDeterministicAndClassifies(t *testing.T) {
	d := t.TempDir()
	os.WriteFile(filepath.Join(d, "b.json"), []byte(`{"schema":"capsule-v1","id":"b","status":"UNKNOWN","source":"x","falsifiers":["f"],"unknowns":["u"]}`), 0644)
	os.WriteFile(filepath.Join(d, "a.json"), []byte(`{"schema":"capsule-v1","id":"a"}`), 0644)
	x, e := Scan(d)
	if e != nil {
		t.Fatal(e)
	}
	y, _ := Scan(d)
	if x.Records[0].Path != "a.json" || x.Valid != 1 || x.Invalid != 1 || x.BenchmarkReady != 1 || x.UnknownFields != 1 {
		t.Fatalf("unexpected %+v", x)
	}
	if stringify(x) != stringify(y) {
		t.Fatal("nondeterministic")
	}
}
func stringify(r Report) string { b, _ := json.Marshal(r); return string(b) }
