package cfg

import (
	"github.com/simplesurance/baur/v2/internal/deepcopy"
)

// InputInclude is a reusable Input definition.
type InputInclude struct {
	IncludeID string `toml:"include_id" comment:"identifier of the include"`

	Files         []FileInputs
	GolangSources []GolangSources `comment:"Inputs specified by resolving dependencies of Golang source files or packages."`

	filepath string
}

func (in *InputInclude) fileInputs() []FileInputs {
	return in.Files
}

func (in *InputInclude) golangSourcesInputs() []GolangSources {
	return in.GolangSources
}

func (in *InputInclude) IsEmpty() bool {
	return len(in.Files) == 0 && len(in.GolangSources) == 0
}

// validate checks if the stored information is valid.
func (in *InputInclude) validate() error {
	if err := validateIncludeID(in.IncludeID); err != nil {
		if in.IncludeID != "" {
			err = fieldErrorWrap(err, in.IncludeID)
		}
		return err
	}

	if in.IsEmpty() {
		return nil
	}

	if err := inputValidate(in); err != nil {
		return err
	}

	return nil
}

func (in *InputInclude) clone() *InputInclude {
	var clone InputInclude

	deepcopy.MustCopy(in, &clone)
	// filepath is assigned manually because filepath is a private field, MustCopy() only clones exported fields
	clone.filepath = in.filepath

	return &clone
}
