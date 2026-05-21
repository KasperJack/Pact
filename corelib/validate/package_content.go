package validate

import (
    "fmt"
    "Pact/corelib/model"
)

func PackageContent(entity model.PackageContent) error {

	ok := isValidVersioningSchema(entity.Package.Versioning)

	if !ok {
		return fmt.Errorf("invalid versioning schema %q", entity.Package.Versioning)
	}

	for _,r := range entity.Releases {

		ok := isValidVersion(entity.Package.Versioning,r.Version)
		if !ok {
        	return fmt.Errorf("invalid version %q for versioning schema %q", r.Version, entity.Package.Versioning)
    }

	}

	return nil
}