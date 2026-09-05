package doctor

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Record struct {
	Path           string   `json:"path"`
	ID             string   `json:"id"`
	Schema         string   `json:"schema"`
	Status         string   `json:"status"`
	Unknowns       []string `json:"unknowns"`
	Falsifiers     []string `json:"falsifiers"`
	Provenance     bool     `json:"provenance"`
	BenchmarkReady bool     `json:"benchmark_ready"`
	Errors         []string `json:"errors"`
}
type Report struct {
	Schema         string   `json:"schema"`
	Root           string   `json:"root"`
	Files          int      `json:"files"`
	Valid          int      `json:"valid"`
	Invalid        int      `json:"invalid"`
	BenchmarkReady int      `json:"benchmark_ready"`
	UnknownFields  int      `json:"unknown_fields"`
	Records        []Record `json:"records"`
}

func check(path, root string) Record {
	r := Record{Path: filepath.ToSlash(strings.TrimPrefix(path, root+string(os.PathSeparator)))}
	b, err := os.ReadFile(path)
	if err != nil {
		r.Errors = []string{fmt.Sprintf("read: %v", err)}
		return r
	}
	var v map[string]any
	if err = json.Unmarshal(b, &v); err != nil {
		r.Errors = []string{fmt.Sprintf("json: %v", err)}
		return r
	}
	r.ID, _ = v["id"].(string)
	r.Schema, _ = v["schema"].(string)
	r.Status, _ = v["status"].(string)
	if u, ok := v["unknowns"].([]any); ok {
		for _, x := range u {
			if s, ok := x.(string); ok {
				r.Unknowns = append(r.Unknowns, s)
			}
		}
	}
	if len(r.Unknowns) == 0 {
		if u, ok := v["unknown"].([]any); ok {
			for _, x := range u {
				if s, ok := x.(string); ok {
					r.Unknowns = append(r.Unknowns, s)
				}
			}
		}
	}
	if f, ok := v["falsifiers"].([]any); ok {
		for _, x := range f {
			if s, ok := x.(string); ok {
				r.Falsifiers = append(r.Falsifiers, s)
			}
		}
	}
	_, hasSource := v["source"]
	_, hasRefs := v["refs"]
	r.Provenance = hasSource || hasRefs || v["source_capsule"] != nil || v["artifact"] != nil
	for _, k := range []string{"schema", "id", "status"} {
		if _, ok := v[k]; !ok {
			r.Errors = append(r.Errors, "missing "+k)
		}
	}
	r.BenchmarkReady = len(r.Errors) == 0 && r.Provenance && len(r.Falsifiers) > 0
	return r
}

func Scan(root string) (Report, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return Report{}, err
	}
	out := Report{Schema: "context-doctor-report-v1", Root: filepath.ToSlash(root)}
	err = filepath.Walk(root, func(path string, info os.FileInfo, e error) error {
		if e != nil {
			return e
		}
		if info.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}
		out.Files++
		r := check(path, root)
		if len(r.Errors) == 0 {
			out.Valid++
		} else {
			out.Invalid++
		}
		if r.BenchmarkReady {
			out.BenchmarkReady++
		}
		if len(r.Unknowns) > 0 {
			out.UnknownFields++
		}
		out.Records = append(out.Records, r)
		return nil
	})
	if err != nil {
		return out, err
	}
	sort.Slice(out.Records, func(i, j int) bool { return out.Records[i].Path < out.Records[j].Path })
	return out, nil
}
