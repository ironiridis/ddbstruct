package ddbstruct

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestIntCantParse(t *testing.T) {
	type z struct{ X int }
	av := &types.AttributeValueMemberN{Value: "pickle"}
	out := &z{}
	expectPanic(t, func() { decInt(out, 0, av) })
	if t.Failed() {
		t.Logf("got %v", out.X)
	}
}

func TestIntTooBig(t *testing.T) {
	type z struct{ X int16 }
	av := &types.AttributeValueMemberN{Value: "65536"}
	out := &z{}
	expectPanic(t, func() { decInt(out, 0, av) })
	if t.Failed() {
		t.Logf("got %v", out.X)
	}
}

func TestIntRoundTrip(t *testing.T) {
	type z struct{ X int }
	in := &z{X: 25}
	av := encInt(in, 0)
	expectT(t, new(types.AttributeValueMemberN), av)
	out := &z{}
	decInt(out, 0, av)
	if in.X != out.X {
		t.Fatalf("expected %d, got %d", in.X, out.X)
	}
}

func TestIntPtrNewRoundTrip(t *testing.T) {
	type z struct{ X *int }
	orig := int(38)
	in := &z{X: &orig}
	av := encInt(in, 0)
	expectT(t, new(types.AttributeValueMemberN), av)
	out := &z{}
	decInt(out, 0, av)
	if out.X == nil {
		t.Fatal("out is nil")
	}
	if in.X == out.X {
		t.Fatal("out is not a new pointer")
	}
	if *in.X != *out.X {
		t.Fatalf("expected %d, got %d", *in.X, *out.X)
	}
}

func TestIntPtrReuseRoundTrip(t *testing.T) {
	type z struct{ X *int }
	orig := int(450)
	in := &z{X: &orig}
	av := encInt(in, 0)
	expectT(t, new(types.AttributeValueMemberN), av)
	out := &z{X: new(int)}
	decInt(out, 0, av)
	if out.X == nil {
		t.Fatal("out is nil")
	}
	if in.X == out.X {
		t.Fatal("out is not a new pointer")
	}
	if *in.X != *out.X {
		t.Fatalf("expected %d, got %d", *in.X, *out.X)
	}
}

func TestUintCantParse(t *testing.T) {
	type z struct{ X uint }
	av := &types.AttributeValueMemberN{Value: "pickle"}
	out := &z{}
	expectPanic(t, func() { decUint(out, 0, av) })
	if t.Failed() {
		t.Logf("got %v", out.X)
	}
}

func TestUintTooBig(t *testing.T) {
	type z struct{ X uint16 }
	av := &types.AttributeValueMemberN{Value: "65536"}
	out := &z{}
	expectPanic(t, func() { decUint(out, 0, av) })
	if t.Failed() {
		t.Logf("got %v", out.X)
	}
}

func TestUintNegative(t *testing.T) {
	type z struct{ X uint16 }
	av := &types.AttributeValueMemberN{Value: "-1"}
	out := &z{}
	expectPanic(t, func() { decUint(out, 0, av) })
	if t.Failed() {
		t.Logf("got %v", out.X)
	}
}

func TestUintRoundTrip(t *testing.T) {
	type z struct{ X uint }
	in := &z{X: 25}
	av := encUint(in, 0)
	expectT(t, new(types.AttributeValueMemberN), av)
	out := &z{}
	decUint(out, 0, av)
	if in.X != out.X {
		t.Fatalf("expected %d, got %d", in.X, out.X)
	}
}

func TestUintPtrNewRoundTrip(t *testing.T) {
	type z struct{ X *uint }
	orig := uint(38)
	in := &z{X: &orig}
	av := encUint(in, 0)
	expectT(t, new(types.AttributeValueMemberN), av)
	out := &z{}
	decUint(out, 0, av)
	if out.X == nil {
		t.Fatal("out is nil")
	}
	if in.X == out.X {
		t.Fatal("out is not a new pouinter")
	}
	if *in.X != *out.X {
		t.Fatalf("expected %d, got %d", *in.X, *out.X)
	}
}

func TestUintPtrReuseRoundTrip(t *testing.T) {
	type z struct{ X *uint }
	orig := uint(450)
	in := &z{X: &orig}
	av := encUint(in, 0)
	expectT(t, new(types.AttributeValueMemberN), av)
	out := &z{X: new(uint)}
	decUint(out, 0, av)
	if out.X == nil {
		t.Fatal("out is nil")
	}
	if in.X == out.X {
		t.Fatal("out is not a new pouinter")
	}
	if *in.X != *out.X {
		t.Fatalf("expected %d, got %d", *in.X, *out.X)
	}
}

