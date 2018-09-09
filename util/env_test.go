package util

import "testing"

func TestEnv(t *testing.T) {
	if IsDevelopment() {
		t.Fatalf("it should be test env")
	}
	if !IsTest() {
		t.Fatalf("it should be test env")
	}
	if IsProduction() {
		t.Fatalf("is should be test env")
	}
}
