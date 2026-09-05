package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/yagizkaterli/lobi-agentic-revolution/internal/doctor"
	"os"
)

func main() {
	root := flag.String("root", ".", "capsule root")
	out := flag.String("out", "", "report path")
	flag.Parse()
	r, e := doctor.Scan(*root)
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(2)
	}
	b, _ := json.MarshalIndent(r, "", "  ")
	b = append(b, '\n')
	if *out != "" {
		if e = os.WriteFile(*out, b, 0644); e != nil {
			fmt.Fprintln(os.Stderr, e)
			os.Exit(2)
		}
	}
	fmt.Print(string(b))
	if r.Invalid > 0 {
		os.Exit(1)
	}
}
