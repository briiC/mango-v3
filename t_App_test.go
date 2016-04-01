package mango

import "testing"

// Parsing datetimes
func Test_NewApplication(t *testing.T) {
	app, err := NewApplication()
	// app.Print()

	if err != nil {
		t.Fatal(app, err)
	}

	// Check default paths
	if app.binPath == "" {
		t.Fatal("binPath is empty")
	}

	// Trim from binPath end to content path end
	if app.ContentPath[len(app.binPath):] != "/test-files/content" {
		t.Fatal("Default ContentPath must end with /test-files/content")
	}
	// Trim from binPath end to public path end
	if app.PublicPath[len(app.binPath):] != "/test-files/public" {
		t.Fatal("Default PublicPath must end with /test-files/public")
	}

}
