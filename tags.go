package ddbstruct

import (
	"fmt"
	"reflect"
	"regexp"
	"sync"
)

type structMetadataCache struct {
	sync.Mutex
	types map[reflect.Type]structMetadata
}

var cache = structMetadataCache{types: map[reflect.Type]structMetadata{}}

type structMetadata struct {
	f     []fieldTagData
	pkidx int
	skidx int
}

func (c *structMetadataCache) get(d interface{}) structMetadata {
	dt := reflect.TypeOf(d)
	if dtk := dt.Kind(); dtk != reflect.Pointer {
		panic(fmt.Errorf("expected pointer to struct, got %T, a %s", d, dtk))
	}
	dte := dt.Elem()
	if dtek := dte.Kind(); dtek != reflect.Struct {
		panic(fmt.Errorf("expected pointer to struct, got %T, a pointer to %s", d, dtek))
	}

	c.Lock()
	defer c.Unlock()

	if r, ok := c.types[dte]; ok {
		return r
	}

	ret := structMetadata{pkidx: 0, skidx: -1}

	for n := 0; n < dte.NumField(); n++ {
		stv, err := parseFieldTag(dte, n)
		if err != nil {
			panic(fmt.Errorf("cannot parse tags on field %d of struct %s: %w", n, dte, err))
		}
		err = stv.typecalc()
		if err != nil {
			panic(fmt.Errorf("unable to typecalc field %d of struct %s: %w", n, dte, err))
		}
		ret.f = append(ret.f, *stv)
	}

	c.types[dt.Elem()] = ret
	return ret
}

var fieldTagMatcher = regexp.MustCompile(
	"^(" + // anchor to beginning of tag string
		"(?P<pk>pk)" + // partition key flag
		"|(?P<sk>sk)" + // sort key flag
		"|(n=(?P<n>.+))" + // field name override
		"|(t=(?P<t>.+))" + // field type override
		"|(def=(?P<def>.+))" + // default value
		"|(?P<opt>opt)" + // optional flag
		")(,|$)") // match the end of the string or a comma

func parseFieldTag(t reflect.Type, idx int) (*fieldTagData, error) {
	var err error
	ret := &fieldTagData{fieldname: t.Field(idx).Name, gotype: t.Field(idx).Type, structidx: idx}
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
				ret.fieldname = v
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
