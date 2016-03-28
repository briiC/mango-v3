package mango

import "testing"

// Parsing datetimes
func Test_NewApplication(t *testing.T) {
	app, err := NewApplication()
	// fmt.Println(app)
	if err != nil {
		t.Fatal(app, err)
	}

	// Check default paths
	if app.BinPath == "" {
		t.Fatal("BinPath is empty")
	}

	// Must not be busy by default
	if app.isBusy {
		t.Fatal("App MUST NOT be busy right now")
	}

	// Trim from BinPath end to content path end
	if app.ContentPath[len(app.BinPath):] != "/test-files/content" {
		t.Fatal("Default ContentPath must end with /test-files/content")
	}
	// Trim from BinPath end to public path end
	if app.PublicPath[len(app.BinPath):] != "/test-files/public" {
		t.Fatal("Default PublicPath must end with /test-files/public")
	}

}
