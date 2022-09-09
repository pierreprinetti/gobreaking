package rules

import (
	"fmt"

	"golang.org/x/tools/go/packages"
)

func PreserveNameTypes(basePackage, newPackage *packages.Package) (result Result) {
	if basePackage == nil {
		return result
	}

	for _, name := range basePackage.Types.Scope().Names() {
		if newObject := newPackage.Types.Scope().Lookup(name); newObject != nil && newObject.Exported() {
			baseObject := basePackage.Types.Scope().Lookup(name)
			if newObject.Type().String() != baseObject.Type().String() {
				result.AddMajor(fmt.Sprintf("Package %q: name %q has changed signature from:\n\t%q\nto:\n\t%q\n", newPackage.Name, baseObject.Name(), baseObject.Type(), newObject.Type()))
			}
		}
	}
	return result
}
