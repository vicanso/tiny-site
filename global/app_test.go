package global

import "testing"

func TestApplicationToggle(t *testing.T) {
	if IsApplicationRunning() {
		t.Fatalf("the application should not be running first")
	}
	StartApplication()
	if !IsApplicationRunning() {
		t.Fatalf("the application should be running")
	}
	PauseApplication()
	if IsApplicationRunning() {
		t.Fatalf("the application should not be running first")
	}
}
