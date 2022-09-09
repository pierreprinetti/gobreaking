package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pierreprinetti/gobreaking/rules"
	"golang.org/x/tools/go/packages"
)

func packageErrors(pkgs []*packages.Package) error {
	var (
		found bool
		msg   strings.Builder
	)
	packages.Visit(pkgs, nil, func(pkg *packages.Package) {
		for _, err := range pkg.Errors {
			found = true
			msg.WriteString(err.Error())
		}
	})
	if found {
		return fmt.Errorf(msg.String())
	}
	return nil
}

func loadPackages(ctx context.Context, dir string) ([]*packages.Package, error) {
	pkgs, err := packages.Load(&packages.Config{
		Context: ctx,
		Dir:     dir,
		Mode:    packages.NeedName | packages.NeedTypes,
		Tests:   false,
	}, filepath.Join(dir, "..."))

	if err != nil {
		return nil, err
	}

	if err := packageErrors(pkgs); err != nil {
		return nil, err
	}

	return pkgs, nil
}

func main() {
	ctx := context.Background()

	flag.Parse()
	args := flag.Args()

	if len(args) != 2 {
		log.Fatal("Expected two arguments: the base tree and the new tree.")
	}

	basePath, err := filepath.Abs(args[0])
	if err != nil {
		log.Fatalf("Error with path %q", args[0])
	}
	newPath, err := filepath.Abs(args[1])
	if err != nil {
		log.Fatalf("Error with path %q", args[1])
	}

	basePackages, err := loadPackages(ctx, basePath)
	if err != nil {
		log.Fatalf(err.Error())
	}

	newPackages, err := loadPackages(ctx, newPath)
	if err != nil {
		log.Fatalf(err.Error())
	}

	var result rules.Result
	result.Add(rules.PreservePackages(basePackages, newPackages))

	basePackagesMap := make(map[string]*packages.Package)
	for _, pkg := range basePackages {
		basePackagesMap[pkg.ID] = pkg
	}
	for _, newPackage := range newPackages {
		for _, test := range [...]func(basePackage, newPackage *packages.Package) rules.Result{
			rules.PreserveNames,     // Removing an exported name (constant, type, variable, function).
			rules.PreserveNameTypes, // Changing the type of an exported name.
			// TODO:
			// rules.PreserveInterfaceMethods,          // Adding or removing a method in an exported interface.
			// rules.PreserveParameters,                // Adding or removing a parameter in an exported function or interface.
			// rules.PreserveParameterTypes,            // Changing the type of a parameter in an exported function or interface.
			// rules.PreserveReturns,                   // Adding or removing a result in an exported function or interface.
			// rules.PreserveReturnTypes,               // Changing the type of a result in an exported function or interface.
			// rules.PreserveStructFields,              // Removing an exported field from an exported struct.
			// rules.PreserveStructFieldTypes,          // Changing the type of an exported field of an exported struct.
			// rules.PreserveStructImplicitAssignment,  // Adding an exported or unexported field to an exported struct containing only exported fields.
			// rules.PreserveStructExportedFieldsOrder, // Repositioning a field in an exported struct containing only exported fields.
		} {
			result.Add(test(basePackagesMap[newPackage.ID], newPackage))
		}
	}

	if os.Getenv("GITHUB_ACTIONS") == "true" {
		f, err := os.Create(os.Getenv("GITHUB_STEP_SUMMARY"))
		if err != nil {
			log.Print(err)
		}
		defer f.Close()
		msg := result.MarkdownString()
		if _, err := fmt.Fprint(f, msg); err != nil {
			log.Print(err)
		}
		msg = strings.Replace(msg, "%", "%25", -1)
		msg = strings.Replace(msg, "\n", "%0A", -1)
		msg = strings.Replace(msg, "\r", "%0D", -1)
		fmt.Printf("::set-output name=verdict::%s", msg)
	} else {
		fmt.Print(result)
	}
}
