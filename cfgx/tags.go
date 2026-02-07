package cfgx

import (
	"reflect"
	"strings"
)

const (
	tagKeyCfg      = "cfg"
	tagKeyDefault  = "default"
	tagKeyRequired = "required"
)

type fieldTags struct {
	Var      string
	Default  string
	HasDef   bool
	Required bool
}

// parseFieldTags reads cfg/default/required tags from a struct field.
// If cfg tag is empty, the field is considered unmanaged by cfgx.
func parseFieldTags(sf reflect.StructField) fieldTags {
	t := sf.Tag

	varName := strings.TrimSpace(t.Get(tagKeyCfg))
	if varName == "" {
		return fieldTags{}
	}

	def, ok := t.Lookup(tagKeyDefault)
	def = strings.TrimSpace(def)

	req := strings.TrimSpace(t.Get(tagKeyRequired))
	required := req == "true" || req == "1" || strings.EqualFold(req, "yes")

	return fieldTags{
		Var:      varName,
		Default:  def,
		HasDef:   ok,
		Required: required,
	}
}
