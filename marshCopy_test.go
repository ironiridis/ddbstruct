package ddbstruct

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestStringRoundTrip(t *testing.T) {
	type z struct{ X string }
	in := &z{X: "example"}
	av := encString(in, 0)
	expectT(t, new(types.AttributeValueMemberS), av)
	out := &z{}
	decString(out, 0, av)
	if in.X != out.X {
		t.Fatalf("expected %q, got %q", in.X, out.X)
	}
}

func TestStringPtrNewRoundTrip(t *testing.T) {
	type z struct{ X *string }
	orig := "example"
	in := &z{X: &orig}
	av := encString(in, 0)
	expectT(t, new(types.AttributeValueMemberS), av)
	out := &z{}
	decString(out, 0, av)
	if out.X == nil {
		t.Fatal("out is nil")
	}
	if in.X == out.X {
		t.Fatal("out is not a new pointer")
	}
	if *in.X != *out.X {
		t.Fatalf("expected %q, got %q", *in.X, *out.X)
	}
}

func TestStringPtrReuseRoundTrip(t *testing.T) {
	type z struct{ X *string }
	orig := "example"
	in := &z{X: &orig}
	av := encString(in, 0)
	expectT(t, new(types.AttributeValueMemberS), av)
	out := &z{X: new(string)}
	decString(out, 0, av)
	if out.X == nil {
		t.Fatal("out is nil")
	}
	if in.X == out.X {
		t.Fatal("out is not a new pointer")
	}
	if *in.X != *out.X {
		t.Fatalf("expected %q, got %q", *in.X, *out.X)
	}
}

func TestBytesRoundTrip(t *testing.T) {
	type z struct{ X []byte }
	in := &z{X: []byte("example")}
	av := encBytes(in, 0)
	expectT(t, new(types.AttributeValueMemberB), av)
	out := &z{}
	decBytes(out, 0, av)
	compareSlice(t, in.X, out.X)
}

func TestBytesPtrNewRoundTrip(t *testing.T) {
	type z struct{ X *[]byte }
	orig := []byte("example")
	in := &z{X: &orig}
	av := encBytes(in, 0)
	expectT(t, new(types.AttributeValueMemberB), av)
	out := &z{}
	decBytes(out, 0, av)
	if out.X == nil {
		t.Fatal("out is nil")
	}
	if in.X == out.X {
		t.Fatal("out is not a new pointer")
	}
	compareSlice(t, *in.X, *out.X)
}

func TestBytesPtrReuseRoundTrip(t *testing.T) {
	type z struct{ X *[]byte }
	orig := []byte("example")
	in := &z{X: &orig}
	av := encBytes(in, 0)
	expectT(t, new(types.AttributeValueMemberB), av)
	out := &z{X: &[]byte{8, 7, 6, 5}}
	decBytes(out, 0, av)
	if out.X == nil {
		t.Fatal("out is nil")
	}
	if in.X == out.X {
		t.Fatal("out is not a new pointer")
	}
	compareSlice(t, *in.X, *out.X)
}

func TestBoolRoundTrip(t *testing.T) {
	type z struct{ X bool }
	in := &z{X: true}
	av := encBool(in, 0)
	expectT(t, new(types.AttributeValueMemberBOOL), av)
	out := &z{}
	decBool(out, 0, av)
	if in.X != out.X {
		t.Fatalf("expected %t, got %t", in.X, out.X)
	}
}

func TestBoolPtrNewRoundTrip(t *testing.T) {
	type z struct{ X *bool }
	orig := bool(true)
	in := &z{X: &orig}
	av := encBool(in, 0)
	expectT(t, new(types.AttributeValueMemberBOOL), av)
	out := &z{}
	decBool(out, 0, av)
	if out.X == nil {
		t.Fatal("out is nil")
	}
	if in.X == out.X {
		t.Fatal("out is not a new pointer")
	}
	if *in.X != *out.X {
		t.Fatalf("expected %t, got %t", *in.X, *out.X)
	}
}

func TestBoolPtrReuseRoundTrip(t *testing.T) {
	type z struct{ X *bool }
	orig := bool(true)
	in := &z{X: &orig}
	av := encBool(in, 0)
	expectT(t, new(types.AttributeValueMemberBOOL), av)
	out := &z{X: new(bool)}
	decBool(out, 0, av)
	if out.X == nil {
		t.Fatal("out is nil")
	}
	if in.X == out.X {
		t.Fatal("out is not a new pointer")
	}
	if *in.X != *out.X {
		t.Fatalf("expected %t, got %t", *in.X, *out.X)
	}
}
