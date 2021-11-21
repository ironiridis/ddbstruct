package ddbstruct

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestTimeEpochRoundTrip(t *testing.T) {
	type z struct{ X time.Time }
	in := &z{X: time.Now()}
	av := encTimeEpoch(in, 0)
	expectT(t, new(types.AttributeValueMemberN), av)
	out := &z{}
	decTimeEpoch(out, 0, av)
	if in.X.Unix() != out.X.Unix() {
		t.Fatalf("expected %v, got %v", in.X, out.X)
	}
}

func TestTimePtrEpochRoundTrip(t *testing.T) {
	type z struct{ X *time.Time }
	orig := time.Now()
	in := &z{X: &orig}
	av := encTimeEpoch(in, 0)
	expectT(t, new(types.AttributeValueMemberN), av)
	out := &z{}
	decTimeEpoch(out, 0, av)
	if out.X == nil {
		t.Fatal("out is nil")
	}
	if in.X == out.X {
		t.Fatal("out is not a new pointer")
	}
	if in.X.Unix() != out.X.Unix() {
		t.Fatalf("expected %v, got %v", *in.X, *out.X)
	}
}

func TestTimeNanoRoundTrip(t *testing.T) {
	type z struct{ X time.Time }
	in := &z{X: time.Now()}
	av := encTimeNano(in, 0)
	expectT(t, new(types.AttributeValueMemberN), av)
	out := &z{}
	decTimeNano(out, 0, av)
	if in.X.UnixNano() != out.X.UnixNano() {
		t.Fatalf("expected %v, got %v", in.X, out.X)
	}
}

func TestTimePtrNanoRoundTrip(t *testing.T) {
	type z struct{ X *time.Time }
	orig := time.Now()
	in := &z{X: &orig}
	av := encTimeNano(in, 0)
	expectT(t, new(types.AttributeValueMemberN), av)
	out := &z{}
	decTimeNano(out, 0, av)
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

func TestDurationStringRoundTrip(t *testing.T) {
	type z struct{ X time.Duration }
	in := &z{X: time.Millisecond*250 + time.Hour*8}
	av := encDurationString(in, 0)
	expectT(t, new(types.AttributeValueMemberS), av)
	out := &z{}
	decDurationString(out, 0, av)
	if in.X != out.X {
		t.Fatalf("expected %v, got %v", in.X, out.X)
	}
}

func TestDurationPtrStringRoundTrip(t *testing.T) {
	type z struct{ X *time.Duration }
	orig := time.Millisecond*350 + time.Hour*12
	in := &z{X: &orig}
	av := encDurationString(in, 0)
	expectT(t, new(types.AttributeValueMemberS), av)
	out := &z{}
	decDurationString(out, 0, av)
	if out.X == nil {
		t.Fatal("out is nil")
	}
	if in.X == out.X {
		t.Fatal("out is not a new pointer")
	}
	if *in.X != *out.X {
		t.Fatalf("expected %v, got %v", *in.X, *out.X)
	}
}

func TestDurationNanoRoundTrip(t *testing.T) {
	type z struct{ X time.Duration }
	in := &z{X: time.Millisecond*250 + time.Hour*8}
	av := encDurationNano(in, 0)
	expectT(t, new(types.AttributeValueMemberN), av)
	out := &z{}
	decDurationNano(out, 0, av)
	if in.X != out.X {
		t.Fatalf("expected %v, got %v", in.X, out.X)
	}
}

func TestDurationPtrNanoRoundTrip(t *testing.T) {
	type z struct{ X *time.Duration }
	orig := time.Millisecond*350 + time.Hour*12
	in := &z{X: &orig}
	av := encDurationNano(in, 0)
	expectT(t, new(types.AttributeValueMemberN), av)
	out := &z{}
	decDurationNano(out, 0, av)
	if out.X == nil {
		t.Fatal("out is nil")
	}
	if in.X == out.X {
		t.Fatal("out is not a new pointer")
	}
	if *in.X != *out.X {
		t.Fatalf("expected %v, got %v", *in.X, *out.X)
	}
}

func TestDurationSecRoundTrip(t *testing.T) {
	type z struct{ X time.Duration }
	in := &z{X: time.Millisecond*250 + time.Hour*8}
	av := encDurationSec(in, 0)
	expectT(t, new(types.AttributeValueMemberN), av)
	out := &z{}
	decDurationSec(out, 0, av)
	if in.X != out.X {
		t.Fatalf("expected %v, got %v", in.X, out.X)
	}
}

func TestDurationPtrSecRoundTrip(t *testing.T) {
	type z struct{ X *time.Duration }
	orig := time.Millisecond*350 + time.Hour*12
	in := &z{X: &orig}
	av := encDurationSec(in, 0)
	expectT(t, new(types.AttributeValueMemberN), av)
	out := &z{}
	decDurationSec(out, 0, av)
	if out.X == nil {
		t.Fatal("out is nil")
	}
	if in.X == out.X {
		t.Fatal("out is not a new pointer")
	}
	if *in.X != *out.X {
		t.Fatalf("expected %v, got %v", *in.X, *out.X)
	}
}
