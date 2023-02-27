package object

import "testing"

func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "Hello world"}
	hello2 := &String{Value: "Hello world"}
	diff1 := &String{Value: "My name is adam"}
	diff2 := &String{Value: "My name is adam"}

	if hello1.HashKey() != hello2.HashKey() ||
		diff1.HashKey() != diff2.HashKey() {
		t.Fatalf("strings with same content have different hash keys")
	}

	if hello1.HashKey() == diff1.HashKey() {
		t.Fatalf("strings with different content have same hash keys")
	}
}
