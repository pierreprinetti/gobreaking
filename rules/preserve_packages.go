package rules

import (
	"fmt"

	"golang.org/x/tools/go/packages"
)

func PreservePackages(basePackages, newPackages []*packages.Package) (result Result) {
	newPackagesMap := make(map[string]*packages.Package)
	for _, pkg := range newPackages {
		newPackagesMap[pkg.ID] = pkg
	}

	for _, pkg := range basePackages {
		if _, ok := newPackagesMap[pkg.ID]; !ok {
			result.AddMajor(fmt.Sprintf("Package missing: %q", pkg.Name))
		}
	}
	return result
}
