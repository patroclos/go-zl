package prompt

import (
	"bufio"
	"bytes"
	"strings"
)

// The prompt package contains everything related to Q.A. prompts and omissables
// the existance of prompts puts pressure on a view-over-edit mentality to experience
// properly embedded prompts.

type EmbeddedPrompt struct {
	Q     string
	A     string
	Omits map[string]string
}

var example = EmbeddedPrompt{
	Q: "What does the {{.animal}} say?",
	A: "meow",
	Omits: map[string]string{"animal": "fox"},
}

func ExtractAll(txt string) []EmbeddedPrompt {
	scn := bufio.NewScanner(bytes.NewReader([]byte(txt)))

	results := make([]EmbeddedPrompt, 0, 4)

	var seenQ bool
	var q string
	for scn.Scan() {
		line := scn.Text()
		if seenQ {
			if !strings.HasPrefix(line, "A. ") {
				seenQ = false
				continue
			}

			results = append(results, EmbeddedPrompt{
				Q: q,
				A: line[3:],
			})
			continue
		}

		if strings.HasPrefix(line, "Q. ") {
			seenQ = true
			q = line[3:]
			continue
		}
	}
	return results
}
