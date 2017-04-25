package main

import (
	"errors"
	"sort"
	"strings"
)

// Flag defines a flag for a command.
// These will be parsed in Go and passed to the Run method in the Context struct.
type Flag struct {
	Name        string `json:"name"`
	Char        string `json:"char"`
	Description string `json:"description"`
	HasValue    bool   `json:"hasValue"`
	Hidden      bool   `json:"hidden"`
	Required    bool   `json:"required"`
}

func (f *Flag) String() string {
	s := " "
	switch {
	case f.Char != "" && f.Name != "":
		s = s + "-" + f.Char + ", --" + f.Name
	case f.Char != "":
		s = s + "-" + f.Char
	case f.Name != "":
		s = s + "--" + f.Name
	}
	if f.HasValue {
		s = s + " " + strings.ToUpper(f.Name)
	}
	return s
}

// ParseFlag parses a flag from argument inputs
func ParseFlag(input string, flags []*Flag) (*Flag, string, error) {
	keyvalue := strings.SplitN(input, "=", 2)
	key := keyvalue[0]
	value := ""
	if len(keyvalue) == 2 {
		value = keyvalue[1]
	}
	if len(key) > 2 && key[1] != '-' {
		return ParseFlag(key[:2]+"="+key[2:], flags)
	}
	for _, flag := range flags {
		if (flag.Char != "" && key == "-"+flag.Char) || key == "--"+flag.Name {
			if flag.HasValue {
				if value == "" {
					return nil, "", errors.New(flag.String() + " needs a value")
				}
				return flag, value, nil
			}
			if value != "" {
				return nil, "", errors.New(flag.String() + " does not take a value")
			}
			return flag, "", nil
		}
	}
	return nil, "", nil
}

// Flags are a list of flags
type Flags []Flag

// Len is for sorting
func (flags Flags) Len() int {
	return len(flags)
}

// Less is for sorting
func (flags Flags) Less(i, j int) bool {
	if flags[i].Char != "" && flags[j].Char == "" {
		return true
	}
	if flags[i].Char == "" && flags[j].Char != "" {
		return false
	}
	return flags[i].Name < flags[j].Name
}

// Swap the flags for sorting
func (flags Flags) Swap(i, j int) {
	flags[i], flags[j] = flags[j], flags[i]
}

// Sort sorts the flags
func (flags Flags) Sort() Flags {
	sort.Sort(flags)
	return flags
}
