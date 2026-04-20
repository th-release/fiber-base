package xss

import "testing"

func TestSanitizeText(t *testing.T) {
	input := ` <script>alert('x')</script><b>Hello</b> `
	sanitized := SanitizeText(input)

	if sanitized != "Hello" {
		t.Fatalf("unexpected sanitized output: %q", sanitized)
	}
}

func TestContainsRisk(t *testing.T) {
	if !ContainsRisk(`<img src=x onerror=alert(1)>`) {
		t.Fatal("expected xss risk to be detected")
	}
	if ContainsRisk("safe text only") {
		t.Fatal("did not expect safe text to be flagged")
	}
}