func TestFloat32CantParse(t *testing.T) {
	type z struct{ X float32 }
	av := &types.AttributeValueMemberN{Value: "pickle"}
	out := &z{}
	expectPanic(t, func() { decFloat(32)(out, 0, av) })
	if t.Failed() {
		t.Logf("got %v", out.X)
	}
}

func TestFloat32TooBig(t *testing.T) {
	type z struct{ X float32 }
	av := &types.AttributeValueMemberN{Value: "1e39"}
	out := &z{}
	expectPanic(t, func() { decFloat(32)(out, 0, av) })
	if t.Failed() {
		t.Logf("got %v", out.X)
	}
}

func TestFloat32RoundTrip(t *testing.T) {
	type z struct{ X float32 }
	in := &z{X: 123.45}
	av := encFloat(32)(in, 0)
	expectT(t, new(types.AttributeValueMemberN), av)
	out := &z{}
	decFloat(32)(out, 0, av)
	if in.X != out.X {
		t.Fatalf("expected %g, got %g", in.X, out.X)
	}
}

func TestFloat32PtrNewRoundTrip(t *testing.T) {
	type z struct{ X *float32 }
	orig := float32(98.765)
	in := &z{X: &orig}
	av := encFloat(32)(in, 0)
	expectT(t, new(types.AttributeValueMemberN), av)
	out := &z{}
	decFloat(32)(out, 0, av)
	if out.X == nil {
		t.Fatal("out is nil")
	}
	if in.X == out.X {
		t.Fatal("out is not a new pointer")
	}
	if *in.X != *out.X {
		t.Fatalf("expected %g, got %g", *in.X, *out.X)
	}
}

func TestFloat32PtrReuseRoundTrip(t *testing.T) {
	type z struct{ X *float32 }
	orig := float32(87.654)
	in := &z{X: &orig}
	av := encFloat(32)(in, 0)
	expectT(t, new(types.AttributeValueMemberN), av)
	out := &z{X: new(float32)}
	decFloat(32)(out, 0, av)
	if out.X == nil {
		t.Fatal("out is nil")
	}
	if in.X == out.X {
		t.Fatal("out is not a new pofloat32er")
	}
	if *in.X != *out.X {
		t.Fatalf("expected %g, got %g", *in.X, *out.X)
	}
}

func TestFloat64CantParse(t *testing.T) {
	type z struct{ X float64 }
	av := &types.AttributeValueMemberN{Value: "pickle"}
	out := &z{}
	expectPanic(t, func() { decFloat(64)(out, 0, av) })
	if t.Failed() {
		t.Logf("got %v", out.X)
	}
}

func TestFloat64TooBig(t *testing.T) {
	type z struct{ X float64 }
	av := &types.AttributeValueMemberN{Value: "1e309"}
	out := &z{}
	expectPanic(t, func() { decFloat(64)(out, 0, av) })
	if t.Failed() {
		t.Logf("got %v", out.X)
	}
}

func TestFloat64RoundTrip(t *testing.T) {
	type z struct{ X float64 }
	in := &z{X: 123.45}
	av := encFloat(64)(in, 0)
	expectT(t, new(types.AttributeValueMemberN), av)
	out := &z{}
	decFloat(64)(out, 0, av)
	if in.X != out.X {
		t.Fatalf("expected %g, got %g", in.X, out.X)
	}
}

func TestFloat64PtrNewRoundTrip(t *testing.T) {
	type z struct{ X *float64 }
	orig := float64(98.765)
	in := &z{X: &orig}
	av := encFloat(64)(in, 0)
	expectT(t, new(types.AttributeValueMemberN), av)
	out := &z{}
	decFloat(64)(out, 0, av)
	if out.X == nil {
		t.Fatal("out is nil")
	}
	if in.X == out.X {
		t.Fatal("out is not a new pointer")
	}
	if *in.X != *out.X {
		t.Fatalf("expected %g, got %g", *in.X, *out.X)
	}
}

func TestFloat64PtrReuseRoundTrip(t *testing.T) {
	type z struct{ X *float64 }
	orig := float64(87.654)
	in := &z{X: &orig}
	av := encFloat(64)(in, 0)
	expectT(t, new(types.AttributeValueMemberN), av)
	out := &z{X: new(float64)}
	decFloat(64)(out, 0, av)
	if out.X == nil {
		t.Fatal("out is nil")
	}
	if in.X == out.X {
		t.Fatal("out is not a new pofloat64er")
	}
	if *in.X != *out.X {
		t.Fatalf("expected %g, got %g", *in.X, *out.X)
	}
}
