package ddbstruct

import (
	"fmt"
	"reflect"
	"sync"
)

type structMetadata struct {
	f  []field
	pk *field
	sk *field
}

type structMetadataCache struct {
	sync.Mutex
	types map[reflect.Type]structMetadata
}

var cache = structMetadataCache{types: map[reflect.Type]structMetadata{}}

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

	ret := structMetadata{}
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
		if stv.pk {
			if ret.pk != nil {
				panic(fmt.Errorf("field %q tagged as pk, but pk is already tagged on field %q", stv.name, ret.pk.name))
			}
			ret.pk = &ret.f[n]
		}
		if stv.sk {
			if ret.sk != nil {
				panic(fmt.Errorf("field %q tagged as sk, but sk is already tagged on field %q", stv.name, ret.sk.name))
			}
			if stv.pk {
				panic(fmt.Errorf("field %q tagged as sk, but also tagged as pk", stv.name))
			}
			ret.sk = &ret.f[n]
		}
	}

	c.types[dt.Elem()] = ret
	return ret
}
