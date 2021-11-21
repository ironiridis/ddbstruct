package ddbstruct

import (
	"math/big"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestTimeJSONRoundTrip(t *testing.T) {
	type z struct{ X time.Time }
	in := &z{X: time.Now()}
	av := encJSON(in, 0)
	expectT(t, new(types.AttributeValueMemberS), av)
	out := &z{}
	decJSON(out, 0, av)
	if in.X.UnixNano() != out.X.UnixNano() {
		t.Fatalf("expected %v, got %v", in.X, out.X)
	}
}

func TestTimePtrJSONRoundTrip(t *testing.T) {
	type z struct{ X *time.Time }
	orig := time.Now()
	in := &z{X: &orig}
	av := encJSON(in, 0)
	expectT(t, new(types.AttributeValueMemberS), av)
	out := &z{}
	decJSON(out, 0, av)
	if out.X == nil {
		t.Fatal("out is nil")
	}
	if in.X == out.X {
		t.Fatal("out is not a new pointer")
	}
	if in.X.UnixNano() != out.X.UnixNano() {
		t.Fatalf("expected %v, got %v", *in.X, *out.X)
	}
}

func TestBigNumPtrJSONRoundTrip(t *testing.T) {
	type z struct{ X *big.Int }
	in := &z{X: big.NewInt(987654321)}
	av := encJSON(in, 0)
	expectT(t, new(types.AttributeValueMemberS), av)
	out := &z{}
	decJSON(out, 0, av)
	if out.X == nil {
		t.Fatal("out is nil")
	}
	if in.X == out.X {
		t.Fatal("out is not a new pointer")
	}
	if in.X.Uint64() != out.X.Uint64() {
		t.Fatalf("expected %s, got %s", in.X, out.X)
	}
}

func TestStructJSONRawRoundTrip(t *testing.T) {
	type s struct {
		S string
		I int
		F float32
		B bool
	}
	type z struct{ X s }
	in := &z{s{S: "string", I: 13, F: 89.3, B: true}}
	av := encJSONRaw(in, 0)
	expectT(t, new(types.AttributeValueMemberS), av)
	out := &z{}
	decJSONRaw(out, 0, av)
	if in.X != out.X {
		t.Fatalf("expected %+v, got %+v", in.X, out.X)
	}
}

func TestStructPtrJSONRawRoundTrip(t *testing.T) {
	type s struct {
		S string
		I int
		F float32
		B bool
	}
	type z struct{ X *s }
	orig := s{S: "string", I: 13, F: 89.3, B: true}
	in := &z{X: &orig}
	av := encJSONRaw(in, 0)
	expectT(t, new(types.AttributeValueMemberS), av)
	out := &z{}
	decJSONRaw(out, 0, av)
	if out.X == nil {
		t.Fatal("out is nil")
	}
	if in.X == out.X {
		t.Fatal("out is not a new pointer")
	}
	if *in.X != *out.X {
		t.Fatalf("expected %+v, got %+v", *in.X, *out.X)
	}
}
