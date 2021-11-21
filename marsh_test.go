package ddbstruct

// This file just provide some testing helpers.

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func expectT(t *testing.T, e, v interface{}) {
	t.Helper()
	if reflect.TypeOf(v) != reflect.TypeOf(e) {
		t.Errorf("expected %T, got %T", e, v)
	}
	switch attr := v.(type) {
	case *types.AttributeValueMemberS:
		t.Logf("string attribute: %q", attr.Value)
	case *types.AttributeValueMemberB:
		t.Logf("binary attribute: %q", attr.Value)
	case *types.AttributeValueMemberBOOL:
		t.Logf("boolean attribute: %t", attr.Value)
	case *types.AttributeValueMemberN:
		t.Logf("numeric attribute: %s", attr.Value)
	default:
		t.Logf("attribute of type %T", v)
	}
}

func compareSlice(t *testing.T, a, b interface{}) {
	t.Helper()
	if a == nil && b == nil {
		// if both are nil, they are, technically, the same
		return
	}
	if reflect.TypeOf(a).Kind() != reflect.Slice {
		t.Fatalf("cannot slice compare because a (%T) is not a slice", a)
		return
	}
	if reflect.TypeOf(b).Kind() != reflect.Slice {
		t.Fatalf("cannot slice compare because b (%T) is not a slice", b)
		return
	}
	if reflect.TypeOf(a).Elem() != reflect.TypeOf(b).Elem() {
		t.Fatalf("cannot slice compare because a (%T) is not a slice of the same type as b (%T) ", a, b)
		return
	}
	if a == nil {
		t.Error("slices differ: a is nil")
		return
	}
	if b == nil {
		t.Error("slices differ: b is nil")
		return
	}
	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)
	if va.Len() != vb.Len() {
		t.Errorf("slices differ: length of a (%d) does not match length of b (%d)", va.Len(), vb.Len())
		return
	}
	for j := 0; j < va.Len(); j++ {
		// weird go-ism: you can compare interface values, and the runtime will only panic sometimes
		if va.Index(j).Interface() != vb.Index(j).Interface() {
			t.Errorf("slices differ at index %d", j)
			return
		}
	}
}

func expectPanic(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		t.Helper()
		if pv := recover(); pv != nil {
			t.Logf("did panic: %v", pv)
			return
		}
	}()
	f()
	t.Error("expected panic")
}
