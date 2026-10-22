// A standalone checker for Magpie routing-group pattern syntax.
// Usage: go run docs/local/magpie-regex-check.go 're:(?i:.+/glm-5[.-]3-flash)' < model-ids.txt
//
// It reproduces compilePattern in internal/provider/groupmatch.go and only
// validates string matching; it does not query Magpie provider state.
package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func compilePattern(pattern string) (*regexp.Regexp, error) {
	if raw, ok := strings.CutPrefix(pattern, "re:"); ok {
		if strings.TrimSpace(raw) == "" {
			return nil, fmt.Errorf("empty regex after re:")
		}
		return regexp.Compile("^(?:" + raw + ")$")
	}
	parts := strings.Split(pattern, "*")
	for i, p := range parts {
		parts[i] = regexp.QuoteMeta(p)
	}
	return regexp.Compile("(?i)^" + strings.Join(parts, ".*") + "$")
}

func main() {
	if len(os.Args) != 2 || strings.TrimSpace(os.Args[1]) == "" {
		fmt.Fprintln(os.Stderr, "Usage: go run docs/local/magpie-regex-check.go 're:<pattern>' < model-ids.txt")
		os.Exit(2)
	}
	re, err := compilePattern(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid Magpie pattern: %v\n", err)
		os.Exit(2)
	}
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	count := 0
	for scanner.Scan() {
		id := strings.TrimSpace(scanner.Text())
		if id == "" || strings.HasPrefix(id, "#") {
			continue
		}
		if re.MatchString(id) {
			fmt.Printf("MATCH\t%s\n", id)
			count++
		} else {
			fmt.Printf("----\t%s\n", id)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "read model IDs: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "matched %d ID(s)\n", count)
}
