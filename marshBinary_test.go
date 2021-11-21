package ddbstruct

import (
	"net/url"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestTimeBinaryRoundTrip(t *testing.T) {
	type z struct{ X time.Time }
	in := &z{X: time.Now()}
	av := encBinary(in, 0)
	expectT(t, new(types.AttributeValueMemberB), av)
	out := &z{}
	decBinary(out, 0, av)
	if in.X.UnixNano() != out.X.UnixNano() {
		t.Fatalf("expected %v, got %v", in.X, out.X)
	}
}

func TestTimePtrBinaryRoundTrip(t *testing.T) {
	type z struct{ X *time.Time }
	orig := time.Now()
	in := &z{X: &orig}
	av := encBinary(in, 0)
	expectT(t, new(types.AttributeValueMemberB), av)
	out := &z{}
	decBinary(out, 0, av)
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

func TestURLPtrBinaryRoundTrip(t *testing.T) {
	type z struct{ X *url.URL }
	orig, err := url.Parse("https://test.example/path/resource?query")
	if err != nil {
		t.Fatalf("failed to parse test URL: %v", err)
	}
	in := &z{X: orig}
	av := encBinary(in, 0)
	expectT(t, new(types.AttributeValueMemberB), av)
	out := &z{}
	decBinary(out, 0, av)
	if out.X == nil {
		t.Fatal("out is nil")
	}
	if in.X == out.X {
		t.Fatal("out is not a new pointer")
	}
	if *in.X != *out.X {
		t.Fatalf("expected %s, got %s", in.X, out.X)
	}
}
