package flags

import (
	"fmt"

	"github.com/spf13/pflag"
)

type FlagValueSet interface {
	VisitAll(fn func(flagValue))
}

type ResourceFlagSet struct {
	flags []ryukFlag
}

func NewResourceFlagSet(flags ...ryukFlag) ResourceFlagSet {
	return ResourceFlagSet{flags: flags}
}

func (f ResourceFlagSet) VisitAll(fn func(flagValue)) {
	for _, flag := range f.flags {
		fn(flag)
	}
}

func NewCreateFlagSet() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("createFlagSet", pflag.ContinueOnError)

	flagSet.String("name", "", "Name of the resource")
	// TODO: Remove this later
	flagSet.Int("int", 1, "Placeholder test")

	return flagSet
}

func NewDeleteFlagSet(resource string) *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("deleteFlagSet", pflag.ContinueOnError)

	description := fmt.Sprintf("Name of the %s", resource)
	flagSet.String("name", "", description)
	// TODO: Remove this later
	flagSet.Int("int", 1, "Placeholder test")

	return flagSet
}
