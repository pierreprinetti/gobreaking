package rules

import (
	"fmt"
	"strings"
)

type Result struct {
	Major []string
}

func (r *Result) Add(r1 Result) {
	r.Major = append(r.Major, r1.Major...)
}

func (r *Result) AddMajor(msg ...string) {
	r.Major = append(r.Major, msg...)
}

const noChangesMsg = `No breaking changes detected.
This tool currently only checks 2 out of 11 rules. Please manually check:
	* Adding or removing a method in an exported interface.
	* Adding or removing a parameter in an exported function or interface.
	* Changing the type of a parameter in an exported function or interface.
	* Adding or removing a result in an exported function or interface.
	* Changing the type of a result in an exported function or interface.
	* Removing an exported field from an exported struct.
	* Changing the type of an exported field of an exported struct.
	* Adding an exported or unexported field to an exported struct containing only exported fields.
	* Repositioning a field in an exported struct containing only exported fields.`

func (r Result) String() string {
	switch {
	case len(r.Major) > 0:
		return fmt.Sprintf("Potential breaking changes detected:\n%s", strings.Join(r.Major, "\n"))
	default:
		return fmt.Sprintf(noChangesMsg)
	}
}

func (r Result) MarkdownString() string {
	switch {
	case len(r.Major) > 0:
		return fmt.Sprintf("Potential breaking changes detected:\n```\n%s\n```", strings.Join(r.Major, "\n"))
	default:
		return fmt.Sprintf(noChangesMsg)
	}
}
