package rules

import (
	"fmt"
	"strings"

	"golang.org/x/tools/go/packages"
)

type ErrPreserveNames struct {
	MissingNames []string
}

func (err ErrPreserveNames) Error() string {
	return fmt.Sprintf("missing names: %s", strings.Join(err.MissingNames, ", "))
}

func PreserveNames(basePackage, newPackage *packages.Package) (result Result) {
	if basePackage == nil {
		return result
	}

	newPackageNames := make(map[string]struct{})
	for _, name := range newPackage.Types.Scope().Names() {
		newPackageNames[name] = struct{}{}
	}

	for _, name := range basePackage.Types.Scope().Names() {
		if obj := basePackage.Types.Scope().Lookup(name); obj.Exported() {
			if _, ok := newPackageNames[name]; !ok {
				result.AddMajor("package %q: missing name: %q", newPackage.Name, name)
			}
		}
	}
	return result
}
