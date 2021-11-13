package ddbstruct

import (
	"fmt"
	"reflect"
	"regexp"
)

var fieldTagMatcher = regexp.MustCompile(
	"^(" + // anchor to beginning of tag string
		"(?P<pk>pk)" + // partition key flag
		"|(?P<sk>sk)" + // sort key flag
		"|(n=(?P<n>.+))" + // field name override
		"|(t=(?P<t>.+))" + // field type override
		"|(def=(?P<def>.+))" + // default value
		"|(?P<opt>opt)" + // optional flag
		")(,|$)") // match the end of the string or a comma

func parseFieldTag(t reflect.Type, idx int) (*field, error) {
	var err error
	ret := &field{name: t.Field(idx).Name, gotype: t.Field(idx).Type, idx: idx}
	tagdata := t.Field(idx).Tag.Get("ddb")
	for len(tagdata) > 0 {
		subexp := fieldTagMatcher.FindStringSubmatch(tagdata)
		if subexp == nil || len(subexp[0]) == 0 {
			return nil, fmt.Errorf("could not parse struct tag %q", tagdata)
		}
		for _, sen := range fieldTagMatcher.SubexpNames() {
			if sen == "" { // unnamed subexpression
				continue
			}
			v := subexp[fieldTagMatcher.SubexpIndex(sen)]
			if v == "" { // this subexpression didn't match anything
				continue
			}
			switch sen {
			case "pk":
				ret.pk = true
			case "sk":
				ret.sk = true
			case "opt":
				ret.optional = true
			case "n":
				ret.name = v
			case "t":
				ret.enctype = v
			case "def":
				ret.defvalue = v
			}
		}
		tagdata = tagdata[len(subexp[0]):]
	}
	return ret, err
}
